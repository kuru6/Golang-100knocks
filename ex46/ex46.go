package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("人数を入力してください:")
	scanner.Scan()
	people_t := scanner.Text()
	people, _ := strconv.Atoi(people_t)
	if people <= 4 {
		fmt.Printf("合計入場料は%d円です", people*600)
	} else if people >= 5 && people <= 19 {
		fmt.Printf("合計入場料は%d円です", people*550)
	} else {
		fmt.Printf("合計入場料は%d円です", people*500)
	}
}
