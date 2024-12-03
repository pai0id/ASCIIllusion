package ui

import "github.com/marcusolsson/tui-go"

type EntryField struct {
	*tui.Entry
	onSubmit func(text string)
}

func NewEntryField(onSubmit func(text string)) *EntryField {
	entry := tui.NewEntry()
	entry.SetFocused(true)

	field := &EntryField{
		Entry:    entry,
		onSubmit: onSubmit,
	}

	entry.OnSubmit(func(e *tui.Entry) {
		if onSubmit != nil {
			onSubmit(e.Text())
		}
		e.SetText("")
	})

	return field
}
