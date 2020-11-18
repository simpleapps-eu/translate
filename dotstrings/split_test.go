package dotstrings

import (
	"testing"

	"bufio"
	"fmt"
	"strings"
)

func ExpectSuccess(c bool, err func(e string)) {
	if !c {
		err("Expected Success\n")
	}
}

func ExpectFailure(c bool, err func(e string)) {
	if c {
		err("Expected Failure\n")
	}
}

func ExpectEqual(a, b string, err func(e string)) {
	if a != b {
		err(fmt.Sprintf("Expected string %q to be equal to %q\n", a, b))
	}
}

func TestDotStringsBasic(t *testing.T) {

	const dotStrings = `
	/* Message Context */
	"Message Identifier" = "Message String";

	/* Message Context */
	"Message Identifier" = "Message String";

	/* Message Context */
	'Message Identifier" = "Message String";
	`
	scanner := bufio.NewScanner(strings.NewReader(dotStrings))
	scanner.Split(Split())

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Identifier", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message String", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Identifier", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message String", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })

	ExpectFailure(scanner.Scan(), func(e string) { t.Error(e) })

	// err := scanner.Err()
	// if err != nil {
	// 	t.Error(err.Error())
	// }
}

func TestDotStringsFuzzy(t *testing.T) {
	const dotStrings = `
	/* Fuzzy */
	/* Message Context */
	"Message Identifier" = "Message String";

	/* Message Context */
	"Message Identifier" = "Message String";

	/* fuzzy */
	/* Message Context */
	'Message Identifier" = "Message String";
	`
	scanner := bufio.NewScanner(strings.NewReader(dotStrings))
	scanner.Split(Split())

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Fuzzy", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Identifier", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message String", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Identifier", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message String", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "fuzzy", func(e string) { t.Error(e) })

	ExpectSuccess(scanner.Scan(), func(e string) { t.Error(e) })
	ExpectEqual(scanner.Text(), "Message Context", func(e string) { t.Error(e) })
}
