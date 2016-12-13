/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package v1

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/loadimpact/k6/api/common"
	"github.com/loadimpact/k6/lib"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
	"io/ioutil"
	"net/http"
	"strconv"
)

var contentType = "application/vnd.api+json"

func NewHandler() http.Handler ***REMOVED***
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(jsonErrorsMiddleware)

	v1 := router.Group("/v1")
	***REMOVED***
		v1.GET("/info", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(lib.Info***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/error", func(c *gin.Context) ***REMOVED***
			_ = c.AbortWithError(500, errors.New("This is an error"))
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			data, err := jsonapi.Marshal(engine.Status)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.PATCH("/status", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)

			var status lib.Status
			data, _ := ioutil.ReadAll(c.Request.Body)
			if err := jsonapi.Unmarshal(data, &status); err != nil ***REMOVED***
				_ = c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			if status.VUsMax.Valid ***REMOVED***
				if status.VUsMax.Int64 < engine.Status.VUs.Int64 ***REMOVED***
					if status.VUsMax.Int64 >= status.VUs.Int64 ***REMOVED***
						if err := engine.SetVUs(status.VUs.Int64); err != nil ***REMOVED***
							_ = c.AbortWithError(http.StatusBadRequest, err)
							return
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						_ = c.AbortWithError(http.StatusBadRequest, lib.ErrMaxTooLow)
						return
					***REMOVED***
				***REMOVED***

				if err := engine.SetMaxVUs(status.VUsMax.Int64); err != nil ***REMOVED***
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
			***REMOVED***
			if status.VUs.Valid ***REMOVED***
				if status.VUs.Int64 > engine.Status.VUsMax.Int64 ***REMOVED***
					_ = c.AbortWithError(http.StatusBadRequest, lib.ErrTooManyVUs)
					return
				***REMOVED***

				if err := engine.SetVUs(status.VUs.Int64); err != nil ***REMOVED***
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
			***REMOVED***
			if status.Running.Valid ***REMOVED***
				engine.SetRunning(status.Running.Bool)
			***REMOVED***

			data, err := jsonapi.Marshal(engine.Status)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			metrics := make([]interface***REMOVED******REMOVED***, 0, len(engine.Metrics))
			for metric, sink := range engine.Metrics ***REMOVED***
				metric.Sample = sink.Format()
				metrics = append(metrics, metric)
			***REMOVED***
			data, err := jsonapi.Marshal(metrics)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics/:id", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			id := c.Param("id")
			for metric, sink := range engine.Metrics ***REMOVED***
				if metric.Name != id ***REMOVED***
					continue
				***REMOVED***
				metric.Sample = sink.Format()
				data, err := jsonapi.Marshal(metric)
				if err != nil ***REMOVED***
					_ = c.AbortWithError(500, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			_ = c.AbortWithError(404, errors.New("Metric not found"))
		***REMOVED***)
		v1.GET("/groups", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			data, err := jsonapi.Marshal(engine.Runner.GetGroups())
			if err != nil ***REMOVED***
				_ = c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/groups/:id", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			id, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			for _, group := range engine.Runner.GetGroups() ***REMOVED***
				if group.ID != id ***REMOVED***
					continue
				***REMOVED***

				data, err := jsonapi.Marshal(group)
				if err != nil ***REMOVED***
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			_ = c.AbortWithError(404, errors.New("Group not found"))
		***REMOVED***)
		v1.GET("/checks", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			data, err := jsonapi.Marshal(engine.Runner.GetChecks())
			if err != nil ***REMOVED***
				_ = c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/checks/:id", func(c *gin.Context) ***REMOVED***
			engine := common.GetEngine(c)
			id, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil ***REMOVED***
				_ = c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			for _, check := range engine.Runner.GetChecks() ***REMOVED***
				if check.ID != id ***REMOVED***
					continue
				***REMOVED***

				data, err := jsonapi.Marshal(check)
				if err != nil ***REMOVED***
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			_ = c.AbortWithError(404, errors.New("Group not found"))
		***REMOVED***)
	***REMOVED***

	return router
***REMOVED***

func jsonErrorsMiddleware(c *gin.Context) ***REMOVED***
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
