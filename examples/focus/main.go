// Example: Focus Management
//
// This example demonstrates focusable elements with tabindex,
// auto-focus, and focus callbacks.
package main

import (
	"log"

	"github.com/elioneto/tuix"
	"github.com/elioneto/tuix/dom"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Focus Management Demo</h1>

			<div id="status" class="status">Focus: (none)</div>

			<div class="section">
				<h2>Tab Navigation</h2>
				<div class="tab-bar">
					<div class="tab focused">General</div>
					<div class="tab">Appearance</div>
					<div class="tab">Advanced</div>
					<div class="tab">About</div>
				</div>

				<div class="tab-content">
					<h3>General Settings</h3>
					<div class="focused-item">
						<div class="item-label">Theme</div>
						<div class="item-value" tabindex="0" autofocus>Dark</div>
					</div>
					<div class="focused-item">
						<div class="item-label">Language</div>
						<div class="item-value" tabindex="0">English</div>
					</div>
					<div class="focused-item">
						<div class="item-label">Font Size</div>
						<div class="item-value" tabindex="0">14px</div>
					</div>
				</div>
			</div>

			<div class="section">
				<h2>Focus Indicators</h2>
				<div class="focus-grid">
					<div class="focus-card" tabindex="0">
						<div class="card-title">Documents</div>
						<div class="card-count">24 items</div>
					</div>
					<div class="focus-card" tabindex="0">
						<div class="card-title">Images</div>
						<div class="card-count">12 items</div>
					</div>
					<div class="focus-card" tabindex="0">
						<div class="card-title">Music</div>
						<div class="card-count">8 items</div>
					</div>
					<div class="focus-card" tabindex="0">
						<div class="card-title">Notes</div>
						<div class="card-count">3 items</div>
					</div>
				</div>
			</div>

			<div class="section">
				<h2>Form Controls</h2>
				<input type="text" value="Focus me" />
				<input type="text" value="Tab between us" />
				<button>Submit</button>
				<select>
					<option>Option A</option>
					<option>Option B</option>
					<option>Option C</option>
				</select>
			</div>

			<p class="hint">Tab/Shift+Tab to navigate • Enter to activate • q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
			color: #c0c0c0;
		}
		h1 {
			color: #00d4aa;
			text-align: center;
			font-weight: bold;
			margin-bottom: 1;
		}
		h2 {
			color: #e94560;
			font-size: 12;
			margin: 0;
			margin-bottom: 1;
		}
		h3 {
			color: #00d4aa;
			margin: 0;
		}
		.section {
			margin-bottom: 1;
			padding: 1;
			border: solid #0f3460;
			background-color: #16213e;
		}

		/* Status bar */
		.status {
			padding: 0 1;
			margin-bottom: 1;
			background-color: #0f3460;
			color: #00d4aa;
		}

		/* Tab bar */
		.tab-bar {
			display: flex;
			gap: 0;
			margin-bottom: 1;
		}
		.tab {
			padding: 0 2;
			background-color: #0f3460;
			color: #888;
		}
		.tab.focused {
			background-color: #00d4aa;
			color: #1a1a2e;
			font-weight: bold;
		}

		/* Focused items */
		.focused-item {
			display: flex;
			align-items: center;
			margin: 0;
		}
		.item-label {
			color: #888;
			width: 14;
		}
		.item-value {
			padding: 0 2;
			color: #c0c0c0;
			border: solid transparent;
		}
		.item-value[focused] {
			border: solid #00d4aa;
			background-color: #0f3460;
			color: #00d4aa;
		}

		/* Focus grid */
		.focus-grid {
			display: flex;
			flex-wrap: wrap;
			gap: 1;
		}
		.focus-card {
			padding: 1;
			width: 12;
			border: solid #0f3460;
			text-align: center;
		}
		.focus-card[focused] {
			border: solid #00d4aa;
			background-color: #0f3460;
		}
		.card-title {
			font-weight: bold;
			color: #c0c0c0;
		}
		.card-count {
			color: #888;
		}

		/* Form controls */
		input {
			margin-bottom: 1;
			width: 30;
		}
		button {
			margin-bottom: 1;
		}
		select {
			margin-bottom: 1;
		}

		.hint { text-align: center; color: #555; margin-top: 1; }
	`)

	// Track focus changes in a status element
	var statusEl *dom.Node

	app.OnInit(func() {
		// Find the status element for updating
		els := app.Document().QuerySelectorAll("#status")
		if len(els) > 0 {
			statusEl = els[0]
		}
	})

	app.OnFocus(func(el *dom.Node) {
		if statusEl != nil {
			tag := el.Data
			id := el.GetAttribute("id")
			class := el.GetAttribute("class")
			text := "Focused: <" + tag + ">"
			if id != "" {
				text += " #" + id
			}
			if class != "" {
				text += " ." + class
			}
			// Update the status element's text via DOM manipulation
			statusEl.Children = nil
			statusEl.Children = append(statusEl.Children, &dom.Node{Type: dom.NodeText, Data: text, Parent: statusEl})
		}
	})

	app.OnBlur(func(el *dom.Node) {
		if statusEl != nil {
			tag := el.Data
			id := el.GetAttribute("id")
			text := "Blurred: <" + tag + ">"
			if id != "" {
				text += " #" + id
			}
		}
	})

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
