package pongo2

import (
	"github.com/astaxie/beego/context"
	p2 "github.com/flosch/pongo2"
	"sync"
)

var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

func Render(begoCtx *context.Context, tmpl string, ctx *p2.Context) {
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

	err := template.ExecuteRW(begoCtx.ResponseWriter, ctx)
	if err != nil {
		panic(err)
	}
}
