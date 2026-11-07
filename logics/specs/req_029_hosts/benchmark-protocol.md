# Local model benchmark protocol (req_029, Need 3)

Run on each machine under test. Nothing here was run in the spike; the laptop used for the spike has no Ollama and no models.
Answers one question per task: go or no-go for local use, per machine.

## Machines
Laptop (CPU only), tower, homelab. Record per machine: CPU, RAM, GPU and VRAM, OS, Ollama version, power source.

## Models
Pick sizes that fit memory (quantised 4-bit, weights roughly 0.6 GB per billion parameters plus context):

| Slot | Candidates (pull exact tags at run time and record the digest) |
|---|---|
| Small | one Gemma, one Qwen, one Llama at about 1 to 4 B |
| Mid | the same three families at about 7 to 14 B |
| Code | one code-tuned model at about 7 B |

Skip any model that does not fit with 25 percent memory free. Record the tag and digest of each.

## Tasks and fixed sets
Five tasks, 20 samples each, stored in a private folder (never committed; samples hold private data):
1. Run summary from a `cdx conversation` excerpt (2 to 6k tokens in, 5 lines out).
2. Mail triage label (one of a fixed list, from subject plus first 500 words).
3. Bank line categorisation (merchant string to a fixed category list).
4. Memo to action (3 to 5 bullets with owner and due date or "none").
5. French short prose (80 words, owner prose rules from `local/owner-prefs.md`).

For each sample: a gold answer written by the owner, or the Claude Opus output as baseline for tasks 1, 4, 5. Prompts are fixed files, same for every model, temperature 0, fixed seed where the runtime allows it.

## Metrics
Per run: load time, first token latency (ms), generation tokens per second, prompt tokens per second, peak memory (RSS and VRAM), wall time, and power draw where a meter or `powertop` or `nvidia-smi` gives it. Per task: quality score.

Quality scoring:
- Tasks 2 and 3: exact match against gold, report accuracy.
- Tasks 1, 4, 5: the owner rates each output 0 (unusable), 1 (needs edits), 2 (good) against the baseline, blind to the model. Alternatively Claude Haiku as judge on the non-private tasks only.

## Go or no-go per task
Go when all hold: accuracy at least 90 percent (tasks 2 and 3) or mean rating at least 1.5 with no more than 2 samples at 0 (tasks 1, 4, 5); one sample finishes within the task's latency budget (interactive 30 s, batch 5 min); memory stays under 75 percent of RAM or VRAM. Otherwise no-go for that machine and model.

## Script outline
`bench.py` (stdlib only, calls the Ollama HTTP API on the machine under test):
1. Read `machine.json` (labels) and `models.txt`.
2. For each model: `POST /api/chat` with `stream: true`, `options.temperature: 0`, `keep_alive: 0` between models so each starts cold; then once warm for the sample loop.
3. Per request record from the response: `load_duration`, `prompt_eval_count`, `prompt_eval_duration`, `eval_count`, `eval_duration`, `total_duration`; first token latency from the stream timestamps.
4. Sample memory every second from `/proc/<ollama pid>/status` (VmRSS) and `nvidia-smi --query-gpu=memory.used --format=csv` when present.
5. Write `results/<machine>/<model>/<task>.jsonl` (sample id, output, timings).
6. `score.py` computes accuracy or builds the blind rating sheet, then a table: machine by model by task with tokens per second, first token ms, memory, quality, go or no-go.

Order of runs: private-data tasks (3, 2, 4) first, then 1 and 5. Repeat each machine twice, once on AC and once on battery (laptop), once with the desktop idle and once under a typical load.

## Output
One table: 3 machines by chosen models by 5 tasks with go or no-go (AC2). Store raw results outside the repo; commit only the summary table with no sample text.
