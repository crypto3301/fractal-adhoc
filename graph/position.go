package graph

import (
	"fmt"
	"math"
)

// Координаты углов K4 в единичном квадрате:
// 0=верх-лево, 1=верх-право, 2=низ-право, 3=низ-лево
var cornerX = [4]float64{0, 1, 1, 0}
var cornerY = [4]float64{0, 0, 1, 1}

// Pow4 возвращает 4^n.
func Pow4(n int) int {
	r := 1
	for i := 0; i < n; i++ {
		r *= 4
	}
	return r
}

// FractalLabel строит иерархическую метку узла.
// Пример для order=3: "2.1.0" означает копия 2 → подкопия 1 → узел 0 в K4.
func FractalLabel(nodeID, order int) string {
	id := nodeID
	s := ""
	for l := 0; l < order; l++ {
		sz := Pow4(order - 1 - l)
		idx := id / sz
		id %= sz
		if l > 0 {
			s += "."
		}
		s += fmt.Sprint(idx)
	}
	return s
}

// FractalPos вычисляет нормированную позицию узла в [0,1]².
//
// На каждом уровне иерархии l выбирается подквадрант по индексу копии.
// Масштаб смещения убывает как 0.5^l, благодаря чему граничные (интерфейсные)
// узлы соседних копий получают одинаковые координаты — фрактальное свойство.
func FractalPos(nodeID, order int) [2]float64 {
	x, y := 0.0, 0.0
	id := nodeID
	for l := 0; l < order; l++ {
		sz := Pow4(order - 1 - l)
		idx := (id / sz) % 4
		exp := l + 1
		if exp >= order {
			exp = order - 1
		}
		sc := math.Pow(0.5, float64(exp))
		x += cornerX[idx] * sc
		y += cornerY[idx] * sc
	}
	return [2]float64{x, y}
}
