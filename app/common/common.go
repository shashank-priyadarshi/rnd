package common

import (
	"fmt"
	"plugin"
)

const (
	PLUGIN_PATH     = "impl1.so"
	PLUGIN_FUNCTION = "Calculator"
)

type Calculator interface {
	Add(...int) int
	Sub(x, y int) int
	Mul(...int) int
}

func Common() (calculator Calculator, err error) {
	plugin, err := plugin.Open(PLUGIN_PATH)
	if err != nil {
		return nil, fmt.Errorf("error opening plugin from path %s: %v", PLUGIN_PATH, err)
	}

	sym, err := plugin.Lookup(PLUGIN_FUNCTION)
	if err != nil {
		return nil, fmt.Errorf("error looking for %s in loaded plugin: %v", PLUGIN_FUNCTION, err)
	}

	calc, ok := sym.(Calculator)
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}

	return calc, nil
}
