# acgdogeSpider

> 为了满足网站内容区分的需要，写了一个并发版的 Go 爬虫爬取 www.acgdoge.net 的各种文章内容

## 特点

完全多携程，包括页面链接的获取，内容的处理，保存数据库，全部是以多协程的方式并发实现的

## 协程分析

```text
====> Main Go-Routine Begin                                    ====>Main Go-Routine Destroy
       ====> Get Page 1
             ====> Use URL Get Content 1->10
                        ===>Format Data ==>Insert In To Mysql
                        ===>Format Data ==>Insert In To Mysql
                        ===>Format Data ==>Insert In To Mysql
                        ===>Format Data ==>Insert In To Mysql
       ====> Get Page 2
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
       ====> Get Page 3
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
       ====> Get Page 4
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
       ====> Get Page 5
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
       ====> Get Page 6
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
            .
            .
            .
       ====> Get Page n
             ====> Use URL Get Content 1->10
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql
                           ===>Format Data ==>Insert In To Mysql

```

## 配置文件说明

必须配置 `config.json` 以确保数据可以正确的采集和进入数据库

```json
{
  "sourceCategoryLink":"http://www.acgdoge.net/archives/category/%e5%8a%a8%e7%94%bb%e6%bc%ab%e7%94%bb", #采集的目录地址
  "pageNumber":100, .. 采集的页数
  "allTime":"2018-10-26 01:16:39", // wordpress 中发布时间
  "categoryId":711, // 采集资源进入数据库后所属于的文章分类目录 id
  "postUserId":121393, // 发布用户的 id
  "hostName":"", // 你的域名
  "mysqlHost":"", //mysql 服务器地址
  "mysqlUser":"", //数据库 用户名
  "mysqlPass":"", //数据库 密码
  "mysqlPort":"", // 数据库 端口
  "mysqlDbName":"" //数据库 名称
}

```

## 最后

请妥善使用此工具，由于正则都是写死的，所以只能爬取 acgdoge 的文章 ，速度大概是 1000篇/20s 左右

此工具对于网站压力比较大，放出的协程经常超过100个，请在网站压力低的时候使用

请尊重 acgdoge 站长发布的原创内容，爬取其内容需要按照协议标明出处


