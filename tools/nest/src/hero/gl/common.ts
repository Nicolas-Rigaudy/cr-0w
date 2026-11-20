// Shared WebGL2 helpers for the live head (head.ts) and the table grid (grid.ts):
// program build, render targets, full-screen passes and the separable blur used by both blooms.

export const COPY_FS = `#version 300 es
precision highp float;uniform sampler2D uSrc;uniform vec2 uRes;out vec4 o;
void main(){o=texture(uSrc,gl_FragCoord.xy/uRes);}`

export const BLUR_FS = `#version 300 es
precision highp float;uniform sampler2D uSrc;uniform vec2 uRes,uDir;uniform float uSigma,uStep;out vec4 o;
void main(){vec2 uv=gl_FragCoord.xy/uRes;vec4 s=vec4(0);float ws=0.;
  for(int i=-12;i<=12;i++){float x=float(i)*uStep;float w=exp(-x*x/(2.*uSigma*uSigma));s+=texture(uSrc,uv+uDir*x/uRes)*w;ws+=w;}
  o=s/ws;}`

export type Prog = { p: WebGLProgram; U: Record<string, WebGLUniformLocation | null> }

export function glKit(gl: WebGL2RenderingContext) {
  const sh = (type: number, src: string) => {
    const s = gl.createShader(type)!; gl.shaderSource(s, src); gl.compileShader(s)
    if (!gl.getShaderParameter(s, gl.COMPILE_STATUS)) throw new Error(gl.getShaderInfoLog(s) ?? 'shader')
    return s
  }
  // Uniform arrays are looked up by their bare name (uPA, not uPA[0]).
  const prog = (vs: string, fs: string): Prog => {
    const p = gl.createProgram()!; gl.attachShader(p, sh(gl.VERTEX_SHADER, vs)); gl.attachShader(p, sh(gl.FRAGMENT_SHADER, fs))
    gl.bindAttribLocation(p, 0, 'p'); gl.linkProgram(p)
    if (!gl.getProgramParameter(p, gl.LINK_STATUS)) throw new Error(gl.getProgramInfoLog(p) ?? 'link')
    const U: Prog['U'] = {}, n = gl.getProgramParameter(p, gl.ACTIVE_UNIFORMS)
    for (let i = 0; i < n; i++) { const nm = gl.getActiveUniform(p, i)!.name.replace('[0]', ''); U[nm] = gl.getUniformLocation(p, nm) }
    return { p, U }
  }
  const tex = (w: number, h: number, filter: number = gl.LINEAR) => {
    const t = gl.createTexture()!; gl.bindTexture(gl.TEXTURE_2D, t)
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, null)
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter); gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter)
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE); gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
    return t
  }
  // Framebuffer over one or more colour textures (several = multiple render targets).
  const fbo = (...ts: WebGLTexture[]) => {
    const f = gl.createFramebuffer()!; gl.bindFramebuffer(gl.FRAMEBUFFER, f)
    ts.forEach((t, i) => gl.framebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0 + i, gl.TEXTURE_2D, t, 0))
    if (ts.length > 1) gl.drawBuffers(ts.map((_, i) => gl.COLOR_ATTACHMENT0 + i))
    return f
  }
  const pass = (p: Prog, f: WebGLFramebuffer | null, w: number, h: number) => { gl.bindFramebuffer(gl.FRAMEBUFFER, f); gl.viewport(0, 0, w, h); gl.useProgram(p.p) }
  const bindT = (p: Prog, name: string, unit: number, t: WebGLTexture) => { gl.activeTexture(gl.TEXTURE0 + unit); gl.bindTexture(gl.TEXTURE_2D, t); gl.uniform1i(p.U[name], unit) }
  // One blur direction (dx, dy) from src into dst at w x h; draw() issues the full-screen geometry.
  const blur = (p: Prog, src: WebGLTexture, dst: WebGLFramebuffer, w: number, h: number, dx: number, dy: number, sigma: number, step: number, draw: () => void) => {
    pass(p, dst, w, h); bindT(p, 'uSrc', 0, src)
    gl.uniform2f(p.U.uRes, w, h); gl.uniform2f(p.U.uDir, dx, dy); gl.uniform1f(p.U.uSigma, sigma); gl.uniform1f(p.U.uStep, step)
    draw()
  }
  return { prog, tex, fbo, pass, bindT, blur }
}
