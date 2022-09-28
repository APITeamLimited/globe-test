package ast

// UnknownPos is a placeholder position when only the source file
// name is known.
func UnknownPos(filename string) *SourcePos ***REMOVED***
	return &SourcePos***REMOVED***Filename: filename***REMOVED***
***REMOVED***

// NoSourceNode is a placeholder AST node that implements numerous
// interfaces in this package. It can be used to represent an AST
// element for a file whose source is not available.
type NoSourceNode struct ***REMOVED***
	pos *SourcePos
***REMOVED***

// NewNoSourceNode creates a new NoSourceNode for the given filename.
func NewNoSourceNode(filename string) NoSourceNode ***REMOVED***
	return NoSourceNode***REMOVED***pos: UnknownPos(filename)***REMOVED***
***REMOVED***

func (n NoSourceNode) Start() *SourcePos ***REMOVED***
	return n.pos
***REMOVED***

func (n NoSourceNode) End() *SourcePos ***REMOVED***
	return n.pos
***REMOVED***

func (n NoSourceNode) LeadingComments() []Comment ***REMOVED***
	return nil
***REMOVED***

func (n NoSourceNode) TrailingComments() []Comment ***REMOVED***
	return nil
***REMOVED***

func (n NoSourceNode) GetSyntax() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetName() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetValue() ValueNode ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) FieldLabel() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) FieldName() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) FieldType() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) FieldTag() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) FieldExtendee() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetGroupKeyword() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetOptions() *CompactOptionsNode ***REMOVED***
	return nil
***REMOVED***

func (n NoSourceNode) RangeStart() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) RangeEnd() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetNumber() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) MessageName() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetInputType() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) GetOutputType() Node ***REMOVED***
	return n
***REMOVED***

func (n NoSourceNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***
