package pkg

import "net/http"

type Context struct {
	Continue   bool                //whether to continue or not
	Req        *http.Request       //access to the native request
	Write      http.ResponseWriter //our response writer
	Resource   *Resource           //the resource we want to query
	Parameters map[string]string   //any request parameters to apply
	Config     *Configuration      //our configuration obj
	Response   *Response           //our response
}

//sets our access control headers
func AccessHeaders(c *Context) {
	c.Write.Header().Set("Access-Control-Allow-Origin", "*")
}

//our filter to checck permissions
func Permissions(c *Context) {
	switch c.Req.Method {
	case "GET":
		if c.Config.GetPermissions["global"] != "allow" {
			c.Continue = false
		}
	case "PUT":
		if c.Config.PutPermissions["global"] != "allow" {
			c.Continue = false
		}
	case "POST":
		if c.Config.PostPermissions["global"] != "allow" {
			c.Continue = false
		}
	case "DELETE":
		if c.Config.PostPermissions["global"] != "allow" {
			c.Continue = false
		}
	}
	if c.Continue == false {
		MessageResponse(c.Write, 401, "Permission denied")
	}
}

