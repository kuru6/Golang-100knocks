package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	texta := scanner.Text()
	numbera, _ := strconv.Atoi(texta)

	scanner.Scan()
	textb := scanner.Text()
	numberb, _ := strconv.Atoi(textb)

	numberc := numbera / numberb
	fmt.Println(numberc)
	fmt.Println(numberc * numberb)
}
