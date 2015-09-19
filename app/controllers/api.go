package controllers

import (
	"github.com/revel/revel"
)

type API struct {
	*revel.Controller
}
/*
func (c API) QueryBlogposts(offset int) revel.Result {
	var posts []*Blogpost
	DB.Order("updated_at DESC").Find(&posts)
	return c.RenderJson(posts)
}*/

func (c API) QueryProjects() revel.Result {
	var projects []*Project
	DB.Order("updated_at DESC").Find(&projects)
	return c.RenderJson(projects)
}

func (c API) QueryLocations() revel.Result {
	var locations []*Location
	DB.Order("created_at DESC").Find(&locations)
	return c.RenderJson(locations)
}
