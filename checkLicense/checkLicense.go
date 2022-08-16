package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// global
var api = "https://license.chinauos.com/v1/openapi/legalquery?_allow_anonymous=true&code="
var licenseStr2 = "激活码,授权产品,使用状态,激活码类型,服务类型,使用时间\n"
var keys = make([]string, 10)
var JsonChan = make(chan string, 10)
var JsonExitChan = make(chan bool, 1)

func main() {
	WriteCSV(licenseStr2)
	// read license
	ReadKey()
	// get info
	GetInfo()

	// start
	go ParseInfo()
	for {
		_, ok := <-JsonExitChan
		if !ok {
			break
		}
	}
}

// inti query native license
func init() {
	// only support  1050u2
}
func ReadKey() {
	sh := "find /home/`w | awk 'NR==3{print $1}'` -name License.txt |xargs cat"
	res := strings.Split(GetBashRet(sh), "\n")
	for _, license := range res {
		if license != " " || license != "" {
			keys = append(keys, license)
		}
	}
}

type TopString struct {
	Code int `json:"code"`
	Row  struct{ Row }
}

type Row struct {
	Codes            int    `json:"code"`
	Status           int    `json:"status"`
	NameCn           string `json:"name_cn"`
	AuthoLimit       int    `json:"autho_limit"`
	Days             int    `json:"days"`
	StartAt          int    `json:"start_at"`
	AsofAt           int    `json:"asof_at"`
	AuthoMode        int    `json:"autho_mode"`
	ManuName         string `json:"manu_name"`
	ManuDeviceName   string `json:"manu_device_name"`
	ServiceType      int    `json:"service_type"`
	ServiceStartTime int    `json:"service_start_time"`
	ServiceEndTime   int    `json:"service_end_time"`
	FlowAuthLimit    int    `json:"flow_auth_limit"`
	FixedActiveTime  int    `json:"fixed_active_time"`
	FixedOverdueTime int    `json:"fixed_overdue_time"`
	DeviceBrand      string `json:"device_brand"`
	FixedYears       int    `json:"fixed_years"`
	FixedMonths      int    `json:"fixed_months"`
	FixedDays        int    `json:"fixed_days"`
	ServiceTimeYear  int    `json:"service_time_year"`
	ServiceTimeMonth int    `json:"service_time_month"`
	ServiceTimeDay   int    `json:"service_time_day"`
	SysActiveDay     int    `json:"sys_active_day"`
	UseStatus        int    `json:"use_status"`
	UseTime          int    `json:"use_time"`
}

// get info
func GetInfo() {
	JsonExitChan <- false
	for _, v := range keys {
		if v != "" {
			url := api + v
			JsonChan <- url
		}
	}
	close(JsonChan)
}

// ParseInfo get license info json
func ParseInfo() {
	for url := range JsonChan {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("http get error", err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		ParseJson(string(body), url)
	}

	for {
		_, ok := <-JsonChan
		if !ok {
			break
		}
	}
	JsonExitChan <- true
	close(JsonExitChan)
}

// ParseJson parse json of license info
func ParseJson(str string, license string) {
	var s2 TopString
	err := json.Unmarshal([]byte(str), &s2)
	if err != nil {
		fmt.Println("error", err)
	}
	license = strings.ReplaceAll(license, api, "")

	fmt.Println("============激活码查询结果===============")
	// autho_mode 1 open 2 oem
	// NameCn  授权产品
	// service_type 1 标准服务  2 金牌服务 3 现场服务 4 无售后服务
	// use_status   1 已使用  2 未使用
	licenseStr := ""
	fmt.Println("激活码：", license)
	licenseStr += license + ","
	fmt.Println("授权产品：", s2.Row.NameCn)
	licenseStr += s2.Row.NameCn + ","
	if s2.Row.UseStatus == 1 {
		fmt.Println("使用状态： 已使用")
		licenseStr += "已使用,"
	} else {
		fmt.Println("使用状态： 未使用")
		licenseStr += "未使用,"
	}
	if s2.Row.AuthoMode == 1 {
		fmt.Println("激活码类型：", "Open")
		licenseStr += "Open,"
	} else {
		fmt.Println("激活码类型：", "OEM")
		licenseStr += "OEM,"
	}
	if s2.Row.ServiceType == 1 {
		fmt.Println("服务类型：", "标准服务")
		licenseStr += "标准服务,"
	} else if s2.Row.ServiceType == 2 {
		fmt.Println("服务类型：", "金牌服务")
		licenseStr += "金牌服务,"
	} else if s2.Row.ServiceType == 3 {
		fmt.Println("服务类型：", "现场服务")
		licenseStr += "现场服务,"
	} else {
		fmt.Println("服务类型：", "无售后服务")
		licenseStr += "无售后服务,"
	}
	fmt.Println("使用时间：", s2.Row.FixedActiveTime)
	licenseStr += string(s2.Row.FixedActiveTime) + "\n"
	WriteCSV(licenseStr)
	fmt.Println("=======================================\n")
}

// GetBashRet get bash Return
func GetBashRet(cmd string) string {
	c := exec.Command("bash", "-c", cmd)
	output, err := c.Output()
	if err != nil {
		fmt.Println("exec shell error, check please ...")
	}
	return string(output)
}

func WriteCSV(str string) {
	filePath := "a.csv"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(str)
	write.Flush()
}
