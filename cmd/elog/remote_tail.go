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
		server, err := cmd.Flags().GetString("server")
		check(err)
		username, err := cmd.Flags().GetString("username")
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
		opts := rtail.Options{
			Address:    server,
			Port:       "22",
			User:       username,
			Filename:   filename,
			Key:        "/Users/eliofrancesconi/.ssh/id_rsa",
			KnownHosts: "/Users/eliofrancesconi/.ssh/known_hosts",
		}

		rTail(opts, convert, includes, excludes)
		return nil
	},
}

func init() {
	remoteTailCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	remoteTailCmd.PersistentFlags().String("format", "json", "format line json, string")
	remoteTailCmd.PersistentFlags().StringSlice("excludes", []string{}, "fields to be excluded")
	remoteTailCmd.PersistentFlags().StringSlice("includes", []string{}, "fields to be excluded")
	remoteTailCmd.PersistentFlags().String("server", "192.168.32.9", "server address")
	remoteTailCmd.PersistentFlags().String("username", "root", "username to connect to server")

	rootCmd.AddCommand(remoteTailCmd)
}

func rTail(opts rtail.Options, writer func(kvs map[string]string) string, in_fields, ex_fields []string) {
	logrus.Infof("tail %s", opts.Filename)
	ctx := context.Background()
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
