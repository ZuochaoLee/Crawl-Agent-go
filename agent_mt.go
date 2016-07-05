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
func fetch(url *string, use_ip, ckaua string) (html string) {
	fmt.Println(*url)
	spec := redis.DefaultSpec().Host("localhost").Port(6379).Db(1).Password("")
	cli, _ := redis.NewSynchClientWithSpec(spec)
	proxy_addr, key, num := "", "", 0
	if use_ip == "1" {
		proxy_addr, key, num = getIp()
	}
	transport := getTransportFieldURL(&proxy_addr)
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

	} else {
		fmt.Println("11")
		if resp.StatusCode == 200 {
			fmt.Println("22")
			robots, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Println("返回失败")
			}
			html = string(robots)
		} else {
			html = ""
		}
	}
	if html == "" {
		if num > 3 {
			cli.Del(key)
			//break
		} else {
			cli.Set(key, []byte(strconv.Itoa(num+1)))
		}
	}
	return
}

//获取动态IP
func getIp() (proxy_addr string, key string, num int) {
	//IP池链接
	spec := redis.DefaultSpec().Host("localhost").Port(6379).Db(1).Password("")
	client, _ := redis.NewSynchClientWithSpec(spec)
	//随机取IP
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
func worker(intime, use_ip, ckaua string, jobs <-chan *string, results chan<- string) {
	for j := range jobs {
		fmt.Println(*j)
		html := fetch(j, use_ip, ckaua)
		fmt.Println("********")
		fmt.Println(html)
		delayRand(intime)
		results <- html
	}
}
func main() {
	url := "http://www.baidu.com/" //目标地址
	use_ip := "1"                  //是否用IP代理 1:使用 0:不使用
	ckaua := "1"                   //是否禁用cookie  1:禁用cookie 随机ua 0:记录cookie 固定ua
	intime := "10"                 //最大时延 ms  不用延时设为0
	//两个channel，一个用来放置工作项，一个用来存放处理结果。
	jobs := make(chan *string, 100)
	results := make(chan string, 100)
	fmt.Println(url)
	for w := 1; w <= 200; w++ {
		//是否使用IP代理
		go worker(intime, use_ip, ckaua, jobs, results)
	}
	// 添加1000个任务后关闭Channel
	for j := 1; j <= 1000; j++ {
		//url=url+strconv.Itoa(j)
		jobs <- &url
	}
	close(jobs)
	//获取所有的处理结果
	for a := 1; a <= 1000; a++ {
		<-results
	}
}
