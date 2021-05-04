package elsch

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	threadIndexName = "thread_test"

	clt, err := NewClient()
	if err != nil {
		panic(err)
	}
	if err := CreateIndices(clt); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := DeleteIndices(clt); err != nil {
		panic(err)
	}
	os.Exit(code)
}
