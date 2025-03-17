package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rivet",
		Short: "Rivet - a fast C/C++ build system",
		Long:  `Rivet is a fast C/C++ build system with cross-compilation and workspace support.`,
	}

	/* TODO: add command impls */

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
