package lib

import (
	"fmt"
)

type DisplayMode string

const (
	Minimal  DisplayMode = "minimal"
	Default  DisplayMode = "default"
	Extended DisplayMode = "extended"
)

var DisplayModeOrder = []DisplayMode{
	Minimal,
	Default,
	Extended,
}

var DisplayModeValueMap = map[string]DisplayMode{
	"minimal":  Minimal,
	"default":  Default,
	"extended": Extended,
}

var DisplayModeStringMap = map[DisplayMode]string{
	Minimal:  "Minimal",
	Default:  "Default",
	Extended: "Extended",
}

func ParseDisplayMode(input string) (DisplayMode, error) {
	result, ok := DisplayModeValueMap[input]
	if !ok {
		return result, fmt.Errorf("Invalid display mode: %s", input)
	}

	return result, nil
}
