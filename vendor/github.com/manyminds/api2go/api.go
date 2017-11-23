package api2go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/manyminds/api2go/routing"
)

const (
	codeInvalidQueryFields  = "API2GO_INVALID_FIELD_QUERY_PARAM"
	defaultContentTypHeader = "application/vnd.api+json"
)

var (
	queryPageRegex   = regexp.MustCompile(`^page\[(\w+)\]$`)
	queryFieldsRegex = regexp.MustCompile(`^fields\[(\w+)\]$`)
)

type information struct ***REMOVED***
	prefix   string
	resolver URLResolver
***REMOVED***

func (i information) GetBaseURL() string ***REMOVED***
	return i.resolver.GetBaseURL()
***REMOVED***

func (i information) GetPrefix() string ***REMOVED***
	return i.prefix
***REMOVED***

type paginationQueryParams struct ***REMOVED***
	number, size, offset, limit string
***REMOVED***

func newPaginationQueryParams(r *http.Request) paginationQueryParams ***REMOVED***
	var result paginationQueryParams

	queryParams := r.URL.Query()
	result.number = queryParams.Get("page[number]")
	result.size = queryParams.Get("page[size]")
	result.offset = queryParams.Get("page[offset]")
	result.limit = queryParams.Get("page[limit]")

	return result
***REMOVED***

func (p paginationQueryParams) isValid() bool ***REMOVED***
	if p.number == "" && p.size == "" && p.offset == "" && p.limit == "" ***REMOVED***
		return false
	***REMOVED***

	if p.number != "" && p.size != "" && p.offset == "" && p.limit == "" ***REMOVED***
		return true
	***REMOVED***

	if p.number == "" && p.size == "" && p.offset != "" && p.limit != "" ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func (p paginationQueryParams) getLinks(r *http.Request, count uint, info information) (result jsonapi.Links, err error) ***REMOVED***
	result = make(jsonapi.Links)

	params := r.URL.Query()
	prefix := ""
	baseURL := info.GetBaseURL()
	if baseURL != "" ***REMOVED***
		prefix = baseURL
	***REMOVED***
	requestURL := fmt.Sprintf("%s%s", prefix, r.URL.Path)

	if p.number != "" ***REMOVED***
		// we have number & size params
		var number uint64
		number, err = strconv.ParseUint(p.number, 10, 64)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		if p.number != "1" ***REMOVED***
			params.Set("page[number]", "1")
			query, _ := url.QueryUnescape(params.Encode())
			result["first"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***

			params.Set("page[number]", strconv.FormatUint(number-1, 10))
			query, _ = url.QueryUnescape(params.Encode())
			result["prev"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***
		***REMOVED***

		// calculate last page number
		var size uint64
		size, err = strconv.ParseUint(p.size, 10, 64)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		totalPages := (uint64(count) / size)
		if (uint64(count) % size) != 0 ***REMOVED***
			// there is one more page with some len(items) < size
			totalPages++
		***REMOVED***

		if number != totalPages ***REMOVED***
			params.Set("page[number]", strconv.FormatUint(number+1, 10))
			query, _ := url.QueryUnescape(params.Encode())
			result["next"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***

			params.Set("page[number]", strconv.FormatUint(totalPages, 10))
			query, _ = url.QueryUnescape(params.Encode())
			result["last"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// we have offset & limit params
		var offset, limit uint64
		offset, err = strconv.ParseUint(p.offset, 10, 64)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		limit, err = strconv.ParseUint(p.limit, 10, 64)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		if p.offset != "0" ***REMOVED***
			params.Set("page[offset]", "0")
			query, _ := url.QueryUnescape(params.Encode())
			result["first"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***

			var prevOffset uint64
			if limit > offset ***REMOVED***
				prevOffset = 0
			***REMOVED*** else ***REMOVED***
				prevOffset = offset - limit
			***REMOVED***
			params.Set("page[offset]", strconv.FormatUint(prevOffset, 10))
			query, _ = url.QueryUnescape(params.Encode())
			result["prev"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***
		***REMOVED***

		// check if there are more entries to be loaded
		if (offset + limit) < uint64(count) ***REMOVED***
			params.Set("page[offset]", strconv.FormatUint(offset+limit, 10))
			query, _ := url.QueryUnescape(params.Encode())
			result["next"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***

			params.Set("page[offset]", strconv.FormatUint(uint64(count)-limit, 10))
			query, _ = url.QueryUnescape(params.Encode())
			result["last"] = jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", requestURL, query)***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

type notAllowedHandler struct ***REMOVED***
	API *API
***REMOVED***

func (n notAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	err := NewHTTPError(nil, "Method Not Allowed", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)

	contentType := defaultContentTypHeader
	if n.API != nil ***REMOVED***
		contentType = n.API.ContentType
	***REMOVED***

	handleError(err, w, r, contentType)
***REMOVED***

type resource struct ***REMOVED***
	resourceType reflect.Type
	source       interface***REMOVED******REMOVED***
	name         string
	api          *API
***REMOVED***

// middlewareChain executes the middleeware chain setup
func (api *API) middlewareChain(c APIContexter, w http.ResponseWriter, r *http.Request) ***REMOVED***
	for _, middleware := range api.middlewares ***REMOVED***
		middleware(c, w, r)
	***REMOVED***
***REMOVED***

// allocateContext creates a context for the api.contextPool, saving allocations
func (api *API) allocateDefaultContext() APIContexter ***REMOVED***
	return &APIContext***REMOVED******REMOVED***
***REMOVED***

func (api *API) addResource(prototype jsonapi.MarshalIdentifier, source interface***REMOVED******REMOVED***) *resource ***REMOVED***
	resourceType := reflect.TypeOf(prototype)
	if resourceType.Kind() != reflect.Struct && resourceType.Kind() != reflect.Ptr ***REMOVED***
		panic("pass an empty resource struct or a struct pointer to AddResource!")
	***REMOVED***

	var ptrPrototype interface***REMOVED******REMOVED***
	var name string

	if resourceType.Kind() == reflect.Struct ***REMOVED***
		ptrPrototype = reflect.New(resourceType).Interface()
		name = resourceType.Name()
	***REMOVED*** else ***REMOVED***
		ptrPrototype = reflect.ValueOf(prototype).Interface()
		name = resourceType.Elem().Name()
	***REMOVED***

	// check if EntityNamer interface is implemented and use that as name
	entityName, ok := prototype.(jsonapi.EntityNamer)
	if ok ***REMOVED***
		name = entityName.GetName()
	***REMOVED*** else ***REMOVED***
		name = jsonapi.Jsonify(jsonapi.Pluralize(name))
	***REMOVED***

	res := resource***REMOVED***
		resourceType: resourceType,
		name:         name,
		source:       source,
		api:          api,
	***REMOVED***

	requestInfo := func(r *http.Request, api *API) *information ***REMOVED***
		var info *information
		if resolver, ok := api.info.resolver.(RequestAwareURLResolver); ok ***REMOVED***
			resolver.SetRequest(*r)
			info = &information***REMOVED***prefix: api.info.prefix, resolver: resolver***REMOVED***
		***REMOVED*** else ***REMOVED***
			info = &api.info
		***REMOVED***

		return info
	***REMOVED***

	prefix := strings.Trim(api.info.prefix, "/")
	baseURL := "/" + name
	if prefix != "" ***REMOVED***
		baseURL = "/" + prefix + baseURL
	***REMOVED***

	api.router.Handle("OPTIONS", baseURL, func(w http.ResponseWriter, r *http.Request, _ map[string]string) ***REMOVED***
		c := api.contextPool.Get().(APIContexter)
		c.Reset()
		api.middlewareChain(c, w, r)
		w.Header().Set("Allow", strings.Join(getAllowedMethods(source, true), ","))
		w.WriteHeader(http.StatusNoContent)
		api.contextPool.Put(c)
	***REMOVED***)

	api.router.Handle("GET", baseURL, func(w http.ResponseWriter, r *http.Request, _ map[string]string) ***REMOVED***
		info := requestInfo(r, api)
		c := api.contextPool.Get().(APIContexter)
		c.Reset()
		api.middlewareChain(c, w, r)

		err := res.handleIndex(c, w, r, *info)
		api.contextPool.Put(c)
		if err != nil ***REMOVED***
			handleError(err, w, r, api.ContentType)
		***REMOVED***
	***REMOVED***)

	if _, ok := source.(ResourceGetter); ok ***REMOVED***
		api.router.Handle("OPTIONS", baseURL+"/:id", func(w http.ResponseWriter, r *http.Request, _ map[string]string) ***REMOVED***
			c := api.contextPool.Get().(APIContexter)
			c.Reset()
			api.middlewareChain(c, w, r)
			w.Header().Set("Allow", strings.Join(getAllowedMethods(source, false), ","))
			w.WriteHeader(http.StatusNoContent)
			api.contextPool.Put(c)
		***REMOVED***)

		api.router.Handle("GET", baseURL+"/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
			info := requestInfo(r, api)
			c := api.contextPool.Get().(APIContexter)
			c.Reset()
			api.middlewareChain(c, w, r)
			err := res.handleRead(c, w, r, params, *info)
			api.contextPool.Put(c)
			if err != nil ***REMOVED***
				handleError(err, w, r, api.ContentType)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	// generate all routes for linked relations if there are relations
	casted, ok := prototype.(jsonapi.MarshalReferences)
	if ok ***REMOVED***
		relations := casted.GetReferences()
		for _, relation := range relations ***REMOVED***
			api.router.Handle("GET", baseURL+"/:id/relationships/"+relation.Name, func(relation jsonapi.Reference) routing.HandlerFunc ***REMOVED***
				return func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
					info := requestInfo(r, api)
					c := api.contextPool.Get().(APIContexter)
					c.Reset()
					api.middlewareChain(c, w, r)
					err := res.handleReadRelation(c, w, r, params, *info, relation)
					api.contextPool.Put(c)
					if err != nil ***REMOVED***
						handleError(err, w, r, api.ContentType)
					***REMOVED***
				***REMOVED***
			***REMOVED***(relation))

			api.router.Handle("GET", baseURL+"/:id/"+relation.Name, func(relation jsonapi.Reference) routing.HandlerFunc ***REMOVED***
				return func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
					info := requestInfo(r, api)
					c := api.contextPool.Get().(APIContexter)
					c.Reset()
					api.middlewareChain(c, w, r)
					err := res.handleLinked(c, api, w, r, params, relation, *info)
					api.contextPool.Put(c)
					if err != nil ***REMOVED***
						handleError(err, w, r, api.ContentType)
					***REMOVED***
				***REMOVED***
			***REMOVED***(relation))

			api.router.Handle("PATCH", baseURL+"/:id/relationships/"+relation.Name, func(relation jsonapi.Reference) routing.HandlerFunc ***REMOVED***
				return func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
					c := api.contextPool.Get().(APIContexter)
					c.Reset()
					api.middlewareChain(c, w, r)
					err := res.handleReplaceRelation(c, w, r, params, relation)
					api.contextPool.Put(c)
					if err != nil ***REMOVED***
						handleError(err, w, r, api.ContentType)
					***REMOVED***
				***REMOVED***
			***REMOVED***(relation))

			if _, ok := ptrPrototype.(jsonapi.EditToManyRelations); ok && relation.Name == jsonapi.Pluralize(relation.Name) ***REMOVED***
				// generate additional routes to manipulate to-many relationships
				api.router.Handle("POST", baseURL+"/:id/relationships/"+relation.Name, func(relation jsonapi.Reference) routing.HandlerFunc ***REMOVED***
					return func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
						c := api.contextPool.Get().(APIContexter)
						c.Reset()
						api.middlewareChain(c, w, r)
						err := res.handleAddToManyRelation(c, w, r, params, relation)
						api.contextPool.Put(c)
						if err != nil ***REMOVED***
							handleError(err, w, r, api.ContentType)
						***REMOVED***
					***REMOVED***
				***REMOVED***(relation))

				api.router.Handle("DELETE", baseURL+"/:id/relationships/"+relation.Name, func(relation jsonapi.Reference) routing.HandlerFunc ***REMOVED***
					return func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
						c := api.contextPool.Get().(APIContexter)
						c.Reset()
						api.middlewareChain(c, w, r)
						err := res.handleDeleteToManyRelation(c, w, r, params, relation)
						api.contextPool.Put(c)
						if err != nil ***REMOVED***
							handleError(err, w, r, api.ContentType)
						***REMOVED***
					***REMOVED***
				***REMOVED***(relation))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if _, ok := source.(ResourceCreator); ok ***REMOVED***
		api.router.Handle("POST", baseURL, func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
			info := requestInfo(r, api)
			c := api.contextPool.Get().(APIContexter)
			c.Reset()
			api.middlewareChain(c, w, r)
			err := res.handleCreate(c, w, r, info.prefix, *info)
			api.contextPool.Put(c)
			if err != nil ***REMOVED***
				handleError(err, w, r, api.ContentType)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	if _, ok := source.(ResourceDeleter); ok ***REMOVED***
		api.router.Handle("DELETE", baseURL+"/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
			c := api.contextPool.Get().(APIContexter)
			c.Reset()
			api.middlewareChain(c, w, r)
			err := res.handleDelete(c, w, r, params)
			api.contextPool.Put(c)
			if err != nil ***REMOVED***
				handleError(err, w, r, api.ContentType)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	if _, ok := source.(ResourceUpdater); ok ***REMOVED***
		api.router.Handle("PATCH", baseURL+"/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) ***REMOVED***
			info := requestInfo(r, api)
			c := api.contextPool.Get().(APIContexter)
			c.Reset()
			api.middlewareChain(c, w, r)
			err := res.handleUpdate(c, w, r, params, *info)
			api.contextPool.Put(c)
			if err != nil ***REMOVED***
				handleError(err, w, r, api.ContentType)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	api.resources = append(api.resources, res)

	return &res
***REMOVED***

func getAllowedMethods(source interface***REMOVED******REMOVED***, collection bool) []string ***REMOVED***
	result := []string***REMOVED***http.MethodOptions***REMOVED***

	if _, ok := source.(ResourceGetter); ok ***REMOVED***
		result = append(result, http.MethodGet)
	***REMOVED***

	if _, ok := source.(ResourceUpdater); ok ***REMOVED***
		result = append(result, http.MethodPatch)
	***REMOVED***

	if _, ok := source.(ResourceDeleter); ok && !collection ***REMOVED***
		result = append(result, http.MethodDelete)
	***REMOVED***

	if _, ok := source.(ResourceCreator); ok && collection ***REMOVED***
		result = append(result, http.MethodPost)
	***REMOVED***

	return result
***REMOVED***

func buildRequest(c APIContexter, r *http.Request) Request ***REMOVED***
	req := Request***REMOVED***PlainRequest: r***REMOVED***
	params := make(map[string][]string)
	pagination := make(map[string]string)
	for key, values := range r.URL.Query() ***REMOVED***
		params[key] = strings.Split(values[0], ",")
		pageMatches := queryPageRegex.FindStringSubmatch(key)
		if len(pageMatches) > 1 ***REMOVED***
			pagination[pageMatches[1]] = values[0]
		***REMOVED***
	***REMOVED***
	req.Pagination = pagination
	req.QueryParams = params
	req.Header = r.Header
	req.Context = c
	return req
***REMOVED***

func (res *resource) marshalResponse(resp interface***REMOVED******REMOVED***, w http.ResponseWriter, status int, r *http.Request) error ***REMOVED***
	filtered, err := filterSparseFields(resp, r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	result, err := json.Marshal(filtered)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	writeResult(w, result, status, res.api.ContentType)
	return nil
***REMOVED***

func (res *resource) handleIndex(c APIContexter, w http.ResponseWriter, r *http.Request, info information) error ***REMOVED***
	if source, ok := res.source.(PaginatedFindAll); ok ***REMOVED***
		pagination := newPaginationQueryParams(r)

		if pagination.isValid() ***REMOVED***
			count, response, err := source.PaginatedFindAll(buildRequest(c, r))
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			paginationLinks, err := pagination.getLinks(r, count, info)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return res.respondWithPagination(response, info, http.StatusOK, paginationLinks, w, r)
		***REMOVED***
	***REMOVED***

	source, ok := res.source.(FindAll)
	if !ok ***REMOVED***
		return NewHTTPError(nil, "Resource does not implement the FindAll interface", http.StatusNotFound)
	***REMOVED***

	response, err := source.FindAll(buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return res.respondWith(response, info, http.StatusOK, w, r)
***REMOVED***

func (res *resource) handleRead(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, info information) error ***REMOVED***
	source, ok := res.source.(ResourceGetter)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceGetter interface", res.name)
	***REMOVED***

	id := params["id"]

	response, err := source.FindOne(id, buildRequest(c, r))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return res.respondWith(response, info, http.StatusOK, w, r)
***REMOVED***

func (res *resource) handleReadRelation(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, info information, relation jsonapi.Reference) error ***REMOVED***
	source, ok := res.source.(ResourceGetter)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceGetter interface", res.name)
	***REMOVED***

	id := params["id"]

	obj, err := source.FindOne(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	document, err := jsonapi.MarshalToStruct(obj.Result(), info)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	rel, ok := document.Data.DataObject.Relationships[relation.Name]
	if !ok ***REMOVED***
		return NewHTTPError(nil, fmt.Sprintf("There is no relation with the name %s", relation.Name), http.StatusNotFound)
	***REMOVED***

	meta := obj.Metadata()
	if len(meta) > 0 ***REMOVED***
		rel.Meta = meta
	***REMOVED***

	return res.marshalResponse(rel, w, http.StatusOK, r)
***REMOVED***

// try to find the referenced resource and call the findAll Method with referencing resource id as param
func (res *resource) handleLinked(c APIContexter, api *API, w http.ResponseWriter, r *http.Request, params map[string]string, linked jsonapi.Reference, info information) error ***REMOVED***
	id := params["id"]
	for _, resource := range api.resources ***REMOVED***
		if resource.name == linked.Type ***REMOVED***
			request := buildRequest(c, r)
			request.QueryParams[res.name+"ID"] = []string***REMOVED***id***REMOVED***
			request.QueryParams[res.name+"Name"] = []string***REMOVED***linked.Name***REMOVED***

			if source, ok := resource.source.(PaginatedFindAll); ok ***REMOVED***
				// check for pagination, otherwise normal FindAll
				pagination := newPaginationQueryParams(r)
				if pagination.isValid() ***REMOVED***
					var count uint
					count, response, err := source.PaginatedFindAll(request)
					if err != nil ***REMOVED***
						return err
					***REMOVED***

					paginationLinks, err := pagination.getLinks(r, count, info)
					if err != nil ***REMOVED***
						return err
					***REMOVED***

					return res.respondWithPagination(response, info, http.StatusOK, paginationLinks, w, r)
				***REMOVED***
			***REMOVED***

			source, ok := resource.source.(FindAll)
			if !ok ***REMOVED***
				return NewHTTPError(nil, "Resource does not implement the FindAll interface", http.StatusNotFound)
			***REMOVED***

			obj, err := source.FindAll(request)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			return res.respondWith(obj, info, http.StatusOK, w, r)
		***REMOVED***
	***REMOVED***

	return NewHTTPError(
		errors.New("Not Found"),
		"No resource handler is registered to handle the linked resource "+linked.Name,
		http.StatusNotFound,
	)
***REMOVED***

func (res *resource) handleCreate(c APIContexter, w http.ResponseWriter, r *http.Request, prefix string, info information) error ***REMOVED***
	source, ok := res.source.(ResourceCreator)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceCreator interface", res.name)
	***REMOVED***

	ctx, err := unmarshalRequest(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Ok this is weird again, but reflect.New produces a pointer, so we need the pure type without pointer,
	// otherwise we would have a pointer pointer type that we don't want.
	resourceType := res.resourceType
	if resourceType.Kind() == reflect.Ptr ***REMOVED***
		resourceType = resourceType.Elem()
	***REMOVED***
	newObj := reflect.New(resourceType).Interface()

	// Call InitializeObject if available to allow implementers change the object
	// before calling Unmarshal.
	if initSource, ok := source.(ObjectInitializer); ok ***REMOVED***
		initSource.InitializeObject(newObj)
	***REMOVED***

	err = jsonapi.Unmarshal(ctx, newObj)
	if err != nil ***REMOVED***
		return NewHTTPError(nil, err.Error(), http.StatusNotAcceptable)
	***REMOVED***

	var response Responder

	if res.resourceType.Kind() == reflect.Struct ***REMOVED***
		// we have to dereference the pointer if user wants to use non pointer values
		response, err = source.Create(reflect.ValueOf(newObj).Elem().Interface(), buildRequest(c, r))
	***REMOVED*** else ***REMOVED***
		response, err = source.Create(newObj, buildRequest(c, r))
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	result, ok := response.Result().(jsonapi.MarshalIdentifier)

	if !ok ***REMOVED***
		return fmt.Errorf("Expected one newly created object by resource %s", res.name)
	***REMOVED***

	if len(prefix) > 0 ***REMOVED***
		w.Header().Set("Location", "/"+prefix+"/"+res.name+"/"+result.GetID())
	***REMOVED*** else ***REMOVED***
		w.Header().Set("Location", "/"+res.name+"/"+result.GetID())
	***REMOVED***

	// handle 200 status codes
	switch response.StatusCode() ***REMOVED***
	case http.StatusCreated:
		return res.respondWith(response, info, http.StatusCreated, w, r)
	case http.StatusNoContent:
		w.WriteHeader(response.StatusCode())
		return nil
	case http.StatusAccepted:
		w.WriteHeader(response.StatusCode())
		return nil
	default:
		return fmt.Errorf("invalid status code %d from resource %s for method Create", response.StatusCode(), res.name)
	***REMOVED***
***REMOVED***

func (res *resource) handleUpdate(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, info information) error ***REMOVED***
	source, ok := res.source.(ResourceUpdater)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceUpdater interface", res.name)
	***REMOVED***

	id := params["id"]
	obj, err := source.FindOne(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ctx, err := unmarshalRequest(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// we have to make the Result to a pointer to unmarshal into it
	updatingObj := reflect.ValueOf(obj.Result())
	if updatingObj.Kind() == reflect.Struct ***REMOVED***
		updatingObjPtr := reflect.New(reflect.TypeOf(obj.Result()))
		updatingObjPtr.Elem().Set(updatingObj)
		err = jsonapi.Unmarshal(ctx, updatingObjPtr.Interface())
		updatingObj = updatingObjPtr.Elem()
	***REMOVED*** else ***REMOVED***
		err = jsonapi.Unmarshal(ctx, updatingObj.Interface())
	***REMOVED***
	if err != nil ***REMOVED***
		return NewHTTPError(nil, err.Error(), http.StatusNotAcceptable)
	***REMOVED***

	identifiable, ok := updatingObj.Interface().(jsonapi.MarshalIdentifier)
	if !ok || identifiable.GetID() != id ***REMOVED***
		conflictError := errors.New("id in the resource does not match servers endpoint")
		return NewHTTPError(conflictError, conflictError.Error(), http.StatusConflict)
	***REMOVED***

	response, err := source.Update(updatingObj.Interface(), buildRequest(c, r))

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch response.StatusCode() ***REMOVED***
	case http.StatusOK:
		updated := response.Result()
		if updated == nil ***REMOVED***
			internalResponse, err := source.FindOne(id, buildRequest(c, r))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			updated = internalResponse.Result()
			if updated == nil ***REMOVED***
				return fmt.Errorf("Expected FindOne to return one object of resource %s", res.name)
			***REMOVED***

			response = internalResponse
		***REMOVED***

		return res.respondWith(response, info, http.StatusOK, w, r)
	case http.StatusAccepted:
		w.WriteHeader(http.StatusAccepted)
		return nil
	case http.StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
		return nil
	default:
		return fmt.Errorf("invalid status code %d from resource %s for method Update", response.StatusCode(), res.name)
	***REMOVED***
***REMOVED***

func (res *resource) handleReplaceRelation(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, relation jsonapi.Reference) error ***REMOVED***
	source, ok := res.source.(ResourceUpdater)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceUpdater interface", res.name)
	***REMOVED***

	var (
		err     error
		editObj interface***REMOVED******REMOVED***
	)

	id := params["id"]

	response, err := source.FindOne(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	body, err := unmarshalRequest(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	inc := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err = json.Unmarshal(body, &inc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	data, ok := inc["data"]
	if !ok ***REMOVED***
		return errors.New("Invalid object. Need a \"data\" object")
	***REMOVED***

	resType := reflect.TypeOf(response.Result()).Kind()
	if resType == reflect.Struct ***REMOVED***
		editObj = getPointerToStruct(response.Result())
	***REMOVED*** else ***REMOVED***
		editObj = response.Result()
	***REMOVED***

	err = processRelationshipsData(data, relation.Name, editObj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if resType == reflect.Struct ***REMOVED***
		_, err = source.Update(reflect.ValueOf(editObj).Elem().Interface(), buildRequest(c, r))
	***REMOVED*** else ***REMOVED***
		_, err = source.Update(editObj, buildRequest(c, r))
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)
	return err
***REMOVED***

func (res *resource) handleAddToManyRelation(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, relation jsonapi.Reference) error ***REMOVED***
	source, ok := res.source.(ResourceUpdater)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceUpdater interface", res.name)
	***REMOVED***

	var (
		err     error
		editObj interface***REMOVED******REMOVED***
	)

	id := params["id"]

	response, err := source.FindOne(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	body, err := unmarshalRequest(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	inc := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err = json.Unmarshal(body, &inc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	data, ok := inc["data"]
	if !ok ***REMOVED***
		return errors.New("Invalid object. Need a \"data\" object")
	***REMOVED***

	newRels, ok := data.([]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return fmt.Errorf("Data must be an array with \"id\" and \"type\" field to add new to-many relationships")
	***REMOVED***

	newIDs := []string***REMOVED******REMOVED***

	for _, newRel := range newRels ***REMOVED***
		casted, ok := newRel.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return errors.New("entry in data object invalid")
		***REMOVED***
		newID, ok := casted["id"].(string)
		if !ok ***REMOVED***
			return errors.New("no id field found inside data object")
		***REMOVED***

		newIDs = append(newIDs, newID)
	***REMOVED***

	resType := reflect.TypeOf(response.Result()).Kind()
	if resType == reflect.Struct ***REMOVED***
		editObj = getPointerToStruct(response.Result())
	***REMOVED*** else ***REMOVED***
		editObj = response.Result()
	***REMOVED***

	targetObj, ok := editObj.(jsonapi.EditToManyRelations)
	if !ok ***REMOVED***
		return errors.New("target struct must implement jsonapi.EditToManyRelations")
	***REMOVED***
	targetObj.AddToManyIDs(relation.Name, newIDs)

	if resType == reflect.Struct ***REMOVED***
		_, err = source.Update(reflect.ValueOf(targetObj).Elem().Interface(), buildRequest(c, r))
	***REMOVED*** else ***REMOVED***
		_, err = source.Update(targetObj, buildRequest(c, r))
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return err
***REMOVED***

func (res *resource) handleDeleteToManyRelation(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string, relation jsonapi.Reference) error ***REMOVED***
	source, ok := res.source.(ResourceUpdater)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceUpdater interface", res.name)
	***REMOVED***

	var (
		err     error
		editObj interface***REMOVED******REMOVED***
	)

	id := params["id"]

	response, err := source.FindOne(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	body, err := unmarshalRequest(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	inc := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err = json.Unmarshal(body, &inc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	data, ok := inc["data"]
	if !ok ***REMOVED***
		return errors.New("Invalid object. Need a \"data\" object")
	***REMOVED***

	newRels, ok := data.([]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return fmt.Errorf("Data must be an array with \"id\" and \"type\" field to add new to-many relationships")
	***REMOVED***

	obsoleteIDs := []string***REMOVED******REMOVED***

	for _, newRel := range newRels ***REMOVED***
		casted, ok := newRel.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return errors.New("entry in data object invalid")
		***REMOVED***
		obsoleteID, ok := casted["id"].(string)
		if !ok ***REMOVED***
			return errors.New("no id field found inside data object")
		***REMOVED***

		obsoleteIDs = append(obsoleteIDs, obsoleteID)
	***REMOVED***

	resType := reflect.TypeOf(response.Result()).Kind()
	if resType == reflect.Struct ***REMOVED***
		editObj = getPointerToStruct(response.Result())
	***REMOVED*** else ***REMOVED***
		editObj = response.Result()
	***REMOVED***

	targetObj, ok := editObj.(jsonapi.EditToManyRelations)
	if !ok ***REMOVED***
		return errors.New("target struct must implement jsonapi.EditToManyRelations")
	***REMOVED***
	targetObj.DeleteToManyIDs(relation.Name, obsoleteIDs)

	if resType == reflect.Struct ***REMOVED***
		_, err = source.Update(reflect.ValueOf(targetObj).Elem().Interface(), buildRequest(c, r))
	***REMOVED*** else ***REMOVED***
		_, err = source.Update(targetObj, buildRequest(c, r))
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return err
***REMOVED***

// returns a pointer to an interface***REMOVED******REMOVED*** struct
func getPointerToStruct(oldObj interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	resType := reflect.TypeOf(oldObj)
	ptr := reflect.New(resType)
	ptr.Elem().Set(reflect.ValueOf(oldObj))
	return ptr.Interface()
***REMOVED***

func (res *resource) handleDelete(c APIContexter, w http.ResponseWriter, r *http.Request, params map[string]string) error ***REMOVED***
	source, ok := res.source.(ResourceDeleter)

	if !ok ***REMOVED***
		return fmt.Errorf("Resource %s does not implement the ResourceDeleter interface", res.name)
	***REMOVED***

	id := params["id"]
	response, err := source.Delete(id, buildRequest(c, r))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch response.StatusCode() ***REMOVED***
	case http.StatusOK:
		data := map[string]interface***REMOVED******REMOVED******REMOVED***
			"meta": response.Metadata(),
		***REMOVED***

		return res.marshalResponse(data, w, http.StatusOK, r)
	case http.StatusAccepted:
		w.WriteHeader(http.StatusAccepted)
		return nil
	case http.StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
		return nil
	default:
		return fmt.Errorf("invalid status code %d from resource %s for method Delete", response.StatusCode(), res.name)
	***REMOVED***
***REMOVED***

func writeResult(w http.ResponseWriter, data []byte, status int, contentType string) ***REMOVED***
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write(data)
***REMOVED***

func (res *resource) respondWith(obj Responder, info information, status int, w http.ResponseWriter, r *http.Request) error ***REMOVED***
	data, err := jsonapi.MarshalToStruct(obj.Result(), info)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	meta := obj.Metadata()
	if len(meta) > 0 ***REMOVED***
		data.Meta = meta
	***REMOVED***

	if objWithLinks, ok := obj.(LinksResponder); ok ***REMOVED***
		baseURL := strings.Trim(info.GetBaseURL(), "/")
		requestURL := fmt.Sprintf("%s%s", baseURL, r.URL.Path)
		links := objWithLinks.Links(r, requestURL)
		if len(links) > 0 ***REMOVED***
			data.Links = links
		***REMOVED***
	***REMOVED***

	return res.marshalResponse(data, w, status, r)
***REMOVED***

func (res *resource) respondWithPagination(obj Responder, info information, status int, links jsonapi.Links, w http.ResponseWriter, r *http.Request) error ***REMOVED***
	data, err := jsonapi.MarshalToStruct(obj.Result(), info)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	data.Links = links
	meta := obj.Metadata()
	if len(meta) > 0 ***REMOVED***
		data.Meta = meta
	***REMOVED***

	return res.marshalResponse(data, w, status, r)
***REMOVED***

func unmarshalRequest(r *http.Request) ([]byte, error) ***REMOVED***
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return data, nil
***REMOVED***

func filterSparseFields(resp interface***REMOVED******REMOVED***, r *http.Request) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	query := r.URL.Query()
	queryParams := parseQueryFields(&query)
	if len(queryParams) < 1 ***REMOVED***
		return resp, nil
	***REMOVED***

	if document, ok := resp.(*jsonapi.Document); ok ***REMOVED***
		wrongFields := map[string][]string***REMOVED******REMOVED***

		// single entry in data
		data := document.Data.DataObject
		if data != nil ***REMOVED***
			errors := replaceAttributes(&queryParams, data)
			for t, v := range errors ***REMOVED***
				wrongFields[t] = v
			***REMOVED***
		***REMOVED***

		// data can be a slice too
		datas := document.Data.DataArray
		for index, data := range datas ***REMOVED***
			errors := replaceAttributes(&queryParams, &data)
			for t, v := range errors ***REMOVED***
				wrongFields[t] = v
			***REMOVED***
			datas[index] = data
		***REMOVED***

		// included slice
		for index, include := range document.Included ***REMOVED***
			errors := replaceAttributes(&queryParams, &include)
			for t, v := range errors ***REMOVED***
				wrongFields[t] = v
			***REMOVED***
			document.Included[index] = include
		***REMOVED***

		if len(wrongFields) > 0 ***REMOVED***
			httpError := NewHTTPError(nil, "Some requested fields were invalid", http.StatusBadRequest)
			for k, v := range wrongFields ***REMOVED***
				for _, field := range v ***REMOVED***
					httpError.Errors = append(httpError.Errors, Error***REMOVED***
						Status: "Bad Request",
						Code:   codeInvalidQueryFields,
						Title:  fmt.Sprintf(`Field "%s" does not exist for type "%s"`, field, k),
						Detail: "Please make sure you do only request existing fields",
						Source: &ErrorSource***REMOVED***
							Parameter: fmt.Sprintf("fields[%s]", k),
						***REMOVED***,
					***REMOVED***)
				***REMOVED***
			***REMOVED***
			return nil, httpError
		***REMOVED***
	***REMOVED***
	return resp, nil
***REMOVED***

func parseQueryFields(query *url.Values) (result map[string][]string) ***REMOVED***
	result = map[string][]string***REMOVED******REMOVED***
	for name, param := range *query ***REMOVED***
		matches := queryFieldsRegex.FindStringSubmatch(name)
		if len(matches) > 1 ***REMOVED***
			match := matches[1]
			result[match] = strings.Split(param[0], ",")
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func filterAttributes(attributes map[string]interface***REMOVED******REMOVED***, fields []string) (filteredAttributes map[string]interface***REMOVED******REMOVED***, wrongFields []string) ***REMOVED***
	wrongFields = []string***REMOVED******REMOVED***
	filteredAttributes = map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

	for _, field := range fields ***REMOVED***
		if attribute, ok := attributes[field]; ok ***REMOVED***
			filteredAttributes[field] = attribute
		***REMOVED*** else ***REMOVED***
			wrongFields = append(wrongFields, field)
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func replaceAttributes(query *map[string][]string, entry *jsonapi.Data) map[string][]string ***REMOVED***
	fieldType := entry.Type
	attributes := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	_ = json.Unmarshal(entry.Attributes, &attributes)
	fields := (*query)[fieldType]
	if len(fields) > 0 ***REMOVED***
		var wrongFields []string
		attributes, wrongFields = filterAttributes(attributes, fields)
		if len(wrongFields) > 0 ***REMOVED***
			return map[string][]string***REMOVED***
				fieldType: wrongFields,
			***REMOVED***
		***REMOVED***
		bytes, _ := json.Marshal(attributes)
		entry.Attributes = bytes
	***REMOVED***

	return nil
***REMOVED***

func handleError(err error, w http.ResponseWriter, r *http.Request, contentType string) ***REMOVED***
	log.Println(err)
	if e, ok := err.(HTTPError); ok ***REMOVED***
		writeResult(w, []byte(marshalHTTPError(e)), e.status, contentType)
		return

	***REMOVED***

	e := NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	writeResult(w, []byte(marshalHTTPError(e)), http.StatusInternalServerError, contentType)
***REMOVED***

// TODO: this can also be replaced with a struct into that we directly json.Unmarshal
func processRelationshipsData(data interface***REMOVED******REMOVED***, linkName string, target interface***REMOVED******REMOVED***) error ***REMOVED***
	hasOne, ok := data.(map[string]interface***REMOVED******REMOVED***)
	if ok ***REMOVED***
		hasOneID, ok := hasOne["id"].(string)
		if !ok ***REMOVED***
			return fmt.Errorf("data object must have a field id for %s", linkName)
		***REMOVED***

		target, ok := target.(jsonapi.UnmarshalToOneRelations)
		if !ok ***REMOVED***
			return errors.New("target struct must implement interface UnmarshalToOneRelations")
		***REMOVED***

		target.SetToOneReferenceID(linkName, hasOneID)
	***REMOVED*** else if data == nil ***REMOVED***
		// this means that a to-one relationship must be deleted
		target, ok := target.(jsonapi.UnmarshalToOneRelations)
		if !ok ***REMOVED***
			return errors.New("target struct must implement interface UnmarshalToOneRelations")
		***REMOVED***

		target.SetToOneReferenceID(linkName, "")
	***REMOVED*** else ***REMOVED***
		hasMany, ok := data.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return fmt.Errorf("invalid data object or array, must be an object with \"id\" and \"type\" field for %s", linkName)
		***REMOVED***

		target, ok := target.(jsonapi.UnmarshalToManyRelations)
		if !ok ***REMOVED***
			return errors.New("target struct must implement interface UnmarshalToManyRelations")
		***REMOVED***

		hasManyIDs := []string***REMOVED******REMOVED***

		for _, entry := range hasMany ***REMOVED***
			data, ok := entry.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return fmt.Errorf("entry in data array must be an object for %s", linkName)
			***REMOVED***
			dataID, ok := data["id"].(string)
			if !ok ***REMOVED***
				return fmt.Errorf("all data objects must have a field id for %s", linkName)
			***REMOVED***

			hasManyIDs = append(hasManyIDs, dataID)
		***REMOVED***

		target.SetToManyReferenceIDs(linkName, hasManyIDs)
	***REMOVED***

	return nil
***REMOVED***
