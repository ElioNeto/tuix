// Example: Test Suite
//
// This example runs visual tests for the tuix library.
// Shows various layout patterns and edge cases for testing.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Visual Test Suite</h1>

			<!-- Test 1: Basic box model -->
			<div class="test">
				<h2>Test 1: Box Model</h2>
				<div class="box-model-test">
					<div class="box-a">Content</div>
				</div>
			</div>

			<!-- Test 2: Flex wrap -->
			<div class="test">
				<h2>Test 2: Flex Wrap</h2>
				<div class="flex-wrap-test">
					<div class="flex-child">One</div>
					<div class="flex-child">Two</div>
					<div class="flex-child">Three</div>
					<div class="flex-child">Four</div>
					<div class="flex-child">Five</div>
					<div class="flex-child">Six</div>
				</div>
			</div>

			<!-- Test 3: Nested flex -->
			<div class="test">
				<h2>Test 3: Nested Flex</h2>
				<div class="nested-flex">
					<div class="col">
						<div class="row-item">A1</div>
						<div class="row-item">A2</div>
					</div>
					<div class="col">
						<div class="row-item">B1</div>
						<div class="row-item">B2</div>
					</div>
				</div>
			</div>

			<!-- Test 4: Justify content -->
			<div class="test">
				<h2>Test 4: Justify Content</h2>
				<div class="justify-row" style="justify-content: flex-start">
					<div class="j-item">start</div>
					<div class="j-item">start</div>
					<div class="j-item">start</div>
				</div>
				<div class="justify-row" style="justify-content: center">
					<div class="j-item">center</div>
					<div class="j-item">center</div>
				</div>
				<div class="justify-row" style="justify-content: flex-end">
					<div class="j-item">end</div>
					<div class="j-item">end</div>
				</div>
				<div class="justify-row" style="justify-content: space-between">
					<div class="j-item">space</div>
					<div class="j-item">between</div>
				</div>
			</div>

			<!-- Test 5: Align items -->
			<div class="test">
				<h2>Test 5: Align Items</h2>
				<div class="align-row" style="align-items: flex-start; height: 5">
					<div class="a-item">top</div>
					<div class="a-item">top</div>
				</div>
				<div class="align-row" style="align-items: center; height: 5">
					<div class="a-item">center</div>
					<div class="a-item">center</div>
				</div>
				<div class="align-row" style="align-items: flex-end; height: 5">
					<div class="a-item">bottom</div>
					<div class="a-item">bottom</div>
				</div>
			</div>

			<!-- Test 6: Flex grow -->
			<div class="test">
				<h2>Test 6: Flex Grow</h2>
				<div class="grow-row">
					<div class="grow-item" style="flex-grow: 1">Grow 1</div>
					<div class="grow-item" style="flex-grow: 2">Grow 2</div>
					<div class="grow-item" style="flex-grow: 1">Grow 1</div>
				</div>
			</div>

			<!-- Test 7: Gap -->
			<div class="test">
				<h2>Test 7: Gap</h2>
				<div class="gap-row">
					<div class="gap-item">A</div>
					<div class="gap-item">B</div>
					<div class="gap-item">C</div>
					<div class="gap-item">D</div>
				</div>
			</div>

			<!-- Test 8: Borders -->
			<div class="test">
				<h2>Test 8: Borders</h2>
				<div class="border-row">
					<div class="b-border">Solid</div>
					<div class="b-dashed">
						<div>Double</div>
						<div>Line</div>
					</div>
					<div class="b-thick">Thick</div>
				</div>
			</div>

			<p class="hint">Press q to quit. All tests are visual (check output).</p>
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
		.test {
			margin-bottom: 1;
			padding: 1;
			border: solid #0f3460;
			background-color: #16213e;
		}

		/* Test 1: Box model */
		.box-model-test {
			padding: 2;
			border: solid #00d4aa;
			background-color: #1a1a2e;
		}
		.box-a {
			background-color: #0f3460;
			color: #c0c0c0;
			padding: 1;
		}

		/* Test 2: Flex wrap */
		.flex-wrap-test {
			display: flex;
			flex-wrap: wrap;
			gap: 1;
		}
		.flex-child {
			padding: 0 2;
			background-color: #0f3460;
			border: solid #00d4aa;
			color: #00d4aa;
		}

		/* Test 3: Nested flex */
		.nested-flex {
			display: flex;
			gap: 1;
		}
		.col {
			display: flex;
			flex-direction: column;
			gap: 1;
			flex: 1;
			padding: 1;
			border: solid #0f3460;
		}
		.row-item {
			padding: 0 2;
			background-color: #533483;
			color: #fff;
		}

		/* Test 4: Justify content */
		.justify-row {
			display: flex;
			gap: 1;
			margin-bottom: 1;
			padding: 0;
			background-color: #0f3460;
		}
		.j-item {
			padding: 0 1;
			background-color: #1a1a2e;
			color: #00d4aa;
			border: solid #00d4aa;
		}

		/* Test 5: Align items */
		.align-row {
			display: flex;
			gap: 1;
			margin-bottom: 1;
			padding: 0;
			background-color: #0f3460;
		}
		.a-item {
			padding: 0 2;
			background-color: #1a1a2e;
			color: #e94560;
			border: solid #e94560;
		}

		/* Test 6: Flex grow */
		.grow-row {
			display: flex;
			gap: 1;
		}
		.grow-item {
			padding: 0 2;
			background-color: #0f3460;
			border: solid #00d4aa;
			color: #c0c0c0;
			text-align: center;
		}

		/* Test 7: Gap */
		.gap-row {
			display: flex;
			gap: 1;
		}
		.gap-item {
			padding: 0 2;
			background-color: #e94560;
			color: #fff;
			font-weight: bold;
		}

		/* Test 8: Borders */
		.border-row {
			display: flex;
			gap: 1;
		}
		.b-border {
			padding: 1;
			border: solid #00d4aa;
			color: #00d4aa;
		}
		.b-dashed {
			padding: 1;
			border: double #e94560;
			color: #e94560;
		}
		.b-thick {
			padding: 1;
			border: solid #c0c0c0;
			color: #c0c0c0;
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
