package common

import (
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Around System
const (
	LOCAL_MAX_PAGE = 30
	MAX_ROW        = 28 //每页数据最大值
)

type MongoManager struct {
	sessCh chan *mgo.Session
}

var g_mongo *MongoManager = nil

func MongoInit(mongodbAddr string) int {
	num := 2
	sessCh := make(chan *mgo.Session, num)

	for i := 0; i < num; i++ {
		sess, err := mgo.Dial(mongodbAddr)
		if err != nil {
			logger.Error("mongo init fail :", err)
			return -1
		}

		//Optional. Switch the session to a monotonic behavior.
		sess.SetMode(mgo.Monotonic, true)
		sessCh <- sess
	}

	g_mongo = &MongoManager{sessCh}

	return 0
}

type LocationInfo struct {
	Uid uint64 `bson:"uid"`
	//	AppKey      string    `bson:"appkey"`
	//	DeveloperID int       `bson:"developerID"`
	Loc      []float64 `bson:"loc"`
	SendTime uint32    `bson:"sendtime"`
}

type LocResult struct {
	Uid  uint64  `json:"uid"`
	Xpos float64 `json:"xpos"`
	Ypos float64 `json:"ypos"`
}

func (location *LocationInfo) NewLocation() int {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du").C("location")

	index := mgo.Index{
		Key:        []string{"$2d:loc"},
		Bits:       26,
		Background: true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		logger.Error("mongodb setup index failed")
		panic(err)
		return -1
	}

	newLoc := bson.M{
		"uid":      location.Uid,
		"loc":      location.Loc,
		"sendtime": location.SendTime,
	}

	logger.Debug("new location: ", newLoc)

	err = c.Insert(newLoc)
	if err != nil {
		logger.Error("Insert mongodb failed")
		return -1
	}

	return 0
}

func (location *LocationInfo) GetLocationInfo() int {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du").C("location")

	selector := bson.M{"uid": location.Uid}

	iter := c.Find(selector).Iter()

	for iter.Next(&location) {
		logger.Info("Uid: ", location.Uid, " AppKey: ", " Loc: ", location.Loc, " send time: ", location.SendTime)
	}

	if err := iter.Close(); err != nil {
		logger.Error("fail to close search.")
		return -1
	}

	return 0
}

func (location *LocationInfo) SaveLocation() int {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du").C("location")

	selector := bson.M{"uid": location.Uid}
	newLoc := bson.M{"uid": location.Uid, "loc": location.Loc, "sendtime": location.SendTime}
	update := bson.M{"$set": newLoc}
	logger.Debug("selector: ", selector, " update: ", update)
	err := c.Update(selector, update)
	if err != nil {
		logger.Error("Update mongodb failed")
		return -1
	}

	return 0
}

func (location *LocationInfo) GetLocation(level uint16, hour uint32, page uint16) []LocResult {
	if level != 1 {
		logger.Error("level ", level, " not completed.")
		return nil
	}

	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du").C("location")

	selector := bson.M{"loc": bson.M{"$near": location.Loc}}
	logger.Debug("selector: ", selector)
	iter := c.Find(selector).Limit(MAX_ROW*LOCAL_MAX_PAGE + 1).Iter()
	result := LocationInfo{}

	retLoc := make([]LocResult, MAX_ROW*30)
	line := 0

	retLoc[line].Uid = location.Uid
	retLoc[line].Xpos = float64(int64(location.Loc[0]*1000000)) / 1000000
	retLoc[line].Ypos = float64(int64(location.Loc[1]*1000000)) / 1000000
	logger.Debug(retLoc[line].Xpos, retLoc[line].Ypos)
	line++

	for iter.Next(&result) && line < MAX_ROW*30 {
		if result.Uid == location.Uid {
			continue
		}

		logger.Debug("Uid: ", result.Uid, " Loc: ", result.Loc, " send time: ", result.SendTime, "line:", line)

		retLoc[line].Uid = result.Uid
		retLoc[line].Xpos = float64(int64(result.Loc[0]*1000000)) / 1000000
		retLoc[line].Ypos = float64(int64(result.Loc[1]*1000000)) / 1000000
		//retlocl[line].Xpos = result.Loc[0]
		//retlocl[line].Ypos = result.Loc[1]

		line++
	}

	return retLoc[:line]
}

//IM System
func (userMsg *UserMsgItem) SaveUserMsg(msg string) error {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du-msg").C("msg")

	logger.Info("mongdb save user msg. Msgid:", userMsg.MsgId, " Fromuid:", userMsg.FromUid, " Touid:", userMsg.ToUid)

	//userMsg.Content = msg

	return c.Insert(userMsg)
}

func (userMsg *UserMsgItem) DelUserMsg() error {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du-msg").C("msg")

	logger.Info("mongdb delete user msg. Msgid:", uint64(userMsg.MsgId), " Fromuid:", userMsg.FromUid, " Touid:", userMsg.ToUid)
	bs := bson.M{"msgid": userMsg.MsgId, "touid": userMsg.ToUid}

	logger.Debug("bson:", bs)

	return c.Remove(bs)
}

func (userMsg *UserMsgItem) GetUserMsg(uid uint64) ([]*UserMsgItem, error) {
	session := <-g_mongo.sessCh

	defer func() {
		g_mongo.sessCh <- session
	}()

	c := session.DB("du-msg").C("msg")

	selector := bson.M{"touid": uid}

	iter := c.Find(selector).Iter()

	msgs := []*UserMsgItem{}

	msg := new(UserMsgItem)
	for iter.Next(&msg) {
		msgs = append(msgs, msg)
		msg = new(UserMsgItem)
	}

	if err := iter.Close(); err != nil {
		logger.Error("fail to close search.")
		return nil, err
	}

	return msgs, nil
}
