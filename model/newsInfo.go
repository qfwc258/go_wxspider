package model

import(
)

type NewsInfo struct{
	Id string
	CreateTime string
	ModifiedTime  string   //修改时间
	Title string   //文章标题
	PublishTime string  //发布时间
	Content  string //内容，存储为HTML文本格式
	Author string // 文章作者
	OfficialAccountName  string   //公众号名称
	NewsType int   //文章类型
}