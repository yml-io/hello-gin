package types

import "gorm.io/gorm"

// 基本数据类型、实现了 Scanner 和 Valuer 接口的自定义类型及其指针或别名
// TODO: 看 golang 中的常见接口
type Person struct {
	gorm.Model
	UserName  string `form:"user_name" json:"user_name" xml:"user_name"`
	NickName  string `form:"nick_name" json:"nick_name" xml:"nick_name"`
	Introduce string `form:"introduce" json:"introduce" xml:"introduce"`
	Sex       string `form:"sex" json:"sex" xml:"sex"`
	Email     string `form:"email" json:"email" xml:"email"`
	Coin      uint64 `form:"coin" json:"coin" xml:"coin" binding:"-"`
	Type      string `form:"type" json:"type" xml:"type"`
	Privacy   string `form:"privacy" json:"privacy" xml:"privacy"`
	Views     uint64 `form:"views" json:"views" xml:"views" binding:"-"`
}

type Post struct {
	gorm.Model
	Title   string `form:"title" json:"title" xml:"title"`
	Content string `form:"content" json:"content" xml:"content"`
	Auth    int    `form:"auth" json:"auth" xml:"auth" binding:"-"`
	Status  string `form:"status" json:"status" xml:"status"`
	Views   uint64 `form:"views" json:"views" xml:"views" binding:"-"`
}

type BlackList struct {
	gorm.Model
	Uid      int
	Strategy string
}

type Follow struct {
	gorm.Model
	Follower int `form:"follower" json:"follower" xml:"follower"`
	Followee int `form:"followee" json:"followee" xml:"followee"`
}

type Favorite struct {
	gorm.Model
	Uid int `form:"uid" json:"uid" xml:"uid"`
	Pid int `form:"pid" json:"pid" xml:"pid"`
}
