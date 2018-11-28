package main

import (
	"encoding/json"
	"fmt"
	"frontend4chain/config"
	"frontend4chain/constant"
	"frontend4chain/db"
	"frontend4chain/handfuncation"
	"frontend4chain/module"
	"frontend4chain/utils"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		defer r.Body.Close()
		// explorer login
		if r.URL.Path == "/explorer/login" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var loginMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &loginMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			name := loginMap["username"].(string)
			password := loginMap["password"].(string)

			tempPro, err := db.GetSetupInfoByUserName(name)
			if err != nil {
				utils.ResponseJson(400, "根据用户名获得项目错误:"+err.Error(), "", w)
				return
			}
			if tempPro == nil {
				utils.ResponseJson(400, "用户名不存在,请确认！", "", w)
				return
			}
			if utils.Md5(password) != tempPro.(module.Define).ExploderPassword {
				utils.ResponseJson(400, "密码错误，请确认！", "", w)
				return
			}

			tstring, err := utils.CreateToken(name, "explorer")
			if err != nil {
				log.Println("login failed " + err.Error())
				utils.ResponseJson(400, "get CreateToken : "+err.Error(), "", w)
				return
			}

			rT := module.Token{}
			rT.Username = name
			rT.Token = tstring
			rBytes, err := json.Marshal(rT)
			if err != nil {
				utils.ResponseJson(400, "token Marshal: "+err.Error(), "", w)
				log.Println("token Marshal:" + err.Error())
				return
			}
			// w.Write(rBytes)
			utils.ResponseJson(200, "", string(rBytes), w)
			return
		}
		//login
		if r.URL.Path == "/login" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var loginMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &loginMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			name := ""
			if loginMap["username"] != nil {
				name = loginMap["username"].(string)
			}

			telno := ""
			if loginMap["telno"] != nil {
				telno = loginMap["telno"].(string)
			}

			email := ""
			if loginMap["email"] != nil {
				email = loginMap["email"].(string)
			}

			password := loginMap["password"].(string)
			// check username
			var userface interface{}
			if name != "" {
				userface, err = db.GetUser(name)
				if err != nil {
					utils.ResponseJson(400, "couchdb "+err.Error(), "", w)
					return
				}
			}

			if email != "" {
				userface, err = db.GetUserByEmail(email)
				if err != nil {
					utils.ResponseJson(400, "couchdb "+err.Error(), "", w)
					return
				}
			}

			if telno != "" {
				userface, err = db.GetUserByTel(telno)
				if err != nil {
					utils.ResponseJson(400, "couchdb "+err.Error(), "", w)
					return
				}
			}

			if userface == nil {
				utils.ResponseJson(400, "用户不存在！", "", w)
				return
			}
			currentUser := userface.(module.User)
			if currentUser.Active == false {
				utils.ResponseJson(400, "激活错误，请打开邮箱的激活地址！ ", "", w)
				return
			}
			// md5 后比较

			if currentUser.Password != utils.Md5(password) {
				utils.ResponseJson(400, "密码不对", "", w)
				return
			}

			tstring, err := utils.CreateToken(name, "project")
			if err != nil {
				log.Println("login failed " + err.Error())
				utils.ResponseJson(400, "get CreateToken : "+err.Error(), "", w)
				return
			}

			rT := module.Token{}
			rT.Username = name
			rT.Role = currentUser.Role
			rT.Token = tstring
			rBytes, err := json.Marshal(rT)
			if err != nil {
				utils.ResponseJson(400, "token Marshal: "+err.Error(), "", w)
				log.Println("token Marshal:" + err.Error())
				return
			}
			// w.Write(rBytes)
			utils.ResponseJson(200, "", string(rBytes), w)
			return
		}
		//register
		if r.URL.Path == "/register" {
			w.Header().Set("Content-Type", "application/json")
			active := true
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			log.Println(string(body))
			registMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(body), &registMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, "解析register body错误："+err.Error(), "", w)
				return
			}

			name := registMap["username"].(string)
			password := registMap["password"].(string)
			telno := ""
			if registMap["telno"] != nil {
				telno = registMap["telno"].(string)
			}
			email := ""
			if registMap["email"] != nil {
				email = registMap["email"].(string)
			}
			random := ""
			if registMap["random"] != nil {
				random = registMap["random"].(string)
			}

			verifycode := ""
			if registMap["checkcode"] != nil {
				verifycode = registMap["checkcode"].(string)
			}
			telverifycode := ""
			if registMap["verifycode"] != nil {
				telverifycode = registMap["verifycode"].(string)

			}

			//验证random code
			randObj, err := db.GetRandomInfo(random)
			if err != nil {
				utils.ResponseJson(400, "获得random错误: "+err.Error(), "", w)
				return
			}
			if strings.ToLower(randObj.(module.VerifyCode).VerifyCode) != strings.ToLower(verifycode) {
				utils.ResponseJson(400, "验证码错误", "", w)
				return
			}

			//验证手机tel code
			if len(strings.TrimSpace(telno)) >= 11 {
				telObj, err := db.GetTelCodeInfo(telno)
				if err != nil {
					utils.ResponseJson(400, "获得手机验证码错误: "+err.Error(), "", w)
					return
				}
				if telObj == nil {
					utils.ResponseJson(400, "手机验证码有误", "", w)
					return
				}
				if strings.ToLower(telObj.(module.TelCode).VerifyCode) != strings.ToLower(telverifycode) {
					utils.ResponseJson(400, "手机验证码错误", "", w)
					return
				}
			}

			exitUser, err := db.GetUser(name)
			if err != nil {
				utils.ResponseJson(400, "获得用户信息错误："+err.Error(), "", w)
				return
			}

			if exitUser == nil {
				// 保存用户信息
				newUser := module.User{}
				newUser.UserName = name
				newUser.Password = utils.Md5(password)
				newUser.Role = "user"
				newUser.TelNo = telno
				newUser.Email = email
				newUser.ID = utils.GetUuid()
				newUser.Active = active
				err = db.SaveUser(newUser, newUser.ID, "")
				if err != nil {
					utils.ResponseJson(400, "保存用户信息： "+err.Error(), "", w)
					return
				}
			} else {
				eUser := module.User{}
				rev, _ := db.GetReadUserInfo(exitUser.(module.User).ID, &eUser)
				eUser.Email = email
				eUser.Password = utils.Md5(password)
				eUser.Role = "user"
				eUser.TelNo = telno
				eUser.Active = active
				err = db.SaveUser(eUser, eUser.ID, rev)
				if err != nil {
					utils.ResponseJson(400, "保存用户信息： "+err.Error(), "", w)
					return
				}
			}

			// 发送邮箱验证码
			if len(strings.TrimSpace(email)) > 0 {
				active = false
				emailString, err := utils.CreateEmailToken(name)
				if err != nil {
					utils.ResponseJson(400, "创建邮箱激活码错误: "+err.Error(), "", w)
					return
				}
				activeUrl := fmt.Sprint("http://192.168.0.237/account/active?s=", emailString)
				activeCode := fmt.Sprintf(constant.ActiveCode, name, emailString, activeUrl)
				if sended := utils.SentMail(strings.TrimSpace(email), activeCode); sended == true {
					log.Println("已经发送邮箱到：", email)
					utils.ResponseJson(200, "激活地址已经发送到："+strings.TrimSpace(email), "", w)
					return
				} else {
					log.Println("已经发送邮箱到：", email, "失败！")
					utils.ResponseJson(400, "邮箱地址错误，请重新输入！", "", w)
					return
				}
			}
			// 如果不是邮箱注册的，返回token

			tstring, err := utils.CreateToken(name, "project")
			if err != nil {
				log.Println("login failed " + err.Error())
				utils.ResponseJson(400, "创建token错误 : "+err.Error(), "", w)
				return
			}

			rT := module.Token{}
			rT.Username = name
			rT.Token = tstring
			rBytes, err := json.Marshal(rT)

			if err != nil {
				utils.ResponseJson(400, "返回json marshal 错误: "+err.Error(), "", w)
				log.Println("token Marshal:" + err.Error())
				return
			}

			utils.ResponseJson(200, "注册完成", string(rBytes), w)
			return
		}
		//checkCode
		if r.URL.Path == "/getcheckcode" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var randomMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &randomMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, "解析post body错误："+err.Error(), "", w)
				return
			}
			random := randomMap["random"].(string)
			code := utils.RandomString(4)
			log.Println(random, code)
			randomObj := module.VerifyCode{}
			randomObj.Random = random
			randomObj.VerifyCode = code
			err = db.SaveRandomInfo(randomObj, utils.GetUuid(), "")
			if err != nil {
				utils.ResponseJson(400, "保存random错误："+err.Error(), "", w)
				return
			}
			listBytes, err := json.Marshal(randomObj)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, "random Marshal 错误："+err.Error(), "", w)
				return
			}
			utils.ResponseJson(200, "", string(listBytes), w)
			return
		}
		// checkuser
		if r.URL.Path == "/checkuser" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var userMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &userMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			username := userMap["username"].(string)
			user, err := db.GetUser(username)
			if err != nil {
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}

			if user == nil {
				// 返回
				utils.ResponseJson(200, "用户名可用", "", w)
				return
			} else {
				// 返回
				utils.ResponseJson(400, "用户名已经被注册，请重新填写用户名", "", w)
				return
			}

		}
		// checktelno
		if r.URL.Path == "/checktelno" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var userMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &userMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			telno := userMap["telno"].(string)
			user, err := db.GetUserByTel(telno)
			if err != nil {
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}

			if user == nil || (user != nil && user.(module.User).Active == false) {
				// 返回
				utils.ResponseJson(200, "手机号码可用", "", w)
				return
			} else {
				// 返回
				utils.ResponseJson(400, "手机号码已经被注册，请重新填写手机号码", "", w)
				return
			}

		}
		// checkemail
		if r.URL.Path == "/checkemail" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var userMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &userMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			email := userMap["email"].(string)
			user, err := db.GetUserByEmail(email)
			if err != nil {
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}

			if user == nil || (user != nil && user.(module.User).Active == false) {
				// 返回
				utils.ResponseJson(200, "邮箱地址可用", "", w)
				return
			} else {
				utils.ResponseJson(400, "邮箱地址已经被注册，请重新填写邮箱地址", "", w)
				return
			}

		}
		// Getverifycode
		if r.URL.Path == "/getverifycode" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var telMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &telMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}

			telno := telMap["telno"].(string)
			// 删除手机号码验证码
			err = db.DeleteAllTelCode(telno)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			// 保存手机号码验证码
			telObj := module.TelCode{}
			telObj.ID = utils.GetUuid()
			telObj.TelNo = telno
			telObj.VerifyCode = "1234"
			err = db.SaveTelCodeInfo(telObj, telObj.ID, "")
			if err != nil {
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			listBytes, _ := json.Marshal(telObj)
			utils.ResponseJson(200, "", string(listBytes), w)
			return
		}
		// modifypassword
		if r.URL.Path == "/activepassword" {

			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var activeMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &activeMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			token := ""
			username := ""
			if activeMap["token"] != nil {
				token = activeMap["token"].(string)
			}
			password := ""
			if activeMap["password"] != nil {
				password = activeMap["password"].(string)
			}

			verifyToken, err := utils.VerifyToken(token)
			if err != nil {
				log.Println("Authentication failed " + err.Error())
				utils.ResponseJson(400, "token expaired: "+err.Error(), "", w)
				return
			}
			if claims, ok := verifyToken.Claims.(jwt.MapClaims); ok && verifyToken.Valid {
				username = claims["iss"].(string)
			} else {
				utils.ResponseJson(400, "修改密码超期，请重新忘记密码操作: "+err.Error(), "", w)
				return
			}
			tUser, _ := db.GetUser(username)
			activeUser := module.User{}

			rev, _ := db.GetReadUserInfo(tUser.(module.User).ID, &activeUser)
			activeUser.Password = utils.Md5(password)
			err = db.SaveUser(activeUser, activeUser.ID, rev)
			if err != nil {
				utils.ResponseJson(400, "密码修改错误: "+err.Error(), "", w)
				return
			}
			// 返回
			utils.ResponseJson(200, "密码修改完成，请重新登录", "", w)
			return
		}
		// forgotpassword
		if r.URL.Path == "/forgotpassword" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var forgotMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &forgotMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			telno := ""
			if forgotMap["telno"] != nil {
				telno = forgotMap["telno"].(string)
			}
			verifycode := ""
			if forgotMap["verifycode"] != nil {
				verifycode = forgotMap["verifycode"].(string)
			}
			newpassword := ""
			if forgotMap["newpassword"] != nil {
				newpassword = forgotMap["newpassword"].(string)
			}
			email := ""
			if forgotMap["email"] != nil {
				email = forgotMap["email"].(string)
			}

			if len(strings.TrimSpace(telno)) >= 11 {

				user, err := db.GetUserByTel(telno)
				if err != nil {
					utils.ResponseJson(400, "获得用户信息错误: "+err.Error(), "", w)
					return
				}
				log.Println("user=======================:")
				log.Println(user)

				if user == nil {
					utils.ResponseJson(400, "手机号码填写错误，请输入注册时填写的手机号码", "", w)
					return
				}
				telObj, err := db.GetTelCodeInfo(telno)
				if err != nil {
					utils.ResponseJson(400, "获得手机验证码错误: "+err.Error(), "", w)
					return
				}
				if telObj == nil {
					utils.ResponseJson(400, "手机验证码有误", "", w)
					return
				}
				if strings.ToLower(telObj.(module.TelCode).VerifyCode) != strings.ToLower(verifycode) {
					utils.ResponseJson(400, "手机验证码错误", "", w)
					return
				}
				if user != nil {
					userobj := user.(module.User)
					mUser := module.User{}
					rev, _ := db.GetReadUserInfo(userobj.ID, &mUser)
					mUser.Password = utils.Md5(newpassword)
					err = db.SaveUser(mUser, mUser.ID, rev)
					if err != nil {
						utils.ResponseJson(400, "保存密码错误: "+err.Error(), "", w)
						return
					}
				}
			}

			if len(strings.TrimSpace(email)) > 0 {
				user, err := db.GetUserByEmail(email)
				if err != nil {
					utils.ResponseJson(400, "获得用户信息错误: "+err.Error(), "", w)
					return
				}
				if user != nil {
					userobj := user.(module.User)
					emailString, err := utils.CreateEmailToken(userobj.UserName)
					if err != nil {
						utils.ResponseJson(400, "创建邮箱激活码错误: "+err.Error(), "", w)
						return
					}
					forgotUrl := fmt.Sprint("http://192.168.0.237/account/forgot?s=", emailString)
					forgothtml := fmt.Sprintf(constant.ForgotCode, userobj.UserName, emailString, forgotUrl)
					if sended := utils.SentMail(strings.TrimSpace(email), forgothtml); sended == true {
						log.Println("已经发送邮箱到：", email)
						utils.ResponseJson(200, "请打开邮箱："+strings.TrimSpace(email)+",打开链接修改密码！", "", w)
						return
					} else {
						log.Println("已经发送邮箱到：", email, "失败！")
						utils.ResponseJson(400, "邮箱地址错误，请重新输入！", "", w)
						return
					}
				} else {
					utils.ResponseJson(400, "该邮箱未注册，请确认", "", w)
					return
				}
			}
			utils.ResponseJson(200, "密码已修改，请重新登录", "", w)
			return
		}

		// 激活用户链接
		if r.URL.Path == "/activeuser" {
			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var activeMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &activeMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			token := ""
			username := ""
			if activeMap["token"] != nil {
				token = activeMap["token"].(string)
			}
			verifyToken, err := utils.VerifyToken(token)
			if err != nil {
				log.Println("Authentication failed " + err.Error())
				utils.ResponseJson(400, "token expaired: "+err.Error(), "", w)
				return
			}
			if claims, ok := verifyToken.Claims.(jwt.MapClaims); ok && verifyToken.Valid {
				username = claims["iss"].(string)
			} else {
				utils.ResponseJson(400, "激活用户超期，请重新注册: "+err.Error(), "", w)
				return
			}
			tUser, _ := db.GetUser(username)
			activeUser := module.User{}

			rev, _ := db.GetReadUserInfo(tUser.(module.User).ID, &activeUser)
			activeUser.Active = true
			err = db.SaveUser(activeUser, activeUser.ID, rev)
			if err != nil {
				utils.ResponseJson(400, "激活用户错误: "+err.Error(), "", w)
				return
			}
			// 返回
			utils.ResponseJson(200, "激活用户完成，请登录", "", w)
			return
		}

		if r.Header.Get("Authorization") != "" {
			token, err := utils.VerifyToken(r.Header.Get("Authorization"))
			if err != nil {
				log.Println("Authentication failed " + err.Error())
				utils.ResponseJson(401, "token expaired: "+err.Error(), "", w)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				log.Println("username: " + claims["iss"].(string))
				r.Header.Set("username", claims["iss"].(string))
				r.Header.Set("subject", claims["sub"].(string))
			} else {
				log.Println("Authentication failed " + err.Error())
				utils.ResponseJson(401, "token expaired: "+err.Error(), "", w)
				return
			}
		}
		// modifypassword
		if r.URL.Path == "/modifypassword" {
			// 必须登陆
			if r.Header.Get("Authorization") == "" {
				utils.ResponseJson(401, "请登录后再修改密码", "", w)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}
			var registMap map[string]interface{}
			err = json.Unmarshal([]byte(body), &registMap)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseJson(400, err.Error(), "", w)
				return
			}

			name := registMap["username"].(string)
			oldpassword := ""
			if registMap["oldpassword"] != nil {
				oldpassword = registMap["oldpassword"].(string)
			}
			password := ""
			if registMap["password"] != nil {
				password = registMap["password"].(string)
			}
			random := ""
			if registMap["random"] != nil {
				random = registMap["random"].(string)
			}
			verifycode := ""
			if registMap["verifycode"] != nil {
				verifycode = registMap["verifycode"].(string)
			}

			//验证random code
			randObj, err := db.GetRandomInfo(random)
			if err != nil {
				utils.ResponseJson(400, "random: "+err.Error(), "", w)
				return
			}
			if strings.ToLower(randObj.(module.VerifyCode).VerifyCode) != strings.ToLower(verifycode) {
				utils.ResponseJson(400, "验证码错误", "", w)
				return
			}
			user, err := db.GetUser(name)
			if err != nil {
				utils.ResponseJson(400, "获得用户信息错误: "+err.Error(), "", w)
				return
			}
			// 修改密码
			if user.(module.User).Password == utils.Md5(oldpassword) {
				userobj := user.(module.User)
				mUser := module.User{}
				rev, _ := db.GetReadUserInfo(userobj.ID, &mUser)
				mUser.Password = utils.Md5(password) //修改密码
				err = db.SaveUser(mUser, mUser.ID, rev)
				if err != nil {
					utils.ResponseJson(400, "修改密码错误: "+err.Error(), "", w)
					return
				}
			} else {
				utils.ResponseJson(400, "老密码输入错误，请重新输入！", "", w)
				return
			}
			// 返回
			utils.ResponseJson(200, "密码修改完成，请重新登录", "", w)
			return
		}

		next.ServeHTTP(w, r)
		log.Printf("Comleted %s in %v", r.URL.Path, time.Since(start))
	})
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	handfuncation.HandlerAll(w, r)

}

func main() {
	config.InitConf(os.Args)
	http.Handle("/", loggingHandler(http.HandlerFunc(handlerIndex)))
	log.Println(config.All().Listen.Ip + ":" + config.All().Listen.Port)
	if err := http.ListenAndServe(config.All().Listen.Ip+":"+config.All().Listen.Port, nil); err != nil {
		log.Fatalln(err)
	}
}
