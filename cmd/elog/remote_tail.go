package main

import (
	"context"
	"easylog/internal/filter"
	"easylog/internal/kvs"
	"easylog/internal/rtail"
	"easylog/internal/writer"
	"easylog/internal/writer/cstring"
	"easylog/internal/writer/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var remoteTailCmd = &cobra.Command{
	Use:   "remote-tail",
	Short: "",
	Long:  ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		filenames, err := cmd.Flags().GetStringArray("filenames")
		check(err)
		format, err := cmd.Flags().GetString("format")
		check(err)
		server, err := cmd.Flags().GetString("server")
		check(err)
		username, err := cmd.Flags().GetString("username")
		check(err)
		filterpath, err := cmd.Flags().GetString("filterpath")
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
		var f *filter.Filter
		if filterpath != "" {
			f, err = filter.New(filterpath)
			check(err)
		}
		defer func() {
			if f != nil {
				err := f.Close()
				check(err)
			}
		}()

		opts := rtail.Options{
			Address:    server,
			Port:       "22",
			User:       username,
			Key:        "/Users/eliofrancesconi/.ssh/id_rsa",
			KnownHosts: "/Users/eliofrancesconi/.ssh/known_hosts",
		}

		rTail(filenames, opts, writer, f)
		return nil
	},
}

func init() {
	remoteTailCmd.PersistentFlags().StringArray("filenames", []string{"AVO:/var/log/containers/nxw-sv__avo.log", "AVCFG:/var/log/containers/nxw-sv__av-cfg-fe.log"}, "filename path")
	remoteTailCmd.PersistentFlags().String("format", "json", "format line json, string")
	remoteTailCmd.PersistentFlags().String("server", "192.168.32.9", "server address")
	remoteTailCmd.PersistentFlags().String("username", "root", "username to connect to server")
	remoteTailCmd.PersistentFlags().String("filterpath", "", "filter file path")
	rootCmd.AddCommand(remoteTailCmd)
}

func rTail(filenames []string, opts rtail.Options, writer writer.IWriter, filter *filter.Filter) {
	logrus.Infof("tail %s", filenames)
	ctx := context.Background()
	lines := make(chan string, 4)
	for _, filename := range filenames {
		go func(f string) {
			err := rtail.Tail(ctx, f, opts, lines)
			check(err)
		}(filename)
	}
	for {
		select {
		case line := <-lines:
			kvs := kvs.ToKVS(line, filter)
			if len(kvs) == 0 {
				continue
			}
			fmt.Println(writer.Write(kvs))
		}
	}
}
