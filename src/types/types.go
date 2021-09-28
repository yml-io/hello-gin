package types

import "gorm.io/gorm"

// 基本数据类型、实现了 Scanner 和 Valuer 接口的自定义类型及其指针或别名
// TODO: golang 中的常见接口
type Person struct {
	gorm.Model
	UserName  string
	NickName  string
	Introduce string
	Sex       string
	Email     string
	Coin      uint64
	Type      string
	Privacy   string
	Views     uint64
}

type Post struct {
	gorm.Model
	Title   string
	Content string
	Auth    uint64
	Status  string
	Views   uint64
}

type BlackList struct {
	gorm.Model
	Uid      uint64
	Strategy string
}

type Follow struct {
	gorm.Model
	Follower uint64
	Followee uint64
}

type Favorite struct {
	gorm.Model
	Uid uint64
	Pid uint64
}
