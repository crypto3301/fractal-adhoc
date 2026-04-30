package web

const Page = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Ad-Hoc Предфрактальная Сеть</title>
<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;600&family=IBM+Plex+Sans:wght@300;400;600&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --black:  #0a0a0a;
  --white:  #ffffff;
  --grey1:  #f4f4f4;
  --grey2:  #e0e0e0;
  --grey3:  #a0a0a0;
  --grey4:  #606060;
  --mono:   'IBM Plex Mono', monospace;
  --sans:   'IBM Plex Sans', sans-serif;
}

html, body { height: 100%; overflow: hidden; background: var(--white); color: var(--black); font-family: var(--sans); font-size: 13px; }
body { display: flex; }

/* ── Sidebar ── */
#side {
  width: 280px; min-width: 280px; height: 100vh;
  border-right: 1px solid var(--grey2);
  display: flex; flex-direction: column;
  background: var(--white);
  overflow-y: auto;
}

/* ── Canvas area ── */
#wrap { flex: 1; position: relative; background: var(--grey1); overflow: hidden; }
canvas { display: block; cursor: crosshair; }

/* ── Sidebar blocks ── */
.side-header { padding: 24px 24px 18px; border-bottom: 1px solid var(--grey2); }
.side-header h1 { font-size: 14px; font-weight: 600; color: var(--black); letter-spacing: -.01em; }
.side-header p  { font-size: 10px; color: var(--grey3); margin-top: 4px; font-family: var(--mono); letter-spacing: .06em; text-transform: uppercase; }

.side-section { padding: 16px 24px; border-bottom: 1px solid var(--grey2); display: flex; flex-direction: column; gap: 10px; }
.sec-label { font-family: var(--mono); font-size: 10px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: var(--grey3); }

/* ── Stats ── */
.stats { display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: var(--grey2); border: 1px solid var(--grey2); }
.stat { background: var(--white); padding: 10px 12px; text-align: center; }
.sv { font-family: var(--mono); font-size: 20px; font-weight: 600; color: var(--black); line-height: 1; }
.sl { font-size: 9px; color: var(--grey3); margin-top: 3px; text-transform: uppercase; letter-spacing: .08em; }

/* ── Tabs ── */
.tabs { display: flex; border: 1px solid var(--black); }
.tab {
  flex: 1; padding: 8px 4px; text-align: center;
  font-family: var(--mono); font-size: 10px; font-weight: 600;
  letter-spacing: .06em; text-transform: uppercase;
  cursor: pointer; border: none; background: var(--white); color: var(--grey4);
  border-right: 1px solid var(--black); transition: background .12s, color .12s;
}
.tab:last-child { border-right: none; }
.tab:hover { background: var(--grey1); color: var(--black); }
.tab.on { background: var(--black); color: var(--white); }

/* ── Inputs ── */
.row { display: flex; gap: 6px; align-items: center; }
.row label { font-size: 11px; color: var(--grey4); white-space: nowrap; }

input[type=number] {
  flex: 1; padding: 8px 10px;
  background: var(--grey1); border: 1px solid var(--grey2); border-radius: 0;
  color: var(--black); font-family: var(--mono); font-size: 12px;
  outline: none; transition: border-color .15s; min-width: 0;
}
input[type=number]:focus { border-color: var(--black); }

select {
  flex: 1; padding: 8px 28px 8px 10px;
  background: var(--grey1)
    url("data:image/svg+xml,%3Csvg width='10' height='6' viewBox='0 0 10 6' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M0 0l5 6 5-6z' fill='%23606060'/%3E%3C/svg%3E")
    no-repeat right 10px center;
  border: 1px solid var(--grey2); border-radius: 0;
  color: var(--black); font-family: var(--mono); font-size: 11px;
  outline: none; cursor: pointer; min-width: 0;
  appearance: none; -webkit-appearance: none;
}
select:focus { border-color: var(--black); }

/* ── Buttons ── */
.btn {
  padding: 8px 14px; border: 1px solid var(--black);
  background: var(--black); color: var(--white);
  font-family: var(--mono); font-size: 10px; font-weight: 600;
  letter-spacing: .06em; text-transform: uppercase;
  cursor: pointer; transition: background .12s, color .12s; border-radius: 0;
  white-space: nowrap;
}
.btn:hover { background: var(--white); color: var(--black); }
.btn-full  { width: 100%; }
.btn-ghost { background: var(--white); color: var(--black); }
.btn-ghost:hover { background: var(--grey1); }
.btn-danger { border-color: var(--black); }

/* ── Legend ── */
.leg { display: flex; flex-direction: column; gap: 5px; }
.li  { display: flex; align-items: center; gap: 10px; font-size: 11px; color: var(--grey4); }
.ld  { width: 8px; height: 8px; flex-shrink: 0; }

/* ── Info panel ── */
#info {
  flex: 1; padding: 16px 24px;
  font-size: 12px; line-height: 1.85; color: var(--grey4);
  overflow-y: auto;
}
#info b    { color: var(--black); display: block; margin-bottom: 3px; }
#info code {
  background: var(--grey1); border: 1px solid var(--grey2);
  padding: 1px 6px; font-family: var(--mono); font-size: 11px; color: var(--black);
}

/* ── Stats bar (bottom) ── */
.stats-bar { display: grid; grid-template-columns: 1fr 1fr; border-top: 1px solid var(--grey2); }
.stat-b    { padding: 12px 24px; border-right: 1px solid var(--grey2); }
.stat-b:last-child { border-right: none; }
.stat-bv   { font-family: var(--mono); font-size: 18px; font-weight: 600; line-height: 1; }
.stat-bl   { font-size: 9px; color: var(--grey3); text-transform: uppercase; letter-spacing: .08em; margin-top: 3px; }

/* ── Canvas label ── */
#lbl {
  position: absolute; top: 18px; left: 50%; transform: translateX(-50%);
  background: var(--white); color: var(--black);
  padding: 6px 18px; font-family: var(--mono); font-size: 11px;
  border: 1px solid var(--grey2); white-space: nowrap;
  pointer-events: none; letter-spacing: .04em;
}
</style>
</head>
<body>

<div id="side">
  <div class="side-header">
    <h1>Ad-Hoc Предфрактальная Сеть</h1>
    <p>Затравка K4 &nbsp;·&nbsp; Адресация a, b, c, d</p>
  </div>

  <!-- Статистика -->
  <div class="side-section">
    <div class="sec-label">Статистика</div>
    <div class="stats">
      <div class="stat"><div class="sv" id="sn">—</div><div class="sl">Узлов</div></div>
      <div class="stat"><div class="sv" id="se">—</div><div class="sl">Рёбер</div></div>
      <div class="stat"><div class="sv" id="so">—</div><div class="sl">Порядок</div></div>
      <div class="stat"><div class="sv">K₄</div><div class="sl">Затравка</div></div>
    </div>
  </div>

  <!-- Режим -->
  <div class="side-section">
    <div class="sec-label">Режим отображения</div>
    <div class="tabs">
      <button class="tab on" id="tab-nodes" onclick="switchMode('nodes')">Узлы</button>
      <button class="tab"    id="tab-adhoc" onclick="switchMode('adhoc')">Ad-Hoc</button>
      <button class="tab"    id="tab-norm"  onclick="switchMode('norm')">Норм.</button>
    </div>
  </div>

  <!-- Параметры -->
  <div class="side-section">
    <div class="sec-label">Параметры</div>
    <div class="row">
      <label>Узлов</label>
      <input type="number" id="nodeCount" value="64" min="64">
    </div>
    <button class="btn btn-full" onclick="changeNodeCount()">Сгенерировать</button>
  </div>

  <!-- Удаление -->
  <div class="side-section">
    <div class="sec-label">Управление узлами</div>
    <select id="removeSelect"><option value="">— выберите узел —</option></select>
    <div class="row">
      <button class="btn btn-full btn-danger" onclick="removeSelectedNode()">Удалить узел</button>
      <button class="btn btn-ghost" onclick="doReset()">Сброс</button>
    </div>
  </div>

  <!-- Легенда -->
  <div class="side-section">
    <div class="sec-label">Иерархия копий</div>
    <div class="leg">
      <div class="li"><div class="ld" style="background:#222"></div>Копия a — верх-лево</div>
      <div class="li"><div class="ld" style="background:#555"></div>Копия b — верх-право</div>
      <div class="li"><div class="ld" style="background:#888"></div>Копия c — низ-лево</div>
      <div class="li"><div class="ld" style="background:#bbb"></div>Копия d — низ-право</div>
    </div>
  </div>

  <!-- Информация -->
  <div id="info">Загрузка графа...</div>

  <!-- Нижняя плашка -->
  <div class="stats-bar">
    <div class="stat-b"><div class="stat-bv" id="bn">—</div><div class="stat-bl">Активных узлов</div></div>
    <div class="stat-b"><div class="stat-bv" id="be">—</div><div class="stat-bl">Активных рёбер</div></div>
  </div>
</div>

<div id="wrap">
  <span id="lbl">Подключение...</span>
  <canvas id="c"></canvas>
</div>

<script>
// ── Globals ─────────────────────────────────────────────────────
var C   = document.getElementById('c');
var ctx = C.getContext('2d');

var graphData   = null;
var activeGraph = null;
var randomPos   = {};
var mode        = 'nodes';
var selected    = null;

// Оттенки серого по первому символу адреса — монохромная схема
var ADDR_COLORS = {a:'#111111', b:'#444444', c:'#777777', d:'#aaaaaa'};

var BACKUP_MAP = {a:['d','b','c'], b:['c','a','d'], c:['b','a','d'], d:['a','b','c']};

// ── Загрузка данных ──────────────────────────────────────────────
function loadGraph(data) {
  var nodes = [], adj = {}, addr = {}, normPos = {};
  data.nodes.forEach(function(n) {
    nodes.push(n.id);
    adj[n.id]     = (data.adj[n.id]     || []).slice();
    addr[n.id]    = n.addr;
    normPos[n.id] = (data.normPos[n.id] || [0.5, 0.5]).slice();
  });
  graphData = {N: data.n, order: data.order, nodes: nodes, adj: adj, addr: addr, normPos: normPos};
  initRandomPositions();
  doReset();
  document.getElementById('nodeCount').value = data.n;
  document.getElementById('so').textContent  = data.order;
  setLbl('Граф загружен — K4 порядка ' + data.order + ', ' + data.n + ' узлов');
}

// ── Позиции ──────────────────────────────────────────────────────
function initRandomPositions() {
  randomPos = {};
  graphData.nodes.forEach(function(n) {
    randomPos[n] = [0.05 + Math.random()*0.90, 0.05 + Math.random()*0.90];
  });
}

function getPos(id) {
  if (mode === 'norm') return activeGraph.normPos[id] || [0.5, 0.5];
  return randomPos[id] || [0.5, 0.5];
}

// ── Resize ───────────────────────────────────────────────────────
function resizeCanvas() {
  var wrap = document.getElementById('wrap');
  C.width  = wrap.clientWidth;
  C.height = wrap.clientHeight;
}

// ── Draw ─────────────────────────────────────────────────────────
function draw() {
  if (!activeGraph) return;
  ctx.clearRect(0, 0, C.width, C.height);

  if (mode === 'adhoc' || mode === 'norm') {
    activeGraph.nodes.forEach(function(i) {
      var p1 = getPos(i);
      var x1 = p1[0] * C.width;
      var y1 = p1[1] * C.height;
      (activeGraph.adj[i] || []).forEach(function(j) {
        if (j <= i) return;
        var p2  = getPos(j);
        var x2  = p2[0] * C.width;
        var y2  = p2[1] * C.height;
        var sel = selected !== null && (i === selected || j === selected);
        ctx.strokeStyle = sel ? '#0a0a0a' : '#c8c8c8';
        ctx.lineWidth   = sel ? 2 : 1;
        ctx.globalAlpha = sel ? 1 : 0.8;
        ctx.beginPath();
        ctx.moveTo(x1, y1);
        if (mode === 'adhoc') {
          var mx = (x1+x2)/2 + Math.sin(i+j) * 24;
          var my = (y1+y2)/2 + Math.cos(i+j) * 24;
          ctx.quadraticCurveTo(mx, my, x2, y2);
        } else {
          ctx.lineTo(x2, y2);
        }
        ctx.stroke();
      });
    });
    ctx.globalAlpha = 1;
  }

  activeGraph.nodes.forEach(function(i) {
    var p = getPos(i);
    var x = p[0] * C.width;
    var y = p[1] * C.height;

    var isSel  = (i === selected);
    var isNeib = selected !== null && (activeGraph.adj[selected] || []).indexOf(i) >= 0;
    var addrStr    = activeGraph.addr[i] || 'a';
    var accentColor = ADDR_COLORS[addrStr[0]] || '#333';

    // Узел: белый круг с чёрной обводкой
    // Выбранный: инвертирован (чёрный круг с белым ID)
    if (isSel) { ctx.shadowBlur = 12; ctx.shadowColor = 'rgba(0,0,0,0.25)'; }

    ctx.fillStyle   = isSel ? '#0a0a0a' : '#ffffff';
    ctx.strokeStyle = isSel ? '#0a0a0a' : (isNeib ? '#0a0a0a' : accentColor);
    ctx.lineWidth   = isSel ? 2 : (isNeib ? 2 : 1.5);
    ctx.beginPath();
    ctx.arc(x, y, 13, 0, Math.PI*2);
    ctx.fill();
    ctx.stroke();
    ctx.shadowBlur = 0;

    ctx.fillStyle    = isSel ? '#ffffff' : '#0a0a0a';
    ctx.font         = 'bold 10px "IBM Plex Mono", monospace';
    ctx.textAlign    = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(i.toString(), x, y);

    // Адрес над выбранным узлом
    if (isSel) {
      ctx.fillStyle = '#0a0a0a';
      ctx.font      = '10px "IBM Plex Mono", monospace';
      ctx.fillText(addrStr, x, y - 22);
    }
  });
}

// ── Find node ────────────────────────────────────────────────────
function findNode(mx, my) {
  for (var k = 0; k < activeGraph.nodes.length; k++) {
    var i = activeGraph.nodes[k];
    var p = getPos(i);
    var dx = p[0]*C.width  - mx;
    var dy = p[1]*C.height - my;
    if (Math.sqrt(dx*dx + dy*dy) < 18) return i;
  }
  return null;
}

// ── Режим ────────────────────────────────────────────────────────
function switchMode(m) {
  mode = m;
  ['nodes','adhoc','norm'].forEach(function(id) {
    document.getElementById('tab-'+id).classList.toggle('on', m === id);
  });
  var labels = {
    nodes: 'Режим: случайное размещение узлов',
    adhoc: 'Режим: Ad-Hoc сеть',
    norm:  'Режим: нормированный вид топологии'
  };
  setLbl(labels[m]);
  draw();
}

// ── Подсчёт рёбер ────────────────────────────────────────────────
function countEdges(adj) {
  var t = 0;
  Object.keys(adj).forEach(function(k) { t += (adj[k]||[]).length; });
  return Math.round(t / 2);
}

// ── Обновление UI ────────────────────────────────────────────────
function updateInfo() {
  if (!activeGraph) return;
  var n = activeGraph.nodes.length;
  var e = countEdges(activeGraph.adj);
  document.getElementById('sn').textContent = n;
  document.getElementById('se').textContent = e;
  document.getElementById('bn').textContent = n;
  document.getElementById('be').textContent = e;

  if (selected !== null && activeGraph.nodes.indexOf(selected) >= 0) {
    var addr   = activeGraph.addr[selected] || '';
    var degree = (activeGraph.adj[selected] || []).length;
    var parts  = addr.split('').map(function(ch, i) {
      return 'Уровень&nbsp;' + (i+1) + ':&nbsp;копия&nbsp;<code>' + ch + '</code>';
    }).join('<br>');
    setInfo(
      '<b>Узел #' + selected + '</b>' +
      'Адрес: <code>' + addr + '</code><br>' +
      'Степень: <code>' + degree + '</code><br><br>' +
      parts
    );
  } else {
    setInfo(
      'Узлов:&nbsp;<code>' + n + '</code><br>' +
      'Рёбер:&nbsp;<code>' + e + '</code><br><br>' +
      '<span style="font-size:11px;color:#a0a0a0">Кликните на узел — увидите его адрес и иерархию</span>'
    );
  }
}

function updateRemoveSelect() {
  var sel = document.getElementById('removeSelect');
  sel.innerHTML = '<option value="">— выберите узел —</option>';
  if (!activeGraph) return;
  activeGraph.nodes.forEach(function(n) {
    var opt = document.createElement('option');
    opt.value       = n;
    opt.textContent = 'Узел ' + n + ' (' + (activeGraph.addr[n]||'') + ')';
    sel.appendChild(opt);
  });
}

// ── Сброс ────────────────────────────────────────────────────────
function doReset() {
  activeGraph = {
    N:       graphData.N,
    order:   graphData.order,
    nodes:   graphData.nodes.slice(),
    adj:     JSON.parse(JSON.stringify(graphData.adj)),
    addr:    Object.assign({}, graphData.addr),
    normPos: Object.assign({}, graphData.normPos)
  };
  selected = null;
  updateRemoveSelect();
  updateInfo();
  document.getElementById('so').textContent = graphData.order;
  draw();
}

// ── Генерация ────────────────────────────────────────────────────
function changeNodeCount() {
  var n = parseInt(document.getElementById('nodeCount').value);
  if (isNaN(n) || n < 64) {
    document.getElementById('nodeCount').value = 64;
    alert('Минимум 64 узла!');
    return;
  }
  setLbl('Генерация графа на ' + n + ' узлов...');
  fetch('/rebuild', {
    method: 'POST', headers: {'Content-Type':'application/json'},
    body: JSON.stringify({n: n})
  })
  .then(function(r) { return r.json(); })
  .then(function()  { return fetch('/graph').then(function(r) { return r.json(); }); })
  .then(loadGraph)
  .catch(function(e) { setLbl('Ошибка: ' + e); });
}

// ── Удаление узла с перестройкой ─────────────────────────────────
function removeSelectedNode() {
  var sel = document.getElementById('removeSelect');
  var val = sel.value;
  if (!val) { alert('Выберите узел!'); return; }
  var nodeId = parseInt(val);
  if (activeGraph.nodes.length <= 1) { alert('Нельзя удалить последний узел!'); return; }
  activeGraph = removeNodeAndRepair(activeGraph, nodeId);
  if (selected === nodeId) selected = null;
  updateRemoveSelect();
  updateInfo();
  setLbl('Узел ' + nodeId + ' удалён. Осталось: ' + activeGraph.nodes.length);
  draw();
}

function getBestBackup(clusterPrefix, brokenPort, aliveNodes, addr) {
  var orig = aliveNodes.filter(function(n) { return (addr[n]||'').slice(-1) === brokenPort; });
  if (orig.length) return orig[0];
  var bps = BACKUP_MAP[brokenPort] || [];
  for (var i = 0; i < bps.length; i++) {
    var cands = aliveNodes.filter(function(n) { return (addr[n]||'').slice(-1) === bps[i]; });
    if (cands.length) return cands[0];
  }
  return aliveNodes[0] || null;
}

function removeNodeAndRepair(g, nodeToRemove) {
  var nodes = g.nodes.filter(function(n) { return n !== nodeToRemove; });
  var nodeSet = {};
  nodes.forEach(function(n) { nodeSet[n] = true; });

  var adj = {};
  nodes.forEach(function(n) {
    adj[n] = (g.adj[n]||[]).filter(function(j) { return j !== nodeToRemove; });
  });

  var addr    = Object.assign({}, g.addr);
  var normPos = Object.assign({}, g.normPos);

  var externalEdges = [];
  g.nodes.forEach(function(u) {
    (g.adj[u]||[]).forEach(function(v) {
      if (v <= u) return;
      var addrU = g.addr[u]||'', addrV = g.addr[v]||'';
      if (addrU.slice(0,-1) !== addrV.slice(0,-1)) {
        var cl = 0;
        while (cl < addrU.length && cl < addrV.length && addrU[cl] === addrV[cl]) cl++;
        externalEdges.push({
          u:u, v:v,
          clusterU: addrU.slice(0, cl+1), clusterV: addrV.slice(0, cl+1),
          portU: addrU.slice(-1),         portV: addrV.slice(-1)
        });
      }
    });
  });

  externalEdges.forEach(function(edge) {
    if (!nodeSet[edge.u] || !nodeSet[edge.v]) {
      var aU = nodes.filter(function(n) { return (addr[n]||'').startsWith(edge.clusterU); });
      var aV = nodes.filter(function(n) { return (addr[n]||'').startsWith(edge.clusterV); });
      if (!aU.length || !aV.length) return;
      var nU = getBestBackup(edge.clusterU, edge.portU, aU, addr);
      var nV = getBestBackup(edge.clusterV, edge.portV, aV, addr);
      if (nU !== null && nV !== null && nU !== nV) {
        adj[nU] = adj[nU]||[]; adj[nV] = adj[nV]||[];
        if (adj[nU].indexOf(nV) < 0) { adj[nU].push(nV); adj[nV].push(nU); }
      }
    }
  });

  return {N:nodes.length, order:g.order, nodes:nodes, adj:adj, addr:addr, normPos:normPos};
}

// ── Helpers ──────────────────────────────────────────────────────
function setLbl(txt)  { document.getElementById('lbl').textContent  = txt; }
function setInfo(html){ document.getElementById('info').innerHTML    = html; }

// ── Events ───────────────────────────────────────────────────────
C.addEventListener('click', function(e) {
  var rect = C.getBoundingClientRect();
  selected = findNode(e.clientX - rect.left, e.clientY - rect.top);
  updateInfo();
  draw();
});

window.addEventListener('resize', function() { resizeCanvas(); draw(); });

// ── Init ─────────────────────────────────────────────────────────
resizeCanvas();
fetch('/graph')
  .then(function(r) { return r.json(); })
  .then(loadGraph)
  .catch(function(e) { setLbl('Ошибка: ' + e); setInfo('Не удалось загрузить граф'); });
</script>
</body>
</html>`
