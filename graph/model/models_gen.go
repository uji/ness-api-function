// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewThread struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Thread struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Closed      bool   `json:"closed"`
}
