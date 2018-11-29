package constant

const (
	//下划线
	PUBLICKEY = "KDKKDKKSDFJASLFJA;LKSDFJADSFJOQWERUQWEOREWQJFDSKLJFASDJF"
	// Couchdb username and password
	COUCHUSERNAME = "couchadmin"
	COUCHPASSWORD = "adminpwd"
	// 本机ip
	CURIP      = "192.168.0.15"
	CURURLROOT = "http://192.168.0.15:8000"
	// SDK
	SDKGITURL = "https://github.com/fieldlee/fabricSdk.git"

	// 状态列表
	SAVEED   = "saved"
	RESAVEED = "resaved"
	DEPLOYED = "deployed"
	CHANED   = "channeled"
	CCED     = "chaincodeed"

	// DBNAME
	DBNAME = "BAAS"
	// TABLENAME
	USERTABLE = "USER"

	// 证书根目录
	ROOTPATH   = "/var/certification"
	CONFIGPATH = "/var/config"
	YAMLPATH   = "/var/yaml"
	SDKPATH    = "/var/sdk"
	//sub service URL
	UPLOADURL   = "/file/upload"     //上传证书接口
	UPLOADTXURL = "/file/uploadById" //上传文件到指定证书目录
	//SDK
	SDKINSTALLURL = "/sdk/install"   //安装fabric sdk
	SHELLEXECURL  = "/shell/execute" //调用sdk接口
	// CHAINCODE
	CCUPLOADZIPURL = "/chaincode/uploadByZip"
	CCUPLOADGITURL = "/chaincode/uploadByGit"
	// Environtment
	SYSTEMINFOURL = "/system/info"
	CHECKIP       = "/environment/ip"
	CHECKURL      = "/environment/check"
	INSTALLURL    = "/environment/install"

	LUNCHURL  = "/program/start"
	STATUSURL = "/program/status"
	// system
	SYSTEMPERFORMURL = "/system/performance"
	DOCKERPERFORMURL = "/docker/performance"
	DOCKEREXEC       = "/docker/execute"
	SERVERPORT       = "1081"

	// SHELL 文件名称
	SHELLRESTART = "restart.sh"
	// SHELLGETIMAGE = "get-images.sh"
	// SHELLSETIMAGE = "set-images-latest.sh"

	// install docker
	INSTALLDOCKER        = "yum install docker -y;systemctl daemon-reload;systemctl restart docker.service;"
	INSTALLNODEJS        = "curl --silent --location https://rpm.nodesource.com/setup_9.x | bash -;yum install -y nodejs;yum install -y gcc-c++;npm install -g pm2;"
	INSTALLDOCKERCOMPOSE = "yum install -y epel-release;yum install -y python-pip;pip install docker-compose;"
	INSTALLGIT           = "yum install -y git;"
	INSTALLJQ            = "yum install -y jq;"

	// check envirement
	CHECKDOCKER        = "docker -v;"
	CHECKDOCKERCOMPOSE = "docker-compose -v;"
	CHECKNODE          = "node -v;"
	CHECKGIT           = "git -v;"
	CHECKJQ            = "jq -v;"

	// check port
	CHECKPORT      = "firewall-cmd --query-port=%s/tcp;"
	ADDPORT        = "firewall-cmd --permanent --add-port=%s/tcp;"
	REMOVEPORT     = "firewall-cmd --permanent --remove-port=%s/tcp;"
	RESTATFIREWALL = "firewall-cmd --reload;"
	STOPFIREWALL   = "systemctl stop firewalld;"

	// ubuntun install 命令
	U_INSTALLDOCKER        = "sudo apt-get install docker.io -y; sudo systemctl daemon-reload;sudo systemctl restart docker;"
	U_INSTALLNODEJS        = "sudo sh node.sh;sudo apt-get install -y nodejs;sudo apt-get install build-essential;sudo npm install -g pm2;"
	U_INSTALLDOCKERCOMPOSE = "sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose;sudo chmod +x /usr/local/bin/docker-compose"
	U_INSTALLGIT           = "sudo apt-get install -y git;"
	U_INSTALLJQ            = "sudo apt-get install -y jq;"

	// check envirement
	U_CHECKDOCKER        = "docker -v;"
	U_CHECKDOCKERCOMPOSE = "docker-compose -v;"
	U_CHECKNODE          = "node -v;"
	U_CHECKGIT           = "git -v;"
	U_CHECKJQ            = "jq -v;"

	// check port
	U_CHECKPORT      = "sudo ufw status verbose | grep %s"
	U_ADDPORT        = "sudo ufw allow %s/tcp"
	U_REMOVEPORT     = "sudo ufw delete allow %s/tcp"
	U_RESTATFIREWALL = "sudo ufw reload"
	U_STOPFIREWALL   = "sudo ufw disable"
	// dbname
	DB_USER    = "user"
	DB_SETUP   = "setup"
	DB_SDK     = "sdk"
	DB_RANDOM  = "random"
	DB_TELCODE = "telcode"
	DB_ANN     = "ann"
	DB_FEED    = "feed"
	DB_Log     = "log"

	ActiveCode = `<div style="padding: 6px;font-size: 16px;">
	<h1 style="margin: 0 0 7px;font-size: 19px;
  color: #333;
  line-height: 28px;">%s，您好！</h1>
	<p style="margin: 7px 0 18px;font-size: 16px;
  color: #333;line-height: 28px;">欢迎注册链佰BAAS平台！为了激活您的账户，请点击下面的按钮验证您的邮箱</p>
	<a target="_blank"
	  style="display:block;width: 140px;height:36px;font-size:16px;line-height:36px;text-align:center;color:#fff;background-color: #01a4ff;text-decoration: none;margin-bottom: 35px;"
	  href="http://192.168.0.237/account/active?s=%s">验证邮箱</a>
	<p style="font-size: 16px;
  color: #333;line-height: 28px;margin: 7px 0;">如果按钮不生效，请复制下面的链接到您的浏览器:</p>
	<span style="font-size:18px;color:#01a4ff;text-decoration: none;" >%s</span>
	<p style="margin: 7px 0 30px;font-size: 16px">如需帮助，请联系<a
	  style="color:#01a4ff;text-decoration: none;" href="">技术支持</a></p>
	<img style="vertical-align: middle;margin-right: 16px;"
		 src="data:image/png;base64,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" width="54" height="31" alt="">
	<a target="_blank" href="//lianbai.io" style="color:#01a4ff;text-decoration: none;" href="">链佰信息科技有限公司</a> <span
	style="color:#151515;">|</span> <a target="_blank"
	style="color:#01a4ff;text-decoration: none;" href="https://lianbai.io">链佰BAAS平台</a>`

	ForgotCode = `<div style="padding: 6px;font-size: 16px;">
	<h1 style="margin: 0 0 7px;font-size: 19px;
  color: #333;
  line-height: 28px;">%s，您好！</h1>
	<p style="margin: 7px 0 18px;font-size: 16px;
  color: #333;line-height: 28px;">感谢对链佰BAAS平台支持！为了重设账户密码，请点击下面的按钮修改密码</p>
	<a target="_blank"
	  style="display:block;width: 140px;height:36px;font-size:16px;line-height:36px;text-align:center;color:#fff;background-color: #01a4ff;text-decoration: none;margin-bottom: 35px;"
	  href="http://192.168.0.237/account/forgot?s=%s">修改密码</a>
	<p style="font-size: 16px;
  color: #333;line-height: 28px;margin: 7px 0;">如果按钮不生效，请复制下面的链接到您的浏览器:</p>
	<span style="font-size:18px;color:#01a4ff;text-decoration: none;" >%s</span>
	<p style="margin: 7px 0 30px;font-size: 16px">如需帮助，请联系<a
	  style="color:#01a4ff;text-decoration: none;" href="">技术支持</a></p>
	<img style="vertical-align: middle;margin-right: 16px;"
		 src="data:image/png;base64,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" width="54" height="31" alt="">
	<a target="_blank" href="//lianbai.io" style="color:#01a4ff;text-decoration: none;" href="">链佰信息科技有限公司</a> <span
	style="color:#151515;">|</span> <a target="_blank"
	style="color:#01a4ff;text-decoration: none;" href="https://lianbai.io">链佰BAAS平台</a>`
)
