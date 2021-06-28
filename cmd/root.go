package cmd

import (
	"fmt"
	"os"

	"github.com/emgage/slug-rename/internal/lib/version"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringP("language", "l", "", "Set substitution language")
}

var rootCmd = &cobra.Command{
	Use:           "slug-rename [files]",
	Short:         "Renames files/directories to a slugified filename",
	SilenceErrors: true,
	Version:       fmt.Sprintf("%s-%s", version.Version, version.Commit),
	Args:          cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// keep dots in filenames
		m_pre := map[string]string{".": "xxxxxdotxxxx"}
		m_post := map[string]string{"xxxxxdotxxxx": "."}

		for _, f := range args {
			t := ""

			if l, err := cmd.Flags().GetString("language"); err == nil && l != "" {
				t = slug.Substitute(slug.MakeLang(slug.Substitute(f, m_pre), l), m_post)
			} else {
				t = slug.Substitute(slug.Make(slug.Substitute(f, m_pre)), m_post)
			}

			if f == t {
				continue
			}

			if err := os.Rename(f, t); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		return nil
	},
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
