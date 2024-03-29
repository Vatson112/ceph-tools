/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ceph-tools",
	Short: "A collection of ceph tools.",
	Long: `This is collection of ceph tools that I use.
See subcommand help for details.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ceph-scripts.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().String("ceph-config", "/etc/ceph/ceph.conf", "Path to ceph config")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logs")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return setupLogs(os.Stdout)
	}
}

func setupLogs(out io.Writer) error {
	if a, _ := rootCmd.Flags().GetBool("debug"); a {

		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	return nil
}
