package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/APITeamLimited/redis/v9/internal"
	"github.com/APITeamLimited/redis/v9/internal/hscan"
	"github.com/APITeamLimited/redis/v9/internal/proto"
	"github.com/APITeamLimited/redis/v9/internal/util"
)

type Cmder interface ***REMOVED***
	Name() string
	FullName() string
	Args() []interface***REMOVED******REMOVED***
	String() string
	stringArg(int) string
	firstKeyPos() int8
	SetFirstKeyPos(int8)

	readTimeout() *time.Duration
	readReply(rd *proto.Reader) error

	SetErr(error)
	Err() error
***REMOVED***

func setCmdsErr(cmds []Cmder, e error) ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		if cmd.Err() == nil ***REMOVED***
			cmd.SetErr(e)
		***REMOVED***
	***REMOVED***
***REMOVED***

func cmdsFirstErr(cmds []Cmder) error ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		if err := cmd.Err(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func writeCmds(wr *proto.Writer, cmds []Cmder) error ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		if err := writeCmd(wr, cmd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func writeCmd(wr *proto.Writer, cmd Cmder) error ***REMOVED***
	return wr.WriteArgs(cmd.Args())
***REMOVED***

func cmdFirstKeyPos(cmd Cmder, info *CommandInfo) int ***REMOVED***
	if pos := cmd.firstKeyPos(); pos != 0 ***REMOVED***
		return int(pos)
	***REMOVED***

	switch cmd.Name() ***REMOVED***
	case "eval", "evalsha":
		if cmd.stringArg(2) != "0" ***REMOVED***
			return 3
		***REMOVED***

		return 0
	case "publish":
		return 1
	case "memory":
		// https://github.com/redis/redis/issues/7493
		if cmd.stringArg(1) == "usage" ***REMOVED***
			return 2
		***REMOVED***
	***REMOVED***

	if info != nil ***REMOVED***
		return int(info.FirstKeyPos)
	***REMOVED***
	return 1
***REMOVED***

func cmdString(cmd Cmder, val interface***REMOVED******REMOVED***) string ***REMOVED***
	b := make([]byte, 0, 64)

	for i, arg := range cmd.Args() ***REMOVED***
		if i > 0 ***REMOVED***
			b = append(b, ' ')
		***REMOVED***
		b = internal.AppendArg(b, arg)
	***REMOVED***

	if err := cmd.Err(); err != nil ***REMOVED***
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	***REMOVED*** else if val != nil ***REMOVED***
		b = append(b, ": "...)
		b = internal.AppendArg(b, val)
	***REMOVED***

	return util.BytesToString(b)
***REMOVED***

//------------------------------------------------------------------------------

type baseCmd struct ***REMOVED***
	ctx    context.Context
	args   []interface***REMOVED******REMOVED***
	err    error
	keyPos int8

	_readTimeout *time.Duration
***REMOVED***

var _ Cmder = (*Cmd)(nil)

func (cmd *baseCmd) Name() string ***REMOVED***
	if len(cmd.args) == 0 ***REMOVED***
		return ""
	***REMOVED***
	// Cmd name must be lower cased.
	return internal.ToLower(cmd.stringArg(0))
***REMOVED***

func (cmd *baseCmd) FullName() string ***REMOVED***
	switch name := cmd.Name(); name ***REMOVED***
	case "cluster", "command":
		if len(cmd.args) == 1 ***REMOVED***
			return name
		***REMOVED***
		if s2, ok := cmd.args[1].(string); ok ***REMOVED***
			return name + " " + s2
		***REMOVED***
		return name
	default:
		return name
	***REMOVED***
***REMOVED***

func (cmd *baseCmd) Args() []interface***REMOVED******REMOVED*** ***REMOVED***
	return cmd.args
***REMOVED***

func (cmd *baseCmd) stringArg(pos int) string ***REMOVED***
	if pos < 0 || pos >= len(cmd.args) ***REMOVED***
		return ""
	***REMOVED***
	arg := cmd.args[pos]
	switch v := arg.(type) ***REMOVED***
	case string:
		return v
	default:
		// TODO: consider using appendArg
		return fmt.Sprint(v)
	***REMOVED***
***REMOVED***

func (cmd *baseCmd) firstKeyPos() int8 ***REMOVED***
	return cmd.keyPos
***REMOVED***

func (cmd *baseCmd) SetFirstKeyPos(keyPos int8) ***REMOVED***
	cmd.keyPos = keyPos
***REMOVED***

func (cmd *baseCmd) SetErr(e error) ***REMOVED***
	cmd.err = e
***REMOVED***

func (cmd *baseCmd) Err() error ***REMOVED***
	return cmd.err
***REMOVED***

func (cmd *baseCmd) readTimeout() *time.Duration ***REMOVED***
	return cmd._readTimeout
***REMOVED***

func (cmd *baseCmd) setReadTimeout(d time.Duration) ***REMOVED***
	cmd._readTimeout = &d
***REMOVED***

//------------------------------------------------------------------------------

type Cmd struct ***REMOVED***
	baseCmd

	val interface***REMOVED******REMOVED***
***REMOVED***

func NewCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	return &Cmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *Cmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *Cmd) SetVal(val interface***REMOVED******REMOVED***) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *Cmd) Val() interface***REMOVED******REMOVED*** ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *Cmd) Result() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *Cmd) Text() (string, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return "", cmd.err
	***REMOVED***
	return toString(cmd.val)
***REMOVED***

func toString(val interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case string:
		return val, nil
	default:
		err := fmt.Errorf("redis: unexpected type=%T for String", val)
		return "", err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Int() (int, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	switch val := cmd.val.(type) ***REMOVED***
	case int64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Int", val)
		return 0, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Int64() (int64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return toInt64(cmd.val)
***REMOVED***

func toInt64(val interface***REMOVED******REMOVED***) (int64, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Int64", val)
		return 0, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Uint64() (uint64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return toUint64(cmd.val)
***REMOVED***

func toUint64(val interface***REMOVED******REMOVED***) (uint64, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case int64:
		return uint64(val), nil
	case string:
		return strconv.ParseUint(val, 10, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Uint64", val)
		return 0, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Float32() (float32, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return toFloat32(cmd.val)
***REMOVED***

func toFloat32(val interface***REMOVED******REMOVED***) (float32, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case int64:
		return float32(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return float32(f), nil
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Float32", val)
		return 0, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Float64() (float64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return toFloat64(cmd.val)
***REMOVED***

func toFloat64(val interface***REMOVED******REMOVED***) (float64, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Float64", val)
		return 0, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Bool() (bool, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return false, cmd.err
	***REMOVED***
	return toBool(cmd.val)
***REMOVED***

func toBool(val interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case int64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Bool", val)
		return false, err
	***REMOVED***
***REMOVED***

func (cmd *Cmd) Slice() ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return nil, cmd.err
	***REMOVED***
	switch val := cmd.val.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		return val, nil
	default:
		return nil, fmt.Errorf("redis: unexpected type=%T for Slice", val)
	***REMOVED***
***REMOVED***

func (cmd *Cmd) StringSlice() ([]string, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ss := make([]string, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toString(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ss[i] = val
	***REMOVED***
	return ss, nil
***REMOVED***

func (cmd *Cmd) Int64Slice() ([]int64, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	nums := make([]int64, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toInt64(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		nums[i] = val
	***REMOVED***
	return nums, nil
***REMOVED***

func (cmd *Cmd) Uint64Slice() ([]uint64, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	nums := make([]uint64, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toUint64(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		nums[i] = val
	***REMOVED***
	return nums, nil
***REMOVED***

func (cmd *Cmd) Float32Slice() ([]float32, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	floats := make([]float32, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toFloat32(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		floats[i] = val
	***REMOVED***
	return floats, nil
***REMOVED***

func (cmd *Cmd) Float64Slice() ([]float64, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	floats := make([]float64, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toFloat64(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		floats[i] = val
	***REMOVED***
	return floats, nil
***REMOVED***

func (cmd *Cmd) BoolSlice() ([]bool, error) ***REMOVED***
	slice, err := cmd.Slice()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	bools := make([]bool, len(slice))
	for i, iface := range slice ***REMOVED***
		val, err := toBool(iface)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		bools[i] = val
	***REMOVED***
	return bools, nil
***REMOVED***

func (cmd *Cmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadReply()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type SliceCmd struct ***REMOVED***
	baseCmd

	val []interface***REMOVED******REMOVED***
***REMOVED***

var _ Cmder = (*SliceCmd)(nil)

func NewSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *SliceCmd ***REMOVED***
	return &SliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *SliceCmd) SetVal(val []interface***REMOVED******REMOVED***) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *SliceCmd) Val() []interface***REMOVED******REMOVED*** ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *SliceCmd) Result() ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *SliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

// Scan scans the results from the map into a destination struct. The map keys
// are matched in the Redis struct fields by the `redis:"field"` tag.
func (cmd *SliceCmd) Scan(dst interface***REMOVED******REMOVED***) error ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return cmd.err
	***REMOVED***

	// Pass the list of keys and values.
	// Skip the first two args for: HMGET key
	var args []interface***REMOVED******REMOVED***
	if cmd.args[0] == "hmget" ***REMOVED***
		args = cmd.args[2:]
	***REMOVED*** else ***REMOVED***
		// Otherwise, it's: MGET field field ...
		args = cmd.args[1:]
	***REMOVED***

	return hscan.Scan(dst, args, cmd.val)
***REMOVED***

func (cmd *SliceCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadSlice()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type StatusCmd struct ***REMOVED***
	baseCmd

	val string
***REMOVED***

var _ Cmder = (*StatusCmd)(nil)

func NewStatusCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *StatusCmd ***REMOVED***
	return &StatusCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *StatusCmd) SetVal(val string) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *StatusCmd) Val() string ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *StatusCmd) Result() (string, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *StatusCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *StatusCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadString()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type IntCmd struct ***REMOVED***
	baseCmd

	val int64
***REMOVED***

var _ Cmder = (*IntCmd)(nil)

func NewIntCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	return &IntCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *IntCmd) SetVal(val int64) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *IntCmd) Val() int64 ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *IntCmd) Result() (int64, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *IntCmd) Uint64() (uint64, error) ***REMOVED***
	return uint64(cmd.val), cmd.err
***REMOVED***

func (cmd *IntCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *IntCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadInt()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type IntSliceCmd struct ***REMOVED***
	baseCmd

	val []int64
***REMOVED***

var _ Cmder = (*IntSliceCmd)(nil)

func NewIntSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *IntSliceCmd ***REMOVED***
	return &IntSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *IntSliceCmd) SetVal(val []int64) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *IntSliceCmd) Val() []int64 ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *IntSliceCmd) Result() ([]int64, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *IntSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *IntSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]int64, n)
	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if cmd.val[i], err = rd.ReadInt(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type DurationCmd struct ***REMOVED***
	baseCmd

	val       time.Duration
	precision time.Duration
***REMOVED***

var _ Cmder = (*DurationCmd)(nil)

func NewDurationCmd(ctx context.Context, precision time.Duration, args ...interface***REMOVED******REMOVED***) *DurationCmd ***REMOVED***
	return &DurationCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
		precision: precision,
	***REMOVED***
***REMOVED***

func (cmd *DurationCmd) SetVal(val time.Duration) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *DurationCmd) Val() time.Duration ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *DurationCmd) Result() (time.Duration, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *DurationCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *DurationCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadInt()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch n ***REMOVED***
	// -2 if the key does not exist
	// -1 if the key exists but has no associated expire
	case -2, -1:
		cmd.val = time.Duration(n)
	default:
		cmd.val = time.Duration(n) * cmd.precision
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type TimeCmd struct ***REMOVED***
	baseCmd

	val time.Time
***REMOVED***

var _ Cmder = (*TimeCmd)(nil)

func NewTimeCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *TimeCmd ***REMOVED***
	return &TimeCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *TimeCmd) SetVal(val time.Time) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *TimeCmd) Val() time.Time ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *TimeCmd) Result() (time.Time, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *TimeCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *TimeCmd) readReply(rd *proto.Reader) error ***REMOVED***
	if err := rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
		return err
	***REMOVED***
	second, err := rd.ReadInt()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	microsecond, err := rd.ReadInt()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = time.Unix(second, microsecond*1000)
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type BoolCmd struct ***REMOVED***
	baseCmd

	val bool
***REMOVED***

var _ Cmder = (*BoolCmd)(nil)

func NewBoolCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	return &BoolCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *BoolCmd) SetVal(val bool) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *BoolCmd) Val() bool ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *BoolCmd) Result() (bool, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *BoolCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *BoolCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadBool()

	// `SET key value NX` returns nil when key already exists. But
	// `SETNX key value` returns bool (0/1). So convert nil to bool.
	if err == Nil ***REMOVED***
		cmd.val = false
		err = nil
	***REMOVED***
	return err
***REMOVED***

//------------------------------------------------------------------------------

type StringCmd struct ***REMOVED***
	baseCmd

	val string
***REMOVED***

var _ Cmder = (*StringCmd)(nil)

func NewStringCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *StringCmd ***REMOVED***
	return &StringCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *StringCmd) SetVal(val string) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *StringCmd) Val() string ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *StringCmd) Result() (string, error) ***REMOVED***
	return cmd.Val(), cmd.err
***REMOVED***

func (cmd *StringCmd) Bytes() ([]byte, error) ***REMOVED***
	return util.StringToBytes(cmd.val), cmd.err
***REMOVED***

func (cmd *StringCmd) Bool() (bool, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return false, cmd.err
	***REMOVED***
	return strconv.ParseBool(cmd.val)
***REMOVED***

func (cmd *StringCmd) Int() (int, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return strconv.Atoi(cmd.Val())
***REMOVED***

func (cmd *StringCmd) Int64() (int64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return strconv.ParseInt(cmd.Val(), 10, 64)
***REMOVED***

func (cmd *StringCmd) Uint64() (uint64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return strconv.ParseUint(cmd.Val(), 10, 64)
***REMOVED***

func (cmd *StringCmd) Float32() (float32, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	f, err := strconv.ParseFloat(cmd.Val(), 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return float32(f), nil
***REMOVED***

func (cmd *StringCmd) Float64() (float64, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return 0, cmd.err
	***REMOVED***
	return strconv.ParseFloat(cmd.Val(), 64)
***REMOVED***

func (cmd *StringCmd) Time() (time.Time, error) ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, cmd.err
	***REMOVED***
	return time.Parse(time.RFC3339Nano, cmd.Val())
***REMOVED***

func (cmd *StringCmd) Scan(val interface***REMOVED******REMOVED***) error ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return cmd.err
	***REMOVED***
	return proto.Scan([]byte(cmd.val), val)
***REMOVED***

func (cmd *StringCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *StringCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadString()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type FloatCmd struct ***REMOVED***
	baseCmd

	val float64
***REMOVED***

var _ Cmder = (*FloatCmd)(nil)

func NewFloatCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *FloatCmd ***REMOVED***
	return &FloatCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *FloatCmd) SetVal(val float64) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *FloatCmd) Val() float64 ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *FloatCmd) Result() (float64, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *FloatCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *FloatCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = rd.ReadFloat()
	return err
***REMOVED***

//------------------------------------------------------------------------------

type FloatSliceCmd struct ***REMOVED***
	baseCmd

	val []float64
***REMOVED***

var _ Cmder = (*FloatSliceCmd)(nil)

func NewFloatSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *FloatSliceCmd ***REMOVED***
	return &FloatSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *FloatSliceCmd) SetVal(val []float64) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *FloatSliceCmd) Val() []float64 ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *FloatSliceCmd) Result() ([]float64, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *FloatSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *FloatSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make([]float64, n)
	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		switch num, err := rd.ReadFloat(); ***REMOVED***
		case err == Nil:
			cmd.val[i] = 0
		case err != nil:
			return err
		default:
			cmd.val[i] = num
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type StringSliceCmd struct ***REMOVED***
	baseCmd

	val []string
***REMOVED***

var _ Cmder = (*StringSliceCmd)(nil)

func NewStringSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *StringSliceCmd ***REMOVED***
	return &StringSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *StringSliceCmd) SetVal(val []string) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *StringSliceCmd) Val() []string ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *StringSliceCmd) Result() ([]string, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *StringSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *StringSliceCmd) ScanSlice(container interface***REMOVED******REMOVED***) error ***REMOVED***
	return proto.ScanSlice(cmd.Val(), container)
***REMOVED***

func (cmd *StringSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]string, n)
	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		switch s, err := rd.ReadString(); ***REMOVED***
		case err == Nil:
			cmd.val[i] = ""
		case err != nil:
			return err
		default:
			cmd.val[i] = s
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type KeyValue struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

type KeyValueSliceCmd struct ***REMOVED***
	baseCmd

	val []KeyValue
***REMOVED***

var _ Cmder = (*KeyValueSliceCmd)(nil)

func NewKeyValueSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *KeyValueSliceCmd ***REMOVED***
	return &KeyValueSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *KeyValueSliceCmd) SetVal(val []KeyValue) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *KeyValueSliceCmd) Val() []KeyValue ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *KeyValueSliceCmd) Result() ([]KeyValue, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *KeyValueSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

// Many commands will respond to two formats:
//  1) 1) "one"
//     2) (double) 1
//  2) 1) "two"
//     2) (double) 2
// OR:
//  1) "two"
//  2) (double) 2
//  3) "one"
//  4) (double) 1
func (cmd *KeyValueSliceCmd) readReply(rd *proto.Reader) error ***REMOVED*** // nolint:dupl
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the n is 0, can't continue reading.
	if n == 0 ***REMOVED***
		cmd.val = make([]KeyValue, 0)
		return nil
	***REMOVED***

	typ, err := rd.PeekReplyType()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	array := typ == proto.RespArray

	if array ***REMOVED***
		cmd.val = make([]KeyValue, n)
	***REMOVED*** else ***REMOVED***
		cmd.val = make([]KeyValue, n/2)
	***REMOVED***

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if array ***REMOVED***
			if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if cmd.val[i].Key, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***

		if cmd.val[i].Value, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type BoolSliceCmd struct ***REMOVED***
	baseCmd

	val []bool
***REMOVED***

var _ Cmder = (*BoolSliceCmd)(nil)

func NewBoolSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *BoolSliceCmd ***REMOVED***
	return &BoolSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *BoolSliceCmd) SetVal(val []bool) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *BoolSliceCmd) Val() []bool ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *BoolSliceCmd) Result() ([]bool, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *BoolSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *BoolSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]bool, n)
	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if cmd.val[i], err = rd.ReadBool(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type MapStringStringCmd struct ***REMOVED***
	baseCmd

	val map[string]string
***REMOVED***

var _ Cmder = (*MapStringStringCmd)(nil)

func NewMapStringStringCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *MapStringStringCmd ***REMOVED***
	return &MapStringStringCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *MapStringStringCmd) Val() map[string]string ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *MapStringStringCmd) SetVal(val map[string]string) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *MapStringStringCmd) Result() (map[string]string, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *MapStringStringCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

// Scan scans the results from the map into a destination struct. The map keys
// are matched in the Redis struct fields by the `redis:"field"` tag.
func (cmd *MapStringStringCmd) Scan(dest interface***REMOVED******REMOVED***) error ***REMOVED***
	if cmd.err != nil ***REMOVED***
		return cmd.err
	***REMOVED***

	strct, err := hscan.Struct(dest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for k, v := range cmd.val ***REMOVED***
		if err := strct.Scan(k, v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (cmd *MapStringStringCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make(map[string]string, n)
	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		value, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		cmd.val[key] = value
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type MapStringIntCmd struct ***REMOVED***
	baseCmd

	val map[string]int64
***REMOVED***

var _ Cmder = (*MapStringIntCmd)(nil)

func NewMapStringIntCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *MapStringIntCmd ***REMOVED***
	return &MapStringIntCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *MapStringIntCmd) SetVal(val map[string]int64) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *MapStringIntCmd) Val() map[string]int64 ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *MapStringIntCmd) Result() (map[string]int64, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *MapStringIntCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *MapStringIntCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make(map[string]int64, n)
	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		nn, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[key] = nn
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type StringStructMapCmd struct ***REMOVED***
	baseCmd

	val map[string]struct***REMOVED******REMOVED***
***REMOVED***

var _ Cmder = (*StringStructMapCmd)(nil)

func NewStringStructMapCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *StringStructMapCmd ***REMOVED***
	return &StringStructMapCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *StringStructMapCmd) SetVal(val map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *StringStructMapCmd) Val() map[string]struct***REMOVED******REMOVED*** ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *StringStructMapCmd) Result() (map[string]struct***REMOVED******REMOVED***, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *StringStructMapCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *StringStructMapCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make(map[string]struct***REMOVED******REMOVED***, n)
	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XMessage struct ***REMOVED***
	ID     string
	Values map[string]interface***REMOVED******REMOVED***
***REMOVED***

type XMessageSliceCmd struct ***REMOVED***
	baseCmd

	val []XMessage
***REMOVED***

var _ Cmder = (*XMessageSliceCmd)(nil)

func NewXMessageSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XMessageSliceCmd ***REMOVED***
	return &XMessageSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XMessageSliceCmd) SetVal(val []XMessage) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XMessageSliceCmd) Val() []XMessage ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XMessageSliceCmd) Result() ([]XMessage, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XMessageSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XMessageSliceCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	cmd.val, err = readXMessageSlice(rd)
	return err
***REMOVED***

func readXMessageSlice(rd *proto.Reader) ([]XMessage, error) ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	msgs := make([]XMessage, n)
	for i := 0; i < len(msgs); i++ ***REMOVED***
		if msgs[i], err = readXMessage(rd); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return msgs, nil
***REMOVED***

func readXMessage(rd *proto.Reader) (XMessage, error) ***REMOVED***
	if err := rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
		return XMessage***REMOVED******REMOVED***, err
	***REMOVED***

	id, err := rd.ReadString()
	if err != nil ***REMOVED***
		return XMessage***REMOVED******REMOVED***, err
	***REMOVED***

	v, err := stringInterfaceMapParser(rd)
	if err != nil ***REMOVED***
		if err != proto.Nil ***REMOVED***
			return XMessage***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return XMessage***REMOVED***
		ID:     id,
		Values: v,
	***REMOVED***, nil
***REMOVED***

func stringInterfaceMapParser(rd *proto.Reader) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	m := make(map[string]interface***REMOVED******REMOVED***, n)
	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		value, err := rd.ReadString()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		m[key] = value
	***REMOVED***
	return m, nil
***REMOVED***

//------------------------------------------------------------------------------

type XStream struct ***REMOVED***
	Stream   string
	Messages []XMessage
***REMOVED***

type XStreamSliceCmd struct ***REMOVED***
	baseCmd

	val []XStream
***REMOVED***

var _ Cmder = (*XStreamSliceCmd)(nil)

func NewXStreamSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XStreamSliceCmd ***REMOVED***
	return &XStreamSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XStreamSliceCmd) SetVal(val []XStream) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XStreamSliceCmd) Val() []XStream ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XStreamSliceCmd) Result() ([]XStream, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XStreamSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XStreamSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	typ, err := rd.PeekReplyType()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var n int
	if typ == proto.RespMap ***REMOVED***
		n, err = rd.ReadMapLen()
	***REMOVED*** else ***REMOVED***
		n, err = rd.ReadArrayLen()
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]XStream, n)
	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if typ != proto.RespMap ***REMOVED***
			if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if cmd.val[i].Stream, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if cmd.val[i].Messages, err = readXMessageSlice(rd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XPending struct ***REMOVED***
	Count     int64
	Lower     string
	Higher    string
	Consumers map[string]int64
***REMOVED***

type XPendingCmd struct ***REMOVED***
	baseCmd
	val *XPending
***REMOVED***

var _ Cmder = (*XPendingCmd)(nil)

func NewXPendingCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XPendingCmd ***REMOVED***
	return &XPendingCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XPendingCmd) SetVal(val *XPending) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XPendingCmd) Val() *XPending ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XPendingCmd) Result() (*XPending, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XPendingCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XPendingCmd) readReply(rd *proto.Reader) error ***REMOVED***
	var err error
	if err = rd.ReadFixedArrayLen(4); err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = &XPending***REMOVED******REMOVED***

	if cmd.val.Count, err = rd.ReadInt(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if cmd.val.Lower, err = rd.ReadString(); err != nil && err != Nil ***REMOVED***
		return err
	***REMOVED***

	if cmd.val.Higher, err = rd.ReadString(); err != nil && err != Nil ***REMOVED***
		return err
	***REMOVED***

	n, err := rd.ReadArrayLen()
	if err != nil && err != Nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val.Consumers = make(map[string]int64, n)
	for i := 0; i < n; i++ ***REMOVED***
		if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
			return err
		***REMOVED***

		consumerName, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		consumerPending, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val.Consumers[consumerName] = consumerPending
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XPendingExt struct ***REMOVED***
	ID         string
	Consumer   string
	Idle       time.Duration
	RetryCount int64
***REMOVED***

type XPendingExtCmd struct ***REMOVED***
	baseCmd
	val []XPendingExt
***REMOVED***

var _ Cmder = (*XPendingExtCmd)(nil)

func NewXPendingExtCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XPendingExtCmd ***REMOVED***
	return &XPendingExtCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XPendingExtCmd) SetVal(val []XPendingExt) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XPendingExtCmd) Val() []XPendingExt ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XPendingExtCmd) Result() ([]XPendingExt, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XPendingExtCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XPendingExtCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]XPendingExt, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if err = rd.ReadFixedArrayLen(4); err != nil ***REMOVED***
			return err
		***REMOVED***

		if cmd.val[i].ID, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***

		if cmd.val[i].Consumer, err = rd.ReadString(); err != nil && err != Nil ***REMOVED***
			return err
		***REMOVED***

		idle, err := rd.ReadInt()
		if err != nil && err != Nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[i].Idle = time.Duration(idle) * time.Millisecond

		if cmd.val[i].RetryCount, err = rd.ReadInt(); err != nil && err != Nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XAutoClaimCmd struct ***REMOVED***
	baseCmd

	start string
	val   []XMessage
***REMOVED***

var _ Cmder = (*XAutoClaimCmd)(nil)

func NewXAutoClaimCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XAutoClaimCmd ***REMOVED***
	return &XAutoClaimCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XAutoClaimCmd) SetVal(val []XMessage, start string) ***REMOVED***
	cmd.val = val
	cmd.start = start
***REMOVED***

func (cmd *XAutoClaimCmd) Val() (messages []XMessage, start string) ***REMOVED***
	return cmd.val, cmd.start
***REMOVED***

func (cmd *XAutoClaimCmd) Result() (messages []XMessage, start string, err error) ***REMOVED***
	return cmd.val, cmd.start, cmd.err
***REMOVED***

func (cmd *XAutoClaimCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XAutoClaimCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch n ***REMOVED***
	case 2, // Redis 6
		3: // Redis 7:
		// ok
	default:
		return fmt.Errorf("redis: got %d elements in XAutoClaim reply, wanted 2/3", n)
	***REMOVED***

	cmd.start, err = rd.ReadString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val, err = readXMessageSlice(rd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n >= 3 ***REMOVED***
		if err := rd.DiscardNext(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XAutoClaimJustIDCmd struct ***REMOVED***
	baseCmd

	start string
	val   []string
***REMOVED***

var _ Cmder = (*XAutoClaimJustIDCmd)(nil)

func NewXAutoClaimJustIDCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XAutoClaimJustIDCmd ***REMOVED***
	return &XAutoClaimJustIDCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XAutoClaimJustIDCmd) SetVal(val []string, start string) ***REMOVED***
	cmd.val = val
	cmd.start = start
***REMOVED***

func (cmd *XAutoClaimJustIDCmd) Val() (ids []string, start string) ***REMOVED***
	return cmd.val, cmd.start
***REMOVED***

func (cmd *XAutoClaimJustIDCmd) Result() (ids []string, start string, err error) ***REMOVED***
	return cmd.val, cmd.start, cmd.err
***REMOVED***

func (cmd *XAutoClaimJustIDCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XAutoClaimJustIDCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch n ***REMOVED***
	case 2, // Redis 6
		3: // Redis 7:
		// ok
	default:
		return fmt.Errorf("redis: got %d elements in XAutoClaimJustID reply, wanted 2/3", n)
	***REMOVED***

	cmd.start, err = rd.ReadString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nn, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make([]string, nn)
	for i := 0; i < nn; i++ ***REMOVED***
		cmd.val[i], err = rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if n >= 3 ***REMOVED***
		if err := rd.DiscardNext(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XInfoConsumersCmd struct ***REMOVED***
	baseCmd
	val []XInfoConsumer
***REMOVED***

type XInfoConsumer struct ***REMOVED***
	Name    string
	Pending int64
	Idle    time.Duration
***REMOVED***

var _ Cmder = (*XInfoConsumersCmd)(nil)

func NewXInfoConsumersCmd(ctx context.Context, stream string, group string) *XInfoConsumersCmd ***REMOVED***
	return &XInfoConsumersCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: []interface***REMOVED******REMOVED******REMOVED***"xinfo", "consumers", stream, group***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XInfoConsumersCmd) SetVal(val []XInfoConsumer) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XInfoConsumersCmd) Val() []XInfoConsumer ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XInfoConsumersCmd) Result() ([]XInfoConsumer, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XInfoConsumersCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XInfoConsumersCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]XInfoConsumer, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if err = rd.ReadFixedMapLen(3); err != nil ***REMOVED***
			return err
		***REMOVED***

		var key string
		for f := 0; f < 3; f++ ***REMOVED***
			key, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			switch key ***REMOVED***
			case "name":
				cmd.val[i].Name, err = rd.ReadString()
			case "pending":
				cmd.val[i].Pending, err = rd.ReadInt()
			case "idle":
				var idle int64
				idle, err = rd.ReadInt()
				cmd.val[i].Idle = time.Duration(idle) * time.Millisecond
			default:
				return fmt.Errorf("redis: unexpected content %s in XINFO CONSUMERS reply", key)
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XInfoGroupsCmd struct ***REMOVED***
	baseCmd
	val []XInfoGroup
***REMOVED***

type XInfoGroup struct ***REMOVED***
	Name            string
	Consumers       int64
	Pending         int64
	LastDeliveredID string
	EntriesRead     int64
	Lag             int64
***REMOVED***

var _ Cmder = (*XInfoGroupsCmd)(nil)

func NewXInfoGroupsCmd(ctx context.Context, stream string) *XInfoGroupsCmd ***REMOVED***
	return &XInfoGroupsCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: []interface***REMOVED******REMOVED******REMOVED***"xinfo", "groups", stream***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XInfoGroupsCmd) SetVal(val []XInfoGroup) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XInfoGroupsCmd) Val() []XInfoGroup ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XInfoGroupsCmd) Result() ([]XInfoGroup, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XInfoGroupsCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XInfoGroupsCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]XInfoGroup, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		group := &cmd.val[i]

		nn, err := rd.ReadMapLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var key string
		for j := 0; j < nn; j++ ***REMOVED***
			key, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			switch key ***REMOVED***
			case "name":
				group.Name, err = rd.ReadString()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			case "consumers":
				group.Consumers, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			case "pending":
				group.Pending, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			case "last-delivered-id":
				group.LastDeliveredID, err = rd.ReadString()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			case "entries-read":
				group.EntriesRead, err = rd.ReadInt()
				if err != nil && err != Nil ***REMOVED***
					return err
				***REMOVED***
			case "lag":
				group.Lag, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			default:
				return fmt.Errorf("redis: unexpected key %q in XINFO GROUPS reply", key)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XInfoStreamCmd struct ***REMOVED***
	baseCmd
	val *XInfoStream
***REMOVED***

type XInfoStream struct ***REMOVED***
	Length               int64
	RadixTreeKeys        int64
	RadixTreeNodes       int64
	Groups               int64
	LastGeneratedID      string
	MaxDeletedEntryID    string
	EntriesAdded         int64
	FirstEntry           XMessage
	LastEntry            XMessage
	RecordedFirstEntryID string
***REMOVED***

var _ Cmder = (*XInfoStreamCmd)(nil)

func NewXInfoStreamCmd(ctx context.Context, stream string) *XInfoStreamCmd ***REMOVED***
	return &XInfoStreamCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: []interface***REMOVED******REMOVED******REMOVED***"xinfo", "stream", stream***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XInfoStreamCmd) SetVal(val *XInfoStream) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XInfoStreamCmd) Val() *XInfoStream ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XInfoStreamCmd) Result() (*XInfoStream, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XInfoStreamCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XInfoStreamCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = &XInfoStream***REMOVED******REMOVED***

	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch key ***REMOVED***
		case "length":
			cmd.val.Length, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "radix-tree-keys":
			cmd.val.RadixTreeKeys, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "radix-tree-nodes":
			cmd.val.RadixTreeNodes, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "groups":
			cmd.val.Groups, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "last-generated-id":
			cmd.val.LastGeneratedID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "max-deleted-entry-id":
			cmd.val.MaxDeletedEntryID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "entries-added":
			cmd.val.EntriesAdded, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "first-entry":
			cmd.val.FirstEntry, err = readXMessage(rd)
			if err != nil && err != Nil ***REMOVED***
				return err
			***REMOVED***
		case "last-entry":
			cmd.val.LastEntry, err = readXMessage(rd)
			if err != nil && err != Nil ***REMOVED***
				return err
			***REMOVED***
		case "recorded-first-entry-id":
			cmd.val.RecordedFirstEntryID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			return fmt.Errorf("redis: unexpected key %q in XINFO STREAM reply", key)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type XInfoStreamFullCmd struct ***REMOVED***
	baseCmd
	val *XInfoStreamFull
***REMOVED***

type XInfoStreamFull struct ***REMOVED***
	Length               int64
	RadixTreeKeys        int64
	RadixTreeNodes       int64
	LastGeneratedID      string
	MaxDeletedEntryID    string
	EntriesAdded         int64
	Entries              []XMessage
	Groups               []XInfoStreamGroup
	RecordedFirstEntryID string
***REMOVED***

type XInfoStreamGroup struct ***REMOVED***
	Name            string
	LastDeliveredID string
	EntriesRead     int64
	Lag             int64
	PelCount        int64
	Pending         []XInfoStreamGroupPending
	Consumers       []XInfoStreamConsumer
***REMOVED***

type XInfoStreamGroupPending struct ***REMOVED***
	ID            string
	Consumer      string
	DeliveryTime  time.Time
	DeliveryCount int64
***REMOVED***

type XInfoStreamConsumer struct ***REMOVED***
	Name     string
	SeenTime time.Time
	PelCount int64
	Pending  []XInfoStreamConsumerPending
***REMOVED***

type XInfoStreamConsumerPending struct ***REMOVED***
	ID            string
	DeliveryTime  time.Time
	DeliveryCount int64
***REMOVED***

var _ Cmder = (*XInfoStreamFullCmd)(nil)

func NewXInfoStreamFullCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *XInfoStreamFullCmd ***REMOVED***
	return &XInfoStreamFullCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *XInfoStreamFullCmd) SetVal(val *XInfoStreamFull) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *XInfoStreamFullCmd) Val() *XInfoStreamFull ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *XInfoStreamFullCmd) Result() (*XInfoStreamFull, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *XInfoStreamFullCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *XInfoStreamFullCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = &XInfoStreamFull***REMOVED******REMOVED***

	for i := 0; i < n; i++ ***REMOVED***
		key, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		switch key ***REMOVED***
		case "length":
			cmd.val.Length, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "radix-tree-keys":
			cmd.val.RadixTreeKeys, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "radix-tree-nodes":
			cmd.val.RadixTreeNodes, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "last-generated-id":
			cmd.val.LastGeneratedID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "entries-added":
			cmd.val.EntriesAdded, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "entries":
			cmd.val.Entries, err = readXMessageSlice(rd)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "groups":
			cmd.val.Groups, err = readStreamGroups(rd)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "max-deleted-entry-id":
			cmd.val.MaxDeletedEntryID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case "recorded-first-entry-id":
			cmd.val.RecordedFirstEntryID, err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			return fmt.Errorf("redis: unexpected key %q in XINFO STREAM FULL reply", key)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func readStreamGroups(rd *proto.Reader) ([]XInfoStreamGroup, error) ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	groups := make([]XInfoStreamGroup, 0, n)
	for i := 0; i < n; i++ ***REMOVED***
		nn, err := rd.ReadMapLen()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		group := XInfoStreamGroup***REMOVED******REMOVED***

		for j := 0; j < nn; j++ ***REMOVED***
			key, err := rd.ReadString()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			switch key ***REMOVED***
			case "name":
				group.Name, err = rd.ReadString()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "last-delivered-id":
				group.LastDeliveredID, err = rd.ReadString()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "entries-read":
				group.EntriesRead, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "lag":
				group.Lag, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "pel-count":
				group.PelCount, err = rd.ReadInt()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "pending":
				group.Pending, err = readXInfoStreamGroupPending(rd)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			case "consumers":
				group.Consumers, err = readXInfoStreamConsumers(rd)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			default:
				return nil, fmt.Errorf("redis: unexpected key %q in XINFO STREAM FULL reply", key)
			***REMOVED***
		***REMOVED***

		groups = append(groups, group)
	***REMOVED***

	return groups, nil
***REMOVED***

func readXInfoStreamGroupPending(rd *proto.Reader) ([]XInfoStreamGroupPending, error) ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pending := make([]XInfoStreamGroupPending, 0, n)

	for i := 0; i < n; i++ ***REMOVED***
		if err = rd.ReadFixedArrayLen(4); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		p := XInfoStreamGroupPending***REMOVED******REMOVED***

		p.ID, err = rd.ReadString()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		p.Consumer, err = rd.ReadString()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		delivery, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p.DeliveryTime = time.Unix(delivery/1000, delivery%1000*int64(time.Millisecond))

		p.DeliveryCount, err = rd.ReadInt()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		pending = append(pending, p)
	***REMOVED***

	return pending, nil
***REMOVED***

func readXInfoStreamConsumers(rd *proto.Reader) ([]XInfoStreamConsumer, error) ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	consumers := make([]XInfoStreamConsumer, 0, n)

	for i := 0; i < n; i++ ***REMOVED***
		if err = rd.ReadFixedMapLen(4); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		c := XInfoStreamConsumer***REMOVED******REMOVED***

		for f := 0; f < 4; f++ ***REMOVED***
			cKey, err := rd.ReadString()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			switch cKey ***REMOVED***
			case "name":
				c.Name, err = rd.ReadString()
			case "seen-time":
				seen, err := rd.ReadInt()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				c.SeenTime = time.Unix(seen/1000, seen%1000*int64(time.Millisecond))
			case "pel-count":
				c.PelCount, err = rd.ReadInt()
			case "pending":
				pendingNumber, err := rd.ReadArrayLen()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				c.Pending = make([]XInfoStreamConsumerPending, 0, pendingNumber)

				for pn := 0; pn < pendingNumber; pn++ ***REMOVED***
					if err = rd.ReadFixedArrayLen(3); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					p := XInfoStreamConsumerPending***REMOVED******REMOVED***

					p.ID, err = rd.ReadString()
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					delivery, err := rd.ReadInt()
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					p.DeliveryTime = time.Unix(delivery/1000, delivery%1000*int64(time.Millisecond))

					p.DeliveryCount, err = rd.ReadInt()
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					c.Pending = append(c.Pending, p)
				***REMOVED***
			default:
				return nil, fmt.Errorf("redis: unexpected content %s "+
					"in XINFO STREAM FULL reply", cKey)
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		consumers = append(consumers, c)
	***REMOVED***

	return consumers, nil
***REMOVED***

//------------------------------------------------------------------------------

type ZSliceCmd struct ***REMOVED***
	baseCmd

	val []Z
***REMOVED***

var _ Cmder = (*ZSliceCmd)(nil)

func NewZSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *ZSliceCmd ***REMOVED***
	return &ZSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *ZSliceCmd) SetVal(val []Z) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *ZSliceCmd) Val() []Z ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *ZSliceCmd) Result() ([]Z, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *ZSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *ZSliceCmd) readReply(rd *proto.Reader) error ***REMOVED*** // nolint:dupl
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the n is 0, can't continue reading.
	if n == 0 ***REMOVED***
		cmd.val = make([]Z, 0)
		return nil
	***REMOVED***

	typ, err := rd.PeekReplyType()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	array := typ == proto.RespArray

	if array ***REMOVED***
		cmd.val = make([]Z, n)
	***REMOVED*** else ***REMOVED***
		cmd.val = make([]Z, n/2)
	***REMOVED***

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		if array ***REMOVED***
			if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if cmd.val[i].Member, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***

		if cmd.val[i].Score, err = rd.ReadFloat(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type ZWithKeyCmd struct ***REMOVED***
	baseCmd

	val *ZWithKey
***REMOVED***

var _ Cmder = (*ZWithKeyCmd)(nil)

func NewZWithKeyCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *ZWithKeyCmd ***REMOVED***
	return &ZWithKeyCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *ZWithKeyCmd) SetVal(val *ZWithKey) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *ZWithKeyCmd) Val() *ZWithKey ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *ZWithKeyCmd) Result() (*ZWithKey, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *ZWithKeyCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *ZWithKeyCmd) readReply(rd *proto.Reader) (err error) ***REMOVED***
	if err = rd.ReadFixedArrayLen(3); err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = &ZWithKey***REMOVED******REMOVED***

	if cmd.val.Key, err = rd.ReadString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if cmd.val.Member, err = rd.ReadString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if cmd.val.Score, err = rd.ReadFloat(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type ScanCmd struct ***REMOVED***
	baseCmd

	page   []string
	cursor uint64

	process cmdable
***REMOVED***

var _ Cmder = (*ScanCmd)(nil)

func NewScanCmd(ctx context.Context, process cmdable, args ...interface***REMOVED******REMOVED***) *ScanCmd ***REMOVED***
	return &ScanCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
		process: process,
	***REMOVED***
***REMOVED***

func (cmd *ScanCmd) SetVal(page []string, cursor uint64) ***REMOVED***
	cmd.page = page
	cmd.cursor = cursor
***REMOVED***

func (cmd *ScanCmd) Val() (keys []string, cursor uint64) ***REMOVED***
	return cmd.page, cmd.cursor
***REMOVED***

func (cmd *ScanCmd) Result() (keys []string, cursor uint64, err error) ***REMOVED***
	return cmd.page, cmd.cursor, cmd.err
***REMOVED***

func (cmd *ScanCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.page)
***REMOVED***

func (cmd *ScanCmd) readReply(rd *proto.Reader) error ***REMOVED***
	if err := rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
		return err
	***REMOVED***

	cursor, err := rd.ReadInt()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.cursor = uint64(cursor)

	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.page = make([]string, n)

	for i := 0; i < len(cmd.page); i++ ***REMOVED***
		if cmd.page[i], err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Iterator creates a new ScanIterator.
func (cmd *ScanCmd) Iterator() *ScanIterator ***REMOVED***
	return &ScanIterator***REMOVED***
		cmd: cmd,
	***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

type ClusterNode struct ***REMOVED***
	ID                 string
	Addr               string
	NetworkingMetadata map[string]string
***REMOVED***

type ClusterSlot struct ***REMOVED***
	Start int
	End   int
	Nodes []ClusterNode
***REMOVED***

type ClusterSlotsCmd struct ***REMOVED***
	baseCmd

	val []ClusterSlot
***REMOVED***

var _ Cmder = (*ClusterSlotsCmd)(nil)

func NewClusterSlotsCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *ClusterSlotsCmd ***REMOVED***
	return &ClusterSlotsCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *ClusterSlotsCmd) SetVal(val []ClusterSlot) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *ClusterSlotsCmd) Val() []ClusterSlot ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *ClusterSlotsCmd) Result() ([]ClusterSlot, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *ClusterSlotsCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *ClusterSlotsCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]ClusterSlot, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		n, err = rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if n < 2 ***REMOVED***
			return fmt.Errorf("redis: got %d elements in cluster info, expected at least 2", n)
		***REMOVED***

		start, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		end, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// subtract start and end.
		nodes := make([]ClusterNode, n-2)

		for j := 0; j < len(nodes); j++ ***REMOVED***
			nn, err := rd.ReadArrayLen()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if nn < 2 || nn > 4 ***REMOVED***
				return fmt.Errorf("got %d elements in cluster info address, expected 2, 3, or 4", n)
			***REMOVED***

			ip, err := rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			port, err := rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			nodes[j].Addr = net.JoinHostPort(ip, port)

			if nn >= 3 ***REMOVED***
				id, err := rd.ReadString()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				nodes[j].ID = id
			***REMOVED***

			if nn >= 4 ***REMOVED***
				metadataLength, err := rd.ReadMapLen()
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				networkingMetadata := make(map[string]string, metadataLength)

				for i := 0; i < metadataLength; i++ ***REMOVED***
					key, err := rd.ReadString()
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					value, err := rd.ReadString()
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					networkingMetadata[key] = value
				***REMOVED***

				nodes[j].NetworkingMetadata = networkingMetadata
			***REMOVED***
		***REMOVED***

		cmd.val[i] = ClusterSlot***REMOVED***
			Start: int(start),
			End:   int(end),
			Nodes: nodes,
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

// GeoLocation is used with GeoAdd to add geospatial location.
type GeoLocation struct ***REMOVED***
	Name                      string
	Longitude, Latitude, Dist float64
	GeoHash                   int64
***REMOVED***

// GeoRadiusQuery is used with GeoRadius to query geospatial index.
type GeoRadiusQuery struct ***REMOVED***
	Radius float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        string
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort      string
	Store     string
	StoreDist string

	// WithCoord+WithDist+WithGeoHash
	withLen int
***REMOVED***

type GeoLocationCmd struct ***REMOVED***
	baseCmd

	q         *GeoRadiusQuery
	locations []GeoLocation
***REMOVED***

var _ Cmder = (*GeoLocationCmd)(nil)

func NewGeoLocationCmd(ctx context.Context, q *GeoRadiusQuery, args ...interface***REMOVED******REMOVED***) *GeoLocationCmd ***REMOVED***
	return &GeoLocationCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: geoLocationArgs(q, args...),
		***REMOVED***,
		q: q,
	***REMOVED***
***REMOVED***

func geoLocationArgs(q *GeoRadiusQuery, args ...interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	args = append(args, q.Radius)
	if q.Unit != "" ***REMOVED***
		args = append(args, q.Unit)
	***REMOVED*** else ***REMOVED***
		args = append(args, "km")
	***REMOVED***
	if q.WithCoord ***REMOVED***
		args = append(args, "withcoord")
		q.withLen++
	***REMOVED***
	if q.WithDist ***REMOVED***
		args = append(args, "withdist")
		q.withLen++
	***REMOVED***
	if q.WithGeoHash ***REMOVED***
		args = append(args, "withhash")
		q.withLen++
	***REMOVED***
	if q.Count > 0 ***REMOVED***
		args = append(args, "count", q.Count)
	***REMOVED***
	if q.Sort != "" ***REMOVED***
		args = append(args, q.Sort)
	***REMOVED***
	if q.Store != "" ***REMOVED***
		args = append(args, "store")
		args = append(args, q.Store)
	***REMOVED***
	if q.StoreDist != "" ***REMOVED***
		args = append(args, "storedist")
		args = append(args, q.StoreDist)
	***REMOVED***
	return args
***REMOVED***

func (cmd *GeoLocationCmd) SetVal(locations []GeoLocation) ***REMOVED***
	cmd.locations = locations
***REMOVED***

func (cmd *GeoLocationCmd) Val() []GeoLocation ***REMOVED***
	return cmd.locations
***REMOVED***

func (cmd *GeoLocationCmd) Result() ([]GeoLocation, error) ***REMOVED***
	return cmd.locations, cmd.err
***REMOVED***

func (cmd *GeoLocationCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.locations)
***REMOVED***

func (cmd *GeoLocationCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.locations = make([]GeoLocation, n)

	for i := 0; i < len(cmd.locations); i++ ***REMOVED***
		// only name
		if cmd.q.withLen == 0 ***REMOVED***
			if cmd.locations[i].Name, err = rd.ReadString(); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		// +name
		if err = rd.ReadFixedArrayLen(cmd.q.withLen + 1); err != nil ***REMOVED***
			return err
		***REMOVED***

		if cmd.locations[i].Name, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if cmd.q.WithDist ***REMOVED***
			if cmd.locations[i].Dist, err = rd.ReadFloat(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if cmd.q.WithGeoHash ***REMOVED***
			if cmd.locations[i].GeoHash, err = rd.ReadInt(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if cmd.q.WithCoord ***REMOVED***
			if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
				return err
			***REMOVED***
			if cmd.locations[i].Longitude, err = rd.ReadFloat(); err != nil ***REMOVED***
				return err
			***REMOVED***
			if cmd.locations[i].Latitude, err = rd.ReadFloat(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

// GeoSearchQuery is used for GEOSearch/GEOSearchStore command query.
type GeoSearchQuery struct ***REMOVED***
	Member string

	// Latitude and Longitude when using FromLonLat option.
	Longitude float64
	Latitude  float64

	// Distance and unit when using ByRadius option.
	// Can use m, km, ft, or mi. Default is km.
	Radius     float64
	RadiusUnit string

	// Height, width and unit when using ByBox option.
	// Can be m, km, ft, or mi. Default is km.
	BoxWidth  float64
	BoxHeight float64
	BoxUnit   string

	// Can be ASC or DESC. Default is no sort order.
	Sort     string
	Count    int
	CountAny bool
***REMOVED***

type GeoSearchLocationQuery struct ***REMOVED***
	GeoSearchQuery

	WithCoord bool
	WithDist  bool
	WithHash  bool
***REMOVED***

type GeoSearchStoreQuery struct ***REMOVED***
	GeoSearchQuery

	// When using the StoreDist option, the command stores the items in a
	// sorted set populated with their distance from the center of the circle or box,
	// as a floating-point number, in the same unit specified for that shape.
	StoreDist bool
***REMOVED***

func geoSearchLocationArgs(q *GeoSearchLocationQuery, args []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	args = geoSearchArgs(&q.GeoSearchQuery, args)

	if q.WithCoord ***REMOVED***
		args = append(args, "withcoord")
	***REMOVED***
	if q.WithDist ***REMOVED***
		args = append(args, "withdist")
	***REMOVED***
	if q.WithHash ***REMOVED***
		args = append(args, "withhash")
	***REMOVED***

	return args
***REMOVED***

func geoSearchArgs(q *GeoSearchQuery, args []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	if q.Member != "" ***REMOVED***
		args = append(args, "frommember", q.Member)
	***REMOVED*** else ***REMOVED***
		args = append(args, "fromlonlat", q.Longitude, q.Latitude)
	***REMOVED***

	if q.Radius > 0 ***REMOVED***
		if q.RadiusUnit == "" ***REMOVED***
			q.RadiusUnit = "km"
		***REMOVED***
		args = append(args, "byradius", q.Radius, q.RadiusUnit)
	***REMOVED*** else ***REMOVED***
		if q.BoxUnit == "" ***REMOVED***
			q.BoxUnit = "km"
		***REMOVED***
		args = append(args, "bybox", q.BoxWidth, q.BoxHeight, q.BoxUnit)
	***REMOVED***

	if q.Sort != "" ***REMOVED***
		args = append(args, q.Sort)
	***REMOVED***

	if q.Count > 0 ***REMOVED***
		args = append(args, "count", q.Count)
		if q.CountAny ***REMOVED***
			args = append(args, "any")
		***REMOVED***
	***REMOVED***

	return args
***REMOVED***

type GeoSearchLocationCmd struct ***REMOVED***
	baseCmd

	opt *GeoSearchLocationQuery
	val []GeoLocation
***REMOVED***

var _ Cmder = (*GeoSearchLocationCmd)(nil)

func NewGeoSearchLocationCmd(
	ctx context.Context, opt *GeoSearchLocationQuery, args ...interface***REMOVED******REMOVED***,
) *GeoSearchLocationCmd ***REMOVED***
	return &GeoSearchLocationCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
		opt: opt,
	***REMOVED***
***REMOVED***

func (cmd *GeoSearchLocationCmd) SetVal(val []GeoLocation) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *GeoSearchLocationCmd) Val() []GeoLocation ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *GeoSearchLocationCmd) Result() ([]GeoLocation, error) ***REMOVED***
	return cmd.val, cmd.err
***REMOVED***

func (cmd *GeoSearchLocationCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *GeoSearchLocationCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make([]GeoLocation, n)
	for i := 0; i < n; i++ ***REMOVED***
		_, err = rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var loc GeoLocation

		loc.Name, err = rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if cmd.opt.WithDist ***REMOVED***
			loc.Dist, err = rd.ReadFloat()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if cmd.opt.WithHash ***REMOVED***
			loc.GeoHash, err = rd.ReadInt()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if cmd.opt.WithCoord ***REMOVED***
			if err = rd.ReadFixedArrayLen(2); err != nil ***REMOVED***
				return err
			***REMOVED***
			loc.Longitude, err = rd.ReadFloat()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			loc.Latitude, err = rd.ReadFloat()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		cmd.val[i] = loc
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type GeoPos struct ***REMOVED***
	Longitude, Latitude float64
***REMOVED***

type GeoPosCmd struct ***REMOVED***
	baseCmd

	val []*GeoPos
***REMOVED***

var _ Cmder = (*GeoPosCmd)(nil)

func NewGeoPosCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *GeoPosCmd ***REMOVED***
	return &GeoPosCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *GeoPosCmd) SetVal(val []*GeoPos) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *GeoPosCmd) Val() []*GeoPos ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *GeoPosCmd) Result() ([]*GeoPos, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *GeoPosCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *GeoPosCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]*GeoPos, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		err = rd.ReadFixedArrayLen(2)
		if err != nil ***REMOVED***
			if err == Nil ***REMOVED***
				cmd.val[i] = nil
				continue
			***REMOVED***
			return err
		***REMOVED***

		longitude, err := rd.ReadFloat()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		latitude, err := rd.ReadFloat()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		cmd.val[i] = &GeoPos***REMOVED***
			Longitude: longitude,
			Latitude:  latitude,
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type CommandInfo struct ***REMOVED***
	Name        string
	Arity       int8
	Flags       []string
	ACLFlags    []string
	FirstKeyPos int8
	LastKeyPos  int8
	StepCount   int8
	ReadOnly    bool
***REMOVED***

type CommandsInfoCmd struct ***REMOVED***
	baseCmd

	val map[string]*CommandInfo
***REMOVED***

var _ Cmder = (*CommandsInfoCmd)(nil)

func NewCommandsInfoCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *CommandsInfoCmd ***REMOVED***
	return &CommandsInfoCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *CommandsInfoCmd) SetVal(val map[string]*CommandInfo) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *CommandsInfoCmd) Val() map[string]*CommandInfo ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *CommandsInfoCmd) Result() (map[string]*CommandInfo, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *CommandsInfoCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *CommandsInfoCmd) readReply(rd *proto.Reader) error ***REMOVED***
	const numArgRedis5 = 6
	const numArgRedis6 = 7
	const numArgRedis7 = 10

	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make(map[string]*CommandInfo, n)

	for i := 0; i < n; i++ ***REMOVED***
		nn, err := rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		switch nn ***REMOVED***
		case numArgRedis5, numArgRedis6, numArgRedis7:
			// ok
		default:
			return fmt.Errorf("redis: got %d elements in COMMAND reply, wanted 6/7/10", nn)
		***REMOVED***

		cmdInfo := &CommandInfo***REMOVED******REMOVED***
		if cmdInfo.Name, err = rd.ReadString(); err != nil ***REMOVED***
			return err
		***REMOVED***

		arity, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdInfo.Arity = int8(arity)

		flagLen, err := rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdInfo.Flags = make([]string, flagLen)
		for f := 0; f < len(cmdInfo.Flags); f++ ***REMOVED***
			switch s, err := rd.ReadString(); ***REMOVED***
			case err == Nil:
				cmdInfo.Flags[f] = ""
			case err != nil:
				return err
			default:
				if !cmdInfo.ReadOnly && s == "readonly" ***REMOVED***
					cmdInfo.ReadOnly = true
				***REMOVED***
				cmdInfo.Flags[f] = s
			***REMOVED***
		***REMOVED***

		firstKeyPos, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdInfo.FirstKeyPos = int8(firstKeyPos)

		lastKeyPos, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdInfo.LastKeyPos = int8(lastKeyPos)

		stepCount, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdInfo.StepCount = int8(stepCount)

		if nn >= numArgRedis6 ***REMOVED***
			aclFlagLen, err := rd.ReadArrayLen()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			cmdInfo.ACLFlags = make([]string, aclFlagLen)
			for f := 0; f < len(cmdInfo.ACLFlags); f++ ***REMOVED***
				switch s, err := rd.ReadString(); ***REMOVED***
				case err == Nil:
					cmdInfo.ACLFlags[f] = ""
				case err != nil:
					return err
				default:
					cmdInfo.ACLFlags[f] = s
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if nn >= numArgRedis7 ***REMOVED***
			if err := rd.DiscardNext(); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := rd.DiscardNext(); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := rd.DiscardNext(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		cmd.val[cmdInfo.Name] = cmdInfo
	***REMOVED***

	return nil
***REMOVED***

//------------------------------------------------------------------------------

type cmdsInfoCache struct ***REMOVED***
	fn func(ctx context.Context) (map[string]*CommandInfo, error)

	once internal.Once
	cmds map[string]*CommandInfo
***REMOVED***

func newCmdsInfoCache(fn func(ctx context.Context) (map[string]*CommandInfo, error)) *cmdsInfoCache ***REMOVED***
	return &cmdsInfoCache***REMOVED***
		fn: fn,
	***REMOVED***
***REMOVED***

func (c *cmdsInfoCache) Get(ctx context.Context) (map[string]*CommandInfo, error) ***REMOVED***
	err := c.once.Do(func() error ***REMOVED***
		cmds, err := c.fn(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Extensions have cmd names in upper case. Convert them to lower case.
		for k, v := range cmds ***REMOVED***
			lower := internal.ToLower(k)
			if lower != k ***REMOVED***
				cmds[lower] = v
			***REMOVED***
		***REMOVED***

		c.cmds = cmds
		return nil
	***REMOVED***)
	return c.cmds, err
***REMOVED***

//------------------------------------------------------------------------------

type SlowLog struct ***REMOVED***
	ID       int64
	Time     time.Time
	Duration time.Duration
	Args     []string
	// These are also optional fields emitted only by Redis 4.0 or greater:
	// https://redis.io/commands/slowlog#output-format
	ClientAddr string
	ClientName string
***REMOVED***

type SlowLogCmd struct ***REMOVED***
	baseCmd

	val []SlowLog
***REMOVED***

var _ Cmder = (*SlowLogCmd)(nil)

func NewSlowLogCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *SlowLogCmd ***REMOVED***
	return &SlowLogCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *SlowLogCmd) SetVal(val []SlowLog) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *SlowLogCmd) Val() []SlowLog ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *SlowLogCmd) Result() ([]SlowLog, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *SlowLogCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *SlowLogCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	cmd.val = make([]SlowLog, n)

	for i := 0; i < len(cmd.val); i++ ***REMOVED***
		nn, err := rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if nn < 4 ***REMOVED***
			return fmt.Errorf("redis: got %d elements in slowlog get, expected at least 4", nn)
		***REMOVED***

		if cmd.val[i].ID, err = rd.ReadInt(); err != nil ***REMOVED***
			return err
		***REMOVED***

		createdAt, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[i].Time = time.Unix(createdAt, 0)

		costs, err := rd.ReadInt()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[i].Duration = time.Duration(costs) * time.Microsecond

		cmdLen, err := rd.ReadArrayLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if cmdLen < 1 ***REMOVED***
			return fmt.Errorf("redis: got %d elements commands reply in slowlog get, expected at least 1", cmdLen)
		***REMOVED***

		cmd.val[i].Args = make([]string, cmdLen)
		for f := 0; f < len(cmd.val[i].Args); f++ ***REMOVED***
			cmd.val[i].Args[f], err = rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if nn >= 5 ***REMOVED***
			if cmd.val[i].ClientAddr, err = rd.ReadString(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if nn >= 6 ***REMOVED***
			if cmd.val[i].ClientName, err = rd.ReadString(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

//-----------------------------------------------------------------------

type MapStringInterfaceCmd struct ***REMOVED***
	baseCmd

	val map[string]interface***REMOVED******REMOVED***
***REMOVED***

var _ Cmder = (*MapStringInterfaceCmd)(nil)

func NewMapStringInterfaceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *MapStringInterfaceCmd ***REMOVED***
	return &MapStringInterfaceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *MapStringInterfaceCmd) SetVal(val map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *MapStringInterfaceCmd) Val() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *MapStringInterfaceCmd) Result() (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *MapStringInterfaceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *MapStringInterfaceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make(map[string]interface***REMOVED******REMOVED***, n)
	for i := 0; i < n; i++ ***REMOVED***
		k, err := rd.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v, err := rd.ReadReply()
		if err != nil ***REMOVED***
			if err == Nil ***REMOVED***
				cmd.val[k] = Nil
				continue
			***REMOVED***
			if err, ok := err.(proto.RedisError); ok ***REMOVED***
				cmd.val[k] = err
				continue
			***REMOVED***
			return err
		***REMOVED***
		cmd.val[k] = v
	***REMOVED***
	return nil
***REMOVED***

//-----------------------------------------------------------------------

type MapStringStringSliceCmd struct ***REMOVED***
	baseCmd

	val []map[string]string
***REMOVED***

var _ Cmder = (*MapStringStringSliceCmd)(nil)

func NewMapStringStringSliceCmd(ctx context.Context, args ...interface***REMOVED******REMOVED***) *MapStringStringSliceCmd ***REMOVED***
	return &MapStringStringSliceCmd***REMOVED***
		baseCmd: baseCmd***REMOVED***
			ctx:  ctx,
			args: args,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (cmd *MapStringStringSliceCmd) SetVal(val []map[string]string) ***REMOVED***
	cmd.val = val
***REMOVED***

func (cmd *MapStringStringSliceCmd) Val() []map[string]string ***REMOVED***
	return cmd.val
***REMOVED***

func (cmd *MapStringStringSliceCmd) Result() ([]map[string]string, error) ***REMOVED***
	return cmd.Val(), cmd.Err()
***REMOVED***

func (cmd *MapStringStringSliceCmd) String() string ***REMOVED***
	return cmdString(cmd, cmd.val)
***REMOVED***

func (cmd *MapStringStringSliceCmd) readReply(rd *proto.Reader) error ***REMOVED***
	n, err := rd.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cmd.val = make([]map[string]string, n)
	for i := 0; i < n; i++ ***REMOVED***
		nn, err := rd.ReadMapLen()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmd.val[i] = make(map[string]string, nn)
		for f := 0; f < nn; f++ ***REMOVED***
			k, err := rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			v, err := rd.ReadString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			cmd.val[i][k] = v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
