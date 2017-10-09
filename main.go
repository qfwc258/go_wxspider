package main

import (
	"encoding/json"
	"strconv"
	"time"
	"github.com/haijiandong/wxSpider/getNewsList"
	"github.com/gorilla/mux"
	"net/http"	
)
const (
	newsTypeNumber = 19 //栏目编号，从0到19
)

var ticker =time.NewTicker(time.Minute*5)
var tickerH=time.NewTicker(time.Hour*24)   //24小时删除一次冗余数据

func main() {
	go takeNews()
	//api
	router:=mux.NewRouter()
	router.HandleFunc("/GetAllList/{type}/{index}",GetAllList).Methods("GET")
	router.HandleFunc("/GetContent/{id}",GetContent).Methods("GET")
	http.ListenAndServe(":8085",router)
}

func takeNews(){
	for _=range ticker.C{
		for i:=0;i<=newsTypeNumber;i++{
			getNewsList.GetNewsList(i)
		}
	}
	//删除5天前的数据
	for _=range tickerH.C{
		getNewsList.DeleteNewsList()
	}
}

//api
func GetAllList(w http.ResponseWriter,req *http.Request){
	parames:=mux.Vars(req)
	types,err:=strconv.Atoi(parames["type"])
	checkErr(err)
	index,err:=strconv.Atoi(parames["index"])
	lists:=getNewsList.GetNewsListByTypeAndIndex(types,index)
	json.NewEncoder(w).Encode(lists)
}
func GetContent(w http.ResponseWriter,req *http.Request){
	parames:=mux.Vars(req)
	id:=parames["id"]
	info:=getNewsList.GetNewsInfoById(id)
	json.NewEncoder(w).Encode(info)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}