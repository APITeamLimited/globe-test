package bpool

// WrapByteSlice wraps a []byte as a ByteSlice
func WrapByteSlice(full []byte, headerLength int) ByteSlice ***REMOVED***
	return ByteSlice***REMOVED***
		full:    full,
		current: full[headerLength:],
		head:    headerLength,
		end:     len(full),
	***REMOVED***
***REMOVED***

// ByteSlice provides a wrapper around []byte with some added convenience
type ByteSlice struct ***REMOVED***
	full    []byte
	current []byte
	head    int
	end     int
***REMOVED***

// ResliceTo reslices the end of the current slice.
func (b ByteSlice) ResliceTo(end int) ByteSlice ***REMOVED***
	return ByteSlice***REMOVED***
		full:    b.full,
		current: b.current[:end],
		head:    b.head,
		end:     b.head + end,
	***REMOVED***
***REMOVED***

// Bytes returns the current slice
func (b ByteSlice) Bytes() []byte ***REMOVED***
	return b.current
***REMOVED***

// BytesWithHeader returns the current slice preceded by the header
func (b ByteSlice) BytesWithHeader() []byte ***REMOVED***
	return b.full[:b.end]
***REMOVED***

// Full returns the full original buffer underlying the ByteSlice
func (b ByteSlice) Full() []byte ***REMOVED***
	return b.full
***REMOVED***

// ByteSlicePool is a bool of byte slices
type ByteSlicePool interface ***REMOVED***
	// Get gets a byte slice from the pool
	GetSlice() ByteSlice
	// Put returns a byte slice to the pool
	PutSlice(ByteSlice)
	// NumPooled returns the number of currently pooled items
	NumPooled() int
***REMOVED***

// NewByteSlicePool creates a new ByteSlicePool bounded to the
// given maxSize, with new byte arrays sized based on width
func NewByteSlicePool(maxSize int, width int) ByteSlicePool ***REMOVED***
	return NewHeaderPreservingByteSlicePool(maxSize, width, 0)
***REMOVED***

// NewHeaderPreservingByteSlicePool creates a new ByteSlicePool bounded to the
// given maxSize, with new byte arrays sized based on width and headerLength
// preserved at the beginning of the slice.
func NewHeaderPreservingByteSlicePool(maxSize int, width int, headerLength int) ByteSlicePool ***REMOVED***
	return &BytePool***REMOVED***
		c: make(chan []byte, maxSize),
		w: width + headerLength,
		h: headerLength,
	***REMOVED***
***REMOVED***

// GetSlice implements the method from interface ByteSlicePool
func (bp *BytePool) GetSlice() ByteSlice ***REMOVED***
	return WrapByteSlice(bp.Get(), bp.h)
***REMOVED***

// PutSlice implements the method from interface ByteSlicePool
func (bp *BytePool) PutSlice(b ByteSlice) ***REMOVED***
	bp.Put(b.Full())
***REMOVED***
