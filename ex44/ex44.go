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
	textchangemoney := scanner.Text()
	changemoney, _ := strconv.Atoi(textchangemoney)
	ryougaego := money / changemoney
	var cent float32
	cent = money / changemoney
	fmt.Printf("%d円は、ドルに両替すると%dドル%fセントになるよ", money, ryougaego, cent)

}
