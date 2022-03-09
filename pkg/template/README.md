# Bhojpur Web - Template Engine

It is a `Django`-syntax like templating language applied by [Bhojpur Web](https://github.com/bhojpur/web).

## Pre-requisites

You need to install/update using `go get`:

```bash
$ go get -u github.com/bhojpur/web/pkg/template
```

## Simple Template

```django
<html>
  <head>
    <title>Our admins and users</title>
  </head>
  {# This is a short example to give you a quick overview of Bhojpur Web template's syntax. #}
  {% macro user_details(user, is_admin=false) %}
  <div class="user_item">
    <!-- Let's indicate a user's good karma -->
    <h2 {% if (user.karma>= 40) || (user.karma > calc_avg_karma(userlist)+5) %} class="karma-good"{% endif %}>

      <!-- This will call user.String() automatically, if available: -->
      {{ user }}
    </h2>

    <!-- It will print a human-readable time duration like "3 weeks ago" -->
    <p>This user registered {{ user.register_date|naturaltime }}.</p>

    <!-- Let's allow the users to write down their biography using markdown;
             we will only show the first 15 words as a preview -->
    <p>The user's biography:</p>
    <p>
      {{ user.biography|markdown|truncatewords_html:15 }}
      <a href="/user/{{ user.id }}/">read more</a>
    </p>

    {% if is_admin %}
    <p>This user is an admin!</p>
    {% endif %}
  </div>
  {% endmacro %}

  <body>
    <!-- Make use of the macro defined above to avoid repetitive HTML code
         since we want to use the same code for admins AND members -->

    <h1>Our admins</h1>
    {% for admin in adminlist %} {{ user_details(admin, true) }} {% endfor %}

    <h1>Our members</h1>
    {% for user in userlist %} {{ user_details(user) }} {% endfor %}
  </body>
</html>
```

## Key Features

- Syntax- and feature-set-compatible with Django
- Advanced `C`-like expressions
- Complex function calls within expressions.
- Easy API to create new `filters` and `tags` (including parsing arguments))
- Additional features:
  - Macros including importing macros from other files
  - Template sandboxing, directory patterns, banned tags/filters

## Caveats

### Filters

- **date** / **time**: The `date` and `time` filter are taking the `Go` specific `time-` and
`date-` format (not the Django's one) currently.
- **stringformat**: `stringformat` does **not** take Python's string format syntax as a parameter,
instead it takes Go. Essentially, `{{ 3.14|stringformat:"pi is %.2f" }}` is `fmt.Sprintf("pi is %.2f", 3.14)`.
- **escape** / **force_escape**: Unlike Django's behaviour, the `escape`-filter is applied immediately.
Therefore, there is no need for a `force_escape`-filter yet.

### Tags

- **for**: All the `forloop` fields (like `forloop.counter`) are written with a capital letter at
the beginning. For example, the `counter` can be accessed by `forloop.Counter` and the parentloop
by `forloop.Parentloop`.
- **now**: takes Go time format (see **date** and **time**-filter).

### Miscellaneous

- **not in-operator**: You can check whether a map/struct/string contains a key/field/substring
by using the in-operator (or the negation of it):
  `{% if key in map %}Key is in map{% else %}Key not in map{% endif %}` or
  `{% if !(key in map) %}Key is NOT in map{% else %}Key is in map{% endif %}`.

## Template API usage examples

Please see the documentation for a full list of provided Template API methods.

### A simple example (template string)

```go
// Compile the template first (i.e. creating the AST)
tpl, err := template.FromString("Hello, {{ name|capfirst }}!")
if err != nil {
    panic(err)
}
// Now, you can render the template with given template.Context how often you want to
out, err := tpl.Execute(template.Context{"name": "Bhojpur Web developer"})
if err != nil {
    panic(err)
}
fmt.Println(out) // Output: Hello, Bhojpur Web developer!
```

## Server-side Usage (based on template file)

```go
package main

import (
    "github.com/bhojpur/web/pkg/template"
    "net/http"
)

// Pre-compiling the templates during an application's startup using the little
// Must()-helper function. (Must() will panic, if FromFile() or FromString() will
// return with an error - that's it). It's faster to pre-compile it anywhere during
// startup and only execute the template later.
var tplExample = template.Must(template.FromFile("example.html"))

func examplePage(w http.ResponseWriter, r *http.Request) {
    // Execute the web template per HTTP request
    err := tplExample.ExecuteWriter(template.Context{"query": r.FormValue("query")}, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/", examplePage)
    http.ListenAndServe(":8080", nil)
}
```
