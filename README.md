# Crawl-Agent-go

#简介：

根据反爬方法，主要有以下反反爬措施

#反反爬策略：
##1、随机ua标记：

防止通过记录客户端ua信息反爬

解决办法：

伪装浏览器的user agent

Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0)
 
Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.2) 

Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)

##2、禁用cookie：

第一种是用户登录情况下

需要保存用户cookie 此时ua信息也要保存固定

第二种是非用户登录情况下

防止通过记录用户cookie反爬

解决办法：
禁用cookie 同时配合随机ua标记

##3、随机IP代理：

防止通过屏蔽IP段反爬

##4、随机访问延迟：

防止通过统计访问间隔反爬

服务接口设计：

服务接口应该是一个下载器 

可以选择设置三种服务（cookie&ua、IP代理、随机访问）

同时，在这优化一下并发多链接

##新增功能：
可以执行javascript代码，对进行简易加密的保护代码进行对称解密

#回退策略：
cookie&ua无成本，每次回退都要随机更换

IP只对接口设置启用IP代理的随机更换，并对失效IP标记，标记多次的移除IP池



#安装使用：

##依赖：
go get github.com/alphazero/Go-Redis

go get github.com/coocood/jas

go get github.com/mcuadros/go-candyjs

##安装：
go get git.oschina.net/buptlee/agent_crawl
