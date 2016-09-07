package team

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/vmihailenco/redis/v2"

	"sirendaou.com/duserver/common"
)

var g_redis *common.RedisManager = nil
var g_msgredis *common.RedisManager = nil

func RedisInit(redisAddr string) int {
	g_redis = common.NewRedisManager(redisAddr)
	if g_redis == nil {
		return -1
	}
	return 0
}

func TeamMsgRedisInit(redisAddr string) int {
	g_msgredis = common.NewRedisManager(redisAddr)
	if g_msgredis == nil {
		return -1
	}
	return 0
}

func CacheCreate(team *common.TeamInfo) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	b, err := json.Marshal(*team)
	if err != nil {
		logger.Error("error:", err)
		return -1
	}

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAM_INFO, team.TeamId)
	rClient.Client.Set(userKey, string(b[0:]))

	userKey = fmt.Sprintf("%s%d", common.SET_USERS_TEAM, team.Uid)
	val := fmt.Sprintf("%d", team.TeamId)

	rClient.Client.SAdd(userKey, val)

	//info version
	val = fmt.Sprintf("%d", time.Now().Unix())
	userKey = fmt.Sprintf("%s%d", common.KEY_TEAM_INFO_VER, team.TeamId)
	rClient.Client.Set(userKey, val)

	return 0
}

func CacheDelete(team *common.TeamInfo) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAM_INFO, team.TeamId)
	rClient.Client.Del(userKey)
	logger.Info("redis delete ", userKey)

	userKey = fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, team.TeamId)

	uidList, err := rClient.Client.SMembers(userKey).Result()

	if err != nil {
		logger.Error("error:", err)
	} else {
		for _, val := range uidList {
			userKey2 := fmt.Sprintf("%s%s", common.SET_USERS_TEAM, val)
			val := fmt.Sprintf("%d", team.TeamId)
			rClient.Client.SRem(userKey2, val)
			logger.Info("redis srem ", userKey2, val)
		}
	}

	rClient.Client.Del(userKey)
	logger.Info("redis delete ", userKey)

	return 0
}

func CacheSetInfo(team *common.TeamInfo) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAM_INFO, team.TeamId)
	val, err := rClient.Client.Get(userKey).Result()

	if err != nil {
		logger.Error("error:", err)
		return -1
	}

	var t common.TeamInfo
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Error("error:", err)
		return -1
	}

	if len(team.TeamName) > 1 {
		t.TeamName = team.TeamName
	}
	if len(team.CoreInfo) > 1 {
		t.CoreInfo = team.CoreInfo
	}

	if len(team.ExInfo) > 1 {
		t.ExInfo = team.ExInfo
	}

	b, err := json.Marshal(t)
	if err != nil {
		logger.Error("error:", err)
		return -1
	}
	userKey = fmt.Sprintf("%s%d", common.KEY_TEAM_INFO, team.TeamId)
	rClient.Client.Set(userKey, string(b[0:]))
	logger.Info("set :", userKey, string(b[0:]))

	//info version
	val = fmt.Sprintf("%d", time.Now().Unix())
	userKey = fmt.Sprintf("%s%d", common.KEY_TEAM_INFO_VER, team.TeamId)
	rClient.Client.Set(userKey, val)

	return 0
}

func RedisQueryInfo(team *common.TeamInfo) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAM_INFO, team.TeamId)
	val, err := rClient.Client.Get(userKey).Result()

	if err != nil {
		logger.Error("error:", err)
		return -1
	}

	var t common.TeamInfo
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Error("error:", err)
		return -1
	}

	team.CoreInfo = t.CoreInfo
	team.ExInfo = t.ExInfo
	team.Uid = t.Uid
	team.MaxCount = common.MAX_MEMBER_NUM_TRAM
	team.TeamName = t.TeamName
	team.TeamType = t.TeamType

	return 0
}

func RedisQueryList(team *common.TeamInfo) (int, []int64) {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_USERS_TEAM, team.Uid)
	strTeamList, err := rClient.Client.SMembers(userKey).Result()

	logger.Debug("SMembers ", userKey, strTeamList)
	if err != nil {
		logger.Error("SMembers ", userKey, "error:", err)
		return -1, nil
	}

	teamList := make([]int64, len(strTeamList))

	for i, val := range strTeamList {
		teamList[i], _ = strconv.ParseInt(val, 10, 64)
	}

	return 0, teamList
}

// 系统预设的群组
func RedisQuerySysList() (int, []int64) {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	strTeamList, err := rClient.Client.SMembers(common.SET_SYS_TEAM).Result()
	if err != nil {
		logger.Error("sys team error:", err)
		return -1, nil
	}

	teamList := make([]int64, len(strTeamList))

	for i, val := range strTeamList {
		teamList[i], _ = strconv.ParseInt(val, 10, 64)
	}

	return 0, teamList
}

func RedisQueryMembers(team *common.TeamInfo) (int, []uint64) {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, team.TeamId)
	strUidList, err := rClient.Client.SMembers(userKey).Result()
	if err != nil {
		logger.Error("error:", err)
		return -1, nil
	}

	uidList := make([]uint64, len(strUidList))

	for i, val := range strUidList {
		uidList[i], _ = strconv.ParseUint(val, 10, 64)
	}

	return 0, uidList
}

func RedisIsMembers(teamId, uid uint64) bool {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, teamId)
	val := fmt.Sprintf("%d", uid)

	logger.Debug("SIsMember ", userKey, val)
	is, err := rClient.Client.SIsMember(userKey, val).Result()

	if err != nil {
		logger.Error("error:", err)
		return false
	}

	return is
}

func CacheAddMember(team *common.TeamInfo, uid uint64) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, team.TeamId)
	val := fmt.Sprintf("%d", uid)

	rClient.Client.SAdd(userKey, val)

	//member version
	val = fmt.Sprintf("%d", time.Now().Unix())
	userKey = fmt.Sprintf("%s%d", common.KEY_TEAM_MEMBER_VER, team.TeamId)
	rClient.Client.Set(userKey, val)

	userKey = fmt.Sprintf("%s%d", common.SET_USERS_TEAM, uid)
	val = fmt.Sprintf("%d", team.TeamId)

	rClient.Client.SAdd(userKey, val)

	return 0
}

func CacheScardMember(team *common.TeamInfo) int64 {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, team.TeamId)

	num := rClient.Client.SCard(userKey).Val()

	return num
}

func CacheRemoveMember(team *common.TeamInfo, uid uint64) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAM_MEMBER, team.TeamId)
	val := fmt.Sprintf("%d", uid)
	rClient.Client.SRem(userKey, val)

	//member version
	val = fmt.Sprintf("%d", time.Now().Unix())
	userKey = fmt.Sprintf("%s%d", common.KEY_TEAM_MEMBER_VER, team.TeamId)
	rClient.Client.Set(userKey, val)

	userKey = fmt.Sprintf("%s%d", common.SET_USERS_TEAM, uid)
	val = fmt.Sprintf("%d", team.TeamId)
	rClient.Client.SRem(userKey, val)

	return 0
}

func CacheAddWBMember(team *common.TeamInfo, uid uint64, Type int) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := ""
	if Type == 1 {
		userKey = fmt.Sprintf("%s%d", common.SET_WHITELIST, team.Uid)
	} else {
		userKey = fmt.Sprintf("%s%d", common.SET_BLACKLIST, team.Uid)
	}

	val := fmt.Sprintf("%d", uid)
	rClient.Client.SAdd(userKey, val)

	// 双向好友
	if Type == 1 {
		key_2 := fmt.Sprintf("%d", val)
		uid_2 := fmt.Sprintf("%d", uid)
		rClient.Client.SAdd(key_2, uid_2)
	}

	return 0
}

func CacheRemoveWBMember(team *common.TeamInfo, uid uint64, Type int) int {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := ""
	if Type == 1 {
		userKey = fmt.Sprintf("%s%d", common.SET_WHITELIST, team.Uid)
	} else {
		userKey = fmt.Sprintf("%s%d", common.SET_BLACKLIST, team.Uid)
	}

	uid_1 := fmt.Sprintf("%d", uid)

	logger.Debug("CacheRemoveWBMember",userKey, uid_1)

	rClient.Client.SRem(userKey, uid_1)

	// 双向好友
	if Type == 1 {
		userKey = fmt.Sprintf("%s%s", common.SET_WHITELIST, uid_1)
		uid_2 := fmt.Sprintf("%d", team.Uid)
		logger.Debug("CacheRemoveWBMember",userKey, uid_2)
		rClient.Client.SRem(userKey, uid_2)
	}

	return 0
}

func RedisQueryWBMembers(team *common.TeamInfo, Type int) (int, []uint64) {
	rClient := <-g_redis.RedisCh

	defer func() {
		g_redis.RedisCh <- rClient
	}()

	userKey := ""
	if Type == 1 {
		userKey = fmt.Sprintf("%s%d", common.SET_WHITELIST, team.Uid)
	} else {
		userKey = fmt.Sprintf("%s%d", common.SET_BLACKLIST, team.Uid)
	}

	strUidList, err := rClient.Client.SMembers(userKey).Result()
	if err != nil {
		logger.Error("error:", err)
		return -1, nil
	}

	if strUidList == nil || len(strUidList) < 1 {
		return 0, []uint64{}
	}

	uidList := make([]uint64, len(strUidList))

	for i, val := range strUidList {
		uidList[i], _ = strconv.ParseUint(val, 10, 64)
	}

	return 0, uidList
}

func CacheSetMsgBuf(msgId uint64, msg []byte) int {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAMMSGBUF, msgId)
	rClient.Client.SetEx(userKey, time.Second*3*86400, string(msg[:]))

	return 0
}

func CacheGetMsgBuf(msgId uint64) []byte {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAMMSGBUF, msgId)
	val, err := rClient.Client.Get(userKey).Result()
	if err != nil {
		logger.Error("error:", err)
		return nil
	}

	return []byte(val)
}

func CacheDelMsgBuf(msgId uint64, msg []byte) int {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.KEY_TEAMMSGBUF, msgId)
	rClient.Client.Del(userKey)

	return 0
}

func CacheAddMsgId(uids []uint64, msgId uint64, score float64) int {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	for _, uid := range uids {
		if uid <= 100000 {
			continue
		}

		userKey := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, uid)
		cnt, err := rClient.Client.ZCard(userKey).Result()
		if err != nil {
			cnt = 0
		}

		if cnt >= common.MAX_TEAM_MSG_PER {
			rClient.Client.ZRemRangeByRank(userKey, 0, 0).Result()
		}

		val := fmt.Sprintf("%d", msgId)
		_, err = rClient.Client.ZAdd(userKey, redis.Z{score, val}).Result()
		if err != nil {
			logger.Info("ZAdd ", userKey, val, score, " fail:", err)
		}
	}

	return 0
}

func CacheRemMsgId(uid, msgId uint64) {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, uid)
	member := fmt.Sprintf("%d", msgId)
	rClient.Client.ZRem(userKey, member)

	return
}

func CacheGetMsgIds(uid uint64) []uint64 {
	rClient := <-g_msgredis.RedisCh

	defer func() {
		g_msgredis.RedisCh <- rClient
	}()

	userKey := fmt.Sprintf("%s%d", common.SET_TEAMMSGID, uid)
	vals, err := rClient.Client.ZRange(userKey, 0, -1).Result()

	if err != nil {
		if err == redis.Nil {
			logger.Info("ZRange ", userKey, 0, -1, " not data")
		} else {
			logger.Info("ZRange ", userKey, 0, -1, "fail", err.Error())
		}

		return nil
	}

	uidList := make([]uint64, len(vals))
	for i, s := range vals {
		uidList[i], err = strconv.ParseUint(s, 10, 64)
	}

	return uidList
}