package main

import (
	"bufio"
	"easylog/internal/converter/cstring"
	"easylog/internal/converter/json"
	"easylog/internal/converter/kvs"
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

		readFile(filename, convert, includes, excludes)
		return nil
	},
}

func init() {
	fileCmd.PersistentFlags().String("filename", "/var/log/containers/nxw-sv__avo.log", "filename path")
	fileCmd.PersistentFlags().String("format", "json", "format line json, string")
	fileCmd.PersistentFlags().StringSlice("excludes", []string{}, "fields to be excluded")
	fileCmd.PersistentFlags().StringSlice("includes", []string{}, "fields to be excluded")

	rootCmd.AddCommand(fileCmd)
}
func readFile(filename string, writer func(map[string]string) string, in_fields, ex_fields []string) (string, error) {
	//Read file
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(writer(kvs.ToKVS(scanner.Text(), in_fields, ex_fields)))
	}

	return "", nil
}
