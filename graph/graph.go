package graph

// Node — узел Ad-Hoc сети.
type Node struct {
	ID    int    `json:"id"`
	Label string `json:"label"` // иерархическая метка, напр. "2.1.0"
}

// Edge — ребро между двумя узлами.
type Edge struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Level int `json:"level"` // уровень иерархии: 1=базовый K4, ORDER=внешний
}

// GraphData — полный граф, передаваемый в браузер.
type GraphData struct {
	Order      int          `json:"order"`
	Nodes      []Node       `json:"nodes"`
	Edges      []Edge       `json:"edges"`
	FractalPos [][2]float64 `json:"fractalPos"` // нормированные позиции [0,1]²
}
