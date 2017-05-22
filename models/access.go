package models

// 名单类型
const (
	ACTypeBlack = iota
	ACTypeWhite
)

//AcessControl  黑白名单控制
type AcessControl struct {
	ID     int64  `xorm:"pk bigint autoincr 'id'"`
	Type   int    `xorm:"int  'type'"`            //类型
	ACList string `xorm:"varchar(200)  'aclist'"` //名单，按逗号分割
}
