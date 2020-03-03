package app

import (
    "github.com/spf13/cobra"
)

func NewRootCmd(name, version string, fn func(c *cobra.Command, args []string) error) *cobra.Command {
    var rootCmd *cobra.Command

    rootCmd = &cobra.Command{
        Use: name + " <yaml>",
        SilenceUsage: true,
        SilenceErrors: true,
        DisableFlagsInUseLine: true,
        RunE: func(c *cobra.Command, args []string) error {
            if len(args) == 0 {
                c.HelpFunc()(c, args)
                return nil
            }

            return fn(c, args)
        },
        Version: version,
    }

    rootCmd.Flags().BoolP("help", "h", false, "Output usage information")

    rootCmd.SetVersionTemplate("{{printf \"%s\" .Version}}\n")
    rootCmd.Flags().BoolP("version", "V", false, "Display this application version")

    return rootCmd
}
