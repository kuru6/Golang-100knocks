package main

import (
	"fmt"
)

func main() {
	var a, b int
	fmt.Print("数値を入力してください:")
	fmt.Scan(&a)
	fmt.Print("数値を入力してください:")
	fmt.Scan(&b)
	fmt.Println("2つの数字の和は", a+b)
	fmt.Println("2つの数字の差は", a-b)
	fmt.Println("2つの数字の積は", a*b)
	fmt.Printf("2つの数字の商は%dで余りは%d", a/b, a%b)
}
