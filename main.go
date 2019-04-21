package main

import (
	"time"
	"runtime"
	"fmt"
	"go-ipscan/tools"
	"os"
	"strings"
	"net"
	"log"
)

var starttime int64

var IpHalf_Str []string

var debugLog *log.Logger
var logFile *os.File

func main() {
	runtime.GOMAXPROCS(4)

	loginit()
	starttime = time.Now().Unix()
	Getlocaladdr()
	for _, ch := range IpHalf_Str {
		tools.Task(ch,debugLog)
	}
	endtime := time.Now().Unix()
	debugLog.Print("Check complete.Use Time ",endtime - starttime," second")
	logFile.Close()
}

func loginit()  {
	fileName := "log.txt"
	logFile,err  := os.Create(fileName)
	if err != nil {
		log.Fatalln("open file error !")
	}
	debugLog = log.New(logFile,"[Debug]",log.Llongfile)
}


func Getlocaladdr(){
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				a := strings.Split(ipnet.IP.String(), ".")
				if(a[3] == "1"){
					continue
				}
				Address_Processing(ipnet.IP.String())
			}
		}
	}
}

func Address_Processing(ip string){
	a := strings.Split(ip, ".")
	if(len(a)!=4){
		debugLog.Print("检查地址错误！")
		os.Exit(1)
	}

	host:=fmt.Sprintf("%s.%s.%s",a[0],a[1],a[2])
	debugLog.Print("获取本机地址:", ip,"------->  处理为网段:",host)
	IpHalf_Str = append(IpHalf_Str, host)
}