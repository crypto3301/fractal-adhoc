package graph

import "math"

// baseCoords задаёт координаты четырёх символов в единичном квадрате.
// Соответствие: a=верх-лево, b=верх-право, c=низ-лево, d=низ-право.
var baseCoords = map[byte][2]float64{
	'a': {0, 1},
	'b': {1, 1},
	'c': {0, 0},
	'd': {1, 0},
}

// calcRawPos вычисляет нормированную позицию узла по его буквенному адресу.
//
// Каждый символ вносит вклад с убывающим масштабом: 1, 0.5, 0.25, ...
// Итоговые координаты нормируются в [0,1]² и инвертируются по Y
// чтобы 'a' (верх-лево) попал в левый верхний угол экрана.
func calcRawPos(addr string) [2]float64 {
	var x, y float64
	for i := 0; i < len(addr); i++ {
		scale := 1.0 / math.Pow(2, float64(i))
		bc := baseCoords[addr[i]]
		x += bc[0] * scale
		y += bc[1] * scale
	}
	maxVal := 2.0 - math.Pow(2, 1-float64(len(addr)))
	return [2]float64{x / maxVal, 1 - y/maxVal}
}

// normalizePositions масштабирует позиции в диапазон [margin, 1-margin].
// Это гарантирует отступ от краёв canvas.
func normalizePositions(raw [][2]float64, n int, margin float64) [][2]float64 {
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)

	for i := 0; i < n; i++ {
		if raw[i][0] < minX {
			minX = raw[i][0]
		}
		if raw[i][0] > maxX {
			maxX = raw[i][0]
		}
		if raw[i][1] < minY {
			minY = raw[i][1]
		}
		if raw[i][1] > maxY {
			maxY = raw[i][1]
		}
	}

	span := 1.0 - 2*margin
	out := make([][2]float64, n)
	for i := 0; i < n; i++ {
		var rx, ry float64
		if maxX > minX {
			rx = (raw[i][0] - minX) / (maxX - minX)
		} else {
			rx = 0.5
		}
		if maxY > minY {
			ry = (raw[i][1] - minY) / (maxY - minY)
		} else {
			ry = 0.5
		}
		out[i] = [2]float64{margin + rx*span, margin + ry*span}
	}
	return out
}
