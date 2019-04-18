package main

import (
	"fmt"
	"net"
	"os"
	"flag"
	"time"
	"log"
	"sync"
	"runtime"
)

type Workdist struct {
	Host	string
}

const (
	taskload		    = 255
	tasknum			= 255
)
var wg sync.WaitGroup

func Task(ip string){
	tasks := make(chan Workdist,taskload)
	wg.Add(tasknum)
	for gr:=1;gr<=tasknum;gr++ {
		go ping(tasks)
	}

	for i:=1;i<256;i++ {
		host:=fmt.Sprintf("%s.%d",ip,i)
		task := Workdist{
			Host:host,
		}
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

var iphalf = flag.String("iphalf", "192.168.0", "example: 192.168.0")
var starttime int64
func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	starttime = time.Now().Unix()
	Task(*iphalf)
	endtime := time.Now().Unix()
	log.Print("Check complete.Use Time ",endtime - starttime," second")
}

func ping(tasks chan Workdist) {
	var size int
	var timeout int64
	defer wg.Done()

	size = 32
	timeout = 1000

	for {
		task,ok := <- tasks
		if !ok {
			return
		}
		host := task.Host
		starttime := time.Now()
		conn, err := net.DialTimeout("ip4:icmp", host, time.Duration(timeout*1000*1000))
		ip := conn.RemoteAddr()

		var seq int16 = 1
		id0, id1 := genidentifier(host)
		const ECHO_REQUEST_HEAD_LEN = 8

		shortT := -1
		longT := -1

		var msg []byte = make([]byte, size+ECHO_REQUEST_HEAD_LEN)
		msg[0] = 8                        // echo
		msg[1] = 0                        // code 0
		msg[2] = 0                        // checksum
		msg[3] = 0                        // checksum
		msg[4], msg[5] = id0, id1         //identifier[0] identifier[1]
		msg[6], msg[7] = gensequence(seq) //sequence[0], sequence[1]

		length := size + ECHO_REQUEST_HEAD_LEN

		check := checkSum(msg[0:length])
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 255)

		conn, err = net.DialTimeout("ip:icmp", host, time.Duration(timeout*1000*1000))

		checkError(err)

		starttime = time.Now()
		conn.SetDeadline(starttime.Add(time.Duration(timeout * 1000 * 1000)))
		_, err = conn.Write(msg[0:length])

		const ECHO_REPLY_HEAD_LEN = 20

		var receive []byte = make([]byte, ECHO_REPLY_HEAD_LEN+length)
		n, err := conn.Read(receive)
		_ = n
		var endduration int = int(int64(time.Since(starttime)) / (1000 * 1000))

		if err != nil || receive[ECHO_REPLY_HEAD_LEN+4] != msg[4] || receive[ECHO_REPLY_HEAD_LEN+5] != msg[5] || receive[ECHO_REPLY_HEAD_LEN+6] != msg[6] || receive[ECHO_REPLY_HEAD_LEN+7] != msg[7] || endduration >= int(timeout) || receive[ECHO_REPLY_HEAD_LEN] == 11 {
		} else {
			if shortT == -1 {
				shortT = endduration
			} else if shortT > endduration {
				shortT = endduration
			}
			if longT == -1 {
				longT = endduration
			} else if longT < endduration {
				longT = endduration
			}
			log.Print("扫描到主机地址:",ip.String())
		}
	}
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func gensequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genidentifier(host string) (byte, byte) {
	return host[0], host[1]
}


