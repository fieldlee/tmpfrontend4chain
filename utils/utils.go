package utils

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"frontend4chain/constant"
	"frontend4chain/db"
	"frontend4chain/module"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
	"github.com/twinj/uuid"
)

var h = md5.New()

func cipherEncode(sourceText string) string {
	h.Write([]byte("xoosdj$%&^@sda!@$sf,.?~"))
	cipherHash := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i := 0; i < loopCount; i++ {
		outData[i] = inputData[i] ^ cipherHash[i%32]
	}
	return fmt.Sprintf("%s", outData)
}

// 发送邮箱
func SentMail(to string, body string) bool {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "services@lianbai.io", "链佰BaaS系统管理员")
	m.SetHeader("To", m.FormatAddress(to, ""))
	m.SetHeader("Subject", "链佰BaaS系统用户名激活")
	m.SetBody("text/html", body)
	// gomail.Dialer{Host: "smtp.qq.com", Port: 465, Username: "1123123123@qq.com", Password: "ocozlhmwsvxhbidd",SSL:true}
	d := gomail.Dialer{Host: "smtp.exmail.qq.com", Port: 465, Username: "services@lianbai.io", Password: "P@ss1234", SSL: true}
	// d := gomail.NewPlainDialer("smtp.exmail.qq.com", 465, "depeng.li@lianbai.io", "FieldLee_0528") // 发送邮件服务器、端口、发件人账号、发件人密码
	err := d.DialAndSend(m)
	if err != nil {
		log.Println("发送失败", err.Error())
		return false
	}
	return true
}
func RT(code int, msg string, data string) module.JsonResult {
	rJson := module.JsonResult{}
	rJson.Code = code
	rJson.Msg = msg
	rJson.Data = data
	return rJson
}

// json convert to map
func ConvertToMap(obj interface{}) map[string]interface{} {
	var returnInterface map[string]interface{}
	rMarshal, _ := json.Marshal(obj)
	json.Unmarshal(rMarshal, &returnInterface)
	return returnInterface
}

// 读取random
func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for index := 0; index < l; index++ {
		if string(randInt(65, 90)) != temp {
			temp = string(randInt(65, 90))
			result.WriteString(temp)
		}
	}
	return result.String()
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func HlcEncode(sourceText string) string {
	h.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	noise := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText + "-" + strconv.Itoa(int(time.Now().Unix())))
	loopCount := len(inputData)
	outData := make([]byte, loopCount*2)
	for i, j := 0, 0; i < loopCount; i, j = i+1, j+1 {
		outData[j] = noise[i%32]
		j++
		outData[j] = inputData[i] ^ noise[i%32]
	}
	return base64.StdEncoding.EncodeToString([]byte(cipherEncode(fmt.Sprintf("%s", outData))))
}

func HlcDecode(sourceText string) string {
	buf, err := base64.StdEncoding.DecodeString(sourceText)
	if err != nil {
		fmt.Println("Decode(%q) failed: %v", sourceText, err)
		return ""
	}
	inputData := []byte(cipherEncode(fmt.Sprintf("%s", buf)))
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i, j := 0, 0; i < loopCount; i, j = i+2, j+1 {
		outData[j] = inputData[i] ^ inputData[i+1]
	}
	str := fmt.Sprintf("%s", outData)
	arr := strings.Split(str, "-")
	if len(arr) < 2 {
		return ""
	}
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	result := strings.Join(arr[:len(arr)-1], "")
	return fmt.Sprintf("%s", result)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// 删除重复的slice
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if Contains(ret, a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// 删除slice 某个元素
func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	} else if (length - 1) >= index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
	}
	return nil, errors.New("error")
}

// 创建token
func CreateToken(username string, subject string) (string, error) {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix() - 3600),
		ExpiresAt: int64(time.Now().Unix() + 60*60*24*7),
		Issuer:    username,
		Subject:   subject,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(constant.PUBLICKEY))
}

// 验证token
func VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(constant.PUBLICKEY), nil
	})
}

// email 验证码token
func CreateEmailToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix() - 3600),
		ExpiresAt: int64(time.Now().Unix() + 60*60*6),
		Issuer:    username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constant.PUBLICKEY))
}

// 运行shell 命令
func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

// 生成uuid
func GetUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.FormatHex)
}

// 验证文件是否存在，如存在删除
func CheckFileAndRemove(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(file)
	return err
}

// 检查目录是否存在 不存在生成
func CheckAndCreatePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// 清空证书路径下的证书文件
func ClearCertficationPath(id string) error {
	newpath := filepath.Join(constant.ROOTPATH, id)
	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		return nil
	}
	err := os.RemoveAll(newpath)
	return err
}

// 生成路径
func CreateCertficationPath(id string) (string, error) {
	newpath := filepath.Join(constant.ROOTPATH, id)
	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		err = os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return newpath, nil
}

//项目默认值修改
func Fill(info module.Define) module.Define {
	// 补全共识
	if len(info.Orders) > 1 {
		info.Consensus = "kafka"
	} else {
		info.Consensus = "solo"
	}
	info.OrderId = "order"
	info.OrderName = "order"
	// info.Status = constant.SAVEED

	info.NetWork = info.Domain[:strings.Index(info.Domain, ".")]
	// orders
	var i = 0
	// orders 最多4个
	// 获得所有的ip
	orderIps := make([]string, 0)
	for _, order := range info.Orders {
		// 默认第一个order 为kafkaip
		if i == 0 {
			info.KafkaIp = order.OrderIp
			info.Orders[i].OrderPort = 7050
		} else {
			var replaceNum = 0
			for _, ip := range orderIps {
				if ip == order.OrderIp {
					replaceNum = replaceNum + 1
				}
			}

			info.Orders[i].OrderPort = 7050 + 1000*replaceNum //不同的端口 根据ip 重复的次数
		}
		orderIps = append(orderIps, order.OrderIp)
		info.Orders[i].OrderId = fmt.Sprint("order", i)
		info.Orders[i].OrderName = fmt.Sprint("order", i)
		info.Orders[i].ContainerId = fmt.Sprint(info.Orders[i].OrderId, ".", info.Domain)
		i = i + 1
	}
	// 循环info
	// 获得所有的ip
	peerIps := make([]string, 0)
	i = 0
	for _, org := range info.Orgs {
		info.Orgs[i].OrgName = org.OrgId
		info.Orgs[i].PeerNumber = len(org.Peers) //计算peernumber
		var j = 0
		for _, peer := range org.Peers {
			var replaceNum = 0
			for _, ip := range peerIps {
				if ip == peer.PeerIp {
					replaceNum = replaceNum + 1
				}
			}
			info.Orgs[i].Peers[j].PeerId = fmt.Sprint("peer", j)
			info.Orgs[i].Peers[j].PostPort = 7051 + 1000*replaceNum  //不同的端口 根据ip 重复的次数
			info.Orgs[i].Peers[j].EventPort = 7053 + 1000*replaceNum //不同的端口 根据ip 重复的次数
			info.Orgs[i].Peers[j].ContainerId = fmt.Sprint(info.Orgs[i].Peers[j].PeerId, ".", org.OrgId, ".", info.Domain)
			// 默认开启couchdb
			info.Orgs[i].Peers[j].JoinCouch = true
			info.Orgs[i].Peers[j].CouchId = fmt.Sprint("couch", "_", org.OrgId, "_", fmt.Sprint("peer", j))
			info.Orgs[i].Peers[j].CouchUsername = constant.COUCHUSERNAME
			info.Orgs[i].Peers[j].CouchPassword = constant.COUCHPASSWORD
			info.Orgs[i].Peers[j].CouchPort = 4984 + 1000*replaceNum //不同的端口 根据ip 重复的次数
			info.Orgs[i].Peers[j].CouchContainerId = fmt.Sprint("couch.", fmt.Sprint("peer", j), ".", org.OrgId, ".", info.Domain)
			if j == 0 {
				info.Orgs[i].AnchorIp = peer.PeerIp
				info.Orgs[i].AnchorPort = 7051 + 1000*replaceNum
				info.Orgs[i].CaIp = peer.PeerIp
				info.Orgs[i].CaId = "ca" //默认ca 为caid
				info.Orgs[i].CaPort = 7054 + 1000*replaceNum
				info.Orgs[i].ContainerId = fmt.Sprint("ca.", info.Orgs[i].OrgId, ".", info.Domain)
				info.Orgs[i].Peers[j].EventPort = 7053 + 1000*replaceNum //默认端口
				info.Orgs[i].Peers[j].PostPort = 7051 + 1000*replaceNum  //默认端口
			}
			peerIps = append(peerIps, peer.PeerIp)
			j = j + 1
		}
		i = i + 1
	}
	return info
}

// 是否存在文件
func IsExistFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

// 根据路径创建文件
func CreateFile(path string) (interface{}, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// 返回json
func ResponseJson(code int, msg string, data string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding,Authorization,X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Method", "GET,HEAD,PUT,PATCH,POST,DELETE")
	w.Header().Set("X-Requested-With", "XmlHttpRequest")
	returnData := module.JsonResult{}
	returnData.Code = code
	returnData.Msg = msg
	returnData.Data = data
	jsonData, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println("请检查返回数据格式是否为标准interface,error:", err.Error())
		w.Write([]byte(`{"code":"500","msg":"系统内部错误！","data":[]}`))
		return
	}
	w.Write(jsonData)
	return
}

// random string
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Md5(decodeString string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(decodeString))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func IsManager(request *http.Request) (bool, string, error) {
	curusername := request.Header.Get("username")
	curuser, err := db.GetUser(curusername)
	if err != nil {
		return false, "", err
	}
	if curuser == nil {
		return false, "", nil
	}
	if strings.TrimSpace(curuser.(module.User).Role) == "manager" {
		return true, curusername, nil
	} else {
		return false, curusername, nil
	}
}
