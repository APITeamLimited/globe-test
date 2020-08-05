package goja

import (
	"hash/maphash"
)

type mapEntry struct ***REMOVED***
	key, value Value

	iterPrev, iterNext *mapEntry
	hNext              *mapEntry
***REMOVED***

type orderedMap struct ***REMOVED***
	hash                *maphash.Hash
	hashTable           map[uint64]*mapEntry
	iterFirst, iterLast *mapEntry
	size                int
***REMOVED***

type orderedMapIter struct ***REMOVED***
	m   *orderedMap
	cur *mapEntry
***REMOVED***

func (m *orderedMap) lookup(key Value) (h uint64, entry, hPrev *mapEntry) ***REMOVED***
	if key == _negativeZero ***REMOVED***
		key = intToValue(0)
	***REMOVED***
	h = key.hash(m.hash)
	for entry = m.hashTable[h]; entry != nil && !entry.key.SameAs(key); hPrev, entry = entry, entry.hNext ***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (m *orderedMap) set(key, value Value) ***REMOVED***
	h, entry, hPrev := m.lookup(key)
	if entry != nil ***REMOVED***
		entry.value = value
	***REMOVED*** else ***REMOVED***
		if key == _negativeZero ***REMOVED***
			key = intToValue(0)
		***REMOVED***
		entry = &mapEntry***REMOVED***key: key, value: value***REMOVED***
		if hPrev == nil ***REMOVED***
			m.hashTable[h] = entry
		***REMOVED*** else ***REMOVED***
			hPrev.hNext = entry
		***REMOVED***
		if m.iterLast != nil ***REMOVED***
			entry.iterPrev = m.iterLast
			m.iterLast.iterNext = entry
		***REMOVED*** else ***REMOVED***
			m.iterFirst = entry
		***REMOVED***
		m.iterLast = entry
		m.size++
	***REMOVED***
***REMOVED***

func (m *orderedMap) get(key Value) Value ***REMOVED***
	_, entry, _ := m.lookup(key)
	if entry != nil ***REMOVED***
		return entry.value
	***REMOVED***

	return nil
***REMOVED***

func (m *orderedMap) remove(key Value) bool ***REMOVED***
	h, entry, hPrev := m.lookup(key)
	if entry != nil ***REMOVED***
		entry.key = nil
		entry.value = nil

		// remove from the doubly-linked list
		if entry.iterPrev != nil ***REMOVED***
			entry.iterPrev.iterNext = entry.iterNext
		***REMOVED*** else ***REMOVED***
			m.iterFirst = entry.iterNext
		***REMOVED***
		if entry.iterNext != nil ***REMOVED***
			entry.iterNext.iterPrev = entry.iterPrev
		***REMOVED*** else ***REMOVED***
			m.iterLast = entry.iterPrev
		***REMOVED***

		// remove from the hashTable
		if hPrev == nil ***REMOVED***
			if entry.hNext == nil ***REMOVED***
				delete(m.hashTable, h)
			***REMOVED*** else ***REMOVED***
				m.hashTable[h] = entry.hNext
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			hPrev.hNext = entry.hNext
		***REMOVED***

		m.size--
		return true
	***REMOVED***

	return false
***REMOVED***

func (m *orderedMap) has(key Value) bool ***REMOVED***
	_, entry, _ := m.lookup(key)
	return entry != nil
***REMOVED***

func (iter *orderedMapIter) next() *mapEntry ***REMOVED***
	if iter.m == nil ***REMOVED***
		// closed iterator
		return nil
	***REMOVED***

	cur := iter.cur
	// if the current item was deleted, track back to find the latest that wasn't
	for cur != nil && cur.key == nil ***REMOVED***
		cur = cur.iterPrev
	***REMOVED***

	if cur != nil ***REMOVED***
		cur = cur.iterNext
	***REMOVED*** else ***REMOVED***
		cur = iter.m.iterFirst
	***REMOVED***

	if cur == nil ***REMOVED***
		iter.close()
	***REMOVED*** else ***REMOVED***
		iter.cur = cur
	***REMOVED***

	return cur
***REMOVED***

func (iter *orderedMapIter) close() ***REMOVED***
	iter.m = nil
	iter.cur = nil
***REMOVED***

func newOrderedMap(h *maphash.Hash) *orderedMap ***REMOVED***
	return &orderedMap***REMOVED***
		hash:      h,
		hashTable: make(map[uint64]*mapEntry),
	***REMOVED***
***REMOVED***

func (m *orderedMap) newIter() *orderedMapIter ***REMOVED***
	iter := &orderedMapIter***REMOVED***
		m: m,
	***REMOVED***
	return iter
***REMOVED***

func (m *orderedMap) clear() ***REMOVED***
	for item := m.iterFirst; item != nil; item = item.iterNext ***REMOVED***
		item.key = nil
		item.value = nil
		if item.iterPrev != nil ***REMOVED***
			item.iterPrev.iterNext = nil
		***REMOVED***
	***REMOVED***
	m.iterFirst = nil
	m.iterLast = nil
	m.hashTable = make(map[uint64]*mapEntry)
	m.size = 0
***REMOVED***
