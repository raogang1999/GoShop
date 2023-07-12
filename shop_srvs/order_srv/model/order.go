package model

import (
	"time"
)

// 购物车
type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index;comment:用户ID"`
	Goods   int32 `gorm:"type:int;index;comment:商品ID"`
	Num     int32 `gorm:"type:int;comment:商品数量"`
	Checked bool
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

// 订单信息
type OrderInfo struct {
	BaseModel
	User    int32  `gorm:"type:int;index;comment:用户ID"`
	OrderSn string `gorm:"type:varchar(64);index;comment:订单号"`
	PayType string `gorm:"type:varchar(20);comment:支付类型,支付宝,微信"`

	//状态
	Status     string `gorm:"type:varchar(20);comment:订单状态，PAYING待支付,TRADE_SUCCESS支付成功,TRADE_FINISHED交易结束,TRADE_CLOSED超时交易关闭，WAIT_BUYER_PAY交易创建"`
	TradeNo    string `gorm:"type:varchar(64);comment: 交易号"` //支付宝的交易号
	OrderMount float32
	PayTime    *time.Time

	//收货人信息
	Address      string `gorm:"type:varchar(255);comment:收货地址"`
	SignerName   string `gorm:"type:varchar(255);comment:收货人姓名"`
	SignerMobile string `gorm:"type:varchar(255);comment:收货人电话"`
	Post         string `gorm:"type:varchar(20);comment:邮编"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

// 订单商品
type OrderGoods struct {
	BaseModel

	Order int32 `gorm:"type:int;index;comment:订单ID"`
	Goods int32 `gorm:"type:int;index;comment:商品ID"`

	//这个商品的快照信息，字段冗余，不符合三范式，但是可以减少表的关联查询，减少跨服务调用
	GoodsName  string  `gorm:"type:varchar(255);comment:商品名称"`
	GoodsImage string  `gorm:"type:varchar(255);comment:商品图片"`
	GoodsPrice float32 `gorm:"type:float;comment:商品价格"`
	Num        int32   `gorm:"type:int;comment:商品数量"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
