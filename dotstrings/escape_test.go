package dotstrings

import "testing"

var tvStringsEscaped = [...]string{
	`ha&llo`,  //0
	`hallo\n`, //1
	`\thallo`, //2
	`\\test`,  //3
	`\\nest`,  //4
	`hallo\ndit\nis\teen\ntest\\test\"test'test`,    //5
	`ha&llo\nd<it\ni>s\teen\ntest\\test\"test'test`, //6
}

var tvStringsEscapedWrong = [...]string{
	`Couldn\'t authenticate with Dropbox: %s`, //1
}

var tvStringsUnescaped = [...]string{
	`ha&llo`, //0
	`hallo
`, //1
	`	hallo`, //2
	`\test`, //3
	`\nest`, //4
	`hallo
dit
is	een
test\test"test'test`, //5
	`ha&llo
d<it
i>s	een
test\test"test'test`, //6
}

func TestStringsEscapedWrong(t *testing.T) {

	for i := 0; i < len(tvStringsEscapedWrong); i++ {
		_, err := StringsUnescape(tvStringsEscapedWrong[i])
		if err == nil {
			t.Errorf("Expected unescaping tvStringsEscapedWrong[%d] to fail", i)
		}
	}
}

func TestStringsEscaping(t *testing.T) {

	for i := 0; i < len(tvStringsEscaped); i++ {
		s, err := StringsUnescape(tvStringsEscaped[i])
		if err != nil {
			t.Errorf("Error while unquoting UTF8 double quoted string %d (%v)\n", i, err)
		}
		if tvStringsUnescaped[i] != s {
			t.Errorf("Unescaped .strings string %d doesn't match test vector\n", i)
		}
	}

}
