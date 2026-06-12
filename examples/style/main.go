package main

import (
	"fmt"

	"github.com/elioneto/tuix"
	"github.com/elioneto/tuix/color"
)

func main() {
	// Demonstration of the Style API (non-interactive — outputs styled text)

	heading := tuix.NewStyle().
		Foreground(color.Hex("#00d4aa")).
		Bold(true).
		Render("Welcome to tuix Style API")

	subtitle := tuix.NewStyle().
		Foreground(color.Gray).
		Italic(true).
		Render("Port of Lip Gloss Style - zero deps")

	fmt.Println(heading)
	fmt.Println(subtitle)
	fmt.Println()

	// Bordered box
	box := tuix.NewStyle().
		Border(tuix.BorderRounded, color.Hex("#00d4aa")).
		Padding(1, 2).
		Render("Hello from a rounded border box!")

	fmt.Println(box)
	fmt.Println()

	// Different border styles
	fmt.Println(tuix.NewStyle().Border(tuix.BorderNormal, color.Hex("#e94560")).Padding(0, 1).Render("Normal"))
	fmt.Println()
	fmt.Println(tuix.NewStyle().Border(tuix.BorderRounded, color.Hex("#2ecc71")).Padding(0, 1).Render("Rounded"))
	fmt.Println()
	fmt.Println(tuix.NewStyle().Border(tuix.BorderThick, color.Hex("#3498db")).Padding(0, 1).Render("Thick"))
	fmt.Println()
	fmt.Println(tuix.NewStyle().Border(tuix.BorderDouble, color.Hex("#f39c12")).Padding(0, 1).Render("Double"))
	fmt.Println()

	// Styled text with background
	notice := tuix.NewStyle().
		Foreground(color.White).
		Background(color.Hex("#e94560")).
		Bold(true).
		Padding(0, 1).
		Render(" Important notice with background ")

	fmt.Println(notice)
	fmt.Println()

	// JoinHorizontal and JoinVertical
	left := tuix.NewStyle().
		Foreground(color.Hex("#00d4aa")).
		Border(tuix.BorderNormal, color.Hex("#00d4aa")).
		Padding(0, 1).
		Render("Left")

	right := tuix.NewStyle().
		Foreground(color.Hex("#e94560")).
		Border(tuix.BorderNormal, color.Hex("#e94560")).
		Padding(0, 1).
		Render("Right")

	fmt.Println("JoinHorizontal:")
	fmt.Println(tuix.JoinHorizontal(left, right))
	fmt.Println()

	fmt.Println("JoinVertical:")
	fmt.Println(tuix.JoinVertical(left, right))

	// Style with margin
	margined := tuix.NewStyle().
		Border(tuix.BorderRounded, color.Hex("#3498db")).
		Padding(1, 2).
		Margin(1, 0).
		Render("Margined box")

	fmt.Println()
	fmt.Println(margined)
	fmt.Println()

	fmt.Println("All outputs use ANSI escape codes - no external libraries!")
}
