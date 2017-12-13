package main

import (
	"fmt"
	"os"

	// 3rd Party packages
	"github.com/marcusolsson/tui-go"
)

func main() {
	quitButton := tui.NewButton("Quit")
	box := tui.NewVBox(
		tui.NewLabel("nuts-cashew"),
		quitButton,
	)
	ui := tui.New(box)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	quitButton.OnActivated(func(b *tui.Button) { ui.Quit() })
	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
