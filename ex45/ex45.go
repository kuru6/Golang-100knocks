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
	distancetext := scanner.Text()
	distance, _ := strconv.Atoi(distancetext)
	ryoukin := 0
	if distance <= 0 {
		ryoukin = 0
	} else if distance <= 1700 {
		ryoukin = 610
	} else {
		distance = (distance - 1700)
		var calc int
		if distance%313 == 0 {
			calc = distance / 313
		} else {
			calc = distance/313 + 1
		}
		ryoukin = 610 + (calc * 80)
	}
	fmt.Println(ryoukin)
}
