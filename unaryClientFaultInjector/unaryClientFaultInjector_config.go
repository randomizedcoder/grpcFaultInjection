package unaryClientFaultInjector

import (
	"fmt"
	"strings"
)

type Mode int32

const (
	Modulus Mode = iota
	Percent Mode = 1
)

type ModeValue struct {
	Mode  Mode
	Value int
}

type UnaryClientInterceptorConfig struct {
	Client ModeValue
	Server ModeValue
	Codes  string
}

func (m Mode) toString() {
	switch m {
	case Modulus:
		fmt.Println("Modulus")
	case Percent:
		fmt.Println("Percent")
	default:
		fmt.Println("Invalid Mode")
	}
}

func StringToMode(str string) (mode Mode) {
	switch strings.ToLower(str) {
	case "m":
		mode = Modulus
	case "mod":
		mode = Modulus
	case "modulus":
		mode = Modulus
	case "p":
		mode = Percent
	case "per":
		mode = Percent
	case "percent":
		mode = Percent
		//default:
	}
	return mode
}
