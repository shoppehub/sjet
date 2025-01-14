package sjet

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/sjet/context"
	"github.com/shoppehub/sjet/engine"
	"github.com/sirupsen/logrus"

	"github.com/shoppehub/sjet/function"
)

var customFunc = make(map[string]CustomFunc)

// 创建文件路径的模板引擎
func CreateWithFile(dir string) *engine.TemplateEngine {
	e := engine.CreateWithFile(dir)
	function.InitGlobalFunc(e)
	return e
}

// 创建内存存储的模板引擎
func CreateWithMem() *engine.TemplateEngine {
	e := engine.CreateWithMem()
	function.InitGlobalFunc(e)
	return e
}

func RenderHTMLTemplate(eng *engine.TemplateEngine, c *gin.Context) {

	result, done := combineHTML(eng, c, nil)
	if done {
		return
	}
	c.Writer.Write([]byte(result))
}
func RenderHTMLTemplateWithReplaceMap(eng *engine.TemplateEngine, c *gin.Context, replaceMap map[string]string) {

	result, done := combineHTML(eng, c, replaceMap)
	if done {
		return
	}
	c.Writer.Write([]byte(result))
}
func combineHTML(eng *engine.TemplateEngine, c *gin.Context, replaceMap map[string]string) (string, bool) {
	templateContext := context.InitTemplateContext(eng, c)

	err := templateContext.FindTemplate(eng)
	if err != nil {
		if gin.Mode() == gin.ReleaseMode {
			logrus.Error("「", gin.Mode(), "环境，自动跳转404」: ", err.Error())
			c.Redirect(http.StatusFound, "/error/404")
		} else {
			logrus.Error("「", gin.Mode(), "环境，不跳转404」: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
		}
		return "", true
	}

	for key, v := range customFunc {
		templateContext.Vars.SetFunc(key, v(c))
	}

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {

			if strings.HasPrefix(err.(string), "redirect::::") {
				p := strings.ReplaceAll(err.(string), "redirect::::", "")
				c.Redirect(302, p)
				return
			}
			if strings.HasPrefix(err.(string), "exit::::") {
				return
			}
			logrus.Error(err)
		}
	}()

	buf := bytes.NewBufferString("")

	err = templateContext.Template.Execute(buf, *templateContext.Vars, nil)
	c.Status(200)
	c.Header("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		c.Writer.WriteString(err.Error())
		return "", true
	}
	result := buf.String()
	if replaceMap != nil {
		for key, value := range replaceMap {
			result = strings.Replace(result, key, value, -1)
		}
	}
	return result, false
}

func RenderMemTemplate(eng *engine.TemplateEngine, templateContext *context.TemplateContext, c *gin.Context, fnName string, fnTemplate string) (string, error) {

	loader := *eng.Loader
	if !loader.Exists(fnName) {
		loader.Set("/"+fnName+".jet", fnTemplate)
	}

	view, err := eng.Views.GetTemplate(fnName)
	// view, err := views.Parse(fnName, fun.Template)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	for key, v := range customFunc {
		templateContext.Vars.SetFunc(key, v(c))
	}

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			if strings.HasPrefix(err.(string), "redirect::::") {
				p := strings.ReplaceAll(err.(string), "redirect::::", "")
				c.Redirect(302, p)
				return
			}
			if strings.HasPrefix(err.(string), "exit::::") {
				return
			}
			logrus.Error(err)
		}
	}()
	var resp bytes.Buffer
	err = view.Execute(&resp, *templateContext.Vars, nil)
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

type CustomFunc func(c *gin.Context) jet.Func

// 注册自定义的函数
func RegCustomFunc(funcName string, v CustomFunc) {
	customFunc[funcName] = v
}
