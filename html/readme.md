# HTML entity

HTML entities utils. Allows to escape/unescape entities a bit faster than standard module.

Designed to use in highload projects.

## Usage

In general way it can be used similar to default module:
```go
import "github.com/koykov/entity/html"

s := "foo & > < bar"
r := html.Escape(s) // "foo &amp; &gt; &lt; bar"

x := "&# &#x &#128;43 &copy = &#169f = &#xa9"
e := html.Unescape(x) // "&# &#x €43 © = ©f = ©"
```

Buffered version:
```go

```
