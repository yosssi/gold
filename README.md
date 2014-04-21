# Gold - Template engine for Go

[![Build Status](http://128.199.249.74/github.com/yosssi/gold/status.png?branch=master)](http://128.199.249.74/github.com/yosssi/gold)
[![GoDoc](http://godoc.org/github.com/yosssi/gold?status.png)](http://godoc.org/github.com/yosssi/gold)

Gold is a template engine for [Go](http://golang.org/). This simplifies HTML coding in Go web application development. This is influenced by [Slim](http://slim-lang.com/) and [Jade](http://jade-lang.com/).

## Example

```gold
doctype html
html lang=en
  head
    title {{.Title}}
  body
    h1 Gold - Template engine for Go
    #container.wrapper
      {{if true}}
        p You can use an expression of Go text/template package in a Gold template.
      {{end}}
      p.
        Gold is a template engine for Go.
        This simplifies HTML coding in Go web application development.
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
		<h1>Gold - Template engine for Go</h1>
		<div id="container" class="wrapper">
			<p>You can use an expression of Go html/template package in a Gold template.</p>
			<p>
				Gold is a template engine for Go.
				This simplifies HTML coding in Go web application development.
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
	// ParseFile parses a Gold template and
	// returns an html/template package's template.
	tpl, err := g.ParseFile("./top.gold")

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{"Title": "Gold"}

	// Call Execute method of the html/template
	// package's template.
	err = tpl.Execute(w, data)
	
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

Following doctypes are available:


```html
doctype html
<!DOCTYPE html>

doctype xml
<?xml version="1.0" encoding="utf-8" ?>

doctype transitional
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">

doctype strict
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">

doctype frameset
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">

doctype 1.1
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

doctype basic
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">

doctype mobile
<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.2//EN" "http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd">
```

You can also use your own literal custom doctype:

```html
doctype html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN"
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN">
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

### JavaScript

```gold
javascript:
  alert(1);
  alert(2);

script type=text/javascript
  alert(3);
  alert(4);
```

becomes

```html
<script type="text/javascript">
	alert(1);
	alert(2);
</script>
<script type="text/javascript">
	alert(3);
	alert(4);
</script>
```

### Style

```gold
style
  h1 {color: red;}
  p {color: blue;}
```

becomes

```html
<style>
	h1 {color: red;}
	p {color: blue;}
</style>
```

### Comments

```gold
div
  p Hello Gold 1
  // p Hello Gold 2
//
  div
    p Hello Gold 3
    p Hello Gold 4
div
    p Hello Gold 5
```

becomes

```html
<div>
	<p>Hello Gold 1</p>
</div>
<div>
	<p>Hello Gold 5</p>
</div>
```

### Includes

Following Gold template includes `./index.gold`.

```gold
p Hello Gold
include ./index name=Foo
```

./index.gold
```gold
p Included tempalte
input type=text value=%{name}
```

### Inheritance

Gold tamplates can inherit other Gold templates as below:

parent.gold

```gold
doctype html
html
  head
    title Gold
  body
    block content
    footer
      block footer
```

child.gold

```gold
extends ./parent

block content
  #container
    | Hello Gold

block footer
  .footer
    | Copyright XXX
```

the above Gold templates generate the following HTML:

```html
<!DOCTYPE html>
<html>
	<head>
		<title>Gold</title>
	</head>
	<body>
		<div id="container">Hello Gold</div>
		<footer>
			<div class="footer">Copyright XXX</div>
		</footer>
	</body>
</html>
```

### Optional Blocks

You can set a default value to the blocks as below:

parent.gold

```gold
doctype html
html
  head
    block title
      title Default Title
```

child1.gold

```gold
extends ./parent

block title
  title Child1 Title
```

child2.gold

```gold
extends ./parent
```

child1.gold template generates the following HTML:

```html
<!DOCTYPE html>
<html>
	<head>
		<title>Child1 Title</title>
	</head>
</html>
```

child2.gold template generates the following HTML:

```html
<!DOCTYPE html>
<html>
	<head>
		<title>Default Title</title>
	</head>
</html>
```


### Expressions

You can embed [text/template](http://golang.org/pkg/text/template/) package's expressions into Gold templates because Gold template wraps this package's Template. [text/template](http://golang.org/pkg/text/template/) package's documentation describes its expressions in detail.

```gold
div
  {{if .IsProduction}}
    p This is a production code.
  {{end}}
div
  {{range .Rows}}
    p {{.}}
  {{end}}
```

### Parse template strings

You can parse template strings and load templates from memory by using the generator's `ParseString` method.

```go
package main

import (
	"net/http"

	"github.com/yosssi/gold"
)

// Create a generator which parses a Gold templates and
// returns a html/template package's template.
// You can have a generator cache templates by passing
// true to NewGenerator function.
var g = gold.NewGenerator(true)

func handler(w http.ResponseWriter, r *http.Request) {

	// template strings
	parent := `
doctype html
html
  head
    title Gold
  body
    block content
    footer
      block footer
`
	child := `
extends parent

block content
  #container
    | Hello Gold

block footer
  .footer
    | Copyright XXX
`

	stringTemplates := map[string]string{"parent": parent, "child": child}

	// ParseString parses a Gold template strings and
	// returns an html/template package's template.
	tpl, err := g.ParseString(stringTemplates, "child")

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{"Title": "Gold"}

	// Call Execute method of the html/template
	// package's template.
	err = tpl.Execute(w, data)

	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

## Docs

* [GoDoc](http://godoc.org/github.com/yosssi/gold)

## Syntax Highlighting

* [vim-gold](https://github.com/yosssi/vim-gold)
