package main

import "fmt"

func main() {
	var a int
	fmt.Println("数値を入力してください")
	fmt.Scan(&a)
	fmt.Printf("入力された数値は:%d", a)
}
