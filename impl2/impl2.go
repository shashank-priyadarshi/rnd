package main

import (
	"github.com/shashank-priyadarshi/go-plugin/app/common"
)

// go build -buildmode=plugin -o ./impl2.so .

type calculator struct {
	common.Calculator
}

func NewCalculator() *calculator {
	return &calculator{}
}
