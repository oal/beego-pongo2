package pongo2

import (
	"github.com/astaxie/beego/context"
	p2 "github.com/flosch/pongo2"
	"sync"
)

type Context map[string]interface{}

var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

// Render takes a Beego context, template name and a Context (map[string]interface{}).
// The template is parsed and cached, and gets executed into beegoCtx's ResponseWriter.
//
// Templates are looked up in `templates/` instead of Beego's default `views/` so that
// Beego doesn't attempt to load and parse our templates with `html/template`.
func Render(beegoCtx *context.Context, tmpl string, ctx Context) {
	mutex.RLock()
	template, ok := templates[tmpl]
	mutex.RUnlock()
	if !ok {
		var err error
		template, err = p2.FromFile("templates/" + tmpl)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		templates[tmpl] = template
		mutex.Unlock()
	}

	err := template.ExecuteRW(beegoCtx.ResponseWriter, p2.Context(ctx))
	if err != nil {
		panic(err)
	}
}
