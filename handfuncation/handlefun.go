package handfuncation

import (
	"encoding/json"
	"frontend4chain/constant"
	"frontend4chain/db"
	"frontend4chain/module"
	"frontend4chain/utils"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

func HandlerAll(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	switch request.Method {
	case "GET":
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	case "POST":
		if strings.Contains(request.URL.Path, "/getlist") { //获得项目
			if request.Header.Get("subject") == "project" {
				listHandler(writer, request)
			} else {
				utils.ResponseJson(403, "该用户无权限，请重新登录", "", writer)
				return
			}
			return
		}
		if strings.Contains(request.URL.Path, "/getproject") { //获得项目
			if request.Header.Get("subject") == "project" {
				getProjectByIDHandler(writer, request)
			} else {
				utils.ResponseJson(403, "该用户无权限，请重新登录", "", writer)
				return
			}
			return
		}
		if strings.Contains(request.URL.Path, "/setup/") {
			if request.Header.Get("subject") == "project" {
				HandlerSetupAll(writer, request)
			} else {
				utils.ResponseJson(403, "该用户无权限，请重新登录", "", writer)
				return
			}
		}
		if strings.Contains(request.URL.Path, "/fabric/") {
			if request.Header.Get("subject") == "project" {
				HandlerFabircAll(writer, request)
			} else {
				utils.ResponseJson(403, "该用户无权限，请重新登录", "", writer)
				return
			}
		}
		// 浏览器查看
		if strings.Contains(request.URL.Path, "/explorer/") {
			HandlerExplorer(writer, request)
		}
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}
}

func listHandler(writer http.ResponseWriter, request *http.Request) {

	listBytes := make([]byte, 1)
	// 管理员
	// if strings.TrimSpace(curuser.(module.User).Role) == "admin" {
	// 	list, err := db.GetReadSetupList()
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
	// 		return
	// 	}
	// 	listBytes, err = json.Marshal(list)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
	// 		return
	// 	}
	// }
	// 普通用户
	isManager, curusername, err := utils.IsManager(request)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得当前人错误："+err.Error(), "", writer)
		return
	}
	if isManager == false {
		list, err := db.GetReadSetupListByUser(curusername)
		var sortableList module.DefineList
		sortableList = list
		sort.Sort(sortableList)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
			return
		}
		listBytes, err = json.Marshal(sortableList)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
			return
		}
	}
	// 维护人员
	if isManager == true {
		list, err := db.GetReadSetupListByStatus(constant.SAVEED)
		var sortableList module.DefineList
		sortableList = list
		sort.Sort(sortableList)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
			return
		}
		listBytes, err = json.Marshal(sortableList)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "获得list"+err.Error(), "", writer)
			return
		}
	}

	utils.ResponseJson(200, "", string(listBytes[:]), writer)
	return
}

func getProjectByIDHandler(writer http.ResponseWriter, request *http.Request) {

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得project参数错误："+err.Error(), "", writer)
		return
	}

	// 获得当前用户信息
	curUser, err := db.GetUser(request.Header.Get("username"))
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "当前用户错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析传参："+err.Error(), "", writer)
		return
	}

	project := module.Define{}
	_, err = db.GetReadSetupInfo(obj["id"].(string), &project)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得项目错误："+err.Error(), "", writer)
		return
	}

	if curUser.(module.User).UserName != project.Manager && curUser.(module.User).Role != "manager" {
		utils.ResponseJson(400, "用户不是该项目管理人员，或不是管理员", "", writer)
		return
	}

	projectBytes, _ := json.Marshal(project)
	utils.ResponseJson(200, "", string(projectBytes[:]), writer)
	return
}
