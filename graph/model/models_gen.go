// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Node interface {
	IsNode()
}

type CloseThread struct {
	ThreadID string `json:"threadID"`
}

type GetThreadsInput struct {
	OffsetTime *string `json:"offsetTime"`
	Closed     *bool   `json:"closed"`
}

type NewThread struct {
	Title string `json:"title"`
}

type OpenThread struct {
	ThreadID string `json:"threadID"`
}

type SignUp struct {
	Name string `json:"name"`
}

type Thread struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Closed    bool   `json:"closed"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (Thread) IsNode() {}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (User) IsNode() {}
