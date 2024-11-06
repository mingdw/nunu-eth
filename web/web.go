package web

import "embed"

//go:embed index.html header.html footer.html favicon.ico html asserts
var HtmlsFs embed.FS
