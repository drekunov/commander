package app

import "context"

// PanelID identifies one of the two file panels a listing is delivered to.
type PanelID int

const (
	PanelLeft PanelID = iota
	PanelRight
)

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
	SetData(panel PanelID, dir string, data []AttributeList)
}
