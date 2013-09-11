package controllers

import (
	//"code.google.com/p/go.crypto/bcrypt"
	"github.com/coopernurse/gorp"
	"github.com/robfig/revel"
    "github.com/graetzer/bnginx/app/models"
	"database/sql"
)

var (
    db *sql.DB
    driver string
    spec string
	Dbm *gorp.DbMap
)

func DBInit() {
    // Read configuration.
    var found bool
    if driver, found = revel.Config.String("db.driver"); !found {
        revel.ERROR.Fatal("No db.driver found.")
    }
    if spec, found = revel.Config.String("db.spec"); !found {
        revel.ERROR.Fatal("No db.spec found.")
    }

    // Open a connection.
    var err error
    db, err = sql.Open(driver, spec)
    if err != nil {
        revel.ERROR.Fatal(err)
    }

    Dbm = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	if revel.DevMode {
		Dbm.TraceOn("[gorp]", revel.INFO)
	}
    
    
    Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
    Dbm.AddTable(models.Post{}).SetKeys(true, "PostId")
	Dbm.AddTable(models.Comment{}).SetKeys(true, "CommentId")
    Dbm.CreateTablesIfNotExists()
    
    users, err := Dbm.Select(models.User{}, "SELECT * FROM User")
    if err != nil {
        revel.ERROR.Fatal("Error checking for users "+err.Error())
    }
    if len(users) == 0 {
        user := &models.User{UserId:0, Name:"Simon", Email:"simon@graetzer.org", Password:models.HashPassword("default"), IsAdmin:true}
        Dbm.Insert(user)
    }
}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
