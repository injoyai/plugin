package plugin

import "time"

type Interface interface {
	Group() string                                                                 //分组
	Version() string                                                               //版本
	Name() string                                                                  //名称
	Memo() string                                                                  //说明
	Func(name string, args map[string]interface{}) (map[string]interface{}, error) //函数
	Run() (<-chan RunInfo, error)                                                  //执行插件,推送些日志信息
	Running() bool                                                                 //是否在运行
}

type RunInfo interface {
	Type() string         //信息类型
	Time() time.Time      //信息时间
	Payload() interface{} //信息内容
}

type FileInfo struct {
	Interface
	Filename string
	Date     time.Time
}

type RpcInfo struct {
	Interface
	Address string
	Date    time.Time
}

type Manage struct {
	Code map[string]Interface //代码插件
	File map[string]*FileInfo //文件插件
	RPC  map[string]*RpcInfo  //rpc插件
}

// FileList 插件列表
func (this *Manage) FileList() []map[string]interface{} {

	return nil
}

func (this *Manage) Install(Type string, key string, f Interface) {

}

/*

 */

func NewManage() {

}
