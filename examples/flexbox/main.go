// Example: Flexbox layout demo
//
// This example demonstrates the flexbox layout engine with various
// flex properties: flex-direction, justify-content, align-items,
// flex-grow, order, and gap.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Flexbox Demo</h1>

			<!-- Row layout with space-between -->
			<div class="section">
				<h2>justify-content: space-between</h2>
				<div class="flex-row space-between">
					<div class="box red">Item A</div>
					<div class="box green">Item B</div>
					<div class="box blue">Item C</div>
				</div>
			</div>

			<!-- Row layout with center -->
			<div class="section">
				<h2>justify-content: center</h2>
				<div class="flex-row center">
					<div class="box red">Center</div>
					<div class="box green">Aligned</div>
				</div>
			</div>

			<!-- Flex-grow: proportional sizing -->
			<div class="section">
				<h2>flex-grow: proportional</h2>
				<div class="flex-row grow">
					<div class="box grow1 red">1x</div>
					<div class="box grow2 green">2x</div>
					<div class="box grow1 blue">1x</div>
				</div>
			</div>

			<!-- Column layout -->
			<div class="section">
				<h2>flex-direction: column</h2>
				<div class="flex-column stretch">
					<div class="box purple">Top</div>
					<div class="box orange">Middle</div>
					<div class="box teal">Bottom</div>
				</div>
			</div>

			<!-- Gap and order -->
			<div class="section">
				<h2>gap + order</h2>
				<div class="flex-row gap-order">
					<div class="box red" style="order: 3">Third</div>
					<div class="box green" style="order: 1">First</div>
					<div class="box blue" style="order: 2">Second</div>
				</div>
			</div>

			<p class="hint">Press q to quit</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 1;
			background-color: #1a1a2e;
			color: #e0e0e0;
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

		/* Flex row container */
		.flex-row {
			display: flex;
			flex-direction: row;
		}
		.space-between {
			justify-content: space-between;
		}
		.center {
			justify-content: center;
		}
		.grow {
			justify-content: flex-start;
		}

		/* Flex column container */
		.flex-column {
			display: flex;
			flex-direction: column;
		}
		.stretch {
			align-items: stretch;
		}

		/* Gap + order */
		.gap-order {
			column-gap: 2;
			justify-content: center;
		}

		/* Flex items */
		.box {
			padding: 1;
			font-weight: bold;
			text-align: center;
		}
		.grow1 { flex-grow: 1; }
		.grow2 { flex-grow: 2; }

		/* Colors */
		.red    { background-color: #e94560; color: #fff; }
		.green  { background-color: #0f3460; color: #00d4aa; }
		.blue   { background-color: #16213e; color: #0f3460; border: solid #00d4aa; }
		.purple { background-color: #533483; color: #fff; }
		.orange { background-color: #e07c24; color: #fff; }
		.teal   { background-color: #008080; color: #fff; }

		.hint {
			text-align: center;
			color: #555;
			margin-top: 1;
		}
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
