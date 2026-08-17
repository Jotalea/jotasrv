package main

import _ "embed"

//go:embed static/highlight.min.js
var hlJS string

//go:embed static/highlight.min.css
var hlCSS string
