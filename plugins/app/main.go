package main

import (
	"fmt"

	"github.com/shashank-priyadarshi/go-plugin/app/common"
)

func main() {
	add := []int{1, 2, 3, 4, 5}
	sub := []int{3, -11}
	mul := []int{1, 2, 3, 4, 5}

	calc, err := common.Common()
	if err != nil {
		fmt.Printf("error opening plugin: %v\n", err)
		return
	}

	result := calc.Add(add...)
	fmt.Println("addition result: ", result)
	result = calc.Sub(sub[0], sub[1])
	fmt.Println("subtraction result: ", result)
	result = calc.Mul(mul...)
	fmt.Println("multiplication result: ", result)
}
