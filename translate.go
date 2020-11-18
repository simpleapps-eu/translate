package translate

import (
	"io"

	"github.com/simpleapps-eu/translate/dotstrings"
	"github.com/simpleapps-eu/translate/plist"
)

func TranslatePlistFile(srcFile io.Reader, translations map[string]dotstrings.Message, tgtFile io.Writer) (n int, err error) {

	// Start loading entries asynchronously
	entryChan, errChan := plist.LoadEntries(srcFile)

	// Start translation entries asynchronously
	entryChan = TranslatePlistEntries(entryChan, translations)

	// Save the translated entries synchronously
	n = plist.SaveEntries(entryChan, tgtFile)

	err, _ = <-errChan
	return
}

func TranslatePlistEntries(entryChan <-chan plist.Entry, translations map[string]dotstrings.Message) <-chan plist.Entry {
	dstChan := make(chan plist.Entry, 3)

	translator := func(srcChan <-chan plist.Entry, dstChan chan<- plist.Entry, translations map[string]dotstrings.Message) {
		defer close(dstChan)
		for entry := range srcChan {
			tran, present := translations[entry.ID]
			if present {
				// Don't output translations to empty string
				if len(tran.Str) > 0 {
					dstChan <- plist.Entry{ID: entry.ID, Str: tran.Str}
				}
				continue
			}
			dstChan <- entry
		}
	}

	go translator(entryChan, dstChan, translations)
	return dstChan
}

func TranslateIDsFile(srcFile io.Reader, translations map[string]dotstrings.Message, translationsFallback map[string]dotstrings.Message, tgtFile io.Writer) (n int, err error) {

	// Load IDs to be translated concurrently
	lineChan, errChan1 := LoadLines(srcFile)

	// Start translation concurrently
	lineChan = TranslateIDs(lineChan, translations, translationsFallback)

	// Start unescaping the translated lines concurrently
	lineChan, errChan2 := unescapeLines(lineChan)

	// Synchronously save the lines to the target file
	n = SaveLines(lineChan, tgtFile)

	err, _ = <-errChan1
	if err != nil {
		return
	}

	err, _ = <-errChan2
	return
}

func unescapeLines(srcChan <-chan string) (<-chan string, <-chan error) {
	dstChan := make(chan string, 3)
	errChan := make(chan error, 1)

	unescaper := func(srcChan <-chan string, dstChan chan<- string, errChan chan<- error) {
		defer close(dstChan)
		defer close(errChan)

		for str := range srcChan {
			if len(str) == 0 {
				dstChan <- ""
			} else {
				s, err := dotstrings.StringsUnescape(str)
				if err != nil {
					s = str
					errChan <- err
					return
				}
				dstChan <- s
			}
		}
	}

	go unescaper(srcChan, dstChan, errChan)
	return dstChan, errChan
}

// TranslateIDs will translate a channel of strings where the strings are
// treated as IDs in a translation map. The translationsFallback map normally
// contains the original source language and is used to fill in the gaps
// where translations haven't been provided yet using strings from the source
// language. The idea being that it is better to show a string instead of an id.
func TranslateIDs(srcChan <-chan string, translations map[string]dotstrings.Message, translationsFallback map[string]dotstrings.Message) <-chan string {
	dstChan := make(chan string, 3)

	translator := func(srcChan <-chan string, dstChan chan<- string, translations map[string]dotstrings.Message, translationsFallback map[string]dotstrings.Message) {
		defer close(dstChan)
		for id := range srcChan {
			tran, present := translations[id]
			if present {
				// Don't output translations to empty string
				if len(tran.Str) > 0 {
					dstChan <- tran.Str
				}
				continue
			}
			tran, present = translationsFallback[id]
			if present {
				// Don't output translations to empty string
				if len(tran.Str) > 0 {
					dstChan <- tran.Str
				}
				continue
			}
			dstChan <- id
		}
	}

	go translator(srcChan, dstChan, translations, translationsFallback)
	return dstChan
}

func TranslateTextFile(srcFile io.Reader, translations map[string]dotstrings.Message, tgtFile io.Writer) (n int, err error) {

	// Load text lines to be translated concurrently
	lineChan, errChan := LoadLines(srcFile)

	// Start translation concurrently
	lineChan = TranslateText(lineChan, translations)

	// Synchronously save the lines to the target file
	n = SaveLines(lineChan, tgtFile)

	err, _ = <-errChan

	return
}

// TranslateText will translate a channel of strings where the strings are
// treated as context values in the translations map. Note that translating this
// way does not handle situations where the Context of different translation
// units is the same. The last translation in the translation map will win out
// and shadow the other entries.
// FIXME: TranslateText does not handle the difference in escaping between lines
//  of text read from the srcChan and the translations map.
func TranslateText(srcChan <-chan string, translations map[string]dotstrings.Message) <-chan string {
	dstChan := make(chan string, 3)

	translator := func(srcChan <-chan string, dstChan chan<- string, translations map[string]dotstrings.Message) {
		defer close(dstChan)

		tm := make(map[string]string)
		for _, tran := range translations {
			tm[tran.Ctx] = tran.Str
		}

		for ctx := range srcChan {
			if len(ctx) == 0 {
				// empty line
				dstChan <- ""
			} else {
				str, present := tm[ctx]
				if !present {
					dstChan <- ctx
				} else {
					// Don't output translations to empty string
					if len(str) > 0 {
						dstChan <- str
					}
				}
			}
		}
	}

	go translator(srcChan, dstChan, translations)
	return dstChan
}

func TranslateMessagesFile(srcFile io.Reader, translations map[string]dotstrings.Message, tgtFile io.Writer) (n int, err error) {

	// Load the messages to be translated asynchronously.
	msgChan, errChan := dotstrings.LoadMessages(dotstrings.NewReaderUTF16(srcFile))

	// Start translation asynchronously
	msgChan = TranslateMessages(msgChan, translations)

	// Save the translated messages synchronously
	n = dotstrings.SaveMessages(msgChan, dotstrings.NewWriterUTF16(tgtFile))

	err, _ = <-errChan
	return
}

// TranslateMessages will translate the strings file entries it takes from srcReader and then using a
// translation it finds in the translations map assemble a translation and write it out as a
// strings file entry to dstWriter.
// In case there is no translation available for a source entry, it
// will write the source entry out to the destination and mark it as Fuzzy.
// While translating TranslateString will detect whether the original text to translate (Context)
// has changed between the entry found in the source and the entry on which the translation was based.
// In this case the translation is still written but marked as fuzzy and the context is also changed to
// the new contetx from the source entry.
func TranslateMessages(srcChan <-chan dotstrings.Message, translations map[string]dotstrings.Message) <-chan dotstrings.Message {
	dstChan := make(chan dotstrings.Message, 3)

	translator := func(srcChan <-chan dotstrings.Message, dstChan chan<- dotstrings.Message, translations map[string]dotstrings.Message) {
		defer close(dstChan)
		for src := range srcChan {
			// type Message struc
			// /* Ctx */
			// "ID" = "Str"
			//
			// src
			// /* Show Help */
			// "help_ad_dialog_help_button" = "Show Help";
			//
			// tm
			// /* Show Help */
			// "help_ad_dialog_help_button" = "Hilfe zeigen";
			//
			if tm, ok := translations[src.ID]; ok {
				// There is a translation available for src.ID
				// Are we still talking about the same translation?
				if tm.Ctx == src.Str {
					// We compare the localized Context (tm.Ctx) to the source String (src.Str)
					// The source Context might actually contain a comment on the actual
					// meaning of the source String. For localized .strings files we copy
					// the source String into the localized Context so translators
					// can always have the String available that they need to translate.
					// The actual localized String value always has the last translation.
					// So when you have a localized entry marked as fuzzy the Context
					// of that entry provides the latest source String to translate and
					// the String of that entry provides the previous translation.
					dstChan <- tm
				} else {
					// No, different, so translation is Fuzzy. But do generate entry
					// with previous translation as basis. We put the src.Str (string to
					// be translated) into Ctx and we put tm.Str (previous translation)
					// into Str.
					dstChan <- dotstrings.Message{Fuzzy: true, ID: src.ID, Ctx: src.Str, Str: tm.Str}
				}
			} else {
				// There is no translation for src.ID so use src as basis but mark it as Missing.
				dstChan <- dotstrings.Message{Fuzzy: true, Missing: true, ID: src.ID, Ctx: src.Str, Str: src.Str}
			}
		}
	}

	go translator(srcChan, dstChan, translations)
	return dstChan
}

// TranslateMessagesXLIFF will asynchronously take a channel of dotstrings messages
// and then using the translations table (loaded from an XLIFF file) translate
// the messages and write the translated messages to another channel.
// This function returns 2 channels, a channel that gets the translated messages and a
// channel of error values that is used by this function to push an error onto before terminating.
// The error channel is one way of delivering errors from an asynchronously called function.
func TranslateMessagesXLIFF(srcChan <-chan dotstrings.Message, translations map[string]string) (<-chan dotstrings.Message, <-chan error) {
	msgChan := make(chan dotstrings.Message, 3)
	errChan := make(chan error, 1)

	translator := func(translations map[string]string, srcChan <-chan dotstrings.Message, dstChan chan<- dotstrings.Message, errChan chan<- error) {
		defer close(dstChan)
		defer close(errChan)

		for m := range srcChan {
			m.Str = dotstrings.StringsEscape(translations[m.ID])
			if m.Str != "" {
				dstChan <- m
			} else {
				// There is no translation for m.ID so use m itself as basis but mark it as Missing.
				dstChan <- dotstrings.Message{Fuzzy: true, Missing: true, ID: m.ID, Ctx: m.Str, Str: m.Str}
			}
		}
	}

	go translator(translations, srcChan, msgChan, errChan)
	return msgChan, errChan
}
