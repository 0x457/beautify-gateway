package models

//Server  后端服务
//TODO:阀值控制、后端check等
type Server struct {
	ID   int64  `xorm:"pk bigint autoincr 'id'"`
	Addr string `xorm:" varchar(20)  'addr'"` //请求地址

	//health check
	CheckURL      string `xorm:" varchar(20)  'check_url'"`
	CheckTimeOut  int64  `xorm:" varchar(20)  'check_timeout'"`  //second
	CheckDuration int64  `xorm:" varchar(20)  'check_duration'"` //检查间隔

	//traffic rate 检测
	//rate limit 流控

}
