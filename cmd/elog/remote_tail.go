package main

import (
	"context"
	"easylog/internal/converter/cstring"
	"easylog/internal/converter/json"
	"easylog/internal/converter/kvs"
	"easylog/internal/rtail"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var remoteTailCmd = &cobra.Command{
	Use:   "remote-tail",
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

		var convert func(map[string]string) string
		switch format {
		case "json":
			convert = json.Convert
		case "string":
			convert = cstring.Convert
		default:
			return fmt.Errorf("format %s not supported", format)
		}
		rTail(filename, convert, includes, excludes)
		return nil
	},
}

func init() {
	remoteTailCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	remoteTailCmd.PersistentFlags().String("format", "json", "format line json, string")
	remoteTailCmd.PersistentFlags().StringSlice("excludes", []string{}, "fields to be excluded")
	remoteTailCmd.PersistentFlags().StringSlice("includes", []string{}, "fields to be excluded")

	rootCmd.AddCommand(remoteTailCmd)
}

func rTail(filename string, writer func(kvs map[string]string) string, in_fields, ex_fields []string) {
	logrus.Infof("tail %s", filename)
	ctx := context.Background()
	opts := rtail.Options{
		Address:    "10.0.4.38",
		Port:       "22",
		User:       "root",
		Filename:   filename,
		Key:        "/Users/eliofrancesconi/.ssh/id_rsa",
		KnownHosts: "/Users/eliofrancesconi/.ssh/known_hosts",
	}
	lines := make(chan string)
	go func() {
		err := rtail.Tail(ctx, opts, lines)
		check(err)
	}()

	for {
		select {
		case line := <-lines:
			fmt.Println(writer(kvs.ToKVS(line, in_fields, ex_fields)))
		}
	}
}
