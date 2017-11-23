package rice

// LocateMethod defines how a box is located.
type LocateMethod int

const (
	LocateFS               = LocateMethod(iota) // Locate on the filesystem according to package path.
	LocateAppended                              // Locate boxes appended to the executable.
	LocateEmbedded                              // Locate embedded boxes.
	LocateWorkingDirectory                      // Locate on the binary working directory
)

// Config allows customizing the box lookup behavior.
type Config struct ***REMOVED***
	// LocateOrder defines the priority order that boxes are searched for. By
	// default, the package global FindBox searches for embedded boxes first,
	// then appended boxes, and then finally boxes on the filesystem.  That
	// search order may be customized by provided the ordered list here. Leaving
	// out a particular method will omit that from the search space. For
	// example, []LocateMethod***REMOVED***LocateEmbedded, LocateAppended***REMOVED*** will never search
	// the filesystem for boxes.
	LocateOrder []LocateMethod
***REMOVED***

// FindBox searches for boxes using the LocateOrder of the config.
func (c *Config) FindBox(boxName string) (*Box, error) ***REMOVED***
	return findBox(boxName, c.LocateOrder)
***REMOVED***

// MustFindBox searches for boxes using the LocateOrder of the config, like
// FindBox does.  It does not return an error, instead it panics when an error
// occurs.
func (c *Config) MustFindBox(boxName string) *Box ***REMOVED***
	box, err := findBox(boxName, c.LocateOrder)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return box
***REMOVED***