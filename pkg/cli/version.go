package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/live-by-unix/anyaccountlogin/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information for AnyAccountLogin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		full, _ := cmd.Flags().GetBool("full")
		json, _ := cmd.Flags().GetBool("json")

		if json {
			if full {
				fmt.Printf(`{"version":"%s","gitCommit":"%s","buildDate":"%s"}`+"\n",
					version.Version, version.GitCommit, version.BuildDate)
			} else {
				fmt.Printf(`{"version":"%s"}`+"\n", version.Version)
			}
		} else {
			if full {
				fmt.Println(version.GetFullVersion())
			} else {
				fmt.Println(version.GetVersion())
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("full", "f", false, "Print full version information including git commit and build date")
	versionCmd.Flags().Bool("json", false, "Output version information in JSON format")
}
