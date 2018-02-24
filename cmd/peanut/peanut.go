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
	"io/ioutil"
	"os"
	"strings"

	// My packages
	"github.com/rsdoiel/nuts"
	"github.com/rsdoiel/opml"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"

	// 3rd Party Packages
	"github.com/rivo/tview"
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

	// App Globals
	UI *tview.Application
)

func manageList(cursor int, list *tview.List, outline opml.OutlineList) (int, int) {
	list.Clear()
	_, _, _, height := list.Box.GetInnerRect()
	//list.AddItem("Subscribe", "", 's', nil)
	//list.AddItem("View All", "", 'a', nil)
	list.AddItem("Next Page", "", 'n', nil)
	list.AddItem("Prev Page", "", 'p', nil)
	list.AddItem("Quit", "", 'q', func() {
		UI.Stop()
	})
	for i := cursor; i < len(outline) && i < height; i++ {
		o := outline[i]
		list.AddItem(strings.TrimSpace(o.Text), o.XMLURL, 0, nil)
	}
	prev := cursor - height
	if prev < 0 {
		prev = 0
	}
	next := cursor + height
	if next >= len(outline) {
		next = len(outline) - height
	}
	return next, prev
}

func main() {
	var (
		next int
		prev int
	)
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

	// Handle inputFName as first of args
	if len(args) > 0 {
		inputFName = args[0]
	}

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

	// Bring in our OPML data so we can populate the list
	src, err := ioutil.ReadAll(app.In)
	cli.ExitOnError(app.Eout, err, quiet)

	// Unmarshal our data
	doc := opml.New()
	err = opml.Unmarshal(src, doc)
	cli.ExitOnError(app.Eout, err, quiet)

	// Setup our TUI environment
	UI = tview.NewApplication()

	list := tview.NewList()
	list.ShowSecondaryText(false)

	UI.SetRoot(list, true).SetFocus(list)

	list.SetSelectedFunc(func(i int, longText, shortText string, r rune) {
		switch longText {
		case "Subscribe":
		case "Next Page":
			next, prev = manageList(next, list, doc.Body.Outline)
			fmt.Printf("DEBUG next (%d), prev (%d)\n", next, prev)
		case "Prev Page":
			next, prev = manageList(prev, list, doc.Body.Outline)
			fmt.Printf("DEBUG next (%d), prev (%d)\n", next, prev)
		case "Quit":
			UI.Stop()
			fmt.Println("All Done!")
		default:
			UI.Stop()
			fmt.Printf("%d %q, %q, %c\n", i, longText, shortText, r)
		}
	})

	// NOTE: Look through list of feeds in opml document and display in feedList upto the height
	if doc.Body != nil && len(doc.Body.Outline) > 0 {
		next, prev = manageList(0, list, doc.Body.Outline)
		fmt.Printf("DEBUG next (%d), prev (%d)\n", next, prev)
	} else {
		list.AddItem("Quiet", "", 'q', func() {
			UI.Stop()
		})
	}

	// Set root window, focus and run.
	err = UI.Run()
	cli.ExitOnError(app.Eout, err, quiet)
}
