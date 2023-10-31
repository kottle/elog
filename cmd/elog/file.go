package main

import (
	"bufio"
	"easylog/internal/filter"
	"easylog/internal/kvs"
	"easylog/internal/writer"
	"easylog/internal/writer/cstring"
	"easylog/internal/writer/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "",
	Long:  ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("filename")
		check(err)
		format, err := cmd.Flags().GetString("format")
		check(err)
		filterpath, err := cmd.Flags().GetString("filterpath")
		check(err)
		themepath, err := cmd.Flags().GetString("themepath")
		check(err)

		var writer writer.IWriter
		switch format {
		case "json":
			writer = json.New("")
		case "string":
			writer, err = cstring.New(themepath)
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
		readFile(filename, writer, f)
		return nil
	},
}

func init() {
	fileCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	fileCmd.PersistentFlags().String("format", "json", "format line json, string")
	fileCmd.PersistentFlags().String("filterpath", "", "filter file path")
	fileCmd.PersistentFlags().String("themepath", "", "theme file path")

	rootCmd.AddCommand(fileCmd)
}
func readFile(filename string, writer writer.IWriter, filter *filter.Filter) (string, error) {
	//Read file
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		kvs := kvs.ToKVS(scanner.Text(), filter)
		if len(kvs) > 0 {
			fmt.Println(writer.Write(kvs))
		}
	}

	return "", nil
}
