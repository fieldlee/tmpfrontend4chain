package handfuncation

import (
	"encoding/json"
	"fmt"
	"frontend4chain/db"
	"frontend4chain/module"
	"frontend4chain/utils"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func HandlerSetupAll(writer http.ResponseWriter, request *http.Request) {
	switch strings.ToLower(request.URL.Path) {
	case "/setup/feedback":
		feedbackHandler(writer, request)
		return
	case "/setup/getfeedback":
		getfeedbackHandler(writer, request)
		return
	case "/setup/announce":
		announceHandler(writer, request)
		return
	case "/setup/getannounce":
		getannounceHandler(writer, request)
		return
	case "/setup/getuser":
		getuserHandler(writer, request)
		return
	case "/setup/userlist":
		userlistHandler(writer, request)
		return
	case "/setup/projectlist":
		projectlistHandler(writer, request)
		return
	case "/setup/loglist":
		loglistHandler(writer, request)
		return
	case "/setup/saveuser":
		saveUserHandler(writer, request)
		return
	case "/setup/emailcode":
		sendemailHandler(writer, request)
		return
	case "/setup/modifypassword":
		modifypasswordHandler(writer, request)
		return
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}
}

// 修改密码
func modifypasswordHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	_, curusername, err := utils.IsManager(request)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得当前用户错误："+err.Error(), "", writer)
		return
	}

	curuser, err := db.GetUser(curusername)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得用户错误："+err.Error(), "", writer)
		return
	}

	user := curuser.(module.User)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得body参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析传参："+err.Error(), "", writer)
		return
	}

	if obj["oldpassword"] != nil {
		oldpassword := obj["oldpassword"].(string)
		if user.Password != utils.Md5(oldpassword) {
			utils.ResponseJson(400, "老密码不对", "", writer)
			return
		}
	} else {
		utils.ResponseJson(400, "老密码不能为空", "", writer)
		return
	}

	if obj["password"] != nil {
		password := obj["password"].(string)
		user.Password = utils.Md5(password)
	} else {
		utils.ResponseJson(400, "新密码不能为空", "", writer)
		return
	}

	tmpUser := module.User{}
	rev, err := db.GetReadUserInfo(user.ID, &tmpUser)
	if err != nil {
		utils.ResponseJson(400, "用户不存在："+err.Error(), "", writer)
		return
	}

	err = db.SaveUser(user, user.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保存用户错误："+err.Error(), "", writer)
		return
	}

	utils.ResponseJson(200, "用户密码修改完成", "", writer)
	return
}

// 个人用户信息修改
func saveUserHandler(writer http.ResponseWriter, request *http.Request) {
	// request.ParseMultipartForm(32 << 20)
	defer request.Body.Close()
	username := request.Header.Get("username")
	curuser, err := db.GetUser(username)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得当前用户错误："+err.Error(), "", writer)
		return
	}
	user := curuser.(module.User)

	// file, handler, err := request.FormFile("uploadfile")
	// if err != nil {
	// 	utils.ResponseJson(400, "上传附件："+err.Error(), "", writer)
	// 	return
	// }
	// defer file.Close()

	// if file != nil {
	// 	filePath := fmt.Sprint("/var/avator/", handler.Filename)
	// 	// 如果有文件，删除
	// 	err = utils.CheckFileAndRemove(filePath)
	// 	if err != nil {
	// 		utils.ResponseJson(400, "删除已有附件："+err.Error(), "", writer)
	// 		return
	// 	}
	// 	// 创建文件路径，拷贝上传附件
	// 	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	// 	if err != nil {
	// 		utils.ResponseJson(400, "上传附件："+err.Error(), "", writer)
	// 		return
	// 	}
	// 	defer f.Close()
	// 	io.Copy(f, file)

	// 	baseImageFile, err := os.Open(filePath)
	// 	if err != nil {
	// 		utils.ResponseJson(400, "上传附件："+err.Error(), "", writer)
	// 		return
	// 	}
	// 	fInfo, err := baseImageFile.Stat()
	// 	if err != nil {
	// 		utils.ResponseJson(400, "上传附件："+err.Error(), "", writer)
	// 		return
	// 	}
	// 	var size int64 = fInfo.Size()
	// 	buf := make([]byte, size)

	// 	// read file content into buffer
	// 	fReader := bufio.NewReader(baseImageFile)
	// 	fReader.Read(buf)

	// 	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	// 	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	// 	user.Avator = imgBase64Str
	// }

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得body参数错误："+err.Error(), "", writer)
		return
	}

	obj := make(map[string]interface{})
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "解析传参："+err.Error(), "", writer)
		return
	}
	// 打印obj
	log.Println("obj", obj)

	if obj["telno"] != nil {
		telno := obj["telno"].(string)
		if obj["telnoCode"] != nil {
			telnoCode := obj["telnoCode"].(string)
			if len(strings.TrimSpace(telno)) >= 11 {
				telObj, err := db.GetTelCodeInfo(telno)
				if err != nil {
					utils.ResponseJson(400, "获得手机验证码错误: "+err.Error(), "", writer)
					return
				}
				if telObj == nil {
					utils.ResponseJson(400, "手机验证码有误", "", writer)
					return
				}
				if strings.ToLower(telObj.(module.TelCode).VerifyCode) != strings.ToLower(telnoCode) {
					utils.ResponseJson(400, "手机验证码错误", "", writer)
					return
				}
				user.TelNo = telno
			}
		}
	}
	if obj["email"] != nil {
		email := obj["email"].(string)
		if obj["emailCode"] != nil {
			emailCode := obj["emailCode"].(string)
			e, err := db.GetTelCodeInfo(email)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, "获得邮箱验证码错误："+err.Error(), "", writer)
				return
			}
			if strings.ToLower(e.(module.TelCode).VerifyCode) != strings.ToLower(emailCode) {
				utils.ResponseJson(400, "邮箱地址的验证码错误请重新输入", "", writer)
				return
			}
			user.Email = email
		}
	}
	// telno := obj["telno"]
	// email := obj["email"]
	// user.Email = email
	// user.TelNo = telno
	// 打印user
	log.Println("user:", user)
	tmpUser := module.User{}
	rev, err := db.GetReadUserInfo(user.ID, &tmpUser)
	if err != nil {
		utils.ResponseJson(400, "用户不存在："+err.Error(), "", writer)
		return
	}

	err = db.SaveUser(user, user.ID, rev)
	if err != nil {
		utils.ResponseJson(400, "保持用户错误："+err.Error(), "", writer)
		return
	}

	utils.ResponseJson(200, "用户信息修改完成", "", writer)
	return
}

// 发送邮件code
func sendemailHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, err.Error(), "", writer)
		return
	}
	var sendMap map[string]interface{}
	err = json.Unmarshal([]byte(body), &sendMap)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, err.Error(), "", writer)
		return
	}
	email := sendMap["email"].(string)

	// tel, _ := db.GetTelCodeInfo(email)
	randomString := utils.RandomString(4)
	sendHtml := fmt.Sprintf(`  尊敬的用户，您好您正在修改您的用户邮箱，请在验证页面输入下方验证码完成后续操作：<b>%s</b> <br>
	
	如当前修改不是来自本人，请尽快登录链佰BAAS平台修改您的密码。<br>

	联系方式:%s`, randomString, "021-12345678")

	sended := utils.SentMail(email, sendHtml)
	if sended {
		// 删除邮件验证码
		err = db.DeleteAllTelCode(email)
		if err != nil {
			utils.ResponseJson(400, err.Error(), "", writer)
			return
		}
		// 保存邮件验证码
		telObj := module.TelCode{}
		telObj.ID = utils.GetUuid()
		telObj.TelNo = email
		telObj.VerifyCode = randomString
		err = db.SaveTelCodeInfo(telObj, telObj.ID, "")
		if err != nil {
			utils.ResponseJson(400, err.Error(), "", writer)
			return
		}
		listBytes, _ := json.Marshal(telObj)
		utils.ResponseJson(200, fmt.Sprintf("验证码已经发送到%s邮箱，请查收", email), string(listBytes), writer)
		return
	} else {
		utils.ResponseJson(400, "请填充正确的邮箱地址", "", writer)
		return
	}

}

// 问题反馈
func getfeedbackHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	list, err := db.GetFeedListInfo()
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得问题反馈错误："+err.Error(), "", writer)
		return
	}
	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得问题反馈错误："+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(listBytes), writer)
	return
}

// 问题反馈
func feedbackHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	feedback := module.Feedback{}
	err = json.Unmarshal(body, &feedback)

	if len(strings.TrimSpace(feedback.ID)) > 0 {

	} else {
		feedback.ID = utils.GetUuid()
	}
	err = db.SaveFeedInfo(feedback, feedback.ID, "")
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "反馈信息保存错误："+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "反馈信息保存完成", "", writer)
	return
}

// 设置通知公告信息
func announceHandler(writer http.ResponseWriter, request *http.Request) {

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	ann := module.Announce{}
	err = json.Unmarshal(body, &ann)

	if len(strings.TrimSpace(ann.ID)) > 0 {

	} else {
		ann.ID = utils.GetUuid()
	}
	err = db.SaveAnnInfo(ann, ann.ID, "")
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "保存通知公告错误："+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "公告信息已发布", "", writer)
	return
}

// 获得通知公告信息
func getannounceHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	list, err := db.GetAnnListInfo()
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得announcelist"+err.Error(), "", writer)
		return
	}
	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得annoucelist"+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(listBytes), writer)
	return
}

// 读取用户信息
func getuserHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	_, curusername, err := utils.IsManager(request)
	if err != nil {
		utils.ResponseJson(400, "获得当前人错误："+err.Error(), "", writer)
		return
	}

	curUser, err := db.GetUser(curusername)
	if err != nil {
		utils.ResponseJson(400, "获得用户错误："+err.Error(), "", writer)
		return
	}
	userBytes, err := json.Marshal(curUser)
	if err != nil {
		utils.ResponseJson(400, "解析用户错误："+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(userBytes), writer)
	return
}

// 读取用户信息列表
func userlistHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	isManager, _, err := utils.IsManager(request)
	if err != nil {
		utils.ResponseJson(400, "获得当前人错误："+err.Error(), "", writer)
		return
	}
	if isManager == false {
		utils.ResponseJson(400, "当前用户不是管理员", "", writer)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	searchObj := make(map[string]interface{})
	err = json.Unmarshal(body, &searchObj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}

	count := 15
	page := 1
	// 获得page 和 count
	if searchObj["page"] != nil {
		switch searchObj["page"].(type) {
		case string:
			page, err = strconv.Atoi(searchObj["page"].(string))
			if err != nil {
				page = 1
			}
		case float64:
			page = int(searchObj["page"].(float64))
		}

	}

	if searchObj["count"] != nil {
		switch searchObj["count"].(type) {
		case string:
			count, err = strconv.Atoi(searchObj["count"].(string))
			if err != nil {
				count = 15
			}
		case float64:
			count = int(searchObj["count"].(float64))
		}

	}

	if searchObj["search"] != nil {
		searchString := searchObj["search"].(string)

		searchlist, err := db.SearchUser(searchString)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "查找userlist"+err.Error(), "", writer)
			return
		}
		sList := []module.User{}
		if searchlist != nil {
			sList = searchlist.([]module.User)
			if len(sList) <= page*count {
				if len(sList) > 0 {
					sList = sList[(page-1)*count:]
				}

			} else {
				sList = sList[(page-1)*count : page*count]
			}
		}

		// get user project number
		listProject := []module.RtUser{}
		for _, userObj := range sList {
			mT := module.RtUser{}
			mT.User = userObj

			projectlist, err := db.GetReadSetupListByUser(userObj.UserName)
			if err != nil {
				mT.ProjectNum = 0
			}
			if projectlist == nil {
				mT.ProjectNum = 0
			}

			mT.ProjectNum = len(projectlist)

			listProject = append(listProject, mT)
		}

		resMap := make(map[string]interface{})
		resMap["count"] = len(searchlist.([]module.User))
		resMap["list"] = listProject
		searchlistBytes, err := json.Marshal(resMap)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "查找userlist"+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "", string(searchlistBytes), writer)
		return
	}

	list, err := db.GetUserList()
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得userlist"+err.Error(), "", writer)
		return
	}
	sList := []module.User{}
	if list != nil {
		sList = list.([]module.User)
		if len(sList) <= page*count {
			if len(sList) > 0 {
				sList = sList[(page-1)*count:]
			}

		} else {
			sList = sList[(page-1)*count : page*count]
		}
	}

	// get user project number
	listProject := []module.RtUser{}
	for _, userObj := range sList {
		mT := module.RtUser{}
		mT.User = userObj

		projectlist, err := db.GetReadSetupListByUser(userObj.UserName)
		if err != nil {
			mT.ProjectNum = 0
		}
		if projectlist == nil {
			mT.ProjectNum = 0
		}
		mT.ProjectNum = len(projectlist)

		listProject = append(listProject, mT)
	}

	resMap := make(map[string]interface{})
	resMap["count"] = len(list.([]module.User))
	resMap["list"] = listProject
	listBytes, err := json.Marshal(resMap)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得userlist"+err.Error(), "", writer)
		return
	}

	utils.ResponseJson(200, "", string(listBytes), writer)
	return
}

// 读取项目信息列表
func projectlistHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	isManager, _, err := utils.IsManager(request)
	if err != nil {
		utils.ResponseJson(400, "获得当前人错误："+err.Error(), "", writer)
		return
	}
	if isManager == false {
		utils.ResponseJson(400, "当前用户不是管理员", "", writer)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	searchObj := make(map[string]interface{})
	err = json.Unmarshal(body, &searchObj)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}

	count := 15
	page := 1
	// 获得page 和 count
	if searchObj["page"] != nil {
		switch searchObj["page"].(type) {
		case string:
			page, err = strconv.Atoi(searchObj["page"].(string))
			if err != nil {
				page = 1
			}
		case float64:
			page = int(searchObj["page"].(float64))
		}
	}

	if searchObj["count"] != nil {
		switch searchObj["count"].(type) {
		case string:
			count, err = strconv.Atoi(searchObj["count"].(string))
			if err != nil {
				count = 15
			}
		case float64:
			count = int(searchObj["count"].(float64))
		}
	}

	if searchObj["search"] != nil {
		searchString := searchObj["search"].(string)

		searchlist, err := db.SearchProject(searchString)
		var sortableList module.DefineList
		sortableList = searchlist
		sort.Sort(sortableList)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "查找项目列表错误："+err.Error(), "", writer)
			return
		}
		// 分页

		if sortableList != nil {
			if len(sortableList) <= page*count {
				if len(sortableList) > 0 {
					sortableList = sortableList[(page-1)*count:]
				}
			} else {
				sortableList = sortableList[(page-1)*count : page*count]
			}
		}

		resMap := make(map[string]interface{})
		resMap["count"] = len(searchlist)
		resMap["list"] = sortableList
		searchlistBytes, err := json.Marshal(resMap)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseJson(400, "查找项目列表错误："+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "", string(searchlistBytes), writer)
		return
	}

	list, err := db.GetReadSetupList()
	var sortableList module.DefineList
	sortableList = list
	sort.Sort(sortableList)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得projectlist"+err.Error(), "", writer)
		return
	}
	// 分页
	// if len(sList) > 0 {
	if len(sortableList) <= page*count {
		if len(sortableList) > 0 {
			sortableList = sortableList[(page-1)*count:]
		}
	} else {
		sortableList = sortableList[(page-1)*count : page*count]
	}
	// }

	resMap := make(map[string]interface{})
	resMap["count"] = len(list)
	resMap["list"] = sortableList
	listBytes, err := json.Marshal(resMap)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得projectlist"+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(listBytes), writer)
	return
}

// 读取日志信息列表
func loglistHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	// 判断管理员
	isManager, _, err := utils.IsManager(request)
	if err != nil {
		utils.ResponseJson(400, "获得当前人错误："+err.Error(), "", writer)
		return
	}
	if isManager == false {
		utils.ResponseJson(400, "当前用户不是管理员", "", writer)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得传参："+err.Error(), "", writer)
		return
	}
	logObj := make(map[string]interface{})
	err = json.Unmarshal(body, &logObj)
	list, err := db.GetLogList()
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得loglist"+err.Error(), "", writer)
		return
	}

	count := 15
	page := 1
	// 获得page 和 count
	if logObj["page"] != nil {
		switch logObj["page"].(type) {
		case string:
			page, err = strconv.Atoi(logObj["page"].(string))
			if err != nil {
				page = 1
			}
		case float64:
			page = int(logObj["page"].(float64))
		}
	}

	if logObj["count"] != nil {
		switch logObj["count"].(type) {
		case string:
			count, err = strconv.Atoi(logObj["count"].(string))
			if err != nil {
				count = 15
			}
		case float64:
			count = int(logObj["count"].(float64))
		}
	}

	// 检索项目操作日志
	if logObj["search"] != nil {
		searchstring := logObj["search"].(string)
		fmt.Println(searchstring)

		list, err = db.SearchLogInfo(searchstring)
		if err != nil {
			utils.ResponseJson(400, "查找loglist"+err.Error(), "", writer)
			return
		}
	}
	// 分页
	sList := list.([]module.Log)
	numCount := len(sList)
	fmt.Println("numCount:", numCount)
	if len(sList) <= page*count {
		if len(sList) > 0 {
			sList = sList[(page-1)*count:]
		}
	} else {
		sList = sList[(page-1)*count : page*count]
	}

	resMap := make(map[string]interface{})
	resMap["count"] = numCount
	resMap["list"] = sList
	listBytes, err := json.Marshal(resMap)
	if err != nil {
		log.Println(err.Error())
		utils.ResponseJson(400, "获得loglist"+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(listBytes), writer)
	return
}
