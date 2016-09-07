package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/vmihailenco/redis/v2"
)

const (
	EXPIRE_TIME = 864000
)

const (
	REDIS_UID_POOL = "du_uid_pool"
)

type RedisItem struct {
	Client    *redis.Client
	RedisAddr string
}

func CreateRedisItem(addr string) *RedisItem {
	client := redis.NewTCPClient(&redis.Options{
		Addr: addr,
	})

	//_, err := client.Ping().Result()
	err := client.Set("test", "test").Err()

	if err != nil {
		client.Close()
		return nil
	}

	return &RedisItem{client, addr}
}

func ReconnectRedis(pclient *RedisItem) int {
	pclient.Client.Close()

	logger.Error("reConnect  Redis:", pclient.RedisAddr)

	pclient.Client = redis.NewTCPClient(&redis.Options{
		Addr: pclient.RedisAddr,
	})

	err := pclient.Client.Set("test", "test").Err()

	if err != nil {
		pclient.Client.Close()
		logger.Error("reConnect  failed:", pclient.RedisAddr, err.Error())
		return -1
	}

	return 0
}

//func GetItemFromPoolTimeout(pool chan *RedisItem) *RedisItem {
//	timeOutCh := make(chan int, 1)
//
//	go func(ch chan int) {
//		time.Sleep(1e9 * 2)
//		ch <- 1
//		close(ch)
//	}(timeOutCh)
//
//	select {
//	case e := <-pool:
//		return e
//	case <-timeOutCh:
//		logger.Error("get item from pool timeout")
//	}
//
//	return nil
//}

const (
	SET_USERS_TEAM      = "userteam_"
	SET_SYS_TEAM        = "systeam_"
	SET_TEAM_MEMBER     = "teammember_"
	KEY_TEAM_INFO       = "teaminfo_"
	KEY_TEAM_INFO_VER   = "tinfov_"
	KEY_TEAM_MEMBER_VER = "tmemberv_"
	SET_WHITELIST       = "WHITE_"
	SET_BLACKLIST       = "BLACK_"
	KEY_TEAMMSGBUF      = "TEAMMSGBUF_"
	SET_TEAMMSGID       = "TEAMMSGID_"
)

type RedisManager struct {
	RedisCh chan *RedisItem
}

func NewRedisManager(redisAddr string) *RedisManager {
	//	logger.Debug("NewRedisManager redisAddr", redisAddr)
	addrs := strings.Split(redisAddr, ",")

	poolsize := 10 * len(addrs)

	AliasRediPool := make(chan *RedisItem, poolsize)

	for _, addr := range addrs {
		for i := 0; i < 10; i++ {
			it := CreateRedisItem(addr)
			if it == nil {
				logger.Error("aliase redis Connect  failed", addr)
				//attrreport.AttrSet(398, 1)
			} else {
				AliasRediPool <- it
			}
		}
	}

	return &RedisManager{AliasRediPool}
}

func (redisMgr *RedisManager) RedisSet(key string, value string) int {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	_, err := rClient.Client.Set(key, value).Result()

	if err != nil {
		logger.Error("redis set ", key, value, " err :", err)
		return -1
	}

	logger.Info("redis set ", key, value, " success. ")

	return 0
}

func (redisMgr *RedisManager) RedisDel(key string) int {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	_, err := rClient.Client.Del(key).Result()

	if err != nil {
		logger.Error("redis del ", key, " err :", err)
		return -1
	}

	logger.Info("redis del ", key, " success. ")

	return 0
}

func (redisMgr *RedisManager) RedisSetEx(key string, dur time.Duration, value string) int {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	_, err := rClient.Client.SetEx(key, dur, value).Result()
	if err != nil {
		logger.Error("redis set key: ", key, " value: ", value, " err :", err)
		return -1
	}

	return 0
}

func (redisMgr *RedisManager) RedisGet(key string) (string, int) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	v, err := rClient.Client.Get(key).Result()

	if err != nil && err != redis.Nil {
		logger.Error("redis get ", key, " err :", err)
		return "", -1
	}

	if err == redis.Nil {
		return "", 0
	} else {
		return v, 0
	}
}

func (redisMgr *RedisManager) RedisHGet(key string, field string) (string, int) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	v, err := rClient.Client.HGet(key, field).Result()

	if err != nil && err != redis.Nil {
		logger.Error("redis get ", key, " ", field, " err :", err)
		return "", -1
	}

	return v, 0
}

func (redisMgr *RedisManager) RedisMGet(keys []string) ([]interface{}, int) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	varr, err := rClient.Client.MGet(keys...).Result()

	if err != nil && err != redis.Nil {
		logger.Error("redis mget err :", err)
		return varr, -1
	}

	return varr, 0
}

func (redisMgr *RedisManager) RedisRPop(key string) string {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	v, err := rClient.Client.RPop(key).Result()

	if err != nil {
		return ""
	}

	return v
}

func (redisMgr *RedisManager) RedisLPush(key, val string) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	rClient.Client.LPush(key, val)

	return
}

func (redisMgr *RedisManager) RedisSAdd(key, val string) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	rClient.Client.SAdd(key, val)

	return
}

func (redisMgr *RedisManager) RedisSDel(key, val string) {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	rClient.Client.SRem(key, val)

	return
}

func (redisMgr *RedisManager) PipelineGetString(keys []string) []string {
	rClient := <-redisMgr.RedisCh

	pipeline := rClient.Client.Pipeline()

	keysNum := 0
	for _, key := range keys {
		if len(key) > 2 {
			logger.Debug("pipeline key ", key)
			pipeline.Get(key)
			keysNum++
		}
	}

	cmds, err := pipeline.Exec()

	logger.Debug("pipe result:", cmds, err)

	result := ""

	valList := make([]string, keysNum)
	n := 0
	if err != nil && err != redis.Nil {
		logger.Debug("redisClient pipeline err %s", err.Error())

		if ReconnectRedis(rClient) != 0 {
			logger.Error("ReconnectRedis fail!!!!!!")
		} else {
			redisMgr.RedisCh <- rClient
		}
	} else {
		redisMgr.RedisCh <- rClient
		logger.Debug("redisClient pipeline ok")

		for _, cmd := range cmds {
			logger.Debug(cmd, " ret result:", cmd.(*redis.StringCmd).Val())
			if cmd.(*redis.StringCmd).Err() != nil {
				logger.Error(cmd.(*redis.StringCmd).Err())
			} else {
				result = cmd.(*redis.StringCmd).Val()
			}
			if len(result) > 0 {
				logger.Debug(n, result)
				valList[n] = result
				n++
			}
		}
	}

	return valList
}

func (redisMgr *RedisManager) RedisZRange(key string) []string {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	userKey := key

	val, err := rClient.Client.ZRange(userKey, 0, -1).Result()

	if err != nil {
		if err == redis.Nil {
			logger.Info("ZRange ", userKey, 0, -1, " not data")
		} else {
			logger.Info("ZRange ", userKey, 0, -1, "fail", err.Error())
		}

		return nil
	}

	return val
}

func (redisMgr *RedisManager) RedisZRange2(key string, cnt int) []string {
	if cnt < 1 {
		return nil
	}

	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	userKey := key

	val, err := rClient.Client.ZRange(userKey, 0, int64(cnt-1)).Result()

	if err != nil {
		if err == redis.Nil {
			logger.Info("ZRange ", userKey, 0, cnt-1, " not data")
		} else {
			logger.Info("ZRange ", userKey, 0, cnt-1, "fail", err.Error())
		}

		return nil
	}

	return val
}

func (redisMgr *RedisManager) RedisZRem(key, val string) int {
	rClient := <-redisMgr.RedisCh

	defer func() {
		redisMgr.RedisCh <- rClient
	}()

	_, err := rClient.Client.ZRem(key, val).Result()

	if err == nil {
		logger.Info("ZRem ", key, val, "success")
	} else {
		logger.Info("ZRem ", key, val, " fail")
	}
	return 0
}

func (redisMgr *RedisManager) RedisStatCacheSet(tail InnerPkgTail) error {
	uidStr := strconv.FormatUint(tail.ToUid, 10)
	key := "du_stat_" + uidStr

	str, err := json.Marshal(tail)
	if err != nil {
		logger.Error("RedisStatCacheSet err=", err)
		return err
	}

	errCode := redisMgr.RedisSet(key, string(str))
	if errCode != 0 {
		return errors.New("RedisGet error errCode")
	}

	return err
}

func (redisMgr *RedisManager) RedisStatCacheGet(uid uint64) (*InnerPkgTail, error) {
	uidStr := strconv.FormatUint(uid, 10)
	key := "du_stat_" + uidStr

	str, errCode := redisMgr.RedisGet(key)
	if errCode != 0 || str == "" {
		logger.Error("RedisGet error key=", key)
		return nil, errors.New("RedisGet error")
	}

	tail := &InnerPkgTail{}
	err := json.Unmarshal([]byte(str), tail)
	if err != nil {
		logger.Error("str:", str)
		logger.Error(err)
		return nil, err
	}

	return tail, nil
}

func (redisMgr *RedisManager) RedisMsgCacheSet(uid uint64, val string) error {
	uidStr := strconv.FormatUint(uid, 10)
	key := "du_msg_cache_" + uidStr
	errCode := redisMgr.RedisSet(key, val)

	if errCode != 0 {
		logger.Error("RedisGet error key=", key)
		return  errors.New("RedisGet error")
	}

	return nil
}

func (redisMgr *RedisManager) RedisMsgCacheGet(uid uint64) (string, error) {
	uidStr := strconv.FormatUint(uid, 10)
	key := "du_msg_cache_" + uidStr
	str, errCode := redisMgr.RedisGet(key)

	if errCode != 0 || str == "" {
		logger.Error("RedisGet error key=", key)
		return "", errors.New("RedisGet error")
	}
	return str, nil
}

func CacheCheckWBMember(redisClient *redis.Client, uid uint64, checkUid uint64, Type int) bool {
	userKey := ""
	if Type == 1 {
		userKey = fmt.Sprintf("%s%d", SET_WHITELIST, uid)
	} else {
		userKey = fmt.Sprintf("%s%d", SET_BLACKLIST, uid)
	}

	val := fmt.Sprintf("%d", checkUid)

	result := redisClient.SIsMember(userKey, val)
	if result == nil {
		if Type == 1 {
			return true
		} else {
			return false
		}
	}

	return result.Val()
}

func (redisMgr *RedisManager) RedisSetupIDCacheGet() (string, error) {
	key := "du_setupid_cache"
	str, errCode := redisMgr.RedisGet(key)

	if errCode != 0 || str == "" {
		errCode = redisMgr.RedisSet(key, "1")

		if errCode != 0 {
			logger.Error("RedisSet error key=", key)
			return "", errors.New("RedisSet error")
		}

		return "1", nil
	}

	val, errCode := strconv.Atoi(str);

	if errCode != 0 {
		logger.Error("strconv.Atoi error str=", str)
		return "", errors.New("Atoi error")
	}

	val++
	str = strconv.Itoa(val)
	errCode = redisMgr.RedisSet(key, str)

	if errCode != 0 {
		logger.Error("RedisSet error key=", key)
		return "", errors.New("RedisSet error")
	}

	return str, nil
}
