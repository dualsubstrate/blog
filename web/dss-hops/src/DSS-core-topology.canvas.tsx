import {
  Button,
  Callout,
  Code,
  CollapsibleSection,
  Divider,
  Grid,
  H1,
  H2,
  H3,
  Pill,
  Row,
  Stack,
  Stat,
  Table,
  Text,
  useCanvasState,
  useEffect,
  useHostTheme,
  useRef,
  useState,
} from "cursor/canvas";

type Via = "Ring" | "cube" | "M";
type HopFilter = "all" | "Ring" | "M" | "cube";
type HopTableTab = "ring" | "cube" | "primes" | "resolution";

type WalkHop = { from: number; to: number; via: Via };

const CORNER_PRIME = [2, 3, 5, 7, 11, 13, 17, 19] as const;
const M_BIND_PRIMES = "137 × 139";

function hopMultiply(from: number, to: number, via: Via): string {
  const a = CORNER_PRIME[from];
  const b = CORNER_PRIME[to];
  if (via === "M") return `${a} × ${M_BIND_PRIMES} × ${b}`;
  return `${a} × ${b}`;
}

const WALK_HOPS: WalkHop[] = [
  { from: 0, to: 1, via: "M" },
  { from: 1, to: 2, via: "Ring" },
  { from: 2, to: 3, via: "M" },
  { from: 3, to: 0, via: "Ring" },
  { from: 4, to: 5, via: "M" },
  { from: 5, to: 6, via: "Ring" },
  { from: 6, to: 7, via: "M" },
  { from: 7, to: 4, via: "Ring" },
];

function hopKey(from: number, to: number): string {
  return `${from}-${to}`;
}

function hopMatchesFilter(via: Via, filter: HopFilter): boolean {
  if (filter === "cube") return false;
  if (filter === "all") return via === "Ring" || via === "M";
  return via === filter;
}

type TetraId = "R0" | "R1";
type TetraFocus = "both" | TetraId;

function tetraShade(id: TetraId, focus: TetraFocus) {
  const on = focus === "both" || focus === id;
  const isolated = focus === "R0" || focus === "R1";
  if (on && isolated) {
    return { front: 0.52, back: 0.08, stroke: 0.85, edge: 1, edgeW: 2.5 };
  }
  if (on) {
    return { front: 0.18, back: 0.04, stroke: 0.35, edge: 0.65, edgeW: 1.8 };
  }
  return { front: 0.04, back: 0.01, stroke: 0.1, edge: 0.14, edgeW: 1.15 };
}

function tetraVertexIds(id: TetraId): [number, number, number, number] {
  return id === "R0" ? [0, 1, 2, 3] : [4, 5, 6, 7];
}

function walkArrowPoints(from: Pt, to: Pt, size: number): string {
  const dx = to.x - from.x;
  const dy = to.y - from.y;
  const len = Math.hypot(dx, dy) || 1;
  const ux = dx / len;
  const uy = dy / len;
  const px = -uy;
  const py = ux;
  const mx = from.x + ux * len * 0.62;
  const my = from.y + uy * len * 0.62;
  const tipX = mx + ux * size;
  const tipY = my + uy * size;
  const baseX = mx - ux * size * 0.28;
  const baseY = my - uy * size * 0.28;
  const wing = size * 0.7;
  return `${tipX},${tipY} ${baseX + px * wing},${baseY + py * wing} ${baseX - px * wing},${baseY - py * wing}`;
}

function perspectiveFade(scale: number, scales: number[]): number {
  let lo = scales[0];
  let hi = scales[0];
  for (const s of scales) {
    if (s < lo) lo = s;
    if (s > hi) hi = s;
  }
  const span = hi - lo;
  if (span < 1e-6) return 1;
  const u = Math.max(0, Math.min(1, (scale - lo) / span));
  return 0.12 + 0.88 * u * u;
}

type NestStage = "corners" | "lattice" | "stack";

const CUBE_CORNERS: Array<{ id: number; x: number; y: number; z: number }> = [
  { id: 0, x: 0, y: 0, z: 0 },
  { id: 1, x: 0, y: 1, z: 1 },
  { id: 2, x: 1, y: 0, z: 1 },
  { id: 3, x: 1, y: 1, z: 0 },
  { id: 4, x: 0, y: 0, z: 1 },
  { id: 5, x: 0, y: 1, z: 0 },
  { id: 6, x: 1, y: 0, z: 0 },
  { id: 7, x: 1, y: 1, z: 1 },
];

const CUBE_EDGES: Array<[number, number]> = [
  [0, 6],
  [0, 5],
  [0, 4],
  [7, 1],
  [7, 2],
  [7, 3],
  [6, 3],
  [6, 2],
  [5, 3],
  [5, 1],
  [4, 2],
  [4, 1],
];

const CUBE_FACES: Array<[number, number, number, number]> = [
  [0, 4, 1, 5],
  [6, 3, 7, 2],
  [0, 6, 2, 4],
  [5, 3, 7, 1],
  [0, 5, 3, 6],
  [4, 1, 7, 2],
];

const CUBE_WALKS: Array<{ from: number; to: number; kind: "cube" }> = [
  { from: 0, to: 4, kind: "cube" },
  { from: 0, to: 6, kind: "cube" },
  { from: 1, to: 5, kind: "cube" },
  { from: 1, to: 7, kind: "cube" },
  { from: 2, to: 4, kind: "cube" },
  { from: 2, to: 6, kind: "cube" },
  { from: 3, to: 5, kind: "cube" },
  { from: 3, to: 7, kind: "cube" },
  { from: 1, to: 4, kind: "cube" },
  { from: 3, to: 6, kind: "cube" },
  { from: 5, to: 0, kind: "cube" },
  { from: 7, to: 2, kind: "cube" },
];

function cubeEdgeMultiply(from: number, to: number): string {
  return `${CORNER_PRIME[from]} × ${CORNER_PRIME[to]}`;
}

function cubeEdgeKind(a: number, b: number): "cube" | undefined {
  const walk = CUBE_WALKS.find(
    (h) => (h.from === a && h.to === b) || (h.from === b && h.to === a),
  );
  return walk?.kind;
}

function hopKeyMatches(a: number, b: number, key: string): boolean {
  return key === hopKey(a, b) || key === hopKey(b, a);
}

const TETRA_R0_EDGES: Array<[number, number]> = [
  [0, 1],
  [0, 2],
  [0, 3],
  [1, 2],
  [1, 3],
  [2, 3],
];

const TETRA_R1_EDGES: Array<[number, number]> = [
  [4, 5],
  [4, 6],
  [4, 7],
  [5, 6],
  [5, 7],
  [6, 7],
];

const TETRA_R0_FACES: Array<[number, number, number]> = [
  [0, 1, 2],
  [0, 1, 3],
  [0, 2, 3],
  [1, 2, 3],
];

const TETRA_R1_FACES: Array<[number, number, number]> = [
  [4, 5, 6],
  [4, 5, 7],
  [4, 6, 7],
  [5, 6, 7],
];

type Mat3 = [
  [number, number, number],
  [number, number, number],
  [number, number, number],
];

function matMul(a: Mat3, b: Mat3): Mat3 {
  return [0, 1, 2].map((i) =>
    [0, 1, 2].map(
      (j) => a[i][0] * b[0][j] + a[i][1] * b[1][j] + a[i][2] * b[2][j],
    ),
  ) as Mat3;
}

function matKey(m: Mat3): string {
  return m.flat().join(",");
}

const MAT_I: Mat3 = [
  [1, 0, 0],
  [0, 1, 0],
  [0, 0, 1],
];
const MAT_RX: Mat3 = [
  [1, 0, 0],
  [0, 0, -1],
  [0, 1, 0],
];
const MAT_RY: Mat3 = [
  [0, 0, 1],
  [0, 1, 0],
  [-1, 0, 0],
];
const MAT_RZ: Mat3 = [
  [0, -1, 0],
  [1, 0, 0],
  [0, 0, 1],
];

function buildCubePoses(): Mat3[] {
  const seen = new Set<string>([matKey(MAT_I)]);
  const list: Mat3[] = [MAT_I];
  const queue: Mat3[] = [MAT_I];
  while (queue.length > 0) {
    const cur = queue.shift();
    if (!cur) break;
    for (const g of [MAT_RX, MAT_RY, MAT_RZ]) {
      const next = matMul(g, cur);
      const key = matKey(next);
      if (!seen.has(key)) {
        seen.add(key);
        list.push(next);
        queue.push(next);
      }
    }
  }
  return list;
}

const CUBE_POSES = buildCubePoses();

function poseAfter(poseIndex: number, gen: Mat3): number {
  const next = matMul(gen, CUBE_POSES[poseIndex] ?? MAT_I);
  const key = matKey(next);
  const idx = CUBE_POSES.findIndex((m) => matKey(m) === key);
  return idx < 0 ? poseIndex : idx;
}

function applyPose(m: Mat3, x: number, y: number, z: number) {
  const x0 = x - 0.5;
  const y0 = y - 0.5;
  const z0 = z - 0.5;
  return {
    x: m[0][0] * x0 + m[0][1] * y0 + m[0][2] * z0 + 0.5,
    y: m[1][0] * x0 + m[1][1] * y0 + m[1][2] * z0 + 0.5,
    z: m[2][0] * x0 + m[2][1] * y0 + m[2][2] * z0 + 0.5,
  };
}

type Pt = {
  x: number;
  y: number;
  depth: number;
  scale: number;
  cx: number;
  cy: number;
};

const CAMERA_Z = 4.2;

function depthAlpha(depth: number, lo = 0.45, hi = 1): number {
  const t = (depth + 0.88) / 1.76;
  const u = Math.max(0, Math.min(1, t));
  return lo + (hi - lo) * u;
}

function facingCamera(a: Pt, b: Pt, c: Pt): boolean {
  const ux = b.cx - a.cx;
  const uy = b.cy - a.cy;
  const vx = c.cx - a.cx;
  const vy = c.cy - a.cy;
  return ux * vy - uy * vx > 0;
}

function edgeHalves(
  a: number,
  b: number,
  project: (x: number, y: number, z: number) => Pt,
): Array<{ pa: Pt; pb: Pt; z: number }> {
  const ca = CUBE_CORNERS[a];
  const cb = CUBE_CORNERS[b];
  const pa = project(ca.x, ca.y, ca.z);
  const pb = project(cb.x, cb.y, cb.z);
  const mid = project((ca.x + cb.x) / 2, (ca.y + cb.y) / 2, (ca.z + cb.z) / 2);
  return [
    { pa, pb: mid, z: (pa.depth + mid.depth) / 2 },
    { pa: mid, pb, z: (mid.depth + pb.depth) / 2 },
  ];
}

function project3(
  x: number,
  y: number,
  z: number,
  ox: number,
  oy: number,
  s: number,
  yaw: number,
  pitch: number,
): Pt {
  const cx = x - 0.5;
  const cy = y - 0.5;
  const cz = z - 0.5;
  const x1 = cx * Math.cos(yaw) - cz * Math.sin(yaw);
  const z1 = cx * Math.sin(yaw) + cz * Math.cos(yaw);
  const y1 = cy * Math.cos(pitch) - z1 * Math.sin(pitch);
  const z2 = cy * Math.sin(pitch) + z1 * Math.cos(pitch);
  const scale = CAMERA_Z / (CAMERA_Z - z2);
  return {
    x: ox + x1 * s * scale,
    y: oy - y1 * s * scale,
    depth: z2,
    scale,
    cx: x1,
    cy: y1,
  };
}

function orientToHop(fromId: number, toId: number): { yaw: number; pitch: number } {
  const A = CUBE_CORNERS[fromId];
  const B = CUBE_CORNERS[toId];
  const mx = (A.x + B.x) / 2 - 0.5;
  const my = (A.y + B.y) / 2 - 0.5;
  const mz = (A.z + B.z) / 2 - 0.5;
  const pitch = 0.48;
  const glance = 0.82;

  function midDepth(yaw: number): number {
    const z1 = mx * Math.sin(yaw) + mz * Math.cos(yaw);
    return my * Math.sin(pitch) + z1 * Math.cos(pitch);
  }

  let yaw = Math.atan2(mx, mz);
  if (midDepth(yaw) < 0) yaw += Math.PI;
  yaw += glance;
  return { yaw, pitch };
}

function shortestAngle(from: number, to: number): number {
  let d = to - from;
  while (d > Math.PI) d -= Math.PI * 2;
  while (d < -Math.PI) d += Math.PI * 2;
  return d;
}

function easeInOutCubic(t: number): number {
  return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

const HOP_TURN_MS = 1600;

function CubeHopView({
  from,
  to,
  via,
  filter,
  onSelect,
}: {
  from: number;
  to: number;
  via: Via;
  filter: HopFilter;
  onSelect: (key: string) => void;
}) {
  const theme = useHostTheme();
  const ox = 260;
  const oy = 135;
  const s = 105;
  const selectedKey = hopKey(from, to);
  const [yaw, setYaw] = useCanvasState("cubeYaw", 0.62);
  const [pitch, setPitch] = useCanvasState("cubePitch", 0.48);
  const [dragging, setDragging] = useState(false);
  const viewRef = useRef({ yaw, pitch });
  const rafRef = useRef<number | null>(null);
  const drag = useRef({ x: 0, y: 0, active: false });
  viewRef.current = { yaw, pitch };

  useEffect(() => {
    const target = orientToHop(from, to);
    const startYaw = viewRef.current.yaw;
    const startPitch = viewRef.current.pitch;
    const dYaw = shortestAngle(startYaw, target.yaw);
    const dPitch = shortestAngle(startPitch, target.pitch);
    if (Math.abs(dYaw) < 0.001 && Math.abs(dPitch) < 0.001) return;
    const t0 = performance.now();
    if (rafRef.current !== null) cancelAnimationFrame(rafRef.current);

    function tick(now: number) {
      const t = Math.min(1, (now - t0) / HOP_TURN_MS);
      const e = easeInOutCubic(t);
      const next = {
        yaw: startYaw + dYaw * e,
        pitch: startPitch + dPitch * e,
      };
      viewRef.current = next;
      setYaw(next.yaw);
      setPitch(next.pitch);
      if (t < 1) rafRef.current = requestAnimationFrame(tick);
      else rafRef.current = null;
    }

    rafRef.current = requestAnimationFrame(tick);
    return () => {
      if (rafRef.current !== null) cancelAnimationFrame(rafRef.current);
    };
  }, [from, to]);

  const r0 = theme.category.blue;
  const r1 = theme.category.orange;
  const accent = theme.accent.primary;
  const onAccent = theme.text.onAccent;
  const p = (x: number, y: number, z: number) =>
    project3(x, y, z, ox, oy, s, yaw, pitch);
  const cornerPt = (id: number) => {
    const c = CUBE_CORNERS[id];
    return p(c.x, c.y, c.z);
  };

  const cubeDim = filter === "Ring" || filter === "M";
  const showTetra = via === "Ring" || via === "M" || cubeDim;
  const tetraIsR0 = from < 4;
  const tetraFocus: TetraFocus = showTetra
    ? tetraIsR0
      ? "R0"
      : "R1"
    : "both";
  const muteCube = showTetra;
  const tetras: Array<{
    id: TetraId;
    faces: Array<[number, number, number]>;
    edges: Array<[number, number]>;
    fill: string;
  }> = [
    { id: "R0", faces: TETRA_R0_FACES, edges: TETRA_R0_EDGES, fill: r0 },
    { id: "R1", faces: TETRA_R1_FACES, edges: TETRA_R1_EDGES, fill: r1 },
  ];
  const layers = [];
  for (const face of CUBE_FACES) {
    const pts = face.map((id) => cornerPt(id));
    const z = pts.reduce((sum, q) => sum + q.depth, 0) / 4;
    const front = facingCamera(pts[0], pts[1], pts[2]);
    layers.push({
      key: `face-${face.join("-")}`,
      z,
      node: (
        <polygon
          points={pts.map((q) => `${q.x},${q.y}`).join(" ")}
          fill={theme.fill.secondary}
          fillOpacity={(front ? 0.16 : 0.05) * depthAlpha(z, 0.55, 1)}
          stroke={muteCube ? "none" : theme.stroke.tertiary}
          strokeWidth={muteCube ? 0 : 1}
        />
      ),
    });
  }
  for (const [a, b] of CUBE_EDGES) {
    const active = via === "cube" && hopKeyMatches(a, b, selectedKey);
    const walk = CUBE_WALKS.find(
      (h) => hopKeyMatches(a, b, hopKey(h.from, h.to)),
    );
    const pa = cornerPt(a);
    const pb = cornerPt(b);
    const z = (pa.depth + pb.depth) / 2;
    const scale = (pa.scale + pb.scale) / 2;
    if (muteCube) {
      layers.push({
        key: `cube-dot-${a}-${b}`,
        z: z - 0.01,
        node: (
          <line
            x1={pa.x}
            y1={pa.y}
            x2={pb.x}
            y2={pb.y}
            stroke={theme.stroke.secondary}
            strokeWidth={1.7 * scale}
            strokeDasharray="5 4"
            opacity={0.72 * depthAlpha(z, 0.55, 1)}
          />
        ),
      });
      continue;
    }
    edgeHalves(a, b, p).forEach((half, i) => {
      const hs = (half.pa.scale + half.pb.scale) / 2;
      layers.push({
        key: `edge-${a}-${b}-${i}`,
        z: half.z,
        node: (
          <g
            opacity={
              (active ? 1 : cubeDim ? 0.35 : 1) * depthAlpha(half.z, 0.4, 1)
            }
            onClick={() => {
              if (walk) onSelect(hopKey(walk.from, walk.to));
            }}
            style={{ cursor: walk ? "pointer" : "default" }}
          >
            <line
              x1={half.pa.x}
              y1={half.pa.y}
              x2={half.pb.x}
              y2={half.pb.y}
              stroke={accent}
              strokeOpacity={0.01}
              strokeWidth={16}
            />
            <line
              x1={half.pa.x}
              y1={half.pa.y}
              x2={half.pb.x}
              y2={half.pb.y}
              stroke={active ? accent : theme.stroke.tertiary}
              strokeWidth={(active ? 4.5 : 1.75) * hs}
            />
          </g>
        ),
      });
    });
  }
  for (const tetra of tetras) {
    if (!showTetra) continue;
    const ghost = tetraFocus !== "both" && tetraFocus !== tetra.id;
    const shade = tetraShade(tetra.id, tetraFocus);
    const scales = tetraVertexIds(tetra.id).map((id) => cornerPt(id).scale);
    const fill = ghost ? theme.fill.tertiary : tetra.fill;
    for (const face of tetra.faces) {
      const pts = face.map((id) => cornerPt(id));
      const z = pts.reduce((sum, q) => sum + q.depth, 0) / 3;
      const scale = pts.reduce((sum, q) => sum + q.scale, 0) / 3;
      const fade = perspectiveFade(scale, scales);
      const front = facingCamera(pts[0], pts[1], pts[2]);
      const faceOp = ghost
        ? (front ? 0.08 : 0.03) * fade
        : (front ? shade.front : shade.back) * fade;
      layers.push({
        key: `tetra-face-${tetra.id}-${face.join("-")}`,
        z: z + 0.001,
        node: (
          <polygon
            points={pts.map((q) => `${q.x},${q.y}`).join(" ")}
            fill={fill}
            fillOpacity={faceOp}
            stroke="none"
          />
        ),
      });
    }
    for (const [a, b] of tetra.edges) {
      const walk = WALK_HOPS.find((h) =>
        hopKeyMatches(a, b, hopKey(h.from, h.to)),
      );
      const hopActive = Boolean(walk && hopKeyMatches(a, b, selectedKey));
      const kindLit = Boolean(
        walk &&
          cubeDim &&
          walk.via === filter &&
          (tetraIsR0 ? walk.from < 4 : walk.from >= 4),
      );
      const pa = cornerPt(a);
      const pb = cornerPt(b);
      const z = (pa.depth + pb.depth) / 2;
      const scale = (pa.scale + pb.scale) / 2;
      const fade = perspectiveFade(scale, scales);
      const edgeStroke = hopActive || kindLit ? tetra.fill : ghost
        ? theme.stroke.tertiary
        : tetra.fill;
      const edgeW = hopActive ? 4.5 : kindLit ? 3.2 : ghost ? 0.55 : 0.7;
      layers.push({
        key: `tetra-edge-${tetra.id}-${a}-${b}`,
        z: z + 0.003,
        node: (
          <g
            opacity={
              hopActive || kindLit
                ? fade
                : ghost
                  ? 0.45 * fade
                  : shade.edge * fade * 0.85
            }
            onClick={() => {
              if (walk) onSelect(hopKey(walk.from, walk.to));
            }}
            style={{ cursor: walk ? "pointer" : "default" }}
          >
            <line
              x1={pa.x}
              y1={pa.y}
              x2={pb.x}
              y2={pb.y}
              stroke={accent}
              strokeOpacity={0.01}
              strokeWidth={16}
            />
            <line
              x1={pa.x}
              y1={pa.y}
              x2={pb.x}
              y2={pb.y}
              stroke={edgeStroke}
              strokeWidth={edgeW * scale}
            />
          </g>
        ),
      });
    }
  }
  {
    const pa = cornerPt(from);
    const pb = cornerPt(to);
    const z = (pa.depth + pb.depth) / 2;
    const scale = (pa.scale + pb.scale) / 2;
    const fill = via === "cube" ? accent : from < 4 ? r0 : r1;
    layers.push({
      key: "walk-arrow",
      z: z + 0.015,
      node: (
        <polygon
          points={walkArrowPoints(pa, pb, 12 * scale)}
          fill={fill}
        />
      ),
    });
  }
  if (via === "M") {
    const A = CUBE_CORNERS[from];
    const B = CUBE_CORNERS[to];
    const mPt = p((A.x + B.x) / 2, (A.y + B.y) / 2, (A.z + B.z) / 2);
    const mark = 7 * mPt.scale;
    const tetraFill = from < 4 ? r0 : r1;
    layers.push({
      key: "m-bind",
      z: mPt.depth + 0.012,
      node: (
        <g opacity={depthAlpha(mPt.depth, 0.55, 1)}>
          <rect
            x={mPt.x - mark}
            y={mPt.y - mark}
            width={mark * 2}
            height={mark * 2}
            transform={`rotate(45 ${mPt.x} ${mPt.y})`}
            fill={tetraFill}
            stroke={theme.bg.editor}
            strokeWidth={1.25}
          />
          <text
            x={mPt.x}
            y={mPt.y}
            dominantBaseline="central"
            textAnchor="middle"
            fill={onAccent}
            fontSize={Math.max(7, 8 * mPt.scale)}
            fontWeight={600}
          >
            M
          </text>
        </g>
      ),
    });
  }
    for (const c of CUBE_CORNERS) {
      const q = p(c.x, c.y, c.z);
      const isEnd = c.id === from || c.id === to;
      const isR0 = c.id < 4;
      layers.push({
        key: `node-${c.id}`,
        z: q.depth + 0.02,
        node: (
          <g opacity={depthAlpha(q.depth, 0.5, 1)}>
            <circle
              cx={q.x}
              cy={q.y}
              r={(isEnd ? 12 : 9.5) * q.scale}
              fill={isR0 ? r0 : r1}
              stroke={isEnd ? accent : theme.bg.editor}
              strokeWidth={isEnd ? 2.5 : 1.5}
            />
            <text
              x={q.x}
              y={q.y}
              dominantBaseline="central"
              textAnchor="middle"
              fill={theme.text.primary}
              fontSize={Math.max(7.5, 9.5 * q.scale)}
              fontWeight={600}
            >
              {c.id}
            </text>
          </g>
        ),
      });
    }
  layers.sort((a, b) => a.z - b.z);

  function onPointerDown(event: {
    currentTarget: { setPointerCapture: (id: number) => void };
    pointerId: number;
    clientX: number;
    clientY: number;
  }) {
    if (rafRef.current !== null) {
      cancelAnimationFrame(rafRef.current);
      rafRef.current = null;
    }
    event.currentTarget.setPointerCapture(event.pointerId);
    drag.current = { x: event.clientX, y: event.clientY, active: true };
    setDragging(true);
  }

  function onPointerMove(event: { clientX: number; clientY: number }) {
    if (!drag.current.active) return;
    const dx = event.clientX - drag.current.x;
    const dy = event.clientY - drag.current.y;
    drag.current.x = event.clientX;
    drag.current.y = event.clientY;
    const next = {
      yaw: viewRef.current.yaw + dx * 0.01,
      pitch: Math.max(-1.2, Math.min(1.2, viewRef.current.pitch - dy * 0.01)),
    };
    viewRef.current = next;
    setYaw(next.yaw);
    setPitch(next.pitch);
  }

  function onPointerUp(event: {
    currentTarget: {
      hasPointerCapture: (id: number) => boolean;
      releasePointerCapture: (id: number) => void;
    };
    pointerId: number;
  }) {
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
    drag.current.active = false;
    setDragging(false);
  }

  return (
    <svg
      viewBox="0 0 520 270"
      width="100%"
      preserveAspectRatio="xMidYMid meet"
      role="img"
      aria-label={`Cube rotating until hop ${from} to ${to} sits in front at an angle. Drag to orbit.`}
      style={{
        touchAction: "none",
        cursor: dragging ? "grabbing" : "grab",
        userSelect: "none",
        maxHeight: 290,
      }}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerUp}
    >
      {layers.map((layer) => (
        <g key={layer.key}>{layer.node}</g>
      ))}
    </svg>
  );
}

function CubeNest({
  stage,
  highlightHopKey,
}: {
  stage: NestStage;
  highlightHopKey: string;
}) {
  const theme = useHostTheme();
  const ox = 270;
  const oy = 150;
  const s = stage === "stack" ? 100 : 114;
  const r0 = theme.category.blue;
  const r1 = theme.category.orange;
  const edgeStroke = theme.stroke.secondary;
  const mute = theme.stroke.tertiary;
  const fillMed = theme.fill.tertiary;
  const accent = theme.accent.primary;
  const onAccent = theme.text.onAccent;
  const label = theme.text.secondary;
  const [yaw, setYaw] = useCanvasState("cubeYaw", 0.62);
  const [pitch, setPitch] = useCanvasState("cubePitch", 0.48);
  const [indexExtra, setIndexExtra] = useCanvasState<
    Array<{ kind: "Roll" | "FF"; label: string; afterStep: number }>
  >("indexExtra", []);
  const [poseIndex, setPoseIndex] = useCanvasState("cubePose", 0);
  const [nestTetra, setNestTetra] = useCanvasState<TetraFocus>(
    "nestTetra",
    "both",
  );
  const [dragging, setDragging] = useState(false);
  const drag = useRef({ x: 0, y: 0, active: false });

  const pose = CUBE_POSES[poseIndex] ?? MAT_I;
  const p = (x: number, y: number, z: number) => {
    const q = applyPose(pose, x, y, z);
    return project3(q.x, q.y, q.z, ox, oy, s, yaw, pitch);
  };

  function recordR(gen: Mat3) {
    const next = poseAfter(poseIndex, gen);
    if (next === poseIndex) return;
    setPoseIndex(next);
    setIndexExtra((prev) => {
      const withoutRoll = prev.filter((e) => e.kind !== "Roll");
      return [
        ...withoutRoll,
        {
          kind: "Roll",
          label: `Roll ${next + 1}/${CUBE_POSES.length}`,
          afterStep: 0,
        },
      ];
    });
  }
  const cornerPt = (id: number) => {
    const c = CUBE_CORNERS[id];
    return p(c.x, c.y, c.z);
  };

  const showLattice = stage === "lattice";
  const showStack = stage === "stack";
  const caption =
    stage === "corners"
      ? "R0 / R1 isolates one tetrahedron; the cube stays as a faint dotted cage. Drag to orbit. Closed tours are on the tours canvas."
      : stage === "lattice"
        ? "Drag to orbit. 8 corners + 12 edges + 6 faces + 1 center = 27. L=0, 1, 2 are slices of this cube."
        : "Drag to orbit. A finished 3×3×3 occupies one corner of the cube above. Same 0–7, one layer up.";

  const nestLayers = [];
  for (const [a, b] of CUBE_EDGES) {
    const active = hopKeyMatches(a, b, highlightHopKey);
    const pa = cornerPt(a);
    const pb = cornerPt(b);
    const z = (pa.depth + pb.depth) / 2;
    const scale = (pa.scale + pb.scale) / 2;
    if (nestTetra !== "both") {
      nestLayers.push({
        key: `cube-dot-${a}-${b}`,
        z: z - 0.01,
        node: (
          <line
            x1={pa.x}
            y1={pa.y}
            x2={pb.x}
            y2={pb.y}
            stroke={theme.stroke.secondary}
            strokeWidth={1.7 * scale}
            strokeDasharray="5 4"
            opacity={0.72 * depthAlpha(z, 0.55, 1)}
          />
        ),
      });
      continue;
    }
    const stroke = active ? accent : mute;
    const width = active ? 2.5 : 1.1;
    edgeHalves(a, b, p).forEach((half, i) => {
      const hs = (half.pa.scale + half.pb.scale) / 2;
      nestLayers.push({
        key: `cube-${a}-${b}-${i}`,
        z: half.z,
        node: (
          <line
            x1={half.pa.x}
            y1={half.pa.y}
            x2={half.pb.x}
            y2={half.pb.y}
            stroke={stroke}
            strokeWidth={width * hs}
            opacity={depthAlpha(half.z, 0.35, 0.9)}
          />
        ),
      });
    });
  }
  if (!showStack) {
    const nestTetras: Array<{
      id: TetraId;
      faces: Array<[number, number, number]>;
      edges: Array<[number, number]>;
      fill: string;
    }> = [
      { id: "R0", faces: TETRA_R0_FACES, edges: TETRA_R0_EDGES, fill: r0 },
      { id: "R1", faces: TETRA_R1_FACES, edges: TETRA_R1_EDGES, fill: r1 },
    ];
    for (const tetra of nestTetras) {
      const ghost = nestTetra !== "both" && nestTetra !== tetra.id;
      const shade = tetraShade(tetra.id, nestTetra);
      const scales = tetraVertexIds(tetra.id).map((id) => cornerPt(id).scale);
      const fill = ghost ? theme.fill.tertiary : tetra.fill;
      const edgeStroke = ghost ? theme.stroke.tertiary : tetra.fill;
      for (const face of tetra.faces) {
        const a = cornerPt(face[0]);
        const b = cornerPt(face[1]);
        const c = cornerPt(face[2]);
        const z = (a.depth + b.depth + c.depth) / 3;
        const scale = (a.scale + b.scale + c.scale) / 3;
        const fade = perspectiveFade(scale, scales);
        const front = facingCamera(a, b, c);
        const faceOp = ghost
          ? (front ? 0.08 : 0.03) * fade
          : (front ? shade.front : shade.back) * fade;
        nestLayers.push({
          key: `${tetra.id}-${face.join("-")}`,
          z,
          node: (
            <polygon
              points={`${a.x},${a.y} ${b.x},${b.y} ${c.x},${c.y}`}
              fill={fill}
              fillOpacity={faceOp}
              stroke="none"
            />
          ),
        });
      }
      for (const [a, b] of tetra.edges) {
        const pa = cornerPt(a);
        const pb = cornerPt(b);
        const z = (pa.depth + pb.depth) / 2;
        const scale = (pa.scale + pb.scale) / 2;
        const fade = perspectiveFade(scale, scales);
        nestLayers.push({
          key: `${tetra.id}-edge-${a}-${b}`,
          z: z + 0.002,
          node: (
            <line
              x1={pa.x}
              y1={pa.y}
              x2={pb.x}
              y2={pb.y}
              stroke={edgeStroke}
              strokeWidth={(ghost ? 0.55 : shade.edgeW) * scale}
              opacity={ghost ? 0.45 * fade : shade.edge * fade}
            />
          ),
        });
      }
    }
    for (const c of CUBE_CORNERS) {
      const q = p(c.x, c.y, c.z);
      const isR0 = c.id < 4;
      const walker = highlightHopKey.split("-").map(Number).includes(c.id);
      const base = showLattice ? 7.5 : 10;
      const onTetra =
        nestTetra === "both" ||
        (nestTetra === "R0" && isR0) ||
        (nestTetra === "R1" && !isR0);
      nestLayers.push({
        key: `node-${c.id}`,
        z: q.depth + 0.02,
        node: (
          <g opacity={depthAlpha(q.depth, 0.5, 1) * (onTetra ? 1 : 0.22)}>
            <circle
              cx={q.x}
              cy={q.y}
              r={base * q.scale}
              fill={isR0 ? r0 : r1}
              stroke={walker ? accent : theme.bg.editor}
              strokeWidth={walker ? 2.5 : 1.75}
            />
            <text
              x={q.x}
              y={q.y}
              dominantBaseline="central"
              textAnchor="middle"
              fill={theme.text.primary}
              fontSize={Math.max(7.5, 9.5 * q.scale)}
              fontWeight={600}
            >
              {c.id}
            </text>
          </g>
        ),
      });
    }
  }
  nestLayers.sort((a, b) => a.z - b.z);

  function onPointerDown(event: {
    currentTarget: { setPointerCapture: (id: number) => void };
    pointerId: number;
    clientX: number;
    clientY: number;
  }) {
    event.currentTarget.setPointerCapture(event.pointerId);
    drag.current = { x: event.clientX, y: event.clientY, active: true };
    setDragging(true);
  }

  function onPointerMove(event: { clientX: number; clientY: number }) {
    if (!drag.current.active) return;
    const dx = event.clientX - drag.current.x;
    const dy = event.clientY - drag.current.y;
    drag.current.x = event.clientX;
    drag.current.y = event.clientY;
    setYaw((prev) => prev + dx * 0.01);
    setPitch((prev) => Math.max(-1.2, Math.min(1.2, prev - dy * 0.01)));
  }

  function onPointerUp(event: {
    currentTarget: { hasPointerCapture: (id: number) => boolean; releasePointerCapture: (id: number) => void };
    pointerId: number;
  }) {
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
    drag.current.active = false;
    setDragging(false);
  }

  return (
    <Stack gap={8}>
      <Row gap={8} align="center" wrap>
        <Text size="small" tone="secondary">
          Cube rolls and faces
        </Text>
        <Button onClick={() => recordR(MAT_RY)}>Record Roll yaw 90</Button>
        <Button onClick={() => recordR(MAT_RX)}>Record Roll pitch 90</Button>
        <Button
          onClick={() => {
            setIndexExtra((prev) => [
              ...prev,
              {
                kind: "FF",
                label: "FF",
                afterStep: 0,
              },
            ]);
          }}
        >
          Record FF
        </Button>
        <Button
          onClick={() => {
            setYaw(0.62);
            setPitch(0.48);
          }}
        >
          Reset view
        </Button>
      </Row>
      {!showStack ? (
        <Row gap={8} wrap>
          <Text size="small" tone="secondary">
            Show shape
          </Text>
          <Pill
            active={nestTetra === "both"}
            onClick={() => setNestTetra("both")}
          >
            Both tetras
          </Pill>
          <Pill active={nestTetra === "R0"} onClick={() => setNestTetra("R0")}>
            R0 0–3
          </Pill>
          <Pill active={nestTetra === "R1"} onClick={() => setNestTetra("R1")}>
            R1 4–7
          </Pill>
        </Row>
      ) : null}
      <Text size="small">
        Drag is camera only (not a Roll). Record Roll snaps to one of{" "}
        {CUBE_POSES.length} cube rolls. A 90° turn that would not change
        roll is ignored.
      </Text>
      <Text size="small" weight="semibold">
        Index: roll {poseIndex + 1}/{CUBE_POSES.length}
        {indexExtra.length > 0
          ? ` · ${indexExtra.map((e) => e.label).join(" · ")}`
          : ""}
      </Text>
      <Callout tone="neutral" title="Roll and FF, not a tour">
        Record Roll writes one of 24 cube rolls (6 faces down × 4 turns), not
        an infinite turn string. Four yaw-90 hops return to roll 1/24. Record
        FF writes face flip of the current face (the face the hop already sits
        on). Camera drag does not write Roll. Closed tours live on the tours
        canvas.
      </Callout>
      <svg
        viewBox="0 0 540 300"
        width="100%"
        preserveAspectRatio="xMidYMid meet"
        role="img"
        aria-label="Rotatable cube with two tetrahedra. Drag to orbit."
        style={{
          touchAction: "none",
          cursor: dragging ? "grabbing" : "grab",
          userSelect: "none",
          maxHeight: 330,
        }}
        onPointerDown={onPointerDown}
        onPointerMove={onPointerMove}
        onPointerUp={onPointerUp}
        onPointerCancel={onPointerUp}
      >
        {showLattice
          ? [0, 0.5, 1].map((yy, i) => {
              const a = p(0, yy, 0);
              const b = p(1, yy, 0);
              const c = p(1, yy, 1);
              const d = p(0, yy, 1);
              const right = Math.max(a.x, b.x, c.x, d.x);
              const midY = (a.y + b.y + c.y + d.y) / 4;
              return (
                <g key={`L${i}`}>
                  <polygon
                    points={`${a.x},${a.y} ${b.x},${b.y} ${c.x},${c.y} ${d.x},${d.y}`}
                    fill={theme.fill.primary}
                    stroke={mute}
                    strokeWidth={1}
                  />
                  <text
                    x={right + 12}
                    y={midY + 4}
                    fill={label}
                    fontSize={11}
                  >{`L=${i}`}</text>
                </g>
              );
            })
          : null}

        {nestLayers.map((layer) => (
          <g key={layer.key}>{layer.node}</g>
        ))}

        {!showStack ? (
          <g transform="translate(20, 20)">
            <line
              x1={0}
              y1={0}
              x2={22}
              y2={0}
              stroke={r0}
              strokeWidth={2}
              opacity={nestTetra === "R1" ? 0.25 : 1}
            />
            <text
              x={30}
              y={4}
              fill={label}
              fontSize={11}
              fontWeight={500}
              opacity={nestTetra === "R1" ? 0.35 : 1}
            >
              R0 tetrahedron 0–3
            </text>
            <line
              x1={0}
              y1={18}
              x2={22}
              y2={18}
              stroke={r1}
              strokeWidth={2}
              opacity={nestTetra === "R0" ? 0.25 : 1}
            />
            <text
              x={30}
              y={22}
              fill={label}
              fontSize={11}
              fontWeight={500}
              opacity={nestTetra === "R0" ? 0.35 : 1}
            >
              R1 tetrahedron 4–7
            </text>
          </g>
        ) : null}

        {showLattice
          ? CUBE_EDGES.map(([a, b]) => {
              const ca = CUBE_CORNERS[a];
              const cb = CUBE_CORNERS[b];
              const m = p(
                (ca.x + cb.x) / 2,
                (ca.y + cb.y) / 2,
                (ca.z + cb.z) / 2,
              );
              return (
                <circle
                  key={`mid-${a}-${b}`}
                  cx={m.x}
                  cy={m.y}
                  r={4 * m.scale}
                  fill={fillMed}
                  stroke={mute}
                  opacity={depthAlpha(m.depth, 0.45, 1)}
                />
              );
            })
          : null}

        {showLattice
          ? [
              [0.5, 0.5, 0],
              [0.5, 0.5, 1],
              [0.5, 0, 0.5],
              [0.5, 1, 0.5],
              [0, 0.5, 0.5],
              [1, 0.5, 0.5],
            ].map((xyz, i) => {
              const q = p(xyz[0], xyz[1], xyz[2]);
              return (
                <circle
                  key={`face-${i}`}
                  cx={q.x}
                  cy={q.y}
                  r={6 * q.scale}
                  fill={theme.fill.secondary}
                  stroke={edgeStroke}
                  opacity={depthAlpha(q.depth, 0.45, 1)}
                />
              );
            })
          : null}

        {showLattice ? (
          <g>
            <rect
              x={p(0.5, 0.5, 0.5).x - 10}
              y={p(0.5, 0.5, 0.5).y - 10}
              width={20}
              height={20}
              transform={`rotate(45 ${p(0.5, 0.5, 0.5).x} ${p(0.5, 0.5, 0.5).y})`}
              fill={accent}
              stroke={accent}
            />
            <text
              x={p(0.5, 0.5, 0.5).x}
              y={p(0.5, 0.5, 0.5).y + 4}
              textAnchor="middle"
              fill={onAccent}
              fontSize={7}
              fontWeight={600}
            >
              M
            </text>
          </g>
        ) : null}

        {showStack
          ? CUBE_CORNERS.map((c) => {
              const size = 0.2;
              const inset = 0.05;
              const x0 = c.x === 0 ? inset : 1 - inset - size;
              const y0 = c.y === 0 ? inset : 1 - inset - size;
              const z0 = c.z === 0 ? inset : 1 - inset - size;
              const highlight = c.id === 0;
              const verts = [
                [0, 0, 0],
                [size, 0, 0],
                [size, size, 0],
                [0, size, 0],
                [0, 0, size],
                [size, 0, size],
                [size, size, size],
                [0, size, size],
              ].map(([x, y, z]) => p(x0 + x, y0 + y, z0 + z));
              const miniEdges: Array<[number, number]> = [
                [0, 1],
                [1, 2],
                [2, 3],
                [3, 0],
                [4, 5],
                [5, 6],
                [6, 7],
                [7, 4],
                [0, 4],
                [1, 5],
                [2, 6],
                [3, 7],
              ];
              return (
                <g key={`mini-${c.id}`}>
                  {miniEdges.map(([ia, ib]) => (
                    <line
                      key={`${c.id}-${ia}-${ib}`}
                      x1={verts[ia].x}
                      y1={verts[ia].y}
                      x2={verts[ib].x}
                      y2={verts[ib].y}
                      stroke={highlight ? accent : mute}
                      strokeWidth={highlight ? 2 : 1}
                    />
                  ))}
                </g>
              );
            })
          : null}

        {showStack
          ? CUBE_CORNERS.map((c) => {
              const q = p(c.x, c.y, c.z);
              return (
                <text
                  key={`lab-${c.id}`}
                  x={q.x}
                  y={q.y - 10}
                  textAnchor="middle"
                  fill={c.id === 0 ? accent : label}
                  fontSize={8}
                  fontWeight={c.id === 0 ? 600 : 400}
                >
                  {c.id === 0 ? "this 3×3×3" : c.id}
                </text>
              );
            })
          : null}
      </svg>
      <Text size="small" tone="tertiary">
        {caption}
      </Text>
    </Stack>
  );
}

export default function DSSCoreTopology() {
  const [nestStage, setNestStage] = useCanvasState<NestStage>(
    "nestStage",
    "corners",
  );
  const [hopFilter, setHopFilter] = useCanvasState<HopFilter>(
    "hopFilter",
    "all",
  );
  const [hopKind, setHopKind] = useCanvasState("hopKind", "1-2");
  const [hopTable, setHopTable] = useCanvasState<HopTableTab>(
    "hopTable",
    "ring",
  );

  const selectedCube = CUBE_WALKS.find((h) => hopKey(h.from, h.to) === hopKind);
  const selectedHop =
    WALK_HOPS.find((h) => hopKey(h.from, h.to) === hopKind) ??
    (selectedCube
      ? {
          from: selectedCube.from,
          to: selectedCube.to,
          via: "cube",
        }
      : WALK_HOPS[1]);
  const cubePills = hopFilter === "cube" || hopTable === "cube";

  function hopTone(from: number, to: number): "info" | undefined {
    if (hopKeyMatches(from, to, hopKind)) return "info";
    return undefined;
  }

  function selectHop(key: string) {
    setHopKind(key);
    const cube = CUBE_WALKS.find((h) => hopKey(h.from, h.to) === key);
    if (cube) {
      setHopTable("cube");
      setHopFilter("cube");
      return;
    }
    const walk = WALK_HOPS.find((h) => hopKey(h.from, h.to) === key);
    if (walk && hopFilter !== "all" && hopFilter !== walk.via) {
      setHopFilter(walk.via);
    }
    setHopTable("ring");
  }

  function selectFilter(next: HopFilter) {
    setHopFilter(next);
    if (next === "cube") setHopTable("cube");
    else setHopTable("ring");
    if (next === "Ring" || next === "M") {
      const current = WALK_HOPS.find((h) => hopKey(h.from, h.to) === hopKind);
      if (!current || current.via !== next) {
        const first = WALK_HOPS.find((h) => h.via === next);
        if (first) setHopKind(hopKey(first.from, first.to));
      }
    }
  }

  return (
    <Stack gap={20}>
      <Stack gap={8}>
        <H1>DSS hops</H1>
        <Text tone="secondary">
          This canvas is the hop set: Ring, M, and 12 cube edges on stops 0–7.
          Cycle A and Cycle B are finished laps that use some of these hops.
          They live on the closed-tours canvas, not here.
        </Text>
      </Stack>

      <Callout tone="info" title="Hops first">
        Allowed moves are here: Ring, M, and 12 cube edges. Closed words
        (Cycle A, Cycle B, Mixed) are a second canvas built on this set. Open
        DSS closed tours to step a lap.
      </Callout>

      <H2>Production rules</H2>
      <Table
        headers={["Rule", "What it does", "Allowed"]}
        rows={[
          [
            "Ring",
            "Continuation on one tetrahedron",
            "1→2, 3→0, 5→6, 7→4",
          ],
          [
            "Cube",
            "12 R0↔R1 wireframe edges",
            "0→4, 0→6, 1→4, 1→5, 1→7, 2→4, 2→6, 3→5, 3→6, 3→7, 5→0, 7→2",
          ],
          ["M", "Bind even to odd", "even → M → odd"],
        ]}
      />
      <Text size="small" tone="tertiary">
        Cycle A and Cycle B pick Ring, some cube edges, and M. They are not
        extra rules on this table.
      </Text>

      <H2>Walks on this cube</H2>
      <Row gap={20} wrap>
        <Stat value="12" label="Cube edges" />
        <Stat value="20" label="Hop kinds (4 Ring + 4 M + 12 cube)" />
      </Row>
      <Text tone="secondary">
        The same cube is the view for every filter. Ring sits on one
        tetrahedron (R0 stops 0–3, or R1 stops 4–7). M is the even-odd bind
        on that tetra edge (diamond on the hop, not at the cube center). Cube
        is the 12 rim edges. FF is face flip of the face that hop already sits
        on, not a stop. Twenty hop kinds total; the store is the path product,
        not 384 × 20. Click a hop pill or a line on the cube. The Tetra (8)
        table on the right is Ring plus M. Ring (4) on the left lights the two
        odd-to-even hops on the selected tetra; the other tetra stays ghost.
      </Text>
      <Grid columns="minmax(0, 1.15fr) minmax(0, 0.85fr)" gap={16}>
        <Stack gap={10}>
          <Row gap={8} wrap>
            <Pill
              active={hopFilter === "all"}
              onClick={() => selectFilter("all")}
            >
              All (8)
            </Pill>
            <Pill
              active={hopFilter === "Ring"}
              onClick={() => selectFilter("Ring")}
            >
              Ring (4)
            </Pill>
            <Pill
              active={hopFilter === "M"}
              onClick={() => selectFilter("M")}
            >
              M (4)
            </Pill>
            <Pill
              active={hopFilter === "cube"}
              onClick={() => selectFilter("cube")}
            >
              Cube (12)
            </Pill>
          </Row>
          <Row gap={6} wrap>
            {(cubePills ? CUBE_WALKS : WALK_HOPS.filter((h) => hopMatchesFilter(h.via, hopFilter))).map(
              (h) => {
                const key = hopKey(h.from, h.to);
                const tag = "kind" in h ? h.kind : h.via;
                return (
                  <Pill
                    key={key}
                    size="sm"
                    active={hopKind === key}
                    onClick={() => selectHop(key)}
                  >
                    {h.from}→{h.to} · {tag}
                  </Pill>
                );
              },
            )}
          </Row>
          <CubeHopView
            from={selectedHop.from}
            to={selectedHop.to}
            via={selectedHop.via}
            filter={hopFilter}
            onSelect={selectHop}
          />
          <Text size="small">
            {selectedHop.from} → {selectedHop.to}
            {selectedCube
              ? " · cube edge, R0 to R1."
              : selectedHop.via === "Ring"
                ? " · Ring stays on one tetrahedron."
                : " · M binds even to odd."}{" "}
            Multiply{" "}
            <Code>
              {selectedCube && selectedCube.kind === "cube"
                ? cubeEdgeMultiply(selectedCube.from, selectedCube.to)
                : hopMultiply(selectedHop.from, selectedHop.to, selectedHop.via)}
            </Code>
            .
          </Text>
        </Stack>
        <Stack gap={10}>
          <Row gap={8} wrap>
            <Pill
              active={hopTable === "ring"}
              onClick={() => setHopTable("ring")}
            >
              Tetra (8)
            </Pill>
            <Pill
              active={hopTable === "cube"}
              onClick={() => {
                setHopTable("cube");
                setHopFilter("cube");
              }}
            >
              Cube
            </Pill>
            <Pill
              active={hopTable === "primes"}
              onClick={() => setHopTable("primes")}
            >
              Primes
            </Pill>
            <Pill
              active={hopTable === "resolution"}
              onClick={() => setHopTable("resolution")}
            >
              Resolution
            </Pill>
          </Row>
          {hopTable === "ring" ? (
            <Table
              headers={["Hop", "Kind", "Multiply"]}
              rows={[
                [
                  <Pill size="sm" active={hopKind === "0-1"} onClick={() => selectHop("0-1")}>
                    0 → 1
                  </Pill>,
                  "M",
                  hopMultiply(0, 1, "M"),
                ],
                [
                  <Pill size="sm" active={hopKind === "1-2"} onClick={() => selectHop("1-2")}>
                    1 → 2
                  </Pill>,
                  "Ring",
                  hopMultiply(1, 2, "Ring"),
                ],
                [
                  <Pill size="sm" active={hopKind === "2-3"} onClick={() => selectHop("2-3")}>
                    2 → 3
                  </Pill>,
                  "M",
                  hopMultiply(2, 3, "M"),
                ],
                [
                  <Pill size="sm" active={hopKind === "3-0"} onClick={() => selectHop("3-0")}>
                    3 → 0
                  </Pill>,
                  "Ring",
                  hopMultiply(3, 0, "Ring"),
                ],
                [
                  <Pill size="sm" active={hopKind === "4-5"} onClick={() => selectHop("4-5")}>
                    4 → 5
                  </Pill>,
                  "M",
                  hopMultiply(4, 5, "M"),
                ],
                [
                  <Pill size="sm" active={hopKind === "5-6"} onClick={() => selectHop("5-6")}>
                    5 → 6
                  </Pill>,
                  "Ring",
                  hopMultiply(5, 6, "Ring"),
                ],
                [
                  <Pill size="sm" active={hopKind === "6-7"} onClick={() => selectHop("6-7")}>
                    6 → 7
                  </Pill>,
                  "M",
                  hopMultiply(6, 7, "M"),
                ],
                [
                  <Pill size="sm" active={hopKind === "7-4"} onClick={() => selectHop("7-4")}>
                    7 → 4
                  </Pill>,
                  "Ring",
                  hopMultiply(7, 4, "Ring"),
                ],
              ]}
              rowTone={[
                hopTone(0, 1),
                hopTone(1, 2),
                hopTone(2, 3),
                hopTone(3, 0),
                hopTone(4, 5),
                hopTone(5, 6),
                hopTone(6, 7),
                hopTone(7, 4),
              ]}
            />
          ) : hopTable === "cube" ? (
            <Table
              headers={["Hop", "Kind", "Multiply"]}
              rows={CUBE_WALKS.map((h) => [
                <Pill
                  size="sm"
                  active={hopKind === hopKey(h.from, h.to)}
                  onClick={() => selectHop(hopKey(h.from, h.to))}
                >
                  {h.from} → {h.to}
                </Pill>,
                "cube edge",
                cubeEdgeMultiply(h.from, h.to),
              ])}
              rowTone={CUBE_WALKS.map((h) => hopTone(h.from, h.to))}
            />
          ) : hopTable === "primes" ? (
            <Stack gap={10}>
              <Table
                headers={["Corner", "Prime"]}
                rows={CORNER_PRIME.map((p, i) => [
                  `${i}${i < 4 ? " R0" : " R1"}`,
                  String(p),
                ])}
                columnAlign={["left", "right"]}
              />
              <Table
                headers={["Bind", "Primes"]}
                rows={[
                  ["M even → odd", "137 × 139"],
                  ["Cube edge (12)", "corner × corner (same as Ring pair)"],
                  ["Roll", "6 faces down × 4 turns = 24 rolls"],
                  ["FF", "59^flip of the current face"],
                  ["S stack", "closed child is the cell payload"],
                ]}
              />
            </Stack>
          ) : (
            <Stack gap={10}>
              <Table
                headers={[
                  "Level",
                  "Stops × rolls",
                  "Walk kinds",
                  "With FF",
                  "One-cell spine",
                ]}
                rows={[
                  ["0 · ~10", "192", "4+4+12=20", "384", "384"],
                  ["1 · ~50", "192", "20", "384", "147,456"],
                  ["2 · ~100", "192", "20", "384", "56,623,104"],
                ]}
                columnAlign={["left", "right", "right", "right", "right"]}
              />
              <Stack gap={4}>
                <Text size="small">
                  <Code>8stp * 6fce * 4trn = 192</Code> stops, faces down,
                  90° turns
                </Text>
                <Text size="small">
                  <Code>8stp * 6fce * 4trn * 2ff = 384</Code> flick of that
                  face after the roll and turn
                </Text>
                <Text size="small">
                  <Code>384^n</Code> one-cell spine
                </Text>
                <Text size="small" tone="tertiary">
                  Walk kinds <Code>4 + 4 + 12 = 20</Code> stay off this
                  product. Store is the hop product along the path, not{" "}
                  <Code>384 × 20</Code>.
                </Text>
              </Stack>
            </Stack>
          )}
        </Stack>
      </Grid>

      <Callout tone="info" title="Everything on one cube">
        <Code>8stp * 6fce * 4trn * 2ff = 384</Code>. Each of the 6 faces can
        flick (0 or 1) after that face-down and turn. FF is not a walk. Ring
        and M stay on a tetrahedron. Cube walks are 12 rim edges.
        Record FF, Roll, and stack S on top. The pointer is the product of
        hops you took, not stops × kinds.
      </Callout>

      <Divider />

      <H2>Cube stack</H2>
      <Text tone="secondary">
        Two nested moves, not one. Filling 0–7 into a 3×3×3 is still one cube.
        Stacking is that cube becoming a node in the cube above.
      </Text>
      <Row gap={8} wrap>
        <Pill
          active={nestStage === "corners"}
          onClick={() => setNestStage("corners")}
        >
          1. Corners (8)
        </Pill>
        <Pill
          active={nestStage === "lattice"}
          onClick={() => setNestStage("lattice")}
        >
          2. One 3×3×3 (27)
        </Pill>
        <Pill
          active={nestStage === "stack"}
          onClick={() => setNestStage("stack")}
        >
          3. Stack layer
        </Pill>
      </Row>
      <CubeNest stage={nestStage} highlightHopKey={hopKind} />
      {nestStage === "corners" ? (
        <Callout tone="neutral" title="Two tetrahedra and 12 edge walks">
          Blue faces join 0–3 (R0). Warm faces join 4–7 (R1). Solid lines are
          the 12 cube edges (same as the Cube tab). Each runs R0 to R1, not a
          tetrahedron side. Matches the hop picker above when you select a cube
          edge.
        </Callout>
      ) : null}
      {nestStage === "lattice" ? (
        <Callout tone="neutral" title="Same cube, filled">
          Twelve edge midpoints are the same 12 cube-edge walks. Face centers
          and the body center (M) make 27. L=0, L=1, and L=2 are three slices
          through this cube, not a second geometry and not yet a stack.
        </Callout>
      ) : null}
      {nestStage === "stack" ? (
        <Callout tone="info" title="The stack is the next cube">
          A closed 0–7 (a finished 3×3×3) occupies one corner of the cube
          above. The other seven corners are the same pattern. You do not add
          state 8. You reuse 0–7 one layer up. Group sizes ~10, ~50, ~100,
          ~1000 are that stack spoken as people rather than as cubes.
        </Callout>
      ) : null}
      <Grid columns={4} gap={12}>
        <Stack gap={4}>
          <H3>~10</H3>
          <Text size="small">Intimate / direct group</Text>
          <Text size="small" tone="tertiary">
            Local cognition
          </Text>
        </Stack>
        <Stack gap={4}>
          <H3>~50</H3>
          <Text size="small">Coordinated group</Text>
          <Text size="small" tone="tertiary">
            Shared memory
          </Text>
        </Stack>
        <Stack gap={4}>
          <H3>~100</H3>
          <Text size="small">Structured / distributed</Text>
          <Text size="small" tone="tertiary">
            Governance
          </Text>
        </Stack>
        <Stack gap={4}>
          <H3>~1000</H3>
          <Text size="small">Network of groups</Text>
          <Text size="small" tone="tertiary">
            Inter-group protocol
          </Text>
        </Stack>
      </Grid>

      <H2>DSS as substrate</H2>
      <Table
        headers={["Grammar piece", "DSS reading"]}
        rows={[
          ["R0 continuation", "Vector substrate: associative, approximate retrieval"],
          ["R1 continuation", "Ledger substrate: persistent, governed state"],
          ["M mediation", "Witnessed bind between otherwise distinct regimes"],
          [
            "Cube edge",
            "R0↔R1 wireframe; 12 edges",
          ],
          ["Closed word", "Auditable completed tetra tour (Cycle A or B)"],
          ["PEO", "Portable evidence-bearing representation of that cycle"],
        ]}
      />

      <H2>Example nested index</H2>
      <Text tone="secondary">
        One pointer walked through corners, one surface, then three stacked
        cubes. Cycle A here is a closed word from the tours canvas, used as
        L0 payload. Corner primes are 2,3,5,7,11,13,17,19 for 0–7. Each M bind is
        137 × 139. A Roll is face-prime × 53^turn (which face is down, then
        90° in that plane). FF is 59^flip of that face. A stack step does not
        add letters: the closed child key is the payload of one parent cell.
      </Text>
      <H3>Level 0 · Cycle A · this cube (~10)</H3>
      <Table
        headers={["Hop", "From", "Constraint", "To", "Index fragment"]}
        rows={[
          ["t00", "0", "M", "1", "2 × 137 × 139 × 3"],
          ["t01", "1", "cube", "4", "3 × 11"],
          ["t02", "4", "M", "5", "11 × 137 × 139 × 13"],
          ["t03", "5", "Ring", "6", "13 × 17"],
          ["t04", "6", "M", "7", "17 × 137 × 139 × 19"],
          ["t05", "7", "cube", "2", "19 × 5"],
          ["t06", "2", "M", "3", "5 × 137 × 139 × 7"],
          ["t07", "3", "Ring", "0", "7 × 2"],
        ]}
      />
      <Text size="small" tone="tertiary">
        Closed L0 key is the product of those fragments (~104 bits). States
        visited: 0,1,4,5,6,7,2,3,0. M four times. Not 0,1,2,3,4,5,6,7.
      </Text>
      <H3>Cube edge · 0→4 after close</H3>
      <Table
        headers={["Hop", "From", "Constraint", "To", "Index fragment"]}
        rows={[
          ["t08", "0", "cube edge", "4", "2 × 11"],
        ]}
      />
      <Text size="small" tone="tertiary">
        One of the 12 cube-edge hops. Same corner primes as a Ring
        pair-multiply; not a tetrahedron continuation.
      </Text>
      <H3>Surface · face Fz- after close</H3>
      <Table
        headers={["Hop", "From", "Constraint", "To", "Index fragment"]}
        rows={[
          [
            "t09",
            "0",
            "Roll turn=1 · FF=0",
            "Fz- then 0",
            "2 × 41 × 53^1 × 59^0 × 2",
          ],
        ]}
      />
      <Text size="small" tone="tertiary">
        Six faces get 23, 29, 31, 37, 41, 43. Roll turn is 0..3; FF is 0 or 1.
        This step is extra completeness on the same cube, not a ninth corner.
      </Text>
      <H3>Level 1 · Cycle B · parent cells (~50)</H3>
      <Text>
        The sealed L0 walk occupies cell 0 of the cube above. Parent hops use
        the same Ring, cube, and M constraints; numbers 0–7 are now cells, each
        cell a finished child cube. A parent cell may also take a cube edge
        (same 0→4 style) between cells.
      </Text>
      <Table
        headers={["Hop", "Cell from", "Constraint", "Cell to", "Index fragment"]}
        rows={[
          ["t00", "0", "M", "1", "2 × 137 × 139 × 3"],
          ["t01", "1", "Ring", "2", "3 × 5"],
          ["t02", "2", "M", "3", "5 × 137 × 139 × 7"],
          ["t03", "3", "cube", "6", "7 × 17"],
          ["t04", "6", "M", "7", "17 × 137 × 139 × 19"],
          ["t05", "7", "Ring", "4", "19 × 11"],
          ["t06", "4", "M", "5", "11 × 137 × 139 × 13"],
          ["t07", "5", "cube", "0", "13 × 2"],
        ]}
      />
      <Text size="small" tone="tertiary">
        Level 1 index = L0 key × parent B key (~220 bits). Indexed transition
        types so far: M, Ring, cube edge, Roll, FF, and S (stack: child key
        rides inside the cell).
      </Text>
      <H3>Level 2 · Cycle A · (~100) and Level 3 · Cycle B · (~1000)</H3>
      <Table
        headers={["Level", "Cycle", "Where the child sits", "What gets indexed"]}
        rows={[
          [
            "2 · ~100",
            "A",
            "cell 4 (cube edge onto R1)",
            "Level 1 key × parent A hops; same eight constraints",
          ],
          [
            "3 · ~1000",
            "B",
            "cell 5",
            "Level 2 key × parent B hops; network of groups",
          ],
        ]}
      />
      <Callout tone="info" title="What the index is">
        Six transition kinds: Ring stay, cube edge, M even-to-odd bind, Roll
        (24 rolls), FF (face flip of the current face), S stack (closed
        child is one cell). The pointer is the product of fragments along that
        path. Roll is 53^turn. FF is 59^flip. Whole-cube camera drag is not
        part of the key.
      </Callout>

      <CollapsibleSection
        title="Engineering tests"
        count={9}
        defaultOpen={false}
      >
        <Stack gap={6}>
          <Text size="small">1. Formalise the grammar.</Text>
          <Text size="small">2. Generate all legal words up to N.</Text>
          <Text size="small">3. Enumerate closed cycles.</Text>
          <Text size="small">
            4. Test whether A and B are uniquely implied (up to start point).
          </Text>
          <Text size="small">
            5. Remove semantic labels and see what structure remains.
          </Text>
          <Text size="small">
            6. Reintroduce labels and test predictive power.
          </Text>
          <Text size="small">
            7. Test recursive application at different scales.
          </Text>
          <Text size="small">
            8. Compare generated transitions against DSS traces.
          </Text>
          <Text size="small">
            9. Show cube hops are not Ring.
          </Text>
          <Text size="small" tone="tertiary">
            Success: the topology generates useful structure without the
            semantic interpretation being hard-coded first.
          </Text>
        </Stack>
      </CollapsibleSection>
    </Stack>
  );
}
