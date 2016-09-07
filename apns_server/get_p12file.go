package apns_server

//import (
//	"database/sql"
//	"fmt"
//	"github.com/anachronistic/apns"
//	"github.com/gosexy/db"
//	_ "github.com/gosexy/db/mysql"
//	"os"
//	"strconv"
//)
//
//type MysqlManager struct {
//	dbCh chan *sql.DB
//}
//
//var g_mysql *MysqlManager = nil
//
//func MysqlInit(host, dbname, user, paswd string, port int) int {
//	var settings = db.DataSource{
//		Host:     host,
//		Database: dbname,
//		User:     user,
//		Password: paswd,
//		Port:     port,
//	}
//
//	num := 8
//	dbCh := make(chan *sql.DB, num)
//
//	for i := 0; i < num; i++ {
//		sess, err := db.Open("mysql", settings)
//		if err != nil {
//			panic(err)
//			return -1
//		}
//
//		drv := sess.Driver().(*sql.DB)
//		dbCh <- drv
//	}
//	g_mysql = &MysqlManager{dbCh}
//
//	return 0
//}
//
//func GetPem(appkey string, node int) (pem1, pem2 string) {
//	sqlStr := ""
//	filename := ""
//	if node == 1 {
//		sqlStr = fmt.Sprintf("select apple_data_test , certificate_pass_test from  t_apps where app_key = '%s'", appkey)
//		filename = "/tmp/" + appkey + "_test.p12"
//	} else {
//		sqlStr = fmt.Sprintf("select apple_data , certificate_pass from  t_apps where app_key = '%s'", appkey)
//		filename = "/tmp/" + appkey + ".p12"
//	}
//
//	fmt.Println(sqlStr)
//
//	drv := <-g_mysql.dbCh
//
//	defer func() {
//		g_mysql.dbCh <- drv
//	}()
//
//	rows, err := drv.Query(sqlStr)
//
//	if err != nil {
//		fmt.Println(sqlStr, err)
//		return "", ""
//	}
//	var data []byte
//	var pass string
//
//	if err != nil {
//		fmt.Println(sqlStr, " err:", err)
//	} else {
//		for rows.Next() {
//			if err = rows.Scan(&data, &pass); err != nil {
//				fmt.Println("Query error.", err)
//				break
//			}
//		}
//
//		fmt.Println(len(data), pass)
//		err = rows.Close()
//
//		if err != nil {
//			fmt.Println("Close rows error.")
//		} else {
//			fmt.Println(sqlStr, " succ:")
//		}
//	}
//
//	if len(data) > 100 {
//		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
//		if err != nil {
//			panic("open file failed")
//			return "", ""
//		}
//		defer f.Close()
//
//		f.Write(data)
//	}
//
//	return appkey, pass
//}
//
//var (
//	//mysql -h192.168.250.27 -udevelop -pdevelop_KKpush developerdb
//	DB_USER3   = "develop"
//	DB_PASSWD3 = "develop_KKpush"
//	DB_HOST3   = "192.168.250.27"
//	DB_NAME3   = "developerdb"
//	DB_POER3   = 3306
//
//	DB_USER2   = "im_read"
//	DB_PASSWD2 = "IM_read_123"
//	DB_HOST2   = "210.14.141.249"
//	DB_NAME2   = "portaldb"
//	DB_POER2   = 4306
//)
//
//func main() {
//	if len(os.Args) < 3 {
//		fmt.Println("please set appkey, modle")
//		fmt.Println(os.Args[0], "[appkey] [mode(1-sandbox 2-product)]")
//		return
//	}
//
//	ret := MysqlInit(DB_HOST2, DB_NAME2, DB_USER2, DB_PASSWD2, DB_POER2)
//	if ret != 0 {
//		fmt.Print("MysqlInit fail")
//		return
//	}
//
//	payload := apns.NewPayload()
//	payload.Alert = "Hello, world!"
//	payload.Badge = 42
//	payload.Sound = "bingbong.aiff"
//
//	pn := apns.NewPushNotification()
//	pn.DeviceToken = "dcafeab37727d87beae12e758b6a5c148570de6daf32400d971aa70988b5eaa3"
//	pn.AddPayload(payload)
//
//	appKey := os.Args[1]
//	mode, _ := strconv.Atoi(os.Args[2])
//	//p1, p2 := GetPem("00b6413a92d4c1c84ad99e0a", 1)
//	p1, p2 := GetPem(appKey, mode)
//	fmt.Printf(p1, p2)
//	return
//
//	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "YOUR_CERT_PEM", "YOUR_KEY_NOENC_PEM")
//	resp := client.Send(pn)
//
//	alert, _ := pn.PayloadString()
//	fmt.Println("  Alert:", alert)
//	fmt.Println("Success:", resp.Success)
//	fmt.Println("  Error:", resp.Error)
//}
