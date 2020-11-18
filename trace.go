package translate

import (
	"fmt"
	"runtime"
)

func trace() {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println(file, ":", line)
	}
}
