package main

import (
	"fmt"
	"os"

	// 3rd Party packages
	"github.com/marcusolsson/tui-go"
)

func main() {
	box := tui.NewVBox(
		tui.NewLabel("nuts-cashew"),
	)
	ui := tui.New(box)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
