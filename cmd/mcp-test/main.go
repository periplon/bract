package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/periplon/bract/internal/dsl"
)

func main() {
	var (
		validate = flag.Bool("validate", false, "Validate script syntax without executing")
		format   = flag.Bool("format", false, "Format the script")
		output   = flag.String("o", "", "Output file for formatted script")
		timeout  = flag.Duration("timeout", 0, "Execution timeout (0 for no timeout)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <script.dsl>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "MCP Test DSL Runner - Execute DSL scripts to test MCP servers\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s test.dsl                  # Run a DSL script\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -validate test.dsl        # Validate script syntax\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format test.dsl          # Format script\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format -o out.dsl in.dsl # Format and save to file\n", os.Args[0])
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	scriptFile := flag.Arg(0)

	// Validate only
	if *validate {
		if err := dsl.ValidateFile(scriptFile); err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Script is valid")
		return
	}

	// Format only
	if *format {
		content, err := os.ReadFile(scriptFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
			os.Exit(1)
		}

		formatted, err := dsl.FormatScript(string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format script: %v\n", err)
			os.Exit(1)
		}

		if *output != "" {
			if err := os.WriteFile(*output, []byte(formatted), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Formatted script written to %s\n", *output)
		} else {
			fmt.Print(formatted)
		}
		return
	}

	// Execute script
	interpreter := dsl.NewInterpreter()

	ctx := context.Background()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	if err := interpreter.ExecuteFile(ctx, scriptFile); err != nil {
		fmt.Fprintf(os.Stderr, "Execution failed: %v\n", err)
		os.Exit(1)
	}

	// Clean up
	if client := interpreter.GetClient(); client != nil {
		if err := client.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close client: %v\n", err)
		}
	}
}

func init() {
	// Set up proper cleanup on interrupt
	setupCleanup()
}

func setupCleanup() {
	// Handle cleanup on interrupt/termination
	// This is handled by context cancellation in main
}
