package pkg

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	mrand "math/rand"
)


func ObjToJson(msg interface{}) string {
	jso, error := json.Marshal(msg)
	if error == nil {
		return string(jso)
	}
	return ""
}


type Body struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Status = -1
		body.Msg = err.Error()
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

func LoginResponse(w http.ResponseWriter, code int, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Status = code
		body.Msg = err.Error()
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

// ResponseData ...
type ResponseData struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	RequestID interface{} `json:"request_id"`
	Data      interface{} `json:"data"`
}

// HTTPPost ...
func HTTPPost(url string, data []byte) (string, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// HTTPPost ...
func HTTPPostFormData(url string, data url.Values) (string, error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// HTTPGet ...
func HTTPGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Println("-===========---", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("response data :", string(body))
	return string(body), nil
}

// JSONMarshal ...
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// ErrorMsg ...
func ErrorMsg(m interface{}) error {
	msg, _ := JSONMarshal(m)
	return errors.New(string(msg))
}

// GetRequestID ...
func GetRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logx.Error(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// MD5 ...
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

// ObjToJSON ...
func ObjToJSON(msg interface{}) string {
	jso, error := json.Marshal(msg)
	if error == nil {
		return string(jso)
	}
	return ""
}

// ObjToJSONRelay ...
func ObjToJSONRelay(msg interface{}) (r []byte) {
	jso, error := json.Marshal(msg)
	if error == nil {
		return jso
	}
	return r
}

// GetCurrentIP ...
func GetCurrentIP(r *http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]

	}
	return ip
}

// ReadFile 读取文件
func ReadFile(path string) string {
	// os 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		logx.Error(err)
	}
	//os.Stdout.Write(data)
	return string(data)
}

//  AES CBC 加密方式

// Padding 对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// UnPadding 对密文删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

func GetRandomNumString(l int) string {
	str := "0123456789abcdefghijklmnpqrstuvmxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetTimeDuration  获取开始时间和结束时间耗费毫秒数
func GetTimeDuration(start time.Time) string {
	duration := time.Now().Sub(start)
	return duration.String()
}

// Interface2Str 接口转string
func Interface2Str(params ...interface{}) string {
	var paramSlice []string
	for _, param := range params {
		switch v := param.(type) {
		case string:
			paramSlice = append(paramSlice, v)
		case int:
			strV := strconv.FormatInt(int64(v), 10)
			paramSlice = append(paramSlice, strV)
		default:
			paramSlice = append(paramSlice, fmt.Sprintf("%v", v))
		}
	}

	res := strings.Join(paramSlice, " ")
	return res
}

func GetUID(r *http.Request) (int64, error) {
	// Get the value of "foo" as a json.Number
	num, ok := r.Context().Value("id").(json.Number)
	if !ok {
		return 0, errors.New("用户未登录")
	}
	// Convert the json.Number to int64
	UID, err := num.Int64()
	if err != nil {
		return 0, err
	}
	return UID, nil
}

// 获取本周的起始时间
// 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func WeekIntervalTime(week int) (startTime, endTime string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
	endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

	return startTime,endTime
}


// PadLeft 函数将整数转换为指定位数的字符串，并在数值前面补充指定字符
func PadLeft(num int, width int, padChar byte) string {
	return fmt.Sprintf("%0*d", width, num)
}


// HttpGetWithHeaders ...
func HttpGetWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
