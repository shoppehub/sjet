package context

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shoppehub/sjet/engine"
)

// 模板引擎配置
type TemplateContext struct {
	Vars    *jet.VarMap
	Context *map[string]interface{}

	Template *jet.Template

	TempatePath string
}

var TemplateRoot = "pages"

func (ctx *TemplateContext) FindTemplate(t *engine.TemplateEngine) error {
	// templatePath := strings.Join([]string{TemplateRoot, ctx.Module, ctx.Page, ctx.TemplName}, "/")

	var view *jet.Template
	var err error

	if view, err = t.Views.GetTemplate(TemplateRoot + "/" + ctx.TempatePath); err != nil {
		if strings.HasSuffix(ctx.TempatePath, "/") {
			ctx.TempatePath += "index"
		} else {
			ctx.TempatePath += "/index"
		}

		if view, err = t.Views.GetTemplate(TemplateRoot + "/" + ctx.TempatePath); err != nil {
			return err
		}
	}

	ctx.Template = view
	templatePath := strings.TrimPrefix(ctx.TempatePath, "/")
	ctx.Vars.Set("namespace", strings.ReplaceAll(templatePath, "/", "_"))

	return nil
}

// 初始化模板
func InitTemplateContext(t *engine.TemplateEngine, c *gin.Context) *TemplateContext {
	vars := make(jet.VarMap)
	handlerGetCtx(&vars, c)

	context := make(map[string]interface{})
	handlerContext(&vars, &context)

	handlerRoute(&vars, c)

	ctxData := TemplateContext{
		Vars:    &vars,
		Context: &context,
	}
	ctxData.TempatePath = strings.TrimPrefix(c.Request.URL.Path, "/")

	// handlerTemplateFile(c, &ctxData)

	return &ctxData
}

// 解析模板路径 /:module/:page/:templ
func handlerTemplateFile(c *gin.Context, ctx *TemplateContext) {
	ctx.TempatePath = strings.TrimPrefix(c.Request.URL.Path, "/")
}

func getParamInContext(key string, c *gin.Context, body *map[string]interface{}) interface{} {
	if val, ok := c.GetQuery(key); ok {
		return val
	}
	if val, ok := c.GetPostForm(key); ok {
		return val
	}
	if val, ok := c.Params.Get(key); ok {
		return val
	}
	bd := *body
	if value, ok := bd[key]; ok {
		return value
	}
	if value, ok := c.Get(key); ok {
		return value
	}
	return ""
}

func handlerGetCtx(vars *jet.VarMap, c *gin.Context) {

	body := make(map[string]interface{})
	if c.Request.Body != nil {
		c.ShouldBindBodyWith(&body, binding.JSON)
	}

	vars.SetFunc("getBody", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(&body)
	})

	vars.SetFunc("getCtx", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()
		return reflect.ValueOf(getParamInContext(key, c, &body))
	})

	vars.SetFunc("getCtxForInt", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c, &body)

		if val == "" {
			return reflect.ValueOf(0)
		}
		val, _ = strconv.ParseInt(val.(string), 10, 64)
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getCtxForFloat", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c, &body)

		if val == "" {
			val = float64(0)
			return reflect.ValueOf(val)
		}
		val, _ = strconv.ParseFloat(val.(string), 64)
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getCtxForBool", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c, &body)

		if val == "" {
			return reflect.ValueOf(false)
		}
		val, _ = strconv.ParseBool(val.(string))
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getURL", func(a jet.Arguments) reflect.Value {
		c.Request.URL.Host = c.Request.Host
		return reflect.ValueOf(c.Request.URL)
	})
}

func handlerContext(vars *jet.VarMap, context *map[string]interface{}) {
	vars.SetFunc("context", func(a jet.Arguments) reflect.Value {
		ctx := *context

		if a.NumOfArguments() == 1 {
			if val, ok := ctx[a.Get(0).String()]; ok {
				return reflect.ValueOf(val)
			}
			return reflect.Value{}
		}

		ctx[a.Get(0).String()] = a.Get(1).Interface()
		*context = ctx
		return reflect.Value{}
	})
}
