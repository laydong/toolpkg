package utils

import (
	"bytes"
	"container/list"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg"
	"github.com/oschwald/geoip2-golang"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	RequestIdKey  = "request_id" // 日志key
)

//获取用户IP地址
func GetClientIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// Find获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func Find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// GetRandomString 获取随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

// GetRandomString6 获取6位随机字符串
func GetRandomString6(n uint64) []byte {
	baseStr := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	base := []byte(baseStr)
	quotient := n
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		mod = quotient % 34
		quotient = quotient / 34
		l.PushFront(base[int(mod)])
	}
	listLen := l.Len()
	if listLen >= 6 {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {
			res = append(res, i.Value.(byte))
		}
		return res
	} else {
		res := make([]byte, 0, 6)
		for i := 0; i < 6; i++ {
			if i < 6-listLen {
				res = append(res, base[0])
			} else {
				res = append(res, l.Front().Value.(byte))
				l.Remove(l.Front())
			}

		}
		return res
	}
}

// GenValidateCode 生成6位随机验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// CreateOrder 生成订单号
func CreateOrder() int64 {
	return int64(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// GetAddressByIP 获取省市区通过ip
func GetAddressByIP(ipA string) string {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipA)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	var province string
	if len(record.Subdivisions) > 0 {
		province = record.Subdivisions[0].Names["zh-CN"]
	}

	return record.Country.Names["zh-CN"] + "-" + province + "-" + record.City.Names["zh-CN"]
}

// InSliceString string是否在[]string里面
func InSliceString(k string, s []string) bool {
	for _, v := range s {
		if k == v {
			return true
		}
	}
	return false
}

// Exists 判断文件或目录是否存在
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func IsNil(obj interface{}) bool {
	type eFace struct {
		data unsafe.Pointer
	}
	if obj == nil {
		return true
	}
	return (*eFace)(unsafe.Pointer(&obj)).data == nil
}

// Base64URLDecode 因为Base64转码后可能包含有+,/,=这些不安全的URL字符串，所以要进行换字符
//'+' -> '-'
//'/' -> '_'
//'=' -> ''
//字符串长度不足4倍的位补"="
func Base64URLDecode(data string) string {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing) //字符串长度不足4倍的位补"="
	data = strings.Replace(data, "_", "/", -1)
	data = strings.Replace(data, "-", "+", -1)
	return data
}

func Base64UrlSafeEncode(data string) string {
	safeUrl := strings.Replace(data, "/", "_", -1)
	safeUrl = strings.Replace(safeUrl, "+", "-", -1)
	safeUrl = strings.Replace(safeUrl, "=", "", -1)
	return safeUrl
}

func Base64Encode(s string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(s))
	return encodeString
}

func Base64Decode(code string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		log.Fatalln(err)
	}
	return string(decodeBytes)
}

// GenerateTraceId 获取链路TraceId
func GenerateTraceId() string {
	return Md5(uuid.NewV4().String())
}

// GetRequestIdKey 获取链路ID
func GetRequestIdKey(c *gin.Context) (requestId string) {
	requestId = c.GetHeader(toolpkg.XtraceKey)
	if requestId != "" {
		c.Set(toolpkg.RequestIdKey, requestId)
	}
	requestId = c.GetString(toolpkg.RequestIdKey)
	if requestId == "" {
		requestId = GenerateTraceId()
		c.Set(toolpkg.RequestIdKey, requestId)
	}
	return
}

// Md5 md5
func Md5(s string) string {
	m := md5.Sum([]byte(s))
	return hex.EncodeToString(m[:])
}

// Merge options(latter ones have higher priority)
func MergeOptions(options ...map[int]interface{}) map[int]interface{} {
	rst := make(map[int]interface{})

	for _, m := range options {
		for k, v := range m {
			rst[k] = v
		}
	}

	return rst
}

// Merge headers(latter ones have higher priority)
func MergeHeaders(headers ...map[string]string) map[string]string {
	rst := make(map[string]string)

	for _, m := range headers {
		for k, v := range m {
			rst[k] = v
		}
	}

	return rst
}

// Add params to a url string.
func AddParams(url_ string, params url.Values) string {
	if len(params) == 0 {
		return url_
	}

	if !strings.Contains(url_, "?") {
		url_ += "?"
	}

	if strings.HasSuffix(url_, "?") || strings.HasSuffix(url_, "&") {
		url_ += params.Encode()
	} else {
		url_ += "&" + params.Encode()
	}

	return url_
}

func ToUrlValues(v interface{}) url.Values {
	switch t := v.(type) {
	case url.Values:
		return t
	case map[string][]string:
		return url.Values(t)
	case map[string]string:
		rst := make(url.Values)
		for k, v := range t {
			rst.Add(k, v)
		}
		return rst
	case nil:
		return make(url.Values)
	default:
		panic("Invalid value")
	}
}

// GetString 只能是map和slice
func GetString(d interface{}) string {
	bytesD, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(bytesD)
	}
}

func GetLocalIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

//@description: 将数组格式化为字符串
//@param: array []interface{}
//@return: string

func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

//@description: 利用反射将结构体转化为map
//@param: obj interface{}
//@return: map[string]interface{}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Tag.Get("mapstructure") != "" {
			data[obj1.Field(i).Tag.Get("mapstructure")] = obj2.Field(i).Interface()
		} else {
			data[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
	return data
}

// Strval interface转string
func Strval(value interface{}) string {

	var key string

	if value == nil {

		return key

	}

	switch value.(type) {

	case float64:

		ft := value.(float64)

		key = strconv.FormatFloat(ft, 'f', -1, 64)

	case float32:

		ft := value.(float32)

		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)

	case int:

		it := value.(int)

		key = strconv.Itoa(it)

	case uint:

		it := value.(uint)

		key = strconv.Itoa(int(it))

	case int8:

		it := value.(int8)

		key = strconv.Itoa(int(it))

	case uint8:

		it := value.(uint8)

		key = strconv.Itoa(int(it))

	case int16:

		it := value.(int16)

		key = strconv.Itoa(int(it))

	case uint16:

		it := value.(uint16)

		key = strconv.Itoa(int(it))

	case int32:

		it := value.(int32)

		key = strconv.Itoa(int(it))

	case uint32:

		it := value.(uint32)

		key = strconv.Itoa(int(it))

	case int64:

		it := value.(int64)

		key = strconv.FormatInt(it, 10)

	case uint64:

		it := value.(uint64)

		key = strconv.FormatUint(it, 10)

	case string:

		key = value.(string)

	case []byte:

		key = string(value.([]byte))

	default:

		newValue, _ := json.Marshal(value)

		key = string(newValue)

	}

	return key

}

func CheckParamsType(v interface{}) int {
	switch v.(type) {
	case url.Values, map[string][]string, map[string]string:
		return 1
	case []byte, string, *bytes.Reader:
		return 2
	case nil:
		return 0
	default:
		return 3
	}
}

func ToReader(v interface{}) *bytes.Reader {
	switch t := v.(type) {
	case []byte:
		return bytes.NewReader(t)
	case string:
		return bytes.NewReader([]byte(t))
	case *bytes.Reader:
		return t
	case nil:
		return bytes.NewReader(nil)
	default:
		panic("Invalid value")
	}
}

// Does the params contain a file?
func CheckParamFile(params url.Values) bool {
	for k, _ := range params {
		if k[0] == '@' {
			return true
		}
	}

	return false
}

// Add a file to a multipart writer.
func AddFormFile(writer *multipart.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)

	return err
}
