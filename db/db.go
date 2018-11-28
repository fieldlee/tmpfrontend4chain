package db

import (
	"encoding/json"
	"fmt"
	"frontend4chain/constant"
	"frontend4chain/module"
	"time"

	"frontend4chain/config"

	"github.com/rhinoman/couchdb-go"
)

var timeout = time.Duration(500 * time.Millisecond)

func CreateDb(dbname string) error {
	conn, err := couchdb.NewConnection(config.All().Couchdb.Ip, config.All().Couchdb.Port, timeout)
	if err != nil {
		return err
	}
	auth := couchdb.BasicAuth{Username: config.All().Couchdb.Username, Password: config.All().Couchdb.Password}
	return conn.CreateDB(dbname, &auth)
}

func CheckDb(dbname string) (bool, error) {
	conn, err := couchdb.NewConnection(config.All().Couchdb.Ip, config.All().Couchdb.Port, timeout)
	if err != nil {
		return false, err
	}
	// auth := couchdb.BasicAuth{Username: config.All().Couchdb.Username, Password: config.All().Couchdb.Password}
	dbList, err := conn.GetDBList()
	if err != nil {
		return false, err
	}

	for i := range dbList {
		if dbname == dbList[i] {
			return true, nil
		}
	}
	return false, nil
}

func GetDb(dbname string) (*couchdb.Database, error) {
	conn, err := couchdb.NewConnection(config.All().Couchdb.Ip, config.All().Couchdb.Port, timeout)
	if err != nil {
		return nil, err
	}
	auth := couchdb.BasicAuth{Username: config.All().Couchdb.Username, Password: config.All().Couchdb.Password}

	if exist, err := CheckDb(dbname); err != nil || exist == false {
		err = CreateDb(dbname)
		if err != nil {
			return nil, err
		}
	}
	return conn.SelectDB(dbname, &auth), nil
}

func SearchUser(search string) (interface{}, error) {
	selector := fmt.Sprintf(`{
		"$or": [
         {
            "username": {
               "$regex": "%s"
            }
         },
         {
            "email": {
               "$regex": "%s"
            }
		 },
         {
            "telno": {
               "$regex": "%s"
            }
         }
      ]
	}`, search, search, search)
	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.UserFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 根据用户名获得用户信息
func GetUser(name string) (interface{}, error) {
	selector := `{"username": "` + name + `"}`

	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.UserFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) >= 1 {
		return findResult.Docs[0], nil
	}

	return nil, nil
}

// 获得全部用户信息
func GetUserList() (interface{}, error) {
	selector := `{"_id":{"$gt":null}}`
	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.UserFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 获得用户rev
func GetReadUserInfo(id string, user *module.User) (string, error) {
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return "", err
	}
	rev, err := db.Read(id, user, nil)
	if err != nil {
		return "", err
	}
	return rev, nil
}

// 根据电话号码获得用户信息
func GetUserByTel(telno string) (interface{}, error) {
	selector := `{"telno": "` + telno + `"}`

	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.UserFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) >= 1 {
		return findResult.Docs[0], nil
	}

	return nil, nil
}

// 根据email获得用户信息
func GetUserByEmail(email string) (interface{}, error) {
	selector := `{"email": "` + email + `"}`

	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.UserFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}
	db, err := GetDb(constant.DB_USER)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) >= 1 {
		return findResult.Docs[0], nil
	}

	return nil, nil
}

// 保存用户信息
func SaveUser(user interface{}, id string, rev string) error {
	userDb, err := GetDb(constant.DB_USER)
	if err != nil {
		return err
	}
	_, err = userDb.Save(user, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 保存fabric 区块链配置信息
func SaveSetupInfo(doc interface{}, id string, rev string) error {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return err
	}
	_, err = setupDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 根据项目id 删除项目
func DeleteProject(id string, rev string) (string, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return "", err
	}

	del, err := setupDb.Delete(id, rev)
	if err != nil {
		return "", err
	}
	return del, nil
}

// 查找项目信息
func SearchProject(search string) ([]module.Define, error) {
	selector := fmt.Sprintf(`{
		"$or": [
         {
            "projectName": {
               "$regex": "%s"
            }
         },
         {
            "manager": {
               "$regex": "%s"
            }
		 },
         {
            "domain": {
               "$regex": "%s"
            }
         }
      ]
	}`, search, search, search)
	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}
	db, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	err = db.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 获得fabric 项目根据id获得的信息
func GetSetupInfo(id string) (interface{}, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	selector := `{"id": "` + id + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}
	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}
	err = setupDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) >= 1 {
		return findResult.Docs[0], nil
	}
	return nil, nil
}

// 获得fabric 项目根据id获得的信息
func GetSetupInfoByUserName(username string) (interface{}, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	selector := `{"exploderUser": "` + username + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}
	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}
	err = setupDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) >= 1 {
		return findResult.Docs[0], nil
	}
	return nil, nil
}

// 获得项目配置信息
func GetReadSetupInfo(id string, setupInfo *module.Define) (string, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return "", err
	}
	rev, err := setupDb.Read(id, setupInfo, nil)
	if err != nil {
		return "", err
	}
	return rev, nil
}

// 获得项目配置列表
func GetReadSetupList() ([]module.Define, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	selector := `{"_id":{"$gt":null}}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = setupDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}

	return findResult.Docs, nil
}

// 根据当前用户获得项目列表
func GetReadSetupListByUser(username string) ([]module.Define, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	selector := `{"manager":"` + username + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = setupDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}

	return findResult.Docs, nil
}

// 根据状态获得项目列表
func GetReadSetupListByStatus(status string) ([]module.Define, error) {
	setupDb, err := GetDb(constant.DB_SETUP)
	if err != nil {
		return nil, err
	}
	selector := `{"status":"` + status + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.SetupFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = setupDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}

	return findResult.Docs, nil
}

// 根据id获得sdk 信息
func GetReadSdkInfo(id string, sdkInfo *module.Sdk) (string, error) {
	sdkDb, err := GetDb(constant.DB_SDK)
	if err != nil {
		return "", err
	}
	rev, err := sdkDb.Read(id, sdkInfo, nil)
	if err != nil {
		return "", err
	}
	return rev, nil
}

// 保存sdk信息
func SaveSdkInfo(doc interface{}, id string, rev string) error {
	sdkDb, err := GetDb(constant.DB_SDK)
	if err != nil {
		return err
	}
	_, err = sdkDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 保存日志信息
func SaveLogInfo(doc interface{}, id string, rev string) error {
	logDb, err := GetDb(constant.DB_Log)
	if err != nil {
		return err
	}
	_, err = logDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 查找日志
func SearchLogInfo(search string) (interface{}, error) {
	logDb, err := GetDb(constant.DB_Log)
	if err != nil {
		return nil, err
	}
	selector := fmt.Sprintf(`{
		"$or": [
         {
            "username": {
               "$regex": "%s"
            }
         },
         {
            "data": {
               "$regex": "%s"
            }
		 }
      ]
	}`, search, search)
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.LogFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = logDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 获得日志列表信息
func GetLogList() (interface{}, error) {
	logDb, err := GetDb(constant.DB_Log)
	if err != nil {
		return nil, err
	}
	selector := `{"_id":{"$gt":null}}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.LogFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 10000}

	err = logDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 获得日志列表信息
func GetLogListByUserName(username string) (interface{}, error) {
	logDb, err := GetDb(constant.DB_Log)
	if err != nil {
		return nil, err
	}
	selector := `{"username":"` + username + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.LogFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 10000}

	err = logDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 获得日志列表信息
func GetLogListByID(id string) (interface{}, error) {
	logDb, err := GetDb(constant.DB_Log)
	if err != nil {
		return nil, err
	}
	selector := `{"id":"` + id + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.LogFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = logDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 保存随机数信息
func SaveRandomInfo(doc interface{}, id string, rev string) error {
	randomDb, err := GetDb(constant.DB_RANDOM)
	if err != nil {
		return err
	}
	_, err = randomDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 获得随机数信息
func GetRandomInfo(random string) (interface{}, error) {
	randomDb, err := GetDb(constant.DB_RANDOM)
	if err != nil {
		return nil, err
	}
	selector := `{"random": "` + random + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.RandomFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}

	err = randomDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs[0], nil
}

// 获得用户rev
func GetTelInfo(id string, telcode *module.TelCode) (string, error) {
	telcodeDb, err := GetDb(constant.DB_TELCODE)
	if err != nil {
		return "", err
	}
	rev, err := telcodeDb.Read(id, telcode, nil)
	if err != nil {
		return "", err
	}
	return rev, nil
}

func DeleteAllTelCode(telno string) error {
	telcodeDb, err := GetDb(constant.DB_TELCODE)
	if err != nil {
		return err
	}
	selector := `{"telno": "` + telno + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return err
	}
	//Get the results from find.
	findResult := module.TelFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}

	err = telcodeDb.Find(&findResult, &params)
	if err != nil {
		return err
	}
	if len(findResult.Docs) > 0 {
		for _, doc := range findResult.Docs {
			rev, err := telcodeDb.Read(doc.ID, &doc, nil)
			if err != nil {
				return err
			}
			err = DeleteTelCodeInfo(doc.ID, rev)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 删除邮箱
func DeleteTelCodeInfo(id string, rev string) error {
	telcodeDb, err := GetDb(constant.DB_TELCODE)
	if err != nil {
		return err
	}
	_, err = telcodeDb.Delete(id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 保存电话号码信息
func SaveTelCodeInfo(doc interface{}, id string, rev string) error {
	telcodeDb, err := GetDb(constant.DB_TELCODE)
	if err != nil {
		return err
	}
	_, err = telcodeDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 获得电话号码code 信息
func GetTelCodeInfo(telno string) (interface{}, error) {
	telcodeDb, err := GetDb(constant.DB_TELCODE)
	if err != nil {
		return nil, err
	}
	selector := `{"telno": "` + telno + `"}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.TelFindResponse{}

	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1}

	err = telcodeDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	if len(findResult.Docs) > 0 {
		return findResult.Docs[0], nil
	}
	return nil, nil
}

// 保存反馈信息
func SaveFeedInfo(doc interface{}, id string, rev string) error {
	feedDb, err := GetDb(constant.DB_FEED)
	if err != nil {
		return err
	}
	_, err = feedDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 保存反馈信息
func GetFeedListInfo() (interface{}, error) {
	feedDb, err := GetDb(constant.DB_FEED)
	if err != nil {
		return nil, err
	}

	selector := `{"_id":{"$gt":null}}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.FeedFindResponse{}
	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}
	err = feedDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}

// 保存通知公告信息
func SaveAnnInfo(doc interface{}, id string, rev string) error {
	annDb, err := GetDb(constant.DB_ANN)
	if err != nil {
		return err
	}
	_, err = annDb.Save(doc, id, rev)
	if err != nil {
		return err
	}
	return nil
}

// 获得通知公告信息
func GetAnnListInfo() (interface{}, error) {
	annDb, err := GetDb(constant.DB_ANN)
	if err != nil {
		return nil, err
	}
	t := string(time.Now().Unix())
	selector := `{"start":{ "$lt":` + t + `},"end":{"$gt":` + t + `}}`
	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	//Get the results from find.
	findResult := module.AnnFindResponse{}
	params := couchdb.FindQueryParams{Selector: &selectorObj, Limit: 1000}
	err = annDb.Find(&findResult, &params)
	if err != nil {
		return nil, err
	}
	return findResult.Docs, nil
}
