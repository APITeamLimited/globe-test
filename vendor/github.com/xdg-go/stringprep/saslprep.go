package stringprep

var mapNonASCIISpaceToASCIISpace = Mapping***REMOVED***
	0x00A0: []rune***REMOVED***0x0020***REMOVED***,
	0x1680: []rune***REMOVED***0x0020***REMOVED***,
	0x2000: []rune***REMOVED***0x0020***REMOVED***,
	0x2001: []rune***REMOVED***0x0020***REMOVED***,
	0x2002: []rune***REMOVED***0x0020***REMOVED***,
	0x2003: []rune***REMOVED***0x0020***REMOVED***,
	0x2004: []rune***REMOVED***0x0020***REMOVED***,
	0x2005: []rune***REMOVED***0x0020***REMOVED***,
	0x2006: []rune***REMOVED***0x0020***REMOVED***,
	0x2007: []rune***REMOVED***0x0020***REMOVED***,
	0x2008: []rune***REMOVED***0x0020***REMOVED***,
	0x2009: []rune***REMOVED***0x0020***REMOVED***,
	0x200A: []rune***REMOVED***0x0020***REMOVED***,
	0x200B: []rune***REMOVED***0x0020***REMOVED***,
	0x202F: []rune***REMOVED***0x0020***REMOVED***,
	0x205F: []rune***REMOVED***0x0020***REMOVED***,
	0x3000: []rune***REMOVED***0x0020***REMOVED***,
***REMOVED***

// SASLprep is a pre-defined stringprep profile for user names and passwords
// as described in RFC-4013.
//
// Because the stringprep distinction between query and stored strings was
// intended for compatibility across profile versions, but SASLprep was never
// updated and is now deprecated, this profile only operates in stored
// strings mode, prohibiting unassigned code points.
var SASLprep Profile = saslprep

var saslprep = Profile***REMOVED***
	Mappings: []Mapping***REMOVED***
		TableB1,
		mapNonASCIISpaceToASCIISpace,
	***REMOVED***,
	Normalize: true,
	Prohibits: []Set***REMOVED***
		TableA1,
		TableC1_2,
		TableC2_1,
		TableC2_2,
		TableC3,
		TableC4,
		TableC5,
		TableC6,
		TableC7,
		TableC8,
		TableC9,
	***REMOVED***,
	CheckBiDi: true,
***REMOVED***
