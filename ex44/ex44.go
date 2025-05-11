package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("所持金を入力して")
	scanner.Scan()
	textmoney := scanner.Text()
	money, _ := strconv.Atoi(textmoney)
	fmt.Printf("レートを入力して")
	scanner.Scan()
	textrate := scanner.Text()
	rate, _ := strconv.Atoi(textrate)
	ryougaego := money / rate
	cent := (money * 100 / rate) % 100
	fmt.Println(1000 % 120)
	fmt.Printf("%d円は、ドルに両替すると%dドル%dセントになる", money, ryougaego, cent)
}
