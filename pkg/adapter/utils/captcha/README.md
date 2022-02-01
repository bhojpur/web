# Bhojpur Web - Adaptor Captcha

An example for use of Captcha

```
package controllers

import (
	bhojpur "github.com/bhojpur/web/pkg/engine"
	"github.com/bhojpur/web/pkg/cache"
	"github.com/bhojpur/web/pkg/utils/captcha"
)

var cpt *captcha.Captcha

func init() {
	// use bhojpur cache system store the captcha data
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
}

type MainController struct {
	bhojpur.Controller
}

func (this *MainController) Get() {
	this.TplName = "index.tpl"
}

func (this *MainController) Post() {
	this.TplName = "index.tpl"

	this.Data["Success"] = cpt.VerifyReq(this.Ctx.Request)
}
```

template usage

```
{{.Success}}
<form action="/" method="post">
	{{create_captcha}}
	<input name="captcha" type="text">
</form>
```