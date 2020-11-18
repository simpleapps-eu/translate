package xliff

// TranslationFile contains meta information about the xliff file.
// There is one entry per xliff file. Every TranslationUnit carries a
// pointer to  its TranslationFile.
type TranslationFile struct {
	Original       string
	SourceLanguage string
	Datatype       string
	TargetLanguage string
}

// TranslationUnit contains information about a single string to be translated.
// There are multiple entries per xliff file.
// Note that text in Source and Target fields is supposed to be properly escaped
// XML text. e.g. the character '&' replaced with '&amp;'
type TranslationUnit struct {
	File   *TranslationFile
	ID     string
	Source string
	Target string
	Note   string
}
