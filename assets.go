package main

import "html/template"

const cssContent = `
:root {
	--base: #1e1e2e;
	--mantle: #181825;
	--surface0: #313244;
	--surface1: #45475a;
	--surface2: #585b70;
	--text: #cdd6f4;
	--subtext0: #a6adc8;
	--subtext1: #bac2de;
	--pink: #f5c2e7;
	--blue: #89b4fa;
	--green: #a6e3a1;
	--red: #f38ba8;
	--yellow: #f9e2af;
	--peach: #fab387;
	--overlay0: #6c7086;
	--overlay1: #7f849c;
	--teal: #94e2d5;
}
* { box-sizing: border-box; }
body {
	margin: 0;
	padding: 20px;
	background-color: var(--base);
	color: var(--text);
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
	line-height: 1.5;
}
.container { max-width: 960px; margin: 0 auto; }

.breadcrumbs {
	font-size: 1.2rem;
	margin-bottom: 20px;
	font-weight: 500;
}
.breadcrumbs a { color: var(--pink); text-decoration: none; }
.breadcrumbs a:hover { text-decoration: underline; }
.breadcrumbs span { color: var(--subtext0); margin: 0 8px; }

.controls {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
	margin-bottom: 20px;
	align-items: center;
}

.upload-zone {
	display: flex;
	gap: 10px;
	background: var(--surface0);
	padding: 10px 15px;
	border-radius: 8px;
	align-items: center;
	flex: 1;
	min-width: 250px;
	border: 2px dashed var(--surface1);
	transition: border-color 0.2s, background 0.2s;
	cursor: pointer;
}
.upload-zone.dragover {
	border-color: var(--pink);
	background: var(--mantle);
}
.upload-zone input[type="file"] { max-width: 200px; }
.upload-zone span { color: var(--subtext0); font-size: 0.9rem; }

button, .btn {
	background-color: var(--pink);
	color: var(--base);
	border: none;
	padding: 8px 16px;
	border-radius: 6px;
	font-weight: bold;
	cursor: pointer;
	text-decoration: none;
	display: inline-block;
	font-size: 0.9rem;
}
button:hover, .btn:hover { opacity: 0.9; }
.btn-secondary { background: var(--surface1); color: var(--text); }
.btn-danger { background: var(--red); }
.btn-small { padding: 4px 10px; font-size: 0.8rem; }

.search-box {
	padding: 8px 12px;
	border-radius: 6px;
	border: 1px solid var(--surface1);
	background: var(--surface0);
	color: var(--text);
	width: 100%;
	max-width: 250px;
}
.search-box:focus { outline: 2px solid var(--pink); }

.create-form {
	display: flex;
	gap: 8px;
	align-items: center;
}
.create-form input[type="text"] {
	padding: 8px 12px;
	border-radius: 6px;
	border: 1px solid var(--surface1);
	background: var(--surface0);
	color: var(--text);
	width: 180px;
}

.file-table {
	width: 100%;
	border-collapse: collapse;
	background: var(--surface0);
	border-radius: 8px;
	overflow: hidden;
}
.file-table th, .file-table td {
	padding: 10px 15px;
	text-align: left;
	border-bottom: 1px solid var(--surface1);
}
.file-table th {
	background: var(--surface1);
	color: var(--subtext0);
	font-weight: 600;
	cursor: pointer;
	user-select: none;
	white-space: nowrap;
}
.file-table th:hover { color: var(--text); }
.file-table th .sort-arrow { margin-left: 4px; font-size: 0.75rem; }
.file-table tr:last-child td { border-bottom: none; }
.file-table tr:hover { background: var(--surface1); }

.file-link {
	color: var(--text);
	text-decoration: none;
	display: flex;
	align-items: center;
	gap: 8px;
	font-weight: 500;
}
.file-link:hover { color: var(--pink); }
.file-link img.file-thumb {
	width: 24px;
	height: 24px;
	object-fit: cover;
	border-radius: 4px;
	flex-shrink: 0;
}
.icon { font-size: 1.2rem; flex-shrink: 0; }
.dir-icon { color: var(--blue); }
.image-icon { color: var(--green); }
.video-icon { color: var(--peach); }
.text-icon { color: var(--teal); }

.actions {
	display: flex;
	gap: 4px;
	align-items: center;
	opacity: 0.3;
	transition: opacity 0.15s;
}
.file-table tr:hover .actions { opacity: 1; }
.actions a, .actions button {
	background: var(--surface1);
	color: var(--subtext0);
	padding: 2px 6px;
	border-radius: 4px;
	font-size: 0.75rem;
	text-decoration: none;
	border: none;
	cursor: pointer;
	font-weight: 500;
}
.actions a:hover, .actions button:hover { color: var(--pink); background: var(--surface2); }
.actions .btn-del:hover { color: var(--red); }

.footer {
	margin-top: 20px;
	padding: 12px 15px;
	background: var(--surface0);
	border-radius: 8px;
	color: var(--subtext0);
	font-size: 0.85rem;
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.footer .version { color: var(--overlay0); }

.confirm-box {
	max-width: 500px;
	margin: 40px auto;
	background: var(--surface0);
	padding: 30px;
	border-radius: 8px;
	text-align: center;
}
.confirm-box h2 { margin-top: 0; color: var(--red); }
.confirm-box p { color: var(--subtext1); margin-bottom: 20px; }
.confirm-box .btn-row { display: flex; gap: 10px; justify-content: center; }

.rename-box, .edit-box {
	max-width: 700px;
	margin: 20px auto;
	background: var(--surface0);
	padding: 20px;
	border-radius: 8px;
}
.rename-box h2, .edit-box h2 { margin-top: 0; color: var(--pink); }
.rename-box input[type="text"] {
	width: 100%;
	padding: 10px 14px;
	border-radius: 6px;
	border: 1px solid var(--surface1);
	background: var(--mantle);
	color: var(--text);
	font-size: 1rem;
	margin-bottom: 15px;
}
.edit-box textarea {
	width: 100%;
	min-height: 500px;
	padding: 14px;
	border-radius: 6px;
	border: 1px solid var(--surface1);
	background: var(--mantle);
	color: var(--text);
	font-family: "JetBrains Mono", "Fira Code", "SF Mono", "Cascadia Code", monospace;
	font-size: 0.9rem;
	resize: vertical;
	line-height: 1.6;
	tab-size: 4;
}

.empty-msg {
	text-align: center;
	color: var(--subtext0);
	padding: 40px;
}

@media (max-width: 600px) {
	.date-col { display: none; }
	.file-table td, .file-table th { padding: 10px; }
	.controls { flex-direction: column; }
	.search-box { max-width: 100%; }
	.actions { opacity: 1; }
	.footer { flex-direction: column; gap: 5px; text-align: center; }
}
`

const htmlFull = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Index of {{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="breadcrumbs">
			{{ range $i, $bc := .Breadcrumbs }}
				<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
				{{ if not $bc.IsLast }}<span>/</span>{{ end }}
			{{ end }}
		</div>

		<div class="controls">
			<div class="upload-zone" id="dropZone">
				<input type="file" name="files" id="fileInput" multiple>
				<span>Drop files here or click to select</span>
			</div>
			<div class="create-form">
				<form action="/mkdir{{ .CurrentPath }}" method="POST" style="display:flex;gap:6px;align-items:center">
					<input type="text" name="dirname" placeholder="New folder" required>
					<button type="submit" class="btn-small">+ Folder</button>
				</form>
				<form action="/mkfile{{ .CurrentPath }}" method="POST" style="display:flex;gap:6px;align-items:center">
					<input type="text" name="filename" placeholder="New file" required>
					<button type="submit" class="btn-small">+ File</button>
				</form>
			</div>
			<input type="search" id="searchInput" class="search-box" placeholder="Filter files...">
		</div>

		<table class="file-table" id="fileTable">
			<thead>
				<tr>
					<th data-sort="name">name <span class="sort-arrow"></span></th>
					<th data-sort="size">size <span class="sort-arrow"></span></th>
					<th class="date-col" data-sort="date">modified <span class="sort-arrow"></span></th>
				</tr>
			</thead>
			<tbody>
				{{ range .Files }}
				<tr data-name="{{ .Name }}" data-size="{{ .Size }}" data-date="{{ .ModTimeISO }}">
					<td>
						<a class="file-link" href="{{ .URLPath }}">
							{{ if .IsImage }}
								<img class="file-thumb" src="{{ .URLPath }}" alt="" loading="lazy">
							{{ else }}
								<span class="icon {{ if .IsDir }}dir-icon{{ else if .IsVideo }}video-icon{{ else }}text-icon{{ end }}">
									{{ if .IsDir }}&#128193;{{ else if .IsVideo }}&#127910;{{ else }}&#128196;{{ end }}
								</span>
							{{ end }}
							<span class="file-name">{{ .Name }}</span>
						</a>
						<div class="actions">
							{{ if .EditURL }}<a href="{{ .EditURL }}" title="Edit">&#9998; edit</a>{{ end }}
							{{ if .ZipURL }}<a href="{{ .ZipURL }}" title="Download as zip">&#128230; zip</a>{{ end }}
							<a href="{{ .URLPath }}" title="Copy link" onclick="copyLink(event, '{{ .URLPath }}')">&#128279; copy</a>
							<a href="{{ .RenameURL }}" title="Rename">&#9999; rename</a>
							<a href="{{ .DeleteURL }}" title="Delete" class="btn-del">&#10005; delete</a>
						</div>
					</td>
					<td>{{ .SizeStr }}</td>
					<td class="date-col"><span class="rel-time" title="{{ .ModTime }}">{{ .ModTimeRel }}</span></td>
				</tr>
				{{ end }}
				{{ if not .Files }}
				<tr><td colspan="3" class="empty-msg">Directory is empty</td></tr>
				{{ end }}
			</tbody>
		</table>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files, {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }}, {{ .Stats.TotalSizeStr }} total{{ end }}</span>
			<span class="version">jotasrv v{{ .Version }}</span>
		</div>
	</div>

	<script>
	function copyLink(e, url) {
		e.preventDefault();
		var full = location.origin + url;
		if (navigator.clipboard) {
			navigator.clipboard.writeText(full);
		} else {
			var ta = document.createElement('textarea');
			ta.value = full;
			document.body.appendChild(ta);
			ta.select();
			document.execCommand('copy');
			document.body.removeChild(ta);
		}
	}

	document.getElementById('searchInput').addEventListener('input', function(e) {
		var term = e.target.value.toLowerCase();
		var rows = document.querySelectorAll('#fileTable tbody tr');
		rows.forEach(function(row) {
			var fn = row.querySelector('.file-name');
			if (!fn) return;
			row.style.display = fn.textContent.toLowerCase().indexOf(term) !== -1 ? '' : 'none';
		});
	});

	var sortState = { col: null, asc: true };
	document.querySelectorAll('#fileTable th[data-sort]').forEach(function(th) {
		th.addEventListener('click', function() {
			var col = th.getAttribute('data-sort');
			if (sortState.col === col) { sortState.asc = !sortState.asc; } else { sortState.col = col; sortState.asc = true; }
			document.querySelectorAll('#fileTable th .sort-arrow').forEach(function(a) { a.textContent = ''; });
			th.querySelector('.sort-arrow').textContent = sortState.asc ? ' \u25B2' : ' \u25BC';
			var tbody = document.querySelector('#fileTable tbody');
			var rows = Array.from(tbody.querySelectorAll('tr[data-name]'));
			rows.sort(function(a, b) {
				var av, bv;
				if (col === 'name') { av = a.dataset.name.toLowerCase(); bv = b.dataset.name.toLowerCase(); }
				else if (col === 'size') { av = parseInt(a.dataset.size)||0; bv = parseInt(b.dataset.size)||0; }
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
			var t = el.getAttribute('title');
			if (!t) return;
			var d = new Date(t);
			var now = new Date();
			var diff = (now - d) / 1000;
			if (diff < 60) { el.textContent = 'just now'; }
			else if (diff < 3600) { var m = Math.floor(diff/60); el.textContent = m === 1 ? '1 minute ago' : m + ' minutes ago'; }
			else if (diff < 86400) { var h = Math.floor(diff/3600); el.textContent = h === 1 ? '1 hour ago' : h + ' hours ago'; }
			else if (diff < 2592000) { var dy = Math.floor(diff/86400); el.textContent = dy === 1 ? '1 day ago' : dy + ' days ago'; }
		});
	}
	setInterval(updateRelTimes, 60000);

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
		var xhr = new XMLHttpRequest();
		xhr.open('POST', '{{ .CurrentPath }}', true);
		xhr.onload = function() { location.reload(); };
		xhr.onerror = function() { alert('Upload failed'); };
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
	<title>Index of {{ .CurrentPath }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="breadcrumbs">
			{{ range $i, $bc := .Breadcrumbs }}
				<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
				{{ if not $bc.IsLast }}<span>/</span>{{ end }}
			{{ end }}
		</div>

		<div class="controls">
			<form action="{{ .CurrentPath }}" method="POST" enctype="multipart/form-data" style="display:flex;gap:8px;align-items:center">
				<input type="file" name="files" multiple>
				<button type="submit">Upload</button>
			</form>
			<div class="create-form">
				<form action="/mkdir{{ .CurrentPath }}" method="POST" style="display:flex;gap:6px;align-items:center">
					<input type="text" name="dirname" placeholder="New folder" required>
					<button type="submit" class="btn-small">+ Folder</button>
				</form>
				<form action="/mkfile{{ .CurrentPath }}" method="POST" style="display:flex;gap:6px;align-items:center">
					<input type="text" name="filename" placeholder="New file" required>
					<button type="submit" class="btn-small">+ File</button>
				</form>
			</div>
		</div>

		<table class="file-table" id="fileTable">
			<thead>
				<tr>
					<th>name</th>
					<th>size</th>
					<th class="date-col">modified</th>
				</tr>
			</thead>
			<tbody>
				{{ range .Files }}
				<tr>
					<td>
						<a class="file-link" href="{{ .URLPath }}">
							<span class="icon {{ if .IsDir }}dir-icon{{ else if .IsVideo }}video-icon{{ else }}text-icon{{ end }}">
								{{ if .IsDir }}&#128193;{{ else if .IsVideo }}&#127910;{{ else }}&#128196;{{ end }}
							</span>
							<span class="file-name">{{ .Name }}</span>
						</a>
						<div class="actions" style="opacity:1">
							{{ if .EditURL }}<a href="{{ .EditURL }}">edit</a>{{ end }}
							{{ if .ZipURL }}<a href="{{ .ZipURL }}">zip</a>{{ end }}
							<a href="{{ .RenameURL }}">rename</a>
							<a href="{{ .DeleteURL }}">delete</a>
						</div>
					</td>
					<td>{{ .SizeStr }}</td>
					<td class="date-col">{{ .ModTime }}</td>
				</tr>
				{{ end }}
				{{ if not .Files }}
				<tr><td colspan="3" class="empty-msg">Directory is empty</td></tr>
				{{ end }}
			</tbody>
		</table>

		<div class="footer">
			<span>{{ .Stats.FileCount }} files, {{ .Stats.DirCount }} folders{{ if .Stats.TotalSizeStr }}, {{ .Stats.TotalSizeStr }} total{{ end }}</span>
			<span class="version">jotasrv v{{ .Version }}</span>
		</div>
	</div>
</body>
</html>`

const editorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Edit {{ .FileName }}</title>
	<style>{{ .CSS }}</style>
	<style>{{ .HighlightCSS }}</style>
</head>
<body>
	<div class="container">
		<div class="breadcrumbs">
			{{ range $i, $bc := .Breadcrumbs }}
				<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
				{{ if not $bc.IsLast }}<span>/</span>{{ end }}
			{{ end }}
			<span>/</span>
			<span style="color:var(--text)">{{ .FileName }}</span>
		</div>

		<div class="edit-box">
			<h2>Editing {{ .FileName }}</h2>
			<form action="{{ .SaveURL }}" method="POST">
				<textarea name="content" spellcheck="false" id="editor">{{ .Content }}</textarea>
				<br><br>
				<button type="submit">Save</button>
				<a href="{{ .FilePath }}" class="btn btn-secondary" style="margin-left:8px">Cancel</a>
			</form>
		</div>
	</div>
	<script>{{ .HighlightJS }}</script>
	<script>hljs.highlightElement(document.getElementById('editor'));</script>
</body>
</html>`

const confirmHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{ .Title }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="confirm-box">
			<h2>{{ .Title }}</h2>
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
	<title>Rename {{ .OldName }}</title>
	<style>{{ .CSS }}</style>
</head>
<body>
	<div class="container">
		<div class="breadcrumbs">
			{{ range $i, $bc := .Breadcrumbs }}
				<a href="{{ $bc.URLPath }}">{{ $bc.Name }}</a>
				{{ if not $bc.IsLast }}<span>/</span>{{ end }}
			{{ end }}
		</div>

		<div class="rename-box">
			<h2>Rename {{ .OldName }}</h2>
			<form action="{{ .NewURL }}" method="POST">
				<input type="text" name="newname" value="{{ .OldName }}" autofocus required>
				<button type="submit">Rename</button>
				<a href="javascript:history.back()" class="btn btn-secondary" style="margin-left:8px">Cancel</a>
			</form>
		</div>
	</div>
</body>
</html>`

var tmplFull = template.Must(template.New("full").Parse(htmlFull))
var tmplLite = template.Must(template.New("lite").Parse(htmlLite))
var tmplEditor = template.Must(template.New("editor").Parse(editorHTML))
var tmplConfirm = template.Must(template.New("confirm").Parse(confirmHTML))
var tmplRename = template.Must(template.New("rename").Parse(renameHTML))
var tmplEditorLite = tmplEditor
