// Example: Z-Index / Stacking Contexts
//
// This example demonstrates z-index layering with overlapping
// positioned elements.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Z-Index Demo</h1>

			<div class="stack">
				<div class="box box1" style="z-index: 1">z-index: 1</div>
				<div class="box box2" style="z-index: 2">z-index: 2 (on top)</div>
				<div class="box box3" style="z-index: 0">z-index: 0 (behind)</div>
			</div>

			<p class="hint">Elements with higher z-index paint on top. Press q to quit.</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 2;
			background-color: #1a1a2e;
			color: #c0c0c0;
		}
		h1 {
			color: #00d4aa;
			text-align: center;
			margin-bottom: 2;
		}
		.stack {
			position: relative;
			height: 12;
		}
		.box {
			position: absolute;
			padding: 2;
			width: 20;
			border: solid #0f3460;
			text-align: center;
		}
		.box1 {
			background-color: #e94560;
			color: #fff;
			top: 0;
			left: 0;
			z-index: 1;
		}
		.box2 {
			background-color: #00d4aa;
			color: #1a1a2e;
			top: 3;
			left: 8;
			z-index: 2;
		}
		.box3 {
			background-color: #0f3460;
			color: #c0c0c0;
			top: 6;
			left: 16;
			z-index: 0;
		}
		.hint { text-align: center; color: #555; margin-top: 2; }
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
