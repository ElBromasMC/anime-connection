package service

import (
	"github.com/wneessen/go-mail"
)

type Email struct {
	client *mail.Client
}

func NewEmailService(client *mail.Client) Email {
	return Email{
		client: client,
	}
}
