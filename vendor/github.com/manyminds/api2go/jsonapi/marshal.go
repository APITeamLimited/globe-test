package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// RelationshipType specifies the type of a relationship.
type RelationshipType int

// The available relationship types.
//
// Note: DefaultRelationship guesses the relationship type based on the
// pluralization of the reference name.
const (
	DefaultRelationship RelationshipType = iota
	ToOneRelationship
	ToManyRelationship
)

// The MarshalIdentifier interface is necessary to give an element a unique ID.
//
// Note: The implementation of this interface is mandatory.
type MarshalIdentifier interface ***REMOVED***
	GetID() string
***REMOVED***

// ReferenceID contains all necessary information in order to reference another
// struct in JSON API.
type ReferenceID struct ***REMOVED***
	ID           string
	Type         string
	Name         string
	Relationship RelationshipType
***REMOVED***

// A Reference information about possible references of a struct.
//
// Note: If IsNotLoaded is set to true, the `data` field will be omitted and only
// the `links` object will be generated. You should do this if there are some
// references, but you do not want to load them. Otherwise, if IsNotLoaded is
// false and GetReferencedIDs() returns no IDs for this reference name, an
// empty `data` field will be added which means that there are no references.
type Reference struct ***REMOVED***
	Type         string
	Name         string
	IsNotLoaded  bool
	Relationship RelationshipType
***REMOVED***

// The MarshalReferences interface must be implemented if the struct to be
// serialized has relationships.
type MarshalReferences interface ***REMOVED***
	GetReferences() []Reference
***REMOVED***

// The MarshalLinkedRelations interface must be implemented if there are
// reference ids that should be included in the document.
type MarshalLinkedRelations interface ***REMOVED***
	MarshalReferences
	MarshalIdentifier
	GetReferencedIDs() []ReferenceID
***REMOVED***

// The MarshalIncludedRelations interface must be implemented if referenced
// structs should be included in the document.
type MarshalIncludedRelations interface ***REMOVED***
	MarshalReferences
	MarshalIdentifier
	GetReferencedStructs() []MarshalIdentifier
***REMOVED***

// The MarshalCustomLinks interface can be implemented if the struct should
// want any custom links.
type MarshalCustomLinks interface ***REMOVED***
	MarshalIdentifier
	GetCustomLinks(string) Links
***REMOVED***

// The MarshalCustomRelationshipMeta interface can be implemented if the struct should
// want a custom meta in a relationship.
type MarshalCustomRelationshipMeta interface ***REMOVED***
	MarshalIdentifier
	GetCustomMeta(string) map[string]Meta
***REMOVED***

// A ServerInformation implementor can be passed to MarshalWithURLs to generate
// the `self` and `related` urls inside `links`.
type ServerInformation interface ***REMOVED***
	GetBaseURL() string
	GetPrefix() string
***REMOVED***

// MarshalWithURLs can be used to pass along a ServerInformation implementor.
func MarshalWithURLs(data interface***REMOVED******REMOVED***, information ServerInformation) ([]byte, error) ***REMOVED***
	document, err := MarshalToStruct(data, information)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return json.Marshal(document)
***REMOVED***

// Marshal wraps data in a Document and returns its JSON encoding.
//
// Data can be a struct, a pointer to a struct or a slice of structs. All structs
// must at least implement the `MarshalIdentifier` interface.
func Marshal(data interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	document, err := MarshalToStruct(data, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return json.Marshal(document)
***REMOVED***

// MarshalToStruct marshals an api2go compatible struct into a jsonapi Document
// structure which then can be marshaled to JSON. You only need this method if
// you want to extract or extend parts of the document. You should directly use
// Marshal to get a []byte with JSON in it.
func MarshalToStruct(data interface***REMOVED******REMOVED***, information ServerInformation) (*Document, error) ***REMOVED***
	if data == nil ***REMOVED***
		return &Document***REMOVED******REMOVED***, nil
	***REMOVED***

	switch reflect.TypeOf(data).Kind() ***REMOVED***
	case reflect.Slice:
		return marshalSlice(data, information)
	case reflect.Struct, reflect.Ptr:
		return marshalStruct(data.(MarshalIdentifier), information)
	default:
		return nil, errors.New("Marshal only accepts slice, struct or ptr types")
	***REMOVED***
***REMOVED***

func recursivelyEmbedIncludes(input []MarshalIdentifier) []MarshalIdentifier ***REMOVED***
	var referencedStructs []MarshalIdentifier

	for _, referencedStruct := range input ***REMOVED***
		included, ok := referencedStruct.(MarshalIncludedRelations)
		if ok ***REMOVED***
			referencedStructs = append(referencedStructs, included.GetReferencedStructs()...)
		***REMOVED***
	***REMOVED***

	if len(referencedStructs) == 0 ***REMOVED***
		return input
	***REMOVED***

	childStructs := recursivelyEmbedIncludes(referencedStructs)
	referencedStructs = append(referencedStructs, childStructs...)
	referencedStructs = append(input, referencedStructs...)

	return referencedStructs
***REMOVED***

func marshalSlice(data interface***REMOVED******REMOVED***, information ServerInformation) (*Document, error) ***REMOVED***
	result := &Document***REMOVED******REMOVED***

	val := reflect.ValueOf(data)
	dataElements := make([]Data, val.Len())
	var referencedStructs []MarshalIdentifier

	for i := 0; i < val.Len(); i++ ***REMOVED***
		k := val.Index(i).Interface()
		element, ok := k.(MarshalIdentifier)
		if !ok ***REMOVED***
			return nil, errors.New("all elements within the slice must implement api2go.MarshalIdentifier")
		***REMOVED***

		err := marshalData(element, &dataElements[i], information)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		included, ok := k.(MarshalIncludedRelations)
		if ok ***REMOVED***
			referencedStructs = append(referencedStructs, included.GetReferencedStructs()...)
		***REMOVED***
	***REMOVED***

	allReferencedStructs := recursivelyEmbedIncludes(referencedStructs)
	includedElements, err := filterDuplicates(allReferencedStructs, information)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	result.Data = &DataContainer***REMOVED***
		DataArray: dataElements,
	***REMOVED***

	if includedElements != nil && len(includedElements) > 0 ***REMOVED***
		result.Included = includedElements
	***REMOVED***

	return result, nil
***REMOVED***

func filterDuplicates(input []MarshalIdentifier, information ServerInformation) ([]Data, error) ***REMOVED***
	alreadyIncluded := map[string]map[string]bool***REMOVED******REMOVED***
	includedElements := []Data***REMOVED******REMOVED***

	for _, referencedStruct := range input ***REMOVED***
		structType := getStructType(referencedStruct)

		if alreadyIncluded[structType] == nil ***REMOVED***
			alreadyIncluded[structType] = make(map[string]bool)
		***REMOVED***

		if !alreadyIncluded[structType][referencedStruct.GetID()] ***REMOVED***
			var data Data
			err := marshalData(referencedStruct, &data, information)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			includedElements = append(includedElements, data)
			alreadyIncluded[structType][referencedStruct.GetID()] = true
		***REMOVED***
	***REMOVED***

	return includedElements, nil
***REMOVED***

func marshalData(element MarshalIdentifier, data *Data, information ServerInformation) error ***REMOVED***
	refValue := reflect.ValueOf(element)
	if refValue.Kind() == reflect.Ptr && refValue.IsNil() ***REMOVED***
		return errors.New("MarshalIdentifier must not be nil")
	***REMOVED***

	attributes, err := json.Marshal(element)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	data.Attributes = attributes
	data.ID = element.GetID()
	data.Type = getStructType(element)

	if information != nil ***REMOVED***
		if customLinks, ok := element.(MarshalCustomLinks); ok ***REMOVED***
			if data.Links == nil ***REMOVED***
				data.Links = make(Links)
			***REMOVED***
			base := getLinkBaseURL(element, information)
			for k, v := range customLinks.GetCustomLinks(base) ***REMOVED***
				if _, ok := data.Links[k]; !ok ***REMOVED***
					data.Links[k] = v
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if references, ok := element.(MarshalLinkedRelations); ok ***REMOVED***
		data.Relationships = getStructRelationships(references, information)
	***REMOVED***

	return nil
***REMOVED***

func isToMany(relationshipType RelationshipType, name string) bool ***REMOVED***
	if relationshipType == DefaultRelationship ***REMOVED***
		return Pluralize(name) == name
	***REMOVED***

	return relationshipType == ToManyRelationship
***REMOVED***

func getMetaForRelation(metaSource MarshalCustomRelationshipMeta, name string, information ServerInformation) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	meta := make(map[string]interface***REMOVED******REMOVED***)
	base := getLinkBaseURL(metaSource, information)
	if metaMap, ok := metaSource.GetCustomMeta(base)[name]; ok ***REMOVED***
		for k, v := range metaMap ***REMOVED***
			if _, ok := meta[k]; !ok ***REMOVED***
				meta[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return meta
***REMOVED***

func getStructRelationships(relationer MarshalLinkedRelations, information ServerInformation) map[string]Relationship ***REMOVED***
	referencedIDs := relationer.GetReferencedIDs()
	sortedResults := map[string][]ReferenceID***REMOVED******REMOVED***
	relationships := map[string]Relationship***REMOVED******REMOVED***

	for _, referenceID := range referencedIDs ***REMOVED***
		sortedResults[referenceID.Name] = append(sortedResults[referenceID.Name], referenceID)
	***REMOVED***

	references := relationer.GetReferences()

	// helper map to check if all references are included to also include empty ones
	notIncludedReferences := map[string]Reference***REMOVED******REMOVED***
	for _, reference := range references ***REMOVED***
		notIncludedReferences[reference.Name] = reference
	***REMOVED***

	for name, referenceIDs := range sortedResults ***REMOVED***
		relationships[name] = Relationship***REMOVED******REMOVED***

		// if referenceType is plural, we need to use an array for data, otherwise it's just an object
		container := RelationshipDataContainer***REMOVED******REMOVED***

		if isToMany(referenceIDs[0].Relationship, referenceIDs[0].Name) ***REMOVED***
			// multiple elements in links
			container.DataArray = []RelationshipData***REMOVED******REMOVED***
			for _, referenceID := range referenceIDs ***REMOVED***
				container.DataArray = append(container.DataArray, RelationshipData***REMOVED***
					Type: referenceID.Type,
					ID:   referenceID.ID,
				***REMOVED***)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			container.DataObject = &RelationshipData***REMOVED***
				Type: referenceIDs[0].Type,
				ID:   referenceIDs[0].ID,
			***REMOVED***
		***REMOVED***

		// set URLs if necessary
		links := getLinksForServerInformation(relationer, name, information)

		// get the custom meta for this relationship
		var meta map[string]interface***REMOVED******REMOVED***
		if customMetaSource, ok := relationer.(MarshalCustomRelationshipMeta); ok ***REMOVED***
			meta = getMetaForRelation(customMetaSource, name, information)
		***REMOVED***

		relationship := Relationship***REMOVED***
			Data:  &container,
			Links: links,
			Meta:  meta,
		***REMOVED***

		relationships[name] = relationship

		// this marks the reference as already included
		delete(notIncludedReferences, referenceIDs[0].Name)
	***REMOVED***

	// check for empty references
	for name, reference := range notIncludedReferences ***REMOVED***
		container := RelationshipDataContainer***REMOVED******REMOVED***

		// Plural empty relationships need an empty array and empty to-one need a null in the json
		if !reference.IsNotLoaded && isToMany(reference.Relationship, reference.Name) ***REMOVED***
			container.DataArray = []RelationshipData***REMOVED******REMOVED***
		***REMOVED***

		links := getLinksForServerInformation(relationer, name, information)

		// get the custom meta for this relationship
		var meta map[string]interface***REMOVED******REMOVED***
		if customMetaSource, ok := relationer.(MarshalCustomRelationshipMeta); ok ***REMOVED***
			meta = getMetaForRelation(customMetaSource, name, information)
		***REMOVED***

		relationship := Relationship***REMOVED***
			Links: links,
			Meta:  meta,
		***REMOVED***

		// skip relationship data completely if IsNotLoaded is set
		if !reference.IsNotLoaded ***REMOVED***
			relationship.Data = &container
		***REMOVED***

		relationships[name] = relationship
	***REMOVED***

	return relationships
***REMOVED***

func getLinkBaseURL(element MarshalIdentifier, information ServerInformation) string ***REMOVED***
	prefix := strings.Trim(information.GetBaseURL(), "/")
	namespace := strings.Trim(information.GetPrefix(), "/")
	structType := getStructType(element)

	if namespace != "" ***REMOVED***
		prefix += "/" + namespace
	***REMOVED***

	return fmt.Sprintf("%s/%s/%s", prefix, structType, element.GetID())
***REMOVED***

func getLinksForServerInformation(relationer MarshalLinkedRelations, name string, information ServerInformation) Links ***REMOVED***
	if information == nil ***REMOVED***
		return nil
	***REMOVED***

	links := make(Links)
	base := getLinkBaseURL(relationer, information)

	links["self"] = Link***REMOVED***Href: fmt.Sprintf("%s/relationships/%s", base, name)***REMOVED***
	links["related"] = Link***REMOVED***Href: fmt.Sprintf("%s/%s", base, name)***REMOVED***

	return links
***REMOVED***

func marshalStruct(data MarshalIdentifier, information ServerInformation) (*Document, error) ***REMOVED***
	var contentData Data

	err := marshalData(data, &contentData, information)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	result := &Document***REMOVED***
		Data: &DataContainer***REMOVED***
			DataObject: &contentData,
		***REMOVED***,
	***REMOVED***

	included, ok := data.(MarshalIncludedRelations)
	if ok ***REMOVED***
		included, err := filterDuplicates(recursivelyEmbedIncludes(included.GetReferencedStructs()), information)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if len(included) > 0 ***REMOVED***
			result.Included = included
		***REMOVED***
	***REMOVED***

	return result, nil
***REMOVED***

func getStructType(data interface***REMOVED******REMOVED***) string ***REMOVED***
	entityName, ok := data.(EntityNamer)
	if ok ***REMOVED***
		return entityName.GetName()
	***REMOVED***

	reflectType := reflect.TypeOf(data)
	if reflectType.Kind() == reflect.Ptr ***REMOVED***
		return Pluralize(Jsonify(reflectType.Elem().Name()))
	***REMOVED***

	return Pluralize(Jsonify(reflectType.Name()))
***REMOVED***
