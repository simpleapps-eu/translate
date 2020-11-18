package translate

import (
	"fmt"
	"io"
)

// SaveLines will take a channel with text lines and stream them to writer in
// unix text format (i.e. adding a linefeed character after every line).
// The function will return when all lines have been written.
// The goroutine feeding lineChan should close the channel once it has finished.
// The closing of the channel indicates to SaveLines that it can finish
// too. The function then returns the number of lines it has written.
// Note that no BOM is being written to the tgtFile to indicate the encoding of
// the text stream (which is most likely utf-8).
func SaveLines(lineChan <-chan string, tgtFile io.Writer) (n int) {
	for line := range lineChan {
		if n > 0 {
			fmt.Fprintln(tgtFile)
		}
		fmt.Fprint(tgtFile, line)
		n++
	}
	// Unless the whole text was just a single line, always end with a newline.
	if n > 1 {
		fmt.Fprintln(tgtFile)
	}
	return
}
