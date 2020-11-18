package translate

import (
	"fmt"

	"github.com/simpleapps-eu/translate/dotstrings"
	"github.com/simpleapps-eu/translate/xliff"
)

type fromType int

const (
	fromSource fromType = iota
	fromTarget
)

func convertMessagesToTranslationUnits(from fromType, msgChan <-chan dotstrings.Message, tf *xliff.TranslationFile) (<-chan xliff.TranslationUnit, <-chan error) {

	dstChan := make(chan xliff.TranslationUnit, 3)
	errChan := make(chan error, 1)

	n := 0
	translator := func(from fromType, msgChan <-chan dotstrings.Message, tf *xliff.TranslationFile, dstChan chan<- xliff.TranslationUnit, errChan chan<- error) {
		defer close(dstChan)
		defer close(errChan)

		hasTargetLanguage := len(tf.TargetLanguage) > 0

		for m := range msgChan {
			id, e := xliff.XMLEscapeStrict(m.ID)
			if e != nil {
				errChan <- fmt.Errorf("Failed to xml escape Message.Id for string %d (%v)", n+1, e)
				return
			}

			var source, target, note string

			switch from {
			case fromSource:
				note, e = dotstrings.StringsUnescape(m.Ctx)
				if e != nil {
					errChan <- fmt.Errorf("Failed to strings unescape Message.Ctx for string %d (%v)", n+1, e)
					return
				}

				note = xliff.XMLEscapeLoose(note)

				source, e = dotstrings.StringsUnescape(m.Str)
				if e != nil {
					errChan <- fmt.Errorf("Failed to strings unescape Message.Str for string %d (%v)", n+1, e)
					return
				}

				source = xliff.XMLEscapeLoose(source)

				// if hasTargetLanguage {
				// 	target = source
				// }

			case fromTarget:
				source, e = dotstrings.StringsUnescape(m.Ctx)
				if e != nil {
					errChan <- fmt.Errorf("Failed to strings unescape Message.Ctx for string %d (%v)", n+1, e)
					return
				}

				source = xliff.XMLEscapeLoose(source)

				if hasTargetLanguage {

					target, e = dotstrings.StringsUnescape(m.Str)
					if e != nil {
						errChan <- fmt.Errorf("Failed to strings unescape Message.Str for string %d (%v)", n+1, e)
						return
					}

					target = xliff.XMLEscapeLoose(target)
				}
			}

			dstChan <- xliff.TranslationUnit{File: tf, ID: id, Source: source, Target: target, Note: note}
			n++
		}
	}

	go translator(from, msgChan, tf, dstChan, errChan)
	return dstChan, errChan
}

// ConvertSourceMessagesToTranslationUnits will convert a channel containing dotstrings
// Messages into XLIFF translation units.
// If the passed in translation file has a TargetLanguage set then a translation unit
// will also contain the Str field from the Message copied into Target field.
func ConvertSourceMessagesToTranslationUnits(srcChan <-chan dotstrings.Message, tf *xliff.TranslationFile) (<-chan xliff.TranslationUnit, <-chan error) {
	return convertMessagesToTranslationUnits(fromSource, srcChan, tf)
}

// ConvertTargetMessagesToTranslationUnits will convert a channel containing dotstrings
// Messages into XLIFF translation units.
// If the passed in translation file has a TargetLanguage set then a translation unit
// will also contain the Str field from the Message copied into Target field.
func ConvertTargetMessagesToTranslationUnits(tgtChan <-chan dotstrings.Message, tf *xliff.TranslationFile) (<-chan xliff.TranslationUnit, <-chan error) {
	return convertMessagesToTranslationUnits(fromTarget, tgtChan, tf)
}

// ConvertTranslationUnitsToSourceMessages will take ID, Source and Note fields of a translation unit and create a message out of it where the
// Note is used as the Ctx, the ID as the ID and the Source as the Str. The channel of messages can then be save to a source .strings file.
func ConvertTranslationUnitsToSourceMessages(xliffChan <-chan xliff.TranslationUnit) <-chan dotstrings.Message {

	msgChan := make(chan dotstrings.Message, 3)

	converter := func(xliffChan <-chan xliff.TranslationUnit, msgChan chan<- dotstrings.Message) {
		defer close(msgChan)
		for x := range xliffChan {
			id := dotstrings.StringsEscape(xliff.XMLUnescape(x.ID))
			source := dotstrings.StringsEscape(xliff.XMLUnescape(x.Source))
			note := dotstrings.StringsEscape(xliff.XMLUnescape(x.Note))
			msgChan <- dotstrings.Message{ID: id, Str: source, Ctx: note}
		}

	}
	go converter(xliffChan, msgChan)
	return msgChan
}

// ConvertTranslationUnitsToTargetMessages will take ID, Source and Target fields of a translation unit and create a message out of it where the
// Source is used as the Ctx, the ID as the ID and the Target as the Str. The channel of messages can then be save to a target .strings file.
func ConvertTranslationUnitsToTargetMessages(xliffChan <-chan xliff.TranslationUnit) <-chan dotstrings.Message {

	msgChan := make(chan dotstrings.Message, 3)

	converter := func(xliffChan <-chan xliff.TranslationUnit, msgChan chan<- dotstrings.Message) {
		defer close(msgChan)
		for x := range xliffChan {
			id := dotstrings.StringsEscape(xliff.XMLUnescape(x.ID))
			source := dotstrings.StringsEscape(xliff.XMLUnescape(x.Source))
			target := dotstrings.StringsEscape(xliff.XMLUnescape(x.Target))
			msgChan <- dotstrings.Message{ID: id, Str: target, Ctx: source}
		}

	}
	go converter(xliffChan, msgChan)
	return msgChan
}
