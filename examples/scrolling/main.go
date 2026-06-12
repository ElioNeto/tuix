// Example: Scrolling and Overflow
//
// This example demonstrates scrolling content areas.
// Use arrow keys, PageUp/PageDown, Home/End to scroll.
// Mouse wheel is also supported.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Scrolling & Overflow Demo</h1>

			<div class="row">
				<div class="panel">
					<h2>Tall List (needs scroll)</h2>
					<div class="scroll-area">
						<div class="list-item">Item 1: Apples</div>
						<div class="list-item">Item 2: Bananas</div>
						<div class="list-item">Item 3: Cherries</div>
						<div class="list-item">Item 4: Dates</div>
						<div class="list-item">Item 5: Elderberries</div>
						<div class="list-item">Item 6: Figs</div>
						<div class="list-item">Item 7: Grapes</div>
						<div class="list-item">Item 8: Honeydew</div>
						<div class="list-item">Item 9: Kiwi</div>
						<div class="list-item">Item 10: Lemons</div>
						<div class="list-item">Item 11: Mangoes</div>
						<div class="list-item">Item 12: Nectarines</div>
					</div>
				</div>

				<div class="panel">
					<h2>Chat Log (needs scroll)</h2>
					<div class="scroll-area">
						<div class="chat-msg"><span class="user">Alice:</span> Hey everyone!</div>
						<div class="chat-msg"><span class="user">Bob:</span> Hi Alice!</div>
						<div class="chat-msg"><span class="user">Charlie:</span> Morning!</div>
						<div class="chat-msg"><span class="user">Alice:</span> How's the project going?</div>
						<div class="chat-msg"><span class="user">Bob:</span> Almost done!</div>
						<div class="chat-msg"><span class="user">Charlie:</span> Just fixing some bugs</div>
						<div class="chat-msg"><span class="user">Alice:</span> Great work team!</div>
						<div class="chat-msg"><span class="user">Bob:</span> Thanks!</div>
						<div class="chat-msg"><span class="user">Charlie:</span> 🎉</div>
					</div>
				</div>
			</div>

			<p class="hint">Press q to quit. ↑↓ PgUp PgDn Home End to scroll. Mouse wheel also works!</p>
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

		/* Side-by-side panels */
		.row {
			display: flex;
			flex-direction: row;
			gap: 1;
		}
		.panel {
			flex: 1;
			padding: 1;
			border: solid #0f3460;
			background-color: #16213e;
		}

		/* Scrollable area */
		.scroll-area {
			overflow-y: scroll;
			height: 10;
		}

		/* List items */
		.list-item {
			padding: 0 1;
			color: #c0c0c0;
		}
		.list-item:nth-child(odd) {
			background-color: #1a1a2e;
		}
		.list-item:nth-child(even) {
			background-color: #16213e;
		}

		/* Chat messages */
		.chat-msg {
			padding: 0 1;
		}
		.user {
			font-weight: bold;
			color: #00d4aa;
		}

		.hint { text-align: center; color: #555; margin-top: 1; }
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
