package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func check(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "helper",
	Long:  `helper details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

func main1() int {
	return Execute()
}

func main() {
	logrus.SetLevel(logrus.ErrorLevel)
	os.Exit(main1())
}
