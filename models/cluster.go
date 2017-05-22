package models

var (
	//Clusters  集群集合
	Clusters = make(map[int64]*Cluster)
)

//Cluster  服务集群
type Cluster struct {
	ID           int64     `xorm:"pk bigint autoincr 'id'"`
	Name         string    `xorm:" varchar(20)  'name'"`
	livedServers []*Server `xorm:"-"` //当前存活的server
	Servers      []*Server `xorm:"-"` //所有的server

	AccessID int64         `xorm:" bigint  'access_id'"`
	Access   *AcessControl `xorm:"-"` //黑白名单控制

	LbName string
}
