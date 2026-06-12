package main

import (
	"time"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
<div id="app">
	<h1>Toast Notifications Demo</h1>
	<p>Press keys to show different types of toasts:</p>
	<div class="keybinds">
		<p><kbd>I</kbd> — Info toast</p>
		<p><kbd>S</kbd> — Success toast</p>
		<p><kbd>W</kbd> — Warning toast</p>
		<p><kbd>E</kbd> — Error toast</p>
		<p><kbd>A</kbd> — Alert dialog</p>
		<p><kbd>C</kbd> — Confirm dialog</p>
		<p><kbd>Q</kbd> — Quit</p>
	</div>
	<div id="log">
		<p>Waiting for input...</p>
	</div>
</div>
`)

	app.SetCSS(`
body {
	background: #1a1a2e;
	color: #e0e0e0;
	padding: 2;
}

h1 {
	color: #00d4aa;
	font-size: 20;
	margin-bottom: 1;
}

.keybinds p {
	margin: 0 0 1 0;
}

kbd {
	color: #ff6b6b;
	font-weight: bold;
}

#log {
	margin-top: 2;
	border: 1;
	border-color: #0f3460;
	padding: 1 2;
	background: #16213e;
	min-height: 5;
}

button {
	background: #16213e;
	color: #00d4aa;
	border: 1;
	border-color: #00d4aa;
	padding: 1 2;
	margin: 1 0;
}
button:focus {
	outline: 1;
	outline-style: solid;
	outline-color: #ff6b6b;
}

/* Dialog styles for Alert/Confirm modals */
.dialog {
	background: #1a1a2e;
	border: 2;
	border-color: #00d4aa;
	padding: 2 3;
	min-width: 40;
	max-width: 60;
}
.dialog h2 { color: #00d4aa; margin-bottom: 1; }
.dialog p { margin-bottom: 1; }
.dialog .actions {
	display: flex;
	gap: 1;
	justify-content: flex-end;
	margin-top: 2;
}
.dialog .actions button { min-width: 10; }
`)

	app.OnRune(func(r rune) {
		switch r {
		case 'i', 'I':
			app.ShowToast(tuix.ToastEntry{
				Message:  "This is an informational message",
				Type:     tuix.ToastInfo,
				Duration: 3 * time.Second,
			})
		case 's', 'S':
			app.ShowToast(tuix.ToastEntry{
				Message:  "Operation completed successfully!",
				Type:     tuix.ToastSuccess,
				Duration: 3 * time.Second,
			})
		case 'w', 'W':
			app.ShowToast(tuix.ToastEntry{
				Message:  "Warning: Low disk space",
				Type:     tuix.ToastWarning,
				Duration: 4 * time.Second,
			})
		case 'e', 'E':
			app.ShowToast(tuix.ToastEntry{
				Message:  "Error: Connection failed",
				Type:     tuix.ToastError,
				Duration: 5 * time.Second,
			})
		case 'a', 'A':
			app.Alert("Alert", "This is an alert dialog.\nPress OK or press Escape to dismiss.")
		case 'c', 'C':
			app.Confirm("Confirm", "Are you sure you want to proceed?", func(ok bool) {
				if ok {
			app.ShowToast(tuix.ToastEntry{
				Message:  "You confirmed!",
				Type:     tuix.ToastSuccess,
				Duration: 3 * time.Second,
			})
		} else {
			app.ShowToast(tuix.ToastEntry{
				Message:  "You cancelled.",
				Type:     tuix.ToastInfo,
				Duration: 3 * time.Second,
			})
		}
			})
		case 'q', 'Q':
			app.Stop()
		}
	})

	app.OnInit(func() {
		// Show a welcome toast
		app.ShowToast(tuix.ToastEntry{
			Message:  "Welcome! Press I, S, W, E for toasts",
			Type:     tuix.ToastInfo,
			Duration: 4 * time.Second,
		})
	})

	app.Run()
}
