package model

import(
)

const(
	热门=iota
	推荐
	段子手
	养生堂
	私房话
	八卦精
	爱生活
	财经迷
	汽车迷
	科技咖
	潮人帮
	辣妈帮
	点赞党
	旅行家
	职场人
	美食家
	古今通
	学霸族
	星座控
	体育迷
)
type NewsList struct{
	Id string
	CreateTime string  //创建时间
	ModifiedTime  string   //修改时间
	Title string  //标题
	PublishTime string  //发布时间(时间戳)
	Description string  //文章描述
	CoverUrl string  //封面图片链接
	Author string  //文章作者
	NewsId string  //文章内容本地存储id
	NewsType int   //文章类型
}