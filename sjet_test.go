package sjet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/conf"
	"github.com/shoppehub/sjet/context"
	"github.com/shoppehub/sjet/engine"
)

type Demo struct {
	Name string
	Pwd  string
}

func (d *Demo) GetName() string {
	return d.Name
}

func TestInt(t *testing.T) {
	conf.Init("")

	RegCustomFunc("demo", func(c *gin.Context) jet.Func {

		return func(a jet.Arguments) reflect.Value {
			fmt.Println(a.Get(0).Type().Kind().String())

			demo := Demo{
				Name: "123",
			}

			return reflect.ValueOf(&demo)
		}
	})

	engine := CreateWithMem()

	router := SetupRouter(engine)

	w := performRequest(router, "GET", "/")
	fmt.Println(w.Body.String())
}

func SetupRouter(engine *engine.TemplateEngine) *gin.Engine {
	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {

		templateContext := context.InitTemplateContext(engine, c)

		template := `
		{{m := array("a","11", "b","12", "c", newObjectId())}}
		{{m = arrayAppend(m, array("xxx")) }}
		{{context("result", m)}}
		`

		result, _ := RenderMemTemplate(engine, templateContext, c, "demo", template)

		fmt.Println(*templateContext.Context)

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})
	return router
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
