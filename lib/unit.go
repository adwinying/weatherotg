package lib

import "fmt"

type TemperatureUnit string

const (
	Celsius    TemperatureUnit = "c"
	Fahrenheit TemperatureUnit = "f"
)

var TemperatureUnitValueMap = map[string]TemperatureUnit{
	"c": Celsius,
	"f": Fahrenheit,
}

func ParseTemperatureUnit(input string) (TemperatureUnit, error) {
	result, ok := TemperatureUnitValueMap[input]
	if !ok {
		return result, fmt.Errorf("Invalid temperature unit: %s", input)
	}

	return result, nil
}
