package graph

import "math"

// Четыре символа — индексы соответствуют порядку [0,1,2,3].
var chars = []byte{'a', 'b', 'c', 'd'}

// charIdx — обратная таблица: символ → индекс.
var charIdx = map[byte]int{'a': 0, 'b': 1, 'c': 2, 'd': 3}

// graphState — внутреннее представление графа во время построения.
type graphState struct {
	nodes []int
	adj   map[int][]int
	addr  map[int]string
}

// buildSeed создаёт затравку G(1) = K4:
// 4 вершины с адресами a,b,c,d, соединённые все со всеми.
func buildSeed() graphState {
	return graphState{
		nodes: []int{0, 1, 2, 3},
		adj: map[int][]int{
			0: {1, 2, 3},
			1: {0, 2, 3},
			2: {0, 1, 3},
			3: {0, 1, 2},
		},
		addr: map[int]string{
			0: "a", 1: "b", 2: "c", 3: "d",
		},
	}
}

// expandGraph выполняет одну итерацию расширения: G(n) → G(n+1).
//
// Правило: каждая старая вершина заменяется четырьмя новыми (K4).
// Адрес новой вершины = адрес_старой + символ_позиции.
// Межкластерные рёбра: для каждого старого ребра (u,v)
//
//	порт u = индекс последнего символа адреса v (и наоборот).
func expandGraph(g graphState) graphState {
	newAdj := map[int][]int{}
	newAddr := map[int]string{}
	var newNodes []int
	nextID := 0

	// nodeMap[oldID][localIdx] = newID
	nodeMap := map[int][4]int{}

	// Создаём 4 новых узла для каждого старого
	for _, vOld := range g.nodes {
		baseAddr := g.addr[vOld]
		var local [4]int
		for li := 0; li < 4; li++ {
			vNew := nextID
			nextID++
			newNodes = append(newNodes, vNew)
			newAddr[vNew] = baseAddr + string(chars[li])
			newAdj[vNew] = []int{}
			local[li] = vNew
		}
		nodeMap[vOld] = local

		// Внутренний K4: соединяем все 6 пар
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				u, v := local[i], local[j]
				newAdj[u] = append(newAdj[u], v)
				newAdj[v] = append(newAdj[v], u)
			}
		}
	}

	// Межкластерные рёбра
	for _, uOld := range g.nodes {
		for _, vOld := range g.adj[uOld] {
			if vOld <= uOld {
				continue // каждое ребро один раз
			}
			// Порт u = индекс последнего символа адреса v (и наоборот)
			uLast := g.addr[uOld][len(g.addr[uOld])-1]
			vLast := g.addr[vOld][len(g.addr[vOld])-1]
			idxU := charIdx[vLast] // порт копии u смотрит в сторону v
			idxV := charIdx[uLast] // порт копии v смотрит в сторону u
			uGate := nodeMap[uOld][idxU]
			vGate := nodeMap[vOld][idxV]
			newAdj[uGate] = append(newAdj[uGate], vGate)
			newAdj[vGate] = append(newAdj[vGate], uGate)
		}
	}

	return graphState{nodes: newNodes, adj: newAdj, addr: newAddr}
}

// Build строит предфрактальный граф с не менее чем N узлами.
//
// Порядок определяется автоматически: минимальный k такой что 4^k >= N.
// Если 4^k > N, берём первые N узлов и фильтруем рёбра.
// Узлы нумеруются 0..N-1 в порядке создания (обход в глубину по копиям).
func Build(N int) GraphData {
	// Минимальный порядок
	order := 1
	for int(math.Pow(4, float64(order))) < N {
		order++
	}

	// Итеративное расширение
	g := buildSeed()
	for i := 1; i < order; i++ {
		g = expandGraph(g)
	}

	// Берём первые N узлов
	nodes := g.nodes
	if N < len(nodes) {
		nodes = nodes[:N]
	}
	actualN := len(nodes)

	// Множество активных узлов для фильтрации рёбер
	nodeSet := map[int]bool{}
	for _, n := range nodes {
		nodeSet[n] = true
	}

	// Строим adj и NormPos как срезы размером actualN
	// (узлы нумерованы 0..actualN-1, поэтому индекс = ID)
	adj := make([][]int, actualN)
	rawPos := make([][2]float64, actualN)
	nodeList := make([]NodeData, actualN)

	for _, n := range nodes {
		// Фильтруем соседей — оставляем только тех кто в nodeSet
		filtered := []int{}
		for _, j := range g.adj[n] {
			if nodeSet[j] {
				filtered = append(filtered, j)
			}
		}
		adj[n] = filtered
		rawPos[n] = calcRawPos(g.addr[n])
		nodeList[n] = NodeData{ID: n, Addr: g.addr[n]}
	}

	normPos := normalizePositions(rawPos, actualN, 0.05)

	return GraphData{
		Order:   order,
		N:       actualN,
		Nodes:   nodeList,
		Adj:     adj,
		NormPos: normPos,
	}
}
