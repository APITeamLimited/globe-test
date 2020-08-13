package parser

import (
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/unistring"
)

type _scope struct ***REMOVED***
	outer           *_scope
	allowIn         bool
	inIteration     bool
	inSwitch        bool
	inFunction      bool
	declarationList []ast.Declaration

	labels []unistring.String
***REMOVED***

func (self *_parser) openScope() ***REMOVED***
	self.scope = &_scope***REMOVED***
		outer:   self.scope,
		allowIn: true,
	***REMOVED***
***REMOVED***

func (self *_parser) closeScope() ***REMOVED***
	self.scope = self.scope.outer
***REMOVED***

func (self *_scope) declare(declaration ast.Declaration) ***REMOVED***
	self.declarationList = append(self.declarationList, declaration)
***REMOVED***

func (self *_scope) hasLabel(name unistring.String) bool ***REMOVED***
	for _, label := range self.labels ***REMOVED***
		if label == name ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	if self.outer != nil && !self.inFunction ***REMOVED***
		// Crossing a function boundary to look for a label is verboten
		return self.outer.hasLabel(name)
	***REMOVED***
	return false
***REMOVED***
