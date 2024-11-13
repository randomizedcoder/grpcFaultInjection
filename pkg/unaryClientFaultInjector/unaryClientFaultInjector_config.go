package unaryClientFaultInjector

import (
	"fmt"
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
	switch str {
	case "Modulus":
		mode = Modulus
	case "Percent":
		mode = Percent
		//default:
	}
	return mode
}
