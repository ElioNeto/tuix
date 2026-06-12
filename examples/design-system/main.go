package main

import (
	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.UseDesignSystem()

	app.SetHTML(`
		<div id="app">
			<div class="navbar">
				<span class="nav-brand">◆ Tuix DS</span>
				<span class="nav-item">Overview</span>
				<span class="nav-item">Components</span>
				<span class="nav-item">Forms</span>
			</div>

			<h1>🎨 Design System</h1>

			<div class="grid-2" style="margin-bottom: 1;">
				<div class="card">
					<h2>Buttons</h2>
					<div class="flex gap-1" style="margin-bottom: 1;">
						<button class="btn btn-primary">Primary</button>
						<button class="btn btn-secondary">Secondary</button>
					</div>
					<div class="flex gap-1">
						<button class="btn btn-danger">Danger</button>
						<button class="btn btn-ghost">Ghost</button>
					</div>
				</div>

				<div class="card">
					<h2>Badges</h2>
					<div class="flex gap-1" style="margin-bottom: 1;">
						<span class="badge badge-primary">New</span>
						<span class="badge badge-success">Done</span>
					</div>
					<div class="flex gap-1">
						<span class="badge badge-warning">Pending</span>
						<span class="badge badge-error">Failed</span>
					</div>
				</div>
			</div>

			<div class="card">
				<h2>Tabs</h2>
				<div class="tabs">
					<span class="tab tab-active">Overview</span>
					<span class="tab">Settings</span>
					<span class="tab">Profile</span>
				</div>
				<p style="margin-top: 1;">Tab content goes here.</p>
			</div>

			<div class="grid-2" style="margin-bottom: 1;">
				<div class="card">
					<h2>List Group</h2>
					<div class="list">
						<div class="list-item">Dashboard</div>
						<div class="list-item">Analytics</div>
						<div class="list-item">Reports</div>
						<div class="list-item">Settings</div>
					</div>
				</div>

				<div class="card">
					<h2>Form Controls</h2>
					<div class="field">
						<label>Text</label>
						<input type="text" class="input" placeholder="Enter name..." />
					</div>
					<div class="field">
						<label>Range</label>
						<input type="range" value="60" />
					</div>
					<div class="field">
						<label>Color</label>
						<input type="color" value="#FF6600" />
					</div>
				</div>
			</div>

			<div class="card">
				<h2>Table</h2>
				<div class="table">
					<div class="table-header">Name          Role         Status  </div>
					<div class="table-row">Alice         Admin        Active  </div>
					<div class="table-row">Bob           Editor       Active  </div>
					<div class="table-row">Charlie       Viewer       Inactive</div>
				</div>
			</div>

			<div class="flex justify-center gap-2" style="margin-top: 1;">
				<button class="btn btn-primary" id="theme-toggle">Toggle Theme (T)</button>
			</div>

			<p class="text-center muted" style="margin-top: 1;">Tab to navigate · t to toggle theme · q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
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
			margin-bottom: 0;
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
		.muted {
			color: #555;
		}
	`)

	// Theme toggle
	isDark := true
	app.OnRune(func(r rune) {
		switch r {
		case 'q':
			app.Stop()
		case 't', 'T':
			if isDark {
				app.SetTheme(tuix.DefaultLightTheme)
				isDark = false
			} else {
				app.SetTheme(tuix.DefaultDarkTheme)
				isDark = true
			}
		}
	})

	app.Run()
}
