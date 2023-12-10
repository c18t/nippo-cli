/*
Copyright © 2023 ɯ̹t͡ɕʲi <xc18tx@gmail.com>
This file is part of CLI application nippo-cli.
*/
package main

import "github.com/c18t/nippo-cli/cmd"

// assign -ldflags on build
var version string

func main() {
	cmd.Version = version
	cmd.Execute()
}
