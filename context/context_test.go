package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {

	str := "/index"

	paths := strings.Split(strings.TrimPrefix(str, "/"), "/")

	pathSize := len(paths)
	fmt.Println(pathSize, paths)

	for i := 1; i < pathSize; i++ {
		fmt.Println(i, paths[i])
	}

	fmt.Println(path.Dir("a/b/c"))

}

func TestHelloWorld(t *testing.T) {
	body := gin.H{
		"Hello": "World",
	}

	router := SetupRouter()

	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["Hello"]

	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["Hello"], value)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {

		fmt.Println("path", c.Request.URL.Path)

		c.JSON(http.StatusOK, gin.H{
			"Hello": "World",
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

func TestRegexp(t *testing.T) {
	src := "ATS-8; Tert-Butyl (4R,6R)-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate; (4R-Cis)-1,1-Dimethylethyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate; (4R,6R)-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetic Acid Tert-Butyl Ester; Tert-Butyl 6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxolane-4-Acetate; (4R,6R)-T-Butyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate; (4R,Cis)-1,1-Dimethylethyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate; (4R,3R)-Tert-Butyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate; (4R-Cis)-1,1-Dimethylethyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate (Ats-8); Tert-Butyl[(4R,6R)-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxan-4-Yl]Acetate; Atorvastatin Calcium Intermediate ats-8; ATS-8: (4R-cis)-1,1-Dimethylethyl-6-Cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate"

	//keyword := `(?i)\b`+function.AddEscapedChar("(4R,6R)-6-cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate")+`\b`
	keyword := `(?i)\(4R,6R\)-6-cyanomethyl-2,2-Dimethyl-1,3-Dioxane-4-Acetate`
	fmt.Println(keyword)
	regExp := regexp.MustCompile(keyword)

	repl := `<span style="color: red;">$0</span>`

	result := regExp.ReplaceAllString(src, repl)
	fmt.Println(result)

	// 特殊字符的查找
	//reg = regexp.MustCompile(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]\|]`)
	//fmt.Printf("%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))
}
