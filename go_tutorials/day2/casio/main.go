package main

import (
	"fmt"

	// ĐÂY LÀ PHẦN QUAN TRỌNG:
	// "mycalculator" là tên module (trong go.mod)
	// "/mathutil" là tên thư mục package
	"casio/math"
)

func main() {
	soA := 10
	soB := 5

	// Dùng package bằng cách gọi TênPackage.TênHàm
	tong := math.Add(soA, soB)

	fmt.Printf("Tổng của %d và %d là: %d\n", soA, soB, tong)
}
