package controllers

import (
	"github.com/graetzer/bnginx/app/routes"
	"github.com/revel/revel"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Admin struct {
	*revel.Controller
	App
}

func (c Admin) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index(0))
	}
	return nil
}

func (c Admin) Index() revel.Result {
	var (
		posts []*Blogpost
		users []*User
	)
	DB.Order("updated_at DESC").Find(&posts)
	DB.Find(&users)
	return c.Render(posts, users)
}

// ==================== Handle Users ====================

func (c Admin) EditUser(email string) revel.Result {
	profile := new(User)
	profile.Id = -1
	if email != "create" {
		profile = c.getUser(email)
		if profile == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}

	return c.Render(profile)
}

func (c Admin) SaveUser(userId int64, name, email, oldPassword, password string) revel.Result {
	c.Validation.Required(userId)
	c.Validation.Required(email)
	c.Validation.MinSize(password, 8)
	c.Validation.MaxSize(password, 300)
	//c.Validation.Match(username, regexp.MustCompile("^\\w*$"))
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		if userId == -1 {
			return c.Redirect(routes.Admin.EditUser("create"))
		} else {
			return c.Redirect(routes.Admin.EditUser(email))
		}
	}

	var user *User
	if userId <= 0 {
		user = new(User)
		if c.getUser(email) != nil {
			c.Flash.Error("This email address is already used")
			return c.Redirect(routes.Admin.EditUser("create"))
		}
		delete(c.Flash.Out, "error")
	} else {
		u := c.connected()
		if u.Id != userId && !u.IsAdmin {
			c.Flash.Error("You have no permission to edit this profile")
			return c.Redirect(routes.Admin.EditUser("create"))
		}
		user = c.getUserById(userId)
		if user == nil || !user.CheckPassword(oldPassword) {
			c.Flash.Error("Db corruption or the password is incorrect")
			return c.Redirect(routes.Admin.EditUser(user.Email))
		}
	}
	user.Name = name
	user.Email = email
	user.SetPassword(password)
	DB.Save(user)

	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeleteUser(email string) revel.Result {
	if user := c.getUser(email); user != nil {
		u := c.connected()
		if u.Id == user.Id || u.IsAdmin {
			DB.Delete(user)
			c.Flash.Success("Deleted user")
		}
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Posts ====================

func (c Admin) EditPost(postId int64) revel.Result {
	post := c.getPostById(postId)
	if post == nil && postId > 0 {
		return c.Redirect(routes.Admin.Index())
	} else if post == nil {
		post = new (Blogpost)
		post.Title = "A new Blogpost"
		post.Body = "### Start with something\n\nE.g.\n\n1. Make a List\n2. Of Interesting\n3. Things"
	}
	return c.Render(post)
}

func (c Admin) SavePost() revel.Result {
	var post Blogpost
	c.Params.Bind(&post, "post")
	u := c.connected()
	if !DB.NewRecord(post) {// Check if the user owns this
		original := c.getPostById(post.Id)
		if original == nil || u.Id != post.UserId && !u.IsAdmin {
			c.Flash.Error("You have no permission to edit this")
			return c.Redirect(routes.Admin.Index())
		}
	}
	post.UserId = u.Id
	DB.Save(&post)
	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeletePost(postId int64) revel.Result {
	post := c.getPostById(postId)
	if post != nil {
		u := c.connected()
		if u.Id == post.UserId || u.IsAdmin {
			DB.Where("PostId = ?", post.Id).Delete(&Comment{})// Delete comments
			DB.Delete(post)
			c.Flash.Success("Deleted post")
		}
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Uploads ====================

func (c Admin) Media() revel.Result {
	basePath := filepath.Join(DataBaseDir(), "uploads/")

	fs, err := ioutil.ReadDir(basePath)
	if err != nil && os.IsNotExist(err) {
		c.Flash.Success("Creating uploads directory")
		os.MkdirAll(basePath, 0777)
	} else if err != nil {
		return c.RenderError(err)
	}

	// Exclude hidden files
	files := make([]os.FileInfo, 0, len(fs))
	for _, f := range fs {
		if !strings.HasPrefix(f.Name(), ".") {
			files = append(files, f)
		}
	}

	uploadPrefix := time.Now().Format("2006_01_02_")
	return c.Render(files, uploadPrefix)
}

func (c Admin) UploadMedia() revel.Result {
	basePath := filepath.Join(DataBaseDir(), "uploads/")

	// Loop through all uploads and save them
	for _, fInfo := range c.Params.Files["file"] {

		fi, err := fInfo.Open()
		if err != nil {
			return c.RenderError(err)
		}
		defer fi.Close()

		// Append the current date, avoid conflicts
		var uploadPrefix = time.Now().Format("2006_01_02_")
		full := filepath.Join(basePath, uploadPrefix+filepath.Base(fInfo.Filename))

		fo, err := os.Create(full)
		if err != nil {
			return c.RenderError(err)
		}
		defer fo.Close()

		if _, err = io.Copy(fo, fi); err != nil {
			return c.RenderError(err)
		}
	}
	return c.RenderJson(struct{ Message string }{"Success"})
}

func (c Admin) DeleteMedia(filename string) revel.Result {
	if c.connected().IsAdmin {
		filepath := filepath.Join(DataBaseDir(), "uploads/", filepath.Base(filename))
		err := os.Remove(filepath)
		if err != nil {
			c.Flash.Error(err.Error())
		}
	}
	return c.Redirect(routes.Admin.Media())
}

// ==================== Handle Comments ====================

func (c Admin) Comments() revel.Result {
	var comments []*Comment
	DB.Order("created_at DESC").Find(&comments)
	return c.Render(comments)
}

func (c Admin) UpdateComment(commentId int64, approved bool) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentId).RecordNotFound() {
		comment.Approved = approved
		DB.Save(&comment)
	}
	return c.Redirect(routes.Admin.Comments())
}

func (c Admin) DeleteComment(commentId int64) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentId).RecordNotFound() {
		DB.Delete(&comment)
	}
	return c.Redirect(routes.Admin.Comments())
}

// ==================== Handle Projects ====================

func (c Admin) Projects() revel.Result {
	var projects []*Project
	DB.Order("updated_at DESC").Find(&projects)
	return c.Render(projects)
}

func (c Admin) EditProject(projectId int64) revel.Result {
	var project Project
	if DB.First(&project, projectId).RecordNotFound() {
		if projectId >= 0 {
			return c.Redirect(routes.Admin.Index())
		} else {
			project = Project{Title:"My new Project"}
		}
	}
	return c.Render(project)
}

func (c Admin) SaveProject() revel.Result {
	var project Project
	c.Params.Bind(&project, "project")
	if !c.connected().IsAdmin {// Check if the user owns this
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Save(&project)
	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeleteProject(projectId int64) revel.Result {
	var project Project
	if DB.First(&project, projectId).RecordNotFound() || !c.connected().IsAdmin {
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Delete(project)
	c.Flash.Success("Deleted Project")
	return c.Redirect(routes.Admin.Projects())
}

// ==================== Handle Locations ====================

func (c Admin) Locations() revel.Result {
	var locations []*Location
	DB.Order("updated_at DESC").Find(&locations)
	return c.Render(locations)
}

func (c Admin) EditLocation(locationId int64) revel.Result {
	var location Location
	if DB.First(&location, locationId).RecordNotFound() {
		if locationId >= 0 {
			return c.Redirect(routes.Admin.Index())
		} else {
			location = Location{Name:"Current Location"}
		}
	}
	return c.Render(location)
}

func (c Admin) SaveLocation() revel.Result {
	var location Location
	c.Params.Bind(&location, "location")
	
	if !c.connected().IsAdmin {// Check if the user owns this
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Save(&location)
	return c.Redirect(routes.Admin.Locations())
}

func (c Admin) DeleteLocation(locationId int64) revel.Result {
	var location Location
	if DB.First(&location, locationId).RecordNotFound() || !c.connected().IsAdmin {
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Delete(location)
	c.Flash.Success("Deleted Location")
	return c.Redirect(routes.Admin.Index())
}
