package templates

// Константа для HTML-шаблона.
const MetricsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Metrics List</title>
</head>
<body>
	<h1>Metrics List</h1>
	<ul>
	{{range .}}
		<li>{{.ID}}: {{.Value}}</li>
	{{else}}
		<li>No metrics available</li>
	{{end}}
	</ul>
</body>
</html>`
