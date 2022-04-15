/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package types

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// ipBlock represents a continuous segment of IP addresses
type ipBlock struct ***REMOVED***
	firstIP, count *big.Int
	ipv6           bool
***REMOVED***

// ipPoolBlock is similar to ipBlock but instead of knowing its count/size it knows the first index
// from which it starts in an IPPool
type ipPoolBlock struct ***REMOVED***
	firstIP, startIndex *big.Int
***REMOVED***

// IPPool represent a slice of IPBlocks
type IPPool struct ***REMOVED***
	list  []ipPoolBlock
	count *big.Int
***REMOVED***

func getIPBlock(s string) (*ipBlock, error) ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(s, "-"):
		return ipBlockFromRange(s)
	case strings.Contains(s, "/"):
		return ipBlockFromCIDR(s)
	default:
		if net.ParseIP(s) == nil ***REMOVED***
			return nil, fmt.Errorf("%s is not a valid IP, IP range or CIDR", s)
		***REMOVED***
		return ipBlockFromRange(s + "-" + s)
	***REMOVED***
***REMOVED***

func ipBlockFromRange(s string) (*ipBlock, error) ***REMOVED***
	ss := strings.SplitN(s, "-", 2)
	ip0, ip1 := net.ParseIP(ss[0]), net.ParseIP(ss[1])
	if ip0 == nil || ip1 == nil ***REMOVED***
		return nil, errors.New("wrong IP range format: " + s)
	***REMOVED***
	if (ip0.To4() == nil) != (ip1.To4() == nil) ***REMOVED*** // XOR
		return nil, errors.New("mixed IP range format: " + s)
	***REMOVED***
	block := ipBlockFromTwoIPs(ip0, ip1)

	if block.count.Sign() <= 0 ***REMOVED***
		return nil, errors.New("negative IP range: " + s)
	***REMOVED***
	return block, nil
***REMOVED***

func ipBlockFromTwoIPs(ip0, ip1 net.IP) *ipBlock ***REMOVED***
	// This code doesn't do any checks on the validity of the arguments, that should be
	// done before and/or after it is called
	var block ipBlock
	block.firstIP = new(big.Int)
	block.count = new(big.Int)
	block.ipv6 = ip0.To4() == nil
	if block.ipv6 ***REMOVED***
		block.firstIP.SetBytes(ip0.To16())
		block.count.SetBytes(ip1.To16())
	***REMOVED*** else ***REMOVED***
		block.firstIP.SetBytes(ip0.To4())
		block.count.SetBytes(ip1.To4())
	***REMOVED***
	block.count.Sub(block.count, block.firstIP)
	block.count.Add(block.count, big.NewInt(1))

	return &block
***REMOVED***

func ipBlockFromCIDR(s string) (*ipBlock, error) ***REMOVED***
	_, pnet, err := net.ParseCIDR(s)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("parseCIDR() failed parsing %s: %w", s, err)
	***REMOVED***
	ip0 := pnet.IP
	// TODO: this is just to copy it, it will probably be better to copy the bytes ...
	ip1 := net.ParseIP(ip0.String())
	if ip1.To4() == nil ***REMOVED***
		ip1 = ip1.To16()
	***REMOVED*** else ***REMOVED***
		ip1 = ip1.To4()
	***REMOVED***
	for i := range ip1 ***REMOVED***
		ip1[i] |= (255 ^ pnet.Mask[i])
	***REMOVED***
	block := ipBlockFromTwoIPs(ip0, ip1)
	// in the case of ipv4 if the network is bigger than 31 the first and last IP are reserved so we
	// need to reduce the addresses by 2 and increment the first ip
	if !block.ipv6 && big.NewInt(2).Cmp(block.count) < 0 ***REMOVED***
		block.count.Sub(block.count, big.NewInt(2))
		block.firstIP.Add(block.firstIP, big.NewInt(1))
	***REMOVED***
	return block, nil
***REMOVED***

func (b ipPoolBlock) getIP(index *big.Int) net.IP ***REMOVED***
	// TODO implement walking ipv6 networks first
	// that will probably require more math ... including knowing which is the next network and ...
	// thinking about it - it looks like it's going to be kind of hard or badly defined
	i := new(big.Int)
	i.Add(b.firstIP, index)
	// TODO use big.Int.FillBytes when golang 1.14 is no longer supported
	return net.IP(i.Bytes())
***REMOVED***

// NewIPPool returns an IPPool slice from the provided string representation that should be comma
// separated list of IPs, IP ranges(ip1-ip2) and CIDRs
func NewIPPool(ranges string) (*IPPool, error) ***REMOVED***
	ss := strings.Split(ranges, ",")
	pool := &IPPool***REMOVED******REMOVED***
	pool.list = make([]ipPoolBlock, len(ss))
	pool.count = new(big.Int)
	for i, bs := range ss ***REMOVED***
		r, err := getIPBlock(bs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		pool.list[i] = ipPoolBlock***REMOVED***
			firstIP:    r.firstIP,
			startIndex: new(big.Int).Set(pool.count), // this is how many there are until now
		***REMOVED***
		pool.count.Add(pool.count, r.count)
	***REMOVED***

	// The list gets reversed here as later it is searched based on when the index we are looking is
	// bigger than startIndex but it will be true always for the first block which is with
	// startIndex 0. This can also be fixed by iterating in reverse but this seems better
	for i := 0; i < len(pool.list)/2; i++ ***REMOVED***
		pool.list[i], pool.list[len(pool.list)-1-i] = pool.list[len(pool.list)-1-i], pool.list[i]
	***REMOVED***
	return pool, nil
***REMOVED***

// GetIP return an IP from a pool of IPBlock slice
func (pool *IPPool) GetIP(index uint64) net.IP ***REMOVED***
	return pool.GetIPBig(new(big.Int).SetUint64(index))
***REMOVED***

// GetIPBig returns an IP from the pool with the provided index that is big.Int
func (pool *IPPool) GetIPBig(index *big.Int) net.IP ***REMOVED***
	index = new(big.Int).Rem(index, pool.count)
	for _, b := range pool.list ***REMOVED***
		if index.Cmp(b.startIndex) >= 0 ***REMOVED***
			return b.getIP(index.Sub(index, b.startIndex))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NullIPPool is a nullable IPPool
type NullIPPool struct ***REMOVED***
	Pool  *IPPool
	Valid bool
	raw   []byte
***REMOVED***

// UnmarshalText converts text data to a valid NullIPPool
func (n *NullIPPool) UnmarshalText(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		*n = NullIPPool***REMOVED******REMOVED***
		return nil
	***REMOVED***
	var err error
	n.Pool, err = NewIPPool(string(data))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n.Valid = true
	n.raw = data
	return nil
***REMOVED***

// MarshalText returns the IPs pool in text form
func (n *NullIPPool) MarshalText() ([]byte, error) ***REMOVED***
	return n.raw, nil
***REMOVED***
