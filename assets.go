package main

import "html/template"

const cssContent = `
:root {
	color-scheme: dark;
	--font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
	--font-mono: "JetBrains Mono", "Fira Code", ui-monospace, "SF Mono", "Cascadia Code", Menlo, Consolas, monospace;

	/* Catppuccin Mocha */
	--base: #1e1e2e;
	--mantle: #181825;
	--crust: #11111b;
	--surface0: #313244;
	--surface1: #45475a;
	--surface2: #585b70;
	--overlay0: #6c7086;
	--overlay1: #7f849c;
	--text: #cdd6f4;
	--subtext0: #a6adc8;
	--subtext1: #bac2de;
	--pink: #f5c2e7;
	--mauve: #cba6f7;
	--blue: #89b4fa;
	--green: #a6e3a1;
	--red: #f38ba8;
	--yellow: #f9e2af;
	--peach: #fab387;
	--teal: #94e2d5;
	--rosewater: #f5e0dc;

	--radius-xs: 4px;
	--radius-sm: 8px;
	--radius-md: 12px;
	--radius-lg: 18px;
	--radius-full: 999px;

	--shadow-sm: 0 1px 3px rgba(17, 17, 27, 0.35);
	--shadow-md: 0 8px 24px rgba(17, 17, 27, 0.45);
	--shadow-toast: 0 10px 30px rgba(17,17,27, 0.6);
	--focus-ring: 0 0 0 3px rgba(203, 166, 247, 0.25);
	--topbar-bg: rgba(24, 24, 37, 0.92);
}

@media (prefers-color-scheme: light) {
	:root {
		color-scheme: light;
		/* Catppuccin Latte */
		--base: #eff1f5;
		--mantle: #e6e9ef;
		--crust: #dce0e8;
		--surface0: #ccd0da;
		--surface1: #bcc0cc;
		--surface2: #acb0be;
		--overlay0: #9ca0b0;
		--overlay1: #8c8fa1;
		--text: #4c4f69;
		--subtext0: #6c6f85;
		--subtext1: #5c5f77;
		--pink: #ea76cb;
		--mauve: #8839ef;
		--blue: #1e66f5;
		--green: #40a02b;
		--red: #d20f39;
		--yellow: #df8e1d;
		--peach: #fe640b;
		--teal: #179299;
		--rosewater: #dc8a78;

		--shadow-sm: 0 1px 3px rgba(76, 79, 105, 0.12);
		--shadow-md: 0 8px 24px rgba(76, 79, 105, 0.16);
		--shadow-toast: 0 10px 30px rgba(76, 79, 105, 0.2);
		--focus-ring: 0 0 0 3px rgba(136, 57, 239, 0.18);
		--topbar-bg: rgba(230, 233, 239, 0.92);
	}
}

* { box-sizing: border-box; }
html { -webkit-text-size-adjust: 100%; }
body {
	margin: 0;
	background: var(--base);
	color: var(--text);
	font-family: var(--font-sans);
	font-size: 15px;
	line-height: 1.55;
	min-height: 100vh;
}
.container { max-width: 1000px; margin: 0 auto; padding: 0 20px 40px; }
a { color: var(--blue); }
::selection { background: var(--mauve); color: var(--crust); }

:focus-visible {
	outline: 2px solid var(--mauve);
	outline-offset: 2px;
	border-radius: var(--radius-xs);
}

/* top bar */
.topbar {
	position: sticky; top: 0; z-index: 10;
	background: var(--topbar-bg);
	backdrop-filter: blur(10px);
	-webkit-backdrop-filter: blur(10px);
	border-bottom: 1px solid var(--surface0);
	margin-bottom: 24px;
}
.topbar-inner {
	max-width: 1000px; margin: 0 auto; padding: 12px 20px;
	display: flex; align-items: center; gap: 14px; flex-wrap: wrap;
}
.brand {
	color: var(--mauve); text-decoration: none; font-weight: 700;
	letter-spacing: .02em; display: flex; align-items: center; gap: 8px; flex: none;
}
.brand svg { width: 18px; height: 18px; stroke: currentColor; fill: none; stroke-width: 2.5; }
.brand:hover { color: var(--pink); }

/* breadcrumbs */
.breadcrumbs {
	display: flex; align-items: center; flex-wrap: wrap; gap: 4px;
	min-width: 0; font-size: 0.9rem; font-family: var(--font-mono);
}
.breadcrumbs a {
	color: var(--subtext1); text-decoration: none; padding: 4px 8px; border-radius: var(--radius-sm);
	max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.breadcrumbs a:hover { background: var(--surface0); color: var(--text); }
.breadcrumbs a.current, .breadcrumbs span.current { color: var(--text); font-weight: 600; padding: 4px 8px; }
.crumb-sep { color: var(--overlay0); user-select: none; }
.cursor {
	display: inline-block; width: 7px; height: 1.05em; margin-left: 4px;
	background: var(--pink); vertical-align: -2px;
	animation: blink 1.1s steps(1) infinite; flex-shrink: 0;
}
@keyframes blink { 50% { opacity: 0; } }

/* controls / toolbar */
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 24px; align-items: stretch; }

.upload-zone {
	display: flex; align-items: center; gap: 12px; flex-wrap: wrap; flex: 2; min-width: 260px;
	background: var(--surface0); border: 1.5px dashed var(--surface2); border-radius: var(--radius-md);
	padding: 10px 16px; cursor: pointer; color: var(--subtext0); font-size: 0.88rem;
	transition: border-color 0.15s ease, background 0.15s ease;
	position: relative; overflow: hidden;
}
.upload-zone:hover { border-color: var(--overlay1); }
.upload-zone.dragover { border-color: var(--mauve); border-style: solid; background: var(--mantle); color: var(--text); }
.upload-zone svg { color: var(--mauve); flex-shrink: 0; width: 22px; height: 22px; stroke: currentColor; fill: none; stroke-width: 2; }
.upload-form { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; background: var(--surface0); border-radius: var(--radius-md); padding: 10px 16px; flex: 2; }

.progress { display: none; position: absolute; bottom: 0; left: 0; width: 100%; height: 3px; background: transparent; }
.progress-bar { height: 100%; width: 0; background: var(--mauve); transition: width 0.15s; }

/* menus & inputs */
details.menu { position: relative; }
details.menu > summary {
	list-style: none; display: inline-flex; align-items: center; gap: 6px;
	background: var(--surface1); color: var(--text);
	border-radius: var(--radius-md); padding: 10px 16px; font-size: 0.88rem; font-weight: 600;
	cursor: pointer; user-select: none; height: 100%; transition: background 0.12s;
}
details.menu > summary::-webkit-details-marker { display: none; }
details.menu > summary svg { width: 16px; height: 16px; stroke: currentColor; fill: none; stroke-width: 2.2; }
details.menu > summary:hover, details.menu[open] > summary { background: var(--surface2); }
.menu-panel {
	position: absolute; right: 0; top: calc(100% + 6px); z-index: 20;
	background: var(--surface0); border: 1px solid var(--surface1); border-radius: var(--radius-md);
	padding: 14px; min-width: 250px; display: flex; flex-direction: column; gap: 10px;
	box-shadow: var(--shadow-md);
}
.menu-panel form { display: flex; gap: 8px; }

input[type="text"], input[type="search"] {
	background: var(--mantle); color: var(--text); border: 1px solid var(--surface1);
	border-radius: var(--radius-sm); padding: 9px 12px; font-size: 0.9rem; font-family: inherit;
}
input[type="text"]:focus, input[type="search"]:focus { outline: none; border-color: var(--mauve); box-shadow: var(--focus-ring); }
input::placeholder { color: var(--overlay0); }
.search-wrap { position: relative; flex: 1; min-width: 160px; max-width: 280px; }
.search-wrap svg { position: absolute; left: 12px; top: 11px; width: 15px; height: 15px; color: var(--overlay0); stroke: currentColor; fill: none; stroke-width: 2; pointer-events: none; }
.search-box { width: 100%; padding-left: 36px !important; border-radius: var(--radius-md) !important; height: 100%; }

/* buttons */
button, .btn {
	background: var(--pink); color: var(--crust); border: none; padding: 9px 16px;
	border-radius: var(--radius-md); font-weight: 600; font-size: 0.85rem;
	cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; gap: 6px;
	transition: opacity 0.12s ease, transform 0.06s ease; white-space: nowrap;
}
button:hover, .btn:hover { opacity: 0.88; }
button:active, .btn:active { transform: scale(0.97); }
.btn-secondary { background: var(--surface1); color: var(--text); }
.btn-danger { background: var(--red); color: var(--crust); }

/* table */
.panel { background: var(--surface0); border-radius: var(--radius-lg); overflow: hidden; box-shadow: var(--shadow-sm); }
.file-table { width: 100%; border-collapse: collapse; }
.file-table th, .file-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid var(--surface1); vertical-align: middle; }
.file-table th {
	background: var(--mantle); color: var(--subtext0); font-weight: 600; font-size: 0.72rem;
	letter-spacing: 0.06em; text-transform: uppercase; cursor: pointer; user-select: none; white-space: nowrap;
}
.file-table th:hover { color: var(--text); }
.file-table th.is-sorted { color: var(--pink); }
.file-table th .sort-icon { display: inline-flex; width: 10px; height: 10px; margin-left: 4px; vertical-align: -1px; opacity: 0.55; }
.file-table th .sort-icon svg { width: 100%; height: 100%; stroke: currentColor; fill: none; stroke-width: 3; }
.file-table th[data-dir="desc"] .sort-icon { transform: rotate(180deg); }
.file-table th.is-sorted .sort-icon { opacity: 1; }
.file-table tbody tr:last-child td { border-bottom: none; }
.file-table tbody tr { transition: background 0.1s; }
.file-table tbody tr:hover { background: var(--surface1); }

.file-link { color: var(--text); text-decoration: none; display: flex; align-items: center; gap: 12px; font-weight: 500; font-family: var(--font-mono); font-size: 0.92rem; word-break: break-all; }
.file-link:hover { color: var(--pink); }
.file-link img.file-thumb { width: 28px; height: 28px; object-fit: cover; border-radius: var(--radius-xs); flex-shrink: 0; border: 1px solid var(--surface1); }
.file-icon { flex-shrink: 0; width: 20px; height: 20px; }
.file-icon svg { width: 100%; height: 100%; stroke-width: 1.7; fill: none; stroke: currentColor; stroke-linecap: round; stroke-linejoin: round; }
.dir-icon { color: var(--blue); }
.image-icon { color: var(--green); }
.video-icon { color: var(--peach); }
.generic-icon { color: var(--teal); }

.col-size, .col-date { font-family: var(--font-mono); font-size: 0.82rem; color: var(--subtext0); white-space: nowrap; font-variant-numeric: tabular-nums; }
.col-size { text-align: right; }

.actions { display: flex; gap: 4px; align-items: center; opacity: 0; transition: opacity 0.12s ease; }
.file-table tr:hover .actions, .file-table tr:focus-within .actions, .actions.always { opacity: 1; }
@media (hover: none) { .actions { opacity: 1; } }

.icon-btn {
	display: inline-flex; align-items: center; justify-content: center; width: 30px; height: 30px;
	border-radius: var(--radius-sm); background: transparent; color: var(--subtext0); border: none; cursor: pointer; text-decoration: none;
}
.icon-btn:hover { background: var(--surface1); color: var(--text); }
.icon-btn.danger:hover { background: rgba(243, 139, 168, 0.15); color: var(--red); }
.icon-btn svg { width: 16px; height: 16px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

/* empty state */
.empty-state { text-align: center; padding: 56px 20px; color: var(--subtext0); }
.empty-state svg { width: 44px; height: 44px; color: var(--overlay0); margin-bottom: 12px; stroke-width: 1.4; fill: none; stroke: currentColor; stroke-linecap: round; stroke-linejoin: round; }
.empty-state .empty-title { font-weight: 600; color: var(--subtext1); font-size: 1.1rem; margin-bottom: 6px; }

/* footer */
.footer { margin-top: 20px; padding: 0 4px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px; color: var(--subtext0); font-size: 0.8rem; font-family: var(--font-mono); }
.footer .version { color: var(--overlay0); }

/* toast */
.toast {
	position: fixed; bottom: 24px; left: 50%; z-index: 50; transform: translateX(-50%) translateY(10px);
	background: var(--surface0); color: var(--text); border: 1px solid var(--surface1);
	border-radius: var(--radius-md); padding: 10px 18px; font-size: 0.88rem; font-weight: 500;
	opacity: 0; pointer-events: none; transition: opacity 0.2s, transform 0.2s; box-shadow: var(--shadow-toast);
}
.toast.show { opacity: 1; transform: translateX(-50%); }

/* dialogs (confirm / rename) */
.dialog-box {
	max-width: 460px; margin: 40px auto 0; background: var(--surface0); border: 1px solid var(--surface1);
	border-radius: var(--radius-lg); padding: 32px; text-align: center; box-shadow: var(--shadow-md);
}
.dialog-box .dialog-icon {
	width: 56px; height: 56px; margin: 0 auto 16px; border-radius: 50%;
	display: flex; align-items: center; justify-content: center;
}
.dialog-box .dialog-icon svg { width: 28px; height: 28px; stroke: currentColor; fill: none; stroke-width: 2; }
.dialog-icon.warn { background: rgba(243, 139, 168, 0.15); color: var(--red); }
.dialog-icon.info { background: rgba(203, 166, 247, 0.15); color: var(--mauve); }
.dialog-box h2 { margin: 0 0 8px; font-size: 1.2rem; }
.dialog-box h2.danger { color: var(--red); }
.dialog-box p { color: var(--subtext0); margin: 0 0 24px; font-size: 0.95rem; word-break: break-word; }
.dialog-box .filename { color: var(--text); font-family: var(--font-mono); font-size: 0.9rem; }
.dialog-box .btn-row { justify-content: center; display: flex; gap: 12px; flex-wrap: wrap; }
.dialog-box input[type="text"] { width: 100%; margin-bottom: 20px; font-size: 1rem; }

/* editor */
.edit-head { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; flex-wrap: wrap; }
.edit-head h2 { margin: 0; font-size: 1.1rem; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: var(--font-mono); }
.edit-head .btn-row { display: flex; gap: 8px; flex-shrink: 0; }
.editor-shell { position: relative; background: var(--crust); border: 1px solid var(--surface1); border-radius: var(--radius-md); overflow: hidden; }
.editor-shell pre { margin: 0; position: absolute; inset: 0; overflow: hidden; padding: 16px; pointer-events: none; }
.editor-shell pre code { font-family: var(--font-mono); font-size: 0.875rem; line-height: 1.6; tab-size: 4; white-space: pre; background: transparent; padding: 0; display: block; }
textarea.code {
	position: relative; z-index: 1; display: block; width: 100%; min-height: 65vh; padding: 16px; border: none; outline: none; resize: vertical;
	background: transparent; color: var(--text); caret-color: var(--rosewater); font-family: var(--font-mono); font-size: 0.875rem; line-height: 1.6; tab-size: 4; white-space: pre; overflow: auto;
}
textarea.code.overlaid { color: transparent; }
textarea.code::selection { background: rgba(88, 91, 112, 0.6); color: transparent; }

/* mobile */
@media (max-width: 640px) {
	.container { padding: 0 12px 32px; }
	.topbar-inner { padding: 12px 14px; }
	.toolbar { flex-direction: column; }
	.search-wrap { max-width: none; }
	details.menu > summary, .menu-panel { width: 100%; }
	.menu-panel { position: static; margin-top: 8px; box-shadow: none; }

	.file-table thead { display: none; }
	.file-table, .file-table tbody { display: block; width: 100%; }
	.file-table tr {
		display: flex; flex-wrap: wrap; background: var(--mantle); border-radius: var(--radius-md);
		margin-bottom: 12px; padding: 12px 16px; border-bottom: none; box-shadow: var(--shadow-sm);
	}
	.file-table tr:last-child { margin-bottom: 0; }
	.file-table td { border-bottom: none; padding: 0; flex: 0 0 auto; }
	.file-table td:first-child { flex: 1 0 100%; padding-bottom: 8px; }
	.file-table td.col-date::before { content: "\2022"; margin: 0 6px; color: var(--overlay0); }
	.actions { opacity: 1; width: 100%; margin-top: 6px; padding-top: 8px; border-top: 1px solid var(--surface0); }
}
`

const iconFolder = `<svg viewBox="0 0 24 24"><path d="M3 7a2 2 0 0 1 2-2h4.5l2 2H19a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z"/></svg>`
const iconImage = `<svg viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="9" cy="9" r="1.5"/><path d="M21 15l-5-5-9 9"/></svg>`
const iconVideo = `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="9"/><polygon points="10 8 16 12 10 16"/></svg>`
const iconFile = `<svg viewBox="0 0 24 24"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M13 2v6h6"/></svg>`
const iconUpload = `<svg viewBox="0 0 24 24"><path d="M16 16l-4-4-4 4"/><path d="M12 12v9"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>`
const iconSearch = `<svg viewBox="0 0 24 24"><circle cx="11" cy="11" r="7"/><path d="M21 21l-4.35-4.35"/></svg>`
const iconPlus = `<svg viewBox="0 0 24 24"><path d="M12 5v14"/><path d="M5 12h14"/></svg>`
const iconLink = `<svg viewBox="0 0 24 24"><path d="M10 13a5 5 0 0 0 7 0l3-3a5 5 0 0 0-7-7l-1.5 1.5"/><path d="M14 11a5 5 0 0 0-7 0l-3 3a5 5 0 0 0 7 7l1.5-1.5"/></svg>`
const iconPencil = `<svg viewBox="0 0 24 24"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.1 2.1 0 0 1 3 3L12 15l-4 1 1-4z"/></svg>`
const iconRename = `<svg viewBox="0 0 24 24"><path d="M20.59 13.41L11.17 4A2 2 0 0 0 9.75 3.4L4.6 3.4a1.2 1.2 0 0 0-1.2 1.2v5.15c0 .53.21 1.04.59 1.42l9.42 9.42a2 2 0 0 0 2.83 0l4.35-4.35a2 2 0 0 0 0-2.83z"/><circle cx="7.9" cy="7.9" r="1.15"/></svg>`
const iconDownload = `<svg viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><path d="M7 10l5 5 5-5"/><path d="M12 15V3"/></svg>`
const iconTrash = `<svg viewBox="0 0 24 24"><path d="M3 6h18"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/></svg>`
const iconChevron = `<svg viewBox="0 0 24 24"><path d="M6 15l6-6 6 6"/></svg>`
const iconWarning = `<svg viewBox="0 0 24 24"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>`
const iconEmpty = `<svg viewBox="0 0 24 24"><path d="M21 8l-9-5-9 5 9 5 9-5z"/><path d="M3 8v8l9 5 9-5V8"/><path d="M12 13v8"/></svg>`
const iconBrand = `<svg viewBox="0 0 24 24"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>`

const breadcrumbBlock = `<nav class="breadcrumbs">
	{{ range $i, $bc := .Breadcrumbs }}
		{{ if $bc.IsLast }}<span class="current">{{ $bc.Name }}</span><span class="cursor"></span>{{ else }}<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a><span class="crumb-sep">/</span>{{ end }}
	{{ end }}
</nav>`

const editorBreadcrumbBlock = `<nav class="breadcrumbs">
	{{ range $i, $bc := .Breadcrumbs }}
		<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a><span class="crumb-sep">/</span>
	{{ end }}
	<span class="current">{{ .FileName }}</span><span class="cursor"></span>
</nav>`

const newMenuBlock = `<details class="menu">
	<summary>` + iconPlus + ` New</summary>
	<div class="menu-panel">
		<form action="/mkdir{{ .CurrentPath }}" method="POST">
			<input type="text" name="dirname" placeholder="Folder name" required>
			<button type="submit" class="btn-small btn-secondary">` + iconFolder + `</button>
		</form>
		<form action="/mkfile{{ .CurrentPath }}" method="POST">
			<input type="text" name="filename" placeholder="File name" required>
			<button type="submit" class="btn-small btn-secondary">` + iconFile + `</button>
		</form>
	</div>
</details>`

const emptyStateBlock = `<tr><td colspan="3">
	<div class="empty-state">
		` + iconEmpty + `
		<div class="empty-title">Nothing here yet</div>
		<div>Drop files above or create a new folder to get started.</div>
	</div>
</td></tr>`

const htmlFull = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>Index of {{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">` + iconBrand + ` jotasrv</a>
			` + breadcrumbBlock + `
		</div>
	</header>

	<div class="container">
		<div class="toolbar">
			<label class="upload-zone" id="dropZone" for="fileInput">
				` + iconUpload + `
				<span id="uploadLabel">Drop files or click to browse</span>
				<input type="file" name="files" id="fileInput" multiple hidden>
				<span class="progress" id="uploadProgress"><span class="progress-bar" id="uploadBar"></span></span>
			</label>
			` + newMenuBlock + `
			<div class="search-wrap">
				` + iconSearch + `
				<input type="search" id="searchInput" class="search-box" placeholder="Filter files ( / )">
			</div>
		</div>

		<div class="panel">
			<table class="file-table" id="fileTable">
				<thead>
					<tr>
						<th data-sort="name">Name <span class="sort-icon">` + iconChevron + `</span></th>
						<th data-sort="size">Size <span class="sort-icon">` + iconChevron + `</span></th>
						<th data-sort="date">Modified <span class="sort-icon">` + iconChevron + `</span></th>
					</tr>
				</thead>
				<tbody>
					{{ range .Files }}
					<tr class="file-row" data-name="{{ .Name }}" data-size="{{ .Size }}" data-date="{{ .ModTimeISO }}" data-dir="{{ if .IsDir }}1{{ else }}0{{ end }}">
						<td>
							<a class="file-link" href="{{ .URLPath }}">
								{{ if .IsDir }}
									<span class="file-icon dir-icon">` + iconFolder + `</span>
								{{ else if .IsImage }}
									<img class="file-thumb" src="{{ .URLPath }}" alt="" loading="lazy" onerror="this.style.display='none';this.nextElementSibling.style.display='inline-flex'">
									<span class="file-icon image-icon" style="display:none">` + iconImage + `</span>
								{{ else if .IsVideo }}
									<span class="file-icon video-icon">` + iconVideo + `</span>
								{{ else }}
									<span class="file-icon generic-icon">` + iconFile + `</span>
								{{ end }}
								<span class="file-name">{{ .Name }}</span>
							</a>
							<div class="actions">
								{{ if .EditURL }}<a class="icon-btn" href="{{ .EditURL }}" title="Edit">` + iconPencil + `</a>{{ end }}
								{{ if .ZipURL }}<a class="icon-btn" href="{{ .ZipURL }}" title="Download zip">` + iconDownload + `</a>{{ end }}
								<a class="icon-btn" href="{{ .URLPath }}" title="Copy link" onclick="copyLink(event, '{{ .URLPath }}')">` + iconLink + `</a>
								<a class="icon-btn" href="{{ .RenameURL }}" title="Rename">` + iconRename + `</a>
								<a class="icon-btn danger" href="{{ .DeleteURL }}" title="Delete">` + iconTrash + `</a>
							</div>
						</td>
						<td class="col-size">{{ .SizeStr }}</td>
						<td class="col-date"><span class="rel-time" title="{{ .ModTime }}" data-iso="{{ .ModTimeISO }}">{{ .ModTimeRel }}</span></td>
					</tr>
					{{ end }}
					{{ if not .Files }}
					` + emptyStateBlock + `
					{{ end }}
				</tbody>
			</table>
		</div>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files &middot; {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }} &middot; {{ .Stats.TotalSizeStr }} total{{ end }}</span>
			<span class="version">jotasrv v{{ .Version }}</span>
		</div>
	</div>

	<div class="toast" id="toast"></div>

	<script>
	var toastTimer;
	function showToast(msg) {
		var t = document.getElementById('toast');
		t.textContent = msg;
		t.classList.add('show');
		clearTimeout(toastTimer);
		toastTimer = setTimeout(function() { t.classList.remove('show'); }, 2200);
	}

	function copyLink(e, url) {
		e.preventDefault();
		var full = location.origin + url;
		function done() { showToast('Link copied!'); }
		if (navigator.clipboard) {
			navigator.clipboard.writeText(full).then(done, done);
		} else {
			var ta = document.createElement('textarea');
			ta.value = full;
			document.body.appendChild(ta);
			ta.select();
			try { document.execCommand('copy'); } catch (err) {}
			document.body.removeChild(ta);
			done();
		}
	}

	var searchInput = document.getElementById('searchInput');
	if (searchInput) {
		searchInput.addEventListener('input', function(e) {
			var term = e.target.value.toLowerCase();
			var rows = document.querySelectorAll('#fileTable tbody tr.file-row');
			rows.forEach(function(row) {
				var fn = row.querySelector('.file-name');
				if (!fn) return;
				row.style.display = fn.textContent.toLowerCase().indexOf(term) !== -1 ? '' : 'none';
			});
		});
		document.addEventListener('keydown', function(e) {
			if (e.key === '/' && document.activeElement.tagName !== 'INPUT' && document.activeElement.tagName !== 'TEXTAREA') {
				e.preventDefault();
				searchInput.focus();
			}
		});
	}

	var sortState = { col: null, asc: true };
	document.querySelectorAll('#fileTable th[data-sort]').forEach(function(th) {
		th.addEventListener('click', function() {
			var col = th.getAttribute('data-sort');
			if (sortState.col === col) { sortState.asc = !sortState.asc; } else { sortState.col = col; sortState.asc = true; }
			document.querySelectorAll('#fileTable th[data-sort]').forEach(function(h) {
				h.classList.remove('is-sorted');
				h.removeAttribute('data-dir');
			});
			th.classList.add('is-sorted');
			th.setAttribute('data-dir', sortState.asc ? 'asc' : 'desc');
			var tbody = document.querySelector('#fileTable tbody');
			var rows = Array.from(tbody.querySelectorAll('tr.file-row'));
			rows.sort(function(a, b) {
				var da = a.getAttribute('data-dir'), db = b.getAttribute('data-dir');
				if (da !== db) { return da === '1' ? -1 : 1; }
				var av, bv;
				if (col === 'name') { av = a.dataset.name.toLowerCase(); bv = b.dataset.name.toLowerCase(); }
				else if (col === 'size') { av = parseInt(a.dataset.size) || 0; bv = parseInt(b.dataset.size) || 0; }
				else { av = a.dataset.date; bv = b.dataset.date; }
				if (av < bv) return sortState.asc ? -1 : 1;
				if (av > bv) return sortState.asc ? 1 : -1;
				return 0;
			});
			rows.forEach(function(r) { tbody.appendChild(r); });
		});
	});

	function updateRelTimes() {
		document.querySelectorAll('.rel-time').forEach(function(el) {
			var iso = el.getAttribute('data-iso');
			if (!iso) return;
			var diff = (Date.now() - new Date(iso).getTime()) / 1000;
			if (diff < 60) { el.textContent = 'just now'; }
			else if (diff < 3600) { var m = Math.floor(diff / 60); el.textContent = m === 1 ? '1 minute ago' : m + ' minutes ago'; }
			else if (diff < 86400) { var h = Math.floor(diff / 3600); el.textContent = h === 1 ? '1 hour ago' : h + ' hours ago'; }
			else if (diff < 2592000) { var dy = Math.floor(diff / 86400); el.textContent = dy === 1 ? '1 day ago' : dy + ' days ago'; }
		});
	}
	setInterval(updateRelTimes, 60000);
	updateRelTimes();

	var dropZone = document.getElementById('dropZone');
	var fileInput = document.getElementById('fileInput');

	dropZone.addEventListener('dragover', function(e) { e.preventDefault(); dropZone.classList.add('dragover'); });
	dropZone.addEventListener('dragleave', function() { dropZone.classList.remove('dragover'); });
	dropZone.addEventListener('drop', function(e) {
		e.preventDefault();
		dropZone.classList.remove('dragover');
		uploadFiles(e.dataTransfer.files);
	});
	fileInput.addEventListener('change', function() { uploadFiles(fileInput.files); });

	function uploadFiles(files) {
		if (!files || files.length === 0) return;
		var fd = new FormData();
		for (var i = 0; i < files.length; i++) { fd.append('files', files[i]); }
		var progress = document.getElementById('uploadProgress');
		var bar = document.getElementById('uploadBar');
		var label = document.getElementById('uploadLabel');
		
		progress.style.display = 'block';
		label.textContent = 'Uploading ' + files.length + (files.length === 1 ? ' file...' : ' files...');
		
		var xhr = new XMLHttpRequest();
		xhr.open('POST', '{{ .CurrentPath }}', true);
		if (xhr.upload) {
			xhr.upload.onprogress = function(e) {
				if (e.lengthComputable) { bar.style.width = Math.round((e.loaded / e.total) * 100) + '%'; }
			};
		}
		xhr.onload = function() {
			if (xhr.status >= 200 && xhr.status < 300) { location.reload(); } 
			else { uploadError('Upload failed (' + xhr.status + ')'); }
		};
		xhr.onerror = function() { uploadError('Connection failed.'); };
		xhr.send(fd);

		function uploadError(msg) {
			progress.style.display = 'none';
			label.textContent = 'Drop files or click to browse';
			showToast(msg);
		}
	}
	</script>
</body>
</html>`

const htmlLite = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>Index of {{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body class="no-js">
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">` + iconBrand + ` jotasrv</a>
			` + breadcrumbBlock + `
		</div>
	</header>

	<div class="container">
		<div class="toolbar">
			<form action="{{ .CurrentPath }}" method="POST" enctype="multipart/form-data" class="upload-form">
				<input type="file" name="files" multiple>
				<button type="submit">Upload</button>
			</form>
			` + newMenuBlock + `
		</div>

		<div class="panel">
			<table class="file-table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Size</th>
						<th>Modified</th>
					</tr>
				</thead>
				<tbody>
					{{ range .Files }}
					<tr class="file-row">
						<td>
							<a class="file-link" href="{{ .URLPath }}">
								{{ if .IsDir }}
									<span class="file-icon dir-icon">` + iconFolder + `</span>
								{{ else if .IsImage }}
									<span class="file-icon image-icon">` + iconImage + `</span>
								{{ else if .IsVideo }}
									<span class="file-icon video-icon">` + iconVideo + `</span>
								{{ else }}
									<span class="file-icon generic-icon">` + iconFile + `</span>
								{{ end }}
								<span class="file-name">{{ .Name }}</span>
							</a>
							<div class="actions always">
								{{ if .EditURL }}<a class="icon-btn" href="{{ .EditURL }}" title="Edit">` + iconPencil + `</a>{{ end }}
								{{ if .ZipURL }}<a class="icon-btn" href="{{ .ZipURL }}" title="Download zip">` + iconDownload + `</a>{{ end }}
								<a class="icon-btn" href="{{ .RenameURL }}" title="Rename">` + iconRename + `</a>
								<a class="icon-btn danger" href="{{ .DeleteURL }}" title="Delete">` + iconTrash + `</a>
							</div>
						</td>
						<td class="col-size">{{ .SizeStr }}</td>
						<td class="col-date">{{ .ModTime }}</td>
					</tr>
					{{ end }}
					{{ if not .Files }}
					` + emptyStateBlock + `
					{{ end }}
				</tbody>
			</table>
		</div>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files &middot; {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }} &middot; {{ .Stats.TotalSizeStr }} total{{ end }}</span>
			<span class="version">jotasrv v{{ .Version }} (lite)</span>
		</div>
	</div>
</body>
</html>`

const editorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>Edit {{ .FileName }}</title>
	<style>{{ .CSS }}</style>
	<style>{{ .HighlightCSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">` + iconBrand + ` jotasrv</a>
			` + editorBreadcrumbBlock + `
		</div>
	</header>

	<div class="container">
		<form id="editForm" action="{{ .SaveURL }}" method="POST">
			<div class="edit-head">
				<h2>Editing {{ .FileName }}</h2>
				<div class="btn-row">
					<button type="submit">Save</button>
					<a href="{{ .FilePath }}" class="btn btn-secondary">Cancel</a>
				</div>
			</div>

			<div class="editor-shell">
				<pre aria-hidden="true"><code id="hlLayer"></code></pre>
				<textarea class="code" name="content" id="editor" spellcheck="false" autocomplete="off" autocapitalize="off">{{ .Content }}</textarea>
			</div>
		</form>
	</div>
	
	<script>{{ .HighlightJS }}</script>
	<script>
	(function() {
		var ta = document.getElementById('editor');
		var hl = document.getElementById('hlLayer');
		var initial = ta.value;
		var canHL = typeof hljs !== 'undefined';

		function render() {
			var text = ta.value;
			if (text.charAt(text.length - 1) === '\n') { text += ' '; }
			if (canHL) {
				try { hl.innerHTML = hljs.highlightAuto(text).value; }
				catch (err) { canHL = false; hl.textContent = ''; }
			}
		}
		function syncScroll() {
			hl.parentNode.scrollTop = ta.scrollTop;
			hl.parentNode.scrollLeft = ta.scrollLeft;
			hl.style.transform = 'translate(' + (-ta.scrollLeft) + 'px,' + (-ta.scrollTop) + 'px)';
		}
		
		if (canHL) {
			ta.classList.add('overlaid');
			render();
			ta.addEventListener('input', render);
			ta.addEventListener('scroll', syncScroll);
		}

		ta.addEventListener('keydown', function(e) {
			if (e.key === 'Tab') {
				e.preventDefault();
				var s = ta.selectionStart, en = ta.selectionEnd;
				ta.value = ta.value.substring(0, s) + '\t' + ta.value.substring(en);
				ta.selectionStart = ta.selectionEnd = s + 1;
				if (canHL) render();
			}
		});

		document.addEventListener('keydown', function(e) {
			if ((e.ctrlKey || e.metaKey) && e.key === 's') {
				e.preventDefault();
				document.getElementById('editForm').submit();
			}
		});

		window.addEventListener('beforeunload', function(e) {
			if (ta.value !== initial) { e.preventDefault(); e.returnValue = ''; }
		});
	})();
	</script>
</body>
</html>`

const editorLiteHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>Edit {{ .FileName }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">` + iconBrand + ` jotasrv</a>
			` + editorBreadcrumbBlock + `
		</div>
	</header>

	<div class="container">
		<form action="{{ .SaveURL }}" method="POST">
			<div class="edit-head">
				<h2>Editing {{ .FileName }}</h2>
				<div class="btn-row">
					<button type="submit">Save</button>
					<a href="{{ .FilePath }}" class="btn btn-secondary">Cancel</a>
				</div>
			</div>
			<div class="editor-shell">
				<textarea class="code" name="content" spellcheck="false" autocomplete="off" autocapitalize="off">{{ .Content }}</textarea>
			</div>
		</form>
	</div>
</body>
</html>`

const confirmHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>{{ .Title }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="dialog-box">
			<div class="dialog-icon warn">` + iconWarning + `</div>
			<h2 class="danger">{{ .Title }}</h2>
			<p>{{ .Message }}</p>
			<div class="btn-row">
				<form action="{{ .ActionURL }}" method="POST">
					<button type="submit" class="btn-danger">Delete</button>
				</form>
				<a href="{{ .CancelURL }}" class="btn btn-secondary">Cancel</a>
			</div>
		</div>
	</div>
</body>
</html>`

const renameHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="color-scheme" content="dark light">
	<title>Rename {{ .OldName }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">` + iconBrand + ` jotasrv</a>
			` + breadcrumbBlock + `
		</div>
	</header>

	<div class="container">
		<div class="dialog-box">
			<div class="dialog-icon info">` + iconPencil + `</div>
			<h2>Rename</h2>
			<p class="filename">{{ .OldName }}</p>
			<form action="{{ .NewURL }}" method="POST">
				<input type="text" name="newname" value="{{ .OldName }}" autofocus required>
				<div class="btn-row">
					<button type="submit">Rename</button>
					<a href="javascript:history.back()" class="btn btn-secondary">Cancel</a>
				</div>
			</form>
		</div>
	</div>
</body>
</html>`

var tmplFull = template.Must(template.New("full").Parse(htmlFull))
var tmplLite = template.Must(template.New("lite").Parse(htmlLite))
var tmplEditor = template.Must(template.New("editor").Parse(editorHTML))
var tmplEditorLite = template.Must(template.New("editorLite").Parse(editorLiteHTML))
var tmplConfirm = template.Must(template.New("confirm").Parse(confirmHTML))
var tmplRename = template.Must(template.New("rename").Parse(renameHTML))
