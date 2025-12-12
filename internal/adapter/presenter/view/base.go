package view

import (
	"github.com/c18t/nippo-cli/internal/adapter/presenter/view/tui"
)

type viewModel struct {
	Input  chan<- interface{}
	Output interface{}
}

func message(output interface{}) bool {
	ret := output != nil
	if ret {
		tui.Print(output)
	}
	return ret
}

func either2(input interface{}, err error) interface{} {
	if err != nil {
		return err
	}
	return input
}
