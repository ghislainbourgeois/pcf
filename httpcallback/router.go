// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package httpcallback

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/omec-project/logger_util"
	"github.com/omec-project/pcf/logger"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/npcf-callback/v1")
	// https://localhost:29507/npcf-callback/v1/route
	for _, route := range routes {
		switch route.Method {
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		}
	}
	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},

	{
		"HTTPNudrNotify",
		strings.ToUpper("Post"),
		"/nudr-notify/:supi",
		HTTPNudrNotify,
	},

	{
		"HTTPAmfStatusChangeNotify",
		strings.ToUpper("Post"),
		"/amfstatus",
		HTTPAmfStatusChangeNotify,
	},
}
