package app

import (
	"fmt"
	"plugin"
)

const (
	PLUGIN_PATH     = "impl1.so"
	PLUGIN_FUNCTION = "NewCalculator"
)

type Calculator interface {
	Add(...int) int
	Sub(x, y int) int
	Mul(...int) int
}

func main() {
	add := []int{1, 2, 3, 4, 5}
	sub := []int{3, -11}
	mul := []int{1, 2, 3, 4, 5}

	plugin, err := plugin.Open(PLUGIN_PATH)
	if err != nil {
		fmt.Printf("error opening plugin from path %s: %v\n", PLUGIN_PATH, err)
		return
	}

	sym, err := plugin.Lookup(PLUGIN_FUNCTION)
	if err != nil {
		fmt.Printf("error looking for %s in loaded plugin: %v\n", PLUGIN_FUNCTION, err)
		return
	}

	calc, ok := sym.(Calculator)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		return
	}

	result := calc.Add(add...)
	fmt.Print("addition result: ", result)
	result = calc.Sub(sub[0], sub[1])
	fmt.Print("subtraction result: ", result)
	result = calc.Mul(mul...)
	fmt.Print("multiplication result: ", result)
}
