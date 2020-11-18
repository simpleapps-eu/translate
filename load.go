package translate

import (
	"bufio"
	"io"
)

// LoadLines reads the text stream from io.Reader and outputs lines on
// a channel it returns. This function will run asynchronously and return
// before the whole text stream has been processed.
func LoadLines(srcFile io.Reader) (<-chan string, <-chan error) {

	lineChan := make(chan string, 3)
	errChan := make(chan error, 1)

	reader := func(srcFile io.Reader, lineChan chan<- string, errChan chan<- error) {
		defer close(lineChan)
		defer close(errChan)
		s := bufio.NewScanner(srcFile)
		for s.Scan() {
			lineChan <- s.Text()
		}
		if err := s.Err(); err != nil {
			errChan <- err
		}
	}

	go reader(srcFile, lineChan, errChan)
	return lineChan, errChan
}
