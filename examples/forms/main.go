// Example: Form Controls (Input Fields)
//
// This example demonstrates interactive form controls.
// Tab between fields, type text, toggle checkboxes and radios,
// cycle select options, and activate buttons!
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Form Controls Demo</h1>

			<form id="login-form">
				<div class="field">
					<label for="username">Username</label>
					<input id="username" type="text" placeholder="Enter your username" />
				</div>

				<div class="field">
					<label for="email">Email</label>
					<input id="email" type="email" placeholder="user@example.com" />
				</div>

				<div class="field">
					<label for="password">Password</label>
					<input id="password" type="password" placeholder="••••••••" />
				</div>

				<div class="field-row">
					<div class="field checkbox">
						<input id="remember" type="checkbox" checked />
						<label for="remember">Remember me</label>
					</div>
				</div>

				<div class="field-row">
					<label>Role:</label>
					<div class="radio-group">
						<label><input type="radio" name="role" checked /> Admin</label>
						<label><input type="radio" name="role" /> Editor</label>
						<label><input type="radio" name="role" /> Viewer</label>
					</div>
				</div>

				<div class="field">
					<label for="bio">Bio</label>
					<textarea id="bio" placeholder="Tell us about yourself...">Hello, I'm a TUI user!</textarea>
				</div>

				<div class="field">
					<label for="lang">Language</label>
					<select id="lang">
						<option>English</option>
						<option selected>Português</option>
						<option>Español</option>
						<option>日本語</option>
					</select>
				</div>

				<div class="button-row">
					<button class="btn primary" type="submit">Login</button>
					<button class="btn secondary" type="reset">Reset</button>
					<button class="btn" type="button" disabled>Disabled</button>
				</div>
			</form>

			<p class="hint">Tab to navigate · Type to edit · Enter/Space to toggle · Esc to unfocus · q to quit</p>
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
		form {
			padding: 1;
			border: solid #0f3460;
			background-color: #16213e;
		}
		.field {
			margin-bottom: 1;
		}
		.field-row {
			margin-bottom: 1;
			display: flex;
			align-items: center;
			gap: 1;
		}
		label {
			color: #e94560;
			font-weight: bold;
			display: block;
		}
		input, textarea, select {
			padding: 0 1;
			background-color: #1a1a2e;
			border: solid #0f3460;
			color: #c0c0c0;
			width: 100%;
		}
		input[type="checkbox"], input[type="radio"] {
			width: auto;
		}
		.checkbox label {
			display: inline;
		}
		.radio-group {
			display: flex;
			gap: 2;
		}
		.button-row {
			display: flex;
			gap: 1;
			justify-content: center;
		}
		.btn {
			padding: 0 2;
			background-color: #0f3460;
			border: solid #00d4aa;
			color: #c0c0c0;
		}
		.btn.primary {
			background-color: #00d4aa;
			color: #1a1a2e;
			font-weight: bold;
		}
		.btn.secondary {
			border-color: #e94560;
			color: #e94560;
		}
		.btn[disabled] {
			border-color: #555;
			color: #555;
		}
		/* Focus styling */
		input[focused], textarea[focused], select[focused] {
			border: solid #00d4aa;
			background-color: #0f3460;
		}
		.btn[focused] {
			border: solid #ffffff;
			background-color: #00d4aa;
			color: #1a1a2e;
		}
		input[type="checkbox"][focused], input[type="radio"][focused] {
			outline: solid #00d4aa;
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
