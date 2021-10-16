package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)
//go build -ldflags="-s -w" -trimpath

var UrlFile string

func Banner(){
	fmt.Println(`用于快速判断url是否存活，读取当前目录下的为url.txt（可指定），可支持以下格式：
-eg:www.baidu.com
-eg:http://www.baidu.com
-eg:https://www.baidu.com`)
}

func Flag(){
	flag.StringVar(&UrlFile,"f","","url.txt文件")
	flag.Parse()
}

func HandleHttps(url string,client *http.Client)(NewUrl string, resp *http.Response,err error){
	//默认以http协议访问
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	resp, err = client.Do(req)
	if  err==nil{
		return url,resp,err
	}else {
		//将所有http无法访问的使用https来访问
		if url[:5] !="https"{
			url = strings.Replace(url,"http","https",1)
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		resp,err:=client.Do(req)
		//log.Print(err)

		if err != nil {
			//此处为http/https皆不能访问，需要contrue
			return url,resp,err
		}else {
			return url,resp,err
		}
	}
}

func handle(url string) (Url string, Title string, Power []string, StatusCode int) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: tr,
	}
	NewUrl,resp,err:=HandleHttps(url,client)
	if err!=nil{
		//log.Print("error Fail")
		return
	}else {
		var (
			title string
			code  int
		)

		Url, title, Power, code = handlebbody(NewUrl, resp)
		defer resp.Body.Close()
		return Url, title, Power, code
	}
}


func handlebbody(url string, resp *http.Response) (Url string, Title string, Power []string, StatusCode int) {
	var (
		code  int
		title string
	)
	code = resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile("<title>(.*)</title>")
	find := re.FindAllStringSubmatch(string(body), -1)
	if len(find) > 0 {
		title = find[0][1]
		if len(title) > 100 {
			title = title[:100]
		}
	} else {
		title = "None"
	}
	Power, ok := resp.Header["Server"]
	if !ok {
		Power = []string{"None"}
	}
	return url, title, Power, code
	//fmt.Printf("%s -> %s -> %s -> %d%s",url,title,Power,code,"\n")

}

func main() {
	Flag()
	if UrlFile ==""{
		Banner()
		flag.Usage()
		os.Exit(0)
	}
	var wg sync.WaitGroup
	start := time.Now()
	file, err := os.Open(UrlFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		i := scanner.Text()
		if i[:4] == "http" || i[:5] == "https" {

		} else {
			i = "http://" + i
		}
		//fmt.Println(i)
		go func(j string) {
			defer wg.Done()
			Url, Title, Power, StatusCode := handle(j)
			if len(Url) == 0 {
				return
			}
			fmt.Printf("%s -> %s -> %s -> %d\n", Url, Title, Power, StatusCode)
		}(i)
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	end := time.Since(start) / 1e9
	fmt.Printf("\n用时%d", end)
}
