package dotstrings

// Message contains the information of a single .strings file entry.
//
// " Ctx "
// "ID" = "Str";
//
// There are 2 types of strings files; "source" and "target" strings files.
// The source type contains the original language of the translation project
// and there is supposed to be only 1 source. A target strings file contains
// a specific translation to another language.
//
// For a source .strings file the Ctx contains a note for the translator.
// For a target .strings file the every Ctx contains the matching Str from the
// source .strings file. This is done so that by comparing source.Str with
// target.Ctx a change in the original translation string can be determined.
type Message struct {
	// Fuzzy is true when the source Str is different from the target Ctx value or
	// when the source ID is missing completely from the target file.
	// Set in messages loaded from a strings file that are preceeded with a 
	// comment that contains the word "fuzzy".
	Fuzzy bool
	// Missing is true when the ID was not found in the target file. When true both
	// the target Ctx and Str will contain the source Str.
	// Set in messages emited by the TranslateMessages function.
	Missing bool
	Ctx     string
	ID      string
	Str     string
}
