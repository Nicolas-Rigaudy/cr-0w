// Live WebGL2 head (item_156, motion items 196 + 198). Port of the approved option B spike
// (vault item126/spike/holo-head.html, holo=0) drawn on its own canvas under the table canvas.
// The cell grid matches the baked heads (216 cols at 4x7 virtual px, focal 993), so the live head
// and the static fallback sit at the same place and size. Margins around the bake grid leave room
// for the fin, the bloom and the arms.

type RGB = [number, number, number]
interface V3 { x: number; y: number; z: number }
export const V = (x: number, y: number, z: number): V3 => ({ x, y, z })
const add = (a: V3, b: V3) => V(a.x + b.x, a.y + b.y, a.z + b.z)
const sub = (a: V3, b: V3) => V(a.x - b.x, a.y - b.y, a.z - b.z)
const mul = (a: V3, s: number) => V(a.x * s, a.y * s, a.z * s)
const dot = (a: V3, b: V3) => a.x * b.x + a.y * b.y + a.z * b.z
const cross = (a: V3, b: V3) => V(a.y * b.z - a.z * b.y, a.z * b.x - a.x * b.z, a.x * b.y - a.y * b.x)
const len = (a: V3) => Math.sqrt(dot(a, a))
const norm = (a: V3) => mul(a, 1 / len(a))
const lerpV = (a: V3, b: V3, t: number) => add(a, mul(sub(b, a), t))
const rotY = (p: V3, a: number) => { const c = Math.cos(a), s = Math.sin(a); return V(p.x * c - p.z * s, p.y, p.x * s + p.z * c) }
const rotX = (p: V3, a: number) => { const c = Math.cos(a), s = Math.sin(a); return V(p.x, p.y * c - p.z * s, p.y * s + p.z * c) }
const D2R = Math.PI / 180

// Bake grid (scripts/hero-heads.mjs: w=863 h=800 focal=993 dy=160).
export const BAKE = { cols: 216, rows: 87, focal: 993, cx: 863 / 2, cy: 133 + 160, cw: 4, ch: 7 }
// Margins in cells: room above for the fin and the bloom, at the sides and below for the arms.
export const MARGIN = { x: 10, top: 8, bottom: 8 }
const COLS = BAKE.cols + 2 * MARGIN.x, ROWS = BAKE.rows + MARGIN.top + MARGIN.bottom
const RAMP = '.,-~:;=!*#$@·'
const CAM = { pos: V(0, 0.7, -6.6), tgt: V(0, -0.2, 0) }

// ---------- arms + fin (spike armPrims / FIN) ----------
export const ST = { mount: (sx: number) => norm(V(sx * 0.72, -0.52, -0.46)), L1: 0.7, L2: 0.66, r1: 0.17, r2: 0.145, rj: 0.2, fn: 2, fl: [0.33, 0.27], fr: 0.06 }
export interface PoseP { w: V3; pole: V3; open: number; hand: V3; hold?: { r?: number; octa?: number } }
export type Pose = 'rest' | 'grab' | 'present' | 'brace' | 'droop' | 'wave' | 'think'
// Arm poses per side sx. One-arm poses use the active arm a (+1 right, -1 left); the other arm rests.
// Rest: claws hang low and angle inwards without touching (owner pick M4 inward).
const one = (sx: number, a: number, p: PoseP) => sx === a ? p : POSE.rest(sx)
export const POSE: Record<Pose, (sx: number, a?: number) => PoseP> = {
  rest: (sx) => ({ w: V(sx * 0.8, -1.64, -0.56), pole: V(sx * 0.7, -0.3, 0.6), open: 0.75, hand: V(-sx * 0.35, -1, -0.45) }),
  // Work the holo table: both claws hover over it, fingers down (HeadCtl taps them).
  grab: (sx) => ({ w: V(sx * 1.1, -1.2, -1.45), pole: V(sx * 0.9, 0, 0.4), open: 0.6, hand: V(-sx * 0.7, -0.65, -0.45) }),
  // Present: open claw sweeps out toward the selected planet's side.
  present: (sx, a = 1) => one(sx, a, { w: V(sx * 1.5, -1.05, -1.1), pole: V(sx * 0.4, -1, 0.3), open: 1.05, hand: V(sx * 0.85, -0.35, -0.6) }),
  brace: (sx) => ({ w: V(sx * 1.85, -0.15, -0.55), pole: V(sx * 0.2, -1, 0.4), open: 1.0, hand: V(sx * 0.6, 0.3, -1) }),
  // Droop after a fail: arms hang slack, claws nearly closed.
  droop: (sx) => ({ w: V(sx * 0.62, -1.78, -0.32), pole: V(sx * 0.6, -0.1, 0.7), open: 0.25, hand: V(-sx * 0.1, -1, 0.15) }),
  wave: (sx, a = 1) => one(sx, a, { w: V(sx * 1.5, 0.05, -0.75), pole: V(sx * 0.6, -1, 0.3), open: 0.8, hand: V(sx * 0.3, 1, -0.3) }),
  // Thinking: claw to the chin.
  think: (sx, a = 1) => one(sx, a, { w: V(sx * 0.25, -1.2, -1.35), pole: V(sx * 1, -1, -0.2), open: 0.35, hand: V(-sx * 0.3, 0.9, -0.3) }),
}
// Pose blend: from -> to at k (0..1). Held objects scale with their pose's weight.
function blendPose(a: PoseP, b: PoseP, k: number): PoseP {
  const hold = { r: (a.hold?.r ?? 0) * (1 - k) + (b.hold?.r ?? 0) * k, octa: (a.hold?.octa ?? 0) * (1 - k) + (b.hold?.octa ?? 0) * k }
  return { w: lerpV(a.w, b.w, k), pole: lerpV(a.pole, b.pole, k), open: a.open + (b.open - a.open) * k, hand: lerpV(a.hand, b.hand, k), hold: hold.r > 0.005 || hold.octa > 0.005 ? hold : undefined }
}
function ik(S: V3, Wt: V3, L1: number, L2: number, pole: V3) {
  let ax = sub(Wt, S); let d = len(ax); ax = mul(ax, 1 / d); d = Math.min(d, L1 + L2 - 1e-3)
  const a = (L1 * L1 - L2 * L2 + d * d) / (2 * d), h = Math.sqrt(Math.max(0, L1 * L1 - a * a))
  const pp = norm(sub(pole, mul(ax, dot(pole, ax))))
  return { E: add(add(S, mul(ax, a)), mul(pp, h)), W: add(S, mul(ax, d)) }
}
interface Prim { t: 's' | 'c'; a: V3; b?: V3; r: number; m: number }
function armPrims(sx: number, ps: PoseP) {
  const M: Prim[] = [], mount = ST.mount(sx), S = mul(mount, 1.16), { E, W: Wr } = ik(S, ps.w, ST.L1, ST.L2, ps.pole)
  const hd = norm(ps.hand)
  const s = (a: V3, r: number, m: number): Prim => ({ t: 's', a, r, m }), c = (a: V3, b: V3, r: number, m: number): Prim => ({ t: 'c', a, b, r, m })
  M.push(s(mul(mount, 1.2), ST.rj * 1.45, 2), c(S, E, ST.r1, 1), s(E, ST.rj, 2), c(E, Wr, ST.r2, 1),
    c(add(Wr, mul(norm(sub(E, Wr)), 0.15)), Wr, ST.r2 * 1.4, 2), s(Wr, ST.r2 * 1.15, 2))
  const B = add(Wr, mul(hd, 0.06))
  let u = cross(hd, V(0, 1, 0)); if (len(u) < 0.2) u = cross(hd, V(1, 0, 0)); u = norm(u)
  const v = cross(hd, u), a1 = 0.2 + ps.open * 0.8, a2 = a1 - 1.15
  for (let k = 0; k < ST.fn; k++) {
    const ph = Math.PI / 2 + k * Math.PI, rk = add(mul(u, Math.cos(ph)), mul(v, Math.sin(ph)))
    const d1 = add(mul(hd, Math.cos(a1)), mul(rk, Math.sin(a1))), d2 = add(mul(hd, Math.cos(a2)), mul(rk, Math.sin(a2)))
    const K = add(B, mul(d1, ST.fl[0])), T = add(K, mul(d2, ST.fl[1]))
    M.push(c(B, K, ST.fr, 1), c(K, T, ST.fr * 0.8, 1), s(T, ST.fr * 0.85, 3))
  }
  const tipC = add(Wr, mul(hd, 0.08 + ST.fl[0] * 0.85 + (ps.hold?.r || 0.1) * 0.3))
  const hold = ps.hold ? { ...ps.hold, c: tipC } : null
  return { M, hold, tipC }
}
const FIN: Prim[] = [
  { t: 'c', a: V(0, 1.2, -0.2), b: V(0, 1.68, 0.42), r: 0.06, m: 11 }, { t: 'c', a: V(0, 1.17, 0.05), b: V(0, 1.52, 0.5), r: 0.06, m: 11 },
  { t: 'c', a: V(0, 1.12, 0.35), b: V(0, 1.36, 0.62), r: 0.06, m: 11 }, { t: 'c', a: V(0, 1.2, -0.2), b: V(0, 1.06, 0.7), r: 0.07, m: 12 },
  { t: 's', a: V(0, 1.68, 0.42), r: 0.08, m: 13 },
]

// ---------- shaders ----------
const VS = `#version 300 es
in vec2 p;void main(){gl_Position=vec4(p,0,1);}`

const CELL_FS = `#version 300 es
precision highp float;precision highp int;
uniform vec3 uPos,uF,uR,uU;uniform float uFocal;uniform vec2 uC;uniform float uRows;
uniform float uYaw,uPitch,uTilt,uBlink,uMouth,uEyeS;uniform int uMood;uniform vec3 uTip,uHoldCol,uEye;uniform float uEyeK;uniform vec2 uGaze;
uniform vec4 uPA[40],uPB[40];uniform int uNP;uniform vec4 uBound;uniform vec4 uOcta;uniform vec4 uHoldS;
uniform float uScan,uShimmer,uFrame;
layout(location=0) out vec4 oCol;layout(location=1) out vec4 oAux;
#define C3(r,g,b) (vec3(r,g,b)/255.)
const vec3 LIGHT=vec3(-0.55,0.75,-0.6)/1.1068;
const vec3 PUR=C3(192.,92.,255.),WHT=C3(255.,236.,255.),CY=C3(46.,242.,255.),FR=C3(79.,195.,255.);
const float PI=3.14159265;
float hash(vec2 p){return fract(sin(dot(p,vec2(127.1,311.7)))*43758.5453);}
float cov(float d){return clamp(0.55-d*0.95,0.,1.);}
float sdRbox(vec2 p,vec2 c,vec2 h,float r){r=min(r,min(h.x,h.y));vec2 q=abs(p-c)-(h-r);return length(max(q,0.))+min(max(q.x,q.y),0.)-r;}
float sdSeg(vec2 p,vec2 a,vec2 b,float th){vec2 pa=p-a,ba=b-a;float h=clamp(dot(pa,ba)/dot(ba,ba),0.,1.);return length(pa-ba*h)-th;}
float halfCov(vec2 p,vec2 n,float c){return cov(c-dot(p,n));}
int ramp(float L){return int(clamp(floor(L*12.),0.,11.));}
float eyeT(vec2 p,float ex,float ey,float hw,float hh,float lid,float tilt){
  float s=ex<0.?-1.:1.,ch=min(hh*1.1,2.6);
  float a=cov(sdRbox(p,vec2(ex,ey),vec2(hw,hh),0.6));if(a<=0.)return 0.;
  vec2 nb=vec2(s*0.6,0.8),pb=vec2(ex+s*(hw-ch),ey+hh),nt=vec2(-s*0.6,-0.8),pt=vec2(ex-s*(hw-ch),ey-hh);
  float c=max(halfCov(p,nb,dot(nb,pb)),halfCov(p,nt,dot(nt,pt)));
  if(lid>0.){float t=radians(tilt);vec2 n=vec2(sin(t),-cos(t));float y0=ey-hh+2.*hh*lid;c=max(c,halfCov(p,n,n.x*ex+n.y*y0));}
  return a*(1.-c);
}
// Mouth: five grille squares on a curve (uMouth > 0 smiles, < 0 frowns).
float grille(vec2 p){float a=0.;for(int i=0;i<5;i++){float x=-4.+2.*float(i);a=max(a,cov(sdRbox(p,vec2(x,8.4+uMouth*(1.-x*x/16.)),vec2(0.55),0.2)));}return a;}
float ring(vec2 p,vec2 c,float r,float th){return cov(abs(length(p-c)-r)-th);}
// Eyes by style (FACE_STYLE): 0 ovals, 1 calm lids, 2 smile arcs, 3 narrowed, 4 > <, 5 worried, 6 x x, 7 slits.
// Gaze shifts them, a blink squeezes them to the lid line.
float eyes(vec2 p){
  const float ey=-1.5;
  p-=uGaze;p.y=ey+(p.y-ey)/uEyeS;
  float lid=ey+0.8,sq=max(1.-uBlink*0.85,0.15);p.y=lid+(p.y-lid)/sq;
  if(uMood==1)return max(eyeT(p,-9.,ey,5.,2.6,0.32,0.),eyeT(p,9.,ey,5.,2.6,0.32,0.));
  if(uMood==2){float l=p.y<ey+1.?ring(p,vec2(-9.,ey+1.6),4.4,0.75):0.,r=p.y<ey+1.?ring(p,vec2(9.,ey+1.6),4.4,0.75):0.;return max(l,r);}
  if(uMood==3)return max(eyeT(p,-9.,ey+0.4,5.2,1.6,0.15,9.),eyeT(p,9.,ey+0.4,5.2,1.6,0.15,-9.));
  if(uMood==4){float l=min(sdSeg(p,vec2(-13.,ey-3.4),vec2(-7.,ey),0.7),sdSeg(p,vec2(-7.,ey),vec2(-13.,ey+3.4),0.7));
    float r=min(sdSeg(p,vec2(13.,ey-3.4),vec2(7.,ey),0.7),sdSeg(p,vec2(7.,ey),vec2(13.,ey+3.4),0.7));return max(cov(l),cov(r));}
  if(uMood==5)return max(eyeT(p,-9.,ey,5.,2.6,0.38,-14.),eyeT(p,9.,ey,5.,2.6,0.38,14.));
  if(uMood==6){float l=min(sdSeg(p,vec2(-12.,ey-3.),vec2(-6.,ey+3.),0.7),sdSeg(p,vec2(-12.,ey+3.),vec2(-6.,ey-3.),0.7));
    float r=min(sdSeg(p,vec2(6.,ey-3.),vec2(12.,ey+3.),0.7),sdSeg(p,vec2(6.,ey+3.),vec2(12.,ey-3.),0.7));return max(cov(l),cov(r));}
  if(uMood==7)return max(eyeT(p,-9.,ey+0.5,6.,1.1,0.,0.),eyeT(p,9.,ey+0.5,6.,1.1,0.,0.));
  return max(eyeT(p,-9.,ey,5.,2.6,0.,0.),eyeT(p,9.,ey,5.,2.6,0.,0.));
}
vec3 rotY(vec3 p,float a){float c=cos(a),s=sin(a);return vec3(p.x*c-p.z*s,p.y,p.x*s+p.z*c);}
vec3 rotX(vec3 p,float a){float c=cos(a),s=sin(a);return vec3(p.x,p.y*c-p.z*s,p.y*s+p.z*c);}
void groove(float lum,out int g,out vec3 c,out float a){g=0;c=C3(44.,32.,96.);a=0.28+0.2*lum;}
void hull(float lum,float under,float k,out int g,out vec3 c,out float a){float L=clamp(lum*0.72+under*0.35+0.03+k,0.,0.76);g=ramp(L);c=mix(mix(C3(52.,36.,112.),C3(220.,196.,255.),L),CY,under*0.7);a=0.5+0.5*L;}
void headShade(vec3 n,vec3 d,out int g,out vec3 c,out float a){
  vec3 loc=rotX(rotY(n,-uYaw),-uPitch);
  float lum=max(0.,dot(n,LIGHT)),rim=pow(1.-max(0.,-dot(n,d)),2.2),under=clamp(-n.y*0.8+0.25,0.,1.)*rim;
  float front=-loc.z,lon=atan(loc.x,-loc.z);
  if(abs(loc.x)>0.917){float r=length(loc.yz);
    if(r<=0.4){
      if(r>0.33){g=11;c=mix(FR,WHT,0.2+lum*0.45);a=0.7+0.3*lum;return;}
      if(r>0.29){groove(lum,g,c,a);return;}
      if(r>0.13){float an=atan(loc.y,loc.z);if(mod(floor((an+PI)/(PI/6.)),2.)>0.5){g=4;c=C3(60.,46.,130.);a=0.5;}else{g=6;c=mix(FR,CY,0.35);a=0.5+0.35*lum;}return;}
      g=9;c=mix(CY,WHT,0.35);a=0.95;return;}}
  if(front>0.1){
    float u=loc.x*23.,v=-loc.y*23.,t=-radians(uTilt);
    vec2 cuv=vec2(u*cos(t)-v*sin(t),u*sin(t)+v*cos(t));
    float vs=sdRbox(vec2(u,v),vec2(0.,0.6),vec2(18.5,12.6),6.);
    if(vs<0.){
      float ke=eyes(cuv),kg=grille(cuv);
      if(ke>0.12&&ke>=kg){g=int(clamp(7.+floor(ke*4.4+0.5),7.,11.));c=mix(mix(PUR,WHT,ke*ke*0.75),uEye,0.9*uEyeK);a=0.55+0.45*ke;return;}
      if(kg>0.12){g=int(clamp(7.+floor(kg*4.4+0.5),7.,11.));c=mix(C3(120.,96.,200.),WHT,kg*0.35);a=0.55+0.45*kg;return;}
      float sh=pow(max(0.,dot(n,normalize(LIGHT+vec3(0,0,-1)))),24.);
      if(sh>0.25){g=2;c=mix(FR,WHT,0.5);a=0.35*sh+0.1;return;}
      float gl=clamp(lum*0.32+0.1,0.,0.34),s=floor(u+0.5)+floor(v+0.5);
      g=(s-3.*trunc(s/3.))==0.?12:ramp(gl);c=mix(C3(70.,50.,150.),FR,0.25);a=0.34+0.3*lum;return;}
    if(vs<0.9){g=6;c=mix(FR,CY,0.5);a=0.75;return;}
  }
  if(abs(lon)<0.055&&loc.y>0.3){g=9;c=mix(FR,WHT,0.25+lum*0.45);a=0.6+0.4*lum;return;}
  if(abs(lon)<0.1&&loc.y>0.3){groove(lum,g,c,a);return;}
  if(loc.y<-0.62&&loc.y>-0.92&&abs(lon)<0.75){float row=floor((-0.62-loc.y)/0.075);if(mod(row,2.)>0.5)groove(lum,g,c,a);else hull(lum,under,0.05,g,c,a);return;}
  if(abs(loc.y-0.62)<0.03){groove(lum,g,c,a);return;}
  hull(lum,under,0.,g,c,a);
}
float sdCap(vec3 p,vec3 a,vec3 b,float r){vec3 pa=p-a,ba=b-a;float h=clamp(dot(pa,ba)/dot(ba,ba),0.,1.);return length(pa-ba*h)-r;}
vec2 scene(vec3 p){float d=1e9,m=0.;
  for(int i=0;i<40;i++){if(i>=uNP)break;vec4 A=uPA[i],B=uPB[i];float e=B.w<0.?length(p-A.xyz)-A.w:sdCap(p,A.xyz,B.xyz,A.w);float mm=abs(B.w);if(e<d){d=e;m=mm;}}
  if(uOcta.w>0.){vec3 q=abs(p-uOcta.xyz);q.y/=1.35;float e=(q.x+q.y+q.z-uOcta.w)*0.5;if(e<d){d=e;m=4.;}}
  if(uHoldS.w>0.){float e=length(p-uHoldS.xyz)-uHoldS.w;if(e<d){d=e;m=5.;}}
  return vec2(d,m);}
vec3 sceneN(vec3 p){const vec2 k=vec2(1,-1)*0.003;
  return normalize(k.xyy*scene(p+k.xyy).x+k.yyx*scene(p+k.yyx).x+k.yxy*scene(p+k.yxy).x+k.xxx*scene(p+k.xxx).x);}
void armShade(vec3 n,vec3 d,float m,out int g,out vec3 c,out float a){
  bool fin=m>9.;if(fin)m-=10.;
  float lum=max(0.,dot(n,LIGHT)),rim=pow(1.-max(0.,-dot(n,d)),2.2),under=clamp(-n.y*0.8+0.25,0.,1.)*rim;
  if(m>4.5){float L=clamp(lum*0.4+0.45+0.08,0.,1.);g=ramp(L);c=mix(mix(uHoldCol,C3(20.,20.,50.),0.5),mix(uHoldCol,WHT,0.24),L);a=0.4+0.6*L;}
  else if(m>3.5){float L=clamp(lum*1.3+0.25,0.,1.);g=ramp(L);c=mix(mix(uHoldCol,C3(20.,20.,50.),0.55),uHoldCol,L);a=0.45+0.55*L;}
  else if(m>2.5){g=4;c=mix(uTip,WHT,0.25);a=0.85;}
  else if(m>1.5){float L=clamp(lum*0.5+under*0.3+rim*0.15+0.14,0.,0.55);g=ramp(L);c=mix(mix(C3(34.,26.,80.),C3(150.,140.,220.),L),CY,under*0.6);a=0.45+0.5*L;}
  else{float L=clamp(lum*0.75+under*0.35+rim*0.2+0.16,0.,0.82);g=ramp(L);c=mix(mix(C3(52.,36.,112.),C3(220.,196.,255.),L),CY,under*0.7);a=0.5+0.5*L;}
  if(fin){c=mix(c,WHT,0.15);a*=1.4;}
}
void main(){
  float col=floor(gl_FragCoord.x),row=uRows-1.-floor(gl_FragCoord.y);
  vec2 s=vec2(col*4.+2.,row*7.+3.5);
  vec3 d=normalize(uF+uR*((s.x-uC.x)/uFocal)+uU*(-(s.y-uC.y)/uFocal));
  float best=1e9;int g=-1;vec3 c=vec3(0);float a=0.;bool onHead=false;
  vec3 oc=uPos;float b=dot(oc,d),h=b*b-(dot(oc,oc)-1.4884);
  if(h>0.){float t=-b-sqrt(h);if(t>0.){best=t;vec3 n=normalize(uPos+d*t);headShade(n,d,g,c,a);onHead=true;}}
  vec3 bo=uPos-uBound.xyz;float bb=dot(bo,d),bh=bb*bb-(dot(bo,bo)-uBound.w*uBound.w);
  if(bh>0.&&uNP>0){float t=max(0.,-bb-sqrt(bh)),tb=min(-bb+sqrt(bh),best);
    for(int i=0;i<96;i++){vec3 p=uPos+d*t;vec2 e=scene(p);if(e.x<0.0025){vec3 n=sceneN(p);armShade(n,d,e.y,g,c,a);best=t;onHead=false;break;}t+=e.x*0.9;if(t>tb)break;}}
  if(g<0){oCol=vec4(0);oAux=vec4(0);return;}
  // Scan band (A2): one ramp step denser and brighter as it rolls down the head.
  if(uScan>-50.){float k=exp(-pow((s.y-uScan)/14.,2.));if(k>0.3&&g<11){g+=1;}a*=1.+0.35*k;c=mix(c,CY,0.25*k);}
  // Shimmer (A4): a few lit cells move one ramp step per frame.
  if(uShimmer>0.&&g<12&&g>0&&hash(vec2(col+uFrame*7.,row-uFrame*3.))<uShimmer)g+=hash(vec2(row,uFrame))<0.5?-1:1;
  oCol=vec4(c,clamp(a,0.,1.));oAux=vec4(float(g)/255.,0,0,1);
}`

// Glyph pass: float cells in device px, one integer atlas slot per glyph, sampled linearly.
const GLYPH_FS = `#version 300 es
precision highp float;
uniform sampler2D uCol,uAux,uAtlas;uniform vec2 uCell,uSlot,uAtlasSize;uniform vec2 uGrid;uniform float uH;
uniform float uTime,uGlitch,uDpr,uAlphaK;
out vec4 o;
float hash(vec2 p){return fract(sin(dot(p,vec2(127.1,311.7)))*43758.5453);}
void main(){
  vec2 p=vec2(gl_FragCoord.x,uH-gl_FragCoord.y);
  if(uGlitch>0.){float sl=floor(p.y/(9.*uDpr)),f=floor(uTime*24.);if(hash(vec2(sl,f))>0.72)p.x-=(hash(vec2(sl,f+3.))-0.5)*44.*uDpr*uGlitch;}
  float cx=floor(p.x/uCell.x);o=vec4(0);
  if(cx<0.||cx>=uGrid.x)return;
  for(int k=-1;k<=1;k++){
    float r=floor(p.y/uCell.y)+float(k);if(r<0.||r>=uGrid.y)continue;
    ivec2 ci=ivec2(cx,uGrid.y-1.-r);vec4 cc=texelFetch(uCol,ci,0);if(cc.a<=0.)continue;
    float g=floor(texelFetch(uAux,ci,0).r*255.+0.5);
    vec2 l=vec2(p.x-(cx+0.5)*uCell.x,p.y-(r+0.5)*uCell.y)+uSlot*0.5;
    if(l.x<0.||l.y<0.||l.x>=uSlot.x||l.y>=uSlot.y)continue;
    float cov=texture(uAtlas,(vec2(g*uSlot.x,0.)+l)/uAtlasSize).a;
    float al=min(1.,cc.a*uAlphaK)*cov;o.rgb+=cc.rgb*al*(1.-o.a);o.a+=al*(1.-o.a);
  }
}`

const COPY_FS = `#version 300 es
precision highp float;uniform sampler2D uSrc;uniform vec2 uRes;out vec4 o;
void main(){o=texture(uSrc,gl_FragCoord.xy/uRes);}`

const BLUR_FS = `#version 300 es
precision highp float;uniform sampler2D uSrc;uniform vec2 uRes,uDir;uniform float uSigma,uStep;out vec4 o;
void main(){vec2 uv=gl_FragCoord.xy/uRes;vec4 s=vec4(0);float ws=0.;
  for(int i=-12;i<=12;i++){float x=float(i)*uStep;float w=exp(-x*x/(2.*uSigma*uSigma));s+=texture(uSrc,uv+uDir*x/uRes)*w;ws+=w;}
  o=s/ws;}`

// Composite: sharp layer + the same two blooms as the table canvas (3 px * 0.6, 14 px * 0.55).
// Output is premultiplied; bloom adds light where the sharp layer is empty.
const COMP_FS = `#version 300 es
precision highp float;
uniform sampler2D uL,uBs,uBw;uniform vec2 uRes;uniform float uGlitch,uDpr,uFlick;
out vec4 o;
const vec3 CY=vec3(46.,242.,255.)/255.,RED=vec3(255.,46.,99.)/255.;
void main(){
  vec2 uv=gl_FragCoord.xy/uRes;
  vec4 L=texture(uL,uv),Bs=texture(uBs,uv)*0.6,Bw=texture(uBw,uv)*0.55;
  vec3 col=L.rgb+Bs.rgb+Bw.rgb;float a=max(L.a,max(Bs.a,Bw.a));
  if(uGlitch>0.){float s=5.*uDpr/uRes.x;
    vec3 sp=0.55*uGlitch*(texture(uL,uv+vec2(s,0)).a*RED*0.85+texture(uL,uv-vec2(s,0)).a*CY*0.75);col+=sp;a=max(a,max(sp.r,max(sp.g,sp.b)));}
  col*=uFlick;a*=uFlick;
  o=vec4(min(col,vec3(1.)),clamp(a,0.,1.));
}`

// ---------- public API ----------
export type Face = 'idle' | 'fail' | 'run'
// Eye style (shader eyes()) and mouth curve per work face.
export const FACE_STYLE: Record<Face, { eyes: number; mouth: number }> = { idle: { eyes: 0, mouth: 0.7 }, fail: { eyes: 6, mouth: -0.6 }, run: { eyes: 3, mouth: 0.4 } }
export interface HeadFrame {
  yaw: number; pitch: number; tilt: number // degrees
  blink: number // 0 open .. 1 closed
  gaze: [number, number] // eye offset in face units (x right, y down), about +-3
  face: Face
  pose: { from: Pose; to: Pose; k: number; side?: [number, number] } // active arm of one-arm poses [from, to], +1 right
  arms?: { dw: [number, number, number][]; open: [number, number] } // per-arm wrist offset and claw open delta, [left, right]
  spark?: { side: number; r: number } // spark at a claw tip (working the table)
  eyes?: number; mouth?: number; eyeS?: number // face overrides while it eases (style, curve, eye height scale)
  eye?: RGB // eye tint (target state colour)
  tip: RGB; hold: RGB
  glitch: number; flicker: number // 0..1, 1 = full brightness
  scan?: number // 0..1 band position down the head, undefined = off
  shimmer: number // share of cells that shimmer
}
// Layout in CSS px: cell size of the head grid (scene CW*k, CH*k) and the font size.
export interface HeadLayout { cw: number; ch: number; font: number; dpr: number }
export interface HeadGL {
  canvas: HTMLCanvasElement
  layout(l: HeadLayout): { w: number; h: number; ox: number; oy: number }
  render(f: HeadFrame, t: number): void
  lost(): boolean
  dispose(): void
}

// Live head allowed here? Software renderers, reduced motion and Save-Data keep the bake.
// localStorage nest.head = 'bake' | 'gl' overrides (testing, manual toggle).
export function liveHeadOk(): boolean {
  let pref: string | null = null
  try { pref = localStorage.getItem('nest.head') } catch { /* storage off */ }
  if (pref === 'bake') return false
  if (pref !== 'gl') {
    if (matchMedia('(prefers-reduced-motion: reduce)').matches) return false
    const c = (navigator as Navigator & { connection?: { saveData?: boolean } }).connection
    if (c?.saveData) return false
  }
  const gl = document.createElement('canvas').getContext('webgl2')
  if (!gl) return false
  const dbg = gl.getExtension('WEBGL_debug_renderer_info')
  const r = String(dbg ? gl.getParameter(dbg.UNMASKED_RENDERER_WEBGL) : gl.getParameter(gl.RENDERER))
  gl.getExtension('WEBGL_lose_context')?.loseContext()
  return pref === 'gl' || !/swiftshader|llvmpipe|software|basic render/i.test(r)
}

export function createHead(canvas: HTMLCanvasElement, onLost?: () => void): HeadGL | null {
  const ctx = canvas.getContext('webgl2', { premultipliedAlpha: true, alpha: true, antialias: false })
  if (!ctx) return null
  const gl: WebGL2RenderingContext = ctx
  let isLost = false
  const lostH = (e: Event) => { e.preventDefault(); isLost = true; onLost?.() }
  canvas.addEventListener('webglcontextlost', lostH)

  const sh = (type: number, src: string) => {
    const s = gl.createShader(type)!; gl.shaderSource(s, src); gl.compileShader(s)
    if (!gl.getShaderParameter(s, gl.COMPILE_STATUS)) throw new Error(gl.getShaderInfoLog(s) ?? 'shader')
    return s
  }
  type Prog = { p: WebGLProgram; U: Record<string, WebGLUniformLocation | null> }
  const prog = (fs: string): Prog => {
    const p = gl.createProgram()!; gl.attachShader(p, sh(gl.VERTEX_SHADER, VS)); gl.attachShader(p, sh(gl.FRAGMENT_SHADER, fs))
    gl.bindAttribLocation(p, 0, 'p'); gl.linkProgram(p)
    if (!gl.getProgramParameter(p, gl.LINK_STATUS)) throw new Error(gl.getProgramInfoLog(p) ?? 'link')
    const U: Prog['U'] = {}, n = gl.getProgramParameter(p, gl.ACTIVE_UNIFORMS)
    for (let i = 0; i < n; i++) { const nm = gl.getActiveUniform(p, i)!.name.replace('[0]', ''); U[nm] = gl.getUniformLocation(p, nm) }
    return { p, U }
  }
  let P: Record<'cell' | 'glyph' | 'copy' | 'blur' | 'comp', Prog>
  try { P = { cell: prog(CELL_FS), glyph: prog(GLYPH_FS), copy: prog(COPY_FS), blur: prog(BLUR_FS), comp: prog(COMP_FS) } }
  catch (e) { console.warn('head: shader', e); return null }
  const vb = gl.createBuffer(); gl.bindBuffer(gl.ARRAY_BUFFER, vb)
  gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([-1, -1, 3, -1, -1, 3]), gl.STATIC_DRAW)
  gl.enableVertexAttribArray(0); gl.vertexAttribPointer(0, 2, gl.FLOAT, false, 0, 0)

  const textures: WebGLTexture[] = [], fbos: WebGLFramebuffer[] = []
  const tex = (w: number, h: number, filter: number) => {
    const t = gl.createTexture()!; textures.push(t); gl.bindTexture(gl.TEXTURE_2D, t)
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, null)
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter); gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter)
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE); gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
    return t
  }
  const fbo = (...ts: WebGLTexture[]) => {
    const f = gl.createFramebuffer()!; fbos.push(f); gl.bindFramebuffer(gl.FRAMEBUFFER, f)
    ts.forEach((t, i) => gl.framebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0 + i, gl.TEXTURE_2D, t, 0))
    if (ts.length > 1) gl.drawBuffers(ts.map((_, i) => gl.COLOR_ATTACHMENT0 + i))
    return f
  }
  const T = { col: tex(COLS, ROWS, gl.NEAREST), aux: tex(COLS, ROWS, gl.NEAREST) }
  const Fcell = fbo(T.col, T.aux)
  let S: { FW: number; FH: number; HW: number; HH: number; QW: number; QH: number; L: WebGLTexture; h1: WebGLTexture; h2: WebGLTexture; q1: WebGLTexture; q2: WebGLTexture; F: Record<string, WebGLFramebuffer> } | undefined
  let atlas: { t: WebGLTexture; slot: [number, number]; size: [number, number] } | undefined
  let lay: HeadLayout | undefined

  const cam = (() => {
    const f = norm(sub(CAM.tgt, CAM.pos)), r = norm(cross(V(0, 1, 0), f)), u = cross(f, r)
    return { f, r, u }
  })()

  const dropTargets = () => {
    if (!S) return
    for (const t of [S.L, S.h1, S.h2, S.q1, S.q2]) { gl.deleteTexture(t); textures.splice(textures.indexOf(t), 1) }
    for (const f of Object.values(S.F)) { gl.deleteFramebuffer(f); fbos.splice(fbos.indexOf(f), 1) }
    S = undefined
  }

  function layout(l: HeadLayout) {
    const w = COLS * l.cw, h = ROWS * l.ch
    const FW = Math.max(4, Math.round(w * l.dpr)), FH = Math.max(4, Math.round(h * l.dpr))
    if (!lay || lay.cw !== l.cw || lay.ch !== l.ch || lay.font !== l.font || lay.dpr !== l.dpr) {
      lay = l
      canvas.width = FW; canvas.height = FH
      dropTargets()
      const HW = Math.max(1, FW >> 1), HH = Math.max(1, FH >> 1), QW = Math.max(1, FW >> 2), QH = Math.max(1, FH >> 2)
      const L = tex(FW, FH, gl.LINEAR), h1 = tex(HW, HH, gl.LINEAR), h2 = tex(HW, HH, gl.LINEAR), q1 = tex(QW, QH, gl.LINEAR), q2 = tex(QW, QH, gl.LINEAR)
      S = { FW, FH, HW, HH, QW, QH, L, h1, h2, q1, q2, F: { L: fbo(L), h1: fbo(h1), h2: fbo(h2), q1: fbo(q1), q2: fbo(q2) } }
      buildAtlas(l)
    }
    return { w, h, ox: MARGIN.x * l.cw, oy: MARGIN.top * l.ch }
  }

  // Ink boost: the floor (desktop hero) and the font size (CSS px) under which it ramps up to full.
  const INK_MIN = 0.35, INK_FONT = 3.6
  // Atlas: same placement as the table canvas (Fira Code 600 with an ink stroke, centred, baseline middle + 0.5 px).
  function buildAtlas(l: HeadLayout) {
    if (atlas) { gl.deleteTexture(atlas.t); textures.splice(textures.indexOf(atlas.t), 1) }
    const cw = l.cw * l.dpr, ch = l.ch * l.dpr, fs = l.font * l.dpr
    const SW = Math.ceil(cw) + 2, SH = Math.ceil(Math.max(ch, fs * 1.4)) + 2
    const c = document.createElement('canvas'); c.width = SW * RAMP.length; c.height = SH
    const x = c.getContext('2d')!
    // Heavier weight plus a stroke keeps a bright core in each cell; it grows as the font drops below the
    // desktop size, where small cells (phone, tablet column) would blur the glyphs into a flat haze.
    const ink = Math.max(INK_MIN, Math.min(1, (INK_FONT - l.font) / INK_FONT * 2))
    x.font = `600 ${fs.toFixed(2)}px 'Fira Code', ui-monospace, monospace`; x.textBaseline = 'middle'; x.textAlign = 'center'; x.fillStyle = '#fff'
    x.strokeStyle = '#fff'; x.lineWidth = ink * 0.09 * fs; x.lineJoin = 'round'
    for (let i = 0; i < RAMP.length; i++) {
      x.fillText(RAMP[i], i * SW + SW / 2, SH / 2 + 0.5 * l.dpr)
      x.strokeText(RAMP[i], i * SW + SW / 2, SH / 2 + 0.5 * l.dpr)
    }
    const t = tex(c.width, c.height, gl.LINEAR)
    gl.bindTexture(gl.TEXTURE_2D, t); gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, c)
    atlas = { t, slot: [SW, SH], size: [c.width, c.height] }
  }

  function prims(f: HeadFrame) {
    const yaw = f.yaw * D2R, pitch = f.pitch * D2R, Wf = (p: V3) => rotY(rotX(p, pitch), yaw)
    const M: Prim[] = [], held: { c: V3; r?: number; octa?: number }[] = []
    for (const sx of [-1, 1]) {
      const ps = blendPose(POSE[f.pose.from](sx, f.pose.side?.[0]), POSE[f.pose.to](sx, f.pose.side?.[1]), f.pose.k), i = sx > 0 ? 1 : 0
      if (f.arms) { const d = f.arms.dw[i]; ps.w = add(ps.w, V(d[0], d[1], d[2])); ps.open = Math.max(0, Math.min(1.2, ps.open + f.arms.open[i])) }
      const r = armPrims(sx, ps)
      M.push(...r.M); if (r.hold) held.push(r.hold)
      if (f.spark && f.spark.side === sx && f.spark.r > 0.004) held.push({ c: r.tipC, r: f.spark.r })
    }
    M.push(...FIN)
    const pa = new Float32Array(160), pb = new Float32Array(160)
    let bc = V(0, 0, 0), n = 0
    const ends = M.map((q) => { const a = Wf(q.a), b = q.t === 's' ? a : Wf(q.b!); return { a, b } })
    ends.forEach(({ a, b }, i) => {
      const q = M[i]
      pa.set([a.x, a.y, a.z, q.r], i * 4); pb.set([b.x, b.y, b.z, q.t === 's' ? -q.m : q.m], i * 4)
      bc = add(add(bc, a), b); n += 2
    })
    bc = mul(bc, 1 / n)
    let R = 0
    ends.forEach(({ a, b }, i) => { R = Math.max(R, len(sub(a, bc)) + M[i].r + 0.05, len(sub(b, bc)) + M[i].r + 0.05) })
    let octa = [0, 0, 0, 0], hs = [0, 0, 0, 0]
    for (const h of held) {
      const c = Wf(h.c)
      if (h.octa) octa = [c.x, c.y, c.z, h.octa]
      if (h.r) hs = [c.x, c.y, c.z, h.r]
      R = Math.max(R, len(sub(c, bc)) + 0.3)
    }
    return { pa, pb, n: M.length, bound: [bc.x, bc.y, bc.z, R], octa, hs }
  }

  const pass = (p: Prog, f: WebGLFramebuffer | null, w: number, h: number) => { gl.bindFramebuffer(gl.FRAMEBUFFER, f); gl.viewport(0, 0, w, h); gl.useProgram(p.p) }
  const bindT = (p: Prog, name: string, unit: number, t: WebGLTexture) => { gl.activeTexture(gl.TEXTURE0 + unit); gl.bindTexture(gl.TEXTURE_2D, t); gl.uniform1i(p.U[name], unit) }
  const blur = (src: WebGLTexture, dst: WebGLFramebuffer, w: number, h: number, dx: number, dy: number, sigma: number, stp: number) => {
    const p = P.blur; pass(p, dst, w, h); bindT(p, 'uSrc', 0, src)
    gl.uniform2f(p.U.uRes, w, h); gl.uniform2f(p.U.uDir, dx, dy); gl.uniform1f(p.U.uSigma, sigma); gl.uniform1f(p.U.uStep, stp)
    gl.drawArrays(gl.TRIANGLES, 0, 3)
  }
  const c3 = (c: RGB) => [c[0] / 255, c[1] / 255, c[2] / 255] as const

  function render(f: HeadFrame, t: number) {
    if (isLost || !S || !atlas || !lay) return
    const sc = prims(f), dpr = lay.dpr
    let p = P.cell; pass(p, Fcell, COLS, ROWS)
    gl.uniform3f(p.U.uPos, CAM.pos.x, CAM.pos.y, CAM.pos.z); gl.uniform3f(p.U.uF, cam.f.x, cam.f.y, cam.f.z)
    gl.uniform3f(p.U.uR, cam.r.x, cam.r.y, cam.r.z); gl.uniform3f(p.U.uU, cam.u.x, cam.u.y, cam.u.z)
    gl.uniform1f(p.U.uFocal, BAKE.focal); gl.uniform2f(p.U.uC, BAKE.cx + MARGIN.x * BAKE.cw, BAKE.cy + MARGIN.top * BAKE.ch); gl.uniform1f(p.U.uRows, ROWS)
    gl.uniform1f(p.U.uYaw, f.yaw * D2R); gl.uniform1f(p.U.uPitch, f.pitch * D2R); gl.uniform1f(p.U.uTilt, f.tilt); gl.uniform1f(p.U.uBlink, f.blink)
    gl.uniform1i(p.U.uMood, f.eyes ?? FACE_STYLE[f.face].eyes); gl.uniform1f(p.U.uMouth, f.mouth ?? FACE_STYLE[f.face].mouth); gl.uniform1f(p.U.uEyeS, f.eyeS ?? 1); gl.uniform2f(p.U.uGaze, f.gaze[0], f.gaze[1])
    gl.uniform3f(p.U.uTip, ...c3(f.tip)); gl.uniform3f(p.U.uHoldCol, ...c3(f.hold))
    gl.uniform3f(p.U.uEye, ...c3(f.eye ?? [0, 0, 0])); gl.uniform1f(p.U.uEyeK, f.eye ? 1 : 0)
    gl.uniform4fv(p.U.uPA, sc.pa); gl.uniform4fv(p.U.uPB, sc.pb); gl.uniform1i(p.U.uNP, sc.n)
    gl.uniform4fv(p.U.uBound, sc.bound); gl.uniform4fv(p.U.uOcta, sc.octa); gl.uniform4fv(p.U.uHoldS, sc.hs)
    // Scan band runs from above the fin to below the chin, in virtual px.
    const top = MARGIN.top * BAKE.ch, span = BAKE.rows * BAKE.ch
    gl.uniform1f(p.U.uScan, f.scan === undefined ? -100 : top + f.scan * span)
    gl.uniform1f(p.U.uShimmer, f.shimmer); gl.uniform1f(p.U.uFrame, Math.floor(t * 30) % 997)
    gl.drawArrays(gl.TRIANGLES, 0, 3)

    p = P.glyph; pass(p, S.F.L, S.FW, S.FH)
    bindT(p, 'uCol', 0, T.col); bindT(p, 'uAux', 1, T.aux); bindT(p, 'uAtlas', 2, atlas.t)
    gl.uniform2f(p.U.uCell, lay.cw * dpr, lay.ch * dpr); gl.uniform2f(p.U.uSlot, ...atlas.slot); gl.uniform2f(p.U.uAtlasSize, ...atlas.size)
    gl.uniform2f(p.U.uGrid, COLS, ROWS); gl.uniform1f(p.U.uH, S.FH)
    gl.uniform1f(p.U.uTime, t); gl.uniform1f(p.U.uGlitch, f.glitch); gl.uniform1f(p.U.uDpr, dpr); gl.uniform1f(p.U.uAlphaK, 1.5)
    gl.drawArrays(gl.TRIANGLES, 0, 3)

    p = P.copy; pass(p, S.F.h1, S.HW, S.HH); bindT(p, 'uSrc', 0, S.L); gl.uniform2f(p.U.uRes, S.HW, S.HH); gl.drawArrays(gl.TRIANGLES, 0, 3)
    pass(p, S.F.q1, S.QW, S.QH); bindT(p, 'uSrc', 0, S.h1); gl.uniform2f(p.U.uRes, S.QW, S.QH); gl.drawArrays(gl.TRIANGLES, 0, 3)
    blur(S.h1, S.F.h2, S.HW, S.HH, 1, 0, 1.5 * dpr, 0.5 * dpr); blur(S.h2, S.F.h1, S.HW, S.HH, 0, 1, 1.5 * dpr, 0.5 * dpr)
    blur(S.q1, S.F.q2, S.QW, S.QH, 1, 0, 3.5 * dpr, 0.6 * dpr); blur(S.q2, S.F.q1, S.QW, S.QH, 0, 1, 3.5 * dpr, 0.6 * dpr)

    p = P.comp; pass(p, null, S.FW, S.FH)
    bindT(p, 'uL', 0, S.L); bindT(p, 'uBs', 1, S.h1); bindT(p, 'uBw', 2, S.q1)
    gl.uniform2f(p.U.uRes, S.FW, S.FH); gl.uniform1f(p.U.uGlitch, f.glitch); gl.uniform1f(p.U.uDpr, dpr); gl.uniform1f(p.U.uFlick, f.flicker)
    gl.clearColor(0, 0, 0, 0); gl.clear(gl.COLOR_BUFFER_BIT)
    gl.drawArrays(gl.TRIANGLES, 0, 3)
  }

  return {
    canvas, layout, render, lost: () => isLost,
    dispose() {
      canvas.removeEventListener('webglcontextlost', lostH)
      if (isLost) return
      textures.forEach((t) => gl.deleteTexture(t)); fbos.forEach((f) => gl.deleteFramebuffer(f))
      Object.values(P).forEach((p) => gl.deleteProgram(p.p)); gl.deleteBuffer(vb)
    },
  }
}
