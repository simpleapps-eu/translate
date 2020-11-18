package dotstrings

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// LoadMessagesMapFromFile uses the given filename to open the messages file
// and reads all messages from it. The function returns a map with the messages
// once all messages have been read.
func LoadMessagesMapFromFile(filename string) (messages map[string]Message, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	messages, err = LoadMessagesMap(NewReaderUTF16(file))
	return
}

// LoadMessagesMap will read all the entries from the messages file.
// If it encounters an error it will return with the error instead of continuing.
// The function returns a map with the messages once all messages have been read.
func LoadMessagesMap(tmReader io.Reader) (messages map[string]Message, err error) {
	messages = make(map[string]Message)
	msgChan, errChan := LoadMessages(tmReader)
	for tm := range msgChan {
		if tm.Fuzzy {
			err = fmt.Errorf("Encountered a fuzzy translation for ID %q", tm.ID)
			return
		}
		if _, present := messages[tm.ID]; present {
			err = fmt.Errorf("Encountered a duplicated ID %q", tm.ID)
			return
		}
		messages[tm.ID] = tm
	}
	err, _ = <-errChan
	return
}

// LoadMessages reads the data provided by io.Reader and outputs messages on
// a channel it returns. This function will run asynchronously and return before
// the whole stream has been processed.
func LoadMessages(r io.Reader) (<-chan Message, <-chan error) {
	c := make(chan Message, 3)
	e := make(chan error, 1)

	s := bufio.NewScanner(r)
	s.Split(Split())

	reader := func(outChan chan<- Message, s *bufio.Scanner, errChan chan<- error) {
		defer close(outChan)
		defer close(errChan)
		for s.Scan() {
			m := Message{}
			if IsFuzzyToken(s.Text()) {
				m.Fuzzy = true
				if !s.Scan() {
					continue
				}
			}
			m.Ctx = s.Text()
			if s.Scan() {
				m.ID = s.Text()
				if s.Scan() {
					m.Str = s.Text()
					outChan <- m
				}
			}
		}
		if e := s.Err(); e != nil {
			errChan <- e
		}
	}

	go reader(c, s, e)
	return c, e
}
