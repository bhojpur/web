<!DOCTYPE html>
<html>
  <head>
    <title>Bhojpur Web - Welcome Template</title>
  </head>
  <body>

	{{template "block"}}
	{{template "header"}}
	{{template "blocks/block.tpl"}}

	<h2>{{ .Title }}</h2>
	<p> This is SomeVar: {{ .SomeVar }}</p>
  </body>
</html>