// A small library that lets you use Pongo2 with Beego
//
// When Render is called, it will populate the render context with Beego's flash messages.
// You can also use {% urlfor "MyController.Action" ":key" "value" %} in your templates, and
// it'll work just like `urlfor` would with `html/template`. It takes one controller argument and
// zero or more key/value pairs to fill the URL.
//
package pongo2

import (
	"net/url"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	p2 "gopkg.in/flosch/pongo2.v3"
)

type Context map[string]interface{}

var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

var devMode bool

// Render takes a Beego context, template name and a Context (map[string]interface{}).
// The template is parsed and cached, and gets executed into beegoCtx's ResponseWriter.
//
// Templates are looked up in `templates/` instead of Beego's default `views/` so that
// Beego doesn't attempt to load and parse our templates with `html/template`.
func Render(beegoCtx *context.Context, tmpl string, ctx Context) {
	mutex.RLock()
	template, ok := templates[tmpl]
	mutex.RUnlock()

	if !ok || devMode {
		var err error
		template, err = p2.FromFile(beego.BConfig.WebConfig.ViewsPath + "/" + tmpl)
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		templates[tmpl] = template
		mutex.Unlock()
	}

	var pCtx p2.Context
	if ctx == nil {
		pCtx = p2.Context{}
	} else {
		pCtx = p2.Context(ctx)
	}

	if xsrf, ok := beegoCtx.GetSecureCookie(beego.BConfig.WebConfig.XSRFKey, "_xsrf"); ok {
		pCtx["_xsrf"] = xsrf
	}

	// Only override "flash" if it hasn't already been set in Context
	if _, ok := ctx["flash"]; !ok {
		if ctx == nil {
			ctx = Context{}
		}
		ctx["flash"] = readFlash(beegoCtx)
	}

	err := template.ExecuteWriter(pCtx, beegoCtx.ResponseWriter)
	if err != nil {
		panic(err)
	}
}

// readFlash is similar to beego.ReadFromRequest except that it takes a *context.Context instead
// of a *beego.Controller, and returns a map[string]string directly instead of a Beego.FlashData
// (which only has a Data field anyway).
func readFlash(ctx *context.Context) map[string]string {
	data := map[string]string{}
	if cookie, err := ctx.Request.Cookie(beego.BConfig.WebConfig.FlashName); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, "\x23"+beego.BConfig.WebConfig.FlashSeparator+"\x23")
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}
		}
		// read one time then delete it
		ctx.SetCookie(beego.BConfig.WebConfig.FlashName, "", -1, "/")
	}
	return data
}

func init() {
	devMode = beego.AppConfig.String("runmode") == "dev"
	beego.BConfig.WebConfig.AutoRender = false
}
