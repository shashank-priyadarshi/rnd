package impl1

import (
	"github.com/shashank-priyadarshi/go-plugin/common"
)

// go build -buildmode=plugin -o ../app/impl1.so .

type calculator struct {
	common.Calculator
}

func NewCalculator() *calculator {
	return &calculator{}
}
