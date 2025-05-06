package main

import "fmt"

func main() {
	var number int
	fmt.Println("数値を入力してください")
	fmt.Scan(&number)
	fmt.Printf("入力された数値*3=%d", number*3)
}
