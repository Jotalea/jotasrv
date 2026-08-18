package main

import "html/template"

const cssContent = `
:root {
	--rosewater:#f5e0dc; --flamingo:#f2cdcd; --pink:#f5c2e7; --mauve:#cba6f7;
	--red:#f38ba8; --maroon:#eba0ac; --peach:#fab387; --yellow:#f9e2af;
	--green:#a6e3a1; --teal:#94e2d5; --sky:#89dceb; --sapphire:#74c7ec;
	--blue:#89b4fa; --lavender:#b4befe;
	--text:#cdd6f4; --subtext1:#bac2de; --subtext0:#a6adc8;
	--overlay2:#9399b2; --overlay1:#7f849c; --overlay0:#6c7086;
	--surface2:#585b70; --surface1:#45475a; --surface0:#313244;
	--base:#1e1e2e; --mantle:#181825; --crust:#11111b;
	--accent: var(--mauve);
	--mono: ui-monospace, "JetBrains Mono", "Cascadia Code", "SF Mono", Menlo, Consolas, monospace;
}
* { box-sizing: border-box; }
html { color-scheme: dark; }
body {
	margin: 0;
	background: var(--base);
	color: var(--text);
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Noto Sans", Helvetica, Arial, sans-serif;
	line-height: 1.5;
	font-size: 15px;
}
a { color: var(--blue); }
.container { max-width: 1000px; margin: 0 auto; padding: 0 20px 40px; }

/* ---------- top bar ---------- */
.topbar {
	position: sticky; top: 0; z-index: 10;
	background: rgba(24,24,37,.92);
	backdrop-filter: blur(10px);
	-webkit-backdrop-filter: blur(10px);
	border-bottom: 1px solid var(--surface0);
}
.topbar-inner {
	max-width: 1000px; margin: 0 auto; padding: 10px 20px;
	display: flex; align-items: center; gap: 14px; flex-wrap: wrap;
}
.brand {
	color: var(--accent); text-decoration: none; font-weight: 700;
	letter-spacing: .02em; display: flex; align-items: center; gap: 8px; flex: none;
}
.brand:hover { color: var(--pink); }
.breadcrumbs { display: flex; align-items: center; flex-wrap: wrap; gap: 2px; min-width: 0; font-size: .92rem; }
.breadcrumbs a {
	color: var(--subtext1); text-decoration: none; padding: 3px 8px; border-radius: 6px;
	max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.breadcrumbs a:hover { background: var(--surface0); color: var(--text); }
.breadcrumbs a.current { color: var(--text); font-weight: 600; }
.crumb-sep { color: var(--overlay0); user-select: none; }

/* ---------- toolbar ---------- */
.toolbar { display: flex; gap: 10px; margin: 20px 0 14px; flex-wrap: wrap; align-items: stretch; }
.toolbar .spacer { flex: 1; }

input[type="text"], input[type="search"] {
	background: var(--mantle); color: var(--text);
	border: 1px solid var(--surface1); border-radius: 8px;
	padding: 8px 12px; font-size: .9rem; font-family: inherit;
}
input[type="text"]:focus, input[type="search"]:focus, textarea:focus {
	outline: none; border-color: var(--accent);
	box-shadow: 0 0 0 3px rgba(203,166,247,.18);
}
input::placeholder { color: var(--overlay0); }
.search-box { flex: 1; min-width: 160px; max-width: 320px; }

button, .btn {
	display: inline-flex; align-items: center; justify-content: center; gap: 6px;
	background: var(--accent); color: var(--crust);
	border: none; border-radius: 8px; padding: 8px 16px;
	font-size: .9rem; font-weight: 600; font-family: inherit;
	cursor: pointer; text-decoration: none; white-space: nowrap;
	transition: filter .12s, background .12s;
}
button:hover, .btn:hover { filter: brightness(1.1); }
.btn-secondary { background: var(--surface0); color: var(--text); border: 1px solid var(--surface1); }
.btn-secondary:hover { background: var(--surface1); filter: none; }
.btn-danger { background: var(--red); }

/* upload */
.upload-zone {
	display: flex; align-items: center; gap: 10px; flex-wrap: wrap;
	flex: 2; min-width: 220px;
	background: var(--mantle); border: 1px dashed var(--surface2); border-radius: 10px;
	padding: 8px 14px; cursor: pointer; color: var(--subtext0); font-size: .88rem;
	transition: border-color .15s, background .15s;
}
.upload-zone:hover { border-color: var(--overlay1); }
.upload-zone.dragover { border-color: var(--accent); background: rgba(203,166,247,.08); color: var(--text); }
.upload-zone .up-icon { color: var(--accent); font-size: 1.1rem; }
.progress { display: none; width: 100%; height: 4px; background: var(--surface0); border-radius: 2px; overflow: hidden; }
.progress-bar { height: 100%; width: 0; background: var(--accent); transition: width .15s; }
.upload-form { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }

/* "+ New" dropdown (CSS-only, works without JS) */
details.menu { position: relative; }
details.menu > summary {
	list-style: none; display: inline-flex; align-items: center; gap: 6px;
	background: var(--surface0); color: var(--text); border: 1px solid var(--surface1);
	border-radius: 8px; padding: 8px 16px; font-size: .9rem; font-weight: 600;
	cursor: pointer; user-select: none; height: 100%;
}
details.menu > summary::-webkit-details-marker { display: none; }
details.menu > summary:hover, details.menu[open] > summary { background: var(--surface1); }
.menu-panel {
	position: absolute; right: 0; top: calc(100% + 6px); z-index: 20;
	background: var(--mantle); border: 1px solid var(--surface1); border-radius: 12px;
	padding: 14px; min-width: 250px; display: flex; flex-direction: column; gap: 10px;
	box-shadow: 0 14px 40px rgba(17,17,27,.65);
}
.menu-panel form { display: flex; gap: 8px; }
.menu-panel input[type="text"] { flex: 1; min-width: 0; }
.menu-panel button { padding: 8px 12px; }

/* ---------- file list ---------- */
.panel { background: var(--mantle); border: 1px solid var(--surface0); border-radius: 12px; overflow: hidden; }
.file-table { width: 100%; border-collapse: collapse; }
.file-table th {
	text-align: left; padding: 10px 16px;
	font-size: .68rem; font-weight: 600; letter-spacing: .1em; text-transform: uppercase;
	color: var(--overlay1); border-bottom: 1px solid var(--surface0);
	user-select: none; white-space: nowrap;
}
.file-table th[data-sort] { cursor: pointer; }
.file-table th[data-sort]:hover { color: var(--text); }
.sort-arrow { color: var(--accent); }
.file-table td { padding: 7px 16px; border-bottom: 1px solid var(--surface0); }
.file-table tbody tr:last-child td { border-bottom: none; }
.file-table tbody tr { transition: background .1s; }
.file-table tbody tr:hover { background: rgba(49,50,68,.45); }

.name-cell { display: flex; align-items: center; gap: 12px; min-width: 0; }
.file-link {
	flex: 1; min-width: 0; display: flex; align-items: center; gap: 12px;
	color: var(--text); text-decoration: none; font-weight: 500;
}
.file-link:hover .file-name { color: var(--accent); }
.file-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.icon-chip {
	flex: none; width: 34px; height: 34px; border-radius: 9px;
	display: flex; align-items: center; justify-content: center;
	font-size: .95rem; background: var(--surface0);
}
.chip-dir  { background: rgba(137,180,250,.15); }
.chip-video{ background: rgba(250,179,135,.15); }
.chip-text { background: rgba(148,226,213,.12); }
.file-thumb {
	flex: none; width: 34px; height: 34px; object-fit: cover;
	border-radius: 9px; background: var(--surface0);
}

.actions { display: flex; gap: 2px; flex: none; opacity: 0; transition: opacity .12s; }
.file-table tr:hover .actions, .file-table tr:focus-within .actions, .actions.always { opacity: 1; }
.actions a {
	color: var(--overlay1); text-decoration: none; font-size: .75rem; font-weight: 500;
	padding: 4px 9px; border-radius: 6px; white-space: nowrap;
}
.actions a:hover { color: var(--text); background: var(--surface0); }
.actions a.act-del:hover { color: var(--red); background: rgba(243,139,168,.12); }

.size-cell, .date-cell {
	color: var(--subtext0); font-size: .82rem; white-space: nowrap;
	font-variant-numeric: tabular-nums;
}
.size-cell { text-align: right; }
th.num { text-align: right; }

.empty-msg { text-align: center; color: var(--overlay0); padding: 56px 20px; }
.empty-msg .big { font-size: 2rem; display: block; margin-bottom: 8px; }

/* ---------- footer ---------- */
.footer {
	margin-top: 14px; padding: 0 4px; display: flex; justify-content: space-between;
	align-items: center; gap: 8px; flex-wrap: wrap;
	color: var(--overlay1); font-size: .8rem;
}
.footer .version { color: var(--overlay0); font-family: var(--mono); }

/* ---------- toast ---------- */
.toast {
	position: fixed; bottom: 24px; left: 50%; z-index: 50;
	transform: translateX(-50%) translateY(10px);
	background: var(--surface0); color: var(--text); border: 1px solid var(--surface1);
	border-radius: 10px; padding: 10px 18px; font-size: .88rem;
	opacity: 0; pointer-events: none; transition: opacity .2s, transform .2s;
	box-shadow: 0 10px 30px rgba(17,17,27,.6);
}
.toast.show { opacity: 1; transform: translateX(-50%); }

/* ---------- dialogs (confirm / rename) ---------- */
.dialog-box {
	max-width: 460px; margin: 60px auto 0;
	background: var(--mantle); border: 1px solid var(--surface0);
	border-radius: 14px; padding: 28px; text-align: center;
}
.dialog-box .dialog-icon {
	width: 52px; height: 52px; margin: 0 auto 14px; border-radius: 50%;
	display: flex; align-items: center; justify-content: center; font-size: 1.4rem;
}
.dialog-icon.warn { background: rgba(243,139,168,.14); }
.dialog-icon.info { background: rgba(203,166,247,.14); }
.dialog-box h2 { margin: 0 0 6px; font-size: 1.15rem; }
.dialog-box h2.danger { color: var(--red); }
.dialog-box p { color: var(--subtext0); margin: 0 0 22px; font-size: .92rem; word-break: break-word; }
.dialog-box .filename { color: var(--text); font-family: var(--mono); font-size: .88rem; }
.btn-row { display: flex; gap: 10px; justify-content: center; flex-wrap: wrap; }
.dialog-box input[type="text"] { width: 100%; padding: 10px 14px; font-size: 1rem; margin-bottom: 18px; font-family: var(--mono); }

/* ---------- editor ---------- */
.edit-head { display: flex; align-items: center; gap: 12px; margin: 20px 0 12px; flex-wrap: wrap; }
.edit-head h2 { margin: 0; font-size: 1.05rem; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.edit-head h2 .fn { color: var(--accent); font-family: var(--mono); }
.editor-shell { position: relative; background: var(--crust); border: 1px solid var(--surface1); border-radius: 12px; overflow: hidden; }
.editor-shell pre {
	margin: 0; position: absolute; inset: 0; overflow: hidden;
	padding: 16px; pointer-events: none;
}
.editor-shell pre code {
	font-family: var(--mono); font-size: .875rem; line-height: 1.6; tab-size: 4;
	white-space: pre; background: transparent; padding: 0; display: block;
}
textarea.code {
	position: relative; z-index: 1; display: block; width: 100%; min-height: 65vh;
	padding: 16px; border: none; outline: none; resize: vertical;
	background: transparent; color: var(--text); caret-color: var(--rosewater);
	font-family: var(--mono); font-size: .875rem; line-height: 1.6; tab-size: 4;
	white-space: pre; overflow: auto;
}
textarea.code.overlaid { color: transparent; }
textarea.code::selection { background: rgba(88,91,112,.6); color: transparent; }

@media (max-width: 640px) {
	.container { padding: 0 12px 32px; }
	.topbar-inner { padding: 10px 12px; }
	.toolbar { flex-direction: column; }
	.search-box { max-width: none; }
	details.menu > summary, .menu-panel { width: 100%; }
	.menu-panel { position: static; margin-top: 8px; box-shadow: none; }

	.file-table thead { display: none; }
	.file-table, .file-table tbody, .file-table tr, .file-table td { display: block; width: 100%; }
	.file-table td { border: none; padding: 2px 14px; }
	.file-table tr { padding: 10px 0 12px; border-bottom: 1px solid var(--surface0); }
	.file-table tbody tr:last-child { border-bottom: none; }
	.size-cell, .date-cell { display: inline-block; width: auto; padding-top: 0; font-size: .76rem; color: var(--overlay1); }
	.size-cell { padding-left: 60px; text-align: left; }
	.actions { opacity: 1; padding-left: 46px; margin-top: 2px; flex-wrap: wrap; }
	.dialog-box { margin-top: 24px; }
}
`

const htmlFull = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="theme-color" content="#181825">
	<title>{{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">&#128194; jotasrv</a>
			<nav class="breadcrumbs">
				{{ range $i, $bc := .Breadcrumbs }}
					<a href="{{ $bc.URLPath }}" class="{{ if $bc.IsLast }}current{{ end }}">{{ $bc.Name }}</a>
					{{ if not $bc.IsLast }}<span class="crumb-sep">/</span>{{ end }}
				{{ end }}
			</nav>
		</div>
	</header>

	<div class="container">
		<div class="toolbar">
			<label class="upload-zone" id="dropZone">
				<input type="file" name="files" id="fileInput" multiple hidden>
				<span class="up-icon">&#8679;</span>
				<span id="uploadLabel">drop files here or click to upload</span>
				<span class="progress" id="uploadProgress"><span class="progress-bar" id="uploadBar"></span></span>
			</label>
			<input type="search" id="searchInput" class="search-box" placeholder="filter&hellip;  ( / )">
			<details class="menu">
				<summary>new</summary>
				<div class="menu-panel">
					<form action="/mkdir{{ .CurrentPath }}" method="POST">
						<input type="text" name="dirname" placeholder="directory" required>
						<button type="submit">&#128193;</button>
					</form>
					<form action="/mkfile{{ .CurrentPath }}" method="POST">
						<input type="text" name="filename" placeholder="file" required>
						<button type="submit">&#128196;</button>
					</form>
				</div>
			</details>
		</div>

		<div class="panel">
			<table class="file-table" id="fileTable">
				<thead>
					<tr>
						<th data-sort="name">name <span class="sort-arrow"></span></th>
						<th data-sort="size" class="num">size <span class="sort-arrow"></span></th>
						<th data-sort="date" class="date-cell">modified <span class="sort-arrow"></span></th>
					</tr>
				</thead>
				<tbody>
					{{ range .Files }}
					<tr data-name="{{ .Name }}" data-size="{{ .Size }}" data-date="{{ .ModTimeISO }}" data-dir="{{ if .IsDir }}1{{ else }}0{{ end }}">
						<td>
							<div class="name-cell">
								<a class="file-link" href="{{ .URLPath }}">
									{{ if .IsImage }}
										<img class="file-thumb" src="{{ .URLPath }}" alt="" loading="lazy">
									{{ else }}
										<span class="icon-chip {{ if .IsDir }}chip-dir{{ else if .IsVideo }}chip-video{{ else }}chip-text{{ end }}">
											{{ if .IsDir }}&#128193;{{ else if .IsVideo }}&#127916;{{ else }}&#128196;{{ end }}
										</span>
									{{ end }}
									<span class="file-name">{{ .Name }}</span>
								</a>
								<div class="actions">
									{{ if .EditURL }}<a href="{{ .EditURL }}" title="Edit">edit</a>{{ end }}
									{{ if .ZipURL }}<a href="{{ .ZipURL }}" title="Download as zip">zip</a>{{ end }}
									<a href="{{ .URLPath }}" title="copy link" onclick="copyLink(event, '{{ .URLPath }}')">copy link</a>
									<a href="{{ .RenameURL }}" title="rename">rename</a>
									<a href="{{ .DeleteURL }}" title="delete" class="act-del">delete</a>
								</div>
							</div>
						</td>
						<td class="size-cell">{{ .SizeStr }}</td>
						<td class="date-cell"><span class="rel-time" title="{{ .ModTime }}" data-iso="{{ .ModTimeISO }}">{{ .ModTimeRel }}</span></td>
					</tr>
					{{ end }}
					{{ if not .Files }}
					<tr><td colspan="3" class="empty-msg"><span class="big">&#128567;</span>this directory is empty</td></tr>
					{{ end }}
				</tbody>
			</table>
		</div>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files &middot; {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }} &middot; {{ .Stats.TotalSizeStr }}{{ end }}</span>
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
		function done() { showToast('Link copied'); }
		if (navigator.clipboard && navigator.clipboard.writeText) {
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

	/* search filter */
	var searchInput = document.getElementById('searchInput');
	searchInput.addEventListener('input', function(e) {
		var term = e.target.value.toLowerCase();
		var rows = document.querySelectorAll('#fileTable tbody tr[data-name]');
		for (var i = 0; i < rows.length; i++) {
			var name = rows[i].getAttribute('data-name').toLowerCase();
			rows[i].style.display = name.indexOf(term) !== -1 ? '' : 'none';
		}
	});
	document.addEventListener('keydown', function(e) {
		if (e.key === '/' && document.activeElement.tagName !== 'INPUT' && document.activeElement.tagName !== 'TEXTAREA') {
			e.preventDefault();
			searchInput.focus();
		}
	});

	/* sort (directories always first) */
	var sortState = { col: null, asc: true };
	var ths = document.querySelectorAll('#fileTable th[data-sort]');
	for (var i = 0; i < ths.length; i++) {
		(function(th) {
			th.addEventListener('click', function() {
				var col = th.getAttribute('data-sort');
				if (sortState.col === col) { sortState.asc = !sortState.asc; }
				else { sortState.col = col; sortState.asc = true; }
				var arrows = document.querySelectorAll('#fileTable .sort-arrow');
				for (var j = 0; j < arrows.length; j++) { arrows[j].textContent = ''; }
				th.querySelector('.sort-arrow').textContent = sortState.asc ? '\u25B4' : '\u25BE';
				var tbody = document.querySelector('#fileTable tbody');
				var rows = Array.prototype.slice.call(tbody.querySelectorAll('tr[data-name]'));
				rows.sort(function(a, b) {
					var da = a.getAttribute('data-dir'), db = b.getAttribute('data-dir');
					if (da !== db) { return da === '1' ? -1 : 1; }
					var av, bv;
					if (col === 'name') { av = a.getAttribute('data-name').toLowerCase(); bv = b.getAttribute('data-name').toLowerCase(); }
					else if (col === 'size') { av = parseInt(a.getAttribute('data-size'), 10) || 0; bv = parseInt(b.getAttribute('data-size'), 10) || 0; }
					else { av = a.getAttribute('data-date'); bv = b.getAttribute('data-date'); }
					if (av < bv) return sortState.asc ? -1 : 1;
					if (av > bv) return sortState.asc ? 1 : -1;
					return 0;
				});
				for (var k = 0; k < rows.length; k++) { tbody.appendChild(rows[k]); }
			});
		})(ths[i]);
	}

	/* relative times */
	function updateRelTimes() {
		var els = document.querySelectorAll('.rel-time');
		for (var i = 0; i < els.length; i++) {
			var iso = els[i].getAttribute('data-iso');
			if (!iso) continue;
			var diff = (Date.now() - new Date(iso).getTime()) / 1000;
			var txt = null;
			if (diff < 60) { txt = 'just now'; }
			else if (diff < 3600) { var m = Math.floor(diff / 60); txt = m + (m === 1 ? ' minute ago' : ' minutes ago'); }
			else if (diff < 86400) { var h = Math.floor(diff / 3600); txt = h + (h === 1 ? ' hour ago' : ' hours ago'); }
			else if (diff < 2592000) { var d = Math.floor(diff / 86400); txt = d + (d === 1 ? ' day ago' : ' days ago'); }
			if (txt) { els[i].textContent = txt; }
		}
	}
	updateRelTimes();
	setInterval(updateRelTimes, 60000);

	/* upload */
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
		label.textContent = 'Uploading ' + files.length + (files.length === 1 ? ' file\u2026' : ' files\u2026');
		var xhr = new XMLHttpRequest();
		xhr.open('POST', '{{ .CurrentPath }}', true);
		if (xhr.upload) {
			xhr.upload.onprogress = function(e) {
				if (e.lengthComputable) { bar.style.width = Math.round(e.loaded / e.total * 100) + '%'; }
			};
		}
		xhr.onload = function() { location.reload(); };
		xhr.onerror = function() {
			progress.style.display = 'none';
			label.textContent = 'Drop files here or click to upload';
			showToast('Upload failed');
		};
		xhr.send(fd);
	}
	</script>
</body>
</html>`

const htmlLite = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="theme-color" content="#181825">
	<title>{{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">&#128194; jotasrv</a>
			<nav class="breadcrumbs">
				{{ range $i, $bc := .Breadcrumbs }}
					<a href="{{ $bc.URLPath }}" class="{{ if $bc.IsLast }}current{{ end }}">{{ $bc.Name }}</a>
					{{ if not $bc.IsLast }}<span class="crumb-sep">/</span>{{ end }}
				{{ end }}
			</nav>
		</div>
	</header>

	<div class="container">
		<div class="toolbar">
			<form class="upload-form upload-zone" action="{{ .CurrentPath }}" method="POST" enctype="multipart/form-data">
				<span class="up-icon">&#8679;</span>
				<input type="file" name="files" multiple>
				<button type="submit">upload</button>
			</form>
			<details class="menu">
				<summary>new</summary>
				<div class="menu-panel">
					<form action="/mkdir{{ .CurrentPath }}" method="POST">
						<input type="text" name="dirname" placeholder="directory" required>
						<button type="submit">&#128193;</button>
					</form>
					<form action="/mkfile{{ .CurrentPath }}" method="POST">
						<input type="text" name="filename" placeholder="file" required>
						<button type="submit">&#128196;</button>
					</form>
				</div>
			</details>
		</div>

		<div class="panel">
			<table class="file-table">
				<thead>
					<tr>
						<th>name</th>
						<th class="num">size</th>
						<th class="date-cell">modified</th>
					</tr>
				</thead>
				<tbody>
					{{ range .Files }}
					<tr>
						<td>
							<div class="name-cell">
								<a class="file-link" href="{{ .URLPath }}">
									<span class="icon-chip {{ if .IsDir }}chip-dir{{ else if .IsVideo }}chip-video{{ else }}chip-text{{ end }}">
										{{ if .IsDir }}&#128193;{{ else if .IsVideo }}&#127916;{{ else }}&#128196;{{ end }}
									</span>
									<span class="file-name">{{ .Name }}</span>
								</a>
								<div class="actions always">
									{{ if .EditURL }}<a href="{{ .EditURL }}">edit</a>{{ end }}
									{{ if .ZipURL }}<a href="{{ .ZipURL }}">zip</a>{{ end }}
									<a href="{{ .RenameURL }}">rename</a>
									<a href="{{ .DeleteURL }}" class="act-del">delete</a>
								</div>
							</div>
						</td>
						<td class="size-cell">{{ .SizeStr }}</td>
						<td class="date-cell">{{ .ModTime }}</td>
					</tr>
					{{ end }}
					{{ if not .Files }}
					<tr><td colspan="3" class="empty-msg"><span class="big">&#128567;</span>this directory is empty</td></tr>
					{{ end }}
				</tbody>
			</table>
		</div>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files &middot; {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }} &middot; {{ .Stats.TotalSizeStr }}{{ end }}</span>
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
	<meta name="theme-color" content="#181825">
	<title>Edit {{ .FileName }}</title>
	<style>{{ .CSS }}</style>
	<style>{{ .HighlightCSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">&#128194; jotasrv</a>
			<nav class="breadcrumbs">
				{{ range $i, $bc := .Breadcrumbs }}
					<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
					<span class="crumb-sep">/</span>
				{{ end }}
				<a class="current" href="{{ .FilePath }}">{{ .FileName }}</a>
			</nav>
		</div>
	</header>

	<div class="container">
		<form action="{{ .SaveURL }}" method="POST" id="editForm">
			<div class="edit-head">
				<h2>Editing <span class="fn">{{ .FileName }}</span></h2>
				<a href="{{ .FilePath }}" class="btn btn-secondary">Cancel</a>
				<button type="submit">Save</button>
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
	})();
	</script>
</body>
</html>`

const editorLiteHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="theme-color" content="#181825">
	<title>Edit {{ .FileName }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">&#128194; jotasrv</a>
			<nav class="breadcrumbs">
				{{ range $i, $bc := .Breadcrumbs }}
					<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
					<span class="crumb-sep">/</span>
				{{ end }}
				<a class="current" href="{{ .FilePath }}">{{ .FileName }}</a>
			</nav>
		</div>
	</header>

	<div class="container">
		<form action="{{ .SaveURL }}" method="POST">
			<div class="edit-head">
				<h2>Editing <span class="fn">{{ .FileName }}</span></h2>
				<a href="{{ .FilePath }}" class="btn btn-secondary">Cancel</a>
				<button type="submit">Save</button>
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
	<meta name="theme-color" content="#181825">
	<title>{{ .Title }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="dialog-box">
			<div class="dialog-icon warn">&#128465;</div>
			<h2 class="danger">{{ .Title }}</h2>
			<p>{{ .Message }}</p>
			<div class="btn-row">
				<a href="{{ .CancelURL }}" class="btn btn-secondary">Cancel</a>
				<form action="{{ .ActionURL }}" method="POST">
					<button type="submit" class="btn-danger">Delete</button>
				</form>
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
	<meta name="theme-color" content="#181825">
	<title>Rename {{ .OldName }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<header class="topbar">
		<div class="topbar-inner">
			<a class="brand" href="/">&#128194; jotasrv</a>
			<nav class="breadcrumbs">
				{{ range $i, $bc := .Breadcrumbs }}
					<a href="{{ $bc.URLPath }}" class="{{ if $bc.IsLast }}current{{ end }}">{{ $bc.Name }}</a>
					{{ if not $bc.IsLast }}<span class="crumb-sep">/</span>{{ end }}
				{{ end }}
			</nav>
		</div>
	</header>

	<div class="container">
		<div class="dialog-box">
			<div class="dialog-icon info">&#9998;</div>
			<h2>Rename</h2>
			<p class="filename">{{ .OldName }}</p>
			<form action="{{ .NewURL }}" method="POST">
				<input type="text" name="newname" value="{{ .OldName }}" autofocus required>
				<div class="btn-row">
					<a href="javascript:history.back()" class="btn btn-secondary">Cancel</a>
					<button type="submit">rename</button>
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
