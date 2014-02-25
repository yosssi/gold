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

## Implementation Example

```go
package main

import (
	"github.com/yosssi/gold"
	"net/http"
)

// Create a generator which parses a Gold templates and
// returns a html/template package's template.
// You can have a generator cache templates by passing
// true to NewGenerator function.
var g = gold.NewGenerator(false)

func handler(w http.ResponseWriter, r *http.Request) {
	// ParseFile parses a Gold templates and
	// returns an html/template package's template.
	tpl, err := g.ParseFile("./top.gold")

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{"Title": "Gold"}

	// Call Execute method of the html/template
	// package's template.
	err := tpl.Execute(w, data)
	
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

## Syntax

### doctype

```gold
doctype html
```

becomes

```html
<!DOCTYPE html>
```

### Creating Simple Tags

```gold
div
	address
	i
	strong
```

becomes

```html
<div>
	<address></address>
	<i></i>
	<strong></strong>
</div>
```

### Putting Texts Inside Tags

```gold
p Welcome to Gold
p
	| You can insert single text.
p.
	You can insert
	multiple texts.
```

becomes

```html
<p>Welcome to Gold</p>
<p>You can insert single text.</p>
<p>
	You can insert
	multiple texts.
</p>
```

### Adding Attributes to Tags

```gold
a href=https://github.com/yosssi/gold target=_blank Gold GitHub Page
button data-action=btnaction style="font-weight: bold; font-size: 1rem;"
  | This is a button
```

becomes

```html
<a href="https://github.com/yosssi/gold" target="_blank">Gold GitHub Page</a>
<button data-action="btnaction" style="font-weight: bold; font-size: 1rem;">This is a button</button>
```

### IDs and Classes

```gold
h1#title.main-title Welcome to Gold
#container
	.wrapper
		| Hello Gold
```

becomes

```html
<h1 id="title" class="main-title">Welcome to Gold</h1>
<div id="container">
	<div class="wrapper">Hello Gold</div>
</div>
```

## APIs
