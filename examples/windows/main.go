// Example: Windows and Panels
//
// This example demonstrates multi-window/panel layouts.
// While true window management is not yet implemented,
// this shows the intended design for split panes and panels.
package main

import (
	"log"

	"github.com/elioneto/tuix"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Windows & Panels Demo</h1>

			<div class="desktop">
				<!-- Title bar -->
				<div class="title-bar">
					<span class="title-text">TUI Window Manager</span>
					<div class="title-actions">
						<span class="action">—</span>
						<span class="action">□</span>
						<span class="action close">✕</span>
					</div>
				</div>

				<!-- Menu bar -->
				<div class="menu-bar">
					<span class="menu-item">File</span>
					<span class="menu-item">Edit</span>
					<span class="menu-item">View</span>
					<span class="menu-item">Help</span>
				</div>

				<!-- Main content area with side panels -->
				<div class="main-area">
					<!-- Left panel: File explorer -->
					<div class="panel left-panel">
						<div class="panel-title">Explorer</div>
						<div class="tree">
							<div class="tree-item folder open">📁 Project</div>
							<div class="tree-item folder indent open">📂 src</div>
							<div class="tree-item indent2">📄 main.go</div>
							<div class="tree-item indent2">📄 app.go</div>
							<div class="tree-item folder indent">📂 docs</div>
							<div class="tree-item indent2">📄 README.md</div>
							<div class="tree-item folder">📂 tests</div>
						</div>
					</div>

					<!-- Center: Editor -->
					<div class="panel editor-panel">
						<div class="panel-title">main.go — tuix</div>
						<div class="editor-content">
							<div class="line"><span class="line-num">1</span> <span class="keyword">package</span> <span class="string">main</span></div>
							<div class="line"><span class="line-num">2</span></div>
							<div class="line"><span class="line-num">3</span> <span class="keyword">import</span> (</div>
							<div class="line"><span class="line-num">4</span> 	<span class="string">"log"</span></div>
							<div class="line"><span class="line-num">5</span> </div>
							<div class="line"><span class="line-num">6</span> 	<span class="string">"github.com/elioneto/tuix"</span></div>
							<div class="line"><span class="line-num">7</span> )</div>
							<div class="line"><span class="line-num">8</span></div>
							<div class="line"><span class="line-num">9</span> <span class="keyword">func</span> <span class="func-name">main</span>() {</div>
							<div class="line highlight-line"><span class="line-num">10</span> 	app := tuix.New()</div>
							<div class="line"><span class="line-num">11</span> 	app.SetHTML(...)</div>
							<div class="line"><span class="line-num">12</span> 	app.SetCSS(...)</div>
							<div class="line"><span class="line-num">13</span> 	log.Fatal(app.Run())</div>
							<div class="line"><span class="line-num">14</span> }</div>
						</div>

						<!-- Status bar -->
						<div class="status-bar">
							<span>UTF-8</span>
							<span>Go</span>
							<span>Ln 10, Col 12</span>
						</div>
					</div>

					<!-- Right panel: Terminal -->
					<div class="panel right-panel">
						<div class="panel-title">Terminal</div>
						<div class="terminal-content">
							<div class="term-line"><span class="prompt">$</span> go build ./...</div>
							<div class="term-line success">✓ build succeeded</div>
							<div class="term-line"><span class="prompt">$</span> go run ./examples/counter</div>
							<div class="term-line">Starting counter app...</div>
							<div class="term-line cursor-line"><span class="prompt">$</span> █</div>
						</div>
					</div>
				</div>

				<!-- Bottom panel: Output -->
				<div class="panel bottom-panel">
					<div class="panel-title">Output</div>
					<div class="output-content">
						<div class="out-line info">[INFO]  Building project...</div>
						<div class="out-line success">[DONE]  Build complete (0.32s)</div>
						<div class="out-line warn">[WARN]  Unused variable 'x' in main.go:12</div>
						<div class="out-line error">[ERROR] Connection timeout (retrying...)</div>
					</div>
				</div>
			</div>

			<p class="hint">Press q to quit. Window management coming soon!</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 0;
			background-color: #1e1e2e;
			color: #c0c0c0;
		}
		h1 {
			color: #00d4aa;
			text-align: center;
			font-weight: bold;
			padding: 0 1;
			background-color: #1e1e2e;
		}

		/* Desktop container */
		.desktop {
			padding: 0;
			border: solid #45475a;
			background-color: #1e1e2e;
		}

		/* Title bar */
		.title-bar {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 0 1;
			background-color: #181825;
			border-bottom: solid #45475a;
		}
		.title-text { color: #00d4aa; font-weight: bold; }
		.title-actions { display: flex; gap: 1; }
		.action { color: #888; }
		.action.close { color: #e94560; }

		/* Menu bar */
		.menu-bar {
			display: flex;
			gap: 1;
			padding: 0 1;
			background-color: #181825;
			border-bottom: solid #45475a;
		}
		.menu-item { color: #c0c0c0; }

		/* Main layout */
		.main-area {
			display: flex;
			height: 14;
		}

		/* Panels */
		.panel {
			border-right: solid #45475a;
			display: flex;
			flex-direction: column;
		}
		.panel-title {
			background-color: #181825;
			color: #cdd6f4;
			font-weight: bold;
			padding: 0 1;
			border-bottom: solid #45475a;
		}

		/* Left panel */
		.left-panel {
			width: 20;
		}
		.tree { padding: 0; }
		.tree-item { padding: 0 1; color: #c0c0c0; }
		.tree-item.folder { color: #f9e2af; }
		.tree-item.indent { padding-left: 3; }
		.tree-item.indent2 { padding-left: 5; color: #a6adc8; }

		/* Editor */
		.editor-panel { flex: 1; }
		.editor-content { padding: 0; }
		.line { padding: 0; }
		.line-num { color: #585b70; width: 3; display: inline; }
		.highlight-line { background-color: #313244; }
		.keyword { color: #cba6f7; }
		.string { color: #a6e3a1; }
		.func-name { color: #89b4fa; }

		/* Right panel */
		.right-panel { width: 24; }
		.terminal-content { padding: 0; }
		.term-line { padding: 0 1; color: #c0c0c0; }
		.prompt { color: #a6e3a1; font-weight: bold; }
		.term-line.success { color: #a6e3a1; }
		.cursor-line { background-color: #313244; }

		/* Bottom panel */
		.bottom-panel {
			border-right: none;
			border-top: solid #45475a;
		}
		.output-content { padding: 0; }
		.out-line { padding: 0 1; }
		.out-line.info { color: #89b4fa; }
		.out-line.success { color: #a6e3a1; }
		.out-line.warn { color: #f9e2af; }
		.out-line.error { color: #e94560; }

		/* Status bar */
		.status-bar {
			display: flex;
			justify-content: space-between;
			padding: 0 2;
			background-color: #181825;
			border-top: solid #45475a;
			color: #585b70;
			font-size: 10;
		}

		.hint { text-align: center; color: #555; padding: 0 1; }
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
