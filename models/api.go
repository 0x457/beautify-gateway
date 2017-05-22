package models

import "regexp"

//Param postion  参数位置
const (
	ParamPosHead     = iota //head
	ParamPosQuery           //url query
	ParamPosPath            //url path
	ParamPosBodyForm        //form表单数据
)

// param type 参数类型
const (
	ParamTypeInt64 = iota //int
	paramTypeFloat64
	paramTypeString
	paramTypeBool
)

//RequestParam  请求入参
type RequestParam struct {
	ID          int64          `xorm:"pk bigint autoincr 'id'"`
	AID         int64          `xorm:" bigint  'aid'"`       //api id
	Name        string         `xorm:" varchar(20)  'name'"` //参数名称
	TransName   string         `xorm:" varchar(20)  'transname'"`
	Type        int            `xorm:" int  'type'"`          //参数类型
	Must        bool           `xorm:" bool  'must'"`         //是否必须
	Position    int            `xorm:" int  'position'"`      //参数位置
	Description string         `xorm:" int  'description'"`   //参数描述
	Regx        string         `xorm:" varchar(100)  'regx'"` //正则信息
	Regexp      *regexp.Regexp `xorm:"-"`                     //正则表达式，用于参数验证
}

//API  api node
type API struct {
	ID int64 `xorm:"pk bigint autoincr 'id'"`

	ClusterID int64    `xorm:" bigint  'cluster_id'"`
	Cluster   *Cluster `xorm:"-"`

	Name     string `xorm:" varchar(20)  'name'"`
	AttrName string `xorm:" varchar(20)  'attr_name'"`
	Method   string `xorm:" varchar(20)  'method'"`

	Params []RequestParam `xorm:"-"`

	Path string `xorm:" varchar(50)  'path'"` //请求路径

	BodyNotForm         bool   //body数据是否为form，如为json等则为true
	PostBodyDescription string // body描述
}
