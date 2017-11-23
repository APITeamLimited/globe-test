package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

var objectSuffix = []byte("***REMOVED***")
var arraySuffix = []byte("[")
var stringSuffix = []byte(`"`)

// A Document represents a JSON API document as specified here: http://jsonapi.org.
type Document struct ***REMOVED***
	Links    Links                  `json:"links,omitempty"`
	Data     *DataContainer         `json:"data"`
	Included []Data                 `json:"included,omitempty"`
	Meta     map[string]interface***REMOVED******REMOVED*** `json:"meta,omitempty"`
***REMOVED***

// A DataContainer is used to marshal and unmarshal single objects and arrays
// of objects.
type DataContainer struct ***REMOVED***
	DataObject *Data
	DataArray  []Data
***REMOVED***

// UnmarshalJSON unmarshals the JSON-encoded data to the DataObject field if the
// root element is an object or to the DataArray field for arrays.
func (c *DataContainer) UnmarshalJSON(payload []byte) error ***REMOVED***
	if bytes.HasPrefix(payload, objectSuffix) ***REMOVED***
		return json.Unmarshal(payload, &c.DataObject)
	***REMOVED***

	if bytes.HasPrefix(payload, arraySuffix) ***REMOVED***
		return json.Unmarshal(payload, &c.DataArray)
	***REMOVED***

	return errors.New("expected a JSON encoded object or array")
***REMOVED***

// MarshalJSON returns the JSON encoding of the DataArray field or the DataObject
// field. It will return "null" if neither of them is set.
func (c *DataContainer) MarshalJSON() ([]byte, error) ***REMOVED***
	if c.DataArray != nil ***REMOVED***
		return json.Marshal(c.DataArray)
	***REMOVED***

	return json.Marshal(c.DataObject)
***REMOVED***

// Link represents a link for return in the document.
type Link struct ***REMOVED***
	Href string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
***REMOVED***

// UnmarshalJSON marshals a string value into the Href field or marshals an
// object value into the whole struct.
func (l *Link) UnmarshalJSON(payload []byte) error ***REMOVED***
	if bytes.HasPrefix(payload, stringSuffix) ***REMOVED***
		return json.Unmarshal(payload, &l.Href)
	***REMOVED***

	if bytes.HasPrefix(payload, objectSuffix) ***REMOVED***
		obj := make(map[string]interface***REMOVED******REMOVED***)
		err := json.Unmarshal(payload, &obj)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var ok bool
		l.Href, ok = obj["href"].(string)
		if !ok ***REMOVED***
			return errors.New(`link object expects a "href" key`)
		***REMOVED***

		l.Meta, _ = obj["meta"].(map[string]interface***REMOVED******REMOVED***)
		return nil
	***REMOVED***

	return errors.New("expected a JSON encoded string or object")
***REMOVED***

// MarshalJSON returns the JSON encoding of only the Href field if the Meta
// field is empty, otherwise it marshals the whole struct.
func (l Link) MarshalJSON() ([]byte, error) ***REMOVED***
	if len(l.Meta) == 0 ***REMOVED***
		return json.Marshal(l.Href)
	***REMOVED***
	return json.Marshal(map[string]interface***REMOVED******REMOVED******REMOVED***
		"href": l.Href,
		"meta": l.Meta,
	***REMOVED***)
***REMOVED***

// Links contains a map of custom Link objects as given by an element.
type Links map[string]Link

// Meta contains unstructured metadata
type Meta map[string]interface***REMOVED******REMOVED***

// Data is a general struct for document data and included data.
type Data struct ***REMOVED***
	Type          string                  `json:"type"`
	ID            string                  `json:"id"`
	Attributes    json.RawMessage         `json:"attributes"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         Links                   `json:"links,omitempty"`
***REMOVED***

// Relationship contains reference IDs to the related structs
type Relationship struct ***REMOVED***
	Links Links                      `json:"links,omitempty"`
	Data  *RelationshipDataContainer `json:"data,omitempty"`
	Meta  map[string]interface***REMOVED******REMOVED***     `json:"meta,omitempty"`
***REMOVED***

// A RelationshipDataContainer is used to marshal and unmarshal single relationship
// objects and arrays of relationship objects.
type RelationshipDataContainer struct ***REMOVED***
	DataObject *RelationshipData
	DataArray  []RelationshipData
***REMOVED***

// UnmarshalJSON unmarshals the JSON-encoded data to the DataObject field if the
// root element is an object or to the DataArray field for arrays.
func (c *RelationshipDataContainer) UnmarshalJSON(payload []byte) error ***REMOVED***
	if bytes.HasPrefix(payload, objectSuffix) ***REMOVED***
		// payload is an object
		return json.Unmarshal(payload, &c.DataObject)
	***REMOVED***

	if bytes.HasPrefix(payload, arraySuffix) ***REMOVED***
		// payload is an array
		return json.Unmarshal(payload, &c.DataArray)
	***REMOVED***

	return errors.New("Invalid json for relationship data array/object")
***REMOVED***

// MarshalJSON returns the JSON encoding of the DataArray field or the DataObject
// field. It will return "null" if neither of them is set.
func (c *RelationshipDataContainer) MarshalJSON() ([]byte, error) ***REMOVED***
	if c.DataArray != nil ***REMOVED***
		return json.Marshal(c.DataArray)
	***REMOVED***
	return json.Marshal(c.DataObject)
***REMOVED***

// RelationshipData represents one specific reference ID.
type RelationshipData struct ***REMOVED***
	Type string `json:"type"`
	ID   string `json:"id"`
***REMOVED***
