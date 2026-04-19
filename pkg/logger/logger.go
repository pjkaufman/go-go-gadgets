package logger

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
)

func WriteError(msg string) {
	fmt.Fprintln(os.Stderr, errorStyle.Render(msg))
	os.Exit(-1)
}

func WriteErrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, errorStyle.Render(format), a...)
	os.Exit(-1)
}

func WriteInfo(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

func WriteInfof(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func WriteWarn(msg string) {
	fmt.Fprintln(os.Stdout, warningStyle.Render(msg))
}

func WriteWarnf(format string, a ...any) {
	fmt.Fprintf(os.Stdout, warningStyle.Render(format), a...)
}

func WriteInfoWithColor(msg string, displayColor color.Color) {
	fmt.Fprintln(os.Stdout, lipgloss.NewStyle().Foreground(displayColor).Render(msg))
}

func GetInputString(prompt string) string {
	fmt.Println(prompt)
	// based on https://stackoverflow.com/a/20895629 since for some reason spaces were not read properly by fmt.Scanln
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		WriteErrorf("failed to read in the string provided: %s\n", err)
	}

	response = strings.TrimRight(response, "\n")

	return response
}

func GetInputInt(prompt string) int {
	fmt.Println(prompt)
	var response int
	_, err := fmt.Scanf("%d", &response)
	if err != nil {
		WriteErrorf("failed to read in the integer from the user: %s\n", err)
	}

	return response
}
