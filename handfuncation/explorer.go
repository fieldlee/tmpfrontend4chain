package handfuncation

import (
	"bytes"
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

func HandlerExplorer(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/explorer/blocks":
		getBlocks(writer, request)
		return
	case "/explorer/transactions":
		getTX(writer, request)
		return
	case "/explorer/asset":
		getAsset(writer, request)
		return
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}
}

func getBlocks(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
	// 	return
	// }
	// 根据登录的用户名获得项目的
	// currentUser := request.Header.Get("username")

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得block："+err.Error(), "", writer)
		return
	}
	explorerParam := module.ExplorderParam{}
	err = json.Unmarshal(body, &explorerParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析项目id："+err.Error(), "", writer)
		return
	}

	tempPro, err := db.GetSetupInfo(explorerParam.ID)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目信息错误："+err.Error(), "", writer)
		return
	}

	tPro := tempPro.(module.Define)

	tIp := tPro.Orders[0].OrderIp
	port := "4000"
	host := fmt.Sprint("http://", tIp, ":", port)
	blockUrl := fmt.Sprint(host, "/")
	jsonMap := make(map[string]interface{})

	jsonBytes, _ := json.Marshal(jsonMap)

	req, err := http.NewRequest("POST", blockUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		fmt.Println("调用浏览器接口错误：" + blockUrl)
	}
	// Content-Type:application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", utils.HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("调用浏览器接口错误")
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	utils.ResponseJson(200, "", string(message), writer)
	return
}

func getTX(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	//blocknum
	txobj := module.TXParam{}
	err = json.Unmarshal(body, &txobj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析参数："+err.Error(), "", writer)
		return
	}

	blockNum := txobj.Blocknum
	// 根据登录的用户名获得项目的
	// currentUser := request.Header.Get("username")

	tempPro, err := db.GetSetupInfo(txobj.ID)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目信息错误："+err.Error(), "", writer)
		return
	}

	tPro := tempPro.(module.Define)

	tIp := tPro.Orders[0].OrderIp
	port := "4000"
	host := fmt.Sprint("http://", tIp, ":", port)
	txUrl := fmt.Sprint(host, "/")

	jsonMap := make(map[string]interface{})
	jsonMap["block"] = blockNum
	jsonBytes, _ := json.Marshal(jsonMap)

	req, err := http.NewRequest("POST", txUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		fmt.Println("调用浏览器接口错误：" + txUrl)
	}
	// Content-Type:application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", utils.HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("调用浏览器接口错误")
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	utils.ResponseJson(200, "", string(message), writer)
	return
}

func getAsset(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// 根据登录的用户名获得项目的
	// currentUser := request.Header.Get("username")

	// tempPro, err := db.GetSetupInfoByUserName(currentUser)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	utils.ResponseJson(400, "获得项目信息错误："+err.Error(), "", writer)
	// 	return
	// }

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "获得资产："+err.Error(), "", writer)
		return
	}
	explorerParam := module.ExplorderParam{}
	err = json.Unmarshal(body, &explorerParam)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(500, "解析项目id："+err.Error(), "", writer)
		return
	}

	tempPro, err := db.GetSetupInfo(explorerParam.ID)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目信息错误："+err.Error(), "", writer)
		return
	}

	tPro := tempPro.(module.Define)

	tIp := tPro.Orders[0].OrderIp
	port := "4000"
	host := fmt.Sprint("http://", tIp, ":", port)
	assetUrl := fmt.Sprint(host, "/")

	jsonMap := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(jsonMap)

	req, err := http.NewRequest("POST", assetUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		fmt.Println("调用浏览器接口错误：" + assetUrl)
	}
	// Content-Type:application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", utils.HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("调用浏览器接口错误")
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	utils.ResponseJson(200, "", string(message), writer)
	return
}
