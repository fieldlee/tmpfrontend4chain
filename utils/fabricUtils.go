package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"frontend4chain/constant"
	"frontend4chain/module"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReplaceCertWithValue(define module.Define, replaceString string) string {
	// -OrderIDForReplace-
	// -OrderPortForReplace-
	// -OrderNameForReplace-
	// -DomainForReplace-
	if strings.HasPrefix(replaceString, "#######OrderList-Start") {
		orderReplaceStr := ""
		for _, order := range define.Orders {
			// -OrderPortForReplace-
			tmpReplaceStr := strings.Replace(replaceString, "-OrderIDForReplace-", order.OrderId, -1)
			tmpReplaceStr = strings.Replace(tmpReplaceStr, "-OrderPortForReplace-", strconv.Itoa(order.OrderPort), -1)
			tmpReplaceStr = strings.Replace(tmpReplaceStr, "-OrderNameForReplace-", order.OrderName, -1)
			tmpReplaceStr = strings.Replace(tmpReplaceStr, "-DomainForReplace-", define.Domain, -1)
			orderReplaceStr = fmt.Sprintln(orderReplaceStr, tmpReplaceStr)
		}
		return orderReplaceStr
	}
	// fmt.Println(replaceString)
	if strings.HasPrefix(replaceString, "#######List-Start") {

		orgReplaceStr := ""
		for _, org := range define.Orgs {
			tmporgReplaceStr := strings.Replace(replaceString, "-OrgIDForReplace-", org.OrgId, -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-OrgNameForReplace-", org.OrgName, -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-AnchorIpForReplace-", org.AnchorIp, -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-AnchorPortForReplace-", strconv.Itoa(org.AnchorPort), -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-DomainForReplace-", define.Domain, -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-CAIpForReplace-", org.CaIp, -1)
			tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-ProjectDIR-", fmt.Sprint(constant.ROOTPATH, "/", define.ID), -1)
			// strings.Replace(newReplacedStr, "-ProjectDIR-", fmt.Sprint(constant.ROOTPATH, "/", define.ID), -1)
			orgReplaceStr = fmt.Sprintln(orgReplaceStr, tmporgReplaceStr)
		}
		return orgReplaceStr
	}

	newReplacedStr := strings.Replace(replaceString, "-OrderIDForReplace-", define.OrderId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrderNameForReplace-", define.OrderName, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", define.Domain, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-ProjectDIR-", fmt.Sprint(constant.ROOTPATH, "/", define.ID), -1)

	// -OrgIDForReplace-
	// -OrgNameForReplace-
	// -AnchorIpForReplace-
	// -AnchorPortForReplace-

	return newReplacedStr
}

// 根据模板生成crypto yaml文件
func CryptoYaml(define module.Define, org module.Org, templatePath string, outPath string) error {
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()

	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(ReplaceCryptoWithValue(define, org, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

func ReplaceCryptoWithValue(define module.Define, org module.Org, replaceString string) string {
	tmporgReplaceStr := strings.Replace(replaceString, "-OrgIDForReplace-", org.OrgId, -1)
	tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-DomainForReplace-", define.Domain, -1)
	tmporgReplaceStr = strings.Replace(tmporgReplaceStr, "-CAIpForReplace-", org.CaIp, -1)
	return tmporgReplaceStr
}

// 根据模板生成yaml文件(config yaml 和 crypto yaml文件)
func YamltoYaml(define module.Define, templatePath string, outPath string) error {
	var orderStart bool = false
	var orgStart bool = false
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()

	inputread := bufio.NewReader(fh)

	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')

		if ferr == io.EOF {

			return nil
		}
		if strings.HasPrefix(input, "#######OrderList-Start") {
			orderStart = true
		}

		if strings.HasPrefix(input, "#######OrderList-End") {
			orderStart = false
		}

		if strings.HasPrefix(input, "#######List-Start") {
			orgStart = true
		}

		if strings.HasPrefix(input, "#######List-End") {
			orgStart = false
		}

		toReplaceString = fmt.Sprint(toReplaceString, input)
		if orderStart == false && orgStart == false {
			outputWriter.WriteString(ReplaceCertWithValue(define, toReplaceString))
			toReplaceString = ""
		}
	}

	return nil
}

// Order replace yaml 文件
func OrderReplaceCertWithValue(define module.Define, order module.Order, replaceString string) string {
	// fmt.Println(replaceString)
	if strings.HasPrefix(replaceString, "########Order_Start") {
		orderReplaceStr := fmt.Sprintln("       - " + filepath.Join(constant.ROOTPATH, define.ID) + ":/etc/hyperledger/configtx")
		orderReplaceStr = fmt.Sprintln(orderReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config")+":/etc/hyperledger/fabric/crypto")
		orderReplaceStr = fmt.Sprintln(orderReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "ordererOrganizations", define.Domain, "orderers", order.OrderId+"."+define.Domain)+":/etc/hyperledger/crypto/orderer")
		for _, org := range define.Orgs {
			orderReplaceStr = fmt.Sprintln(orderReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "peerOrganizations", org.OrgId+"."+define.Domain, "peers", org.Peers[0].PeerId+"."+org.OrgId+"."+define.Domain)+":/etc/hyperledger/crypto/peer"+org.OrgId) //
		}
		return orderReplaceStr
	}

	if strings.HasPrefix(replaceString, "########Kafka_Start") {
		if define.Consensus == "kafka" {
			newReplacedStr := strings.Replace(replaceString, "-DomainForReplace-", define.Domain, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "-KafkaIpForReplace-", define.KafkaIp, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "########Kafka_Start", "", -1)
			newReplacedStr = strings.Replace(newReplacedStr, "########Kafka_End", "", -1)
			return newReplacedStr
		} else {
			return ""
		}
	}

	if strings.HasPrefix(replaceString, "########Kafka_OneIp_Start") {
		if define.Consensus == "kafka" && define.KafkaIp == order.OrderIp {
			newReplacedStr := strings.Replace(replaceString, "-DomainForReplace-", define.Domain, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "########Kafka_OneIp_Start", "", -1)
			newReplacedStr = strings.Replace(newReplacedStr, "########Kafka_OneIp_End", "", -1)
			return newReplacedStr
		} else {
			return ""
		}
	}
	// -OrderIdForReplace-.-DomainForReplace-
	newReplacedStr := strings.Replace(replaceString, "-OrderIdForReplace-", order.OrderId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrderPortForReplace-", strconv.Itoa(order.OrderPort), -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", define.Domain, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrderIDForReplace-", define.OrderId, -1)
	// -OrderPortForReplace-
	// -OrderIdForReplace-
	// -OrderIDForReplace
	// -ContainerIdForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-ContainerIdForReplace-", order.ContainerId, -1)
	return newReplacedStr
}

// 生成order yaml 文件
func OrderYaml(define module.Define, order module.Order, templatePath string, outPath string) error {
	var orderStart bool = false
	var kafkaStart bool = false
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()

	inputread := bufio.NewReader(fh)

	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()

	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')

		if ferr == io.EOF {
			return nil
		}
		if strings.HasPrefix(input, "########Order_Start") {
			orderStart = true
		}

		if strings.HasPrefix(input, "########Order_End") {
			orderStart = false
		}

		if strings.HasPrefix(input, "########Kafka_Start") || strings.HasPrefix(input, "########Kafka_OneIp_Start") {
			kafkaStart = true
		}

		if strings.HasPrefix(input, "########Kafka_End") || strings.HasPrefix(input, "########Kafka_OneIp_End") {
			kafkaStart = false
		}

		toReplaceString = fmt.Sprint(toReplaceString, input)
		if orderStart == false && kafkaStart == false {
			outputWriter.WriteString(OrderReplaceCertWithValue(define, order, toReplaceString))
			toReplaceString = ""
		}
	}

	return nil
}

// KafkaReplaceCertWithValue

func KafkaReplaceCertWithValue(define module.Define, replaceString string) string {
	// -OrderIdForReplace-.-DomainForReplace-
	newReplacedStr := strings.Replace(replaceString, "-DomainForReplace-", define.Domain, -1)

	return newReplacedStr
}

// 生成kafka yaml 文件
func KafkaYaml(define module.Define, templatePath string, outPath string) error {
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(KafkaReplaceCertWithValue(define, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

// CaReplaceCertWithValue
func CaReplaceCertWithValue(define module.Define, org module.Org, replaceString string) string {
	// -OrgCAIDForReplace-.-OrgIDForReplace-.-DomainForReplace-
	// -CAPEMFILENAMEForReplace-
	// -OrgCAIDForReplace-.-OrgIDForReplace-.-DomainForReplace
	newReplacedStr := strings.Replace(replaceString, "-OrgCAIDForReplace-", org.CaId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrgIDForReplace-", org.OrgId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", define.Domain, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-ContainerIdForReplace-", org.ContainerId, -1)
	// 	替换证书名称
	caCertFilePath := filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "peerOrganizations", org.OrgId+"."+define.Domain, "ca")
	rd, err := ioutil.ReadDir(caCertFilePath)
	if err != nil {
		fmt.Println("read dir error!")
		return ""
	}
	for _, file := range rd {
		if strings.HasSuffix(file.Name(), "_sk") {
			newReplacedStr = strings.Replace(newReplacedStr, "-CAPEMFILENAMEForReplace-", file.Name(), -1)
			break
		}
	}
	// -CAPATHForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-CAPATHForReplace-", caCertFilePath, -1)
	// -OrgNameForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-OrgNameForReplace-", org.OrgName, -1)
	// -OrgCAPortForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-OrgCAPortForReplace-", strconv.Itoa(org.CaPort), -1)
	return newReplacedStr
}

// 生成ca yaml 文件
func CaYaml(define module.Define, org module.Org, templatePath string, outPath string) error {
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(CaReplaceCertWithValue(define, org, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

// Cli yaml
func CliReplaceCertWithValue(define module.Define, org module.Org, replaceString string) string {

	newReplacedStr := strings.Replace(replaceString, "-OrgIDForReplace-", org.OrgId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", define.Domain, -1)
	// -CryptoPathForReplace-
	cryptoPath := filepath.Join(constant.ROOTPATH, define.ID, "crypto-config")
	newReplacedStr = strings.Replace(newReplacedStr, "-CryptoPathForReplace-", cryptoPath, -1)
	//-OrderDomainForReplace-:-OrderIPForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-OrderDomainForReplace-", fmt.Sprint(define.Orders[0].OrderId, ".", define.Domain), -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrderIPForReplace-", define.Orders[0].OrderIp, -1)
	return newReplacedStr
}

// 生成cli yaml 文件
func CliYaml(define module.Define, org module.Org, templatePath string, outPath string) error {
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(CliReplaceCertWithValue(define, org, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

// PeerReplaceCertWithValue
func PeerReplaceCertWithValue(define module.Define, org module.Org, peer module.Peer, replaceString string) string {
	// -PeerIdForReplace-.-OrgIDForReplace-.-DomainForReplace- -CouchPortForReplace-

	if strings.HasPrefix(replaceString, "#########COUCH_Start") {
		if peer.JoinCouch {
			newReplacedStr := strings.Replace(replaceString, "-CouchIDForReplace-", peer.CouchId, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "-CouchUserForReplace-", peer.CouchUsername, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "-CouchPasswordForReplace-", peer.CouchPassword, -1)
			newReplacedStr = strings.Replace(newReplacedStr, "-CouchPortForReplace-", strconv.Itoa(peer.CouchPort), -1)
			newReplacedStr = strings.Replace(newReplacedStr, "-CouchContainerIdForReplace-", peer.CouchContainerId, -1)
			return newReplacedStr
		} else {
			return ""
		}
	}

	// 	#########DependOn_Start
	//     depends_on:
	//       - -CouchIDForReplace-
	// #########DependOn_End

	if strings.HasPrefix(replaceString, "#########DependOn_Start") {
		hadhost := false
		dependReplaceStr := ""
		if peer.JoinCouch {
			if hadhost == false {
				dependReplaceStr = fmt.Sprintln("    depends_on:")
				hadhost = true
			}
			dependReplaceStr = fmt.Sprintln(dependReplaceStr, "      - "+peer.CouchContainerId)

		}

		for _, order := range define.Orders {
			if order.OrderIp == peer.PeerIp {
				if hadhost == false {
					dependReplaceStr = fmt.Sprintln("    depends_on:")
					hadhost = true
				}
				dependReplaceStr = fmt.Sprintln(dependReplaceStr, "      - "+order.ContainerId)
			}
		}

		return dependReplaceStr
	}

	// #########OrderExtra_Start
	// -OrderIDForReplace-.-DomainForReplace-:-OrderIPForReplace-
	if strings.HasPrefix(replaceString, "#########OrderExtra_Start") {
		extraReplaceStr := ""
		hadhost := false
		for _, order := range define.Orders {
			if peer.PeerIp != order.OrderIp {
				if hadhost == false {
					extraReplaceStr = fmt.Sprintln("    extra_hosts:")
					hadhost = true
				}
				extraReplaceStr = fmt.Sprintln(extraReplaceStr, "      - "+order.ContainerId+":"+order.OrderIp)
			}
		}
		return extraReplaceStr
	}

	// #########PeerList_Start
	if strings.HasPrefix(replaceString, "#########PeerList_Start") {
		listReplaceStr := fmt.Sprintln("       - /var/run/:/host/var/run/")
		listReplaceStr = fmt.Sprintln(listReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "peerOrganizations", org.OrgId+"."+define.Domain, "peers", peer.PeerId+"."+org.OrgId+"."+define.Domain)+":/etc/hyperledger/crypto/peer")
		listReplaceStr = fmt.Sprintln(listReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "peerOrganizations", org.OrgId+"."+define.Domain, "peers", peer.PeerId+"."+org.OrgId+"."+define.Domain, "msp")+":/etc/hyperledger/fabric/msp")
		listReplaceStr = fmt.Sprintln(listReplaceStr, "      - "+filepath.Join(constant.ROOTPATH, define.ID, "crypto-config", "peerOrganizations", org.OrgId+"."+define.Domain, "peers", peer.PeerId+"."+org.OrgId+"."+define.Domain, "tls")+":/etc/hyperledger/fabric/tls")
		return listReplaceStr
	}

	newReplacedStr := strings.Replace(replaceString, "-PeerIdForReplace-", peer.PeerId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-OrgIDForReplace-", org.OrgId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", define.Domain, -1)
	// -ContainerIdForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-ContainerIdForReplace-", peer.ContainerId, -1)
	// -PeerPortForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-PeerPortForReplace-", strconv.Itoa(peer.PostPort), -1)
	// -PeerEventPortForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-PeerEventPortForReplace-", strconv.Itoa(peer.EventPort), -1)
	// -NetWorkForReplace-
	newReplacedStr = strings.Replace(newReplacedStr, "-NetWorkForReplace-", define.ID, -1) //define.NetWork
	return newReplacedStr
}

// peer 生成
func PeerYaml(define module.Define, org module.Org, peer module.Peer, templatePath string, outPath string) error {

	var peerlistStart bool = false
	var extraStart bool = false
	var couchStart bool = false
	var depengStart bool = false
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {
			return nil
		}

		if strings.HasPrefix(input, "#########COUCH_Start") {
			couchStart = true
		}

		if strings.HasPrefix(input, "#########COUCH_End") {
			couchStart = false
		}

		if strings.HasPrefix(input, "#########PeerList_Start") {
			peerlistStart = true
		}

		if strings.HasPrefix(input, "#########PeerList_End") {
			peerlistStart = false
		}

		if strings.HasPrefix(input, "#########OrderExtra_Start") {
			extraStart = true
		}

		if strings.HasPrefix(input, "#########OrderExtra_End") {
			extraStart = false
		}

		if strings.HasPrefix(input, "#########DependOn_Start") {
			depengStart = true
		}

		if strings.HasPrefix(input, "#########DependOn_End") {
			depengStart = false
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		if extraStart == false && peerlistStart == false && couchStart == false && depengStart == false {
			outputWriter.WriteString(PeerReplaceCertWithValue(define, org, peer, toReplaceString))
			toReplaceString = ""
		}
	}
	return nil
}

func CouchReplaceCertWithValue(peer module.Peer, replaceString string) string {
	// -CouchIdForReplace-
	// -CouchUserForReplace-
	// -CouchPasswordForReplace-
	// -CouchPortForReplace-
	newReplacedStr := strings.Replace(replaceString, "-CouchIdForReplace-", peer.CouchId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-CouchUserForReplace-", peer.CouchUsername, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-CouchPasswordForReplace-", peer.CouchPassword, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-CouchPortForReplace-", strconv.Itoa(peer.CouchPort), -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-ContainerIdForReplace-", peer.CouchContainerId, -1)
	return newReplacedStr
}

func CouchYaml(peer module.Peer, templatePath string, outPath string) error {
	// 读取template yaml 文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)

	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(CouchReplaceCertWithValue(peer, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

func CryptoReplaceCertWithValue(project module.Define, org module.Org, replaceString string) string {
	newReplacedStr := strings.Replace(replaceString, "-OrgIDForReplace-", org.OrgId, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-DomainForReplace-", project.Domain, -1)
	newReplacedStr = strings.Replace(newReplacedStr, "-ProjectDIR-", filepath.Join(constant.ROOTPATH, project.ID), -1)
	return newReplacedStr
}

func AddOrgCrytojson(project module.Define, org module.Org, templatePath string, outPath string) error {
	// 读取template yaml文件
	fh, ferr := os.Open(templatePath)
	if ferr != nil {
		return ferr
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)
	// 写人yaml文件
	if isExist, _ := IsExistFile(outPath); isExist == true {
		err := os.Remove(outPath)
		if err != nil {
			fmt.Printf("An error occurred with file creation\n")
			return err
		}
	}
	// 创建yaml文件
	CreateFile(outPath)
	outputFile, outputError := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()
	// 	写人内容到yaml文件
	toReplaceString := ""
	for {
		input, ferr := inputread.ReadString('\n')
		if ferr == io.EOF {

			return nil
		}
		toReplaceString = fmt.Sprint(toReplaceString, input)
		outputWriter.WriteString(CryptoReplaceCertWithValue(project, org, toReplaceString))
		toReplaceString = ""
	}
	return nil
}

// 发送数据
func SendFile(fileName string, ip string, port string, id string, meschan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			meschan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	uploadUrl := fmt.Sprint(host, constant.UPLOADURL)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("certTar", id+".tar") //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		panic(err)
	}

	fd, err := os.Open(fileName)
	if err != nil {
		fmt.Println("d")
		panic(err)
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)

	if err != nil {
		fmt.Println("e")
		panic(err)
	}
	//
	fwId, err := w.CreateFormField("id")
	if err != nil {
		panic(err)
	}
	fwId.Write([]byte(id))
	w.Close()

	req, err := http.NewRequest("POST", uploadUrl, buf)
	if err != nil {
		fmt.Println("f")
		panic(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	meschan <- string(message)
}

// 发送channel 文件
func SendChannelFile(filename string, filepath string, ip string, port string, id string, sendchan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			sendchan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	uploadUrl := fmt.Sprint(host, constant.UPLOADTXURL)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("file", filename) //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		panic(err)
	}

	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Println("d")
		panic(err)
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)

	if err != nil {
		fmt.Println("e")
		panic(err)
	}
	//
	fwId, err := w.CreateFormField("id")
	if err != nil {
		panic(err)
	}
	fwId.Write([]byte(id))
	w.Close()

	req, err := http.NewRequest("POST", uploadUrl, buf)
	if err != nil {
		fmt.Println("f")
		panic(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	sendchan <- string(message)
}

// 发送压缩包智能合约文件
func SendCCFile(filename string, filepath string, ip string, port string, ccpath string, ccname string, sendchan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			sendchan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	uploadUrl := fmt.Sprint(host, constant.CCUPLOADZIPURL)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("chaincodeZip", filename) //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		panic(err)
	}
	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Println("d")
		panic(err)
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		fmt.Println("e")
		panic(err)
	}
	// ccpath
	fwpath, err := w.CreateFormField("chaincodePath")
	if err != nil {
		panic(err)
	}
	fwpath.Write([]byte(ccpath))

	// chaincode name id
	fwCcname, err := w.CreateFormField("chaincodeName")
	if err != nil {
		panic(err)
	}
	fwCcname.Write([]byte(ccname))
	w.Close()

	req, err := http.NewRequest("POST", uploadUrl, buf)
	if err != nil {
		fmt.Println("f")
		panic(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	sendchan <- string(message)
}

func GetSysEnv(ip string, port string, meschan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("===========GetSysEnv,recover")
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			meschan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	sysUrl := fmt.Sprint(host, constant.SYSTEMINFOURL)

	req, err := http.NewRequest("POST", sysUrl, bytes.NewReader([]byte("")))
	if err != nil {
		fmt.Println("LunchDockerEnv -- b")
		panic(err)
	}
	// Content-Type:application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("LunchDockerEnv -- c")
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	meschan <- string(message)
}

func CheckIPEnv(ip string, port string, command string, meschan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			meschan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	chekcUrl := fmt.Sprint(host, constant.CHECKURL)

	jsonMap := make(map[string]interface{})
	jsonMap["command"] = command

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("LunchDockerEnv -- a")
		panic(err)
	}

	req, err := http.NewRequest("POST", chekcUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		fmt.Println("LunchDockerEnv -- b")
		panic(err)
	}
	// Content-Type:application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("LunchDockerEnv -- c")
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	meschan <- string(message)
}

func InstallIPEnv(ip string, port string, command string, meschan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			meschan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	installUrl := fmt.Sprint(host, constant.INSTALLURL)

	jsonMap := make(map[string]interface{})
	jsonMap["command"] = command

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("LunchDockerEnv -- a")
		panic(err)
	}

	req, err := http.NewRequest("POST", installUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		fmt.Println("LunchDockerEnv -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("LunchDockerEnv -- c")
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	meschan <- string(message)
}

func LunchDockerEnv(ip string, port string, command string, id string, meschan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			meschan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	lunchUrl := fmt.Sprint(host, constant.LUNCHURL)

	jsonMap := make(map[string]interface{})
	jsonMap["command"] = command
	jsonMap["id"] = id //id 没有用上删除

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("LunchDockerEnv -- a")
		panic(err)
	}

	req, err := http.NewRequest("POST", lunchUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("LunchDockerEnv -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("LunchDockerEnv -- c")
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	meschan <- string(message)
}

// 发送sdk 配置文件和git 路径
func SdkSendConfig(ip string, port string, id string, config []byte, sdkchan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			sdkchan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	configUrl := fmt.Sprint(host, constant.SDKINSTALLURL)

	jsonMap := make(map[string]interface{})
	jsonMap["sdkPath"] = constant.SDKPATH
	jsonMap["id"] = id
	jsonMap["sdkGit"] = constant.SDKGITURL
	jsonMap["sdkConfig"] = config

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("sdkSendConfig -- a")
		panic(err)
	}

	req, err := http.NewRequest("POST", configUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("sdkSendConfig -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	encode := HlcEncode(constant.CURIP)
	fmt.Sprintln("hlcEncode =========:", encode)
	req.Header.Set("Authorization", encode)
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("sdkSendConfig -- c")
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	sdkchan <- string(message)
}

// 发送sdk 配置文件和git 路径
func ChaincodeConfig(ip string, port string, ccpath string, ccname string, ccgit string, ccchan chan string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			ccchan <- "{'code':400,'msg':'','data':{}}"
		}
	}()

	host := fmt.Sprint("http://", ip, ":", port)
	configUrl := fmt.Sprint(host, constant.CCUPLOADGITURL)

	jsonMap := make(map[string]interface{})
	jsonMap["chaincodePath"] = ccpath
	jsonMap["chaincodeName"] = ccname
	jsonMap["chaincodeGit"] = ccgit

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("sdkSendConfig -- a")
		panic(err)
	}

	req, err := http.NewRequest("POST", configUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("sdkSendConfig -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("sdkSendConfig -- c")
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	ccchan <- string(message)
}

////调用shell 命令
func ReqShellServer(ip string, cmd string) []byte {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	// 请求服务器，
	host := fmt.Sprint("http://", ip, ":", constant.SERVERPORT)
	shellUrl := fmt.Sprint(host, constant.SHELLEXECURL)
	fmt.Println(cmd)
	jsonMap := make(map[string]interface{})
	jsonMap["shellContent"] = []byte(cmd)
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("shell -- a")
		panic(err)
	}
	req, err := http.NewRequest("POST", shellUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("shell -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("shell -- c")
		panic(err)
	}
	defer res.Body.Close()

	message, _ := ioutil.ReadAll(res.Body)
	return message
}

////调用shell 命令
func GetProjectPerformance(ip string, msgchan chan string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	// 请求服务器，
	host := fmt.Sprint("http://", ip, ":", constant.SERVERPORT)
	perUrl := fmt.Sprint(host, constant.SYSTEMPERFORMURL)

	// jsonMap := make(map[string]interface{})
	// jsonBytes, err := json.Marshal(jsonMap)
	// if err != nil {
	// 	fmt.Println("shell -- a")
	// 	panic(err)
	// }
	req, err := http.NewRequest("POST", perUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("shell -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("shell -- c")
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	msgMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(message), &msgMap)
	if err != nil {
		msgMap["ip"] = ip
		msgMap["code"] = 400
	} else {
		msgMap["ip"] = ip
	}
	msgBytes, _ := json.Marshal(msgMap)
	msgchan <- string(msgBytes)
}

////调用shell 命令
func GetDockerPerformance(ip string, containerlist []string, msgchan chan string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	// 请求服务器，
	host := fmt.Sprint("http://", ip, ":", constant.SERVERPORT)
	perUrl := fmt.Sprint(host, constant.DOCKERPERFORMURL)

	jsonMap := make(map[string]interface{})
	jsonMap["dockers"] = containerlist
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("shell -- a")
		panic(err)
	}
	req, err := http.NewRequest("POST", perUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("shell -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("shell -- c")
		panic(err)
	}
	defer res.Body.Close()

	message, _ := ioutil.ReadAll(res.Body)
	msgchan <- string(message)
}

////调用shell 命令
func DockerExec(ip string, containerid string, typestring string, msgchan chan string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	// 请求服务器，
	host := fmt.Sprint("http://", ip, ":", constant.SERVERPORT)
	perUrl := fmt.Sprint(host, constant.DOCKEREXEC)

	jsonMap := make(map[string]interface{})
	jsonMap["container_id"] = containerid
	jsonMap["type"] = typestring
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("shell -- a")
		panic(err)
	}
	req, err := http.NewRequest("POST", perUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("shell -- b")
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", HlcEncode(constant.CURIP))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("shell -- c")
		panic(err)
	}
	defer res.Body.Close()

	message, _ := ioutil.ReadAll(res.Body)
	msgchan <- string(message)
}
