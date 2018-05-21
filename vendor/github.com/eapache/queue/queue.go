/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue struct ***REMOVED***
	buf               []interface***REMOVED******REMOVED***
	head, tail, count int
***REMOVED***

// New constructs and returns a new Queue.
func New() *Queue ***REMOVED***
	return &Queue***REMOVED***
		buf: make([]interface***REMOVED******REMOVED***, minQueueLen),
	***REMOVED***
***REMOVED***

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int ***REMOVED***
	return q.count
***REMOVED***

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() ***REMOVED***
	newBuf := make([]interface***REMOVED******REMOVED***, q.count<<1)

	if q.tail > q.head ***REMOVED***
		copy(newBuf, q.buf[q.head:q.tail])
	***REMOVED*** else ***REMOVED***
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	***REMOVED***

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
***REMOVED***

// Add puts an element on the end of the queue.
func (q *Queue) Add(elem interface***REMOVED******REMOVED***) ***REMOVED***
	if q.count == len(q.buf) ***REMOVED***
		q.resize()
	***REMOVED***

	q.buf[q.tail] = elem
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
***REMOVED***

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Queue) Peek() interface***REMOVED******REMOVED*** ***REMOVED***
	if q.count <= 0 ***REMOVED***
		panic("queue: Peek() called on empty queue")
	***REMOVED***
	return q.buf[q.head]
***REMOVED***

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue) Get(i int) interface***REMOVED******REMOVED*** ***REMOVED***
	// If indexing backwards, convert to positive index.
	if i < 0 ***REMOVED***
		i += q.count
	***REMOVED***
	if i < 0 || i >= q.count ***REMOVED***
		panic("queue: Get() called with index out of range")
	***REMOVED***
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
***REMOVED***

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Queue) Remove() interface***REMOVED******REMOVED*** ***REMOVED***
	if q.count <= 0 ***REMOVED***
		panic("queue: Remove() called on empty queue")
	***REMOVED***
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) ***REMOVED***
		q.resize()
	***REMOVED***
	return ret
***REMOVED***
