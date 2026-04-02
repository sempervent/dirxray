package cli

import (
	"fmt"
	"os"

	"dirxray/internal/analyze"
	"dirxray/internal/scan"
	"dirxray/internal/tui"
	"github.com/spf13/cobra"
)

// Execute runs the dirxray CLI.
func Execute() error {
	var (
		path           string
		hidden         bool
		maxDepth       int
		followSymlinks bool
		noGitignore    bool
		debug          bool
	)

	root := &cobra.Command{
		Use:   "dirxray [path]",
		Short: "Terminal-native forensic explainer for directories",
		Long:  `dirxray scans a path and presents evidence-backed heuristics in a TUI (Bubble Tea).`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := path
			if target == "" && len(args) > 0 {
				target = args[0]
			}
			if target == "" {
				target = "."
			}
			if debug {
				fmt.Fprintf(os.Stderr, "dirxray: target=%q\n", target)
			}
			abs, err := scan.ValidatePath(target)
			if err != nil {
				return err
			}
			sc, err := scan.Scan(scan.Options{
				Root:            abs,
				Hidden:          hidden,
				MaxDepth:        maxDepth,
				FollowSymlinks:  followSymlinks,
				NoGitignore:     noGitignore,
			})
			if err != nil {
				return err
			}
			ar, err := analyze.Run(abs, sc)
			if err != nil {
				return err
			}
			return tui.Run(abs, sc, ar)
		},
	}

	root.Flags().StringVar(&path, "path", "", "directory to scan (default: . or first arg)")
	root.Flags().BoolVar(&hidden, "hidden", false, "include hidden files and directories")
	root.Flags().IntVar(&maxDepth, "max-depth", 0, "maximum recursion depth (0 = unlimited)")
	root.Flags().BoolVar(&followSymlinks, "follow-symlinks", false, "follow symbolic links when scanning")
	root.Flags().BoolVar(&noGitignore, "no-gitignore", false, "do not apply .gitignore rules")
	root.Flags().BoolVar(&debug, "debug", false, "print debug messages to stderr")

	return root.Execute()
}
