package view

import (
	"fmt"
)

type viewModel struct {
	Input  chan<- interface{}
	Output interface{}
}

func message(output interface{}) bool {
	ret := output != nil
	if ret {
		fmt.Printf("%v", output)
	}
	return ret
}

func either2(input interface{}, err error) interface{} {
	if err != nil {
		return err
	}
	return input
}
