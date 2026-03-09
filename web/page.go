package web

// Page — встроенная HTML-страница с визуализацией предфрактальной Ad-Hoc сети.
const Page = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Ad-Hoc Предфрактальная Сеть</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{display:flex;height:100vh;overflow:hidden;background:#0d1117;color:#e6edf3;font-family:'Segoe UI',system-ui,sans-serif;font-size:13px}
#side{width:280px;min-width:280px;background:#161b22;border-right:1px solid #30363d;padding:16px;display:flex;flex-direction:column;gap:12px;overflow-y:auto}
h1{font-size:14px;font-weight:700;color:#4ecca3;padding-bottom:10px;border-bottom:1px solid #30363d;letter-spacing:.3px}
.stats{display:grid;grid-template-columns:1fr 1fr;gap:6px}
.stat{background:#0d1117;border-radius:6px;padding:8px;text-align:center;border:1px solid #21262d}
.sv{font-size:22px;font-weight:700;color:#4ecca3}
.sl{font-size:10px;color:#8b949e;margin-top:2px;text-transform:uppercase;letter-spacing:.6px}
.sec{font-size:10px;color:#8b949e;text-transform:uppercase;letter-spacing:.9px;margin-top:2px}
button{width:100%;padding:9px 12px;border:1px solid #30363d;background:#21262d;color:#e6edf3;border-radius:6px;cursor:pointer;text-align:left;transition:all .15s;font-size:12px;margin-bottom:3px}
button:hover:not(:disabled){background:#1f6feb22;border-color:#388bfd;color:#58a6ff}
button.on{background:#4ecca322;border-color:#4ecca3;color:#4ecca3}
button:disabled{opacity:.3;cursor:not-allowed}
#info{background:#0d1117;border:1px solid #30363d;border-radius:6px;padding:12px;font-size:12px;line-height:1.75;flex:1;color:#8b949e;min-height:120px}
#info b{color:#e6edf3;display:block;margin-bottom:4px}
#info code{background:#21262d;padding:1px 6px;border-radius:3px;color:#f0883e;font-size:11px;font-family:monospace}
.leg{display:flex;flex-direction:column;gap:5px}
.li{display:flex;align-items:center;gap:8px;font-size:11px;color:#8b949e}
.ld{width:10px;height:10px;border-radius:50%;flex-shrink:0;box-shadow:0 0 4px currentColor}
#wrap{flex:1;position:relative;overflow:hidden}
canvas{display:block}
#lbl{position:absolute;top:14px;left:50%;transform:translateX(-50%);background:rgba(13,17,23,.92);border:1px solid #30363d;color:#4ecca3;padding:5px 16px;border-radius:20px;font-size:12px;pointer-events:none;white-space:nowrap;z-index:10;box-shadow:0 2px 12px #0008}
</style>
</head>
<body>
<div id="side">
  <h1>&#9672; Предфрактальная Ad-Hoc Сеть</h1>

  <div class="stats">
    <div class="stat"><div class="sv" id="sn">—</div><div class="sl">Узлов</div></div>
    <div class="stat"><div class="sv" id="se">—</div><div class="sl">Рёбер</div></div>
    <div class="stat"><div class="sv" id="so">—</div><div class="sl">Порядок</div></div>
    <div class="stat"><div class="sv">K₄</div><div class="sl">Затравка</div></div>
  </div>

  <div class="sec">Этапы построения</div>
  <button id="b1" disabled onclick="doPhase1()">▶ 1. Случайное размещение узлов</button>
  <button id="b2" disabled onclick="doPhase2()">▶ 2. Объединение в Ad-Hoc сеть</button>
  <button id="b3" disabled onclick="doPhase3()">▶ 3. Нормированный вид топологии</button>
  <button onclick="doReset()">↺ Сбросить</button>

  <div class="sec">Копии верхнего уровня</div>
  <div class="leg">
    <div class="li"><div class="ld" style="background:#4ecca3;color:#4ecca3"></div>Копия 0 — верх-лево</div>
    <div class="li"><div class="ld" style="background:#e94560;color:#e94560"></div>Копия 1 — верх-право</div>
    <div class="li"><div class="ld" style="background:#f5a623;color:#f5a623"></div>Копия 2 — низ-право</div>
    <div class="li"><div class="ld" style="background:#8b8fea;color:#8b8fea"></div>Копия 3 — низ-лево</div>
  </div>

  <div class="sec">Информация об узле</div>
  <div id="info">Наведите курсор на узел для просмотра его фрактального индекса и иерархии</div>
</div>

<div id="wrap">
  <span id="lbl">Загрузка данных графа...</span>
  <canvas id="c"></canvas>
</div>

<script>
// ─── Globals ───────────────────────────────────────────────────────────────
var C   = document.getElementById('c');
var ctx = C.getContext('2d');
var W = 0, H = 0;
var g      = null;   // GraphData from server
var pos    = [];     // current screen positions [[x,y],...]
var colors = [];     // color per node
var phase  = 0;      // 0=idle 1=random 2=network 3=normalized
var edgeCnt = 0;     // animated edge counter
var hov    = -1;     // hovered node index
var raf    = null;   // animation frame handle
var PAD    = 60;

var COLS = ['#4ecca3','#e94560','#f5a623','#8b8fea'];

// ─── Resize ────────────────────────────────────────────────────────────────
function resize() {
    var p = C.parentElement;
    W = p.clientWidth; H = p.clientHeight;
    C.width = W; C.height = H;
}

// ─── Coordinate helpers ────────────────────────────────────────────────────
function toCanvas(fx, fy) {
    return [PAD + fx*(W-2*PAD), PAD + fy*(H-2*PAD)];
}

function fractalPositions() {
    return g.fractalPos.map(function(p){ return toCanvas(p[0], p[1]); });
}

function randPositions() {
    return g.nodes.map(function(){
        return [PAD + Math.random()*(W-2*PAD), PAD + Math.random()*(H-2*PAD)];
    });
}

// ─── Draw ──────────────────────────────────────────────────────────────────
function draw() {
    ctx.clearRect(0, 0, W, H);

    if (phase >= 3) drawGrid();
    if (phase >= 2) drawEdges();
    drawNodes();
}

function drawGrid() {
    ctx.save();
    ctx.strokeStyle = '#21262d';
    ctx.lineWidth = 1;
    var divs = 4;
    for (var x = 0; x <= divs; x++) {
        var cx = PAD + (x/divs)*(W-2*PAD);
        ctx.beginPath(); ctx.moveTo(cx, PAD); ctx.lineTo(cx, H-PAD); ctx.stroke();
    }
    for (var y = 0; y <= divs; y++) {
        var cy = PAD + (y/divs)*(H-2*PAD);
        ctx.beginPath(); ctx.moveTo(PAD, cy); ctx.lineTo(W-PAD, cy); ctx.stroke();
    }
    ctx.restore();
}

function drawEdges() {
    var n = Math.min(Math.floor(edgeCnt), g.edges.length);
    ctx.save();
    for (var i = 0; i < n; i++) {
        var e  = g.edges[i];
        var lv = e.level;
        ctx.globalAlpha  = 0.10 + 0.13 * lv;
        ctx.lineWidth    = 0.3  + 0.35 * lv;
        ctx.strokeStyle  = colors[e.from];
        var p0 = pos[e.from], p1 = pos[e.to];
        ctx.beginPath();
        ctx.moveTo(p0[0], p0[1]);
        ctx.lineTo(p1[0], p1[1]);
        ctx.stroke();
    }
    ctx.restore();
}

function drawNodes() {
    var r = (phase >= 3) ? 6 : 8;
    for (var i = 0; i < g.nodes.length; i++) {
        var p = pos[i];
        var c = colors[i];
        ctx.save();

        // Highlight hovered node
        if (i === hov) {
            ctx.shadowBlur  = 22;
            ctx.shadowColor = c;
            ctx.fillStyle   = '#ffffff';
            ctx.beginPath();
            ctx.arc(p[0], p[1], r+5, 0, Math.PI*2);
            ctx.fill();
        }

        ctx.shadowBlur  = (phase >= 2) ? 8 : 0;
        ctx.shadowColor = c;
        ctx.fillStyle   = c;
        ctx.beginPath();
        ctx.arc(p[0], p[1], r, 0, Math.PI*2);
        ctx.fill();

        // Node ID label (small, inside dot)
        ctx.fillStyle     = '#0d1117';
        ctx.font          = 'bold 7px monospace';
        ctx.textAlign     = 'center';
        ctx.textBaseline  = 'middle';
        ctx.shadowBlur    = 0;
        ctx.fillText(i, p[0], p[1]);

        ctx.restore();
    }
}

// ─── Render loop ───────────────────────────────────────────────────────────
function startLoop() {
    if (raf) cancelAnimationFrame(raf);
    function loop() { draw(); raf = requestAnimationFrame(loop); }
    loop();
}

// ─── Phase handlers ────────────────────────────────────────────────────────
function doReset() {
    if (raf) { cancelAnimationFrame(raf); raf = null; }
    phase = 0; edgeCnt = 0; pos = []; hov = -1;
    ctx.clearRect(0, 0, W, H);
    setBtn('b1', false, false);
    setBtn('b2', true,  false);
    setBtn('b3', true,  false);
    setLbl('Нажмите «1. Случайное размещение узлов»');
    setInfo('Граф сброшен. Начните с шага 1.');
}

function doPhase1() {
    phase = 1; edgeCnt = 0;
    pos = randPositions();
    setBtn('b1', false, true);
    setBtn('b2', false, false);
    setLbl('Шаг 1: Случайное размещение узлов (' + g.nodes.length + ' узлов)');
    startLoop();
}

function doPhase2() {
    phase = 2; edgeCnt = 0;
    setBtn('b2', false, true);
    setBtn('b3', false, false);
    var total = g.edges.length;
    setLbl('Шаг 2: Формирование предфрактальной топологии...');

    function addEdges() {
        edgeCnt = Math.min(edgeCnt + Math.max(1, total / 120), total);
        if (edgeCnt < total) {
            requestAnimationFrame(addEdges);
        } else {
            setLbl('Шаг 2: Ad-Hoc сеть сформирована. Рёбер: ' + total);
        }
    }
    requestAnimationFrame(addEdges);
}

function doPhase3() {
    phase = 3;
    var target = fractalPositions();
    var src    = pos.map(function(p){ return p.slice(); });
    setBtn('b3', false, true);
    setLbl('Шаг 3: Нормированный вид предфрактальной топологии');

    var t = 0;
    function step() {
        t = Math.min(t + 0.018, 1);
        var e = t < 0.5 ? 2*t*t : -1 + (4-2*t)*t; // ease in-out
        for (var i = 0; i < pos.length; i++) {
            pos[i][0] = src[i][0] + (target[i][0] - src[i][0]) * e;
            pos[i][1] = src[i][1] + (target[i][1] - src[i][1]) * e;
        }
        if (t < 1) requestAnimationFrame(step);
    }
    requestAnimationFrame(step);
}

// ─── Hover / info ──────────────────────────────────────────────────────────
C.addEventListener('mousemove', function(e) {
    if (!g || !pos.length) return;
    var rect = C.getBoundingClientRect();
    var mx = e.clientX - rect.left;
    var my = e.clientY - rect.top;
    var best = -1, bd = 900;
    for (var i = 0; i < pos.length; i++) {
        var dx = pos[i][0]-mx, dy = pos[i][1]-my;
        var d = dx*dx + dy*dy;
        if (d < bd) { bd = d; best = i; }
    }
    if (bd < 400 && best >= 0) {
        if (hov !== best) {
            hov = best;
            showNodeInfo(best);
        }
    } else {
        hov = -1;
    }
});

function showNodeInfo(id) {
    var nd = g.nodes[id];
    var fp = g.fractalPos[id];
    var levels = nd.label.split('.').map(function(l, i){
        return '<br>&nbsp;&nbsp;' + (i+1) + '. Копия&nbsp;<code>' + l + '</code>';
    }).join('');
    setInfo(
        '<b>Узел #' + nd.id + '</b>' +
        'Фрактальный индекс: <code>' + nd.label + '</code><br>' +
        'Нормир. позиция: <code>(' + fp[0].toFixed(3) + ',&nbsp;' + fp[1].toFixed(3) + ')</code><br>' +
        'Иерархия:' + levels
    );
}

// ─── UI helpers ────────────────────────────────────────────────────────────
function setBtn(id, disabled, active) {
    var b = document.getElementById(id);
    b.disabled = disabled;
    b.classList.toggle('on', active);
}

function setLbl(txt) {
    document.getElementById('lbl').textContent = txt;
}

function setInfo(html) {
    document.getElementById('info').innerHTML = html;
}

function nodeColor(node) {
    return COLS[Math.floor(node.id / (g.nodes.length / 4)) % 4];
}

// ─── Init ──────────────────────────────────────────────────────────────────
window.addEventListener('resize', function(){ resize(); draw(); });
resize();

fetch('/graph')
    .then(function(r){ return r.json(); })
    .then(function(data){
        g      = data;
        colors = g.nodes.map(nodeColor);

        document.getElementById('sn').textContent = g.nodes.length;
        document.getElementById('se').textContent = g.edges.length;
        document.getElementById('so').textContent = g.order;

        setLbl('Нажмите «1. Случайное размещение узлов»');
        setBtn('b1', false, false);
        setInfo(
            'Граф загружен: K₄ порядка ' + g.order +
            ', ' + g.nodes.length + ' узлов, ' + g.edges.length + ' рёбер.'
        );
        draw();
    });
</script>
</body>
</html>`
