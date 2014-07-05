beego-pongo2
============

A tiny little helper for using Pongo2 with Beego.

Documentation: http://godoc.org/github.com/oal/beego-pongo2

## Usage

```go
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/oal/beego-pongo2"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	pongo2.Render(this.Ctx, "page.html", pongo2.Context{
		"ints": []int{1, 2, 3, 4, 5},
	})
}
```