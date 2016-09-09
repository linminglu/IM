package main

import (
	//"log"
	//"time"
	"fmt"
	"time"

	"IM/common/json"
	"IM/common/rpc"
)

type UserService struct{}

func init() {
	registerService(new(UserService))
}

func (this *UserService) GetUserStatusKey(uid uint64) string {
	return fmt.Sprintf("im_user_status_key_%d", uid)
}

func (this *UserService) SetStatus(arg *rpc.ArgType, repley *rpc.ReplyType) error {
	obj := json.NewWithMap(arg.Args)
	uid, err := obj.Int("userId")
	if err != nil {
		return err
	}
	key := this.GetUserStatusKey(uint64(uid))
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	g_cache.Put(key, body, time.Hour)

	return nil
}

func (this *UserService) GetStatus(arg *rpc.ArgType, reply *rpc.ReplyType) error {
	obj := json.NewWithMap(arg.Args)
	uid, err := obj.Int("userId")
	if err != nil {
		return err
	}
	key := this.GetUserStatusKey(uint64(uid))

	body := g_cache.Get(key).([]byte)
	obj, err = json.Unmarshal(body)
	if err != nil {
		return err
	}
	reply.Reply = obj.Map()
	return nil
}
