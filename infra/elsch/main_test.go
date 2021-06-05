package elsch

import (
	"testing"
)

func CreateIndexForTest(t *testing.T, c *Client) {
	if err := CreateIndices(c); err != nil {
		t.Fatal(err)
	}
}

func DeleteIndexForTest(t *testing.T, c *Client) {
	if err := DeleteIndices(c); err != nil {
		t.Fatal(err)
	}
}
