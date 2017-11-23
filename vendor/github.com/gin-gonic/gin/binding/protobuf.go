// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/golang/protobuf/proto"

	"io/ioutil"
	"net/http"
)

type protobufBinding struct***REMOVED******REMOVED***

func (protobufBinding) Name() string ***REMOVED***
	return "protobuf"
***REMOVED***

func (protobufBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = proto.Unmarshal(buf, obj.(proto.Message)); err != nil ***REMOVED***
		return err
	***REMOVED***

	//Here it's same to return validate(obj), but util now we cann't add `binding:""` to the struct
	//which automatically generate by gen-proto
	return nil
	//return validate(obj)
***REMOVED***
