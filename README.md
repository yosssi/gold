# Gold - Template engine for Golang

[![Build Status](https://drone.io/github.com/yosssi/gold/status.png)](https://drone.io/github.com/yosssi/gold/latest)
[![Coverage Status](https://coveralls.io/repos/yosssi/gold/badge.png?branch=master)](https://coveralls.io/r/yosssi/gold?branch=master)
[![GoDoc](https://godoc.org/github.com/yosssi/gold?status.png)](https://godoc.org/github.com/yosssi/gold)

Gold is a template engine for [Golang](http://golang.org/). This simplifies HTML coding in Golang web application development. This is influenced by [Slim](http://slim-lang.com/) and [Jade](http://jade-lang.com/).

## Example

```gold
doctype html
html lang=en
	head
		title {{.Title}}
	body
		h1 Gold - Template engine for Golang
		#container.wrapper
			{{if true}}
				p You can use an expression of Golang html/template package in a Gold template.
			{{end}}
			p.
				Gold is a template engine for Golang.
				This simplifies HTML coding in Golang web application development.
		javascript:
			msg = 'Welcome to Gold!';
			alert(msg);
```


becomes

```html
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Gold</title>
	</head>
	<body>
		<h1>Gold - Template engine for Golang</h1>
		<div id="container" class="wrapper">
			<p>You can use an expression of Golang html/template package in a Gold template.</p>
			<p>
				Gold is a template engine for Golang.
				This simplifies HTML coding in Golang web application development.
			</p>
		</div>
		<script type="text/javascript">
			msg = 'Welcome to Gold!';
			alert(msg);
		</script>
	</body>
</html>
```
