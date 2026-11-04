#!/usr/bin/env python3
"""crow_term: terminal actions behind one backend (open a window, focus or type into a session's window).

Backends, picked by CROW_TERMINAL (kitty, tmux, print) or, when unset, the first that fits:
  kitty  kitty on PATH. Opens windows with `kitty`, focuses and types through remote control sockets.
  tmux   inside a tmux session ($TMUX). Opens tmux windows, focuses and types into panes.
  print  anything else. Prints the command to run in a new terminal; cannot focus or type.

As a script (for shell callers such as bin/crow):
  crow_term.py backend
  crow_term.py open --title T [--cwd D] [--set K=V]... -- CMD...
      runs CMD detached from the calling session (clean_env plus the --set vars); on the print
      backend prints the line to run instead (the --set vars included).
"""
import glob, json, os, pwd, shlex, shutil, subprocess, sys

BACKENDS = ("kitty", "tmux", "print")


class UnknownBackend(SystemExit):
    pass


def choose(env=None):
    """Backend name: CROW_TERMINAL when set, else kitty, tmux, print in that order."""
    env = os.environ if env is None else env
    name = env.get("CROW_TERMINAL", "").strip().lower()
    if name:
        if name not in BACKENDS:
            raise UnknownBackend(f"crow: CROW_TERMINAL={name} is not a backend ({', '.join(BACKENDS)})")
        return name
    if shutil.which("kitty"):
        return "kitty"
    if env.get("TMUX") and shutil.which("tmux"):
        return "tmux"
    return "print"


def clean_env(extra=None):
    """Detach from the calling session: real HOME, no inherited Claude/cdx variables, plus extra."""
    env = {k: v for k, v in os.environ.items() if not k.startswith(("CLAUDE", "ANTHROPIC_", "CDX_"))}
    env["HOME"] = pwd.getpwuid(os.getuid()).pw_dir
    env.update(extra or {})
    return env


def parent(pid, proc="/proc"):
    """Parent pid from /proc, else from `ps` (macOS), else None."""
    try:
        with open(f"{proc}/{pid}/stat") as fh:
            return int(fh.read().rsplit(")", 1)[1].split()[1])
    except (OSError, ValueError, IndexError):
        pass
    if os.path.isdir(proc):
        return None
    try:
        r = subprocess.run(["ps", "-o", "ppid=", "-p", str(pid)], capture_output=True, text=True, timeout=5)
        return int(r.stdout.strip())
    except (OSError, ValueError, subprocess.TimeoutExpired):
        return None


def ancestors(pid, proc="/proc"):
    """pid and its parents, nearest first."""
    out = []
    while pid and pid > 1 and pid not in out:
        out.append(pid)
        pid = parent(pid, proc)
    return out


def env_changes(env):
    """(vars set or changed, vars removed) in env compared with this process."""
    if env is None:
        return {}, []
    added = {k: v for k, v in env.items() if os.environ.get(k) != v}
    return added, sorted(k for k in os.environ if k not in env)


def detach(argv, env, cwd=None):
    """Run argv in its own session, stdio closed; `setsid -f` where it exists (Linux)."""
    quiet = dict(env=env, cwd=cwd, stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    if shutil.which("setsid", path=(env or os.environ).get("PATH")):
        subprocess.run(["setsid", "-f", *argv], check=True, **quiet)
    else:
        subprocess.Popen(argv, start_new_session=True, **quiet)


# ---------- kitty ----------

def kitty_find_window(kitty_ls, pid, chain=None):
    """Id of the kitty window holding pid (its own process, a foreground process or an ancestor), else None."""
    chain = chain if chain is not None else ancestors(pid)
    for osw in kitty_ls:
        for tab in osw.get("tabs", []):
            for win in tab.get("windows", []):
                pids = {win.get("pid")} | {p.get("pid") for p in win.get("foreground_processes", [])}
                if pids & set(chain):
                    return win["id"]
    return None


def kitty_send_args(sock, wid, text):
    """kitten command that types text and Enter into window wid."""
    return ["kitten", "@", "--to", sock, "send-text", "--match", f"id:{wid}", "--", text + "\r"]


def kitty_sockets():
    runtime = os.environ.get("XDG_RUNTIME_DIR") or f"/run/user/{os.getuid()}"
    return ["unix:" + s for s in sorted(glob.glob(f"{runtime}/kitty-*"))]


def kitty_locate(pid):
    """(socket, window id) of the window holding pid, else None."""
    for sock in kitty_sockets():
        try:
            r = subprocess.run(["kitten", "@", "--to", sock, "ls"], capture_output=True, text=True, timeout=5)
            wid = kitty_find_window(json.loads(r.stdout) if r.returncode == 0 else [], pid)
        except (OSError, ValueError, subprocess.TimeoutExpired):
            continue
        if wid is not None:
            return sock, wid
    return None


def kitty_open_args(argv, title, cwd=None):
    return ["kitty", "--title", title, *(["--directory", cwd] if cwd else []), *argv]


# ---------- tmux ----------

def tmux_find_pane(panes, chain):
    """Pane id from `list-panes -F '#{pane_id} #{pane_pid}'` lines whose shell is in chain, else None."""
    for line in panes.splitlines():
        parts = line.split()
        if len(parts) == 2 and parts[1].isdigit() and int(parts[1]) in chain:
            return parts[0]
    return None


def tmux_locate(pid):
    try:
        r = subprocess.run(["tmux", "list-panes", "-a", "-F", "#{pane_id} #{pane_pid}"],
                           capture_output=True, text=True, timeout=5)
    except (OSError, subprocess.TimeoutExpired):
        return None
    return tmux_find_pane(r.stdout, ancestors(pid)) if r.returncode == 0 else None


def tmux_open_args(argv, title, cwd=None, env=None):
    """tmux new-window running argv; env changes go through `env` (the tmux server has its own env)."""
    added, removed = env_changes(env)
    added.pop("PATH", None)
    wrap = ["env", *(a for k in removed for a in ("-u", k)), *(f"{k}={v}" for k, v in sorted(added.items()))]
    return ["tmux", "new-window", "-n", title, *(["-c", cwd] if cwd else []),
            *(wrap + argv if len(wrap) > 1 else argv)]


def tmux_send_args(pane, text):
    return [["tmux", "send-keys", "-t", pane, "-l", text], ["tmux", "send-keys", "-t", pane, "Enter"]]


# ---------- print ----------

def print_line(argv, cwd=None, env=None, show_env=()):
    """Shell line to run argv by hand: `cd DIR && K=V cmd...` (only the show_env vars)."""
    shown = [f"{k}={shlex.quote(env[k])}" for k in show_env if env and env.get(k)]
    return " ".join([*([f"cd {shlex.quote(cwd)} &&"] if cwd else []), *shown, shlex.join(argv)])


# ---------- the actions ----------

def open_window(argv, title, cwd=None, env=None, show_env=(), backend=None):
    """Run argv in a new terminal window. Returns 'opened' or, on the print backend, the line to run."""
    backend = backend or choose()
    if backend == "kitty":
        detach(kitty_open_args(argv, title, cwd), env, cwd)
        return "opened"
    if backend == "tmux":
        subprocess.run(tmux_open_args(argv, title, cwd, env), check=True, env=env,
                       stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        return "opened"
    return print_line(argv, cwd, env, show_env)


def focus(pid, backend=None):
    """Bring the window holding pid to the front. A short description of it, else None."""
    backend = backend or choose()
    if backend == "kitty" and (hit := kitty_locate(pid)):
        subprocess.run(["kitten", "@", "--to", hit[0], "focus-window", "--match", f"id:{hit[1]}"], check=True)
        return f"window {hit[1]}"
    if backend == "tmux" and (pane := tmux_locate(pid)):
        for cmd in (["select-window", "-t", pane], ["select-pane", "-t", pane], ["switch-client", "-t", pane]):
            subprocess.run(["tmux", *cmd], capture_output=True, timeout=5)
        return f"pane {pane}"
    return None


def send(pid, text, backend=None):
    """Type text and Enter into the window holding pid. True when a window took it."""
    backend = backend or choose()
    if backend == "kitty" and (hit := kitty_locate(pid)):
        subprocess.run(kitty_send_args(*hit, text), check=True, timeout=5)
        return True
    if backend == "tmux" and (pane := tmux_locate(pid)):
        for cmd in tmux_send_args(pane, text):
            subprocess.run(cmd, check=True, timeout=5)
        return True
    return False


def main(argv):
    try:
        if argv[:1] == ["backend"]:
            print(choose())
            return 0
        if argv[:1] != ["open"] or "--" not in argv:
            sys.exit(__doc__)
        i = argv.index("--")
        opts, cmd = argv[1:i], argv[i + 1:]
        title, cwd, extra = "crow", None, {}
        while opts:
            key, val, opts = opts[0], opts[1] if len(opts) > 1 else sys.exit(__doc__), opts[2:]
            if key == "--title":
                title = val
            elif key == "--cwd":
                cwd = val
            elif key == "--set" and "=" in val:
                k, v = val.split("=", 1)
                extra[k] = v
            else:
                sys.exit(__doc__)
        if not cmd:
            sys.exit(__doc__)
        out = open_window(cmd, title, cwd, clean_env(extra), list(extra))
    except (OSError, subprocess.CalledProcessError) as e:
        sys.exit(f"crow: could not open a window ({e})")
    if out != "opened":
        print(out)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
