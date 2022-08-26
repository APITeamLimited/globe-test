package redis

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
)

type Scripter interface ***REMOVED***
	Eval(ctx context.Context, script string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd
	ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *StringCmd
***REMOVED***

var (
	_ Scripter = (*Client)(nil)
	_ Scripter = (*Ring)(nil)
	_ Scripter = (*ClusterClient)(nil)
)

type Script struct ***REMOVED***
	src, hash string
***REMOVED***

func NewScript(src string) *Script ***REMOVED***
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script***REMOVED***
		src:  src,
		hash: hex.EncodeToString(h.Sum(nil)),
	***REMOVED***
***REMOVED***

func (s *Script) Hash() string ***REMOVED***
	return s.hash
***REMOVED***

func (s *Script) Load(ctx context.Context, c Scripter) *StringCmd ***REMOVED***
	return c.ScriptLoad(ctx, s.src)
***REMOVED***

func (s *Script) Exists(ctx context.Context, c Scripter) *BoolSliceCmd ***REMOVED***
	return c.ScriptExists(ctx, s.hash)
***REMOVED***

func (s *Script) Eval(ctx context.Context, c Scripter, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	return c.Eval(ctx, s.src, keys, args...)
***REMOVED***

func (s *Script) EvalSha(ctx context.Context, c Scripter, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	return c.EvalSha(ctx, s.hash, keys, args...)
***REMOVED***

// Run optimistically uses EVALSHA to run the script. If script does not exist
// it is retried using EVAL.
func (s *Script) Run(ctx context.Context, c Scripter, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	r := s.EvalSha(ctx, c, keys, args...)
	if err := r.Err(); err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") ***REMOVED***
		return s.Eval(ctx, c, keys, args...)
	***REMOVED***
	return r
***REMOVED***
