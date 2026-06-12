package main

import (
	"fmt"
	"log"

	"github.com/elioneto/tuix/ascii"
)

func main() {
	// Available fonts
	fmt.Println("Available fonts:", ascii.AvailableFonts())
	fmt.Println()

	// Generate ASCII art with different fonts
	text := "Hello!"
	examples := []struct{ name, text string }{
		{"graffiti", text},
		{"standard", text},
		{"big", "ASCII"},
		{"block", "BLOCK"},
		{"shadow", "SHADOW"},
	}

	for _, ex := range examples {
		art, err := ascii.Generate(ex.text, ex.name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("=== %s (%q) ===\n", ex.name, ex.text)
		fmt.Println(art)
		fmt.Println()
	}

	// Using the Must helper for one-liners
	art := ascii.Must(ascii.Generate("TUIX", "graffiti"))
	fmt.Println("=== Must helper (graffiti) ===")
	fmt.Println(art)
}
