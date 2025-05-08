package main

import "fmt"

func main() {
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 8; j++ {
			fmt.Print("とんで")
		}
		for k := 0; k <= 2; k++ {
			fmt.Print("まわって")
		}
		fmt.Print("まわる\n")
	}
}
