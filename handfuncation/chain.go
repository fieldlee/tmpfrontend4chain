package handfuncation

import (
	"encoding/json"
	"fmt"
	"frontend4chain/constant"
	"frontend4chain/db"
	"frontend4chain/module"
	"frontend4chain/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func HandlerChainAll(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/chain/createChannel":
		createChannel(writer, request)
		return
	case "/chain/joinChannel":
		joinChannel(writer, request)
		return
	case "/chain/installChaincode":
		installChaincode(writer, request)
		return
	case "/chain/instaniateChaincode":
		instaniateChaincode(writer, request)
		return
	case "/chain/upgradeChaincode":
		upgradeChaincode(writer, request)
		return
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}
}

// 创建channel
func createChannel(writer http.ResponseWriter, request *http.Request) {

	createChannel := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// 读取传送的数据
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得createChannel传参："+err.Error(), "", writer)
		return
	}

	chainParam := module.ChainParam{}
	err = json.Unmarshal(body, &chainParam)

	log.Println(chainParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析chain传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	_, err = db.GetReadSetupInfo(chainParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}

	createChannel = fmt.Sprint(createChannel, fmt.Sprintf(`TOKEN=$(curl -s -X POST \
		http://localhost:4000/login \
		-H "content-type: application/x-www-form-urlencoded" \
		-d 'username=%s&password=password&orgName=%s&channelName=%s')
	  echo $TOKEN
	  TOKEN=$(echo $TOKEN | jq ".token" | sed "s/\"//g")`, defindeInfo.Orgs[0].OrgId, defindeInfo.Orgs[0].OrgId, chainParam.ChannelId))

	channelpath := fmt.Sprint(constant.ROOTPATH, "/", defindeInfo.ID)

	createChannel = fmt.Sprint(createChannel, fmt.Sprintf(`curl -s -X POST \
		  http://localhost:4000/channels \
		  -H "authorization: Bearer $TOKEN" \
		  -H "content-type: application/json" \
		  -d '{
			"channelName":"%s",
			"channelConfigPath":"%s"}'`, chainParam.ChannelId, channelpath))

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, createChannel) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "生成channel 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "生成channel 错误", "", writer)
		return
	}
	utils.ResponseJson(200, "生成channel完成", "", writer)
	return
}

// join channel
func joinChannel(writer http.ResponseWriter, request *http.Request) {
	joinChannel := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// 读取传送的数据
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得createChannel传参："+err.Error(), "", writer)
		return
	}

	chainParam := module.ChainParam{}
	err = json.Unmarshal(body, &chainParam)

	log.Println(chainParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析chain传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	_, err = db.GetReadSetupInfo(chainParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}

	for _, org := range defindeInfo.Orgs {
		joinChannel = fmt.Sprintln(joinChannel, fmt.Sprintf(`%s_TOKEN=$(curl -s -X POST http://localhost:4000/login \ 
			-H "content-type: application/x-www-form-urlencoded" \
			-d 'username=%s&password=password&orgName=%s&channelName=%s')
		  %s_TOKEN=$(echo $%s_TOKEN | jq ".token" | sed "s/\"//g")`, org.OrgId, org.OrgId, org.OrgId, org.OrgId, org.OrgId))

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
			-d '{"peers": [%s]}'`, org.OrgId, peerStr))
	}

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, joinChannel) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "joinchannel 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "joinchannel 错误", "", writer)
		return
	}
	utils.ResponseJson(200, "joinchannel", "", writer)
	return
}

// installChaincode
func installChaincode(writer http.ResponseWriter, request *http.Request) {
	installChaincode := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// 读取传送的数据
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得createChannel传参："+err.Error(), "", writer)
		return
	}

	chainParam := module.ChainCodeParam{}
	err = json.Unmarshal(body, &chainParam)

	log.Println(chainParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析chain传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	_, err = db.GetReadSetupInfo(chainParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}
	for _, org := range defindeInfo.Orgs {
		installChaincode = fmt.Sprintln(installChaincode, fmt.Sprintf(`%s_TOKEN=$(curl -s -X POST http://localhost:4000/login \ 
			-H "content-type: application/x-www-form-urlencoded" \
			-d 'username=%s&password=password&orgName=%s')
		  %s_TOKEN=$(echo $%s_TOKEN | jq ".token" | sed "s/\"//g")`, org.OrgId, org.OrgId, org.OrgId, org.OrgId, org.OrgId))

		peerStr := ""
		for _, peer := range org.Peers {
			if peerStr == "" {
				peerStr = fmt.Sprint(peerStr, fmt.Sprintf(`"%s"`, peer.PeerId))
			} else {
				peerStr = fmt.Sprint(peerStr, ",", fmt.Sprintf(`"%s"`, peer.PeerId))
			}
		}
		// chaincodePath :=
		installChaincode = fmt.Sprintln(installChaincode, fmt.Sprintf(`curl -s -X POST \
			http://localhost:4000/chaincodes \
			-H "authorization: Bearer $%s_TOKEN" \
			-H "content-type: application/json" \
			-d '{
			  "peers": [%s],
			  "chaincodeName":"%s",
			  "chaincodePath":"%s",
			  "chaincodeVersion":"v%s"
		  }'`, org.OrgId, peerStr, chainParam.ChainCodeName, chainParam.ChainCodeName, chainParam.ChainCodeVersion))
	}

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, installChaincode) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	utils.ResponseJson(200, "installChaincode", "", writer)
	return
}

// 初始化智能合约
func instaniateChaincode(writer http.ResponseWriter, request *http.Request) {
	instantiateChainCode := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// 读取传送的数据
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得createChannel传参："+err.Error(), "", writer)
		return
	}

	chainParam := module.ChainCodeParam{}
	err = json.Unmarshal(body, &chainParam)

	log.Println(chainParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析chain传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	_, err = db.GetReadSetupInfo(chainParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}

	instantiateChainCode = fmt.Sprintln(instantiateChainCode, fmt.Sprintf(`TOKEN=$(curl -s -X POST \
		http://localhost:4000/login \
		-H "content-type: application/x-www-form-urlencoded" \
		-d 'username=%s&password=password&orgName=%s')
	  echo $TOKEN
	  TOKEN=$(echo $TOKEN | jq ".token" | sed "s/\"//g")`, defindeInfo.Orgs[0].OrgId, defindeInfo.Orgs[0].OrgId))

	instantiateChainCode = fmt.Sprintln(instantiateChainCode, fmt.Sprintf(`curl -s -X POST \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "chaincodeName":"%s",
		  "chaincodeVersion":"v%s",
		  "args":[]
	  }'`, defindeInfo.Orgs[0].OrgId, chainParam.ChainCodeName, chainParam.ChainCodeVersion))

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, instantiateChainCode) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	utils.ResponseJson(200, "installChaincode", "", writer)
	return
}

// 升级智能合约
func upgradeChaincode(writer http.ResponseWriter, request *http.Request) {
	upgradeChaincode := fmt.Sprintln(
		`jq --version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
			echo
			exit 1
		fi`)
	// 读取传送的数据
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得createChannel传参："+err.Error(), "", writer)
		return
	}

	chainParam := module.ChainCodeParam{}
	err = json.Unmarshal(body, &chainParam)

	log.Println(chainParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析chain传参："+err.Error(), "", writer)
		return
	}

	defindeInfo := module.Define{}
	_, err = db.GetReadSetupInfo(chainParam.ID, &defindeInfo)
	if err != nil {
		utils.ResponseJson(400, "根据ID get define："+err.Error(), "", writer)
		return
	}

	upgradeChaincode = fmt.Sprintln(upgradeChaincode, fmt.Sprintf(`TOKEN=$(curl -s -X POST \
		http://localhost:4000/login \
		-H "content-type: application/x-www-form-urlencoded" \
		-d 'username=%s&password=password&orgName=%s')
	  echo $TOKEN
	  TOKEN=$(echo $TOKEN | jq ".token" | sed "s/\"//g")`, defindeInfo.Orgs[0].OrgId, defindeInfo.Orgs[0].OrgId))

	upgradeChaincode = fmt.Sprintln(upgradeChaincode, fmt.Sprintf(`curl -s -X PUT \
		http://localhost:4000/channels/chaincodes \
		-H "authorization: Bearer $TOKEN" \
		-H "content-type: application/json" \
		-d '{
		  "chaincodeName":"%s",
		  "chaincodeVersion":"v%s",
		  "args":[]
	  }'`, defindeInfo.Orgs[0].OrgId, chainParam.ChainCodeName, chainParam.ChainCodeVersion))

	message := utils.ReqShellServer(defindeInfo.Orders[0].OrderIp, upgradeChaincode) //请求执行shell命令
	msgMap := make(map[string]interface{})
	err = json.Unmarshal(message, &msgMap)
	if err != nil {
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	if int(msgMap["code"].(float64)) != 200 { // interface{} into int
		utils.ResponseJson(400, "installChaincode 错误", "", writer)
		return
	}
	utils.ResponseJson(200, "installChaincode", "", writer)
	return
}
