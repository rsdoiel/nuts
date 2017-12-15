package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	// 3rd Party packages
	"github.com/marcusolsson/tui-go"
)

func main() {
	args := os.Args[:]
	args[0] = path.Base(os.Args[0])

	box := tui.NewVBox(
		tui.NewLabel(strings.Join(args, " ")),
	)
	ui := tui.New(box)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
