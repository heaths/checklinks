// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/heaths/checklinks/internal/find"
	"github.com/heaths/checklinks/internal/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "checklinks",
		Long: "Check links in files with optional URL replacements.",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args)
		},
	}

	replacements map[string]string
	verbose      int
)

func init() {
	rootCmd.Flags().StringToStringVarP(&replacements, "replace", "r", map[string]string{}, "optional `regex=replacement` replacements in URLs")
	rootCmd.Flags().CountVarP(&verbose, "verbose", "v", "log verbose output; repeat to increase verbosity``") // empty backticks to remove help argument
}

func Execute() error {
	return rootCmd.Execute()
}

func run(args []string) error {
	log.SetLevel(log.Verbosity(verbose))

	// Check that all patterns are valid.
	for _, pattern := range args {
		if !doublestar.ValidatePathPattern(pattern) {
			return fmt.Errorf("invalid pattern: %s", pattern)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	ctx := context.Background()
	rootFS := os.DirFS(cwd)
	matches := find.Find(ctx, rootFS, args)
	if err != nil {
		log.Fatal("find: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case match, ok := <-matches:
			if !ok {
				log.Debug("done scanning")
				return nil
			}

			fmt.Printf("%s: %s\n", match.Path, match.URL)
		}
	}
}
