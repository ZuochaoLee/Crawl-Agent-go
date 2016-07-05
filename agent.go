package crawlagent

import (
	"fmt"
	//"log"
	redis "github.com/alphazero/Go-Redis"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	//"strings"
	"encoding/json"
	"math/rand"
	//"github.com/PuerkitoBio/goquery"
)

//指定代理ip
func getTransportFieldURL(proxy_addr *string) (transport *http.Transport) {
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(*proxy_addr)
	transport = &http.Transport{Proxy: http.ProxyURL(url_proxy)}
	return
}

//
func fetch(url, proxy_addr *string, ckaua string) (html string) {
	transport := getTransportFieldURL(proxy_addr)
	client := &http.Client{Transport: transport}
	//新建请求
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		// log.Fatal(err.Error())
		fmt.Println("新建失败")
		html = ""
		return
	}
	//设置请求头
	if ckaua == "0" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:6.0.2) Gecko/20100101 Firefox/6.0.2")
		req.Header.Set("Cookie", "name=anny")
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:6.0.2) Gecko/20100101 Firefox/6.0.2")
	}
	//
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败")
		html = ""
		return
	}
	//
	if resp.StatusCode == 200 {
		robots, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("返回失败")
		}
		html = string(robots)

	} else {
		html = ""
	}
	return
}

//获取动态IP
func getIp() (proxy_addr string, key string, num int) {
	//IP池链接
	spec := redis.DefaultSpec().Host("localhost").Port(6379).Db(1).Password("")
	client, err := redis.NewSynchClientWithSpec(spec)
	//随机取IP
	if err != nil {
		println("数据库连接失败")
	}
	key, _ = client.Randomkey()
	js, _ := client.Get(key)
	num, _ = strconv.Atoi(string(js))
	var dat map[string]string
	k := []byte(key)
	json.Unmarshal(k, &dat)
	proxy_addr = "http://" + dat["ip"] + ":" + dat["port"]
	fmt.Println(dat["ip"] + "-->" + string(js))
	return
}

//随机延时 ms
func delayRand(intime string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tm, _ := strconv.Atoi(string(intime))
	tt := r.Intn(tm)
	t := time.Duration(tt)
	time.Sleep(1000000 * t)
}
func main() {
	url := "http://www.baidu.com/" //目标地址
	use_ip := "1"                  //是否用IP代理 1:使用 0:不使用
	ckaua := "1"                   //是否禁用cookie  1:禁用cookie 随机ua 0:记录cookie 固定ua
	intime := "10"                 //最大时延 ms  不用延时设为0
	//spec := redis.DefaultSpec().Host("123.57.61.107").Port(6380).Db(11).Password("16c51b2287ed4bd2:zhugeZHAOFANG1116")
	spec := redis.DefaultSpec().Host("localhost").Port(6379).Db(1).Password("")
	client, _ := redis.NewSynchClientWithSpec(spec)
	html := ""
	proxy_addr := ""
	key := ""
	num := 0
	i := 1
	for true {
		//是否使用IP代理
		if use_ip == "1" {
			proxy_addr, key, num = getIp()
		}
		fetch(&url, &proxy_addr, ckaua)
		if html != "" {
			fmt.Println(html)
			//break
		} else {
			if num > 3 {
				client.Del(key)
				//break
			} else {
				client.Set(key, []byte(strconv.Itoa(num+1)))
			}
		}
		i++
		fmt.Println(i)

	}
	delayRand(intime)
}
