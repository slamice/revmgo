package revmgo

import (
	"errors"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
)

var (
	Session  *mgo.Session // Global mgo Session
	Dial     string       // http://godoc.org/labix.org/v2/mgo#Dial
	Method   string       // clone, copy, new http://godoc.org/labix.org/v2/mgo#Session.New
	Database string       // (optional) database name to connect to
)

func AppInit() {
	// Read configuration.
	var found bool
	if Dial, found = revel.Config.String("revmgo.dial"); !found {
		// Default to 'localhost'
		Dial = "localhost"
	}
	if Method, found = revel.Config.String("db.spec"); !found {
		Method = "clone"
	} else if err := MethodError(Method); err != nil {
		revel.ERROR.Panic(err)
	}
	if Database, found = revel.Config.String("db.name"); !found {
		Database = ""
	} else if err := MethodError(Method); err != nil {
		revel.ERROR.Panic(err)
	}

	var err error
	if Session == nil {
		fmt.Println("dial: " + revel.Config.String("revmgo.dial"))
		fmt.Println("spec: " + revel.Config.String("db.spec"))
		fmt.Println("name: " + revel.Config.String("db.name"))
		// Read configuration.
		if Session, err = mgo.Dial(Dial); err != nil {
			revel.ERROR.Panic(err)
		}
	}
}

func ControllerInit() {
	revel.InterceptMethod((*MongoController).Begin, revel.BEFORE)
	revel.InterceptMethod((*MongoController).End, revel.FINALLY)
}

type MongoController struct {
	*revel.Controller
	MongoSession *mgo.Session  // named MongoSession to avoid collision with revel.Session
	Database     *mgo.Database // shortcut to the database you want to connect to by default
}

// Connect to mgo if we haven't already and return a copy/new/clone of the session and database connection
func (c *MongoController) Begin() revel.Result {
	switch Method {
	case "clone":
		c.MongoSession = Session.Clone()
	case "copy":
		c.MongoSession = Session.Copy()
	case "new":
		c.MongoSession = Session.New()
	}
	if Database != "" {
		c.Database = c.MongoSession.DB(Database)
	}

	return nil
}

// Close the controller session if we have an active one.
func (c *MongoController) End() revel.Result {
	c.MongoSession.Close()
	return nil
}

func MethodError(m string) error {
	switch m {
	case "clone", "copy", "new":
		return nil
	}
	return errors.New("revmgo: Invalid session instantiation method '%s'")
}
