package main

import "math/rand"

func main() {
	m := rand.Intn(3) + 1
	if m == 1 {
		fmt.Prtint("ダイヤ")
	} else if m == 2 {
		fmt.Prtint("ハート")
	} else if m == 3 {
		fmt.Prtint("スペード")
	} else if m == 4 {
		fmt.Prtint("クローバー")
	}
	n := rand.Intn(12) + 1
	fmt.Println(n)
}
