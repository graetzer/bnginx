package controllers

import (
	"github.com/robfig/revel"
	"github.com/russross/blackfriday"
	"html/template"
)

func init() {
	revel.OnAppStart(DBInit)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.PANIC)
	revel.InterceptMethod((*App).addUser, revel.BEFORE)
	revel.InterceptMethod((*App).addPages, revel.BEFORE)
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)

	
	revel.TemplateFuncs["markdown"] = func(input string) template.HTML { 
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(input))))
	}
	
	revel.TemplateFuncs["lt"] = func(a []interface{},b int) bool { 
		return len(a) < b
	}
}
