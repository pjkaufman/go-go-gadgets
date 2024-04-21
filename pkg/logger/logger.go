package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func WriteError(msg string) {
	color.New(color.FgRed).Fprintln(os.Stderr, msg)
	os.Exit(-1)
}

func WriteInfo(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

func WriteWarn(msg string) {
	color.New(color.FgYellow).Fprintln(os.Stdout, msg)
}

func WriteInfoWithColor(msg string, displayColor color.Attribute) {
	color.New(displayColor).Fprintln(os.Stdout, msg)
}

func GetInputString(prompt string) string {
	fmt.Println(prompt)
	// based on https://stackoverflow.com/a/20895629 since for some reason spaces were not read properly by fmt.Scanln
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		WriteError(fmt.Sprintf("failed to read in the string provided: %s", err))
	}

	response = strings.TrimRight(response, "\n")

	return response
}

func GetInputInt(prompt string) int {
	fmt.Println(prompt)
	var response int
	fmt.Scanf("%d", &response)

	return response
}
