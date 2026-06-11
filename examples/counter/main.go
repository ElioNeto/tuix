package main

import (
	"fmt"
	"log"

	"github.com/elioneto/tuix"
	"github.com/elioneto/tuix/dom"
	"github.com/elioneto/tuix/terminal"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Tuix Counter</h1>
			<div class="counter-box">
				<div class="counter" id="count">0</div>
				<div class="buttons">
					<button id="dec" class="btn">-1</button>
					<button id="inc" class="btn">+1</button>
					<button id="reset" class="btn secondary">Reset</button>
				</div>
			</div>
			<div class="footer">
				Press q or Ctrl+C to quit
			</div>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 10px;
			background-color: #1a1a2e;
			color: #e0e0e0;
		}

		h1 {
			text-align: center;
			color: #00d4aa;
			font-weight: bold;
			margin-bottom: 20px;
		}

		.counter-box {
			text-align: center;
			margin: 20px;
			padding: 20px;
			border: solid;
			border-color: #16213e;
			background-color: #0f3460;
		}

		.counter {
			font-size: 48;
			color: #e94560;
			font-weight: bold;
			margin: 20px;
		}

		.buttons {
			margin: 10px;
		}

		.btn {
			padding: 5px 10px;
			margin: 5px;
			border: solid;
			border-color: #00d4aa;
			color: #00d4aa;
			background-color: #16213e;
		}

		.btn.secondary {
			border-color: #e94560;
			color: #e94560;
		}

		.footer {
			text-align: center;
			color: #888;
			margin-top: 10px;
		}
	`)

	count := 0

	app.OnInit(func() {
		app.SetTitle("Tuix Counter")
	})

	app.OnKey(func(key terminal.Key) {
		switch key {
		case terminal.KeyCtrlC, terminal.KeyRune:
			// Handled by event loop
		}
	})

	app.OnRune(func(r rune) {
		switch r {
		case 'q':
			app.Stop()
		case 'i', '+':
			count++
			updateCount(app, count)
		case 'd', '-':
			count--
			updateCount(app, count)
		case 'r':
			count = 0
			updateCount(app, count)
		}
	})

	app.OnMouse(func(btn terminal.MouseButton, x, y int) {
		box := app.RootBox()
		if box == nil {
			return
		}

		clicked := box.FindBoxAtPoint(x, y)
		if clicked == nil || clicked.Node == nil {
			return
		}

		id := clicked.Node.ID()
		switch id {
		case "inc":
			count++
			updateCount(app, count)
		case "dec":
			count--
			updateCount(app, count)
		case "reset":
			count = 0
			updateCount(app, count)
		}
	})

	app.OnClose(func() {
		fmt.Println("Goodbye!")
	})

	log.Fatal(app.Run())
}

func updateCount(app *tuix.App, count int) {
	// Update the DOM node's text content
	doc := app.Document()
	elements := doc.QuerySelectorAll("#count")
	for _, el := range elements {
		// Clear children and add new text node
		el.Children = nil
		el.AppendChild(dom.Text(fmt.Sprintf("%d", count)))
	}

	// Re-render
	app.Rebuild()
}
