package rest_server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"

	"sirendaou.com/duserver/common"
)

type Team struct {
	TeamId    uint64 `json:"teamid,omitempty"`
	MaxUsers  int    `json:"maxusers,omitempty"`
	GroupName string `json:"groupname,omitempty"`
	Desc      string `json:"desc,omitempty"`
}

type ListResp struct {
	State int       `json:"state"`
	Msg   string    `json:"msg"`
	Teams [](*Team) `json:"grouplist,omitempty"`
}

func (h *Handler) ListSysTeam(res http.ResponseWriter, req *http.Request) {
	errCode, teamList := common.DBGetSysTeamList()
	if errCode != 0 || teamList == nil || len(teamList) == 0 {
		logger.Error("group list is empty")
		res.Write([]byte(`{"state":500,"msg":"server error"}`))
		return
	}

	teams := [](*Team){}
	for _, t := range teamList {
		team := &Team{TeamId: t.TeamId, MaxUsers: t.MaxCount, GroupName: t.TeamName, Desc: t.CoreInfo}
		teams = append(teams, team)
	}

	resp := &ListResp{State: 0, Msg: "ok", Teams:teams}
	result, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err)
		res.Write([]byte(`{"state":500,"msg":"server error"}`))
		return
	} else {
		res.Write(result)
		return
	}

	return
}

func (h *Handler) CreateSysTeam(res http.ResponseWriter, req *http.Request) {
	maxCount := 100

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		res.Write([]byte(`{"state":400,"msg":"server error"}`))
		return
	}

	logger.Debug("rest api users req:", string(body))

	js, err := simplejson.NewJson(body)
	if err != nil {
		res.Write([]byte(`{"state":400,"msg":"server error"}`))
		logger.Error(err)
		return
	}

	arr, err := js.Array()
	if err != nil || len(arr) > maxCount {
		res.Write([]byte(`{"state":400,"msg":"request json error or json array too long"}`))
		logger.Error(err)
		return
	}

	teamInfo := &common.TeamInfo{Uid: 1}

	tid := teamInfo.DBGetNewTeamID()
	logger.Debug("tid=", tid)
	if tid == 0 {
		logger.Error("get team id error")
		res.Write([]byte(`{"state":500,"msg":"server erro"}`))
		return
	}

	//	req json : '[{"maxusers":5000, "groupname":"team name", "desc": "coreinfo"}, {"maxusers":5000, "groupname":"team name", "desc": "coreinfo"}]'
	strSql := "insert into t_team_info (creater, teamid, name, maxnum, coreinfo, is_sys, create_date) values "
	for i := 0; i < len(arr); i++ {
		maxusers, _ := js.GetIndex(i).Get("maxusers").Uint64()
		groupname, _ := js.GetIndex(i).Get("groupname").String()
		coreinfo, _ := js.GetIndex(i).Get("desc").String()
		s := fmt.Sprintf("(1, %d, '%s',%d, '%s',  1, now()), ", tid+uint64(i), groupname, maxusers, coreinfo)
		strSql += s
	}

	strSql = strSql[:len(strSql)-2]

	logger.Debug("strSql=", strSql)

	err = common.ProcExec(strSql)
	if err != nil {
		res.Write([]byte(`{"state":401,"msg":"create team fail"}`))
		logger.Error(err)
		return
	}

	res.Write([]byte(`{"state":0,"msg":"ok"}`))

	return
}
