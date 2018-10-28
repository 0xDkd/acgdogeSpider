package main

import (
	"acgdogeSpider/ReadJson"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//爬虫总协程的通道，每个爬虫一个协程
var spiderC = make(chan int, 0)

//每个爬虫，每个页面的通道
var pageC = make(chan bool, 0)

//每次内容处理的通道
var contentC = make(chan bool, 0)

//每次储存的内容
var storeC = make(chan bool, 0)

//只连接一次数据库
// var db *sql.DB
// var dbErr error

func main() {
	//开始获取配置
	jsonParse := ReadJson.NewJsonStruct()
	v := ReadJson.Config{}
	jsonParse.Load("./config.json", &v)
	fmt.Println("采集开始:==>")
	link_str := v.MysqlUser + ":" + v.MysqlPass + "@" + "tcp(" + v.MysqlHost + ":" + v.MysqlPort + ")/" + v.MysqlDbName
	//仅仅只连接一次数据库
	fmt.Println(link_str)
	db, err := sql.Open("mysql", link_str)
	checkErr(err)
	for i := 0; i < v.PageNumber; i++ {
		go spider(v.SourceCategoryLink, i, spiderC, db, v)
	}
	for i := 0; i < v.PageNumber; i++ {
		fmt.Println("Spider Number:", <-spiderC)
	}

}

//爬虫协程，可以确定爬取的 url
func spider(url string, id int, c chan<- int, db *sql.DB, v ReadJson.Config) {
	durl := url + "/page/" + strconv.Itoa(id)
	//获取到页面的内容
	res, err := httpGet(durl)
	checkErr(err)
	pat := `<h2 class="post_h"><a href="(?s:(.*?))"`
	reGetUrl := regexp.MustCompile(pat)
	checkRegex(reGetUrl)
	//获取文章的 url
	postUrls := reGetUrl.FindAllStringSubmatch(res, -1)
	//开协程爬取文章
	for _, data := range postUrls {
		go getPostContent(data[1], pageC, db, v)
	}
	//接受通道消息
	len := len(postUrls)
	for i := 0; i < len; i++ {
		fmt.Println(<-pageC)
	}
	c <- id
}

//开协程爬取文章内容
func getPostContent(url string, end chan<- bool, db *sql.DB, v ReadJson.Config) {
	res, err := httpGet(url)
	checkErr(err)
	//提取有用的部分，并且全部放到数组中，可以直接并发处理
	pat := `<p style="text-align: center;"><a href="(?s:(.*?))" title="(?s:(.*?))" ><img src="http:\/\/www\.acgdoge\.net\/wp-content\/plugins\/lazy-load\/images\/1x1.trans.gif" data-lazy-src="(?s:(.*?))" alt="(?s:(.*?))"><noscript><img src="(?s:(.*?))" alt="(?s:(.*?))" \/><\/noscript><\/a><\/p>([\s\S]*?)<\/div>`
	re := regexp.MustCompile(pat)
	//content 就已经是有用信息的集合数组
	content := re.FindAllStringSubmatch(res, -1)
	//7->文章内容
	//3-> 头图
	//4->标题
	for _, data := range content {
		//并发处理需要的内容
		go stringFmt(data, contentC)
		go storeContent(data, storeC, db, v)
	}
	for i, _ := range content {
		fmt.Println("Get ContentC:", <-contentC)
		fmt.Println("Get storeC:", <-storeC)
		i = i + 1 - 1
	}
	end <- true
}

//Http 内容下载,返回 result 是页面的内容
func httpGet(url string) (result string, err error) {
	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()

	buf := make([]byte, 10*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

//内容处理:过滤掉某些不必要的内容，主要是正则循环，考虑并发
func stringFmt(data []string, end chan<- bool) {
	pat := `(?m)<p style="text-align: center;"><img src="http:\/\/www\.acgdoge\.net/wp-content/plugins/lazy-load\/images\/1x1\.trans\.gif" data-lazy-src="[\s\S]*?" alt="[\s\S]*?"><noscript>`
	patRemove := `<\/noscript><\/a><\/p>`                  //移除 noscript 部分
	patRemove2 := `<p><span id="[\s\S]*?"></span></p>`     //移除部分 span
	patRemove3 := `<h2 class="post_h_quote">[\s\S]*?</h2>` //移除最后的耳机标题
	c_re := regexp.MustCompile(pat)
	c_re2 := regexp.MustCompile(patRemove)
	c_re3 := regexp.MustCompile(patRemove2)
	c_re4 := regexp.MustCompile(patRemove3)
	data[4] = strings.Replace(data[4], "/t", "", -1) //移除 标题中的 table 键
	data[7] = c_re.ReplaceAllString(data[7], "")
	data[7] = c_re2.ReplaceAllString(data[7], "")
	data[7] = c_re3.ReplaceAllString(data[7], "")
	data[7] = c_re4.ReplaceAllString(data[7], "")
	//fmt.Println(data[7])
	end <- true
}

//数据库模型生成，采取单例，不可以并发链接但可以并发操作
func storeContent(data []string, end chan<- bool, db *sql.DB, v ReadJson.Config) {
	// 插入数据
	sql := `INSERT  wp_posts SET  post_modified_gmt=?,post_date_gmt=?,post_author=?,post_date=?,post_content=?,post_title=?,post_excerpt=?,post_status=?,comment_status=?,ping_status=?,post_name=?,to_ping=?,pinged=?,post_modified=?,post_content_filtered=?,post_parent=?,menu_order=?,post_type=?,comment_count=?,post_mime_type=?,guid=?`
	stmt, err := db.Prepare(sql)
	checkErr(err)

	res, err := stmt.Exec(v.AllTime, v.AllTime, v.PostUserId, v.AllTime, data[7], data[4], "", "publish", "open", "close", "post_name", "", "", v.AllTime, "", "0", "0", "post", "0", "", "")
	if err != nil {
		fmt.Println(err)
		fmt.Println(data[7])
		return
	}
	//checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	sql = `INSERT wp_term_relationships SET object_id=?,term_taxonomy_id=?,term_order=?`
	stmt, err = db.Prepare(sql)
	checkErr(err)
	res, err = stmt.Exec(id, v.CategoryId, 0)

	// 更新数据
	stmt, err = db.Prepare("update wp_posts set guid=? where ID=?")
	checkErr(err)

	res, err = stmt.Exec(v.HostName+"/"+strconv.FormatInt(id, 10)+".html", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	end <- true
}

//错误检查输出
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func checkRegex(cRegex *regexp.Regexp) {
	if cRegex == nil {
		fmt.Println("Can not complied regex rule", cRegex)
		return
	}

}
