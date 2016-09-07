package api

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	// "strconv"
	"time"
)

var contentType = "application/vnd.api+json"

type Server struct ***REMOVED***
	Engine *lib.Engine
	Cancel context.CancelFunc

	Info lib.Info
***REMOVED***

// Run runs the API server.
// I'm not sure how idiomatic this is, probably not particularly...
func (s *Server) Run(ctx context.Context, addr string) ***REMOVED***
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(s.logRequestsMiddleware)
	router.Use(s.jsonErrorsMiddleware)

	router.Use(static.Serve("/", static.LocalFile("web/dist", true)))
	router.GET("/ping", func(c *gin.Context) ***REMOVED***
		c.Data(http.StatusNoContent, "", nil)
	***REMOVED***)
	v1 := router.Group("/v1")
	***REMOVED***
		v1.GET("/info", func(c *gin.Context) ***REMOVED***
			if err := jsonapi.MarshalOnePayload(c.Writer, &s.Info); err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
		***REMOVED***)
		v1.GET("/error", func(c *gin.Context) ***REMOVED***
			c.AbortWithError(500, errors.New("This is an error"))
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			if err := jsonapi.MarshalOnePayload(c.Writer, &s.Engine.Status); err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
		***REMOVED***)
		// v1.GET("/metrics", func(c *gin.Context) ***REMOVED***
		// 	metrics := make(map[string]Metric)
		// 	for m, sink := range s.Engine.Metrics ***REMOVED***
		// 		metrics[m.Name] = Metric***REMOVED***
		// 			Name:     m.Name,
		// 			Type:     MetricType(m.Type),
		// 			Contains: ValueType(m.Contains),
		// 			Data:     sink.Format(),
		// 		***REMOVED***
		// 	***REMOVED***
		// 	c.JSON(200, metrics)
		// ***REMOVED***)
		// v1.GET("/metrics/:name", func(c *gin.Context) ***REMOVED***
		// 	name := c.Param("name")
		// 	for m, sink := range s.Engine.Metrics ***REMOVED***
		// 		if m.Name != name ***REMOVED***
		// 			continue
		// 		***REMOVED***

		// 		c.JSON(200, Metric***REMOVED***
		// 			Name:     m.Name,
		// 			Type:     MetricType(m.Type),
		// 			Contains: ValueType(m.Contains),
		// 			Data:     sink.Format(),
		// 		***REMOVED***)
		// 		return
		// 	***REMOVED***
		// 	c.AbortWithError(404, errors.New("No such metric"))
		// ***REMOVED***)
		// v1.POST("/abort", func(c *gin.Context) ***REMOVED***
		// 	s.Cancel()
		// 	c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		// ***REMOVED***)
		// v1.POST("/scale", func(c *gin.Context) ***REMOVED***
		// 	vus, err := strconv.ParseInt(c.Query("vus"), 10, 64)
		// 	if err != nil ***REMOVED***
		// 		c.AbortWithError(http.StatusBadRequest, err)
		// 		return
		// 	***REMOVED***

		// 	if err := s.Engine.Scale(vus); err != nil ***REMOVED***
		// 		c.AbortWithError(http.StatusInternalServerError, err)
		// 		return
		// 	***REMOVED***

		// 	c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		// ***REMOVED***)
	***REMOVED***
	router.NoRoute(func(c *gin.Context) ***REMOVED***
		c.JSON(404, gin.H***REMOVED***"error": "Not Found"***REMOVED***)
	***REMOVED***)

	srv := graceful.Server***REMOVED***NoSignalHandling: true, Server: &http.Server***REMOVED***Addr: addr, Handler: router***REMOVED******REMOVED***
	go srv.ListenAndServe()

	<-ctx.Done()
	srv.Stop(10 * time.Second)
	<-srv.StopChan()
***REMOVED***

func (s *Server) logRequestsMiddleware(c *gin.Context) ***REMOVED***
	path := c.Request.URL.Path
	c.Next()
	log.WithField("status", c.Writer.Status()).Debugf("%s %s", c.Request.Method, path)
***REMOVED***

func (s *Server) jsonErrorsMiddleware(c *gin.Context) ***REMOVED***
	c.Header("Content-Type", contentType)
	c.Next()
	if len(c.Errors) > 0 ***REMOVED***
		var errors ErrorResponse
		for _, err := range c.Errors ***REMOVED***
			errors.Errors = append(errors.Errors, Error***REMOVED***
				Title: err.Error(),
			***REMOVED***)
		***REMOVED***
		data, _ := json.Marshal(errors)
		c.Data(c.Writer.Status(), contentType, data)
	***REMOVED***
***REMOVED***
