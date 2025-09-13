package main

import (
	"fmt"
	"math"
)

/*
Разработать программу нахождения расстояния между двумя точками на плоскости.
Точки представлены в виде структуры Point с инкапсулированными (приватными) полями x, y (типа float64) и конструктором.
Расстояние рассчитывается по формуле между координатами двух точек.

Подсказка: используйте функцию-конструктор NewPoint(x, y), Point и метод Distance(other Point) float64.
*/

// Point - structure for storing a point coordinates in 2-dimensional space
type Point struct {
	x float64
	y float64
}

// Distance - calculates distance from the 1st point to the 2nd one
func (firstPoint Point) Distance(secondPoint Point) float64 {
	return math.Sqrt(math.Pow(secondPoint.x-firstPoint.x, 2) + math.Pow(secondPoint.y-firstPoint.y, 2))
}

// NewPoint - returns new Point-structure with provided coordinates
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func main() {
	// проверим работу на пифагоровой тройке 3-4-5
	point0 := NewPoint(0, 0)
	point1 := NewPoint(0, 3)
	point2 := NewPoint(4, 0)
	fmt.Println("Расстояние от 0-вой точки до (0;3):", point0.Distance(point1))
	fmt.Println("Расстояние от 0-вой точки до (4;0):", point0.Distance(point2))
	fmt.Println("Расстояние от (4;0) до (0;3) (ожидаем 5):", point1.Distance(point2))
}
