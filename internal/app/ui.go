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
	Confirm(ctx context.Context, title, message string) bool
	Input(ctx context.Context, title, footer, message string) string
	Password(ctx context.Context, title, message string) string
	Select(ctx context.Context, title, message string, options []string) string
	SelectMultiple(ctx context.Context, title, message string, options []string) []string
}

type Panel interface {
	SetData(data []AttributeList)
}
