package impl1

import (
	"github.com/shashank-priyadarshi/go-plugin/app"
)

// go build -buildmode=plugin -o ../app/impl1.so .

type calculator struct {
	app.Calculator
}

func NewCalculator() *calculator {
	return &calculator{}
}
