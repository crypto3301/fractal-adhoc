package graph

// buildEdges рекурсивно строит рёбра предфрактального графа G(order)
// с базовым смещением узлов base.
//
// Правило построения:
//   G(1) = K4 — полносвязный четырёхугольник (6 рёбер).
//   G(n) = 4 копии G(n-1), соединённых по шаблону K4:
//     для каждого ребра (i,j) в K4:
//       крайний-j узел копии i  ↔  крайний-i узел копии j.
//
// Возвращает четыре «угловых» (граничных) узла данного подграфа [c0,c1,c2,c3].
func buildEdges(order, base int, edges *[]Edge) [4]int {
	if order == 1 {
		// Базовый случай: K4 — C(4,2)=6 рёбер
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				*edges = append(*edges, Edge{base + i, base + j, 1})
			}
		}
		// Угловые узлы — сами вершины K4
		return [4]int{base, base + 1, base + 2, base + 3}
	}

	sz := Pow4(order - 1)

	// Рекурсивно строим 4 подграфа и собираем их угловые узлы
	var sub [4][4]int
	for i := 0; i < 4; i++ {
		sub[i] = buildEdges(order-1, base+i*sz, edges)
	}

	// Межкопийные рёбра уровня order: для каждого ребра K4 (i,j)
	// соединяем угловой-j узел копии i с угловым-i узлом копии j
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			*edges = append(*edges, Edge{sub[i][j], sub[j][i], order})
		}
	}

	// Угловые узлы данного подграфа: угловой-i копии i
	var ext [4]int
	for i := 0; i < 4; i++ {
		ext[i] = sub[i][i]
	}
	return ext
}

// Build собирает полный GraphData предфрактального графа K4 заданного порядка.
func Build(order int) GraphData {
	n := Pow4(order)

	nodes := make([]Node, n)
	for i := range nodes {
		nodes[i] = Node{
			ID:    i,
			Label: FractalLabel(i, order),
		}
	}

	var edges []Edge
	buildEdges(order, 0, &edges)

	fp := make([][2]float64, n)
	for i := range nodes {
		fp[i] = FractalPos(i, order)
	}

	return GraphData{
		Order:      order,
		Nodes:      nodes,
		Edges:      edges,
		FractalPos: fp,
	}
}
