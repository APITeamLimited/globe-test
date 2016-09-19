package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/loadimpact/speedboat/lib"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"strconv"
	// "strconv"
	"time"
)

var (
	contentType = "application/vnd.api+json"
	webBox      = rice.MustFindBox("../web/dist")
)

type Server struct ***REMOVED***
	Engine *lib.Engine
	Info   lib.Info
***REMOVED***

// Run runs the API server.
// I'm not sure how idiomatic this is, probably not particularly...
func (s *Server) Run(ctx context.Context, addr string) ***REMOVED***
	indexData, err := webBox.Bytes("index.html")
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't load index.html; web UI unavailable")
	***REMOVED***

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(s.logRequestsMiddleware)
	router.Use(s.jsonErrorsMiddleware)

	// router.Use(static.Serve("/", static.LocalFile("web/dist", true)))
	router.GET("/ping", func(c *gin.Context) ***REMOVED***
		c.Data(http.StatusNoContent, "", nil)
	***REMOVED***)
	v1 := router.Group("/v1")
	***REMOVED***
		v1.GET("/info", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(s.Info)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/error", func(c *gin.Context) ***REMOVED***
			c.AbortWithError(500, errors.New("This is an error"))
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(s.Engine.Status)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.PATCH("/status", func(c *gin.Context) ***REMOVED***
			var status lib.Status
			data, _ := ioutil.ReadAll(c.Request.Body)
			if err := jsonapi.Unmarshal(data, &status); err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			if status.VUsMax.Valid ***REMOVED***
				if status.VUsMax.Int64 < s.Engine.Status.VUs.Int64 ***REMOVED***
					if status.VUsMax.Int64 >= status.VUs.Int64 ***REMOVED***
						s.Engine.SetVUs(status.VUs.Int64)
					***REMOVED*** else ***REMOVED***
						c.AbortWithError(http.StatusBadRequest, lib.ErrMaxTooLow)
						return
					***REMOVED***
				***REMOVED***

				if err := s.Engine.SetMaxVUs(status.VUsMax.Int64); err != nil ***REMOVED***
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
			***REMOVED***
			if status.VUs.Valid ***REMOVED***
				if status.VUs.Int64 > s.Engine.Status.VUsMax.Int64 ***REMOVED***
					c.AbortWithError(http.StatusBadRequest, lib.ErrTooManyVUs)
					return
				***REMOVED***

				if err := s.Engine.SetVUs(status.VUs.Int64); err != nil ***REMOVED***
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
			***REMOVED***
			if status.Running.Valid ***REMOVED***
				s.Engine.SetRunning(status.Running.Bool)
			***REMOVED***

			data, err := jsonapi.Marshal(s.Engine.Status)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics", func(c *gin.Context) ***REMOVED***
			metrics := make([]interface***REMOVED******REMOVED***, 0, len(s.Engine.Metrics))
			for metric, sink := range s.Engine.Metrics ***REMOVED***
				metric.Sample = sink.Format()
				metrics = append(metrics, metric)
			***REMOVED***
			data, err := jsonapi.Marshal(metrics)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics/:id", func(c *gin.Context) ***REMOVED***
			id := c.Param("id")
			for metric, sink := range s.Engine.Metrics ***REMOVED***
				if metric.Name != id ***REMOVED***
					continue
				***REMOVED***
				metric.Sample = sink.Format()
				data, err := jsonapi.Marshal(metric)
				if err != nil ***REMOVED***
					c.AbortWithError(500, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			c.AbortWithError(404, errors.New("Metric not found"))
		***REMOVED***)
		v1.GET("/groups", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(s.Engine.Runner.GetGroups())
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/groups/:id", func(c *gin.Context) ***REMOVED***
			id, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			for _, group := range s.Engine.Runner.GetGroups() ***REMOVED***
				if group.ID != id ***REMOVED***
					continue
				***REMOVED***

				data, err := jsonapi.Marshal(group)
				if err != nil ***REMOVED***
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			c.AbortWithError(404, errors.New("Group not found"))
		***REMOVED***)
		v1.GET("/tests", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(s.Engine.Runner.GetTests())
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/tests/:id", func(c *gin.Context) ***REMOVED***
			id, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			for _, test := range s.Engine.Runner.GetTests() ***REMOVED***
				if test.ID != id ***REMOVED***
					continue
				***REMOVED***

				data, err := jsonapi.Marshal(test)
				if err != nil ***REMOVED***
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			c.AbortWithError(404, errors.New("Group not found"))
		***REMOVED***)
	***REMOVED***
	router.NoRoute(func(c *gin.Context) ***REMOVED***
		requestPath := c.Request.URL.Path
		bytes, err := webBox.Bytes(requestPath)
		if err != nil ***REMOVED***
			log.WithError(err).Debug("Falling back to index.html")
			if indexData == nil ***REMOVED***
				c.String(404, "Web UI is unavailable - see console output.")
				return
			***REMOVED***
			c.Data(200, "text/html; charset=utf-8", indexData)
			return
		***REMOVED***

		mimeType := mime.TypeByExtension(path.Ext(requestPath))
		if mimeType == "" ***REMOVED***
			mimeType = "application/octet-stream"
		***REMOVED***
		c.Data(200, mimeType, bytes)
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
		var errors api2go.HTTPError
		for _, err := range c.Errors ***REMOVED***
			errors.Errors = append(errors.Errors, api2go.Error***REMOVED***
				Title: err.Error(),
			***REMOVED***)
		***REMOVED***
		data, _ := json.Marshal(errors)
		c.Data(c.Writer.Status(), contentType, data)
	***REMOVED***
***REMOVED***
