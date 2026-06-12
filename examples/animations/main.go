// Example: CSS Animations
//
// This example demonstrates the built-in animation system:
// - .animate-spin: Rotating spinner character (|/-\)
// - .animate-pulse: Alternating dim/bright state
// - .animate-blink: Toggle visible/hidden state
package main

import (
	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>🎬 CSS Animations</h1>

			<div class="card">
				<h2>Spinner</h2>
				<span>Loading </span><span class="animate-spin">▌</span>
				<span> </span>
				<p>Class: <code>.animate-spin</code> — cycles through | / - \ characters</p>
			</div>

			<div class="card">
				<h2>Pulse</h2>
				<span class="animate-pulse">● Live</span>
				<p>Class: <code>.animate-pulse</code> — toggles dim/bright state via [pulsing] attribute</p>
			</div>

			<div class="card">
				<h2>Blink</h2>
				<span class="animate-blink">● Recording</span>
				<p>Class: <code>.animate-blink</code> — toggles visible/hidden via [blinking] attribute</p>
			</div>

			<div class="card">
				<h2>Progress (determinate)</h2>
				<progress value="65" max="100"></progress>
				<progress value="15" max="100"></progress>
				<progress value="90" max="100"></progress>
			</div>

			<div class="card">
				<h2>Progress (indeterminate)</h2>
				<progress></progress>
			</div>

			<p class="text-center muted">Animations tick at 200ms intervals · q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
			color: #c0c0c0;
			width: 100%;
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
			margin: 1;
			border: solid;
			border-color: #00d444;
			padding: 1;
		}
		code {
			color: #00d4aa;
		}
		.muted {
			color: #555;
		}
		.animate-spin {
			font-weight: bold;
		}
		.animate-pulse {
			color: #661dbfff;
		}
		.animate-pulse[pulsing] {
			color: #00a2ffff;
			font-weight: bold;
		}
		.animate-blink {
			color: #c70d16ff;
		}
		.animate-blink[blinking] {
			color: #dee3deff;
		}
	
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	app.Run()
}
