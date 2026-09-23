package assets

import "embed"

// FS holds static assets served under /assets/ and bundled with the binary.
//
//go:embed tailwindcss-browser-v4.3.3.js datastar-v1.0.4.js
var FS embed.FS
