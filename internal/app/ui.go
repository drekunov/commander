package app

import "context"

type UI interface {
	Dialog
	Panel
}

type Dialog interface {
	Info(ctx context.Context, title, footer, message string)
	Error(title, message string)
	Warning(title, message string)
	Confirm(title, message string) bool
	Input(ctx context.Context, title, footer, message string) string
	Password(title, message string) string
	Select(title, message string, options []string) string
	SelectMultiple(title, message string, options []string) []string
}

type Panel interface {
	SetData(data []AttributeList)
}
