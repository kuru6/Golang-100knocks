package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	j := 0
	k := 0
	for i := 0; i <= 4; {
		if j+k == 5 {
			break
		}
		fmt.Println("あなたの手を入力してください[グー=0][チョキ=1][パー=2]")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		human, _ := strconv.Atoi(scanner.Text())
		computer := rand.Intn(3)
		fmt.Printf("コンピューターは%dを出しました\n", computer)
		if human == computer {
			fmt.Println("あいこです")
		} else if human == 0 && computer == 1 || human == 1 && computer == 2 || human == 2 && computer == 0 {
			fmt.Println("人間の勝ちです")
			j++
			fmt.Printf("%d勝,%d敗", j, k)

		} else {
			fmt.Println("コンピューターの勝ちです")
			k++
			fmt.Printf("%d勝,%d敗", j, k)
		}
	}
	if j > k {
		fmt.Println("人間の勝利")
	} else {
		fmt.Println("コンピューターの勝利")
	}
}
