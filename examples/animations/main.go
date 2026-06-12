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

	app.UseDesignSystem()

	app.SetHTML(`
		<div id="app">
			<h1>🎬 CSS Animations</h1>

			<div class="card">
				<h2>Spinner</h2>
				<div style="margin-bottom: 1;">
					<span class="animate-spin">Loading </span><span class="animate-spin">▌</span>
				</div>
				<p>Class: <code>.animate-spin</code> — cycles through | / - \ characters</p>
			</div>

			<div class="card">
				<h2>Pulse</h2>
				<div style="margin-bottom: 1;">
					<span class="animate-pulse">● Live</span>
				</div>
				<p>Class: <code>.animate-pulse</code> — toggles dim/bright state via [pulsing] attribute</p>
			</div>

			<div class="card">
				<h2>Blink</h2>
				<div style="margin-bottom: 1;">
					<span class="animate-blink" style="color: #e94560;">● Recording</span>
				</div>
				<p>Class: <code>.animate-blink</code> — toggles visible/hidden via [blinking] attribute</p>
			</div>

			<div class="card">
				<h2>Progress (determinate)</h2>
				<progress value="65" max="100"></progress>
			</div>

			<div class="card">
				<h2>Progress (indeterminate)</h2>
				<progress></progress>
			</div>

			<p class="text-center muted" style="margin-top: 1;">Animations tick at 200ms intervals · q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
			color: #0f3460;
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
		code {
			color: #00d4aa;
		}
		.muted {
			color: #555;
		}
		/* Pulse and blink are now text-based (no CSS attribute selectors) */
		.animate-pulse {
			color: #00ff00;
			font-weight: bold;
		}
		.animate-blink {
			color: #00ff00;
		}
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	app.Run()
}
