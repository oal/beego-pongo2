package pongo2

import (
	"github.com/astaxie/beego/context"
	p2 "github.com/flosch/pongo2"
	"sync"
)

type Context map[string]interface{}

var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

func Render(begoCtx *context.Context, tmpl string, ctx Context) {
	mutex.RLock()
	template, ok := templates[tmpl]
	mutex.RUnlock()
	if !ok {
		var err error
		template, err = p2.FromFile("views/" + tmpl)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		templates[tmpl] = template
		mutex.Unlock()
	}

	pongoContext := p2.Context(ctx)
	err := template.ExecuteRW(begoCtx.ResponseWriter, &pongoContext)
	if err != nil {
		panic(err)
	}
}
