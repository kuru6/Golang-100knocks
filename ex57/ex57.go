package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type subject struct {
	english  []int
	math     []int
	japanese []int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	people_num, _ := strconv.Atoi(scanner.Text())
	if people_num <= 0 {
		return
	}
	s := inputdata(people_num)
	outputdata(s, people_num)
}

func inputdata(people_num int) subject {
	var s subject
	scanner := bufio.NewScanner(os.Stdin)
	s.english = make([]int, people_num)
	s.math = make([]int, people_num)
	s.japanese = make([]int, people_num)

	for i := 0; i < people_num; i++ {
		scanner.Scan()
		people_data := scanner.Text()
		parts := strings.Split(people_data, " ")
		s.english[i], _ = strconv.Atoi(parts[0])
		s.math[i], _ = strconv.Atoi(parts[1])
		s.japanese[i], _ = strconv.Atoi(parts[2])
	}
	return s
}

func outputdata(s subject, people_num int) {
	fmt.Println("合計点数")
	var total_english_score int
	var total_math_score int
	var total_japanese_score int
	for i := 0; i < people_num; i++ {
		fmt.Printf("%d: %d\n", i+1, s.english[i]+s.math[i]+s.japanese[i])
	}
	for i := 0; i < people_num; i++ {
		total_english_score += s.english[i]
		total_math_score += s.math[i]
		total_japanese_score += s.japanese[i]
	}
	fmt.Printf("英語の平均点は%d\n", total_english_score/people_num)
	fmt.Printf("数学の平均点は%d\n", total_math_score/people_num)
	fmt.Printf("国語の平均点は%d\n", total_japanese_score/people_num)
}
