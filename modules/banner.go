// Package modules provides functions for displaying banners and messages in the console
package modules

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Console represents a console for displaying formatted text
type Console struct {
	logger *Logger
}

// NewConsole creates a new Console instance
func NewConsole() *Console {
	return &Console{
		logger: GetLogger(),
	}
}

// Print prints a message to the console
func (c *Console) Print(message string, style ...string) {
	// Apply styling if provided
	if len(style) > 0 {
		switch style[0] {
		case "bold red":
			color.New(color.FgRed, color.Bold).Print(message)
		case "bold green":
			color.New(color.FgGreen, color.Bold).Print(message)
		case "bold yellow":
			color.New(color.FgYellow, color.Bold).Print(message)
		case "bold blue":
			color.New(color.FgBlue, color.Bold).Print(message)
		case "bold magenta":
			color.New(color.FgMagenta, color.Bold).Print(message)
		case "bold cyan":
			color.New(color.FgCyan, color.Bold).Print(message)
		case "bold white":
			color.New(color.FgWhite, color.Bold).Print(message)
		case "yellow":
			color.New(color.FgYellow).Print(message)
		case "red":
			color.New(color.FgRed).Print(message)
		case "green":
			color.New(color.FgGreen).Print(message)
		case "blue":
			color.New(color.FgBlue).Print(message)
		case "magenta":
			color.New(color.FgMagenta).Print(message)
		case "cyan":
			color.New(color.FgCyan).Print(message)
		default:
			fmt.Print(message)
		}
	} else {
		fmt.Print(message)
	}

	c.logger.Debug("Console print: %s", message)
}

// Println prints a message to the console with a newline
func (c *Console) Println(message string, style ...string) {
	c.Print(message+"\n", style...)
}

// PrintWelcomeBanner prints the welcome banner
func PrintWelcomeBanner(console *Console, username string) {
	logger := GetLogger()
	logger.Info("Printing welcome banner for user: %s", username)

	// Print the ASCII art banner
	fmt.Println("                      ***")
	fmt.Println("                      ***")
	fmt.Println("                      ***")
	fmt.Println("                      ***")
	fmt.Println("                      ***")
	fmt.Println("                      ***")
	fmt.Println("                      ***")

	// Print the OTTO banner with colors
	color.New(color.FgYellow, color.Bold).Println("    ██████  ████████ ████████  ██████ ")
	color.New(color.FgRed, color.Bold).Println("    ██    ██    ██       ██    ██    ██")
	color.New(color.FgGreen, color.Bold).Println("    ██    ██    ██       ██    ██    ██")
	color.New(color.FgBlue, color.Bold).Println("    ██    ██    ██       ██    ██    ██")
	color.New(color.FgMagenta, color.Bold).Println("    ██    ██    ██       ██    ██    ██")
	color.New(color.FgCyan, color.Bold).Println("     ██████     ██       ██     ██████ ")
	fmt.Println()

	// Print welcome message
	color.New(color.FgWhite, color.Bold).Println("    Your Personal AI :) ")
	fmt.Println()
	color.New(color.FgGreen, color.Bold).Printf("    Welcome, %s!\n\n", username)
	color.New(color.FgRed, color.Bold).Println("    Enter /help to get help. ")

	logger.Debug("Welcome banner printed")
}

// PrintSeparator prints a separator line
func PrintSeparator(console *Console) {
	logger := GetLogger()
	logger.Info("Printing separator")

	// Print a newline first
	fmt.Println()

	// Create and print the separator
	f1 := color.New(color.FgBlue, color.Bold).Sprint("~")
	f2 := color.New(color.FgYellow, color.Bold).Sprint("*")
	separator := strings.Repeat(f1+f2, 76)
	fmt.Println(separator)
	fmt.Println() // End with another newline

	logger.Debug("Separator printed")
}

// PrintCustomBanner prints a custom banner
func PrintCustomBanner(console *Console, text string, style ...string) {
	logger := GetLogger()
	logger.Info("Printing custom banner: '%s'", text)

	// Default style
	bannerStyle := "bold green"
	if len(style) > 0 {
		bannerStyle = style[0]
	}

	// Create a box around the text
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Print the top border
	border := "+" + strings.Repeat("-", maxLen+2) + "+"
	console.Println(border, bannerStyle)

	// Print each line with padding
	for _, line := range lines {
		padding := strings.Repeat(" ", maxLen-len(line))
		console.Println("| "+line+padding+" |", bannerStyle)
	}

	// Print the bottom border
	console.Println(border, bannerStyle)

	logger.Debug("Custom banner printed")
}

// PrintErrorMessage prints an error message
func PrintErrorMessage(console *Console, message string) {
	logger := GetLogger()
	logger.Info("Printing error message: '%s'", message)
	PrintCustomBanner(console, message, "bold red")
	logger.Debug("Error message printed")
}

// PrintSuccessMessage prints a success message
func PrintSuccessMessage(console *Console, message string) {
	logger := GetLogger()
	logger.Info("Printing success message: '%s'", message)
	PrintCustomBanner(console, message, "bold green")
	logger.Debug("Success message printed")
}
