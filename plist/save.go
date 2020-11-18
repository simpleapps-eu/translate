package plist

import (
	"fmt"
	"io"
)

// SaveEntries is a synchronous function that will take a channel with Plist
// entries and stream them to a writer in XML format. The function will return
// when all entries have been written. The goroutine feeding entryChan should
// close the channel once it has finished. The closing of the channel indicates
// to SaveEntries that it can finish too. The function then returns the number
// of entries it has written.
func SaveEntries(entryChan <-chan Entry, tgtFile io.Writer) (n int) {

	fmt.Fprintln(tgtFile, plistPrefix)
	for entry := range entryChan {
		fmt.Fprintf(tgtFile, plistEntry, entry.ID, entry.Str)
		n++
	}
	fmt.Fprintln(tgtFile, plistPostfix)
	return
}

var (
	plistPrefix = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>`

	plistEntry = "\t<key>%s</key>\n\t<string>%s</string>\n"

	plistPostfix = "</dict>\n</plist>"
)
