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
	Base
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
	DB.Order("created_at DESC").Find(&posts)
	DB.Find(&users)
	return c.Render(posts, users)
}

// ==================== Handle Users ====================

func (c Admin) EditUser(email string) revel.Result {
	profile := new(User)
	profile.ID = -1
	if email != "create" {
		profile = c.getUser(email)
		if profile == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}

	return c.Render(profile)
}

func (c Admin) SaveUser(userID int64, name, email, oldPassword, password string) revel.Result {
	c.Validation.Required(userID)
	c.Validation.Required(email)
	c.Validation.MinSize(password, 8)
	c.Validation.MaxSize(password, 300)
	//c.Validation.Match(username, regexp.MustCompile("^\\w*$"))
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		if userID == 0 {
			return c.Redirect(routes.Admin.EditUser("create"))
		} else {
			return c.Redirect(routes.Admin.EditUser(email))
		}
	}

	var user *User
	if userID <= 0 {
		user = new(User)
		if c.getUser(email) != nil {
			c.Flash.Error("This email address is already used")
			return c.Redirect(routes.Admin.EditUser("create"))
		}
		delete(c.Flash.Out, "error")
	} else {
		u := c.connected()
		if u.ID != userID && !u.IsAdmin {
			c.Flash.Error("You have no permission to edit this profile")
			return c.Redirect(routes.Admin.EditUser("create"))
		}
		user = c.getUserByID(userID)
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
		if u.ID == user.ID || u.IsAdmin {
			DB.Delete(user)
			c.Flash.Success("Deleted user")
		}
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Posts ====================

func (c Admin) EditPost(postID int64) revel.Result {
	post := c.getPostByID(postID)
	if post == nil && postID > 0 {
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
		original := c.getPostByID(post.ID)
		if original == nil || u.ID != post.UserID && !u.IsAdmin {
			c.Flash.Error("You have no permission to edit this")
			return c.Redirect(routes.Admin.Index())
		}
	}
	post.UserID = u.ID
	DB.Save(&post)
	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeletePost(postID int64) revel.Result {
	post := c.getPostByID(postID)
	if post != nil {
		u := c.connected()
		if u.ID == post.UserID || u.IsAdmin {
			DB.Where("post_id = ?", post.ID).Delete(&Comment{})// Delete comments
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

func (c Admin) UpdateComment(commentID int64, approved bool) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentID).RecordNotFound() {
		comment.Approved = approved
		DB.Save(&comment)
	}
	return c.Redirect(routes.Admin.Comments())
}

func (c Admin) DeleteComment(commentID int64) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentID).RecordNotFound() {
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

func (c Admin) EditProject(projectID int64) revel.Result {
	var project Project
	if DB.First(&project, projectID).RecordNotFound() {
		if projectID > 0 {
			return c.Redirect(routes.Admin.Projects())
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
		return c.Redirect(routes.Admin.Projects())
	}
	DB.Save(&project)
	return c.Redirect(routes.Admin.Projects())
}

func (c Admin) DeleteProject(projectID int64) revel.Result {
	var project Project
	if DB.First(&project, projectID).RecordNotFound() || !c.connected().IsAdmin {
		c.Flash.Error("You have no permission to delete this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Delete(project)
	c.Flash.Success("Deleted Project")
	return c.Redirect(routes.Admin.Projects())
}

// ==================== Handle About places and stays ====================

func (c Admin) About() revel.Result {
	var (
		places []Place
		stays []Stay
	)
	DB.Order("started_at DESC").Find(&stays)
	DB.Find(&places)
	return c.Render(places, stays)
}

func (c Admin) EditPlace(placeID int64) revel.Result {
	var place Place
	if DB.First(&place, placeID).RecordNotFound() {
		if placeID > 0 {
			return c.Redirect(routes.Admin.About())
		} else {
			place = Place{ID: 0, Name:"Current Location"}
		}
	}
	return c.Render(place)
}

func (c Admin) SavePlace() revel.Result {
	var place Place
	c.Params.Bind(&place, "place")

	if !c.connected().IsAdmin {// Check if the user owns this
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.About())
	}
	DB.Save(&place)
	return c.Redirect(routes.Admin.About())
}

func (c Admin) DeletePlace(placeID int64) revel.Result {
	var place Place
	if DB.First(&place, placeID).RecordNotFound() || !c.connected().IsAdmin {
		c.Flash.Error("You have no permission to delete this")
		return c.Redirect(routes.Admin.About())
	}
	var stays []Stay
	DB.Model(&place).Related(&stays).Delete(Stay{})
	DB.Delete(place)
	c.Flash.Success("Deleted Location")
	return c.Redirect(routes.Admin.About())
}

// ==================== Handle Stays ====================

func (c Admin) EditStay(stayID int64) revel.Result {
	var (
		stay Stay
		places []Place
	)
	if DB.First(&stay, stayID).RecordNotFound() {
		if stayID > 0 {
			return c.Redirect(routes.Admin.About())
		} else {
			stay.StartedAt = time.Now()
			stay.EndedAt = stay.StartedAt
		}
	}
	DB.Find(&places)
	return c.Render(stay, places)
}

func (c Admin) SaveStay() revel.Result {
	var stay Stay
	c.Params.Bind(&stay, "stay")

	if !c.connected().IsAdmin {// Check if the user owns this
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.About())
	}
	DB.Save(&stay)
	return c.Redirect(routes.Admin.About())
}

func (c Admin) DeleteStay(stayID int64) revel.Result {
	var stay Stay
	if DB.First(&stay, stayID).RecordNotFound() || !c.connected().IsAdmin {
		c.Flash.Error("You have no permission to edit this")
		return c.Redirect(routes.Admin.Index())
	}
	DB.Delete(stay)
	c.Flash.Success("Deleted Stay")
	return c.Redirect(routes.Admin.About())
}
