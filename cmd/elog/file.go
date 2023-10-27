package main

import (
	"bufio"
	"easylog/internal/common"
	"easylog/internal/converter/cstring"
	"easylog/internal/converter/json"
	"easylog/internal/converter/kvs"
	"easylog/internal/filter"
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
		var convert func(common.KVS) string
		switch format {
		case "json":
			convert = json.Convert
		case "string":
			convert = cstring.Convert
		default:
			return fmt.Errorf("format %s not supported", format)
		}

		readFile(filename, convert, nil)
		return nil
	},
}

func init() {
	fileCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	fileCmd.PersistentFlags().String("format", "json", "format line json, string")

	rootCmd.AddCommand(fileCmd)
}
func readFile(filename string, writer func(common.KVS) string, filter *filter.Filter) (string, error) {
	//Read file
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(writer(kvs.ToKVS(scanner.Text(), filter)))
	}

	return "", nil
}
