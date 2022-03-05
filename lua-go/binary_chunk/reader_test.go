package binary_chunk

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func errorHandler(t *testing.T) {
	if a := recover(); a != nil {
		t.Fatal(a)
	}
}

func TestHeaderCheck(t *testing.T) {
	defer errorHandler(t)
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		t.Fatal("read binary file fail")
	}
	fmt.Print(Undump(data))
}
