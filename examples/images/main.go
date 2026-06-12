// Example: Image Rendering
//
// This example demonstrates image rendering using ASCII/Unicode art.
// Actual image loading and rendering is coming soon!
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Image Rendering Demo</h1>

			<div class="section">
				<h2>Logo (ASCII Art)</h2>
				<div class="image logo" alt="TUIX Logo">
					<pre>
    ████████╗██╗   ██╗██╗██╗  ██╗
    ╚══██╔══╝██║   ██║██║╚██╗██╔╝
       ██║   ██║   ██║██║ ╚███╔╝
       ██║   ██║   ██║██║ ██╔██╗
       ██║   ╚██████╔╝██║██╔╝ ██╗
       ╚═╝    ╚═════╝ ╚═╝╚═╝  ╚═╝
					</pre>
				</div>
			</div>

			<div class="section">
				<h2>Avatar Placeholder</h2>
				<div class="avatar-row">
					<div class="avatar" style="background-color: #e94560">
						<span class="avatar-initials">JD</span>
					</div>
					<div class="avatar" style="background-color: #00d4aa; color: #1a1a2e">
						<span class="avatar-initials">AL</span>
					</div>
					<div class="avatar" style="background-color: #533483">
						<span class="avatar-initials">MK</span>
					</div>
					<div class="avatar" style="background-color: #e07c24">
						<span class="avatar-initials">RS</span>
					</div>
					<div class="avatar" style="background-color: #0f3460; border: solid #00d4aa">
						<span class="avatar-initials">+3</span>
					</div>
				</div>
			</div>

			<div class="section">
				<h2>Photo Gallery Grid</h2>
				<div class="gallery-grid">
					<div class="gallery-item">
						<div class="gallery-img" style="background-color: #2d1b69">
							<span class="gallery-icon">🏔️</span>
						</div>
						<div class="gallery-label">Mountains</div>
					</div>
					<div class="gallery-item">
						<div class="gallery-img" style="background-color: #1b3a5c">
							<span class="gallery-icon">🌊</span>
						</div>
						<div class="gallery-label">Ocean</div>
					</div>
					<div class="gallery-item">
						<div class="gallery-img" style="background-color: #3d1c1c">
							<span class="gallery-icon">🌅</span>
						</div>
						<div class="gallery-label">Sunset</div>
					</div>
					<div class="gallery-item">
						<div class="gallery-img" style="background-color: #1c3d1c">
							<span class="gallery-icon">🌲</span>
						</div>
						<div class="gallery-label">Forest</div>
					</div>
				</div>
			</div>

			<div class="section">
				<h2>Icon Library</h2>
				<div class="icon-row">
					<span class="icon" title="Save">💾</span>
					<span class="icon" title="Edit">✏️</span>
					<span class="icon" title="Delete">🗑️</span>
					<span class="icon" title="Search">🔍</span>
					<span class="icon" title="Settings">⚙️</span>
					<span class="icon" title="User">👤</span>
					<span class="icon" title="Bell">🔔</span>
					<span class="icon" title="Star">⭐</span>
				</div>
			</div>

			<p class="hint">Press q to quit. Image loading coming soon!</p>
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

		/* Logo */
		.image.logo {
			display: flex;
			justify-content: center;
			padding: 1;
		}
		.image.logo pre {
			color: #00d4aa;
			font-family: monospace;
		}

		/* Avatars */
		.avatar-row {
			display: flex;
			gap: 2;
			justify-content: center;
		}
		.avatar {
			width: 5;
			height: 3;
			display: flex;
			align-items: center;
			justify-content: center;
			border-radius: 50%;
		}
		.avatar-initials {
			font-weight: bold;
			color: #fff;
			font-size: 12;
		}

		/* Gallery */
		.gallery-grid {
			display: flex;
			flex-wrap: wrap;
			gap: 1;
		}
		.gallery-item {
			width: 12;
			text-align: center;
		}
		.gallery-img {
			height: 5;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 20;
			border: solid #0f3460;
		}
		.gallery-label {
			color: #888;
			font-size: 10;
		}

		/* Icons */
		.icon-row {
			display: flex;
			gap: 2;
			justify-content: center;
		}
		.icon {
			font-size: 16;
			cursor: pointer;
		}
		.icon:hover {
			background-color: #0f3460;
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
