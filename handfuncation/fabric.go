package handfuncation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend4chain/constant"
	"frontend4chain/db"
	"frontend4chain/module"
	"frontend4chain/utils"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func HandlerFabircAll(writer http.ResponseWriter, request *http.Request) {

	switch request.URL.Path {
	case "/fabric/deleteproject":
		deleteHandler(writer, request)
		return
	case "/fabric/checkip":
		checkipHandler(writer, request)
		return
	case "/fabric/saveprojectpassword":
		passwordHandler(writer, request)
		return
	case "/fabric/saveproject":
		defineHandle(writer, request, constant.SAVEED)
		return
	case "/fabric/resaveproject":
		defineHandle(writer, request, constant.RESAVEED)
		return
	case "/fabric/deployenv": //合并证书生成 ，yaml 文件生成，以及发送文件
		deployHandle(writer, request)
		return
	case "/fabric/checkenv": //检查linux 环境
		checkEnvHandle(writer, request)
		return
	case "/fabric/installenv": //install linux 环境
		installEnvHandle(writer, request)
		return
	case "/fabric/generatechan":
		channelgenHandle(writer, request)
		return
	case "/fabric/installchaincode":
		chaincodeHandler(writer, request) //chaincode handle
		return
	case "/fabric/uploadchaincode":
		uploadChaincode(writer, request) //chaincode handle
		return
	case "/fabric/addorg":
		addOrg(writer, request)
		return
	case "/fabric/setexplorer":
		setExplorer(writer, request) //设置用户名和密码
		return
	case "/fabric/projectservers":
		fixProject(writer, request) //项目的服务器信息
		return
	case "/fabric/serverinfo":
		fixIp(writer, request) //fix project info
		return
	case "/fabric/dockeraction":
		fixDocker(writer, request) //fix project info
		return
	case "/fabric/topology":
		getTopology(writer, request) //获得网络拓扑图
		return
	case "/fabric/allblocktx":
		getBlockTx(writer, request) //获得项目所有的channel的blocks和transcations
		return
	case "/fabric/chatblocktx":
		getChatBlockTx(writer, request) //获得项目所有Chat的channel的blocks和transcations
		return
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}

}

// 检查服务器是否安装其他项目
func checkipHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得所有Chat的Block和Trascations参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	if obj["ip"] != nil {
		dockerip := obj["ip"].(string)
		msgchan := make(chan []byte)
		go func(ip string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					msgchan <- []byte("{'code':400}")
					return
				}
			}()
			// 请求服务器，
			host := fmt.Sprint("http://", ip, ":", constant.SERVERPORT)
			perUrl := fmt.Sprint(host, constant.CHECKIP)
			req, err := http.NewRequest("POST", perUrl, bytes.NewBuffer([]byte("")))
			if err != nil {
				fmt.Println("shell -- b")
				panic(err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", utils.HlcEncode(constant.CURIP))
			var client http.Client
			res, err := client.Do(req)
			if err != nil {
				fmt.Println("shell -- c")
				panic(err)
			}
			defer res.Body.Close()
			message, _ := ioutil.ReadAll(res.Body)
			msgchan <- message
		}(dockerip)
		// 等待check 结果
		for {
			select {
			case ipcheckinfo := <-msgchan:
				msgMap := make(map[string]interface{})
				err = json.Unmarshal(ipcheckinfo, &msgMap)
				if int(msgMap["code"].(float64)) != 200 {
					utils.ResponseJson(400, "服务器已经配置其他项目，请选择其他服务器！", "", writer)
					return
				}
				utils.ResponseJson(200, "服务器可用", "", writer)
				return
			}
		}

	} else {
		utils.ResponseJson(400, "请输入IP", "", writer)
		return
	}
}

// 修改项目密码
func passwordHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得所有Chat的Block和Trascations参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	project := module.Define{}
	rev, err := db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	if obj["password"] != nil {
		project.ProjectPassword = utils.Md5(obj["password"].(string))
		err = db.SaveSetupInfo(project, project.ID, rev)
		if err != nil {
			utils.ResponseJson(400, "保持项目错误："+err.Error(), "", writer)
			return
		}
	} else {
		utils.ResponseJson(400, "请设置项目密码", "", writer)
		return
	}
	utils.ResponseJson(200, "项目密码已经保持", "", writer)
	return
}

// 获得项目所有的channel的blocks和transcations FOR chat
func getChatBlockTx(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得所有Chat的Block和Trascations参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	project := module.Define{}
	_, err = db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}
	// 没有生产智能合约，返回错误
	if project.HasApp != true {
		utils.ResponseJson(400, "合约未部署", "{'tx':{},'block':{}}", writer)
		return
	}
	// 获得order 的ip地址
	ipOrder := project.Orders[0].OrderIp
	msgchan := make(chan []byte)
	go func(ip string) {
		defer func() { // 必须要先声明defer，否则不能捕获到panic异常
			if err := recover(); err != nil {
				fmt.Println(err) // 这里的err其实就是panic传入的内容，55
				msgchan <- []byte("{\"block\":{},\"tx\":{}}")
			}
		}()
		url := fmt.Sprint("http://", ip, ":4000/blockchat")
		fmt.Println("url:", url)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
		if err != nil {
			fmt.Println("shell -- b")
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		var client http.Client
		client.Timeout = time.Second * 5
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("shell -- c")
			panic(err)
		}
		defer res.Body.Close()
		message, _ := ioutil.ReadAll(res.Body)
		msgchan <- message
	}(ipOrder)
	for {
		select {
		case txinfo := <-msgchan:
			utils.ResponseJson(200, "获得块和交易完成", string(txinfo), writer)
			return
		}
	}

}

// 获得项目所有的channel的blocks和transcations
func getBlockTx(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得Block和Trascations参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	project := module.Define{}
	_, err = db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}
	// 没有生产智能合约，返回错误
	if project.HasApp != true {
		utils.ResponseJson(400, "合约未部署", "{'tx':0,'block':0}", writer)
		return
	}
	// 获得order 的ip地址
	ipOrder := project.Orders[0].OrderIp
	msgchan := make(chan []byte)
	go func(ip string) {
		defer func() { // 必须要先声明defer，否则不能捕获到panic异常
			if err := recover(); err != nil {
				fmt.Println(err) // 这里的err其实就是panic传入的内容，55
				msgchan <- []byte("{\"blockHeight\": 0,\"txNum\":0}")
			}
		}()
		url := fmt.Sprint("http://", ip, ":4000/blocktxnum")
		fmt.Println("url:", url)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
		if err != nil {
			fmt.Println("shell -- b")
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		var client http.Client
		client.Timeout = time.Second * 5
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("shell -- c")
			panic(err)
		}
		defer res.Body.Close()
		message, _ := ioutil.ReadAll(res.Body)
		msgchan <- message
	}(ipOrder)
	for {
		select {
		case txinfo := <-msgchan:

			fmt.Println("txinfo:", string(txinfo))
			txobj := make(map[string]interface{})
			err := json.Unmarshal(txinfo, &txobj)
			if err != nil {
				utils.ResponseJson(400, "", "读取交易和块错误："+err.Error(), writer)
				return
			}
			txNum := 0
			blockNum := 0
			for key, value := range txobj {
				if reflect.TypeOf(value).Kind() == reflect.Float64 {
					intV := int(value.(float64))
					if key == "blockHeight" {
						blockNum = blockNum + intV
					}
					if key == "txNum" {
						txNum = txNum + intV
					}
				}
			}
			respMap := make(map[string]interface{})
			respMap["tx"] = txNum
			respMap["block"] = blockNum
			responseByte, _ := json.Marshal(respMap)

			utils.ResponseJson(200, "获得块和交易完成", string(responseByte), writer)
			return
		}
	}

}

// delete project
func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得project参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	if obj["id"] == nil {
		utils.ResponseJson(500, "请输入项目id", "", writer)
		return
	}

	project := module.Define{}
	rev, err := db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	isManager, curusername, err := utils.IsManager(request)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "当前用户错误："+err.Error(), "", writer)
		return
	}
	// 角色和人员验证
	if curusername != project.Manager && !isManager {
		utils.ResponseJson(400, "当前用户不是项目的所有人，或者不是管理员", "", writer)
		return
	}

	if project.Status == constant.SAVEED || project.Status == constant.RESAVEED {
		delResult, err := db.DeleteProject(project.ID, rev)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "删除项目错误："+err.Error(), "", writer)
			return
		} else {
			log.Println(delResult)
			utils.ResponseJson(200, "项目删除完成", "", writer)
			return
		}
	} else {
		utils.ResponseJson(400, "项目已经部署，不能删除！", "", writer)
		return
	}
}

// fixProject(writer, request)
func getTopology(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得topology参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	project := module.Define{}
	_, err = db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	topMap := make(map[string]interface{})

	orderlist := make([]string, 0)
	for _, order := range project.Orders {
		orderlist = append(orderlist, order.OrderId)
	}
	topMap["orders"] = orderlist

	// 获得peers
	peerlist := make([]string, 0)

	for _, org := range project.Orgs {
		for _, peer := range org.Peers {
			peerlist = append(peerlist, fmt.Sprint(peer.PeerId, ".", org.OrgId))
		}
	}
	topMap["peers"] = peerlist
	returnBtyes, _ := json.Marshal(topMap)
	utils.ResponseJson(200, "获得拓扑图信息", string(returnBtyes), writer)
	return
}

// fixProject(writer, request)
func fixProject(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得project参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	if obj["id"] == nil {
		utils.ResponseJson(500, "请输入项目id", "", writer)
		return
	}

	project := module.Define{}
	_, err = db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	// 获得所有的ip
	iplist := make([]string, 0)

	for _, order := range project.Orders {
		iplist = append(iplist, order.OrderIp)
	}
	for _, org := range project.Orgs {
		for _, peer := range org.Peers {
			iplist = append(iplist, peer.PeerIp)
		}
	}
	iplist = utils.RemoveDuplicatesAndEmpty(iplist) // 删除重复
	// 获得project ip string
	perChan := make(chan string, len(iplist))
	// goroutine send file to server
	for _, ip := range iplist {
		go utils.GetProjectPerformance(ip, perChan)
	}
	// get post return
	returnMap := make(map[string]interface{})
	ipMap := make(map[string]interface{})
ReturnProjectInfo:
	for {
		iplen := len(iplist)
		select {
		case msg, ok := <-perChan:
			if !ok {
				return
			}
			msgMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(msg), &msgMap)
			if err != nil {
				utils.ResponseJson(400, "获得服务器信息错误："+err.Error(), "", writer)
				return
			}
			if int(msgMap["code"].(float64)) != 200 { // interface{} into int
				utils.ResponseJson(400, "获得服务器信息错误", "", writer)
				return
			}
			//获得ip上的docker containers
			// 获得所有的ip
			ip := msgMap["ip"]
			containerList := make([]string, 0)
			for _, order := range project.Orders {
				if ip == order.OrderIp {
					containerList = append(containerList, order.ContainerId)
				}
			}
			for _, org := range project.Orgs {
				if ip == org.CaIp {
					containerList = append(containerList, org.ContainerId)
				}
				for _, peer := range org.Peers {
					if ip == peer.PeerIp {
						containerList = append(containerList, peer.ContainerId)
						containerList = append(containerList, peer.CouchContainerId)
					}
				}
			}
			tmpData := msgMap["data"].(map[string]interface{})
			tmpData["dockers"] = containerList
			ipMap[ip.(string)] = tmpData

			// 获得iplenght
			iplen--
			if iplen <= 0 {
				break ReturnProjectInfo
			}

		}
	}
	returnMap["servers"] = ipMap
	returnMap["projectName"] = project.ProjectName
	returnBtyes, _ := json.Marshal(returnMap)
	utils.ResponseJson(200, "获得服务器信息", string(returnBtyes), writer)
	return
}

// fixProject(writer, request)
func fixIp(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得project参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}

	if obj["id"] == nil {
		utils.ResponseJson(500, "请输入项目id", "", writer)
		return
	}
	if obj["ip"] == nil {
		utils.ResponseJson(500, "请输入对应的ip", "", writer)
		return
	}
	projectId := obj["id"].(string)
	ip := obj["ip"].(string)
	project := module.Define{}
	_, err = db.GetReadSetupInfo(projectId, &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	// 获得所有的ip
	containerList := make([]string, 0)

	for _, order := range project.Orders {
		if ip == order.OrderIp {
			containerList = append(containerList, order.ContainerId)
		}
	}
	for _, org := range project.Orgs {
		if ip == org.CaIp {
			containerList = append(containerList, org.ContainerId)
		}
		for _, peer := range org.Peers {
			if ip == peer.PeerIp {
				containerList = append(containerList, peer.ContainerId)
				containerList = append(containerList, peer.CouchContainerId)
			}
		}
	}

	// 获得project ip string
	log.Println("======================containerList:")
	log.Println(containerList)
	perChan := make(chan string, 1)
	go utils.GetDockerPerformance(ip, containerList, perChan)
	msg := <-perChan //获得chan 传值
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		utils.ResponseJson(400, "获得docker信息错误："+err.Error(), "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "获得docker信息错误", "", writer)
		return
	}
	returnBtyes, _ := json.Marshal(msgMap["data"])
	utils.ResponseJson(200, "获得docker信息", string(returnBtyes), writer)
	return
}

// docker exec
func fixDocker(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得docker操作参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}
	if obj["container_name"] == nil {
		utils.ResponseJson(500, "请输入对应的container name", "", writer)
		return
	}
	if obj["ip"] == nil {
		utils.ResponseJson(500, "请输入对应的ip", "", writer)
		return
	}
	if obj["type"] == nil {
		utils.ResponseJson(500, "请输入对应的操作类型", "", writer)
		return
	}
	// 获得project ip string
	perChan := make(chan string, 1)
	go utils.DockerExec(obj["ip"].(string), obj["container_name"].(string), obj["type"].(string), perChan)
	msg := <-perChan //获得chan 传值
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		utils.ResponseJson(400, "操作docker信息错误："+err.Error(), "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "操作docker信息错误", "", writer)
		return
	}

	utils.ResponseJson(200, "操作docker完成", "", writer)
	return
}

// save log
func SaveLog(user string, body string, url string) {
	log := module.Log{}
	log.ID = utils.GetUuid()
	log.Time = time.Now().Unix()
	log.Interaction = url
	log.Data = body
	log.UserName = user
	db.SaveLogInfo(log, log.ID, "")
}

func setExplorer(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得define传参1："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}
	project := module.Define{}
	rev, err := db.GetReadSetupInfo(obj["pid"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}
	project.ExploderUser = obj["username"].(string)
	project.ExploderPassword = utils.Md5(obj["password"].(string))
	err = db.SaveSetupInfo(project, project.ID, rev)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "保存浏览器的用户名和密码错误："+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "保存浏览器的用户名和密码成功", "", writer)
	return
}

func addOrg(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得addorg传参："+err.Error(), "", writer)
		return
	}

	addOrgObj := module.AddOrgParam{}
	err = json.Unmarshal(body, &addOrgObj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析传参："+err.Error(), "", writer)
		return
	}
	project := module.Define{}
	rev, err := db.GetReadSetupInfo(addOrgObj.PID, &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	// 补全addorg
	OrgObjParam := addOrgObj.AddOrgs[0]
	OrgObj := module.Org{}
	OrgObj.ContainerId = fmt.Sprint("ca.", OrgObjParam.OrgId, ".", project.Domain)
	OrgObj.CaPort = 7054 //默认7054
	OrgObj.CaId = "ca"

	OrgObj.CaUser = "admin"
	OrgObj.CaPwd = "adminpw"
	OrgObj.OrgName = OrgObjParam.OrgId
	OrgObj.OrgId = OrgObjParam.OrgId
	OrgObj.PeerNumber = len(OrgObjParam.Peers)
	// 补全端口信息
	for i, p := range OrgObjParam.Peers {
		newpeer := module.Peer{}
		newpeer.PeerIp = p.PeerIp
		newpeer.PeerId = fmt.Sprint("peer", strconv.Itoa(i)) //增加peerid
		newpeer.PostPort = 7051 + i*1000
		newpeer.EventPort = 7053 + i*1000
		newpeer.JoinCouch = true
		newpeer.CouchPort = 4984 + i*1000
		newpeer.CouchUsername = constant.COUCHUSERNAME
		newpeer.CouchPassword = constant.COUCHPASSWORD
		newpeer.CouchId = fmt.Sprint("couch", "_", OrgObj.OrgId, "_", newpeer.PeerId)
		newpeer.CouchContainerId = fmt.Sprint("couch.", fmt.Sprint("peer", i), ".", OrgObj.OrgId, ".", project.Domain)
		newpeer.ContainerId = fmt.Sprint(newpeer.PeerId, ".", OrgObj.OrgId, ".", project.Domain)
		OrgObj.Peers = append(OrgObj.Peers, newpeer)
	}
	// 根据peer顺序修改orgobj的ip
	OrgObj.CaIp = OrgObj.Peers[0].PeerIp
	OrgObj.AnchorIp = OrgObj.Peers[0].PeerIp
	OrgObj.AnchorPort = OrgObj.Peers[0].PostPort
	fmt.Println(OrgObj)
	project.Orgs = append(project.Orgs, OrgObj) //将现有项目的组织，新增填写的组织
	// channel 里包含的组织，增加填写的组织
	for i, channel := range project.AddChannels {
		if channel.ChannelId == addOrgObj.ChannelId {
			project.AddChannels[i].IncludeOrgs = append(project.AddChannels[i].IncludeOrgs, OrgObj.OrgId)
		}
	}

	// err = db.SaveSetupInfo(project, project.ID, rev)
	// if err != nil {
	// 	utils.ResponseJson(400, "保存新的组织到项目错误："+err.Error(), "", writer)
	// 	return
	// }

	// 根据传入的org id 和 peer信息生成 yaml文件
	projectPath := filepath.Join(constant.ROOTPATH, project.ID)
	// 生成组织crypto yaml文件
	addOrgOutPath := filepath.Join(projectPath, fmt.Sprint(OrgObj.OrgId, "_", "crypto.yaml"))
	configtx := filepath.Join(constant.YAMLPATH, "addorg_crypto.yaml")
	err = utils.CryptoYaml(project, OrgObj, configtx, addOrgOutPath)
	if err != nil {
		utils.ResponseJson(400, "根据模板生成crypto yaml文件出错："+err.Error(), "", writer)
		return
	}
	// 根据yaml文件生成证书文件
	// export FABRIC_CFG_PATH
	// exportShellString := fmt.Sprint("export FABRIC_CFG_PATH=$PWD")
	cdShellString := fmt.Sprint("cd ", projectPath)
	shellString := fmt.Sprint(cdShellString, "; ../cryptogen generate --config=./", fmt.Sprint(OrgObj.OrgId, "_", "crypto.yaml"))
	fmt.Println(shellString)
	err, outstring, outerr := utils.Shellout(shellString)
	fmt.Println(outstring)
	fmt.Println(outerr)
	if err != nil {
		utils.ResponseJson(400, "执行 cryptogen命令："+err.Error(), "", writer)
		return
	}

	// //////////////////////========根据configtx 生成yaml================================================
	// 根据yaml生成json
	addConfigTx := filepath.Join(constant.YAMLPATH, "configtx_add.yaml")
	err = utils.CheckAndCreatePath(fmt.Sprint(projectPath, "/", OrgObj.OrgId))
	if err != nil {
		utils.ResponseJson(400, "生成文件的路径错误："+err.Error(), "", writer)
		return
	}
	addConfigOutpath := filepath.Join(fmt.Sprint(projectPath, "/", OrgObj.OrgId), "configtx.yaml")
	err = utils.AddOrgCrytojson(project, OrgObj, addConfigTx, addConfigOutpath)
	if err != nil {
		utils.ResponseJson(400, "根据模板生成configtx yaml文件出错："+err.Error(), "", writer)
		return
	}

	cdShellString = fmt.Sprint("cd ", filepath.Join(projectPath, OrgObj.OrgId))
	shellString = fmt.Sprint(cdShellString, ";export FABRIC_CFG_PATH=$PWD;", fmt.Sprintf("../../configtxgen -printOrg %s > ../crypto-config/%s.json", OrgObj.OrgId, OrgObj.OrgId))
	fmt.Println(shellString)
	err, outstring, outerr = utils.Shellout(shellString)
	fmt.Println(outstring)
	fmt.Println(outerr)
	if err != nil {
		utils.ResponseJson(400, "执行cryptogen命令，生成json文件错误："+err.Error(), "", writer)
		return
	}

	// 生成yaml文件到制定文件夹
	//////////////////////===================生成yaml文件开始=======================================
	// 生成order yaml 文件
	caoutpath := filepath.Join(projectPath, "ca_"+OrgObj.OrgId+".yaml")
	err = utils.CheckFileAndRemove(caoutpath)
	if err != nil {
		utils.ResponseJson(400, "生成 ca yaml 文件，clear 文件错误："+err.Error(), "", writer)
		return
	}
	err = utils.CaYaml(project, OrgObj, filepath.Join(constant.YAMLPATH, "ca_demo.yaml"), caoutpath)
	if err != nil {
		utils.ResponseJson(400, "生成 ca yaml 文件："+err.Error(), "", writer)
		return
	}
	// 		生成peer
	for _, peer := range OrgObj.Peers {
		peerOutpath := filepath.Join(projectPath, "peer_"+peer.PeerId+"."+OrgObj.OrgId+".yaml")
		err := utils.CheckFileAndRemove(peerOutpath)
		if err != nil {
			utils.ResponseJson(400, "生成peer 文件路径："+err.Error(), "", writer)
			return
		}
		err = utils.PeerYaml(project, OrgObj, peer, filepath.Join(constant.YAMLPATH, "peer_demo.yaml"), peerOutpath)
		if err != nil {
			utils.ResponseJson(400, "生成 peer yaml 文件："+err.Error(), "", writer)
			return
		}
		// 生成couch db
		if peer.JoinCouch {
			couchOutpath := filepath.Join(projectPath, "couch_"+peer.PeerId+"."+OrgObj.OrgId+".yaml")
			err := utils.CheckFileAndRemove(couchOutpath)
			if err != nil {
				utils.ResponseJson(400, "生成couch peer 文件路径："+err.Error(), "", writer)
				return
			}
			err = utils.CouchYaml(peer, filepath.Join(constant.YAMLPATH, "couch_demo.yaml"), couchOutpath)
			if err != nil {
				utils.ResponseJson(400, "生成 couch peer yaml 文件："+err.Error(), "", writer)
				return
			}
		}
	}
	// 生成cli yaml 文件
	clioutpath := filepath.Join(projectPath, "cli_"+OrgObj.OrgId+".yaml")
	err = utils.CheckFileAndRemove(clioutpath)
	if err != nil {
		utils.ResponseJson(400, "生成 cli yaml 文件，clear 文件错误："+err.Error(), "", writer)
		return
	}
	err = utils.CliYaml(project, OrgObj, filepath.Join(constant.YAMLPATH, "cli_demo.yaml"), clioutpath)
	if err != nil {
		utils.ResponseJson(400, "生成 cli yaml 文件："+err.Error(), "", writer)
		return
	}
	//////////////////////===================生成yaml文件结束=======================================

	// 发送证书文件到指定服务器

	// 压缩文件tar包
	relativePath := filepath.Join("./", project.ID)
	gizPath := filepath.Join(constant.ROOTPATH, project.ID+".tar")
	if _, err := os.Stat(gizPath); os.IsNotExist(err) {
		os.Remove(gizPath)
	}
	gizShell := fmt.Sprint("cd ", constant.ROOTPATH, ";tar -cvf ", gizPath, " ", relativePath)
	// log.Println(gizShell)
	err, _, outerr = utils.Shellout(gizShell)
	if err != nil {
		err, _, outerr = utils.Shellout(gizShell)
		utils.ResponseJson(400, "压缩过程："+outerr, "", writer)
		return
	}
	// 获得所有的ip
	iplist := make([]string, 0)
	iplist = append(iplist, OrgObj.CaIp) //ca ip
	for _, peer := range OrgObj.Peers {
		iplist = append(iplist, peer.PeerIp) // peer ip
	}
	iplist = utils.RemoveDuplicatesAndEmpty(iplist) // 删除重复
	// 发送文件
	uploadChan := make(chan string, len(iplist))
	// goroutine send file to server
	for _, ip := range iplist {
		go utils.SendFile(gizPath, ip, constant.SERVERPORT, project.ID, uploadChan)
	}
	// get post return
	for i := range iplist {
		fmt.Println(i)
		msg := <-uploadChan //获得chan 传值
		msgMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &msgMap)
		if err != nil {
			utils.ResponseJson(400, "证书和文件发送失败，请重新发送："+err.Error(), "", writer)
			return
		}
		if int(msgMap["code"].(float64)) != 200 { // interface{} into int
			utils.ResponseJson(400, "证书和文件发送失败，请重新发送", "", writer)
			return
		}
	}
	///////////////////////////////====================发送文件结束=================================

	///////////////////////////////====================启动docker开始==============================
	for _, ip := range iplist {
		composefile := ""

		if strings.TrimSpace(OrgObj.CaIp) == strings.TrimSpace(ip) {
			composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("ca_%s.yaml", OrgObj.OrgId))
			composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("cli_%s.yaml", OrgObj.OrgId))
		}
		for _, peer := range OrgObj.Peers {
			if strings.TrimSpace(peer.PeerIp) == strings.TrimSpace(ip) {
				composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("peer_%s.%s.yaml", peer.PeerId, OrgObj.OrgId))
				if peer.JoinCouch {
					composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("couch_%s.%s.yaml", peer.PeerId, OrgObj.OrgId))
				}
			}
		}

		lunchChan := make(chan string)
		// 	docker compose 命令
		downcompose := fmt.Sprintf("docker-compose %s  down;", composefile)
		restartsh := fmt.Sprint("sudo chmod 777 ./restart.sh;", "./restart.sh;")
		upcompose := fmt.Sprintf("docker-compose %s up -d ;", composefile)
		fmt.Println(fmt.Sprint(downcompose, restartsh, upcompose))
		//  rmps, rmimages, rmkeystore,
		go utils.LunchDockerEnv(ip, constant.SERVERPORT, fmt.Sprint(downcompose, restartsh, upcompose), project.ID, lunchChan)
		// select 处理channel
		select {
		case msg, ok := <-lunchChan:
			if ok {
				fmt.Println(msg)
				msgMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(msg), &msgMap)
				if err != nil {
					utils.ResponseJson(400, "启动docker 环境失败，请重试:"+err.Error(), "", writer)
					return
				}
				if int(msgMap["code"].(float64)) != 200 { // interface{} into int
					utils.ResponseJson(400, "启动docker 环境失败，请重试", "", writer)
					return
				}
			}
		}
	}
	// 启动docker针对的containers

	// 进入cli container 配置 str=$"/n";sstr=$(echo -e $str);echo "$sstr"
	clishell := fmt.Sprintf("docker exec -it %s bash -c  '", fmt.Sprint(OrgObj.OrgId, "cli"))
	clishell = fmt.Sprintln(clishell, `apt update && apt install -y jq;configtxlator start & `) //fmt.Sprintf(`sudo docker exec -it %s bash -c  '`, fmt.Sprint(OrgObj.OrgId, "cli")),
	clishell = fmt.Sprintln(clishell, `CONFIGTXLATOR_URL=http://127.0.0.1:7059;`)
	orderCaPath := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/%s/orderers/%s.%s/msp/tlscacerts/tlsca.%s-cert.pem", project.Domain, project.Orders[0].OrderId, project.Domain, project.Domain)
	// channelName := addOrgObj.ChannelId
	// clishell = fmt.Sprint(clishell, fmt.Sprintf("export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/%s/orderers/%s.%s/msp/tlscacerts/tlsca.%s-cert.pem; export CHANNEL_NAME=%s; ", project.Domain, project.Orders[0].OrderId, project.Domain, project.Domain, addOrgObj.ChannelId))
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`export CORE_PEER_LOCALMSPID="%s";`, project.Orgs[0].OrgId))
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s.%s/peers/peer0.%s.%s/tls/ca.crt;`, project.Orgs[0].OrgId, project.Domain, project.Orgs[0].OrgId, project.Domain))
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s.%s/users/Admin@%s.%s/msp;`, project.Orgs[0].OrgId, project.Domain, project.Orgs[0].OrgId, project.Domain))
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`export CORE_PEER_ADDRESS=peer0.%s.%s:%s;`, project.Orgs[0].OrgId, project.Domain, strconv.Itoa(project.Orgs[0].Peers[0].PostPort)))
	orderDomain := fmt.Sprint(project.Orders[0].OrderId, ".", project.Domain, ":", project.Orders[0].OrderPort)
	// orderDomain := fmt.Sprint(project.Orders[0].OrderIp, ":", project.Orders[0].OrderPort)
	clishell = fmt.Sprint(clishell, fmt.Sprintf("peer channel fetch config config_block.pb -o %s -c %s --tls --cafile %s;", orderDomain, addOrgObj.ChannelId, orderCaPath))
	clishell = fmt.Sprint(clishell, "curl -X POST --data-binary @config_block.pb http://127.0.0.1:7059/protolator/decode/common.Block > config_block.json;")
	clishell = fmt.Sprint(clishell, "jq .data.data[0].payload.data.config config_block.json > config.json;")
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`jq -s ".[0]*{\"channel_group\":{\"groups\":{\"Application\":{\"groups\":{\"%s\":.[1]}}}}}" config.json ./crypto/%s.json >&updated_config.json;`, OrgObj.OrgId, OrgObj.OrgId))
	clishell = fmt.Sprint(clishell, `curl -X POST --data-binary @config.json http://127.0.0.1:7059/protolator/encode/common.Config > config.pb;`)
	clishell = fmt.Sprint(clishell, `curl -X POST --data-binary @updated_config.json http://127.0.0.1:7059/protolator/encode/common.Config > updated_config.pb;`)
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`curl -X POST -F original=@config.pb -F updated=@updated_config.pb http://127.0.0.1:7059/configtxlator/compute/update-from-configs -F channel=%s > config_update.pb;`, addOrgObj.ChannelId))
	clishell = fmt.Sprint(clishell, `curl -X POST --data-binary @config_update.pb http://127.0.0.1:7059/protolator/decode/common.ConfigUpdate > config_update.json;`)
	clishell = fmt.Sprint(clishell, fmt.Sprintf(`echo "{\"payload\":{\"header\":{\"channel_header\":{\"channel_id\":\"%s\" ,\"type\":2}},\"data\":{\"config_update\":\"$(cat config_update.json)\"}}}" > config_update_as_envelope.json;`, addOrgObj.ChannelId))
	clishell = fmt.Sprint(clishell, `curl -X POST --data-binary @config_update_as_envelope.json http://127.0.0.1:7059/protolator/encode/common.Envelope > config_update_as_envelope.pb;`)
	clishell = fmt.Sprint(clishell, `peer channel signconfigtx -f config_update_as_envelope.pb;`)
	clishell = fmt.Sprint(clishell, `'`)
	fmt.Println(clishell)
	return
	message := utils.ReqShellServer(OrgObj.CaIp, clishell) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "新增组织过程中错误："+err.Error(), "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "新增组织过程发生错误", "", writer)
		return
	}
	// install chaincode，
	/////////////////////////////==========================新增组织安装chaincode========================
	jqShell := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
		echo
		exit 1
	fi`)
	// login shell
	loginShell := ""
	defaultToken := ""
	for _, org := range project.Orgs {
		for _, channel := range project.AddChannels {
			if channel.ChannelId == addOrgObj.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					if defaultToken == "" {
						defaultToken = org.OrgId
					}
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(curl -s -X POST http://localhost:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=%s&password=password&orgName=%s')`, org.OrgId, org.OrgId, org.OrgId))
					// &channelName=%s  channel.ChannelId
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(echo $%s_TOKEN | jq ".token" | sed "s/\"//g")`, org.OrgId, org.OrgId))
				}
			}
		}
	}

	joinChannel := ""
	peerStr := ""
	for _, peer := range OrgObj.Peers {
		if peerStr == "" {
			peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
		} else {
			peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
		}
	}
	joinChannel = fmt.Sprintln(joinChannel, fmt.Sprintf(`curl -s -X POST \
		http://localhost:4000/channels/peers \
		-H "authorization: Bearer $%s_TOKEN" \
		-H "content-type: application/json" \
		-d '{"peers": [%s],"channelName":"%s"}'`, OrgObj.OrgId, peerStr, addOrgObj.ChannelId))
	installChaincode := ""
	// 根据chaincode
	installCC := module.CC{}
	for _, channel := range project.AddChannels {
		if channel.ChannelId == addOrgObj.ChannelId {
			for _, cc := range channel.ChainCodes {
				if cc.Using == true {
					installCC = cc
				}
			}
		}
	}
	installChaincode = fmt.Sprintln(installChaincode, fmt.Sprintf(`curl -s -X POST \
			http://localhost:4000/chaincodes \
			-H "authorization: Bearer $%s_TOKEN" \
			-H "content-type: application/json" \
			-d '{
			  "peers": [%s],
			  "channelName":"%s",
			  "chaincodeName":"%s",
			  "chaincodePath":"%s",
			  "chaincodeVersion":"%s"
		  }'`, OrgObj.OrgId, peerStr, addOrgObj.ChannelId, installCC.CCName, installCC.CCName, installCC.CCVersion))
	//如果升级智能合约，修改allChannel
	allChannel := fmt.Sprintln(jqShell, loginShell, joinChannel, installChaincode)

	message = utils.ReqShellServer(project.Orders[0].OrderIp, allChannel) //请求执行shell命令
	msgMap = make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "新增组织和安装链码错误："+err.Error(), "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "新增组织和安装链码错误", "", writer)
		return
	}
	//==============================================新增组织完成===============================
	//////////////////////////////========================配置sdk开始=======================================
	// config 配置信息 , 有些默认信息以后修改 TODO
	config := module.Config{}
	config.CaUser = "admin"
	config.CaSecret = "adminpw"
	config.CcSrcPath = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/", "cc")
	config.ProjectName = project.ProjectName
	// config.ChannelConfigPath = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/", defindeInfo.ChannelName, ".tx")
	config.EventWaitTime = "100000"
	config.ExpireTime = "360000"
	config.Host = "localhost"
	config.Port = "4000"
	config.CurOrgId = ""
	config.KeyValueStore = "/var/fabric-client-kvs"
	config.Consensus = project.Consensus

	// config channels
	for _, channel := range project.AddChannels {
		channelMap := make(map[string]interface{})
		channelMap["channelId"] = channel.ChannelId
		channelMap["includes"] = channel.IncludeOrgs
		channelMap["channelConfigPath"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/", channel.ChannelId, ".tx")
		config.Channels = append(config.Channels, channelMap)
	}

	networkconfig := make(map[string]interface{})
	if project.Consensus == "kafka" {
		for _, order := range project.Orders {
			tmpOrder := make(map[string]string)
			tmpOrder["url"] = fmt.Sprint("grpcs://", order.OrderIp, ":", strconv.Itoa(order.OrderPort))
			tmpOrder["server-hostname"] = fmt.Sprint(order.OrderId, ".", project.Domain)
			tmpOrder["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/crypto-config/ordererOrganizations/", project.Domain, "/orderers/", order.OrderId, ".", project.Domain, "/tls/ca.crt")
			networkconfig[order.OrderId] = tmpOrder
			// add orderer to orderer list
			config.Orderers = append(config.Orderers, order.OrderId)
		}
	} else {
		if len(project.Orders) == 1 {
			order := project.Orders[0]
			tmpOrder := make(map[string]string)
			tmpOrder["url"] = fmt.Sprint("grpcs://", order.OrderIp, ":", strconv.Itoa(order.OrderPort))
			tmpOrder["server-hostname"] = fmt.Sprint(order.OrderId, ".", project.Domain)
			tmpOrder["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/crypto-config/ordererOrganizations/", project.Domain, "/orderers/", order.OrderId, ".", project.Domain, "/tls/ca.crt")
			networkconfig[order.OrderId] = tmpOrder
			// add orderer to orderer list
			config.Orderers = append(config.Orderers, order.OrderId)
		} else {
			utils.ResponseJson(400, "order未配置", "", writer)
			return
		}
	}

	for _, org := range project.Orgs {
		config.Orgs = append(config.Orgs, org.OrgId) // 记录组织
		tmpOrg := make(map[string]interface{})
		tmpOrg["aliasName"] = "org1"
		tmpOrg["name"] = org.OrgId
		tmpOrg["mspid"] = org.OrgId
		tmpOrg["ca"] = fmt.Sprint("https://", org.CaIp, ":", strconv.Itoa(org.CaPort))
		tmpAdmin := make(map[string]string)
		// /var/certification/f18fafe1e2b4494696e1dac580ab6c53/crypto-config/peerOrganizations/nxia.jiake.com/users/Admin@nxia.jiake.com/msp/keystore
		tmpAdmin["key"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", project.Domain, "/users/Admin@", org.OrgId, ".", project.Domain, "/msp/keystore")
		tmpAdmin["cert"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", project.Domain, "/users/Admin@", org.OrgId, ".", project.Domain, "/msp/signcerts")
		tmpOrg["admin"] = tmpAdmin

		tmpPeerObj := make(map[string]interface{})
		for _, peer := range org.Peers {
			// 		"requests": "grpcs://localhost:7051",
			// 		"events": "grpcs://localhost:7053",
			// 		"server-hostname": "peer0.nxia.jiake.com",
			// 		"tls_cacerts": peerOrganizations/nxia.jiake.com/peers/peer0.nxia.jiake.com/tls/ca.crt
			tmpPeer := make(map[string]string)
			tmpPeer["server-hostname"] = fmt.Sprint(peer.PeerId, ".", org.OrgId, ".", project.Domain)
			tmpPeer["requests"] = fmt.Sprint("grpcs://", peer.PeerIp, ":", strconv.Itoa(peer.PostPort))
			tmpPeer["events"] = fmt.Sprint("grpcs://", peer.PeerIp, ":", strconv.Itoa(peer.EventPort))
			tmpPeer["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", project.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", project.Domain, "/peers/", peer.PeerId, ".", org.OrgId, ".", project.Domain, "/tls/ca.crt")
			tmpPeerObj[peer.PeerId] = tmpPeer
		}
		tmpOrg["peers"] = tmpPeerObj //peers 初始化
		networkconfig[org.OrgId] = tmpOrg
	}
	config.NetworkConfig = networkconfig

	//根据ip 分发
	sdkChan := make(chan string)
	for _, ip := range iplist {
		if ip == project.Orders[0].OrderIp { //如果是order ip，那么为管理端，sdk负责搭建链和更新链码
			config.Manager = true
		} else {
			config.Manager = false
		}
		// 根据共识判断orderid
		if project.Consensus == "solo" {
			config.OrderId = project.Orders[0].OrderId // 根据如果是solo order只有一个
		} else {
			// 如果是kafka，orderid 随机分配
			randNumber := rand.Intn(len(project.Orders))
			if randNumber >= len(project.Orders) {
				randNumber = randNumber - 1
			}
			config.OrderId = project.Orders[randNumber].OrderId
		}
		// 计算当前组织的id
		// config.CurOrgId
		for _, org := range project.Orgs {
			if ip == org.CaIp {
				config.CurOrgId = org.OrgId
				break
			}
		}
		// 分发sdk
		configBytes, err := json.Marshal(config)
		if err != nil {
			utils.ResponseJson(400, "config json："+err.Error(), "", writer)
			return
		}
		// 请求sdk 目录
		fmt.Sprintln("===============================ip:", ip)
		go utils.SdkSendConfig(ip, constant.SERVERPORT, project.ID, configBytes, sdkChan)
	}
	// 获得chan
	for i, _ := range iplist {
		fmt.Println(i)
		msg := <-sdkChan //获得lunch chan 传值
		fmt.Println(msg)
		msgMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &msgMap)
		if err != nil {
			utils.ResponseJson(400, "配置sdk错误", "", writer)
			return
		}
		if int(msgMap["code"].(float64)) != 200 { // interface{} into int
			utils.ResponseJson(400, "配置sdk错误", "", writer)
			return
		}
	}
	//////////////////////////////=========================配置sdk结束==================================

	//////////////////////////////==========================开启服务器端口================================
	installChan := make(chan string, 1)
	for _, ip := range iplist {
		firewallCommand := ""
		firewallCommand = fmt.Sprint(firewallCommand, fmt.Sprintf(constant.ADDPORT, "4000"))

		if OrgObj.CaIp == ip {
			firewallCommand = fmt.Sprint(firewallCommand, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(OrgObj.CaPort)))

		}
		for _, peer := range OrgObj.Peers {
			if peer.PeerIp == ip {
				firewallCommand = fmt.Sprint(firewallCommand, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(peer.EventPort)))
				firewallCommand = fmt.Sprint(firewallCommand, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(peer.PostPort)))
				firewallCommand = fmt.Sprint(firewallCommand, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(peer.CouchPort)))
			}
		}
		// ===================重启端口
		// 	重启firewall
		firewallCommand = fmt.Sprint(firewallCommand, constant.RESTATFIREWALL)
		go utils.InstallIPEnv(ip, constant.SERVERPORT, firewallCommand, installChan)

		select {
		case msg, ok := <-installChan:
			if ok {
				fmt.Println("开启服务器端口，重启firewall：")
				fmt.Println(msg)
				msgMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(msg), &msgMap)
				if err != nil {
					utils.ResponseJson(400, "开启服务器端口未安装完整，请重新安装"+err.Error(), "", writer)
					return
				}
				if int(msgMap["code"].(float64)) != 200 { // interface{} into int
					utils.ResponseJson(400, "开启服务器端口未安装完整，请重新安装", "", writer)
					return
				}
			}
		}
	}
	//////////////////////////////==========================开启服务器端口================================
	//////////////////////////////===========================保持项目信息================================

	err = db.SaveSetupInfo(project, project.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存新的组织到项目错误："+err.Error(), "", writer)
		return
	}
	//////////////////////////////===========================保持项目信息================================

	utils.ResponseJson(200, "新增组织完成，请确认", "", writer)
	return
}

// 获得define 传参处理
func defineHandle(writer http.ResponseWriter, request *http.Request, status string) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得define传参1："+err.Error(), "", writer)
		return
	}
	// log.Println(body) // 打印body
	curUser := request.Header.Get("username")
	SaveLog(curUser, string(body), request.URL.Path) //保存
	// 根据当前用户的角色判断status
	// currentUser, _ := db.GetUser(curUser)
	// if currentUser != nil {
	// 	if currentUser.(module.User).Role == "manager" {
	// 		status = constant.RESAVEED
	// 	} else {
	// 		status = constant.SAVEED
	// 	}
	// }

	defineParam := module.DefineParam{}
	err = json.Unmarshal(body, &defineParam)

	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析define传参2："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	rev := ""
	if defineParam.ID != "" {
		rev, err = db.GetReadSetupInfo(defineParam.ID, &defindeInfo)
		if err != nil {
			utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
			return
		}
	}
	//根据post 的数据赋值，保持
	defindeInfo.ProjectName = defineParam.ProjectName
	defindeInfo.CreateTime = time.Now().Unix()
	defindeInfo.Domain = defineParam.Domain
	defindeInfo.NetWork = defineParam.NetWork
	defindeInfo.Consensus = defineParam.Consensus
	defindeInfo.OrderId = defineParam.OrderId
	defindeInfo.OrderName = defineParam.OrderName
	defindeInfo.KafkaIp = defineParam.KafkaIp
	if status == constant.RESAVEED {
		defindeInfo.Status = constant.RESAVEED
		if defindeInfo.Manager == "" {
			defindeInfo.Manager = request.Header.Get("username") //当前人
		}
	} else {
		defindeInfo.Manager = request.Header.Get("username") //当前人
		defindeInfo.Status = constant.SAVEED
	}
	defindeInfo.ProjectPassword = defineParam.ProjectPassword
	// 将param 转换为 define
	defindeInfo.Orders = defindeInfo.Orders[:0] //清空orders
	for _, orderTpm := range defineParam.Orders {
		// add 默认端口 7050
		Order := module.Order{OrderParam: orderTpm, ContainerId: ""}
		// defindeInfo.Orders
		defindeInfo.Orders = append(defindeInfo.Orders, Order)
	}
	// 将param 转换为define
	defindeInfo.Orgs = defindeInfo.Orgs[:0] //清空orgs
	for _, orgTmp := range defineParam.Orgs {
		Org := module.Org{}
		Org.ContainerId = ""
		Org.CaUser = "admin"
		Org.CaPwd = "adminpw"
		Org.CaId = orgTmp.CaId
		Org.CaIp = orgTmp.CaIp
		Org.CaPort = orgTmp.CaPort
		Org.OrgId = orgTmp.OrgId
		Org.OrgName = orgTmp.OrgName
		Org.PeerNumber = orgTmp.PeerNumber

		// Org.Peers = Org.Peers[:0]
		for _, peerTmp := range orgTmp.Peers {
			peer := module.Peer{PeerParam: peerTmp, ContainerId: "", CouchContainerId: "", CouchId: "couch_" + peerTmp.PeerId + "_" + orgTmp.OrgId}
			Org.Peers = append(Org.Peers, peer)
		}
		// 默认第一个peer 为 anchor peer
		if len(Org.Peers) > 0 {
			Org.AnchorIp = Org.Peers[0].PeerIp
			Org.AnchorPort = Org.Peers[0].PostPort
		}
		defindeInfo.Orgs = append(defindeInfo.Orgs, Org)
	}
	if defineParam.ID == "" {
		defineParam.ID = utils.GetUuid()
		defindeInfo.ID = defineParam.ID
	}
	// 修改默认值
	defindeInfo = utils.Fill(defindeInfo)
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保持setup文档："+err.Error(), "", writer)
		return
	}
	if status == constant.SAVEED {
		utils.ResponseJson(200, "项目提交完成", "", writer)
	} else {
		utils.ResponseJson(200, "项目审核完成", "", writer)
	}

	return
}

// 生成channel tx 文件
func channelgenHandle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得channel传参："+err.Error(), "", writer)
		return
	}

	//保存日志
	curUser := request.Header.Get("username")
	SaveLog(curUser, string(body), request.URL.Path) //保存日志

	channelTxParam := module.ChannelTx{}
	err = json.Unmarshal(body, &channelTxParam)

	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析channel传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	rev, err := db.GetReadSetupInfo(channelTxParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "channel根据ID get define："+err.Error(), "", writer)
		return
	}

	addChannel := module.Channel{}
	addChannel.ChannelId = channelTxParam.ChannelId
	addChannel.IncludeOrgs = channelTxParam.IncludeOrgs
	addChannel.CreateTime = time.Now().Unix()
	//判断channel 是否存在
	hasTheChannel := false
	for _, channel := range defindeInfo.AddChannels {
		if channel.ChannelId == addChannel.ChannelId {
			hasTheChannel = true
			if channel.ChainCodes != nil && len(channel.ChainCodes) > 0 {
				utils.ResponseJson(400, fmt.Sprint("通道", addChannel.ChannelId, "已经存在，请输入其它ID信息"), "", writer)
				return
			}
		}
	}
	if hasTheChannel == false {
		defindeInfo.AddChannels = append(defindeInfo.AddChannels, addChannel)
	}
	// 发送channel file 到各个子服务器上
	projectPath := fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID)
	// export FABRIC_CFG_PATH
	exportShellString := fmt.Sprint("export FABRIC_CFG_PATH=$PWD")
	cdShellString := fmt.Sprint("cd ", projectPath)

	configtxgenchannelShellString := fmt.Sprint(cdShellString, ";", exportShellString, "; ../configtxgen -profile ProjectOrgsChannel -outputCreateChannelTx ./", channelTxParam.ChannelId, ".tx", " -channelID ", channelTxParam.ChannelId)
	fmt.Println(configtxgenchannelShellString)
	err, outstring, outerr := utils.Shellout(configtxgenchannelShellString)
	fmt.Println(outstring)
	fmt.Println(outerr)
	if err != nil {
		utils.ResponseJson(400, "执行 configtxgen 创建channel命令"+err.Error(), "", writer)
		return
	}
	// 发送文件到对应的服务
	channelChan := make(chan string, 1)
	ipOrder := defindeInfo.Orders[0].OrderIp //order ip
	go utils.SendChannelFile(fmt.Sprint(channelTxParam.ChannelId, ".tx"), fmt.Sprint(projectPath, "/", channelTxParam.ChannelId, ".tx"), ipOrder, constant.SERVERPORT, channelTxParam.ID, channelChan)

	msg := <-channelChan //获得chan 传值
	// 解析go 发送的值
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		utils.ResponseJson(400, "channel tx 文件发送失败，请重试", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "channel tx 文件发送失败，请重试", "", writer)
		return
	}

	//////////////////////////////========================配置sdk开始=======================================
	// 获得所有的ip
	iplist := make([]string, 0)
	iplist = append(iplist, defindeInfo.KafkaIp) //kafkaip
	for _, org := range defindeInfo.Orgs {
		iplist = append(iplist, org.CaIp) // ca ip
		for _, peer := range org.Peers {
			iplist = append(iplist, peer.PeerIp) // peer ip
		}
	}
	for _, order := range defindeInfo.Orders {
		iplist = append(iplist, order.OrderIp) // order ip
	}
	iplist = utils.RemoveDuplicatesAndEmpty(iplist) // 删除重复

	// config 配置信息 , 有些默认信息以后修改 TODO
	config := module.Config{}
	config.CaUser = "admin"
	config.CaSecret = "adminpw"
	config.CcSrcPath = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/", "cc")
	config.ProjectName = defindeInfo.ProjectName
	// config.ChannelConfigPath = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/", defindeInfo.ChannelName, ".tx")
	config.EventWaitTime = "100000"
	config.ExpireTime = "360000"
	config.Host = "localhost"
	config.Port = "4000"
	config.CurOrgId = ""
	config.ProjectPassword = defindeInfo.ProjectPassword
	config.KeyValueStore = "/var/fabric-client-kvs"
	config.Consensus = defindeInfo.Consensus
	// config channels
	for _, channel := range defindeInfo.AddChannels {
		channelMap := make(map[string]interface{})
		channelMap["channelId"] = channel.ChannelId
		channelMap["includes"] = channel.IncludeOrgs
		channelMap["channelConfigPath"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/", channel.ChannelId, ".tx")
		config.Channels = append(config.Channels, channelMap)
	}

	networkconfig := make(map[string]interface{})
	if defindeInfo.Consensus == "kafka" {
		for _, order := range defindeInfo.Orders {
			tmpOrder := make(map[string]string)
			tmpOrder["url"] = fmt.Sprint("grpcs://", order.OrderIp, ":", strconv.Itoa(order.OrderPort))
			tmpOrder["server-hostname"] = fmt.Sprint(order.OrderId, ".", defindeInfo.Domain)
			tmpOrder["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/crypto-config/ordererOrganizations/", defindeInfo.Domain, "/orderers/", order.OrderId, ".", defindeInfo.Domain, "/tls/ca.crt")
			networkconfig[order.OrderId] = tmpOrder
			// add orderer to orderer list
			config.Orderers = append(config.Orderers, order.OrderId)
		}
	} else {
		if len(defindeInfo.Orders) == 1 {
			order := defindeInfo.Orders[0]
			tmpOrder := make(map[string]string)
			tmpOrder["url"] = fmt.Sprint("grpcs://", order.OrderIp, ":", strconv.Itoa(order.OrderPort))
			tmpOrder["server-hostname"] = fmt.Sprint(order.OrderId, ".", defindeInfo.Domain)
			tmpOrder["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/crypto-config/ordererOrganizations/", defindeInfo.Domain, "/orderers/", order.OrderId, ".", defindeInfo.Domain, "/tls/ca.crt")
			networkconfig[order.OrderId] = tmpOrder
			// add orderer to orderer list
			config.Orderers = append(config.Orderers, order.OrderId)
		} else {
			utils.ResponseJson(400, "order未配置", "", writer)
			return
		}
	}

	for _, org := range defindeInfo.Orgs {
		config.Orgs = append(config.Orgs, org.OrgId) // 记录组织
		tmpOrg := make(map[string]interface{})
		tmpOrg["aliasName"] = "org1"
		tmpOrg["name"] = org.OrgId
		tmpOrg["mspid"] = org.OrgId
		tmpOrg["ca"] = fmt.Sprint("https://", org.CaIp, ":", strconv.Itoa(org.CaPort))
		tmpAdmin := make(map[string]string)
		// /var/certification/f18fafe1e2b4494696e1dac580ab6c53/crypto-config/peerOrganizations/nxia.jiake.com/users/Admin@nxia.jiake.com/msp/keystore
		tmpAdmin["key"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", defindeInfo.Domain, "/users/Admin@", org.OrgId, ".", defindeInfo.Domain, "/msp/keystore")
		tmpAdmin["cert"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", defindeInfo.Domain, "/users/Admin@", org.OrgId, ".", defindeInfo.Domain, "/msp/signcerts")
		tmpOrg["admin"] = tmpAdmin

		tmpPeerObj := make(map[string]interface{})
		for _, peer := range org.Peers {
			// "requests": "grpcs://localhost:7051",
			// 		"events": "grpcs://localhost:7053",
			// 		"server-hostname": "peer0.nxia.jiake.com",
			// 		"tls_cacerts": peerOrganizations/nxia.jiake.com/peers/peer0.nxia.jiake.com/tls/ca.crt
			tmpPeer := make(map[string]string)
			tmpPeer["server-hostname"] = fmt.Sprint(peer.PeerId, ".", org.OrgId, ".", defindeInfo.Domain)
			tmpPeer["requests"] = fmt.Sprint("grpcs://", peer.PeerIp, ":", strconv.Itoa(peer.PostPort))
			tmpPeer["events"] = fmt.Sprint("grpcs://", peer.PeerIp, ":", strconv.Itoa(peer.EventPort))
			tmpPeer["tls_cacerts"] = fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/crypto-config/peerOrganizations/", org.OrgId, ".", defindeInfo.Domain, "/peers/", peer.PeerId, ".", org.OrgId, ".", defindeInfo.Domain, "/tls/ca.crt")
			tmpPeerObj[peer.PeerId] = tmpPeer
		}
		tmpOrg["peers"] = tmpPeerObj //peers 初始化
		networkconfig[org.OrgId] = tmpOrg
	}
	config.NetworkConfig = networkconfig

	//根据ip 分发
	sdkChan := make(chan string)
	for _, ip := range iplist {
		if ip == defindeInfo.Orders[0].OrderIp { //如果是order ip，那么为管理端，sdk负责搭建链和更新链码
			config.Manager = true
		} else {
			config.Manager = false
		}
		// 根据共识判断orderid
		if defindeInfo.Consensus == "solo" {
			config.OrderId = defindeInfo.Orders[0].OrderId // 根据如果是solo order只有一个
		} else {
			// 如果是kafka，orderid 随机分配
			randNumber := rand.Intn(len(defindeInfo.Orders))
			if randNumber >= len(defindeInfo.Orders) {
				randNumber = randNumber - 1
			}
			config.OrderId = defindeInfo.Orders[randNumber].OrderId
		}
		// 计算当前组织的id
		// config.CurOrgId
		for _, org := range defindeInfo.Orgs {
			if ip == org.CaIp {
				config.CurOrgId = org.OrgId
				break
			}
		}
		// 分发sdk
		configBytes, err := json.Marshal(config)
		if err != nil {
			utils.ResponseJson(400, "config json："+err.Error(), "", writer)
			return
		}
		// 请求sdk 目录
		fmt.Sprintln("===============================ip:", ip)
		go utils.SdkSendConfig(ip, constant.SERVERPORT, defindeInfo.ID, configBytes, sdkChan)
	}
	// 获得chan
	for i, _ := range iplist {
		fmt.Println(i)
		msg := <-sdkChan //获得lunch chan 传值
		fmt.Println(msg)
		msgMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &msgMap)
		if err != nil {
			utils.ResponseJson(400, "配置sdk错误", "", writer)
			return
		}
		if int(msgMap["code"].(float64)) != 200 { // interface{} into int
			utils.ResponseJson(400, "配置sdk错误", "", writer)
			return
		}

	}
	//////////////////////////////=========================配置sdk结束==================================
	//////////////////////////////=====================================================================
	defindeInfo.Status = constant.CHANED
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保持项目文档错误："+err.Error(), "", writer)
		return
	}
	//////////////////////////////=====================================================================
	utils.ResponseJson(200, "生成通道完成", "", writer)
	return
}

// 部署
func deployHandle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得generate传参："+err.Error(), "", writer)
		return
	}

	//保存日志
	curUser := request.Header.Get("username")
	SaveLog(curUser, string(body), request.URL.Path) //保存日志

	generateParam := module.Genarate{}
	err = json.Unmarshal(body, &generateParam)

	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析generate传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	rev, err := db.GetReadSetupInfo(generateParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "generate 根据ID get define："+err.Error(), "", writer)
		return
	}

	// 	生成某个项目的证书路径
	err = utils.ClearCertficationPath(defindeInfo.ID)
	if err != nil {
		utils.ResponseJson(400, "删除项目目录下所有证书："+err.Error(), "", writer)
		return
	}
	projectPath, err := utils.CreateCertficationPath(defindeInfo.ID)
	if err != nil {
		utils.ResponseJson(400, "generate 生成路径错误："+err.Error(), "", writer)
		return
	}

	// 证书生成
	var configtx string
	var outpath string
	var cryptgen string
	var cryptgenoutpath string
	if defindeInfo.Consensus == "solo" {
		configtx = filepath.Join(constant.YAMLPATH, "configtx_solo.yaml")
	} else {
		configtx = filepath.Join(constant.YAMLPATH, "configtx_kafka.yaml")
	}
	outpath = filepath.Join(projectPath, "configtx.yaml")
	err = utils.YamltoYaml(defindeInfo, configtx, outpath)
	if err != nil {
		utils.ResponseJson(400, "根据模板生成yaml文件出错："+err.Error(), "", writer)

		return
	}
	// crypto generate
	cryptgen = filepath.Join(constant.YAMLPATH, "cryptogen.yaml")
	cryptgenoutpath = filepath.Join(projectPath, "cryptogen.yaml")

	err = utils.YamltoYaml(defindeInfo, cryptgen, cryptgenoutpath)
	if err != nil {
		utils.ResponseJson(400, "根据模板生成crypto generate yaml文件出错："+err.Error(), "", writer)

		return
	}

	// cryptogen generate --config=./cryptogen.yaml
	// configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./genesis.block
	// configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel.tx -channelID agelichannel

	configPath := filepath.Join(projectPath, "crypto-config")
	err = utils.CheckAndCreatePath(configPath)
	if err != nil {
		utils.ResponseJson(400, "创建crypto-config 目录："+err.Error(), "", writer)

		return
	}

	// export FABRIC_CFG_PATH
	exportShellString := fmt.Sprint("export FABRIC_CFG_PATH=$PWD")
	cdShellString := fmt.Sprint("cd ", projectPath)
	// cryptoToolPath := filepath.Join(constant.ROOTPATH, "cryptogen")
	// cryptoFilePath := filepath.Join(projectPath, "cryptogen.yaml")
	// exportShellString, ";",
	shellString := fmt.Sprint(cdShellString, "; ../cryptogen generate --config=./cryptogen.yaml")
	fmt.Println(shellString)
	err, outstring, outerr := utils.Shellout(shellString)
	fmt.Println(outstring)
	fmt.Println(outerr)
	if err != nil {
		utils.ResponseJson(400, "执行 cryptogen命令："+err.Error(), "", writer)

		return
	}

	// configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./genesis.block
	// configtxgenToolPath := filepath.Join(constant.ROOTPATH, "configtxgen")
	// exportShellString, ";",
	configtxgenShellString := fmt.Sprint(cdShellString, ";", exportShellString, "; ../configtxgen -profile ProjectOrgsOrdererGenesis -outputBlock ./genesis.block")
	fmt.Println(configtxgenShellString)
	err, outstring, outerr = utils.Shellout(configtxgenShellString)
	fmt.Println(outstring)
	fmt.Println(outerr)
	if err != nil {
		utils.ResponseJson(400, "执行 configtxgen 命令："+err.Error(), "", writer)
		return
	}

	//////////////////////===================生成yaml文件开始=======================================
	// 生成order yaml 文件
	for _, ordertmp := range defindeInfo.Orders {
		outpath := filepath.Join(projectPath, "order_"+ordertmp.OrderId+".yaml")
		err := utils.CheckFileAndRemove(outpath)
		if err != nil {
			utils.ResponseJson(400, "clear 文件错误："+err.Error(), "", writer)
			return
		}
		err = utils.OrderYaml(defindeInfo, ordertmp, filepath.Join(constant.YAMLPATH, "order_demo.yaml"), outpath)
		if err != nil {
			utils.ResponseJson(400, "生成yaml order 文件："+err.Error(), "", writer)
			return
		}
	}
	if defindeInfo.Consensus == "kafka" {
		kafkapath := filepath.Join(projectPath, "kafka.yaml")
		err := utils.CheckFileAndRemove(kafkapath)
		if err != nil {
			utils.ResponseJson(400, "clear 文件错误："+err.Error(), "", writer)
			return
		}
		err = utils.KafkaYaml(defindeInfo, filepath.Join(constant.YAMLPATH, "kafka_demo.yaml"), kafkapath)
		if err != nil {
			utils.ResponseJson(400, "生成kafka yaml ："+err.Error(), "", writer)
			return
		}
	}
	// org for ca
	for _, org := range defindeInfo.Orgs {
		caoutpath := filepath.Join(projectPath, "ca_"+org.OrgId+".yaml")
		err := utils.CheckFileAndRemove(caoutpath)
		if err != nil {
			utils.ResponseJson(400, "生成 ca yaml 文件，clear 文件错误："+err.Error(), "", writer)
			return
		}
		err = utils.CaYaml(defindeInfo, org, filepath.Join(constant.YAMLPATH, "ca_demo.yaml"), caoutpath)
		if err != nil {
			utils.ResponseJson(400, "生成 ca yaml 文件："+err.Error(), "", writer)
			return
		}
		// 		生成peer
		for _, peer := range org.Peers {
			peerOutpath := filepath.Join(projectPath, "peer_"+peer.PeerId+"."+org.OrgId+".yaml")
			err := utils.CheckFileAndRemove(peerOutpath)
			if err != nil {
				utils.ResponseJson(400, "生成peer 文件路径："+err.Error(), "", writer)
				return
			}
			err = utils.PeerYaml(defindeInfo, org, peer, filepath.Join(constant.YAMLPATH, "peer_demo.yaml"), peerOutpath)
			if err != nil {
				utils.ResponseJson(400, "生成 peer yaml 文件："+err.Error(), "", writer)
				return
			}
			// 生成couch db
			if peer.JoinCouch {
				couchOutpath := filepath.Join(projectPath, "couch_"+peer.PeerId+"."+org.OrgId+".yaml")
				err := utils.CheckFileAndRemove(couchOutpath)
				if err != nil {
					utils.ResponseJson(400, "生成couch peer 文件路径："+err.Error(), "", writer)
					return
				}
				err = utils.CouchYaml(peer, filepath.Join(constant.YAMLPATH, "couch_demo.yaml"), couchOutpath)
				if err != nil {
					utils.ResponseJson(400, "生成 couch peer yaml 文件："+err.Error(), "", writer)
					return
				}
			}
		}
	}
	//////////////////////===================生成yaml文件结束=======================================

	//////////////////////====================发送文件开始 ==========================================
	relativePath := filepath.Join("./", defindeInfo.ID)

	// copy restart.sh 文件到项目文件夹下面
	cpShell := fmt.Sprint("cd ", constant.ROOTPATH, ";cp ./", constant.SHELLRESTART, " ", relativePath, "/", constant.SHELLRESTART, ";")
	err, _, outerr = utils.Shellout(cpShell)
	if err != nil {
		utils.ResponseJson(400, "拷贝restart.sh文件到项目目录中："+err.Error(), "", writer)
		return
	}

	// 压缩文件tar包
	gizPath := filepath.Join(constant.ROOTPATH, defindeInfo.ID+".tar")
	if _, err := os.Stat(gizPath); os.IsNotExist(err) {
		os.Remove(gizPath)
	}
	gizShell := fmt.Sprint("cd ", constant.ROOTPATH, ";tar -cvf ", gizPath, " ", relativePath)
	// log.Println(gizShell)
	err, _, outerr = utils.Shellout(gizShell)
	if err != nil {
		err, _, outerr = utils.Shellout(gizShell)
		utils.ResponseJson(400, "压缩过程："+outerr, "", writer)
		return
	}
	// 获得所有的ip
	iplist := make([]string, 0)
	iplist = append(iplist, defindeInfo.KafkaIp) //kafkaip
	for _, org := range defindeInfo.Orgs {
		iplist = append(iplist, org.CaIp) // ca ip
		for _, peer := range org.Peers {
			iplist = append(iplist, peer.PeerIp) // peer ip
		}
	}
	for _, order := range defindeInfo.Orders {
		iplist = append(iplist, order.OrderIp) // order ip
	}

	iplist = utils.RemoveDuplicatesAndEmpty(iplist) // 删除重复
	// 发送文件
	uploadChan := make(chan string, len(iplist))
	// goroutine send file to server
	for _, ip := range iplist {
		fmt.Println("=======================ip:", ip)
		go utils.SendFile(gizPath, ip, constant.SERVERPORT, defindeInfo.ID, uploadChan)
	}
	// get post return
	for i := range iplist {
		fmt.Println(i)
		msg := <-uploadChan //获得chan 传值
		msgMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &msgMap)
		if err != nil {
			utils.ResponseJson(400, "证书和文件发送失败，请重新发送："+err.Error(), "", writer)
			return
		}
		if int(msgMap["code"].(float64)) != 200 { // interface{} into int
			utils.ResponseJson(400, "证书和文件发送失败，请重新发送", "", writer)
			return
		}
	}
	///////////////////////////////====================发送文件结束=================================

	///////////////////////////////====================启动docker开始==============================
	for _, ip := range iplist {
		composefile := ""
		// 	order 端口安装 kafka
		if (defindeInfo.Consensus == "kafka") && (strings.TrimSpace(defindeInfo.KafkaIp) == strings.TrimSpace(ip)) {
			composefile = fmt.Sprint(composefile, " -f kafka.yaml ")
		}
		// 	order 端口安装 firewall
		for _, order := range defindeInfo.Orders {
			if strings.TrimSpace(order.OrderIp) == strings.TrimSpace(ip) {
				composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("order_%s.yaml", order.OrderId))
			}
		}
		// 	peer 端口安装 firewall
		for _, org := range defindeInfo.Orgs {
			if strings.TrimSpace(org.CaIp) == strings.TrimSpace(ip) {
				composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("ca_%s.yaml", org.OrgId))
			}
			for _, peer := range org.Peers {
				if strings.TrimSpace(peer.PeerIp) == strings.TrimSpace(ip) {
					composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("peer_%s.%s.yaml", peer.PeerId, org.OrgId))
					if peer.JoinCouch {
						composefile = fmt.Sprint(composefile, " -f ", fmt.Sprintf("couch_%s.%s.yaml", peer.PeerId, org.OrgId))
					}
				}
			}
		}

		lunchChan := make(chan string)
		// 	docker compose 命令
		downcompose := fmt.Sprintf("sudo docker-compose %s  down;", composefile)
		// rmps := `CONTAINER_IDS=$(docker ps -aq);if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" = " " ]; then echo "========== No images available for deletion ===========" else docker rm -f $CONTAINER_IDS fi;`
		// rmimages := `DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}');if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" = " " ]; then echo "========== No images available for deletion ===========" else docker rmi -f $DOCKER_IMAGE_IDS fi;`
		// rmkeystore := fmt.Sprint("rm -rf ./fabric-client-kv*;")

		restartsh := fmt.Sprint("sudo chmod 777 ./restart.sh;", "./restart.sh;")
		upcompose := fmt.Sprintf("sudo docker-compose %s up -d ;", composefile)
		fmt.Println(fmt.Sprint(downcompose, restartsh, upcompose))
		//  rmps, rmimages, rmkeystore,
		go utils.LunchDockerEnv(ip, constant.SERVERPORT, fmt.Sprint(downcompose, restartsh, upcompose), defindeInfo.ID, lunchChan)

		msg := <-lunchChan //获得lunch chan 传值
		fmt.Println(msg)
		msgMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &msgMap)
		if err != nil {
			utils.ResponseJson(400, "启动docker 环境失败，请重试:"+err.Error(), "", writer)
			return
		}
		if int(msgMap["code"].(float64)) != 200 { // interface{} into int
			utils.ResponseJson(400, "启动docker 环境失败，请重试", "", writer)
			return
		}
	}
	//////////////////////////////======================启动docker结束====================================
	//========================== 修改状态，保存=======================
	defindeInfo.Status = constant.DEPLOYED
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保持setup文档："+err.Error(), "", writer)
		return
	}
	//=============================================================
	utils.ResponseJson(200, "服务器部署完成", "", writer)
	return
}

// check 环境
func checkEnvHandle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())

		utils.ResponseJson(500, "获得check环境传参："+err.Error(), "", writer)
		return
	}
	ipParam := module.IPANDID{}
	err = json.Unmarshal(body, &ipParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析check环境参数："+err.Error(), "", writer)

		return
	}
	// check envirment
	sysChan := make(chan string, 1)

	go utils.GetSysEnv(ipParam.IP, constant.SERVERPORT, sysChan)
	sysMsg := <-sysChan //获得system chan 传值
	fmt.Println("==========get sys env msg:")
	fmt.Println(sysMsg)
	sysMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(sysMsg), &sysMap)
	if sysMap["code"] == nil {
		utils.ResponseJson(400, "环境未配置完整", "", writer)
		return
	}
	if int(sysMap["code"].(float64)) != 200 {
		data := sysMap["data"].(map[string]interface{})
		if data["type"].(int) == 1 { // 1 centos 0 ubuntu 2 其它
			// 发送文件
			checkChan := make(chan string, 5)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.CHECKDOCKER, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.CHECKDOCKERCOMPOSE, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.CHECKNODE, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.CHECKGIT, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.CHECKJQ, checkChan)
			goruntineNum := 5
		CentosChan:
			for {
				select {
				case msg := <-checkChan:
					msgMap := make(map[string]interface{})
					err = json.Unmarshal([]byte(msg), &msgMap)
					if err != nil {
						utils.ResponseJson(400, "环境未配置完整，请重新安装："+err.Error(), "", writer)

						return
					}
					if int(msgMap["code"].(float64)) != 200 { // interface{} into int
						utils.ResponseJson(400, "环境未配置完整，请重新安装", "", writer)
						return
					}
					// 处理完后减一
					goruntineNum--
					if goruntineNum <= 0 {
						break CentosChan
					}
				}
			}
		}
		// ubuntu 系统
		if data["type"].(int) == 0 { // 1 centos 0 ubuntu 2 其它
			// 发送文件
			checkChan := make(chan string, 5)

			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_CHECKDOCKER, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_CHECKDOCKERCOMPOSE, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_CHECKNODE, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_CHECKGIT, checkChan)
			go utils.CheckIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_CHECKJQ, checkChan)
			goruntineNum := 5
		UbuntuChan:
			for {
				select {
				case msg := <-checkChan:
					msgMap := make(map[string]interface{})
					err = json.Unmarshal([]byte(msg), &msgMap)
					if err != nil {
						utils.ResponseJson(400, "环境未配置完整，请重新安装："+err.Error(), "", writer)
						return
					}
					if int(msgMap["code"].(float64)) != 200 { // interface{} into int
						utils.ResponseJson(400, "环境未配置完整，请重新安装", "", writer)
						return
					}
					// 处理完后减一
					goruntineNum--
					if goruntineNum <= 0 {
						break UbuntuChan
					}
				}
			}
		}

		if data["type"].(int) == 2 { // 1 centos 0 ubuntu 2 其它
			utils.ResponseJson(400, "系统不是centos或者ubuntu，请更换系统OS", "", writer)
			return
		}
	}

	utils.ResponseJson(200, "环境已经配置完成", "", writer)
	return
}

// install 环境
func installEnvHandle(writer http.ResponseWriter, request *http.Request) {
	goroutineNum := 0
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得install环境传参："+err.Error(), "", writer)
		return
	}
	ipParam := module.IPANDID{}
	err = json.Unmarshal(body, &ipParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析install环境参数："+err.Error(), "", writer)
		return
	}

	// check envirment
	sysChan := make(chan string, 1)
	go utils.GetSysEnv(ipParam.IP, constant.SERVERPORT, sysChan)
	sysMsg := <-sysChan //获得system chan 传值

	sysMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(sysMsg), &sysMap)
	if int(sysMap["code"].(float64)) != 200 {
		data := sysMap["data"].(map[string]interface{})
		if data["type"].(int) == 1 { // 1 centos 0 ubuntu 2 其它
			// install docker docker-compose nodejs
			installChan := make(chan string, 4)

			// go routine start ===========================================================
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.INSTALLDOCKER, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.INSTALLDOCKERCOMPOSE, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.INSTALLNODEJS, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.INSTALLGIT, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.INSTALLJQ, installChan)
			// 根据IP 获得该项目中所有的端口号
			defindeInfo := module.Define{}
			if ipParam.ID != "" {
				_, err = db.GetReadSetupInfo(ipParam.ID, &defindeInfo)
				if err != nil {
					utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
					return
				}
			}
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "4000"), installChan)
			goroutineNum = goroutineNum + 1
			// 	order 端口安装 kafka
			if strings.TrimSpace(defindeInfo.KafkaIp) == strings.TrimSpace(ipParam.IP) {
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9092"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9093"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9192"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9193"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9292"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9293"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9392"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "9393"), installChan)
				goroutineNum = goroutineNum + 8
			}
			// 	order 端口安装 firewall
			for _, order := range defindeInfo.Orders {
				if strings.TrimSpace(order.OrderIp) == strings.TrimSpace(ipParam.IP) {
					go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(order.OrderPort)), installChan)
					goroutineNum = goroutineNum + 1
				}
			}
			// 	peer 端口安装 firewall
			for _, org := range defindeInfo.Orgs {
				if strings.TrimSpace(org.CaIp) == strings.TrimSpace(ipParam.IP) {
					go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(org.CaPort)), installChan)
					goroutineNum = goroutineNum + 1
				}
				for _, peer := range org.Peers {
					if strings.TrimSpace(peer.PeerIp) == strings.TrimSpace(ipParam.IP) {
						go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(peer.PostPort)), installChan)
						go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, strconv.Itoa(peer.EventPort)), installChan)
						goroutineNum = goroutineNum + 2
						if peer.JoinCouch {
							go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.ADDPORT, "5984"), installChan) //默认couchdb 端口
							goroutineNum = goroutineNum + 1
						}
					}
				}
			}
			// 	循环得到chan 的值
		PortChan:
			for {
				select {
				case msg := <-installChan:
					msgMap := make(map[string]interface{})
					err = json.Unmarshal([]byte(msg), &msgMap)
					if err != nil {
						utils.ResponseJson(400, "环境未安装完整，请重新安装："+err.Error(), "", writer)
						return
					}
					if int(msgMap["code"].(float64)) != 200 { // interface{} into int
						utils.ResponseJson(400, "环境未安装完整，请重新安装", "", writer)
						return
					}
					// 处理后减一
					goroutineNum--
					if goroutineNum <= 0 {
						break PortChan
					}
				}
			}

			// 	重启firewall
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.STOPFIREWALL, installChan) //临时关闭防火墙constant.RESTATFIREWALL

			msg := <-installChan //获得chan 传值
			fmt.Println("重启firewall：")
			fmt.Println(msg)
			msgMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(msg), &msgMap)
			if err != nil {
				utils.ResponseJson(400, "环境未安装完整，请重新安装"+err.Error(), "", writer)
				return
			}
			if int(msgMap["code"].(float64)) != 200 { // interface{} into int
				utils.ResponseJson(400, "环境未安装完整，请重新安装", "", writer)
				return
			}
		}

		if data["type"].(int) == 0 { // 1 centos 0 ubuntu 2 其它
			// install docker docker-compose nodejs
			installChan := make(chan string, 4)

			// go routine start ===========================================================
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_INSTALLDOCKER, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_INSTALLDOCKERCOMPOSE, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_INSTALLNODEJS, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_INSTALLGIT, installChan)
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_INSTALLJQ, installChan)
			// 根据IP 获得该项目中所有的端口号
			defindeInfo := module.Define{}
			if ipParam.ID != "" {
				_, err = db.GetReadSetupInfo(ipParam.ID, &defindeInfo)
				if err != nil {
					utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
					return
				}
			}
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "4000"), installChan)
			goroutineNum = goroutineNum + 1
			// 	order 端口安装 kafka
			if strings.TrimSpace(defindeInfo.KafkaIp) == strings.TrimSpace(ipParam.IP) {
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9092"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9093"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9192"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9193"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9292"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9293"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9392"), installChan)
				go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "9393"), installChan)
				goroutineNum = goroutineNum + 8
			}
			// 	order 端口安装 firewall
			for _, order := range defindeInfo.Orders {
				if strings.TrimSpace(order.OrderIp) == strings.TrimSpace(ipParam.IP) {
					go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, strconv.Itoa(order.OrderPort)), installChan)
					goroutineNum = goroutineNum + 1
				}
			}
			// 	peer 端口安装 firewall
			for _, org := range defindeInfo.Orgs {
				if strings.TrimSpace(org.CaIp) == strings.TrimSpace(ipParam.IP) {
					go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, strconv.Itoa(org.CaPort)), installChan)
					goroutineNum = goroutineNum + 1
				}
				for _, peer := range org.Peers {
					if strings.TrimSpace(peer.PeerIp) == strings.TrimSpace(ipParam.IP) {
						go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, strconv.Itoa(peer.PostPort)), installChan)
						go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, strconv.Itoa(peer.EventPort)), installChan)
						goroutineNum = goroutineNum + 2
						if peer.JoinCouch {
							go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, fmt.Sprintf(constant.U_ADDPORT, "5984"), installChan) //默认couchdb 端口
							goroutineNum = goroutineNum + 1
						}
					}
				}
			}
			// 	循环得到chan 的值
		UbuntuPortChan:
			for {
				select {
				case msg := <-installChan:
					msgMap := make(map[string]interface{})
					err = json.Unmarshal([]byte(msg), &msgMap)
					if err != nil {
						utils.ResponseJson(400, "环境未安装完整，请重新安装："+err.Error(), "", writer)
						return
					}
					if int(msgMap["code"].(float64)) != 200 { // interface{} into int
						utils.ResponseJson(400, "环境未安装完整，请重新安装", "", writer)
						return
					}
					// 处理后减一
					goroutineNum--
					if goroutineNum <= 0 {
						break UbuntuPortChan
					}
				}
			}

			// 	重启firewall
			go utils.InstallIPEnv(ipParam.IP, constant.SERVERPORT, constant.U_STOPFIREWALL, installChan) //临时关闭防火墙 constant.U_RESTATFIREWALL

			msg := <-installChan //获得chan 传值
			fmt.Println("重启firewall：")
			fmt.Println(msg)
			msgMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(msg), &msgMap)
			if err != nil {
				utils.ResponseJson(400, "环境未安装完整，请重新安装"+err.Error(), "", writer)
				return
			}
			if int(msgMap["code"].(float64)) != 200 { // interface{} into int
				utils.ResponseJson(400, "环境未安装完整，请重新安装", "", writer)
				return
			}
		}

		if data["type"].(int) == 2 { //其他系统
			utils.ResponseJson(400, "系统不是centos或者ubuntu，请更换系统OS", "", writer)
			return
		}
	}

	// 完成
	utils.ResponseJson(200, "环境已经安装完成", "", writer)
	return

}

// chaincode配置文件和工程以及环境
func chaincodeHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得chaincodeHandler传参："+err.Error(), "", writer)
		return
	}

	//保存日志
	curUser := request.Header.Get("username")
	SaveLog(curUser, string(body), request.URL.Path) //保存日志

	ccParam := module.ChainCode{}
	err = json.Unmarshal(body, &ccParam)

	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析chaincodeHandler传参："+err.Error(), "", writer)
		return
	}

	log.Println("ccParam====:", ccParam)

	defindeInfo := module.Define{}
	rev, err := db.GetReadSetupInfo(ccParam.PID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "channel根据ID get define："+err.Error(), "", writer)
		return
	}

	ccChan := make(chan string)
	// chaincode 目录
	go utils.ChaincodeConfig(defindeInfo.Orders[0].OrderIp, constant.SERVERPORT, fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/cc/src"), ccParam.CCName, ccParam.CCGitUrl, ccChan)

	msg := <-ccChan //获得lunch chan 传值
	fmt.Println(msg)
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		utils.ResponseJson(400, "配置cc错误："+err.Error(), "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "配置cc错误", "", writer)
		return
	}
	// get path
	if msgMap["data"] != nil {
		pathData := msgMap["data"].(map[string]interface{})
		if pathData["path"] != nil {
			ccParam.CCPath = strings.Replace(pathData["path"].(string), fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/cc/src"), "", -1)
		} else {
			utils.ResponseJson(400, "发送智能合约错误", "", writer)
			return
		}
	} else {
		utils.ResponseJson(400, "发送智能合约错误", "", writer)
		return
	}
	/////////////////////////////////////////更新channel 下面的智能合约=================================================
	upgrade := false
	hasInstallChannel := false
	// 保存cc
	chIndex := 0
	ccIndex := 0
	for i, channel := range defindeInfo.AddChannels {
		// channel id 相同
		if ccParam.ChannelId == channel.ChannelId {
			for j, cc := range channel.ChainCodes {
				// 判断是否正在使用,channel 是否已经安装
				if defindeInfo.AddChannels[i].ChainCodes[j].Using == true && hasInstallChannel == false {
					hasInstallChannel = true
				}
				// defindeInfo.AddChannels[i].ChainCodes[j].Using = false
				if cc.CCName == ccParam.CCName && cc.CCVersion == ccParam.CCVersion {
					// 返回
					utils.ResponseJson(400, "版本已存在", "", writer)
					return
				}
				if cc.CCName == ccParam.CCName {
					defindeInfo.AddChannels[i].ChainCodes[j].CCGitUrl = ccParam.CCGitUrl
					defindeInfo.AddChannels[i].ChainCodes[j].CCVersion = ccParam.CCVersion
					defindeInfo.AddChannels[i].ChainCodes[j].CCPath = ccParam.CCPath
					defindeInfo.AddChannels[i].ChainCodes[j].Using = false
					chIndex = i
					ccIndex = j
					upgrade = true
					break
				}
			}
			// 如果不是更新新增chaincode
			if upgrade == false {
				ccObj := module.CC{}
				ccObj.CCGitUrl = ccParam.CCGitUrl
				ccObj.CCName = ccParam.CCName
				ccObj.CCVersion = ccParam.CCVersion
				ccObj.CCPath = ccParam.CCPath
				ccObj.Using = false
				chIndex = i
				// channel chaincode index
				defindeInfo.AddChannels[i].ChainCodes = append(defindeInfo.AddChannels[i].ChainCodes, ccObj)
				if len(defindeInfo.AddChannels[i].ChainCodes) > 0 {
					ccIndex = len(defindeInfo.AddChannels[i].ChainCodes) - 1
				} else {
					ccIndex = 0
				}
			}
			break
		}
	}
	// 保持 智能合约信息
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存项目定义信息错误："+err.Error(), "", writer)
		return
	}

	/////////////////////////////////////////更新channel 下面的智能合约  end=================================================
	/////////////////////////////==========================启动链========================
	jqShell := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// login shell
	loginShell := ""
	defaultToken := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == ccParam.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					if defaultToken == "" {
						defaultToken = org.OrgId
					}
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(curl -s -X POST http://localhost:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=%s&password=password&orgName=%s')`, org.OrgId, org.OrgId, org.OrgId))
					//&channelName=%s , channel.ChannelId
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(echo $%s_TOKEN | jq ".token" | sed "s/\"//g")`, org.OrgId, org.OrgId))
				}
			}
		}
	}

	// 读取传送的数据
	createChannel := ""

	createChannel = fmt.Sprintln(createChannel, fmt.Sprintf(`curl -s -X POST \
		  http://localhost:4000/channels \
		  -H "authorization: Bearer $%s_TOKEN" \
		  -H "content-type: application/json" \
		  -d '{"channelName":"%s"}'`, defaultToken, ccParam.ChannelId))

	joinChannel := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == ccParam.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					peerStr := ""
					for _, peer := range org.Peers {
						if peerStr == "" {
							peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
						} else {
							peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
						}
					}
					joinChannel = fmt.Sprintln(joinChannel, fmt.Sprintf(`curl -s -X POST \
				http://localhost:4000/channels/peers \
				-H "authorization: Bearer $%s_TOKEN" \
				-H "content-type: application/json" \
				-d '{"peers": [%s],"channelName":"%s"}'`, org.OrgId, peerStr, ccParam.ChannelId))
				}
			}
		}
	}

	installChaincode := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == ccParam.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					peerStr := ""
					for _, peer := range org.Peers {
						if peerStr == "" {
							peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
						} else {
							peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
						}
					}
					installChaincode = fmt.Sprintln(installChaincode, fmt.Sprintf(`curl -s -X POST \
			http://localhost:4000/chaincodes \
			-H "authorization: Bearer $%s_TOKEN" \
			-H "content-type: application/json" \
			-d '{
			  "peers": [%s],
			  "channelName":"%s",
			  "chaincodeName":"%s",
			  "chaincodePath":"%s",
			  "chaincodeVersion":"%s"
		  }'`, org.OrgId, peerStr, ccParam.ChannelId, ccParam.CCName, ccParam.CCPath, ccParam.CCVersion))
				}
			}
		}
	}

	instantiateChainCode := ""
	instantiateChainCode = fmt.Sprintln(instantiateChainCode, fmt.Sprintf(`curl -s -X POST \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $%s_TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "channelName":"%s",
		  "chaincodeName":"%s",
		  "chaincodeVersion":"%s",
		  "args":[]
	  }'`, defaultToken, ccParam.ChannelId, ccParam.CCName, ccParam.CCVersion))

	upgradeChaincode := ""

	upgradeChaincode = fmt.Sprintln(upgradeChaincode, fmt.Sprintf(`curl -s -X PUT \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $%s_TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "channelName":"%s",
		  "chaincodeName":"%s",
		  "chaincodeVersion":"%s",
		  "args":[]
	  }'`, defaultToken, ccParam.ChannelId, ccParam.CCName, ccParam.CCVersion))

	allChannel := fmt.Sprintln(jqShell, loginShell, createChannel, fmt.Sprintln("sleep 5"), joinChannel, installChaincode, instantiateChainCode)
	//如果升级智能合约，修改allChannel
	if upgrade == true {
		allChannel = fmt.Sprintln(jqShell, loginShell, installChaincode, upgradeChaincode)
	}
	// 已经安装了channel 和 join org
	if hasInstallChannel == true && upgrade == false {
		allChannel = fmt.Sprintln(jqShell, loginShell, installChaincode, instantiateChainCode)
	}

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, allChannel) //请求执行shell命令
	msgMap = make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		log.Println("生成channel 错误")
		utils.ResponseJson(400, "生成channel 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		log.Println("生成channel 错误")
		utils.ResponseJson(400, "生成channel 错误", "", writer)
		return
	}
	// 重新获得rev
	rev, err = db.GetReadSetupInfo(ccParam.PID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}
	//////////////////////////////=========================启动链=========================
	/////========================修改状态=======================================
	defindeInfo.Status = constant.CCED //修改状态，安装完智能合约
	defindeInfo.HasApp = true          //已经部署了智能合约
	log.Println("========================部署智能合约返回============================================")
	defindeInfo.AddChannels[chIndex].ChainCodes[ccIndex].Using = true
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存项目定义信息错误："+err.Error(), "", writer)
		return
	}
	/////========================修改状态=======================================
	// 完成
	utils.ResponseJson(200, "智能合约安装完成", "", writer)
	return
}

// UPLOAD FILE
func uploadChaincode(writer http.ResponseWriter, request *http.Request) {

	request.ParseMultipartForm(32 << 20)
	file, handler, err := request.FormFile("uploadfile")
	if err != nil {
		utils.ResponseJson(500, "上传附件："+err.Error(), "", writer)
		return
	}
	defer file.Close()

	projectId := request.FormValue("pid")
	ccid := request.FormValue("id")
	version := request.FormValue("ccVersion")
	channelId := request.FormValue("channelId")
	ccname := request.FormValue("ccName")

	oldCC := module.ChainCode{}
	oldCC.ID = ccid
	oldCC.PID = projectId
	oldCC.ChannelId = channelId
	oldCC.CCVersion = version
	oldCC.CCName = ccname

	// 根据IP 获得该项目中所有的端口号
	defindeInfo := module.Define{}
	// if ipParam.ID != "" {
	rev, err := db.GetReadSetupInfo(projectId, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}
	ccPath := fmt.Sprint(constant.ROOTPATH, "/", projectId, "/cc/src")
	err = utils.CheckAndCreatePath(ccPath)
	if err != nil {
		utils.ResponseJson(400, "生成cc路径："+err.Error(), "", writer)
		return
	}
	filePath := fmt.Sprint(ccPath, "/", handler.Filename)
	// 如果有文件，删除
	err = utils.CheckFileAndRemove(filePath)
	if err != nil {
		utils.ResponseJson(400, "删除已有附件："+err.Error(), "", writer)
		return
	}
	// 创建文件路径，拷贝上传附件
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		utils.ResponseJson(400, "上传附件："+err.Error(), "", writer)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	//保存chaincode 路径信息
	oldCC.CCPath = filePath

	//保存日志
	ccstring, _ := json.Marshal(oldCC)
	curUser := request.Header.Get("username")
	SaveLog(curUser, string(ccstring), request.URL.Path) //保存日志

	////////////////***************************上传智能合约
	hostIp := defindeInfo.Orders[0].OrderIp
	ccChan := make(chan string, 1)
	go utils.SendCCFile(handler.Filename, filePath, hostIp, constant.SERVERPORT, ccPath, ccname, ccChan)
	msg := <-ccChan //获得chan 传值
	// 解析go 发送的值
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		utils.ResponseJson(400, "智能合约压缩文件发送失败，请重试", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "智能合约压缩文件发送失败，请重试", "", writer)
		return
	}
	// get path
	if msgMap["data"] != nil {
		pathData := msgMap["data"].(map[string]interface{})
		if pathData["path"] != nil {
			oldCC.CCPath = strings.Replace(pathData["path"].(string), fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID, "/cc/src/"), "", -1)
		} else {
			utils.ResponseJson(400, "发送智能合约错误", "", writer)
			return
		}
	} else {
		utils.ResponseJson(400, "发送智能合约错误", "", writer)
		return
	}
	/////////////////**************************
	/////////////////////////////////////////更新channel 下面的智能合约=================================================
	upgrade := false
	hasInstallChannel := false
	// 保存cc
	chIndex := 0
	ccIndex := 0
	for i, channel := range defindeInfo.AddChannels {
		// channel id 相同
		if oldCC.ChannelId == channel.ChannelId {
			for j, cc := range channel.ChainCodes {
				// 判断是否正在使用,channel 是否已经安装
				if defindeInfo.AddChannels[i].ChainCodes[j].Using == true && hasInstallChannel == false {
					hasInstallChannel = true
				}
				// defindeInfo.AddChannels[i].ChainCodes[j].Using = false
				if cc.CCName == oldCC.CCName && cc.CCVersion == oldCC.CCVersion {
					// 返回
					utils.ResponseJson(400, "版本已存在", "", writer)
					return
				}
				if cc.CCName == oldCC.CCName {
					defindeInfo.AddChannels[i].ChainCodes[j].CCGitUrl = oldCC.CCGitUrl
					defindeInfo.AddChannels[i].ChainCodes[j].CCVersion = oldCC.CCVersion
					defindeInfo.AddChannels[i].ChainCodes[j].CCPath = oldCC.CCPath
					defindeInfo.AddChannels[i].ChainCodes[j].Using = false
					chIndex = i
					ccIndex = j
					upgrade = true
					break
				}
			}
			// 如果不是更新新增chaincode
			if upgrade == false {
				ccObj := module.CC{}
				ccObj.CCGitUrl = oldCC.CCGitUrl
				ccObj.CCName = oldCC.CCName
				ccObj.CCVersion = oldCC.CCVersion
				ccObj.CCPath = oldCC.CCPath
				ccObj.Using = false
				chIndex = i
				defindeInfo.AddChannels[i].ChainCodes = append(defindeInfo.AddChannels[i].ChainCodes, ccObj)
				if len(defindeInfo.AddChannels[i].ChainCodes) > 0 {
					ccIndex = len(defindeInfo.AddChannels[i].ChainCodes) - 1
				} else {
					ccIndex = 0
				}
			}
			break
		}
	}
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存项目定义信息错误："+err.Error(), "", writer)
		return
	}
	/////////////////////////////////////////更新channel 下面的智能合约  end=================================================

	/////////////////////////////==========================启动链========================
	jqShell := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// login shell
	loginShell := ""
	defaultToken := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == oldCC.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					if defaultToken == "" {
						defaultToken = org.OrgId
					}
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(curl -s -X POST http://localhost:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=%s&password=password&orgName=%s')`, org.OrgId, org.OrgId, org.OrgId))
					//&channelName=%s , channel.ChannelId
					loginShell = fmt.Sprintln(loginShell, fmt.Sprintf(`%s_TOKEN=$(echo $%s_TOKEN | jq ".token" | sed "s/\"//g")`, org.OrgId, org.OrgId))
				}
			}
		}
	}

	// 读取传送的数据
	createChannel := ""

	createChannel = fmt.Sprintln(createChannel, fmt.Sprintf(`curl -s -X POST \
		  http://localhost:4000/channels \
		  -H "authorization: Bearer $%s_TOKEN" \
		  -H "content-type: application/json" \
		  -d '{"channelName":"%s"}'`, defaultToken, oldCC.ChannelId))

	joinChannel := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == oldCC.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					peerStr := ""
					for _, peer := range org.Peers {
						if peerStr == "" {
							peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
						} else {
							peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
						}
					}
					joinChannel = fmt.Sprintln(joinChannel, fmt.Sprintf(`curl -s -X POST \
				http://localhost:4000/channels/peers \
				-H "authorization: Bearer $%s_TOKEN" \
				-H "content-type: application/json" \
				-d '{"peers": [%s],"channelName":"%s"}'`, org.OrgId, peerStr, oldCC.ChannelId))
				}
			}
		}
	}

	installChaincode := ""
	for _, org := range defindeInfo.Orgs {
		for _, channel := range defindeInfo.AddChannels {
			if channel.ChannelId == oldCC.ChannelId {
				if utils.Contains(channel.IncludeOrgs, org.OrgId) {
					peerStr := ""
					for _, peer := range org.Peers {
						if peerStr == "" {
							peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
						} else {
							peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
						}
					}
					installChaincode = fmt.Sprintln(installChaincode, fmt.Sprintf(`curl -s -X POST \
			http://localhost:4000/chaincodes \
			-H "authorization: Bearer $%s_TOKEN" \
			-H "content-type: application/json" \
			-d '{
			  "peers": [%s],
			  "channelName":"%s",
			  "chaincodeName":"%s",
			  "chaincodePath":"%s",
			  "chaincodeVersion":"%s"
		  }'`, org.OrgId, peerStr, oldCC.ChannelId, oldCC.CCName, oldCC.CCPath, oldCC.CCVersion))
				}
			}
		}
	}

	instantiateChainCode := ""
	instantiateChainCode = fmt.Sprintln(instantiateChainCode, fmt.Sprintf(`curl -s -X POST \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $%s_TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "channelName":"%s",
		  "chaincodeName":"%s",
		  "chaincodeVersion":"%s",
		  "args":[]
	  }'`, defaultToken, oldCC.ChannelId, oldCC.CCName, oldCC.CCVersion))

	upgradeChaincode := ""

	upgradeChaincode = fmt.Sprintln(upgradeChaincode, fmt.Sprintf(`curl -s -X PUT \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $%s_TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "channelName":"%s",
		  "chaincodeName":"%s",
		  "chaincodeVersion":"%s",
		  "args":[]
	  }'`, defaultToken, oldCC.ChannelId, oldCC.CCName, oldCC.CCVersion))

	allChannel := fmt.Sprintln(jqShell, loginShell, createChannel, fmt.Sprintln("sleep 5"), joinChannel, installChaincode, instantiateChainCode)
	//如果升级智能合约，修改allChannel
	if upgrade == true {
		allChannel = fmt.Sprintln(jqShell, loginShell, installChaincode, upgradeChaincode)
	}
	// 已经安装了channel 和 join org
	if hasInstallChannel == true && upgrade == false {
		allChannel = fmt.Sprintln(jqShell, loginShell, installChaincode, instantiateChainCode)
	}
	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, allChannel) //请求执行shell命令
	msgMap = make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		log.Println("========================生成channel 错误 && 安装智能合约错误==============================")
		utils.ResponseJson(400, "生成channel 错误 && 安装智能合约错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		log.Println("========================生成channel 错误 && 安装智能合约错误==============================")
		utils.ResponseJson(400, "生成channel 错误 && 安装智能合约错误", "", writer)
		return
	}
	// 重新获得rev
	rev, err = db.GetReadSetupInfo(projectId, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}
	//////////////////////////////=========================启动链=========================
	////////////////////////=============修改状态===============================================
	defindeInfo.Status = constant.CCED //修改状态， 已经安装智能合约
	defindeInfo.HasApp = true          //已经部署了智能合约
	log.Println("========================部署智能合约返回============================================")
	defindeInfo.AddChannels[chIndex].ChainCodes[ccIndex].Using = true
	err = db.SaveSetupInfo(defindeInfo, defindeInfo.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存项目定义信息错误："+err.Error(), "", writer)
		return
	}
	////////////////////////=============修改状态===============================================
	utils.ResponseJson(200, "智能合约安装完成", "", writer)
	return
}
