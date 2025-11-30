package app

type UI interface { //nolint: iface
	Dialogs
}

type Dialogs interface { //nolint: iface
	ShowModal(title, message string)
	Info(title, footer, message string)
	Error(title, message string)
	Warning(title, message string)
	Confirm(title, message string) bool
	Input(title, footer, message string) string
	Password(title, message string) string
	Select(title, message string, options []string) string
	SelectMultiple(title, message string, options []string) []string
}
