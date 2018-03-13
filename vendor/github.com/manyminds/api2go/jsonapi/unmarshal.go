package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// The UnmarshalIdentifier interface must be implemented to set the ID during
// unmarshalling.
type UnmarshalIdentifier interface ***REMOVED***
	SetID(string) error
***REMOVED***

// The UnmarshalToOneRelations interface must be implemented to unmarshal
// to-one relations.
type UnmarshalToOneRelations interface ***REMOVED***
	SetToOneReferenceID(name, ID string) error
***REMOVED***

// The UnmarshalToManyRelations interface must be implemented to unmarshal
// to-many relations.
type UnmarshalToManyRelations interface ***REMOVED***
	SetToManyReferenceIDs(name string, IDs []string) error
***REMOVED***

// The EditToManyRelations interface can be optionally implemented to add and
// delete to-many relationships on a already unmarshalled struct. These methods
// are used by our API for the to-many relationship update routes.
//
// There are 3 HTTP Methods to edit to-many relations:
//
//	PATCH /v1/posts/1/comments
//	Content-Type: application/vnd.api+json
//	Accept: application/vnd.api+json
//
//	***REMOVED***
//	  "data": [
//		***REMOVED*** "type": "comments", "id": "2" ***REMOVED***,
//		***REMOVED*** "type": "comments", "id": "3" ***REMOVED***
//	  ]
//	***REMOVED***
//
// This replaces all of the comments that belong to post with ID 1 and the
// SetToManyReferenceIDs method will be called.
//
//	POST /v1/posts/1/comments
//	Content-Type: application/vnd.api+json
//	Accept: application/vnd.api+json
//
//	***REMOVED***
//	  "data": [
//		***REMOVED*** "type": "comments", "id": "123" ***REMOVED***
//	  ]
//	***REMOVED***
//
// Adds a new comment to the post with ID 1.
// The AddToManyIDs method will be called.
//
//	DELETE /v1/posts/1/comments
//	Content-Type: application/vnd.api+json
//	Accept: application/vnd.api+json
//
//	***REMOVED***
//	  "data": [
//		***REMOVED*** "type": "comments", "id": "12" ***REMOVED***,
//		***REMOVED*** "type": "comments", "id": "13" ***REMOVED***
//	  ]
//	***REMOVED***
//
// Deletes comments that belong to post with ID 1.
// The DeleteToManyIDs method will be called.
type EditToManyRelations interface ***REMOVED***
	AddToManyIDs(name string, IDs []string) error
	DeleteToManyIDs(name string, IDs []string) error
***REMOVED***

// Unmarshal parses a JSON API compatible JSON and populates the target which
// must implement the `UnmarshalIdentifier` interface.
func Unmarshal(data []byte, target interface***REMOVED******REMOVED***) error ***REMOVED***
	if target == nil ***REMOVED***
		return errors.New("target must not be nil")
	***REMOVED***

	if reflect.TypeOf(target).Kind() != reflect.Ptr ***REMOVED***
		return errors.New("target must be a ptr")
	***REMOVED***

	ctx := &Document***REMOVED******REMOVED***

	err := json.Unmarshal(data, ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if ctx.Data == nil ***REMOVED***
		return errors.New(`Source JSON is empty and has no "attributes" payload object`)
	***REMOVED***

	if ctx.Data.DataObject != nil ***REMOVED***
		return setDataIntoTarget(ctx.Data.DataObject, target)
	***REMOVED***

	if ctx.Data.DataArray != nil ***REMOVED***
		targetSlice := reflect.TypeOf(target).Elem()
		if targetSlice.Kind() != reflect.Slice ***REMOVED***
			return fmt.Errorf("Cannot unmarshal array to struct target %s", targetSlice)
		***REMOVED***
		targetType := targetSlice.Elem()
		targetPointer := reflect.ValueOf(target)
		targetValue := targetPointer.Elem()

		for _, record := range ctx.Data.DataArray ***REMOVED***
			// check if there already is an entry with the same id in target slice,
			// otherwise create a new target and append
			var targetRecord, emptyValue reflect.Value
			for i := 0; i < targetValue.Len(); i++ ***REMOVED***
				marshalCasted, ok := targetValue.Index(i).Interface().(MarshalIdentifier)
				if !ok ***REMOVED***
					return errors.New("existing structs must implement interface MarshalIdentifier")
				***REMOVED***
				if record.ID == marshalCasted.GetID() ***REMOVED***
					targetRecord = targetValue.Index(i).Addr()
					break
				***REMOVED***
			***REMOVED***

			if targetRecord == emptyValue || targetRecord.IsNil() ***REMOVED***
				targetRecord = reflect.New(targetType)
				err := setDataIntoTarget(&record, targetRecord.Interface())
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				targetValue = reflect.Append(targetValue, targetRecord.Elem())
			***REMOVED*** else ***REMOVED***
				err := setDataIntoTarget(&record, targetRecord.Interface())
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		targetPointer.Elem().Set(targetValue)
	***REMOVED***

	return nil
***REMOVED***

func setDataIntoTarget(data *Data, target interface***REMOVED******REMOVED***) error ***REMOVED***
	castedTarget, ok := target.(UnmarshalIdentifier)
	if !ok ***REMOVED***
		return errors.New("target must implement UnmarshalIdentifier interface")
	***REMOVED***

	if data.Type == "" ***REMOVED***
		return errors.New("invalid record, no type was specified")
	***REMOVED***

	err := checkType(data.Type, castedTarget)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if data.Attributes != nil ***REMOVED***
		err = json.Unmarshal(data.Attributes, castedTarget)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := castedTarget.SetID(data.ID); err != nil ***REMOVED***
		return err
	***REMOVED***

	return setRelationshipIDs(data.Relationships, castedTarget)
***REMOVED***

// extracts all found relationships and set's them via SetToOneReferenceID or
// SetToManyReferenceIDs
func setRelationshipIDs(relationships map[string]Relationship, target UnmarshalIdentifier) error ***REMOVED***
	for name, rel := range relationships ***REMOVED***
		// if Data is nil, it means that we have an empty toOne relationship
		if rel.Data == nil ***REMOVED***
			castedToOne, ok := target.(UnmarshalToOneRelations)
			if !ok ***REMOVED***
				return fmt.Errorf("struct %s does not implement UnmarshalToOneRelations", reflect.TypeOf(target))
			***REMOVED***

			castedToOne.SetToOneReferenceID(name, "")
			continue
		***REMOVED***

		// valid toOne case
		if rel.Data.DataObject != nil ***REMOVED***
			castedToOne, ok := target.(UnmarshalToOneRelations)
			if !ok ***REMOVED***
				return fmt.Errorf("struct %s does not implement UnmarshalToOneRelations", reflect.TypeOf(target))
			***REMOVED***
			err := castedToOne.SetToOneReferenceID(name, rel.Data.DataObject.ID)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// valid toMany case
		if rel.Data.DataArray != nil ***REMOVED***
			castedToMany, ok := target.(UnmarshalToManyRelations)
			if !ok ***REMOVED***
				return fmt.Errorf("struct %s does not implement UnmarshalToManyRelations", reflect.TypeOf(target))
			***REMOVED***
			IDs := make([]string, len(rel.Data.DataArray))
			for index, relData := range rel.Data.DataArray ***REMOVED***
				IDs[index] = relData.ID
			***REMOVED***
			err := castedToMany.SetToManyReferenceIDs(name, IDs)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func checkType(incomingType string, target UnmarshalIdentifier) error ***REMOVED***
	actualType := getStructType(target)
	if incomingType != actualType ***REMOVED***
		return fmt.Errorf("Type %s in JSON does not match target struct type %s", incomingType, actualType)
	***REMOVED***

	return nil
***REMOVED***
