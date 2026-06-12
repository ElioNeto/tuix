package main

import (
	"github.com/elioneto/tuix"
	"github.com/elioneto/tuix/terminal"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
<div id="app">
	<h1>Modal Demo</h1>
	<p>Press <strong>M</strong> to open the modal, or tab to the button and press Enter.</p>
	<button id="open-modal">Open Modal</button>
	<div id="content">
		<p>This is the main content behind the modal.</p>
		<p>When the modal is open, Tab focus is trapped inside it.</p>
		<p>Press <strong>Escape</strong> or click/tab to the X button and press Enter to close.</p>
		<div class="form-fields">
			<p>Name: <input type="text" value="John" id="name"/></p>
			<p>Email: <input type="text" value="john@example.com" id="email"/></p>
		</div>
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

#content {
	margin-top: 2;
	border: 1;
	border-color: #0f3460;
	padding: 1 2;
	background: #16213e;
}

.form-fields input {
	background: #0f3460;
	color: #e0e0e0;
	border: 1;
	border-color: #0f3460;
	padding: 0 1;
	min-width: 30;
}
.form-fields input:focus {
	outline: 1;
	outline-style: solid;
	outline-color: #00d4aa;
}

/* Modal styles — used when ShowModal is called */
.modal-wrapper {
	position: fixed;
	top: 0; left: 0; right: 0; bottom: 0;
	background: rgba(0, 0, 0, 128);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 100;
}

.dialog {
	background: #1a1a2e;
	border: 2;
	border-color: #00d4aa;
	padding: 2 3;
	min-width: 40;
	max-width: 60;
}

.dialog h2 {
	color: #00d4aa;
	margin-bottom: 1;
}

.dialog p {
	margin-bottom: 1;
}

.dialog .close-btn {
	float: right;
	background: none;
	border: 1;
	border-color: #e74c3c;
	color: #e74c3c;
	padding: 0 1;
	cursor: pointer;
}
.dialog .close-btn:focus {
	outline: 1;
	outline-color: #ff6b6b;
}

.dialog input {
	background: #0f3460;
	color: #e0e0e0;
	border: 1;
	border-color: #0f3460;
	padding: 0 1;
	width: 100%;
}
.dialog input:focus {
	outline: 1;
	outline-style: solid;
	outline-color: #00d4aa;
}

.dialog .actions {
	display: flex;
	gap: 1;
	justify-content: flex-end;
	margin-top: 2;
}
.dialog .actions button {
	min-width: 10;
}
`)

	modalHTML := `
<div class="dialog">
	<div style="display: flex; justify-content: space-between; align-items: center;">
		<h2>Modal Dialog</h2>
		<button class="close-btn" id="modal-close">X</button>
	</div>
	<p>This is a modal dialog with a focus trap.</p>
	<p>
		<label>Name:</label><br/>
		<input type="text" value="" placeholder="Enter your name" id="modal-name"/>
	</p>
	<p>
		<label>Email:</label><br/>
		<input type="text" value="" placeholder="Enter your email" id="modal-email"/>
	</p>
	<div class="actions">
		<button id="modal-cancel">Cancel</button>
		<button id="modal-ok">OK</button>
	</div>
</div>
`

	// Handle Enter/Space on buttons via events
	app.OnEvent(func(event terminal.Event) {
		if event.Type == terminal.EventKey {
			if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
				focused := app.FocusedElement()
				if focused != nil {
					id := focused.GetAttribute("id")
					switch id {
					case "open-modal":
						app.ShowModal(modalHTML)
					case "modal-close", "modal-cancel":
						app.CloseModal()
					case "modal-ok":
						app.CloseModal()
					}
				}
			}
		}
	})

	// Also open modal when 'm' or 'M' is pressed
	app.OnRune(func(r rune) {
		if r == 'm' || r == 'M' {
			app.ShowModal(modalHTML)
		}
	})

	app.Run()
}
