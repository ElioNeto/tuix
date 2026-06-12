package main

import (
	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.UseDesignSystem()

	app.SetHTML(`
		<div id="app">
			<h1>🎨 Design System</h1>

			<div class="card">
				<h2>Button Variants</h2>
				<div class="flex gap-2" style="margin-bottom: 1;">
					<button class="btn btn-primary">Primary</button>
					<button class="btn btn-secondary">Secondary</button>
					<button class="btn btn-danger">Danger</button>
					<button class="btn btn-ghost">Ghost</button>
				</div>
				<div class="flex gap-2">
					<button class="btn btn-sm">Small</button>
					<button class="btn">Default</button>
					<button class="btn btn-lg">Large</button>
				</div>
			</div>

			<div class="card">
				<h2>Badges</h2>
				<div class="flex gap-2">
					<span class="badge badge-primary">Primary</span>
					<span class="badge badge-success">Success</span>
					<span class="badge badge-warning">Warning</span>
					<span class="badge badge-error">Error</span>
				</div>
			</div>

			<div class="card">
				<h2>Form Elements</h2>
				<div class="field">
					<label>Text Input</label>
					<input type="text" class="input" placeholder="Type something..." />
				</div>
				<div class="field">
					<label>Colored Input</label>
					<input type="color" value="#FF6600" />
				</div>
				<div class="field">
					<label>Range</label>
					<input type="range" value="60" />
				</div>
			</div>

			<p class="text-center text-muted">Tab to navigate · q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
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
		.field {
			margin-bottom: 1;
		}
		label {
			display: block;
			color: #555;
			font-weight: bold;
			margin-bottom: 1;
		}
		.text-muted {
			color: #555;
		}
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	app.Run()
}
