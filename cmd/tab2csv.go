// Copyright © 2016 Wei Shen <shenwei356@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/csv"
	"fmt"
	"runtime"

	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// tab2csvCmd represents the tab2csv command
var tab2csvCmd = &cobra.Command{
	Use:   "tab2csv",
	Short: "convert tabular format to CSV",
	Long: `convert tabular format to CSV

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		files := getFileList(args)
		runtime.GOMAXPROCS(config.NumCPUs)

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		writer := csv.NewWriter(outfh)
		writer.Comma = ','

		for _, file := range files {
			csvReader, err := newCSVReaderByConfig(config, file)
			checkError(err)
			csvReader.Reader.Comma = '\t'
			csvReader.Run()

			printMetaLine := true
			for chunk := range csvReader.Ch {
				checkError(chunk.Err)

				if printMetaLine && len(csvReader.MetaLine) > 0 {
					outfh.WriteString(fmt.Sprintf("sep=%s\n", string(writer.Comma)))
					printMetaLine = false
				}

				for _, record := range chunk.Data {
					checkError(writer.Write(record))
				}
			}
		}
		writer.Flush()
		checkError(writer.Error())
	},
}

func init() {
	RootCmd.AddCommand(tab2csvCmd)
}
