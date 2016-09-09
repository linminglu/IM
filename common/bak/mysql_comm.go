package common

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/mysql"
	"time"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"crypto/aes"
	"crypto/cipher"
	"sirendaou.com/duserver/db_server"
	"encoding/json"
)

type MysqlManager struct {
	dbCh chan *sql.DB
}

type TokenLoginBody struct {
	Uid uint64 `json:"uid,omitempty"`
	PlatformType string   `json:"platformtype,omitempty"`
	CoreToken string `json:"coretoken,omitempty"`
}

var g_mysql *MysqlManager = nil

func MysqlInit(host, dbname, user, paswd string) int {
	var settings = db.DataSource{
		Host:     host,
		Database: dbname,
		User:     user,
		Password: paswd,
	}

	//	num := 8
	num := 2
	dbCh := make(chan *sql.DB, num)

	for i := 0; i < num; i++ {
		sess, err := db.Open("mysql", settings)

		if err != nil {
			panic(err)
			return -1
		}

		drv := sess.Driver().(*sql.DB)

		sqlStr := fmt.Sprintf("update t_user_info set uid = uid where uid = 100000 ")
		_, err = drv.Exec(sqlStr)

		if err != nil {
			panic(err)
			return -2
		}

		dbCh <- drv
	}

	g_mysql = &MysqlManager{dbCh}

	return 0
}

//增 删 改
func ProcExec(sql string) error {
	logger.Debug("procExec sql=", sql)

	drv := <-g_mysql.dbCh

	logger.Debug(sql)
	_, err := drv.Exec(sql)

	if err != nil {
		logger.Error("err:", err)
	} else {
		logger.Debug("success:")
	}

	g_mysql.dbCh <- drv

	return err
}

// 查
func procSQL(sql string) *sql.Rows {
	logger.Debug("procSQL sql=", sql)

	drv := <-g_mysql.dbCh

	rows, err := drv.Query(sql)

	if err != nil {
		logger.Error("err:", err)
	} else {
		logger.Debug("success:")
	}

	g_mysql.dbCh <- drv

	return rows
}

//User System
func (user *UserInfo) DBInsertUser() error {
	logger.Debug("DBInsertUser")
	sqlStr := fmt.Sprintf("insert into t_user_info (id, uid, did, reg_date, update_date, baseinfo, exinfo, phonenum, password , platform) values (0 , %d, '%s', now(), now(), '', '', '%s', '%s','%s')",
		user.Uid, user.Did, user.PhoneNum, user.Pwd, user.Platform)

	return ProcExec(sqlStr)
}

func (user *UserInfo) DBInsertPhoneNum() error {
	sqlStr := fmt.Sprintf("insert into t_bind_phone (phonenum, uid, bind_date) values ('%s','%d',now())", user.PhoneNum, user.Uid)

	logger.Debug(sqlStr)

	return ProcExec(sqlStr)
}

func (user *UserInfo) DBUpdateUserInfo() int {
	sqlStr := ""
	if len(user.BaseInfo) > 1 && len(user.ExInfo) > 1 {
		sqlStr = fmt.Sprintf("UPDATE t_user_info SET did = '%s', baseinfo = '%s', exinfo = '%s', bv = %d, v = %d  WHERE uid = %d",
			user.Did, user.BaseInfo, user.ExInfo, user.BV, user.V, user.Uid)
	} else if len(user.BaseInfo) > 1 {
		sqlStr = fmt.Sprintf("UPDATE t_user_info SET did = '%s', baseinfo = '%s' ,  bv = %d WHERE uid = %d", user.Did, user.BaseInfo, user.BV, user.Uid)
	} else if len(user.ExInfo) > 1 {
		sqlStr = fmt.Sprintf("UPDATE t_user_info SET  exinfo = '%s' ,  v = %d WHERE uid = %d", user.ExInfo, user.V, user.Uid)
	} else {
		return 0
	}

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (user *UserInfo) DBInsertCS(image, tel, name, email string) int {
	sqlStr := fmt.Sprintf("insert into t_customservice (id, uid, account, password, phonenum, nick_name, image_id, email, tel ) values (0 , %d, '%s','%s', '%s', '%s', '%s', '%s','%s') on duplicate key update enable =1, del = 1",
		user.Uid, user.Pwd, user.PhoneNum, name, image, email, tel)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (user *UserInfo) DBUpdateCSEnable(enable int) int {
	sqlStr := fmt.Sprintf("update  t_customservice set enable = %d where uid = %d;", enable, user.Uid)

	err := ProcExec(sqlStr)

	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (user *UserInfo) DBUpdateCSDelete() int {
	sqlStr := fmt.Sprintf("update  t_customservice set del = 1 where uid = %d;", user.Uid)

	err := ProcExec(sqlStr)

	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (user *UserInfo) DBGetCSList(page, cnt int) (int, []CSInfo) {
	sqlStr := ""

	if len(user.PhoneNum) == 24 && user.Uid > 0 {
		sqlStr = fmt.Sprintf("select uid, account, password, Phonenum, nick_name, image_id, email, tel, enable from t_customservice  where uid = %d and del = 0;", user.Uid)
	} else if len(user.PhoneNum) == 24 {
		sqlStr = fmt.Sprintf("select uid, account, password, Phonenum, nick_name, image_id, email, tel, enable from t_customservice where Phonenum = '%s' and del = 0 and enable=1 limit %d , %d;",
			user.PhoneNum, page*cnt, cnt)
	} else {
		return 1, nil
	}

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)

	result := []CSInfo{}

	if rows != nil {
		for rows.Next() {
			var uid uint64 = 0
			var enable int = 0
			var account string = ""
			var password string = ""
			var Phonenum string = ""
			var nick_name string = ""
			var image_id string = ""
			var email string = ""
			var tel string = ""

			err := rows.Scan(&uid, &account, &password, &Phonenum, &nick_name, &image_id, &email, &tel, &enable)
			if err != nil {
				logger.Error(err)
			}
			/*type CSInfo struct {
			  	Uid     uint64
			  	Passwd  string
			  	Phonenum  string
			  	Account string
			  	Name    string
			  	Image   string
			  	Tel     string
			  	Email   string
			  	Enable  int
			  }
			*/
			cs := CSInfo{uid, password, Phonenum, account, nick_name, image_id, tel, email, enable, 0}

			result = append(result, cs)
		}

		rows.Close()
	}

	return 0, result
}

func (user *UserInfo) DBGetCSVerList(v int) (int, []CSVerInfo) {
	sqlStr := ""

	if len(user.PhoneNum) == 24 {
		sqlStr = fmt.Sprintf("select uid, UNIX_TIMESTAMP(reg_date) from t_customservice where phonenum = '%s' and UNIX_TIMESTAMP(reg_date) > %d;", user.PhoneNum, v)
	} else {
		return 1, nil
	}

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)

	result := []CSVerInfo{}

	if rows != nil {
		for rows.Next() {
			var uid uint64 = 0
			var v int = 0
			err := rows.Scan(&uid, &v)

			if err != nil {
				logger.Error(err)
			}

			cs := CSVerInfo{uid, v}

			result = append(result, cs)
		}

		rows.Close()
	}

	return 0, result
}

func GetUserInfoByUid(uid uint64) (*UserInfo, error) {
	sqlStr := fmt.Sprintf("select uid, password from t_user_info where uid=%d", uid)

	logger.Debug("sqlStr=", sqlStr)

	drv := <-g_mysql.dbCh
	defer func() {
		g_mysql.dbCh <- drv
	}()
	rows, err := drv.Query(sqlStr)
	if err != nil {
		logger.Error(sqlStr, " err:", err)
		return nil, err
	}
	defer rows.Close()
	if rows != nil {
		if rows.Next() {
			//			var id uint64 = 0
			var uid uint64 = 0
			//			var enable int = 0
			//			var del int = 0
			//			var v int = 0
			//			var account string = ""
			var password string = ""
			//			var phonenum string = ""
			//			var nick_name string = ""
			//			var image_id string = ""
			//			var email string = ""
			//			var tel string = ""

			err := rows.Scan(&uid, &password)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			//			if enable == 0 || del > 0 {
			//				enable = 0
			//			}
			userInfo := &UserInfo{
				Uid: uid,
				Pwd: password,
			}
			return userInfo, nil
		} else {
			return nil, sql.ErrNoRows
		}
	}

	return nil, err
}

func GetUserInfoByPhoneNum(phonenum string) (*UserInfo, error) {
	sqlStr := fmt.Sprintf("select uid, phonenum, password from t_user_info where phonenum=%s", phonenum)
	logger.Debug(sqlStr)

	drv := <-g_mysql.dbCh
	defer func() {
		g_mysql.dbCh <- drv
	}()

	rows, err := drv.Query(sqlStr)
	if err != nil {
		logger.Error(sqlStr, " err:", err)
		return nil, err
	}
	defer rows.Close()
	logger.Debug(rows)

	if rows.Next() {
		//			var id uint64 = 0
		var uid uint64 = 0
		//			var enable int = 0
		//			var del int = 0
		//			var v int = 0
		//			var account string = ""
		var password string = ""
		var phonenum string = ""
		//			var nick_name string = ""
		//			var image_id string = ""
		//			var email string = ""
		//			var tel string = ""

		err := rows.Scan(&uid, &phonenum, &password)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		//			if enable == 0 || del > 0 {
		//				enable = 0
		//			}
		userInfo := &UserInfo{
			Uid:      uid,
			Pwd:      password,
			PhoneNum: phonenum,
		}
		return userInfo, nil
	} else {
		return nil, sql.ErrNoRows
	}

	return nil, err
}

func GetUidByPhoneNum(phonenum string) (uid uint64, err error) {
	sqlStr := fmt.Sprintf("select uid from t_bind_phone where phonenum=%s", phonenum)
	logger.Debug("sqlStr=", sqlStr)

	drv := <-g_mysql.dbCh

	rows, err := drv.Query(sqlStr)
	if err != nil {
		logger.Error(sqlStr, " err:", err)
		return 0, err
	}
	defer rows.Close()

	g_mysql.dbCh <- drv

	for rows.Next() {
		var uid uint64 = 0
		err := rows.Scan(&uid)
		if err != nil {
			logger.Error(err)
			return 0, err
		}
		return uid, nil
	}

	return uid, nil
}

func (user *UserInfo) DBGetCSListByUidList(uids []uint64) (int, []CSInfo) {
	sqlStr := ""
	tempstr := ""
	s := ""
	for i, uid := range uids {
		tempstr = fmt.Sprintf("%d", uid)
		if i == 0 {
			s = s + tempstr
		} else {
			s = s + "," + tempstr
		}
	}

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)

	result := []CSInfo{}

	if rows != nil {
		for rows.Next() {
			var uid uint64 = 0
			var enable int = 0
			var del int = 0
			var v int = 0
			var account string = ""
			var password string = ""
			var phonenum string = ""
			var nick_name string = ""
			var image_id string = ""
			var email string = ""
			var tel string = ""

			err := rows.Scan(&uid, &account, &password, &phonenum, &nick_name, &image_id, &email, &tel, &enable, &del, &v)
			if err != nil {
				logger.Error(err)
			}

			if enable == 0 || del > 0 {
				enable = 0
			}
			cs := CSInfo{uid, password, phonenum, account, nick_name, image_id, tel, email, enable, v}

			result = append(result, cs)
		}

		rows.Close()
	}

	return 0, result
}

func (user *UserInfo) DBGetCSLastVer() int {
	sqlStr := ""
	if len(user.PhoneNum) == 24 {
		sqlStr = fmt.Sprintf("select  UNIX_TIMESTAMP(reg_date) from t_customservice where phonenum = '%s' order by reg_date desc limit 1;", user.PhoneNum)
	} else {
		return 0xffffffff
	}

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)
	var v int = 0
	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&v)
			if err != nil {
				logger.Error(err)
			}
		}
		rows.Close()
	}

	return v
}

func (user *UserInfo) DBCSLogin(phonenum, account string) (uint64, string) {
	sqlStr := fmt.Sprintf("select uid , password from t_customservice where phonenum = '%s' and account = '%s' and enable=1 and del = 0", phonenum, account)

	rows := procSQL(sqlStr)

	var uid uint64 = 0
	var passwd string = ""
	if rows != nil {
		for rows.Next() {

			err := rows.Scan(&uid, &passwd)
			if err != nil {
				logger.Error(err)
			}
		}
		rows.Close()
	}

	return uid, passwd
}

func DBUpdateUserInfo(uid uint64, passwd string) int {
	sqlStr := fmt.Sprintf("UPDATE t_user_info SET password = '%s' WHERE uid = %d", passwd, uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBUpdateSetupId(uid uint64, setupid uint64) int {
	sqlStr := fmt.Sprintf("UPDATE t_user_info SET setupid = %d WHERE uid = %d", setupid, uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBGetUserInfo(uidArr []uint64, propertyList []string) map[uint64]string {
	sqlStr := "SELECT"
	for _, key := range propertyList {
		sqlStr += fmt.Sprintf(" %s,", key)
	}
	sqlStr = strings.TrimRight(sqlStr, ",")
	sqlStr += " FROM t_user_info WHERE uid in ( 0 "
	for _, key := range uidArr {
		sqlStr += fmt.Sprintf(",%d", key)
	}
	sqlStr += ");"

	rows := procSQL(sqlStr)

	infoMap := make(map[uint64]string)
	if rows != nil {
		infos := make([]sql.RawBytes, len(propertyList))
		infoAddrs := make([]interface{}, len(propertyList))
		for i := range infos {
			infoAddrs[i] = &infos[i]
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(infoAddrs...)
			if err == nil {
				//info := fmt.Sprintf("\"uid\":%d,", uidArr[i])
				info := ""
				for k, v := range propertyList {
					if v == "uid" {
						info += fmt.Sprintf(" \"%s\":%s,", v, string(infos[k]))
					} else {
						info += fmt.Sprintf(" \"%s\":\"%s\",", v, string(infos[k]))
					}
				}
				info = strings.TrimRight(info, ",")
				infoMap[uint64(i)] = info
			} else {
				logger.Error(err)
			}
			i++
		}

		err := rows.Close()
		if err != nil {
			logger.Error("Close rows error.")
		} else {
			logger.Info(sqlStr, " success:")
		}
	}

	return infoMap
}

func DBUniInsertToken(uid uint64, token string) int {
	sqlStr := fmt.Sprintf("insert into t_devicetoken (id, uid, token) values (0 , %d, '%s') on duplicate key update token = '%s'", uid, token, token)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBClearExtraToken(uid uint64, token string) int {
	sqlStr := ""

	if uid > 0 {
		sqlStr = fmt.Sprintf("update t_devicetoken set token = '' where token = '%s' and uid != %d ", token, uid)
	} else {
		sqlStr = fmt.Sprintf("update t_devicetoken set token = '' where token = '%s'", token)
	}

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBGetUidByToken(token string) uint64 {
	sqlStr := fmt.Sprintf("select uid from t_devicetoken where token = '%d'", token)

	rows := procSQL(sqlStr)

	var uid uint64 = 0
	if rows != nil {
		for rows.Next() {
			var uid uint64 = 0

			err := rows.Scan(&uid)
			if err != nil {
				logger.Error(err)
			}
		}
		rows.Close()
	}

	return uid
}

func DBAddReport(uid uint64, msg string) int {
	sqlStr := fmt.Sprintf("insert into t_report (id, uid, msg) values (0 , %d, '%s')", uid, msg)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBCSUpdateLogoName(Phonenum, logo, name string) int {
	sqlStr := fmt.Sprintf("update t_apps set cs_logourl='%s', cs_name='%s' where app_key = '%s'", logo, name, Phonenum)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBCSGetLogoName(Phonenum string) (string, string) {
	sqlStr := fmt.Sprintf("select cs_logourl,  cs_name from t_apps where app_key = '%s';", Phonenum)

	rows := procSQL(sqlStr)
	var logo string = ""
	var name string = ""
	if rows != nil {
		for rows.Next() {

			err := rows.Scan(&logo, &name)
			if err != nil {
				logger.Error(err)
			}
		}

		rows.Close()
	}

	return logo, name
}

func DBCSInsertMsg(from, to, msg, phoneNum string, msgId int) int {
	sqlStr := fmt.Sprintf(`insert into t_csmsg_history (phonenum, fromcid, tocid, msg, msgid) value('%s','%s','%s','%s',%d)`, phoneNum, from, to, msg, msgId)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBCSQueryMsg(phoneNum, account string, page, page_size int) []string {
	sqlStr := fmt.Sprintf("select msg from t_csmsg_history where phonenum = '%s' and (fromcid = '%s' or tocid = '%s') order by id desc limit %d, %d;", phoneNum, account, account, page*page_size, page_size)

	rows := procSQL(sqlStr)
	var msg string = ""
	msgs := []string{}
	if rows != nil {
		for rows.Next() {

			err := rows.Scan(&msg)
			if err != nil {
				logger.Error(err)
			} else {
				msgs = append(msgs, msg)
			}
		}

		rows.Close()
	}

	return msgs
}

func DBCSQueryMsgByTime(phoneNum, account string, starttime, endtime int) []string {
	sqlStr := fmt.Sprintf("select msg from t_csmsg_history where phonenum = '%s' and (fromcid = '%s' or tocid = '%s') and unix_timestamp(itime) > %d and unix_timestamp(itime) < %d ;",
		phoneNum, account, account, starttime, endtime)

	rows := procSQL(sqlStr)
	var msg string = ""
	msgs := []string{}
	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&msg)
			if err != nil {
				logger.Error(err)
			} else {
				msgs = append(msgs, msg)
			}
		}
		rows.Close()
	}

	return msgs
}

func DBCSResetPwd(uid uint64, passwd string) int {
	sqlStr := fmt.Sprintf("update t_customservice set password='%s'  where uid = %d", passwd, uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
		return 1
	}

	return 0
}

//Team System
func (team *TeamInfo) DBCreate() int {
	sqlStr := fmt.Sprintf("insert into t_team_info (creater, teamid, name, type, maxnum, create_date, coreinfo, exinfo) values (%d, %d,'%s', %d , %d, now(), '','') on duplicate key update creater = creater", team.Uid, team.TeamId, team.TeamName, team.TeamType, team.MaxCount)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBDelete() int {
	sqlStr := fmt.Sprintf("update t_team_info set del_flag = 1 where teamid = %d and creater = %d", team.TeamId, team.Uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBSetInfo() int {
	sqlStr := "update t_team_info set  "

	isUpdate := false
	if len(team.CoreInfo) > 1 {
		teamStr := " coreinfo = '" + team.CoreInfo + "',"
		sqlStr += teamStr
		isUpdate = true
	}

	if len(team.ExInfo) > 1 {
		teamStr := " exinfo = '" + team.ExInfo + "',"
		sqlStr += teamStr
		isUpdate = true
	}

	if len(team.TeamName) > 1 {
		teamStr := " name = '" + team.TeamName + "',"
		sqlStr += teamStr
		isUpdate = true
	}

	teamStr := fmt.Sprintf(" maxnum = maxnum where teamid = %d", team.TeamId)
	sqlStr += teamStr

	if !isUpdate {
		return 1
	}

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBAddMember(uid uint64) int {
	sqlStr := fmt.Sprintf("insert into t_team_list (teamid, uid, itime) value (%d, %d, now()) on duplicate key update uid = uid", team.TeamId, uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBGetMembers() ([]uint64, error) {
	sqlStr := fmt.Sprintf("select uid from t_team_list where teamid = %d;", team.TeamId)
	rows := procSQL(sqlStr)
	var uid uint64
	uids := []uint64{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&uid)
			if err != nil {
				logger.Error(err)
				return nil, err
			} else {
				uids = append(uids, uid)
			}
		}
	}
	return uids, nil
}

// 获取好友列表
func (team *TeamInfo) DBGetFriendUids(myUid uint64) ([]uint64, error) {
	sqlStr := fmt.Sprintf("select fuid as uid from t_whitelist where uid = %d and del_flag=0;", myUid)
	rows := procSQL(sqlStr)
	var uid uint64
	uids := []uint64{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&uid)
			if err != nil {
				logger.Error(err)
				return nil, err
			} else {
				uids = append(uids, uid)
			}
		}
	}
	return uids, nil
}

func (team *TeamInfo) DBRemoveMember(uid uint64) int {
	sqlStr := fmt.Sprintf("delete from  t_team_list where teamid = %d and uid = %d", team.TeamId, uid)

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBGetNewTeamID() uint64 {
	sqlStr := fmt.Sprintf("select del_flag from t_team_info where creater=%d", team.Uid)

	drv := <-g_mysql.dbCh
	defer func() {
		g_mysql.dbCh <- drv
	}()

	rows, err := drv.Query(sqlStr)
	defer func() {
		if rows != nil {
			err = rows.Close()
		}
	}()

	count := 0
	okCount := 0
	delFlag := 1
	if err != nil {
		logger.Error(sqlStr, " err:", err)
		return 0
	} else {
		for rows.Next() {
			if err = rows.Scan(&delFlag); err != nil {
				logger.Error("Query error.")
				return 0
			}
			if delFlag == 0 {
				okCount++
			}
			count++
		}
	}

	logger.Info("uid:", team.Uid, " team num:", count, " ok num:", okCount)

	if okCount >= MAX_TEAM_NUM_PER {
		return 1
	}

	var tid uint64 = uint64(count) + 1
	//	tid = tid + (team.Uid & 0xffffffffff000000)
	tid = tid + (team.Uid * 1000)

	return tid
}

func (team *TeamInfo) DBWBAdd(Uid uint64, Type int) int {
	sqlStr := ""
	if Type == 1 {
		sqlStr = fmt.Sprintf("insert into t_whitelist (uid, fuid, itime, del_flag, last_modify_date) values(%d, %d, now(), 0, unix_timestamp()) on duplicate key update del_flag =0;",
			team.Uid, Uid)
	} else {
		sqlStr = fmt.Sprintf("insert into t_blacklist (uid, fuid, itime, del_flag, last_modify_date) values(%d, %d, now(), 0, unix_timestamp()) on duplicate key update del_flag =0;",
			team.Uid, Uid)
	}

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	// 双向好友
	if Type == 1 {
		sqlStr = fmt.Sprintf("insert into t_whitelist (uid, fuid, itime, del_flag, last_modify_date) values(%d, %d, now(), 0, unix_timestamp()) on duplicate key update del_flag =0;",
			Uid, team.Uid)
	} else {
		sqlStr = fmt.Sprintf("insert into t_blacklist (uid, fuid, itime, del_flag, last_modify_date) values(%d, %d, now(), 0, unix_timestamp()) on duplicate key update del_flag =0;",
			Uid, team.Uid)
	}

	err = ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func (team *TeamInfo) DBWBDelete(Uid uint64, Type int) int {
	sqlStr := ""
	if Type == 1 {
		sqlStr = fmt.Sprintf("update t_whitelist set del_flag = 1 where uid = %d and fuid = %d", team.Uid, Uid)
	} else {
		sqlStr = fmt.Sprintf("update t_blacklist set del_flag = 1 where uid = %d and fuid = %d", team.Uid, Uid)
	}

	err := ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	// 双向关系
	if Type == 1 {
		sqlStr = fmt.Sprintf("update t_whitelist set del_flag = 1 where uid = %d and fuid = %d", Uid, team.Uid)
	} else {
		sqlStr = fmt.Sprintf("update t_blacklist set del_flag = 1 where uid = %d and fuid = %d", Uid, team.Uid)
	}

	err = ProcExec(sqlStr)
	if err != nil {
		logger.Error(sqlStr, "err:", err)
	}

	return 0
}

func DBGetSysTeamList() (int, []TeamInfo) {
	sqlStr := "select teamid, name, maxnum, coreinfo, exinfo from t_team_info where is_sys=1;"

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)

	result := []TeamInfo{}

	if rows != nil {
		for rows.Next() {
			var teamid uint64 = 0
			var name string = ""
			var maxnum int = 0
			var coreinfo string = ""
			var exinfo string = ""

			err := rows.Scan(&teamid, &name, &maxnum, &coreinfo, &exinfo)
			if err != nil {
				logger.Error(err)
			}
			cs := TeamInfo{TeamId: teamid, TeamName: name, CoreInfo: coreinfo, ExInfo: exinfo, MaxCount: maxnum}

			result = append(result, cs)
		}
		rows.Close()
	}

	return 0, result
}

func DBGetSysTeamIdList() (int, []int64) {
	sqlStr := "select teamid from t_team_info where is_sys=1;"

	logger.Debug("sqlStr:", sqlStr)

	rows := procSQL(sqlStr)

	result := []int64{}
	if rows != nil {
		for rows.Next() {
			var teamid int64 = 0
			err := rows.Scan(&teamid)
			if err != nil {
				logger.Error(err)
			}
			result = append(result, teamid)
		}
		rows.Close()
	}

	return 0, result
}

func DBCreateMobileLoginToken(platformtype string, uid uint64, passwd string, setupid string, deviceid string) string {
	if uid == 0 {
		logger.Error("error: uid shouldn't be 0")
		return ""
	}

	if passwd == nil || passwd == "" {
		logger.Error("error: passwd shouldn't be nil or empty string")
		return ""
	}

	if setupid == nil {
		setupid = ""
	}

	if deviceid == nil {
		deviceid = ""
	}

	// 精确到纳秒
	t := time.Now().UnixNano()

	str := fmt.Sprintf("time:%lld,uid:%llu,psd:%s,setupid:%s,deviceid:%s,platformtype:%s", t, uid, passwd, setupid, deviceid, platformtype)

	h := md5.New()
	h.Write([]byte(str))

	coreToken := hex.EncodeToString(h.Sum(nil))

	tokenLoginBody := TokenLoginBody{
		Uid: int64(uid),
		PlatformType: platformtype,
		CoreToken:coreToken,
	}

	tokenLoginBodyBuf, err := json.Marshal(tokenLoginBody)

	// 密钥
	var keyText = "astaxie12798akljzmknm.ahkjkljl;k"
	var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	c, err := aes.NewCipher([]byte(keyText))
	if err != nil {
		logger.Error("NewCipher error:", err)
		return ""
	}

	// 加密
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	tokenEncrypted := make([]byte, len(tokenLoginBodyBuf))
	cfb.XORKeyStream(tokenEncrypted, tokenLoginBodyBuf)

	if strings.HasPrefix(platformtype, "p") || strings.HasPrefix(platformtype, "P") {
		sqlStr := fmt.Sprintf("UPDATE t_user_login_info SET password = '%s', pc_login_core_token = '%s', pc_login_token_encrypted = '%s', pc_setup_id = '%s', pc_device_id = '%s', pc_time = '%llu' WHERE uid = %llu", passwd, coreToken, tokenEncrypted, setupid, deviceid, t)

		err := ProcExec(sqlStr)

		if err != nil {
			logger.Error(sqlStr, "err:", err)
		}
	} else {
		platformtype = strings.ToLower(platformtype)
		platformtype = platformtype[0:1]

		sqlStr := fmt.Sprintf("UPDATE t_user_login_info SET password = '%s', mobile_login_core_token = '%s', mobile_login_token_encrypted = '%s', mobile_setup_id = '%s', mobile_device_id = '%s', mobile_time = '%llu', mobile_type = '%s' WHERE uid = %llu", passwd, coreToken, tokenEncrypted, setupid, deviceid, t, platformtype)

		err := ProcExec(sqlStr)

		if err != nil {
			logger.Error(sqlStr, "err:", err)
		}
	}

	return tokenEncrypted
}

func DBGetMobileCoreTokenWithUID(uid uint64) string {
	if uid == 0 {
		logger.Error("error: uid shouldn't be 0")
		return 0
	}

	sqlStr := fmt.Sprintf("select mobile_login_core_token where uid = %llu;", uid)
	rows := procSQL(sqlStr)

	var result string = ""

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&result)
			if err != nil {
				logger.Error(err)
				break
			}
		}

		rows.Close()
	}

	return result
}

func DBGetMobileTokenTimeWithUID(uid uint64) uint64 {
	if uid == 0 {
		logger.Error("error: uid shouldn't be 0")
		return 0
	}

	sqlStr := fmt.Sprintf("select mobile_time where uid = %llu;", uid)
	rows := procSQL(sqlStr)

	var result int64 = 0

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&result)
			if err != nil {
				logger.Error(err)
				break
			}
		}

		rows.Close()
	}

	return result
}
