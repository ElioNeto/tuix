// Example: Image Rendering
//
// This example demonstrates image rendering using the Kitty/Sixel terminal
// graphics protocol. If your terminal supports it, images will display inline.
// Falls back to ASCII art if the protocol is not available.
package main

import (
	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>🖼️ Image Rendering</h1>

			<div class="card">
				<h2>Kitty / Sixel Protocol Image</h2>
				<img src="./examples/ascii-image/testcard.png" width="40" height="20" alt="Test Card" />
			</div>

			<div class="card">
				<h2>Image Info</h2>
				<p>Source: examples/ascii-image/testcard.png</p>
				<p>Terminals with image support: Kitty, WezTerm, iTerm2, xterm with Sixel</p>
				<p>Fallback: ASCII art conversion when no image protocol is available</p>
			</div>

			<p class="hint">q to quit</p>
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
			margin-bottom: 1;
		}
		h2 {
			color: #e94560;
			margin-bottom: 1;
		}
		.card {
			margin-bottom: 1;
		}
		.hint {
			text-align: center;
			color: #555;
			margin-top: 1;
		}
		img {
			alt: "Test Card Image";
		}
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	app.Run()
}
