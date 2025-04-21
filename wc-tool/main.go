/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"wc-tool/cmd"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cmd.ReadFromStdin()
	} else {
		cmd.Execute()
	}
}
