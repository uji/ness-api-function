package graph

import (
	"time"

	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
	"github.com/uji/ness-api-function/graph/model"
)

func cnvUser(user *usr.User) *model.User {
	return &model.User{
		ID:        string(user.UserID()),
		Name:      string(user.Name()),
		CreatedAt: user.CreatedAt().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt().Format(time.RFC3339),
	}
}

func cnvThread(thread thread.Thread) *model.Thread {
	return &model.Thread{
		ID:        thread.ID(),
		Title:     thread.Title(),
		Closed:    thread.Closed(),
		CreatedAt: thread.CreatedAt().Format(time.RFC3339),
		UpdatedAt: thread.UpdatedAt().Format(time.RFC3339),
	}
}
