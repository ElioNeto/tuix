// Example: Transitions and Animations
//
// This example demonstrates animated effects.
// CSS transition/animation support is coming soon!
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Animations Demo</h1>

			<div class="section">
				<h2>Loading Spinner</h2>
				<div class="spinner-container">
					<div class="spinner">
						<span class="spinner-dot">●</span>
						<span class="spinner-dot">●</span>
						<span class="spinner-dot">●</span>
						<span class="spinner-dot">●</span>
					</div>
					<span class="loading-text">Loading...</span>
				</div>
			</div>

			<div class="section">
				<h2>Progress Bar</h2>
				<div class="progress-bar">
					<div class="progress-fill" style="width: 65%"></div>
				</div>
				<div class="progress-label">65% complete</div>
			</div>

			<div class="section">
				<h2>Notification / Toast</h2>
				<div class="toast">
					<span class="toast-icon">✓</span>
					<span class="toast-msg">File saved successfully!</span>
				</div>
				<div class="toast warning">
					<span class="toast-icon">⚠</span>
					<span class="toast-msg">Low disk space</span>
				</div>
				<div class="toast error">
					<span class="toast-icon">✗</span>
					<span class="toast-msg">Connection lost</span>
				</div>
			</div>

			<div class="section">
				<h2>Pulse Effect</h2>
				<div class="pulse-container">
					<div class="pulse-dot pulse"></div>
					<span>Live</span>
				</div>
			</div>

			<div class="section">
				<h2>Slide-in Panel</h2>
				<div class="slide-panel">
					<div class="panel-header">Notifications</div>
					<div class="panel-body">
						<div class="notification slide-in">📧 New email from Alice</div>
						<div class="notification slide-in">🔔 Merge request approved</div>
						<div class="notification slide-in">📢 System update available</div>
					</div>
				</div>
			</div>

			<p class="hint">Press q to quit. Animations coming soon!</p>
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
			background-color: #16213e;
		}

		/* Spinner */
		.spinner-container {
			display: flex;
			align-items: center;
			gap: 1;
		}
		.spinner {
			display: flex;
			gap: 0;
		}
		.spinner-dot {
			color: #00d4aa;
			animation: pulse 1s infinite;
		}
		.loading-text {
			color: #888;
			font-style: italic;
		}

		/* Progress bar */
		.progress-bar {
			background-color: #0f3460;
			height: 1;
		}
		.progress-fill {
			background-color: #00d4aa;
			height: 1;
		}
		.progress-label {
			text-align: center;
			color: #00d4aa;
			font-size: 10;
		}

		/* Toast notifications */
		.toast {
			display: flex;
			align-items: center;
			gap: 1;
			padding: 0 2;
			margin-bottom: 1;
			background-color: #0f3460;
			border-left: solid #00d4aa;
			transition: all 0.3s ease;
		}
		.toast.warning {
			border-left-color: #e07c24;
		}
		.toast.error {
			border-left-color: #e94560;
		}
		.toast-icon {
			font-weight: bold;
		}
		.toast .toast-icon { color: #00d4aa; }
		.toast.warning .toast-icon { color: #e07c24; }
		.toast.error .toast-icon { color: #e94560; }

		/* Pulse */
		.pulse-container {
			display: flex;
			align-items: center;
			gap: 1;
		}
		.pulse-dot {
			width: 1;
			height: 1;
			background-color: #00d4aa;
		}
		.pulse-dot.pulse {
			background-color: #00ff88;
		}

		/* Slide panel */
		.slide-panel {
			border: solid #0f3460;
		}
		.panel-header {
			background-color: #0f3460;
			color: #00d4aa;
			font-weight: bold;
			padding: 0 1;
		}
		.panel-body {
			padding: 1;
		}
		.notification {
			padding: 0 1;
			color: #c0c0c0;
		}
		.notification.slide-in {
			border-left: solid #00d4aa;
			background-color: #1a1a2e;
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
