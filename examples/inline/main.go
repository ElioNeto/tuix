// Example: Inline Formatting Context
//
// This example shows inline elements flowing within a block container,
// with inline text wrapping, spans with different colors/styles,
// and mixed inline/block content.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Inline Formatting Demo</h1>

			<div class="section">
				<h2>Text with inline spans</h2>
				<p>
					This is a paragraph with
					<span class="bold">bold text</span>,
					<span class="italic">italic text</span>,
					<span class="color">colored text</span>,
					and <span class="highlight">highlighted</span> words.
				</p>
			</div>

			<div class="section">
				<h2>Navigation bar (inline links)</h2>
				<div class="nav">
					<a class="nav-item active">Home</a>
					<a class="nav-item">About</a>
					<a class="nav-item">Services</a>
					<a class="nav-item">Contact</a>
				</div>
			</div>

			<div class="section">
				<h2>Mixed inline and block</h2>
				<p>Some text before the divider.</p>
				<div class="divider"></div>
				<p>Some text after the <strong>divider</strong> with <em>emphasis</em>.</p>
			</div>

			<div class="section">
				<h2>Tag badges</h2>
				<div class="badge-row">
					<span class="badge red">Go</span>
					<span class="badge green">TUI</span>
					<span class="badge blue">HTML</span>
					<span class="badge orange">CSS</span>
					<span class="badge purple">Flexbox</span>
				</div>
			</div>

			<p class="hint">Press q to quit</p>
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
		.section {
			margin-bottom: 1;
			padding: 1;
			border: solid #0f3460;
		}
		p {
			margin: 0;
		}

		/* Inline spans */
		.bold { font-weight: bold; color: #ffffff; }
		.italic { font-style: italic; color: #ffd700; }
		.color { color: #00d4aa; }
		.highlight { background-color: #533483; color: #fff; padding: 0; }

		/* Navigation bar */
		.nav {
			background-color: #16213e;
			padding: 0;
		}
		.nav-item {
			padding: 0 2;
			color: #888;
		}
		.nav-item.active {
			color: #00d4aa;
			font-weight: bold;
			background-color: #0f3460;
		}

		/* Divider */
		.divider {
			height: 1;
			background-color: #0f3460;
			margin: 1 0;
		}

		/* Badges */
		.badge-row {
			display: flex;
			flex-direction: row;
			gap: 1;
		}
		.badge {
			padding: 0 1;
			font-weight: bold;
			color: #fff;
		}
		.badge.red    { background-color: #e94560; }
		.badge.green  { background-color: #00d4aa; color: #1a1a2e; }
		.badge.blue   { background-color: #0f3460; border: solid #00d4aa; }
		.badge.orange { background-color: #e07c24; }
		.badge.purple { background-color: #533483; }

		.hint { text-align: center; color: #555; margin-top: 1; }
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
