// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewThread struct {
	Title string `json:"title"`
}

type Thread struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Closed bool   `json:"closed"`
}
