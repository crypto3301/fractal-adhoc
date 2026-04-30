package graph

// NodeData — узел с числовым ID и буквенным иерархическим адресом.
// Адрес строится из символов a,b,c,d — по одному на каждый уровень.
// Пример: "abca" — копия a → подкопия b → подподкопия c → вершина a.
type NodeData struct {
	ID   int    `json:"id"`
	Addr string `json:"addr"`
}

// GraphData — полный граф, отправляемый в браузер.
// Adj и NormPos — срезы, индексированные числовым ID узла (0..N-1).
type GraphData struct {
	Order   int          `json:"order"`   // порядок предфрактала
	N       int          `json:"n"`       // реальное число узлов
	Nodes   []NodeData   `json:"nodes"`   // список узлов
	Adj     [][]int      `json:"adj"`     // adj[id] = список соседей
	NormPos [][2]float64 `json:"normPos"` // normPos[id] = [x, y] в [0.05, 0.95]
}
