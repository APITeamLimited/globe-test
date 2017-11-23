package highlight

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

// A Group represents a syntax group
type Group uint8

// Groups contains all of the groups that are defined
// You can access them in the map via their string name
var Groups map[string]Group
var numGroups Group

// String returns the group name attached to the specific group
func (g Group) String() string ***REMOVED***
	for k, v := range Groups ***REMOVED***
		if v == g ***REMOVED***
			return k
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// A Def is a full syntax definition for a language
// It has a filetype, information about how to detect the filetype based
// on filename or header (the first line of the file)
// Then it has the rules which define how to highlight the file
type Def struct ***REMOVED***
	FileType string
	ftdetect []*regexp.Regexp
	rules    *rules
***REMOVED***

// A Pattern is one simple syntax rule
// It has a group that the rule belongs to, as well as
// the regular expression to match the pattern
type pattern struct ***REMOVED***
	group Group
	regex *regexp.Regexp
***REMOVED***

// rules defines which patterns and regions can be used to highlight
// a filetype
type rules struct ***REMOVED***
	regions  []*region
	patterns []*pattern
	includes []string
***REMOVED***

// A region is a highlighted region (such as a multiline comment, or a string)
// It belongs to a group, and has start and end regular expressions
// A region also has rules of its own that only apply when matching inside the
// region and also rules from the above region do not match inside this region
// Note that a region may contain more regions
type region struct ***REMOVED***
	group      Group
	limitGroup Group
	parent     *region
	start      *regexp.Regexp
	end        *regexp.Regexp
	skip       *regexp.Regexp
	rules      *rules
***REMOVED***

func init() ***REMOVED***
	Groups = make(map[string]Group)
***REMOVED***

// ParseDef parses an input syntax file into a highlight Def
func ParseDef(input []byte) (s *Def, err error) ***REMOVED***
	// This is just so if we have an error, we can exit cleanly and return the parse error to the user
	defer func() ***REMOVED***
		if e := recover(); e != nil ***REMOVED***
			err = e.(error)
		***REMOVED***
	***REMOVED***()

	var rules map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	if err = yaml.Unmarshal(input, &rules); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	s = new(Def)

	for k, v := range rules ***REMOVED***
		if k == "filetype" ***REMOVED***
			filetype := v.(string)

			s.FileType = filetype
		***REMOVED*** else if k == "detect" ***REMOVED***
			ftdetect := v.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
			if len(ftdetect) >= 1 ***REMOVED***
				syntax, err := regexp.Compile(ftdetect["filename"].(string))
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				s.ftdetect = append(s.ftdetect, syntax)
			***REMOVED***
			if len(ftdetect) >= 2 ***REMOVED***
				header, err := regexp.Compile(ftdetect["header"].(string))
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				s.ftdetect = append(s.ftdetect, header)
			***REMOVED***
		***REMOVED*** else if k == "rules" ***REMOVED***
			inputRules := v.([]interface***REMOVED******REMOVED***)

			rules, err := parseRules(inputRules, nil)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			s.rules = rules
		***REMOVED***
	***REMOVED***

	return s, err
***REMOVED***

// ResolveIncludes will sort out the rules for including other filetypes
// You should call this after parsing all the Defs
func ResolveIncludes(defs []*Def) ***REMOVED***
	for _, d := range defs ***REMOVED***
		resolveIncludesInDef(defs, d)
	***REMOVED***
***REMOVED***

func resolveIncludesInDef(defs []*Def, d *Def) ***REMOVED***
	for _, lang := range d.rules.includes ***REMOVED***
		for _, searchDef := range defs ***REMOVED***
			if lang == searchDef.FileType ***REMOVED***
				d.rules.patterns = append(d.rules.patterns, searchDef.rules.patterns...)
				d.rules.regions = append(d.rules.regions, searchDef.rules.regions...)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, r := range d.rules.regions ***REMOVED***
		resolveIncludesInRegion(defs, r)
		r.parent = nil
	***REMOVED***
***REMOVED***

func resolveIncludesInRegion(defs []*Def, region *region) ***REMOVED***
	for _, lang := range region.rules.includes ***REMOVED***
		for _, searchDef := range defs ***REMOVED***
			if lang == searchDef.FileType ***REMOVED***
				region.rules.patterns = append(region.rules.patterns, searchDef.rules.patterns...)
				region.rules.regions = append(region.rules.regions, searchDef.rules.regions...)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, r := range region.rules.regions ***REMOVED***
		resolveIncludesInRegion(defs, r)
		r.parent = region
	***REMOVED***
***REMOVED***

func parseRules(input []interface***REMOVED******REMOVED***, curRegion *region) (*rules, error) ***REMOVED***
	rules := new(rules)

	for _, v := range input ***REMOVED***
		rule := v.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		for k, val := range rule ***REMOVED***
			group := k

			switch object := val.(type) ***REMOVED***
			case string:
				if k == "include" ***REMOVED***
					rules.includes = append(rules.includes, object)
				***REMOVED*** else ***REMOVED***
					// Pattern
					r, err := regexp.Compile(object)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					groupStr := group.(string)
					if _, ok := Groups[groupStr]; !ok ***REMOVED***
						numGroups++
						Groups[groupStr] = numGroups
					***REMOVED***
					groupNum := Groups[groupStr]
					rules.patterns = append(rules.patterns, &pattern***REMOVED***groupNum, r***REMOVED***)
				***REMOVED***
			case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
				// region
				region, err := parseRegion(group.(string), object, curRegion)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				rules.regions = append(rules.regions, region)
			default:
				return nil, fmt.Errorf("Bad type %T", object)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return rules, nil
***REMOVED***

func parseRegion(group string, regionInfo map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, prevRegion *region) (*region, error) ***REMOVED***
	var err error

	region := new(region)
	if _, ok := Groups[group]; !ok ***REMOVED***
		numGroups++
		Groups[group] = numGroups
	***REMOVED***
	groupNum := Groups[group]
	region.group = groupNum
	region.parent = prevRegion

	region.start, err = regexp.Compile(regionInfo["start"].(string))

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	region.end, err = regexp.Compile(regionInfo["end"].(string))

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// skip is optional
	if _, ok := regionInfo["skip"]; ok ***REMOVED***
		region.skip, err = regexp.Compile(regionInfo["skip"].(string))

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// limit-color is optional
	if _, ok := regionInfo["limit-group"]; ok ***REMOVED***
		groupStr := regionInfo["limit-group"].(string)
		if _, ok := Groups[groupStr]; !ok ***REMOVED***
			numGroups++
			Groups[groupStr] = numGroups
		***REMOVED***
		groupNum := Groups[groupStr]
		region.limitGroup = groupNum

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		region.limitGroup = region.group
	***REMOVED***

	region.rules, err = parseRules(regionInfo["rules"].([]interface***REMOVED******REMOVED***), region)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return region, nil
***REMOVED***
