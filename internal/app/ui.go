package app

type UI interface {
	Dialogs
}

type Dialogs interface {
	ShowModal(title, message string)
	Info(title, message string)
	Error(title, message string)
	Warning(title, message string)
	Confirm(title, message string) bool
	Input(title, message string) string
	Password(title, message string) string
	Select(title, message string, options []string) string
	SelectMultiple(title, message string, options []string) []string
}
