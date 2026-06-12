// Example: Hover Effects
//
// This example demonstrates hover effects using the :hover pseudo-class.
// Move your mouse over elements to see hover styles applied!
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Hover Effects Demo</h1>

			<div class="section">
				<h2>Buttons</h2>
				<div class="button-group">
					<button class="btn">Submit</button>
					<button class="btn danger">Delete</button>
					<button class="btn ghost">Cancel</button>
					<button class="btn outline">Learn More</button>
				</div>
			</div>

			<div class="section">
				<h2>Cards</h2>
				<div class="card-grid">
					<div class="card">
						<div class="card-title">Dashboard</div>
						<div class="card-desc">View analytics and reports</div>
					</div>
					<div class="card">
						<div class="card-title">Team</div>
						<div class="card-desc">Manage team members</div>
					</div>
					<div class="card">
						<div class="card-title">Settings</div>
						<div class="card-desc">Configure preferences</div>
					</div>
					<div class="card">
						<div class="card-title">Files</div>
						<div class="card-desc">Browse and upload files</div>
					</div>
				</div>
			</div>

			<div class="section">
				<h2>Navigation Menu</h2>
				<div class="menu">
					<div class="menu-item active">Home</div>
					<div class="menu-item">Profile</div>
					<div class="menu-item">Projects</div>
					<div class="menu-item">Reports</div>
					<div class="menu-item disabled">Archived</div>
				</div>
			</div>

			<div class="section">
				<h2>Link List</h2>
				<div class="link-list">
					<a class="link">Getting Started Guide</a>
					<a class="link">API Reference</a>
					<a class="link">Troubleshooting</a>
					<a class="link">Community Forum</a>
				</div>
			</div>

			<p class="hint">Hover over elements with your mouse • q to quit</p>
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

		/* Buttons — :hover changes background and text color */
		.button-group {
			display: flex;
			gap: 1;
		}
		.btn {
			padding: 0 3;
			border: solid #00d4aa;
			color: #c0c0c0;
			background-color: transparent;
		}
		.btn:hover {
			background-color: #00d4aa;
			color: #1a1a2e;
			font-weight: bold;
		}
		.btn.danger {
			border-color: #e94560;
			color: #e94560;
		}
		.btn.danger:hover {
			background-color: #e94560;
			color: #fff;
		}
		.btn.ghost {
			border-color: transparent;
		}
		.btn.ghost:hover {
			border-color: #c0c0c0;
		}
		.btn.outline:hover {
			background-color: #0f3460;
		}

		/* Cards — :hover changes border and background */
		.card-grid {
			display: flex;
			flex-wrap: wrap;
			gap: 1;
		}
		.card {
			padding: 1;
			border: solid #0f3460;
			background-color: #16213e;
			width: 14;
			text-align: center;
		}
		.card:hover {
			border: solid #00d4aa;
			background-color: #0f3460;
		}
		.card-title { font-weight: bold; color: #ffffff; }
		.card-desc { color: #888; font-size: 10; }

		/* Menu — :hover highlights the menu item */
		.menu {
			display: flex;
			flex-direction: column;
		}
		.menu-item {
			padding: 0 2;
			color: #c0c0c0;
			border-left: solid transparent;
		}
		.menu-item:hover {
			background-color: #0f3460;
			border-left: solid #00d4aa;
			color: #00d4aa;
		}
		.menu-item.active {
			color: #00d4aa;
			font-weight: bold;
		}
		.menu-item.disabled {
			color: #555;
		}
		.menu-item.disabled:hover {
			background-color: transparent;
			border-left: solid transparent;
			color: #555;
		}

		/* Links — :hover adds background and brightens text */
		.link-list {
			display: flex;
			flex-direction: column;
		}
		.link {
			color: #00d4aa;
			padding: 0 1;
		}
		.link:hover {
			background-color: #0f3460;
			color: #ffffff;
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
