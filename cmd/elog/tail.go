package main

import (
	"easylog/internal/converter/cstring"
	"easylog/internal/converter/json"
	"easylog/internal/converter/kvs"
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
		var convert func(map[string]string) string
		switch format {
		case "json":
			convert = json.Convert
		case "string":
			convert = cstring.Convert
		default:
			return fmt.Errorf("format %s not supported", format)
		}
		tail(filename, convert)
		return nil
	},
}

func init() {
	tailCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	tailCmd.PersistentFlags().String("format", "json", "format line json, string")

	rootCmd.AddCommand(tailCmd)
}

func tail(filename string, writer func(kvs map[string]string) string) {
	logrus.Infof("tail %s", filename)
	t, err := follower.New(filename, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	check(err)

	for line := range t.Lines() {
		fmt.Println(writer(kvs.ToKVS(line.String())))
	}
	logrus.Infof("tail %s done", filename)
}
