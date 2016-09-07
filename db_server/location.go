package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type RegLoc struct {
	Xpos  float64 `json:"xpos,omitempty"`
	Ypos  float64 `json:"ypos,omitempty"`
	Level uint16  `json:"level,omitempty"`
	Hour  uint32  `json:"hour,omitempty"`
	Page  uint16  `json:"page,omitempty"`
}

type LocResp struct {
	UserLocations []common.LocResult `json:"userlocations"`
}

func (h *DBHandler) RequestUserLocation(head common.PkgHead, jsonBody []byte, tail *common.InnerPkgTail) ([]byte, uint32) {
	var req RegLoc
	err := json.Unmarshal(jsonBody, &req)

	if err != nil {
		logger.Error("Unmarshal error:", err)
		return []byte(""), common.ERROR_CLIENT_BUG
	}

	sendTime := int(time.Now().Unix())

	logger.Info("Uid ", head.Uid, " xpos ", req.Xpos, " ypos ", req.Ypos, " time ", sendTime, " level ", req.Level, " hour ", req.Hour, " page ", req.Page)

	locInfo := &common.LocationInfo{head.Uid, []float64{req.Xpos, req.Ypos}, uint32(sendTime)}

	result := 0
	if req.Page == 0 {
		result = locInfo.SaveLocation()
		if result != 0 {
			return []byte(""), common.ERR_CODE_SYS
		}
	}

	var locArray []common.LocResult
	var saveBuf bytes.Buffer

	result = locInfo.GetLocationInfo()
	if result != 0 {
		return []byte(""), common.ERR_CODE_SYS
	}

	resp := LocResp{}

	redisKey := fmt.Sprintf("location_%d", head.Uid)
	if req.Level != 0 {
		if req.Page == 0 {
			locArray = locInfo.GetLocation(req.Level, req.Hour, req.Page)
			if locArray == nil {
				return []byte(""), common.ERR_CODE_SYS
			}
			logger.Debug("locArray size:", len(locArray))
			enc := gob.NewEncoder(&saveBuf)
			err := enc.Encode(locArray)
			if err != nil {
				return []byte(""), common.ERR_CODE_SYS
			}

			val := string(saveBuf.Bytes()[:])
			logger.Info("redis set local key:", redisKey, "len :", len(val))
			h.LocRedis.RedisSetEx(redisKey, 600*time.Second, val)

			if len(locArray) < common.MAX_ROW {
				resp.UserLocations = locArray
			} else {
				resp.UserLocations = locArray[:common.MAX_ROW]
			}
		} else {
			res, errcode := h.LocRedis.RedisGet(redisKey)
			if errcode != 0 {
				return []byte(""), common.ERR_CODE_SYS
			}

			logger.Info("redis local key:", redisKey, "len :", len(res))

			saveBuf2 := bytes.NewBufferString(res)
			dec := gob.NewDecoder(saveBuf2)
			err = dec.Decode(&locArray)
			if err != nil {
				logger.Error("decode error 1:", err)
				return []byte(""), common.ERR_CODE_SYS
			}

			logger.Debug("localarry size:", len(locArray))

			if len(locArray) <= int(req.Page)*common.MAX_ROW {
				return []byte("{\"userlocations\":[]}"), 0
			} else if len(locArray) <= int(req.Page+1)*common.MAX_ROW {
				resp.UserLocations = locArray[int(req.Page)*common.MAX_ROW:]
			} else {
				resp.UserLocations = locArray[req.Page*common.MAX_ROW : (req.Page+1)*common.MAX_ROW]
			}
		}

		respBuf, _ := json.Marshal(resp)
		return respBuf, 0
	}

	return []byte(""), 0
}
