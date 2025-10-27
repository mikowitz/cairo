package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mikowitz/cairo/examples"
)

func main() {
	// Default output path
	outputPath := "complex_shapes.png"

	// Allow specifying custom output path as first argument
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	// Generate the complex shapes image
	fmt.Println("Generating complex shapes image demonstrating all Context functionality...")
	if err := examples.GenerateComplexShapes(outputPath); err != nil {
		log.Fatalf("Failed to generate complex shapes: %v", err)
	}

	fmt.Printf("✓ Successfully generated complex shapes at %s\n", outputPath)
	fmt.Println("\nThis image demonstrates all 20 implemented Context methods:")
	fmt.Println("  • Paint, Save, Restore, Status")
	fmt.Println("  • SetSourceRGB, SetSourceRGBA")
	fmt.Println("  • MoveTo, LineTo, Rectangle, ClosePath, NewPath, NewSubPath")
	fmt.Println("  • GetCurrentPoint, HasCurrentPoint")
	fmt.Println("  • Fill, FillPreserve, Stroke, StrokePreserve")
	fmt.Println("  • SetLineWidth, GetLineWidth")
}
