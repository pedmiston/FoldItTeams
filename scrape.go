package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// Scrape reads solution files and writes to JSON.
func Scrape(in io.Reader, dst io.Writer, errDst io.Writer) {
	out, n := ReadAll(in)

	// Pull data from the channel and encode it to JSON.
	encoder := json.NewEncoder(dst)
	for i := 0; i < n; i++ {
		result := <-out
		if result.Err != nil {
			fmt.Fprintf(errDst, "%s,%v\n", result.Data.Filepath(), result.Err)
			continue
		}
		if err := encoder.Encode(result.Data); err != nil {
			fmt.Fprintf(errDst, "%s,%v\n", result.Data.Filepath(), result.Err)
		}
	}
}
