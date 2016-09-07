package file_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"testing"

	qio "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

type RespUpload struct {
	Errcode int    `json:"errcode"`
	Token   string `json:"token"`
	Key     string `json:"key"`
	Url     string `json:"url"`
}

func TestUploadFile(t *testing.T) {
	val := make(url.Values)
	val.Add("fromuid", "100057")
	val.Add("hash", "hash_tdsestls_25431245431")
	val.Add("vercode", "12345678912345678912345678912345")
	resp, err := http.PostForm("http://192.168.20.51:8889/file/upload", val)
	//resp, err := http.PostForm("http://localhost:8889/file/upload", val)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))
	respMap := &RespUpload{}
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		fmt.Println("--------0", err)
	}
	uploadPost(respMap.Url, respMap.Token, respMap.Key)
}

func uploadPost(urlStr, token, key string) {
	body := new(bytes.Buffer)
	multiWrite := multipart.NewWriter(body)
	multiWrite.WriteField("token", token)
	multiWrite.WriteField("key", key)
	part, _ := multiWrite.CreateFormFile("file", "abc")
	file, err := os.Open("./abc")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	io.Copy(part, file)
	multiWrite.Close()

	req, err := http.NewRequest("POST", urlStr, body)
	fmt.Println("requst body:\n", body.String())
	if err != nil {
		fmt.Println("NewRequst failed:", err)
		return
	}
	req.Header.Add("Host", "upload.qiniu.com")
	req.Header.Add("Content-Type", fmt.Sprint("multipart/form-data; boundary=", multiWrite.Boundary()))
	req.Header.Add("Content-Length", fmt.Sprint(body.Len()))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("DO failed:", err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll failed:", err)
		return
	}
	fmt.Println("Post Ret:", string(respBody))
	m := make(map[string]string)
	json.Unmarshal(respBody, &m)
	fmt.Println(m)

	url := downloadUrl("7xlnqu.com2.z0.glb.qiniucdn.com", m["fid"])
	fmt.Println("url:", url)
	resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println("==========:", string(respBody))
}

func upload(token, key string) {
	localFile := "./abc"

	var ret qio.PutRet
	var extra = &qio.PutExtra{
	//Params:    params,
	//MimeType:  mieType,
	//Crc32:     crc32,
	//CheckCrc:  CheckCrc,

	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识
	// localFile 为本地文件名
	// extra     为上传文件的额外信息，详情见 io.PutExtra，可选    localFile := "./abc"
	err := qio.PutFile(nil, &ret, token, key, localFile, extra)

	if err != nil {
		//上传产生错误
		fmt.Println("io.PutFile failed:", err)
		return

	}
	//上传成功，处理返回值
	fmt.Println(ret.Hash, ret.Key)

	url := downloadUrl("7xlmc9.com1.z0.glb.clouddn.com", key)
	fmt.Println("url:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
}

func downloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)

}
