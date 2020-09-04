package mailer

import "unapu.com/mailer"

var (
	Fallback = mailer.Fallback
	WithReporter = mailer.WithReporter
)

type (
	Email           = mailer.Email
	SenderInterface = mailer.Sender
	SenderFunc      = mailer.SenderFunc
	SendHandler     = mailer.SendHandler
	SendhandlerFunc = mailer.SendHandlerFunc
)
