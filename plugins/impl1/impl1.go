package main

import (
	"github.com/shashank-priyadarshi/go-plugin/app/common"
)

// go build -buildmode=plugin -o ../app/impl1.so .

var Calculator calculator

type calculator struct {
	common.Calculator
}

func (c *calculator) Add(args ...int) int {
	var result int

	for _, arg := range args {
		result += arg
	}

	return result
}

func (c *calculator) Sub(x, y int) int {
	return x - y
}

func (c *calculator) Mul(args ...int) int {
	var result int

	for _, arg := range args {
		result *= arg
	}

	return result
}
