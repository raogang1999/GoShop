package model

const (
	LEAVING_MESSAGES = iota + 1
	COMPLAINT
	INQUIRY
	POST_SALE
	WANT_TO_BUY
)

// 留言
type LeavingMessages struct {
	BaseModel

	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100)"`

	Message string
	File    string `gorm:"type:varchar(200)"` //保存文件的url
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

type Address struct {
	BaseModel

	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

// 收藏
type UserFav struct {
	BaseModel

	User  int32 `gorm:"type:int;index:idx_user_goods,unique"` //索引，联合索引
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"` //索引，联合索引
}

func (UserFav) TableName() string {
	return "userfav"
}
