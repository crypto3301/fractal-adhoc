package web

// Page — минималистичный чёрно-белый дизайн
const Page = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Ad-Hoc Предфрактальная Сеть</title>
    <style>
        :root {
            --bg: #000000;
            --panel: #0a0a0a;
            --text: #dddddd;
            --gray: #aaaaaa;
            --accent: #cccccc;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            display: flex;
            height: 100vh;
            overflow: hidden;
            background: var(--bg);
            color: var(--text);
            font-family: system-ui, -apple-system, sans-serif;
            font-size: 14px;
        }

        #sidebar {
            width: 280px;
            background: var(--panel);
            border-right: 1px solid #222222;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 28px;
            overflow-y: auto;
        }

        h1 {
            font-size: 18px;
            font-weight: 500;
            color: var(--accent);
            letter-spacing: 0.5px;
        }

        /* Минималистичные карточки статистики */
        .stats {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
        }

        .stat {
            padding: 12px 0;
        }

        .stat-value {
            font-size: 32px;
            font-weight: 600;
            color: white;
            line-height: 1;
        }

        .stat-label {
            font-size: 12px;
            color: var(--gray);
            margin-top: 4px;
        }

        .section {
            color: var(--gray);
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 10px;
        }

        button {
            width: 100%;
            padding: 13px 16px;
            margin-bottom: 8px;
            background: transparent;
            color: var(--text);
            border: 1px solid #333333;
            border-radius: 4px;
            font-size: 13.5px;
            cursor: pointer;
            text-align: left;
            transition: all 0.2s;
        }

        button:hover:not(:disabled) {
            border-color: var(--accent);
            color: white;
        }

        button.active {
            border-color: var(--accent);
            color: white;
            background: rgba(255,255,255,0.05);
        }

        button:disabled {
            opacity: 0.4;
            cursor: not-allowed;
        }

        .legend {
            display: flex;
            flex-direction: column;
            gap: 11px;
            font-size: 13.5px;
            color: var(--gray);
        }

        .legend-item {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .color-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            flex-shrink: 0;
        }

        #info {
            flex: 1;
            background: transparent;
            border: 1px solid #222222;
            border-radius: 4px;
            padding: 16px;
            font-size: 13.2px;
            line-height: 1.65;
            color: var(--gray);
            overflow-y: auto;
        }

        #canvas-container {
            flex: 1;
            position: relative;
            background: #000000;
        }

        canvas {
            display: block;
        }

        #status {
            position: absolute;
            top: 24px;
            left: 50%;
            transform: translateX(-50%);
            background: rgba(0,0,0,0.85);
            color: var(--accent);
            padding: 9px 22px;
            border-radius: 20px;
            font-size: 13px;
            border: 1px solid #222;
            white-space: nowrap;
            z-index: 10;
        }
    </style>
</head>
<body>
    <div id="sidebar">
        <h1>Ad-Hoc Предфрактальная Сеть</h1>

        <div class="stats">
            <div class="stat">
                <div class="stat-value" id="nodes">—</div>
                <div class="stat-label">Узлов</div>
            </div>
            <div class="stat">
                <div class="stat-value" id="edges">—</div>
                <div class="stat-label">Рёбер</div>
            </div>
            <div class="stat">
                <div class="stat-value" id="order">—</div>
                <div class="stat-label">Порядок</div>
            </div>
            <div class="stat">
                <div class="stat-value">K₄</div>
                <div class="stat-label">Затравка</div>
            </div>
        </div>

        <div>
            <div class="section">Этапы</div>
            <button id="btn1" onclick="startPhase(1)">1. Случайное размещение узлов</button>
            <button id="btn2" onclick="startPhase(2)" disabled>2. Построение Ad-Hoc сети</button>
            <button id="btn3" onclick="startPhase(3)" disabled>3. Нормализованный вид</button>
            <button onclick="resetAll()">Сбросить</button>
        </div>

        <div>
            <div class="section">Информация об узле</div>
            <div id="info">Наведите курсор на узел для просмотра информации</div>
        </div>
    </div>

    <div id="canvas-container">
        <div id="status">Загрузка графа...</div>
        <canvas id="canvas"></canvas>
    </div>

    <script>
        // ... (JS остаётся прежним, только минимальные правки)
        const canvas = document.getElementById('canvas');
        const ctx = canvas.getContext('2d');
        let W = 0, H = 0;
        let graph = null;
        let positions = [];
        let targetPositions = [];
        let nodeColors = [];
        let phase = 0;
        let edgeProgress = 0;
        let hoveredNode = -1;

        const QUAD_COLORS = ['#cccccc', '#bbbbbb', '#aaaaaa', '#999999'];

        function resize() {
            const cont = canvas.parentElement;
            W = cont.clientWidth;
            H = cont.clientHeight;
            canvas.width = W;
            canvas.height = H;
        }

        function toScreen(x, y) {
            const pad = 80;
            return [pad + x*(W - 2*pad), pad + y*(H - 2*pad)];
        }

        function getNodeColor(id) {
            return QUAD_COLORS[Math.floor(id / (graph.nodes.length / 4)) % 4];
        }

        function draw() {
            ctx.clearRect(0, 0, W, H);
            if (phase >= 3) drawGrid();
            if (phase >= 2) drawEdges();
            drawNodes();
        }

        function drawGrid() {
            ctx.strokeStyle = '#1a1a1a';
            ctx.lineWidth = 1;
            const s = 4;
            for (let i = 0; i <= s; i++) {
                const x = 80 + (i/s)*(W-160);
                ctx.beginPath(); ctx.moveTo(x,80); ctx.lineTo(x,H-80); ctx.stroke();
                const y = 80 + (i/s)*(H-160);
                ctx.beginPath(); ctx.moveTo(80,y); ctx.lineTo(W-80,y); ctx.stroke();
            }
        }

        function drawEdges() {
            const n = Math.floor(edgeProgress);
            ctx.lineCap = 'round';
            for (let i = 0; i < n; i++) {
                const e = graph.edges[i];
                const a = positions[e.from];
                const b = positions[e.to];
                ctx.strokeStyle = nodeColors[e.from];
                ctx.globalAlpha = 0.25 + e.level * 0.2;
                ctx.lineWidth = 1.2;
                ctx.beginPath();
                ctx.moveTo(a[0], a[1]);
                ctx.lineTo(b[0], b[1]);
                ctx.stroke();
            }
            ctx.globalAlpha = 1;
        }

        function drawNodes() {
            const r = phase >= 3 ? 5.5 : 7.5;
            for (let i = 0; i < graph.nodes.length; i++) {
                const p = positions[i];
                const col = nodeColors[i];

                if (i === hoveredNode) {
                    ctx.shadowBlur = 25;
                    ctx.shadowColor = '#ffffff';
                    ctx.fillStyle = '#ffffff';
                    ctx.beginPath(); ctx.arc(p[0], p[1], r+7, 0, Math.PI*2); ctx.fill();
                }

                ctx.shadowBlur = 10;
                ctx.shadowColor = col;
                ctx.fillStyle = col;
                ctx.beginPath(); ctx.arc(p[0], p[1], r, 0, Math.PI*2); ctx.fill();

                ctx.shadowBlur = 0;
                ctx.fillStyle = '#000000';
                ctx.font = 'bold 8px monospace';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillText(i.toString(), p[0], p[1]);
            }
        }

        function startPhase(n) {
            if (n === 1) {
                phase = 1;
                positions = graph.nodes.map(() => [80 + Math.random()*(W-160), 80 + Math.random()*(H-160)]);
                updateButtons();
                document.getElementById('status').textContent = "Фаза 1 — случайное размещение";
            } else if (n === 2) {
                phase = 2;
                edgeProgress = 0;
                updateButtons();
                document.getElementById('status').textContent = "Фаза 2 — построение сети...";

                const total = graph.edges.length;
                let last = Date.now();

                function anim() {
                    const now = Date.now();
                    edgeProgress += (now - last) / 12;
                    last = now;
                    if (edgeProgress < total) requestAnimationFrame(anim);
                    else {
                        edgeProgress = total;
                        document.getElementById('status').textContent = "Фаза 2 завершена";
                    }
                    draw();
                }
                anim();
            } else if (n === 3) {
                phase = 3;
                updateButtons();
                document.getElementById('status').textContent = "Фаза 3 — нормализация...";

                targetPositions = graph.fractalPos.map(p => toScreen(p[0], p[1]));
                const start = positions.map(p => [p[0], p[1]]);
                let t = 0;

                function move() {
                    t = Math.min(t + 0.028, 1);
                    const e = t < 0.5 ? 2*t*t : 1 - Math.pow(-2*t + 2, 2)/2;
                    for (let i = 0; i < positions.length; i++) {
                        positions[i][0] = start[i][0] + (targetPositions[i][0] - start[i][0]) * e;
                        positions[i][1] = start[i][1] + (targetPositions[i][1] - start[i][1]) * e;
                    }
                    draw();
                    if (t < 1) requestAnimationFrame(move);
                    else document.getElementById('status').textContent = "Фаза 3 завершена";
                }
                move();
            }
        }

        function updateButtons() {
            document.getElementById('btn1').disabled = phase >= 1;
            document.getElementById('btn2').disabled = phase < 1 || phase >= 2;
            document.getElementById('btn3').disabled = phase < 2 || phase >= 3;

            document.getElementById('btn1').classList.toggle('active', phase >= 1);
            document.getElementById('btn2').classList.toggle('active', phase >= 2);
            document.getElementById('btn3').classList.toggle('active', phase >= 3);
        }

        function resetAll() {
            phase = 0;
            edgeProgress = 0;
            positions = [];
            hoveredNode = -1;
            ctx.clearRect(0, 0, W, H);
            updateButtons();
            document.getElementById('status').textContent = "Готов";
            document.getElementById('info').innerHTML = "Наведите курсор на узел";
        }

        function findHovered(mx, my) {
            let best = -1, minD = Infinity;
            for (let i = 0; i < positions.length; i++) {
                const p = positions[i];
                const d = (p[0]-mx)**2 + (p[1]-my)**2;
                if (d < minD) { minD = d; best = i; }
            }
            return minD < 200 ? best : -1;
        }

        canvas.addEventListener('mousemove', function(e) {
            if (!graph || positions.length === 0) return;
            const rect = canvas.getBoundingClientRect();
            const mx = e.clientX - rect.left;
            const my = e.clientY - rect.top;
            const newH = findHovered(mx, my);

            if (newH !== hoveredNode) {
                hoveredNode = newH;
                if (hoveredNode >= 0) showNodeInfo(hoveredNode);
                else document.getElementById('info').innerHTML = "Наведите курсор на узел";
                draw();
            }
        });

        function showNodeInfo(id) {
            const node = graph.nodes[id];
            const fp = graph.fractalPos[id];
            const levels = node.label.split('.').map((l,i) => 
                "<div>Уровень " + (i+1) + ": <code>" + l + "</code></div>"
            ).join('');

            document.getElementById('info').innerHTML = 
                "<b>Узел #" + node.id + "</b><br>" +
                "Индекс: <code>" + node.label + "</code><br>" +
                "Позиция: <code>(" + fp[0].toFixed(3) + ", " + fp[1].toFixed(3) + ")</code><br><br>" +
                "<b>Иерархия:</b><br>" + levels;
        }

        function init() {
            resize();
            window.addEventListener('resize', () => { resize(); if (phase >= 1) draw(); });

            fetch('/graph')
                .then(r => r.json())
                .then(data => {
                    graph = data;
                    nodeColors = graph.nodes.map((_,i) => getNodeColor(i));

                    document.getElementById('nodes').textContent = graph.nodes.length;
                    document.getElementById('edges').textContent = graph.edges.length;
                    document.getElementById('order').textContent = graph.order;

                    document.getElementById('status').textContent = "Граф загружен";
                    resetAll();
                    draw();
                })
                .catch(() => document.getElementById('status').textContent = "Ошибка загрузки");
        }

        window.onload = init;
    </script>
</body>
</html>`
