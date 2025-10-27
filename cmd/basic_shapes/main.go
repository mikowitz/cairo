package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mikowitz/cairo/examples"
)

func main() {
	// Default output path
	outputPath := "output.png"

	// Allow specifying custom output path as first argument
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	// Generate the image
	if err := examples.GenerateBasicShapes(outputPath); err != nil {
		log.Fatalf("Failed to generate basic shapes: %v", err)
	}

	fmt.Printf("Successfully generated image at %s\n", outputPath)
}
