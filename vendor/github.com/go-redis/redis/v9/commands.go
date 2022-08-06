package redis

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/go-redis/redis/v9/internal"
)

// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
// For example:
//
//    rdb.Set(ctx, key, value, redis.KeepTTL)
const KeepTTL = -1

func usePrecise(dur time.Duration) bool ***REMOVED***
	return dur < time.Second || dur%time.Second != 0
***REMOVED***

func formatMs(ctx context.Context, dur time.Duration) int64 ***REMOVED***
	if dur > 0 && dur < time.Millisecond ***REMOVED***
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1ms",
			dur, time.Millisecond,
		)
		return 1
	***REMOVED***
	return int64(dur / time.Millisecond)
***REMOVED***

func formatSec(ctx context.Context, dur time.Duration) int64 ***REMOVED***
	if dur > 0 && dur < time.Second ***REMOVED***
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1s",
			dur, time.Second,
		)
		return 1
	***REMOVED***
	return int64(dur / time.Second)
***REMOVED***

func appendArgs(dst, src []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	if len(src) == 1 ***REMOVED***
		return appendArg(dst, src[0])
	***REMOVED***

	dst = append(dst, src...)
	return dst
***REMOVED***

func appendArg(dst []interface***REMOVED******REMOVED***, arg interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	switch arg := arg.(type) ***REMOVED***
	case []string:
		for _, s := range arg ***REMOVED***
			dst = append(dst, s)
		***REMOVED***
		return dst
	case []interface***REMOVED******REMOVED***:
		dst = append(dst, arg...)
		return dst
	case map[string]interface***REMOVED******REMOVED***:
		for k, v := range arg ***REMOVED***
			dst = append(dst, k, v)
		***REMOVED***
		return dst
	case map[string]string:
		for k, v := range arg ***REMOVED***
			dst = append(dst, k, v)
		***REMOVED***
		return dst
	default:
		return append(dst, arg)
	***REMOVED***
***REMOVED***

type Cmdable interface ***REMOVED***
	Pipeline() Pipeliner
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)

	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner

	Command(ctx context.Context) *CommandsInfoCmd
	ClientGetName(ctx context.Context) *StringCmd
	Echo(ctx context.Context, message interface***REMOVED******REMOVED***) *StringCmd
	Ping(ctx context.Context) *StatusCmd
	Quit(ctx context.Context) *StatusCmd
	Del(ctx context.Context, keys ...string) *IntCmd
	Unlink(ctx context.Context, keys ...string) *IntCmd
	Dump(ctx context.Context, key string) *StringCmd
	Exists(ctx context.Context, keys ...string) *IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	ExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd
	ExpireNX(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	ExpireXX(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	ExpireGT(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	ExpireLT(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	Keys(ctx context.Context, pattern string) *StringSliceCmd
	Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *StatusCmd
	Move(ctx context.Context, key string, db int) *BoolCmd
	ObjectRefCount(ctx context.Context, key string) *IntCmd
	ObjectEncoding(ctx context.Context, key string) *StringCmd
	ObjectIdleTime(ctx context.Context, key string) *DurationCmd
	Persist(ctx context.Context, key string) *BoolCmd
	PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	PExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd
	PTTL(ctx context.Context, key string) *DurationCmd
	RandomKey(ctx context.Context) *StringCmd
	Rename(ctx context.Context, key, newkey string) *StatusCmd
	RenameNX(ctx context.Context, key, newkey string) *BoolCmd
	Restore(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd
	RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd
	Sort(ctx context.Context, key string, sort *Sort) *StringSliceCmd
	SortStore(ctx context.Context, key, store string, sort *Sort) *IntCmd
	SortInterfaces(ctx context.Context, key string, sort *Sort) *SliceCmd
	Touch(ctx context.Context, keys ...string) *IntCmd
	TTL(ctx context.Context, key string) *DurationCmd
	Type(ctx context.Context, key string) *StatusCmd
	Append(ctx context.Context, key, value string) *IntCmd
	Decr(ctx context.Context, key string) *IntCmd
	DecrBy(ctx context.Context, key string, decrement int64) *IntCmd
	Get(ctx context.Context, key string) *StringCmd
	GetRange(ctx context.Context, key string, start, end int64) *StringCmd
	GetSet(ctx context.Context, key string, value interface***REMOVED******REMOVED***) *StringCmd
	GetEx(ctx context.Context, key string, expiration time.Duration) *StringCmd
	GetDel(ctx context.Context, key string) *StringCmd
	Incr(ctx context.Context, key string) *IntCmd
	IncrBy(ctx context.Context, key string, value int64) *IntCmd
	IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd
	MGet(ctx context.Context, keys ...string) *SliceCmd
	MSet(ctx context.Context, values ...interface***REMOVED******REMOVED***) *StatusCmd
	MSetNX(ctx context.Context, values ...interface***REMOVED******REMOVED***) *BoolCmd
	Set(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *StatusCmd
	SetArgs(ctx context.Context, key string, value interface***REMOVED******REMOVED***, a SetArgs) *StatusCmd
	SetEx(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *StatusCmd
	SetNX(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *BoolCmd
	SetXX(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *BoolCmd
	SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd
	StrLen(ctx context.Context, key string) *IntCmd
	Copy(ctx context.Context, sourceKey string, destKey string, db int, replace bool) *IntCmd

	GetBit(ctx context.Context, key string, offset int64) *IntCmd
	SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd
	BitCount(ctx context.Context, key string, bitCount *BitCount) *IntCmd
	BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpNot(ctx context.Context, destKey string, key string) *IntCmd
	BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmd
	BitField(ctx context.Context, key string, args ...interface***REMOVED******REMOVED***) *IntSliceCmd

	Scan(ctx context.Context, cursor uint64, match string, count int64) *ScanCmd
	ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *ScanCmd
	SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd

	HDel(ctx context.Context, key string, fields ...string) *IntCmd
	HExists(ctx context.Context, key, field string) *BoolCmd
	HGet(ctx context.Context, key, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *MapStringStringCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd
	HKeys(ctx context.Context, key string) *StringSliceCmd
	HLen(ctx context.Context, key string) *IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
	HSet(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd
	HMSet(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *BoolCmd
	HSetNX(ctx context.Context, key, field string, value interface***REMOVED******REMOVED***) *BoolCmd
	HVals(ctx context.Context, key string) *StringSliceCmd
	HRandField(ctx context.Context, key string, count int) *StringSliceCmd
	HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd

	BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmd
	LIndex(ctx context.Context, key string, index int64) *StringCmd
	LInsert(ctx context.Context, key, op string, pivot, value interface***REMOVED******REMOVED***) *IntCmd
	LInsertBefore(ctx context.Context, key string, pivot, value interface***REMOVED******REMOVED***) *IntCmd
	LInsertAfter(ctx context.Context, key string, pivot, value interface***REMOVED******REMOVED***) *IntCmd
	LLen(ctx context.Context, key string) *IntCmd
	LPop(ctx context.Context, key string) *StringCmd
	LPopCount(ctx context.Context, key string, count int) *StringSliceCmd
	LPos(ctx context.Context, key string, value string, args LPosArgs) *IntCmd
	LPosCount(ctx context.Context, key string, value string, count int64, args LPosArgs) *IntSliceCmd
	LPush(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd
	LPushX(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	LRem(ctx context.Context, key string, count int64, value interface***REMOVED******REMOVED***) *IntCmd
	LSet(ctx context.Context, key string, index int64, value interface***REMOVED******REMOVED***) *StatusCmd
	LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd
	RPop(ctx context.Context, key string) *StringCmd
	RPopCount(ctx context.Context, key string, count int) *StringSliceCmd
	RPopLPush(ctx context.Context, source, destination string) *StringCmd
	RPush(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd
	RPushX(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd
	LMove(ctx context.Context, source, destination, srcpos, destpos string) *StringCmd
	BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *StringCmd

	SAdd(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd
	SCard(ctx context.Context, key string) *IntCmd
	SDiff(ctx context.Context, keys ...string) *StringSliceCmd
	SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd
	SInter(ctx context.Context, keys ...string) *StringSliceCmd
	SInterStore(ctx context.Context, destination string, keys ...string) *IntCmd
	SIsMember(ctx context.Context, key string, member interface***REMOVED******REMOVED***) *BoolCmd
	SMIsMember(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *BoolSliceCmd
	SMembers(ctx context.Context, key string) *StringSliceCmd
	SMembersMap(ctx context.Context, key string) *StringStructMapCmd
	SMove(ctx context.Context, source, destination string, member interface***REMOVED******REMOVED***) *BoolCmd
	SPop(ctx context.Context, key string) *StringCmd
	SPopN(ctx context.Context, key string, count int64) *StringSliceCmd
	SRandMember(ctx context.Context, key string) *StringCmd
	SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd
	SRem(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd
	SUnion(ctx context.Context, keys ...string) *StringSliceCmd
	SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmd

	XAdd(ctx context.Context, a *XAddArgs) *StringCmd
	XDel(ctx context.Context, stream string, ids ...string) *IntCmd
	XLen(ctx context.Context, stream string) *IntCmd
	XRange(ctx context.Context, stream, start, stop string) *XMessageSliceCmd
	XRangeN(ctx context.Context, stream, start, stop string, count int64) *XMessageSliceCmd
	XRevRange(ctx context.Context, stream string, start, stop string) *XMessageSliceCmd
	XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *XMessageSliceCmd
	XRead(ctx context.Context, a *XReadArgs) *XStreamSliceCmd
	XReadStreams(ctx context.Context, streams ...string) *XStreamSliceCmd
	XGroupCreate(ctx context.Context, stream, group, start string) *StatusCmd
	XGroupCreateMkStream(ctx context.Context, stream, group, start string) *StatusCmd
	XGroupSetID(ctx context.Context, stream, group, start string) *StatusCmd
	XGroupDestroy(ctx context.Context, stream, group string) *IntCmd
	XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *IntCmd
	XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *IntCmd
	XReadGroup(ctx context.Context, a *XReadGroupArgs) *XStreamSliceCmd
	XAck(ctx context.Context, stream, group string, ids ...string) *IntCmd
	XPending(ctx context.Context, stream, group string) *XPendingCmd
	XPendingExt(ctx context.Context, a *XPendingExtArgs) *XPendingExtCmd
	XClaim(ctx context.Context, a *XClaimArgs) *XMessageSliceCmd
	XClaimJustID(ctx context.Context, a *XClaimArgs) *StringSliceCmd
	XAutoClaim(ctx context.Context, a *XAutoClaimArgs) *XAutoClaimCmd
	XAutoClaimJustID(ctx context.Context, a *XAutoClaimArgs) *XAutoClaimJustIDCmd
	XTrimMaxLen(ctx context.Context, key string, maxLen int64) *IntCmd
	XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *IntCmd
	XTrimMinID(ctx context.Context, key string, minID string) *IntCmd
	XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *IntCmd
	XInfoGroups(ctx context.Context, key string) *XInfoGroupsCmd
	XInfoStream(ctx context.Context, key string) *XInfoStreamCmd
	XInfoStreamFull(ctx context.Context, key string, count int) *XInfoStreamFullCmd
	XInfoConsumers(ctx context.Context, key string, group string) *XInfoConsumersCmd

	BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd
	BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd

	ZAdd(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd
	ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd
	ZCard(ctx context.Context, key string) *IntCmd
	ZCount(ctx context.Context, key, min, max string) *IntCmd
	ZLexCount(ctx context.Context, key, min, max string) *IntCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd
	ZInter(ctx context.Context, store *ZStore) *StringSliceCmd
	ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd
	ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd
	ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd
	ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd
	ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd
	ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd
	ZRank(ctx context.Context, key, member string) *IntCmd
	ZRem(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd
	ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRevRank(ctx context.Context, key, member string) *IntCmd
	ZScore(ctx context.Context, key, member string) *FloatCmd
	ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd
	ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd
	ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd
	ZUnion(ctx context.Context, store ZStore) *StringSliceCmd
	ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd
	ZDiff(ctx context.Context, keys ...string) *StringSliceCmd
	ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd
	ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd

	PFAdd(ctx context.Context, key string, els ...interface***REMOVED******REMOVED***) *IntCmd
	PFCount(ctx context.Context, keys ...string) *IntCmd
	PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd

	BgRewriteAOF(ctx context.Context) *StatusCmd
	BgSave(ctx context.Context) *StatusCmd
	ClientKill(ctx context.Context, ipPort string) *StatusCmd
	ClientKillByFilter(ctx context.Context, keys ...string) *IntCmd
	ClientList(ctx context.Context) *StringCmd
	ClientPause(ctx context.Context, dur time.Duration) *BoolCmd
	ClientUnpause(ctx context.Context) *BoolCmd
	ClientID(ctx context.Context) *IntCmd
	ClientUnblock(ctx context.Context, id int64) *IntCmd
	ClientUnblockWithError(ctx context.Context, id int64) *IntCmd
	ConfigGet(ctx context.Context, parameter string) *MapStringStringCmd
	ConfigResetStat(ctx context.Context) *StatusCmd
	ConfigSet(ctx context.Context, parameter, value string) *StatusCmd
	ConfigRewrite(ctx context.Context) *StatusCmd
	DBSize(ctx context.Context) *IntCmd
	FlushAll(ctx context.Context) *StatusCmd
	FlushAllAsync(ctx context.Context) *StatusCmd
	FlushDB(ctx context.Context) *StatusCmd
	FlushDBAsync(ctx context.Context) *StatusCmd
	Info(ctx context.Context, section ...string) *StringCmd
	LastSave(ctx context.Context) *IntCmd
	Save(ctx context.Context) *StatusCmd
	Shutdown(ctx context.Context) *StatusCmd
	ShutdownSave(ctx context.Context) *StatusCmd
	ShutdownNoSave(ctx context.Context) *StatusCmd
	SlaveOf(ctx context.Context, host, port string) *StatusCmd
	SlowLogGet(ctx context.Context, num int64) *SlowLogCmd
	Time(ctx context.Context) *TimeCmd
	DebugObject(ctx context.Context, key string) *StringCmd
	ReadOnly(ctx context.Context) *StatusCmd
	ReadWrite(ctx context.Context) *StatusCmd
	MemoryUsage(ctx context.Context, key string, samples ...int) *IntCmd

	Eval(ctx context.Context, script string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd
	ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd
	ScriptFlush(ctx context.Context) *StatusCmd
	ScriptKill(ctx context.Context) *StatusCmd
	ScriptLoad(ctx context.Context, script string) *StringCmd

	Publish(ctx context.Context, channel string, message interface***REMOVED******REMOVED***) *IntCmd
	PubSubChannels(ctx context.Context, pattern string) *StringSliceCmd
	PubSubNumSub(ctx context.Context, channels ...string) *StringIntMapCmd
	PubSubNumPat(ctx context.Context) *IntCmd

	ClusterSlots(ctx context.Context) *ClusterSlotsCmd
	ClusterNodes(ctx context.Context) *StringCmd
	ClusterMeet(ctx context.Context, host, port string) *StatusCmd
	ClusterForget(ctx context.Context, nodeID string) *StatusCmd
	ClusterReplicate(ctx context.Context, nodeID string) *StatusCmd
	ClusterResetSoft(ctx context.Context) *StatusCmd
	ClusterResetHard(ctx context.Context) *StatusCmd
	ClusterInfo(ctx context.Context) *StringCmd
	ClusterKeySlot(ctx context.Context, key string) *IntCmd
	ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *StringSliceCmd
	ClusterCountFailureReports(ctx context.Context, nodeID string) *IntCmd
	ClusterCountKeysInSlot(ctx context.Context, slot int) *IntCmd
	ClusterDelSlots(ctx context.Context, slots ...int) *StatusCmd
	ClusterDelSlotsRange(ctx context.Context, min, max int) *StatusCmd
	ClusterSaveConfig(ctx context.Context) *StatusCmd
	ClusterSlaves(ctx context.Context, nodeID string) *StringSliceCmd
	ClusterFailover(ctx context.Context) *StatusCmd
	ClusterAddSlots(ctx context.Context, slots ...int) *StatusCmd
	ClusterAddSlotsRange(ctx context.Context, min, max int) *StatusCmd

	GeoAdd(ctx context.Context, key string, geoLocation ...*GeoLocation) *IntCmd
	GeoPos(ctx context.Context, key string, members ...string) *GeoPosCmd
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *GeoLocationCmd
	GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *IntCmd
	GeoRadiusByMember(ctx context.Context, key, member string, query *GeoRadiusQuery) *GeoLocationCmd
	GeoRadiusByMemberStore(ctx context.Context, key, member string, query *GeoRadiusQuery) *IntCmd
	GeoSearch(ctx context.Context, key string, q *GeoSearchQuery) *StringSliceCmd
	GeoSearchLocation(ctx context.Context, key string, q *GeoSearchLocationQuery) *GeoSearchLocationCmd
	GeoSearchStore(ctx context.Context, key, store string, q *GeoSearchStoreQuery) *IntCmd
	GeoDist(ctx context.Context, key string, member1, member2, unit string) *FloatCmd
	GeoHash(ctx context.Context, key string, members ...string) *StringSliceCmd
***REMOVED***

type StatefulCmdable interface ***REMOVED***
	Cmdable
	Auth(ctx context.Context, password string) *StatusCmd
	AuthACL(ctx context.Context, username, password string) *StatusCmd
	Select(ctx context.Context, index int) *StatusCmd
	SwapDB(ctx context.Context, index1, index2 int) *StatusCmd
	ClientSetName(ctx context.Context, name string) *BoolCmd
	Hello(ctx context.Context, ver int, username, password, clientName string) *MapStringInterfaceCmd
***REMOVED***

var (
	_ Cmdable = (*Client)(nil)
	_ Cmdable = (*Tx)(nil)
	_ Cmdable = (*Ring)(nil)
	_ Cmdable = (*ClusterClient)(nil)
)

type cmdable func(ctx context.Context, cmd Cmder) error

type statefulCmdable func(ctx context.Context, cmd Cmder) error

//------------------------------------------------------------------------------

func (c statefulCmdable) Auth(ctx context.Context, password string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "auth", password)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// AuthACL Perform an AUTH command, using the given user and pass.
// Should be used to authenticate the current connection with one of the connections defined in the ACL list
// when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
func (c statefulCmdable) AuthACL(ctx context.Context, username, password string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "auth", username, password)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "wait", numSlaves, int(timeout/time.Millisecond))
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c statefulCmdable) Select(ctx context.Context, index int) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "select", index)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c statefulCmdable) SwapDB(ctx context.Context, index1, index2 int) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "swapdb", index1, index2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ClientSetName assigns a name to the connection.
func (c statefulCmdable) ClientSetName(ctx context.Context, name string) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "client", "setname", name)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// Hello Set the resp protocol used.
func (c statefulCmdable) Hello(ctx context.Context,
	ver int, username, password, clientName string) *MapStringInterfaceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 7)
	args = append(args, "hello", ver)
	if password != "" ***REMOVED***
		if username != "" ***REMOVED***
			args = append(args, "auth", username, password)
		***REMOVED*** else ***REMOVED***
			args = append(args, "auth", "default", password)
		***REMOVED***
	***REMOVED***
	if clientName != "" ***REMOVED***
		args = append(args, "setname", clientName)
	***REMOVED***
	cmd := NewMapStringInterfaceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) Command(ctx context.Context) *CommandsInfoCmd ***REMOVED***
	cmd := NewCommandsInfoCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ClientGetName returns the name of the connection.
func (c cmdable) ClientGetName(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "client", "getname")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Echo(ctx context.Context, message interface***REMOVED******REMOVED***) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "echo", message)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Ping(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "ping")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Quit(_ context.Context) *StatusCmd ***REMOVED***
	panic("not implemented")
***REMOVED***

func (c cmdable) Del(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "del"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Unlink(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "unlink"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Dump(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "dump", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Exists(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "exists"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	return c.expire(ctx, key, expiration, "")
***REMOVED***

func (c cmdable) ExpireNX(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	return c.expire(ctx, key, expiration, "NX")
***REMOVED***

func (c cmdable) ExpireXX(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	return c.expire(ctx, key, expiration, "XX")
***REMOVED***

func (c cmdable) ExpireGT(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	return c.expire(ctx, key, expiration, "GT")
***REMOVED***

func (c cmdable) ExpireLT(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	return c.expire(ctx, key, expiration, "LT")
***REMOVED***

func (c cmdable) expire(
	ctx context.Context, key string, expiration time.Duration, mode string,
) *BoolCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 3, 4)
	args[0] = "expire"
	args[1] = key
	args[2] = formatSec(ctx, expiration)
	if mode != "" ***REMOVED***
		args = append(args, mode)
	***REMOVED***

	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "expireat", key, tm.Unix())
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Keys(ctx context.Context, pattern string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "keys", pattern)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(
		ctx,
		"migrate",
		host,
		port,
		key,
		db,
		formatMs(ctx, timeout),
	)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Move(ctx context.Context, key string, db int) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "move", key, db)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ObjectRefCount(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "object", "refcount", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ObjectEncoding(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "object", "encoding", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ObjectIdleTime(ctx context.Context, key string) *DurationCmd ***REMOVED***
	cmd := NewDurationCmd(ctx, time.Second, "object", "idletime", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Persist(ctx context.Context, key string) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "persist", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "pexpire", key, formatMs(ctx, expiration))
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(
		ctx,
		"pexpireat",
		key,
		tm.UnixNano()/int64(time.Millisecond),
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PTTL(ctx context.Context, key string) *DurationCmd ***REMOVED***
	cmd := NewDurationCmd(ctx, time.Millisecond, "pttl", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RandomKey(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "randomkey")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Rename(ctx context.Context, key, newkey string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "rename", key, newkey)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RenameNX(ctx context.Context, key, newkey string) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "renamenx", key, newkey)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Restore(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(
		ctx,
		"restore",
		key,
		formatMs(ctx, ttl),
		value,
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(
		ctx,
		"restore",
		key,
		formatMs(ctx, ttl),
		value,
		"replace",
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type Sort struct ***REMOVED***
	By            string
	Offset, Count int64
	Get           []string
	Order         string
	Alpha         bool
***REMOVED***

func (sort *Sort) args(key string) []interface***REMOVED******REMOVED*** ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"sort", key***REMOVED***
	if sort.By != "" ***REMOVED***
		args = append(args, "by", sort.By)
	***REMOVED***
	if sort.Offset != 0 || sort.Count != 0 ***REMOVED***
		args = append(args, "limit", sort.Offset, sort.Count)
	***REMOVED***
	for _, get := range sort.Get ***REMOVED***
		args = append(args, "get", get)
	***REMOVED***
	if sort.Order != "" ***REMOVED***
		args = append(args, sort.Order)
	***REMOVED***
	if sort.Alpha ***REMOVED***
		args = append(args, "alpha")
	***REMOVED***
	return args
***REMOVED***

func (c cmdable) Sort(ctx context.Context, key string, sort *Sort) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, sort.args(key)...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SortStore(ctx context.Context, key, store string, sort *Sort) *IntCmd ***REMOVED***
	args := sort.args(key)
	if store != "" ***REMOVED***
		args = append(args, "store", store)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SortInterfaces(ctx context.Context, key string, sort *Sort) *SliceCmd ***REMOVED***
	cmd := NewSliceCmd(ctx, sort.args(key)...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Touch(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, len(keys)+1)
	args[0] = "touch"
	for i, key := range keys ***REMOVED***
		args[i+1] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) TTL(ctx context.Context, key string) *DurationCmd ***REMOVED***
	cmd := NewDurationCmd(ctx, time.Second, "ttl", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Type(ctx context.Context, key string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "type", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Append(ctx context.Context, key, value string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "append", key, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Decr(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "decr", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) DecrBy(ctx context.Context, key string, decrement int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "decrby", key, decrement)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (c cmdable) Get(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "get", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GetRange(ctx context.Context, key string, start, end int64) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "getrange", key, start, end)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GetSet(ctx context.Context, key string, value interface***REMOVED******REMOVED***) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "getset", key, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GetEx An expiration of zero removes the TTL associated with the key (i.e. GETEX key persist).
// Requires Redis >= 6.2.0.
func (c cmdable) GetEx(ctx context.Context, key string, expiration time.Duration) *StringCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 4)
	args = append(args, "getex", key)
	if expiration > 0 ***REMOVED***
		if usePrecise(expiration) ***REMOVED***
			args = append(args, "px", formatMs(ctx, expiration))
		***REMOVED*** else ***REMOVED***
			args = append(args, "ex", formatSec(ctx, expiration))
		***REMOVED***
	***REMOVED*** else if expiration == 0 ***REMOVED***
		args = append(args, "persist")
	***REMOVED***

	cmd := NewStringCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GetDel redis-server version >= 6.2.0.
func (c cmdable) GetDel(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "getdel", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Incr(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "incr", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) IncrBy(ctx context.Context, key string, value int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "incrby", key, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd ***REMOVED***
	cmd := NewFloatCmd(ctx, "incrbyfloat", key, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) MGet(ctx context.Context, keys ...string) *SliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "mget"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// MSet is like Set but accepts multiple values:
//   - MSet("key1", "value1", "key2", "value2")
//   - MSet([]string***REMOVED***"key1", "value1", "key2", "value2"***REMOVED***)
//   - MSet(map[string]interface***REMOVED******REMOVED******REMOVED***"key1": "value1", "key2": "value2"***REMOVED***)
func (c cmdable) MSet(ctx context.Context, values ...interface***REMOVED******REMOVED***) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1, 1+len(values))
	args[0] = "mset"
	args = appendArgs(args, values)
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// MSetNX is like SetNX but accepts multiple values:
//   - MSetNX("key1", "value1", "key2", "value2")
//   - MSetNX([]string***REMOVED***"key1", "value1", "key2", "value2"***REMOVED***)
//   - MSetNX(map[string]interface***REMOVED******REMOVED******REMOVED***"key1": "value1", "key2": "value2"***REMOVED***)
func (c cmdable) MSetNX(ctx context.Context, values ...interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1, 1+len(values))
	args[0] = "msetnx"
	args = appendArgs(args, values)
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// Set Redis `SET key value [expiration]` command.
// Use expiration for `SETEx`-like behavior.
//
// Zero expiration means the key has no expiration time.
// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
func (c cmdable) Set(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 3, 5)
	args[0] = "set"
	args[1] = key
	args[2] = value
	if expiration > 0 ***REMOVED***
		if usePrecise(expiration) ***REMOVED***
			args = append(args, "px", formatMs(ctx, expiration))
		***REMOVED*** else ***REMOVED***
			args = append(args, "ex", formatSec(ctx, expiration))
		***REMOVED***
	***REMOVED*** else if expiration == KeepTTL ***REMOVED***
		args = append(args, "keepttl")
	***REMOVED***

	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SetArgs provides arguments for the SetArgs function.
type SetArgs struct ***REMOVED***
	// Mode can be `NX` or `XX` or empty.
	Mode string

	// Zero `TTL` or `Expiration` means that the key has no expiration time.
	TTL      time.Duration
	ExpireAt time.Time

	// When Get is true, the command returns the old value stored at key, or nil when key did not exist.
	Get bool

	// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
	// otherwise you will receive an error: (error) ERR syntax error.
	KeepTTL bool
***REMOVED***

// SetArgs supports all the options that the SET command supports.
// It is the alternative to the Set function when you want
// to have more control over the options.
func (c cmdable) SetArgs(ctx context.Context, key string, value interface***REMOVED******REMOVED***, a SetArgs) *StatusCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"set", key, value***REMOVED***

	if a.KeepTTL ***REMOVED***
		args = append(args, "keepttl")
	***REMOVED***

	if !a.ExpireAt.IsZero() ***REMOVED***
		args = append(args, "exat", a.ExpireAt.Unix())
	***REMOVED***
	if a.TTL > 0 ***REMOVED***
		if usePrecise(a.TTL) ***REMOVED***
			args = append(args, "px", formatMs(ctx, a.TTL))
		***REMOVED*** else ***REMOVED***
			args = append(args, "ex", formatSec(ctx, a.TTL))
		***REMOVED***
	***REMOVED***

	if a.Mode != "" ***REMOVED***
		args = append(args, a.Mode)
	***REMOVED***

	if a.Get ***REMOVED***
		args = append(args, "get")
	***REMOVED***

	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SetEx Redis `SETEx key expiration value` command.
func (c cmdable) SetEx(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "setex", key, formatSec(ctx, expiration), value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SetNX Redis `SET key value [expiration] NX` command.
//
// Zero expiration means the key has no expiration time.
// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
func (c cmdable) SetNX(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *BoolCmd ***REMOVED***
	var cmd *BoolCmd
	switch expiration ***REMOVED***
	case 0:
		// Use old `SETNX` to support old Redis versions.
		cmd = NewBoolCmd(ctx, "setnx", key, value)
	case KeepTTL:
		cmd = NewBoolCmd(ctx, "set", key, value, "keepttl", "nx")
	default:
		if usePrecise(expiration) ***REMOVED***
			cmd = NewBoolCmd(ctx, "set", key, value, "px", formatMs(ctx, expiration), "nx")
		***REMOVED*** else ***REMOVED***
			cmd = NewBoolCmd(ctx, "set", key, value, "ex", formatSec(ctx, expiration), "nx")
		***REMOVED***
	***REMOVED***

	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SetXX Redis `SET key value [expiration] XX` command.
//
// Zero expiration means the key has no expiration time.
// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
func (c cmdable) SetXX(ctx context.Context, key string, value interface***REMOVED******REMOVED***, expiration time.Duration) *BoolCmd ***REMOVED***
	var cmd *BoolCmd
	switch expiration ***REMOVED***
	case 0:
		cmd = NewBoolCmd(ctx, "set", key, value, "xx")
	case KeepTTL:
		cmd = NewBoolCmd(ctx, "set", key, value, "keepttl", "xx")
	default:
		if usePrecise(expiration) ***REMOVED***
			cmd = NewBoolCmd(ctx, "set", key, value, "px", formatMs(ctx, expiration), "xx")
		***REMOVED*** else ***REMOVED***
			cmd = NewBoolCmd(ctx, "set", key, value, "ex", formatSec(ctx, expiration), "xx")
		***REMOVED***
	***REMOVED***

	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "setrange", key, offset, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) StrLen(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "strlen", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Copy(ctx context.Context, sourceKey, destKey string, db int, replace bool) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"copy", sourceKey, destKey, "DB", db***REMOVED***
	if replace ***REMOVED***
		args = append(args, "REPLACE")
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) GetBit(ctx context.Context, key string, offset int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "getbit", key, offset)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd ***REMOVED***
	cmd := NewIntCmd(
		ctx,
		"setbit",
		key,
		offset,
		value,
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type BitCount struct ***REMOVED***
	Start, End int64
***REMOVED***

func (c cmdable) BitCount(ctx context.Context, key string, bitCount *BitCount) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"bitcount", key***REMOVED***
	if bitCount != nil ***REMOVED***
		args = append(
			args,
			bitCount.Start,
			bitCount.End,
		)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) bitOp(ctx context.Context, op, destKey string, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 3+len(keys))
	args[0] = "bitop"
	args[1] = op
	args[2] = destKey
	for i, key := range keys ***REMOVED***
		args[3+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmd ***REMOVED***
	return c.bitOp(ctx, "and", destKey, keys...)
***REMOVED***

func (c cmdable) BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmd ***REMOVED***
	return c.bitOp(ctx, "or", destKey, keys...)
***REMOVED***

func (c cmdable) BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmd ***REMOVED***
	return c.bitOp(ctx, "xor", destKey, keys...)
***REMOVED***

func (c cmdable) BitOpNot(ctx context.Context, destKey string, key string) *IntCmd ***REMOVED***
	return c.bitOp(ctx, "not", destKey, key)
***REMOVED***

func (c cmdable) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 3+len(pos))
	args[0] = "bitpos"
	args[1] = key
	args[2] = bit
	switch len(pos) ***REMOVED***
	case 0:
	case 1:
		args[3] = pos[0]
	case 2:
		args[3] = pos[0]
		args[4] = pos[1]
	default:
		panic("too many arguments")
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BitField(ctx context.Context, key string, args ...interface***REMOVED******REMOVED***) *IntSliceCmd ***REMOVED***
	a := make([]interface***REMOVED******REMOVED***, 0, 2+len(args))
	a = append(a, "bitfield")
	a = append(a, key)
	a = append(a, args...)
	cmd := NewIntSliceCmd(ctx, a...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) Scan(ctx context.Context, cursor uint64, match string, count int64) *ScanCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"scan", cursor***REMOVED***
	if match != "" ***REMOVED***
		args = append(args, "match", match)
	***REMOVED***
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	cmd := NewScanCmd(ctx, c, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *ScanCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"scan", cursor***REMOVED***
	if match != "" ***REMOVED***
		args = append(args, "match", match)
	***REMOVED***
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	if keyType != "" ***REMOVED***
		args = append(args, "type", keyType)
	***REMOVED***
	cmd := NewScanCmd(ctx, c, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"sscan", key, cursor***REMOVED***
	if match != "" ***REMOVED***
		args = append(args, "match", match)
	***REMOVED***
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	cmd := NewScanCmd(ctx, c, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"hscan", key, cursor***REMOVED***
	if match != "" ***REMOVED***
		args = append(args, "match", match)
	***REMOVED***
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	cmd := NewScanCmd(ctx, c, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"zscan", key, cursor***REMOVED***
	if match != "" ***REMOVED***
		args = append(args, "match", match)
	***REMOVED***
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	cmd := NewScanCmd(ctx, c, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) HDel(ctx context.Context, key string, fields ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(fields))
	args[0] = "hdel"
	args[1] = key
	for i, field := range fields ***REMOVED***
		args[2+i] = field
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HExists(ctx context.Context, key, field string) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "hexists", key, field)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HGet(ctx context.Context, key, field string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "hget", key, field)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HGetAll(ctx context.Context, key string) *MapStringStringCmd ***REMOVED***
	cmd := NewMapStringStringCmd(ctx, "hgetall", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "hincrby", key, field, incr)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd ***REMOVED***
	cmd := NewFloatCmd(ctx, "hincrbyfloat", key, field, incr)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HKeys(ctx context.Context, key string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "hkeys", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HLen(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "hlen", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// HMGet returns the values for the specified fields in the hash stored at key.
// It returns an interface***REMOVED******REMOVED*** to distinguish between empty string and nil value.
func (c cmdable) HMGet(ctx context.Context, key string, fields ...string) *SliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(fields))
	args[0] = "hmget"
	args[1] = key
	for i, field := range fields ***REMOVED***
		args[2+i] = field
	***REMOVED***
	cmd := NewSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string***REMOVED***"key1", "value1", "key2", "value2"***REMOVED***)
//   - HSet("myhash", map[string]interface***REMOVED******REMOVED******REMOVED***"key1": "value1", "key2": "value2"***REMOVED***)
//
// Note that it requires Redis v4 for multiple field/value pairs support.
func (c cmdable) HSet(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "hset"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// HMSet is a deprecated version of HSet left for compatibility with Redis 3.
func (c cmdable) HMSet(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "hmset"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HSetNX(ctx context.Context, key, field string, value interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "hsetnx", key, field, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) HVals(ctx context.Context, key string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "hvals", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// HRandField redis-server version >= 6.2.0.
func (c cmdable) HRandField(ctx context.Context, key string, count int) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "hrandfield", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// HRandFieldWithValues redis-server version >= 6.2.0.
func (c cmdable) HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd ***REMOVED***
	cmd := NewKeyValueSliceCmd(ctx, "hrandfield", key, count, "withvalues")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys)+1)
	args[0] = "blpop"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys)+1)
	args[0] = "brpop"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	args[len(keys)+1] = formatSec(ctx, timeout)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmd ***REMOVED***
	cmd := NewStringCmd(
		ctx,
		"brpoplpush",
		source,
		destination,
		formatSec(ctx, timeout),
	)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LIndex(ctx context.Context, key string, index int64) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "lindex", key, index)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LInsert(ctx context.Context, key, op string, pivot, value interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "linsert", key, op, pivot, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LInsertBefore(ctx context.Context, key string, pivot, value interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "linsert", key, "before", pivot, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LInsertAfter(ctx context.Context, key string, pivot, value interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "linsert", key, "after", pivot, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LLen(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "llen", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LPop(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "lpop", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LPopCount(ctx context.Context, key string, count int) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "lpop", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type LPosArgs struct ***REMOVED***
	Rank, MaxLen int64
***REMOVED***

func (c cmdable) LPos(ctx context.Context, key string, value string, a LPosArgs) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"lpos", key, value***REMOVED***
	if a.Rank != 0 ***REMOVED***
		args = append(args, "rank", a.Rank)
	***REMOVED***
	if a.MaxLen != 0 ***REMOVED***
		args = append(args, "maxlen", a.MaxLen)
	***REMOVED***

	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LPosCount(ctx context.Context, key string, value string, count int64, a LPosArgs) *IntSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"lpos", key, value, "count", count***REMOVED***
	if a.Rank != 0 ***REMOVED***
		args = append(args, "rank", a.Rank)
	***REMOVED***
	if a.MaxLen != 0 ***REMOVED***
		args = append(args, "maxlen", a.MaxLen)
	***REMOVED***
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LPush(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "lpush"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LPushX(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "lpushx"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(
		ctx,
		"lrange",
		key,
		start,
		stop,
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LRem(ctx context.Context, key string, count int64, value interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "lrem", key, count, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LSet(ctx context.Context, key string, index int64, value interface***REMOVED******REMOVED***) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "lset", key, index, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(
		ctx,
		"ltrim",
		key,
		start,
		stop,
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RPop(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "rpop", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RPopCount(ctx context.Context, key string, count int) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "rpop", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RPopLPush(ctx context.Context, source, destination string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "rpoplpush", source, destination)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RPush(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "rpush"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) RPushX(ctx context.Context, key string, values ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(values))
	args[0] = "rpushx"
	args[1] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LMove(ctx context.Context, source, destination, srcpos, destpos string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "lmove", source, destination, srcpos, destpos)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BLMove(
	ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration,
) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "blmove", source, destination, srcpos, destpos, formatSec(ctx, timeout))
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) SAdd(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(members))
	args[0] = "sadd"
	args[1] = key
	args = appendArgs(args, members)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SCard(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "scard", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SDiff(ctx context.Context, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "sdiff"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "sdiffstore"
	args[1] = destination
	for i, key := range keys ***REMOVED***
		args[2+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SInter(ctx context.Context, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "sinter"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SInterStore(ctx context.Context, destination string, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "sinterstore"
	args[1] = destination
	for i, key := range keys ***REMOVED***
		args[2+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SIsMember(ctx context.Context, key string, member interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "sismember", key, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SMIsMember Redis `SMISMEMBER key member [member ...]` command.
func (c cmdable) SMIsMember(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *BoolSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(members))
	args[0] = "smismember"
	args[1] = key
	args = appendArgs(args, members)
	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SMembers Redis `SMEMBERS key` command output as a slice.
func (c cmdable) SMembers(ctx context.Context, key string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "smembers", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SMembersMap Redis `SMEMBERS key` command output as a map.
func (c cmdable) SMembersMap(ctx context.Context, key string) *StringStructMapCmd ***REMOVED***
	cmd := NewStringStructMapCmd(ctx, "smembers", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SMove(ctx context.Context, source, destination string, member interface***REMOVED******REMOVED***) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "smove", source, destination, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SPop Redis `SPOP key` command.
func (c cmdable) SPop(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "spop", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SPopN Redis `SPOP key count` command.
func (c cmdable) SPopN(ctx context.Context, key string, count int64) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "spop", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SRandMember Redis `SRANDMEMBER key` command.
func (c cmdable) SRandMember(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "srandmember", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// SRandMemberN Redis `SRANDMEMBER key count` command.
func (c cmdable) SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "srandmember", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SRem(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(members))
	args[0] = "srem"
	args[1] = key
	args = appendArgs(args, members)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SUnion(ctx context.Context, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "sunion"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "sunionstore"
	args[1] = destination
	for i, key := range keys ***REMOVED***
		args[2+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

// XAddArgs accepts values in the following formats:
//   - XAddArgs.Values = []interface***REMOVED******REMOVED******REMOVED***"key1", "value1", "key2", "value2"***REMOVED***
//   - XAddArgs.Values = []string("key1", "value1", "key2", "value2")
//   - XAddArgs.Values = map[string]interface***REMOVED******REMOVED******REMOVED***"key1": "value1", "key2": "value2"***REMOVED***
//
// Note that map will not preserve the order of key-value pairs.
// MaxLen/MaxLenApprox and MinID are in conflict, only one of them can be used.
type XAddArgs struct ***REMOVED***
	Stream     string
	NoMkStream bool
	MaxLen     int64 // MAXLEN N
	MinID      string
	// Approx causes MaxLen and MinID to use "~" matcher (instead of "=").
	Approx bool
	Limit  int64
	ID     string
	Values interface***REMOVED******REMOVED***
***REMOVED***

// XAdd a.Limit has a bug, please confirm it and use it.
// issue: https://github.com/redis/redis/issues/9046
func (c cmdable) XAdd(ctx context.Context, a *XAddArgs) *StringCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 11)
	args = append(args, "xadd", a.Stream)
	if a.NoMkStream ***REMOVED***
		args = append(args, "nomkstream")
	***REMOVED***
	switch ***REMOVED***
	case a.MaxLen > 0:
		if a.Approx ***REMOVED***
			args = append(args, "maxlen", "~", a.MaxLen)
		***REMOVED*** else ***REMOVED***
			args = append(args, "maxlen", a.MaxLen)
		***REMOVED***
	case a.MinID != "":
		if a.Approx ***REMOVED***
			args = append(args, "minid", "~", a.MinID)
		***REMOVED*** else ***REMOVED***
			args = append(args, "minid", a.MinID)
		***REMOVED***
	***REMOVED***
	if a.Limit > 0 ***REMOVED***
		args = append(args, "limit", a.Limit)
	***REMOVED***
	if a.ID != "" ***REMOVED***
		args = append(args, a.ID)
	***REMOVED*** else ***REMOVED***
		args = append(args, "*")
	***REMOVED***
	args = appendArg(args, a.Values)

	cmd := NewStringCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XDel(ctx context.Context, stream string, ids ...string) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"xdel", stream***REMOVED***
	for _, id := range ids ***REMOVED***
		args = append(args, id)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XLen(ctx context.Context, stream string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "xlen", stream)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XRange(ctx context.Context, stream, start, stop string) *XMessageSliceCmd ***REMOVED***
	cmd := NewXMessageSliceCmd(ctx, "xrange", stream, start, stop)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XRangeN(ctx context.Context, stream, start, stop string, count int64) *XMessageSliceCmd ***REMOVED***
	cmd := NewXMessageSliceCmd(ctx, "xrange", stream, start, stop, "count", count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XRevRange(ctx context.Context, stream, start, stop string) *XMessageSliceCmd ***REMOVED***
	cmd := NewXMessageSliceCmd(ctx, "xrevrange", stream, start, stop)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XRevRangeN(ctx context.Context, stream, start, stop string, count int64) *XMessageSliceCmd ***REMOVED***
	cmd := NewXMessageSliceCmd(ctx, "xrevrange", stream, start, stop, "count", count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type XReadArgs struct ***REMOVED***
	Streams []string // list of streams and ids, e.g. stream1 stream2 id1 id2
	Count   int64
	Block   time.Duration
***REMOVED***

func (c cmdable) XRead(ctx context.Context, a *XReadArgs) *XStreamSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 6+len(a.Streams))
	args = append(args, "xread")

	keyPos := int8(1)
	if a.Count > 0 ***REMOVED***
		args = append(args, "count")
		args = append(args, a.Count)
		keyPos += 2
	***REMOVED***
	if a.Block >= 0 ***REMOVED***
		args = append(args, "block")
		args = append(args, int64(a.Block/time.Millisecond))
		keyPos += 2
	***REMOVED***
	args = append(args, "streams")
	keyPos++
	for _, s := range a.Streams ***REMOVED***
		args = append(args, s)
	***REMOVED***

	cmd := NewXStreamSliceCmd(ctx, args...)
	if a.Block >= 0 ***REMOVED***
		cmd.setReadTimeout(a.Block)
	***REMOVED***
	cmd.SetFirstKeyPos(keyPos)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XReadStreams(ctx context.Context, streams ...string) *XStreamSliceCmd ***REMOVED***
	return c.XRead(ctx, &XReadArgs***REMOVED***
		Streams: streams,
		Block:   -1,
	***REMOVED***)
***REMOVED***

func (c cmdable) XGroupCreate(ctx context.Context, stream, group, start string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "xgroup", "create", stream, group, start)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "xgroup", "create", stream, group, start, "mkstream")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XGroupSetID(ctx context.Context, stream, group, start string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "xgroup", "setid", stream, group, start)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XGroupDestroy(ctx context.Context, stream, group string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "xgroup", "destroy", stream, group)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "xgroup", "createconsumer", stream, group, consumer)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "xgroup", "delconsumer", stream, group, consumer)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type XReadGroupArgs struct ***REMOVED***
	Group    string
	Consumer string
	Streams  []string // list of streams and ids, e.g. stream1 stream2 id1 id2
	Count    int64
	Block    time.Duration
	NoAck    bool
***REMOVED***

func (c cmdable) XReadGroup(ctx context.Context, a *XReadGroupArgs) *XStreamSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 10+len(a.Streams))
	args = append(args, "xreadgroup", "group", a.Group, a.Consumer)

	keyPos := int8(4)
	if a.Count > 0 ***REMOVED***
		args = append(args, "count", a.Count)
		keyPos += 2
	***REMOVED***
	if a.Block >= 0 ***REMOVED***
		args = append(args, "block", int64(a.Block/time.Millisecond))
		keyPos += 2
	***REMOVED***
	if a.NoAck ***REMOVED***
		args = append(args, "noack")
		keyPos++
	***REMOVED***
	args = append(args, "streams")
	keyPos++
	for _, s := range a.Streams ***REMOVED***
		args = append(args, s)
	***REMOVED***

	cmd := NewXStreamSliceCmd(ctx, args...)
	if a.Block >= 0 ***REMOVED***
		cmd.setReadTimeout(a.Block)
	***REMOVED***
	cmd.SetFirstKeyPos(keyPos)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XAck(ctx context.Context, stream, group string, ids ...string) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"xack", stream, group***REMOVED***
	for _, id := range ids ***REMOVED***
		args = append(args, id)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XPending(ctx context.Context, stream, group string) *XPendingCmd ***REMOVED***
	cmd := NewXPendingCmd(ctx, "xpending", stream, group)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type XPendingExtArgs struct ***REMOVED***
	Stream   string
	Group    string
	Idle     time.Duration
	Start    string
	End      string
	Count    int64
	Consumer string
***REMOVED***

func (c cmdable) XPendingExt(ctx context.Context, a *XPendingExtArgs) *XPendingExtCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 9)
	args = append(args, "xpending", a.Stream, a.Group)
	if a.Idle != 0 ***REMOVED***
		args = append(args, "idle", formatMs(ctx, a.Idle))
	***REMOVED***
	args = append(args, a.Start, a.End, a.Count)
	if a.Consumer != "" ***REMOVED***
		args = append(args, a.Consumer)
	***REMOVED***
	cmd := NewXPendingExtCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

type XAutoClaimArgs struct ***REMOVED***
	Stream   string
	Group    string
	MinIdle  time.Duration
	Start    string
	Count    int64
	Consumer string
***REMOVED***

func (c cmdable) XAutoClaim(ctx context.Context, a *XAutoClaimArgs) *XAutoClaimCmd ***REMOVED***
	args := xAutoClaimArgs(ctx, a)
	cmd := NewXAutoClaimCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XAutoClaimJustID(ctx context.Context, a *XAutoClaimArgs) *XAutoClaimJustIDCmd ***REMOVED***
	args := xAutoClaimArgs(ctx, a)
	args = append(args, "justid")
	cmd := NewXAutoClaimJustIDCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func xAutoClaimArgs(ctx context.Context, a *XAutoClaimArgs) []interface***REMOVED******REMOVED*** ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 8)
	args = append(args, "xautoclaim", a.Stream, a.Group, a.Consumer, formatMs(ctx, a.MinIdle), a.Start)
	if a.Count > 0 ***REMOVED***
		args = append(args, "count", a.Count)
	***REMOVED***
	return args
***REMOVED***

type XClaimArgs struct ***REMOVED***
	Stream   string
	Group    string
	Consumer string
	MinIdle  time.Duration
	Messages []string
***REMOVED***

func (c cmdable) XClaim(ctx context.Context, a *XClaimArgs) *XMessageSliceCmd ***REMOVED***
	args := xClaimArgs(a)
	cmd := NewXMessageSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XClaimJustID(ctx context.Context, a *XClaimArgs) *StringSliceCmd ***REMOVED***
	args := xClaimArgs(a)
	args = append(args, "justid")
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func xClaimArgs(a *XClaimArgs) []interface***REMOVED******REMOVED*** ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 5+len(a.Messages))
	args = append(args,
		"xclaim",
		a.Stream,
		a.Group, a.Consumer,
		int64(a.MinIdle/time.Millisecond))
	for _, id := range a.Messages ***REMOVED***
		args = append(args, id)
	***REMOVED***
	return args
***REMOVED***

// xTrim If approx is true, add the "~" parameter, otherwise it is the default "=" (redis default).
// example:
//		XTRIM key MAXLEN/MINID threshold LIMIT limit.
//		XTRIM key MAXLEN/MINID ~ threshold LIMIT limit.
// The redis-server version is lower than 6.2, please set limit to 0.
func (c cmdable) xTrim(
	ctx context.Context, key, strategy string,
	approx bool, threshold interface***REMOVED******REMOVED***, limit int64,
) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 7)
	args = append(args, "xtrim", key, strategy)
	if approx ***REMOVED***
		args = append(args, "~")
	***REMOVED***
	args = append(args, threshold)
	if limit > 0 ***REMOVED***
		args = append(args, "limit", limit)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// XTrimMaxLen No `~` rules are used, `limit` cannot be used.
// cmd: XTRIM key MAXLEN maxLen
func (c cmdable) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *IntCmd ***REMOVED***
	return c.xTrim(ctx, key, "maxlen", false, maxLen, 0)
***REMOVED***

// XTrimMaxLenApprox LIMIT has a bug, please confirm it and use it.
// issue: https://github.com/redis/redis/issues/9046
// cmd: XTRIM key MAXLEN ~ maxLen LIMIT limit
func (c cmdable) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *IntCmd ***REMOVED***
	return c.xTrim(ctx, key, "maxlen", true, maxLen, limit)
***REMOVED***

// XTrimMinID No `~` rules are used, `limit` cannot be used.
// cmd: XTRIM key MINID minID
func (c cmdable) XTrimMinID(ctx context.Context, key string, minID string) *IntCmd ***REMOVED***
	return c.xTrim(ctx, key, "minid", false, minID, 0)
***REMOVED***

// XTrimMinIDApprox LIMIT has a bug, please confirm it and use it.
// issue: https://github.com/redis/redis/issues/9046
// cmd: XTRIM key MINID ~ minID LIMIT limit
func (c cmdable) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *IntCmd ***REMOVED***
	return c.xTrim(ctx, key, "minid", true, minID, limit)
***REMOVED***

func (c cmdable) XInfoConsumers(ctx context.Context, key string, group string) *XInfoConsumersCmd ***REMOVED***
	cmd := NewXInfoConsumersCmd(ctx, key, group)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XInfoGroups(ctx context.Context, key string) *XInfoGroupsCmd ***REMOVED***
	cmd := NewXInfoGroupsCmd(ctx, key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) XInfoStream(ctx context.Context, key string) *XInfoStreamCmd ***REMOVED***
	cmd := NewXInfoStreamCmd(ctx, key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// XInfoStreamFull XINFO STREAM FULL [COUNT count]
// redis-server >= 6.0.
func (c cmdable) XInfoStreamFull(ctx context.Context, key string, count int) *XInfoStreamFullCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 6)
	args = append(args, "xinfo", "stream", key, "full")
	if count > 0 ***REMOVED***
		args = append(args, "count", count)
	***REMOVED***
	cmd := NewXInfoStreamFullCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

// Z represents sorted set member.
type Z struct ***REMOVED***
	Score  float64
	Member interface***REMOVED******REMOVED***
***REMOVED***

// ZWithKey represents sorted set member including the name of the key where it was popped.
type ZWithKey struct ***REMOVED***
	Z
	Key string
***REMOVED***

// ZStore is used as an arg to ZInter/ZInterStore and ZUnion/ZUnionStore.
type ZStore struct ***REMOVED***
	Keys    []string
	Weights []float64
	// Can be SUM, MIN or MAX.
	Aggregate string
***REMOVED***

func (z ZStore) len() (n int) ***REMOVED***
	n = len(z.Keys)
	if len(z.Weights) > 0 ***REMOVED***
		n += 1 + len(z.Weights)
	***REMOVED***
	if z.Aggregate != "" ***REMOVED***
		n += 2
	***REMOVED***
	return n
***REMOVED***

func (z ZStore) appendArgs(args []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	for _, key := range z.Keys ***REMOVED***
		args = append(args, key)
	***REMOVED***
	if len(z.Weights) > 0 ***REMOVED***
		args = append(args, "weights")
		for _, weights := range z.Weights ***REMOVED***
			args = append(args, weights)
		***REMOVED***
	***REMOVED***
	if z.Aggregate != "" ***REMOVED***
		args = append(args, "aggregate", z.Aggregate)
	***REMOVED***
	return args
***REMOVED***

// BZPopMax Redis `BZPOPMAX key [key ...] timeout` command.
func (c cmdable) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys)+1)
	args[0] = "bzpopmax"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// BZPopMin Redis `BZPOPMIN key [key ...] timeout` command.
func (c cmdable) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys)+1)
	args[0] = "bzpopmin"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd(ctx, args...)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZAddArgs WARN: The GT, LT and NX options are mutually exclusive.
type ZAddArgs struct ***REMOVED***
	NX      bool
	XX      bool
	LT      bool
	GT      bool
	Ch      bool
	Members []Z
***REMOVED***

func (c cmdable) zAddArgs(key string, args ZAddArgs, incr bool) []interface***REMOVED******REMOVED*** ***REMOVED***
	a := make([]interface***REMOVED******REMOVED***, 0, 6+2*len(args.Members))
	a = append(a, "zadd", key)

	// The GT, LT and NX options are mutually exclusive.
	if args.NX ***REMOVED***
		a = append(a, "nx")
	***REMOVED*** else ***REMOVED***
		if args.XX ***REMOVED***
			a = append(a, "xx")
		***REMOVED***
		if args.GT ***REMOVED***
			a = append(a, "gt")
		***REMOVED*** else if args.LT ***REMOVED***
			a = append(a, "lt")
		***REMOVED***
	***REMOVED***
	if args.Ch ***REMOVED***
		a = append(a, "ch")
	***REMOVED***
	if incr ***REMOVED***
		a = append(a, "incr")
	***REMOVED***
	for _, m := range args.Members ***REMOVED***
		a = append(a, m.Score)
		a = append(a, m.Member)
	***REMOVED***
	return a
***REMOVED***

func (c cmdable) ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, c.zAddArgs(key, args, false)...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd ***REMOVED***
	cmd := NewFloatCmd(ctx, c.zAddArgs(key, args, true)...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZAdd Redis `ZADD key score member [score member ...]` command.
func (c cmdable) ZAdd(ctx context.Context, key string, members ...Z) *IntCmd ***REMOVED***
	return c.ZAddArgs(ctx, key, ZAddArgs***REMOVED***
		Members: members,
	***REMOVED***)
***REMOVED***

// ZAddNX Redis `ZADD key NX score member [score member ...]` command.
func (c cmdable) ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd ***REMOVED***
	return c.ZAddArgs(ctx, key, ZAddArgs***REMOVED***
		NX:      true,
		Members: members,
	***REMOVED***)
***REMOVED***

// ZAddXX Redis `ZADD key XX score member [score member ...]` command.
func (c cmdable) ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd ***REMOVED***
	return c.ZAddArgs(ctx, key, ZAddArgs***REMOVED***
		XX:      true,
		Members: members,
	***REMOVED***)
***REMOVED***

func (c cmdable) ZCard(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zcard", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZCount(ctx context.Context, key, min, max string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zcount", key, min, max)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZLexCount(ctx context.Context, key, min, max string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zlexcount", key, min, max)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd ***REMOVED***
	cmd := NewFloatCmd(ctx, "zincrby", key, increment, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 3+store.len())
	args = append(args, "zinterstore", destination, len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZInter(ctx context.Context, store *ZStore) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 2+store.len())
	args = append(args, "zinter", len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 3+store.len())
	args = append(args, "zinter", len(store.Keys))
	args = store.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(members))
	args[0] = "zmscore"
	args[1] = key
	for i, member := range members ***REMOVED***
		args[2+i] = member
	***REMOVED***
	cmd := NewFloatSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***
		"zpopmax",
		key,
	***REMOVED***

	switch len(count) ***REMOVED***
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	***REMOVED***

	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***
		"zpopmin",
		key,
	***REMOVED***

	switch len(count) ***REMOVED***
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	***REMOVED***

	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZRangeArgs is all the options of the ZRange command.
// In version> 6.2.0, you can replace the(cmd):
//		ZREVRANGE,
//		ZRANGEBYSCORE,
//		ZREVRANGEBYSCORE,
//		ZRANGEBYLEX,
//		ZREVRANGEBYLEX.
// Please pay attention to your redis-server version.
//
// Rev, ByScore, ByLex and Offset+Count options require redis-server 6.2.0 and higher.
type ZRangeArgs struct ***REMOVED***
	Key string

	// When the ByScore option is provided, the open interval(exclusive) can be set.
	// By default, the score intervals specified by <Start> and <Stop> are closed (inclusive).
	// It is similar to the deprecated(6.2.0+) ZRangeByScore command.
	// For example:
	//		ZRangeArgs***REMOVED***
	//			Key: 				"example-key",
	//	 		Start: 				"(3",
	//	 		Stop: 				8,
	//			ByScore:			true,
	//	 	***REMOVED***
	// 	 	cmd: "ZRange example-key (3 8 ByScore"  (3 < score <= 8).
	//
	// For the ByLex option, it is similar to the deprecated(6.2.0+) ZRangeByLex command.
	// You can set the <Start> and <Stop> options as follows:
	//		ZRangeArgs***REMOVED***
	//			Key: 				"example-key",
	//	 		Start: 				"[abc",
	//	 		Stop: 				"(def",
	//			ByLex:				true,
	//	 	***REMOVED***
	//		cmd: "ZRange example-key [abc (def ByLex"
	//
	// For normal cases (ByScore==false && ByLex==false), <Start> and <Stop> should be set to the index range (int).
	// You can read the documentation for more information: https://redis.io/commands/zrange
	Start interface***REMOVED******REMOVED***
	Stop  interface***REMOVED******REMOVED***

	// The ByScore and ByLex options are mutually exclusive.
	ByScore bool
	ByLex   bool

	Rev bool

	// limit offset count.
	Offset int64
	Count  int64
***REMOVED***

func (z ZRangeArgs) appendArgs(args []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	// For Rev+ByScore/ByLex, we need to adjust the position of <Start> and <Stop>.
	if z.Rev && (z.ByScore || z.ByLex) ***REMOVED***
		args = append(args, z.Key, z.Stop, z.Start)
	***REMOVED*** else ***REMOVED***
		args = append(args, z.Key, z.Start, z.Stop)
	***REMOVED***

	if z.ByScore ***REMOVED***
		args = append(args, "byscore")
	***REMOVED*** else if z.ByLex ***REMOVED***
		args = append(args, "bylex")
	***REMOVED***
	if z.Rev ***REMOVED***
		args = append(args, "rev")
	***REMOVED***
	if z.Offset != 0 || z.Count != 0 ***REMOVED***
		args = append(args, "limit", z.Offset, z.Count)
	***REMOVED***
	return args
***REMOVED***

func (c cmdable) ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 9)
	args = append(args, "zrange")
	args = z.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 10)
	args = append(args, "zrange")
	args = z.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd ***REMOVED***
	return c.ZRangeArgs(ctx, ZRangeArgs***REMOVED***
		Key:   key,
		Start: start,
		Stop:  stop,
	***REMOVED***)
***REMOVED***

func (c cmdable) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd ***REMOVED***
	return c.ZRangeArgsWithScores(ctx, ZRangeArgs***REMOVED***
		Key:   key,
		Start: start,
		Stop:  stop,
	***REMOVED***)
***REMOVED***

type ZRangeBy struct ***REMOVED***
	Min, Max      string
	Offset, Count int64
***REMOVED***

func (c cmdable) zRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy, withScores bool) *StringSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***zcmd, key, opt.Min, opt.Max***REMOVED***
	if withScores ***REMOVED***
		args = append(args, "withscores")
	***REMOVED***
	if opt.Offset != 0 || opt.Count != 0 ***REMOVED***
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd ***REMOVED***
	return c.zRangeBy(ctx, "zrangebyscore", key, opt, false)
***REMOVED***

func (c cmdable) ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd ***REMOVED***
	return c.zRangeBy(ctx, "zrangebylex", key, opt, false)
***REMOVED***

func (c cmdable) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"zrangebyscore", key, opt.Min, opt.Max, "withscores"***REMOVED***
	if opt.Offset != 0 || opt.Count != 0 ***REMOVED***
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	***REMOVED***
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 10)
	args = append(args, "zrangestore", dst)
	args = z.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRank(ctx context.Context, key, member string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zrank", key, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRem(ctx context.Context, key string, members ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(members))
	args[0] = "zrem"
	args[1] = key
	args = appendArgs(args, members)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(
		ctx,
		"zremrangebyrank",
		key,
		start,
		stop,
	)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zremrangebyscore", key, min, max)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zremrangebylex", key, min, max)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "zrevrange", key, start, stop)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd ***REMOVED***
	cmd := NewZSliceCmd(ctx, "zrevrange", key, start, stop, "withscores")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) zRevRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy) *StringSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***zcmd, key, opt.Max, opt.Min***REMOVED***
	if opt.Offset != 0 || opt.Count != 0 ***REMOVED***
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd ***REMOVED***
	return c.zRevRangeBy(ctx, "zrevrangebyscore", key, opt)
***REMOVED***

func (c cmdable) ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd ***REMOVED***
	return c.zRevRangeBy(ctx, "zrevrangebylex", key, opt)
***REMOVED***

func (c cmdable) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"zrevrangebyscore", key, opt.Max, opt.Min, "withscores"***REMOVED***
	if opt.Offset != 0 || opt.Count != 0 ***REMOVED***
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	***REMOVED***
	cmd := NewZSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZRevRank(ctx context.Context, key, member string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "zrevrank", key, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZScore(ctx context.Context, key, member string) *FloatCmd ***REMOVED***
	cmd := NewFloatCmd(ctx, "zscore", key, member)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZUnion(ctx context.Context, store ZStore) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 2+store.len())
	args = append(args, "zunion", len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 3+store.len())
	args = append(args, "zunion", len(store.Keys))
	args = store.appendArgs(args)
	args = append(args, "withscores")
	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 3+store.len())
	args = append(args, "zunionstore", dest, len(store.Keys))
	args = store.appendArgs(args)
	cmd := NewIntCmd(ctx, args...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZRandMember redis-server version >= 6.2.0.
func (c cmdable) ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "zrandmember", key, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZRandMemberWithScores redis-server version >= 6.2.0.
func (c cmdable) ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd ***REMOVED***
	cmd := NewZSliceCmd(ctx, "zrandmember", key, count, "withscores")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZDiff redis-server version >= 6.2.0.
func (c cmdable) ZDiff(ctx context.Context, keys ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "zdiff"
	args[1] = len(keys)
	for i, key := range keys ***REMOVED***
		args[i+2] = key
	***REMOVED***

	cmd := NewStringSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZDiffWithScores redis-server version >= 6.2.0.
func (c cmdable) ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 3+len(keys))
	args[0] = "zdiff"
	args[1] = len(keys)
	for i, key := range keys ***REMOVED***
		args[i+2] = key
	***REMOVED***
	args[len(keys)+2] = "withscores"

	cmd := NewZSliceCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ZDiffStore redis-server version >=6.2.0.
func (c cmdable) ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 3+len(keys))
	args = append(args, "zdiffstore", destination, len(keys))
	for _, key := range keys ***REMOVED***
		args = append(args, key)
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) PFAdd(ctx context.Context, key string, els ...interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2, 2+len(els))
	args[0] = "pfadd"
	args[1] = key
	args = appendArgs(args, els)
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PFCount(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "pfcount"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "pfmerge"
	args[1] = dest
	for i, key := range keys ***REMOVED***
		args[2+i] = key
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) BgRewriteAOF(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "bgrewriteaof")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) BgSave(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "bgsave")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientKill(ctx context.Context, ipPort string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "client", "kill", ipPort)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// ClientKillByFilter is new style syntax, while the ClientKill is old
//
//   CLIENT KILL <option> [value] ... <option> [value]
func (c cmdable) ClientKillByFilter(ctx context.Context, keys ...string) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(keys))
	args[0] = "client"
	args[1] = "kill"
	for i, key := range keys ***REMOVED***
		args[2+i] = key
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientList(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "client", "list")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientPause(ctx context.Context, dur time.Duration) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "client", "pause", formatMs(ctx, dur))
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientUnpause(ctx context.Context) *BoolCmd ***REMOVED***
	cmd := NewBoolCmd(ctx, "client", "unpause")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientID(ctx context.Context) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "client", "id")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientUnblock(ctx context.Context, id int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "client", "unblock", id)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClientUnblockWithError(ctx context.Context, id int64) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "client", "unblock", id, "error")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ConfigGet(ctx context.Context, parameter string) *MapStringStringCmd ***REMOVED***
	cmd := NewMapStringStringCmd(ctx, "config", "get", parameter)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ConfigResetStat(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "config", "resetstat")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ConfigSet(ctx context.Context, parameter, value string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "config", "set", parameter, value)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ConfigRewrite(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "config", "rewrite")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) DBSize(ctx context.Context) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "dbsize")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) FlushAll(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "flushall")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) FlushAllAsync(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "flushall", "async")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) FlushDB(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "flushdb")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) FlushDBAsync(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "flushdb", "async")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Info(ctx context.Context, section ...string) *StringCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"info"***REMOVED***
	if len(section) > 0 ***REMOVED***
		args = append(args, section[0])
	***REMOVED***
	cmd := NewStringCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) LastSave(ctx context.Context) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "lastsave")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Save(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "save")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) shutdown(ctx context.Context, modifier string) *StatusCmd ***REMOVED***
	var args []interface***REMOVED******REMOVED***
	if modifier == "" ***REMOVED***
		args = []interface***REMOVED******REMOVED******REMOVED***"shutdown"***REMOVED***
	***REMOVED*** else ***REMOVED***
		args = []interface***REMOVED******REMOVED******REMOVED***"shutdown", modifier***REMOVED***
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	if err := cmd.Err(); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			// Server quit as expected.
			cmd.err = nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Server did not quit. String reply contains the reason.
		cmd.err = errors.New(cmd.val)
		cmd.val = ""
	***REMOVED***
	return cmd
***REMOVED***

func (c cmdable) Shutdown(ctx context.Context) *StatusCmd ***REMOVED***
	return c.shutdown(ctx, "")
***REMOVED***

func (c cmdable) ShutdownSave(ctx context.Context) *StatusCmd ***REMOVED***
	return c.shutdown(ctx, "save")
***REMOVED***

func (c cmdable) ShutdownNoSave(ctx context.Context) *StatusCmd ***REMOVED***
	return c.shutdown(ctx, "nosave")
***REMOVED***

func (c cmdable) SlaveOf(ctx context.Context, host, port string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "slaveof", host, port)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) SlowLogGet(ctx context.Context, num int64) *SlowLogCmd ***REMOVED***
	cmd := NewSlowLogCmd(context.Background(), "slowlog", "get", num)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) Sync(_ context.Context) ***REMOVED***
	panic("not implemented")
***REMOVED***

func (c cmdable) Time(ctx context.Context) *TimeCmd ***REMOVED***
	cmd := NewTimeCmd(ctx, "time")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) DebugObject(ctx context.Context, key string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "debug", "object", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ReadOnly(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "readonly")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ReadWrite(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "readwrite")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) MemoryUsage(ctx context.Context, key string, samples ...int) *IntCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"memory", "usage", key***REMOVED***
	if len(samples) > 0 ***REMOVED***
		if len(samples) != 1 ***REMOVED***
			panic("MemoryUsage expects single sample count")
		***REMOVED***
		args = append(args, "SAMPLES", samples[0])
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) Eval(ctx context.Context, script string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmdArgs := make([]interface***REMOVED******REMOVED***, 3+len(keys), 3+len(keys)+len(args))
	cmdArgs[0] = "eval"
	cmdArgs[1] = script
	cmdArgs[2] = len(keys)
	for i, key := range keys ***REMOVED***
		cmdArgs[3+i] = key
	***REMOVED***
	cmdArgs = appendArgs(cmdArgs, args)
	cmd := NewCmd(ctx, cmdArgs...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmdArgs := make([]interface***REMOVED******REMOVED***, 3+len(keys), 3+len(keys)+len(args))
	cmdArgs[0] = "evalsha"
	cmdArgs[1] = sha1
	cmdArgs[2] = len(keys)
	for i, key := range keys ***REMOVED***
		cmdArgs[3+i] = key
	***REMOVED***
	cmdArgs = appendArgs(cmdArgs, args)
	cmd := NewCmd(ctx, cmdArgs...)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(hashes))
	args[0] = "script"
	args[1] = "exists"
	for i, hash := range hashes ***REMOVED***
		args[2+i] = hash
	***REMOVED***
	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ScriptFlush(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "script", "flush")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ScriptKill(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "script", "kill")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ScriptLoad(ctx context.Context, script string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "script", "load", script)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

// Publish posts the message to the channel.
func (c cmdable) Publish(ctx context.Context, channel string, message interface***REMOVED******REMOVED***) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "publish", channel, message)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PubSubChannels(ctx context.Context, pattern string) *StringSliceCmd ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"pubsub", "channels"***REMOVED***
	if pattern != "*" ***REMOVED***
		args = append(args, pattern)
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PubSubNumSub(ctx context.Context, channels ...string) *StringIntMapCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(channels))
	args[0] = "pubsub"
	args[1] = "numsub"
	for i, channel := range channels ***REMOVED***
		args[2+i] = channel
	***REMOVED***
	cmd := NewStringIntMapCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) PubSubNumPat(ctx context.Context) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "pubsub", "numpat")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) ClusterSlots(ctx context.Context) *ClusterSlotsCmd ***REMOVED***
	cmd := NewClusterSlotsCmd(ctx, "cluster", "slots")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterNodes(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "cluster", "nodes")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterMeet(ctx context.Context, host, port string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "meet", host, port)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterForget(ctx context.Context, nodeID string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "forget", nodeID)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterReplicate(ctx context.Context, nodeID string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "replicate", nodeID)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterResetSoft(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "reset", "soft")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterResetHard(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "reset", "hard")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterInfo(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "cluster", "info")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterKeySlot(ctx context.Context, key string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "cluster", "keyslot", key)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "cluster", "getkeysinslot", slot, count)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterCountFailureReports(ctx context.Context, nodeID string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "cluster", "count-failure-reports", nodeID)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterCountKeysInSlot(ctx context.Context, slot int) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "cluster", "countkeysinslot", slot)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterDelSlots(ctx context.Context, slots ...int) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(slots))
	args[0] = "cluster"
	args[1] = "delslots"
	for i, slot := range slots ***REMOVED***
		args[2+i] = slot
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterDelSlotsRange(ctx context.Context, min, max int) *StatusCmd ***REMOVED***
	size := max - min + 1
	slots := make([]int, size)
	for i := 0; i < size; i++ ***REMOVED***
		slots[i] = min + i
	***REMOVED***
	return c.ClusterDelSlots(ctx, slots...)
***REMOVED***

func (c cmdable) ClusterSaveConfig(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "saveconfig")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterSlaves(ctx context.Context, nodeID string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "cluster", "slaves", nodeID)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterFailover(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "cluster", "failover")
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterAddSlots(ctx context.Context, slots ...int) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(slots))
	args[0] = "cluster"
	args[1] = "addslots"
	for i, num := range slots ***REMOVED***
		args[2+i] = num
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) ClusterAddSlotsRange(ctx context.Context, min, max int) *StatusCmd ***REMOVED***
	size := max - min + 1
	slots := make([]int, size)
	for i := 0; i < size; i++ ***REMOVED***
		slots[i] = min + i
	***REMOVED***
	return c.ClusterAddSlots(ctx, slots...)
***REMOVED***

//------------------------------------------------------------------------------

func (c cmdable) GeoAdd(ctx context.Context, key string, geoLocation ...*GeoLocation) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+3*len(geoLocation))
	args[0] = "geoadd"
	args[1] = key
	for i, eachLoc := range geoLocation ***REMOVED***
		args[2+3*i] = eachLoc.Longitude
		args[2+3*i+1] = eachLoc.Latitude
		args[2+3*i+2] = eachLoc.Name
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GeoRadius is a read-only GEORADIUS_RO command.
func (c cmdable) GeoRadius(
	ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery,
) *GeoLocationCmd ***REMOVED***
	cmd := NewGeoLocationCmd(ctx, query, "georadius_ro", key, longitude, latitude)
	if query.Store != "" || query.StoreDist != "" ***REMOVED***
		cmd.SetErr(errors.New("GeoRadius does not support Store or StoreDist"))
		return cmd
	***REMOVED***
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GeoRadiusStore is a writing GEORADIUS command.
func (c cmdable) GeoRadiusStore(
	ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery,
) *IntCmd ***REMOVED***
	args := geoLocationArgs(query, "georadius", key, longitude, latitude)
	cmd := NewIntCmd(ctx, args...)
	if query.Store == "" && query.StoreDist == "" ***REMOVED***
		cmd.SetErr(errors.New("GeoRadiusStore requires Store or StoreDist"))
		return cmd
	***REMOVED***
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GeoRadiusByMember is a read-only GEORADIUSBYMEMBER_RO command.
func (c cmdable) GeoRadiusByMember(
	ctx context.Context, key, member string, query *GeoRadiusQuery,
) *GeoLocationCmd ***REMOVED***
	cmd := NewGeoLocationCmd(ctx, query, "georadiusbymember_ro", key, member)
	if query.Store != "" || query.StoreDist != "" ***REMOVED***
		cmd.SetErr(errors.New("GeoRadiusByMember does not support Store or StoreDist"))
		return cmd
	***REMOVED***
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

// GeoRadiusByMemberStore is a writing GEORADIUSBYMEMBER command.
func (c cmdable) GeoRadiusByMemberStore(
	ctx context.Context, key, member string, query *GeoRadiusQuery,
) *IntCmd ***REMOVED***
	args := geoLocationArgs(query, "georadiusbymember", key, member)
	cmd := NewIntCmd(ctx, args...)
	if query.Store == "" && query.StoreDist == "" ***REMOVED***
		cmd.SetErr(errors.New("GeoRadiusByMemberStore requires Store or StoreDist"))
		return cmd
	***REMOVED***
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoSearch(ctx context.Context, key string, q *GeoSearchQuery) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 13)
	args = append(args, "geosearch", key)
	args = geoSearchArgs(q, args)
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoSearchLocation(
	ctx context.Context, key string, q *GeoSearchLocationQuery,
) *GeoSearchLocationCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 16)
	args = append(args, "geosearch", key)
	args = geoSearchLocationArgs(q, args)
	cmd := NewGeoSearchLocationCmd(ctx, q, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoSearchStore(ctx context.Context, key, store string, q *GeoSearchStoreQuery) *IntCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 15)
	args = append(args, "geosearchstore", store, key)
	args = geoSearchArgs(&q.GeoSearchQuery, args)
	if q.StoreDist ***REMOVED***
		args = append(args, "storedist")
	***REMOVED***
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoDist(
	ctx context.Context, key string, member1, member2, unit string,
) *FloatCmd ***REMOVED***
	if unit == "" ***REMOVED***
		unit = "km"
	***REMOVED***
	cmd := NewFloatCmd(ctx, "geodist", key, member1, member2, unit)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoHash(ctx context.Context, key string, members ...string) *StringSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(members))
	args[0] = "geohash"
	args[1] = key
	for i, member := range members ***REMOVED***
		args[2+i] = member
	***REMOVED***
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***

func (c cmdable) GeoPos(ctx context.Context, key string, members ...string) *GeoPosCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(members))
	args[0] = "geopos"
	args[1] = key
	for i, member := range members ***REMOVED***
		args[2+i] = member
	***REMOVED***
	cmd := NewGeoPosCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
***REMOVED***
