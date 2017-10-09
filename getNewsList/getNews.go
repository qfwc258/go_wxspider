package getNewsList

import (
	"github.com/haijiandong/wxSpider/model"
	"strconv"
	"strings"
	"time"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/haijiandong/wxSpider/common"
	
	"github.com/PuerkitoBio/goquery"
)

const (
	newsTypeNumber = 19 //栏目编号，从0到19
	urlString      = "http://weixin.sogou.com/wapindex/wap/0612/wap_{0}/{1}.html"
)

func GetNewsList(types int) {
	//抓取最近的10页新闻数据
	for i := 0; i < 10; i++ {
		url := common.Format(urlString, types, i) //相当于c#中的string.format()
		getDataWhitUrl(url, types)
	}
}

func DeleteNewsList(){
	//删除5天以前的数据
	db, dbErr := sql.Open("sqlite3", "./wxnews.db")
	checkErr(dbErr)
	delList,delErr:=db.Prepare("DELETE FROM NewsList WHERE date('now', '-5 day') >= date(CreateTime)")
	checkErr(delErr)
	delInfo,delErr:=db.Prepare("DELETE FROM NewsInfo WHERE date('now', '-5 day') >= date(CreateTime)")
	checkErr(delErr)
	_,err1:= delList.Exec()
	checkErr(err1)
	_,err2:= delInfo.Exec()
	checkErr(err2)
}

func getDataWhitUrl(url string, newsType int) {
	//sql
	db, dbErr := sql.Open("sqlite3", "./wxnews.db")
	checkErr(dbErr)
	ifStop := false
	result := common.GetMethod(url)
	if result != "" {
		//解析文件获得经过处理的并且可以被path访问的Node类型的node
		doc, docErr := goquery.NewDocumentFromReader(strings.NewReader(result))
		checkErr(docErr)
		doc.Find("li").Each(func(i int,s *goquery.Selection){
			title := s.Find("h4").Text()
			summary := s.Find("p").Text()
			author:=s.Find("span.s2").Text()
			timespan:=s.Find("span.s3").AttrOr("data-lastmodified",strconv.FormatInt(time.Now().Unix(),10))
			temporaryUrl:=s.Find("div.pic a").AttrOr("href","")
			image:=s.Find("div.pic a img").AttrOr("src","")

			//查询是否插入了这条数据
			var  count int
			sqlErr :=db.QueryRow("select count(*) from NewsList where Title=? and NewsType=?;",title,newsType).Scan(&count);
			checkErr(sqlErr)
			if count>0{
				ifStop=true
			}else{
				newsInfoId:=getNewsInfo(temporaryUrl,newsType)
				if newsInfoId!=""{
					insertSql,insErr:=db.Prepare(`INSERT INTO NewsList(
							Id,CreateTime,ModifiedTime,Title,PublishTime,Description,CoverUrl,Author,NewsId,NewsType
						) 
						values(
							?,?,?,?,?,?,?,?,?,?
						)`)
					checkErr(insErr)
					defer insertSql.Close()
					_,resErr:=insertSql.Exec(common.GetGuid(),time.Now().Format("2006-01-02 15:04:05"),time.Now().Format("2006-01-02 15:04:05"),title,timespan,summary,image,author,newsInfoId,newsType)	
					checkErr(resErr)
				}
			}
		})
	}
}

func getNewsInfo(url string,newsType int)string{
	//sql
	db, dbErr := sql.Open("sqlite3", "./wxnews.db")
	checkErr(dbErr)
	result := common.GetMethod(url)
	if result != "" {
		//解析文件获得经过处理的并且可以被path访问的Node类型的node
		doc, docErr := goquery.NewDocumentFromReader(strings.NewReader(result))
		checkErr(docErr)
		content:=doc.Find("div#img-content")
		//解析具体内容
		title:=strings.TrimSpace(content.Find("h2.rich_media_title").Text())
		publishTime:=strings.TrimSpace(content.Find("em#post-date").Text())
		author:=strings.TrimSpace(content.Find("a#post-user").Text())
		officialAccountName:=strings.TrimSpace(content.Find("a#post-user").Text())
		contentDom:=strings.TrimSpace(content.Find("div.rich_media_content ").Text())
		var contentHtml string
		if contentDom!=""{
			contentHtmlRes,htmlErr:=content.Find("div.rich_media_content ").Html()
			checkErr(htmlErr)
			contentHtml=strings.TrimSpace(strings.Replace(strings.Replace(contentHtmlRes,"data-src","src",-1),"preview.html","player.html",-1))
		}
		//查询是否插入了这条数据
		var  count int
		sqlErr :=db.QueryRow("select count(*) from NewsInfo where Title=? and NewsType=?;",title,newsType).Scan(&count);
		checkErr(sqlErr)
		if count>0{
			return ""
		}else{
			insertSql,insErr:=db.Prepare(`INSERT INTO NewsInfo(
					Id,CreateTime,ModifiedTime,Title,PublishTime,Content,Author,OfficialAccountName,NewsType
				) 
				values(
					?,?,?,?,?,?,?,?,?
				)`)
			checkErr(insErr)
			defer insertSql.Close()
			id:=common.GetGuid()
			_,resErr:=insertSql.Exec(id,time.Now().Format("2006-01-02 15:04:05"),time.Now().Format("2006-01-02 15:04:05"),title,publishTime,contentHtml,author,officialAccountName,newsType)
			checkErr(resErr)
			return id
		}
	}
	return ""
}

//api handle
func GetNewsListByTypeAndIndex(types int,index int)[]model.NewsList{
	//sql
	db, dbErr := sql.Open("sqlite3", "./wxnews.db")
	checkErr(dbErr)
	skipNum := index * 10
	//初始行的偏移量是0
	if index == 0{
		skipNum = 0;
	}
	listSql,sqlErr:=db.Query(`SELECT * FROM NewsList where NewsType=? order by PublishTime desc limit 10 offset ? ;`,types,skipNum)
	checkErr(sqlErr)
	defer listSql.Close()
	news:=make([]model.NewsList,10,11)
	i:=0
	for listSql.Next(){
		lists:=model.NewsList{}
		sErr:=listSql.Scan(&lists.Id,&lists.CreateTime,&lists.Title,&lists.PublishTime,&lists.Description,&lists.CoverUrl,&lists.Author,&lists.NewsId,&lists.NewsType,&lists.ModifiedTime)
		checkErr(sErr)
		news[i]=lists
		i++
	}
	return news
}
func GetNewsInfoById(id string)model.NewsInfo{
	//sql
	db, dbErr := sql.Open("sqlite3", "./wxnews.db")
	checkErr(dbErr)
	info:=model.NewsInfo{}
	sqlErr:=db.QueryRow("SELECT * FROM NewsInfo where Id=?",id).Scan(&info.Id,&info.CreateTime,&info.ModifiedTime,&info.Title,&info.PublishTime,&info.Content,&info.Author,&info.OfficialAccountName,&info.NewsType)
	checkErr(sqlErr)
	return info
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}