//
// peanut - an experiment in tui (Text User Interface) fruits and seeds
//
// Copyright (c) 2018, R. S. Doiel
// All rights not granted herein are expressly reserved by R. S. Doiel.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"fmt"
	"os"

	// My packages
	"github.com/rsdoiel/nuts"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"

	// 3rd Party Packages
	"github.com/jroimartin/gocui"
	//"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
)

var (
	description = ``

	examples = ``

	// Standard Options
	showHelp             bool
	showLicense          bool
	showVersion          bool
	inputFName           string
	quiet                bool
	generateMarkdownDocs bool

	// App Options
)

func rootLayout(ui *gocui.Gui) error {
	maxX, maxY := ui.Size()
	if view, err := ui.SetView("Oh nuts!", (maxX/2)-20, maxY/2, (maxX/2)+20, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(view, "We have a peanut! --> %q", err)
	}
	return nil
}

func quit(ui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	// Configure and instanciate the command line interface
	app := cli.NewCli(nuts.Version)

	// Describe the non-option parameters
	app.AddParams("[STARTUP_DIRECTORY]")

	// Add help docs
	app.AddHelp("description", []byte(description))

	// Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.StringVar(&inputFName, "i,input", "", "input filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error output")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "generate markdown docs")

	// Parse and evaluate options before starting up
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out = os.Stdout

	// Process Options
	if showHelp {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, "%s", app.Help(args...))
			os.Exit(0)
		}
		app.Usage(app.Out)
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
		os.Exit(0)
	}

	// Create our root window
	ui, err := gocui.NewGui(gocui.Output256)
	cli.ExitOnError(app.Eout, err, false)
	defer ui.Close()

	// Setup our root window manager
	ui.SetManagerFunc(rootLayout)

	// Add the key/event binding
	err = ui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit)
	cli.ExitOnError(app.Eout, err, false)

	// Startup the main loop
	err = ui.MainLoop()
	cli.ExitOnError(app.Eout, err, false)
}
