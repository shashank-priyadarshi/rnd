package main

import (
	"github.com/shashank-priyadarshi/go-plugin/app/common"
)

// go build -buildmode=plugin -o ../app/impl1.so .

var Calculator calculator

type calculator struct {
	common.Calculator
}

func NewCalculator() *calculator {
	return &calculator{}
}
