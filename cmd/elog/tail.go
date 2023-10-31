package main

import (
	"easylog/internal/kvs"
	"easylog/internal/writer"
	"easylog/internal/writer/cstring"
	"easylog/internal/writer/json"
	"fmt"
	"io"

	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "",
	Long:  ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("filename")
		check(err)
		format, err := cmd.Flags().GetString("format")
		check(err)
		includes, err := cmd.Flags().GetStringSlice("includes")
		check(err)
		excludes, err := cmd.Flags().GetStringSlice("excludes")
		check(err)

		var writer writer.IWriter
		switch format {
		case "json":
			writer = json.New("")
		case "string":
			writer, err = cstring.New("")
			check(err)
		default:
			return fmt.Errorf("format %s not supported", format)
		}
		tail(filename, writer, includes, excludes)
		return nil
	},
}

func init() {
	tailCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	tailCmd.PersistentFlags().String("format", "json", "format line json, string")
	tailCmd.PersistentFlags().StringSlice("excludes", []string{}, "fields to be excluded")
	tailCmd.PersistentFlags().StringSlice("includes", []string{}, "fields to be excluded")

	rootCmd.AddCommand(tailCmd)
}

func tail(filename string, writer writer.IWriter, in_fields, ex_fields []string) {
	logrus.Infof("tail %s", filename)
	t, err := follower.New(filename, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	check(err)

	for line := range t.Lines() {
		fmt.Println(writer.Write(kvs.ToKVS(line.String(), nil)))
	}
	logrus.Infof("tail %s done", filename)
}
