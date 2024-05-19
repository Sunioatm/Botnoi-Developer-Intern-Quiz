package main

import (
	"fmt"
	"strings"
)

func main() {
	var input_num int
	fmt.Print("Enter a number: ")
	fmt.Scan(&input_num)
	generatePyramid(input_num)

}

func generatePyramid(input_num int) {
	var i int
	for i = 0; i < input_num; i++ {
		fmt.Println(strings.Repeat("*", i+1))
	}
	for i = i - 1; i >= 1; i-- {
		fmt.Println(strings.Repeat("*", i))
	}
}
