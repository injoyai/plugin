package plugin

import (
	"fmt"
	"time"
)

//type Manage struct {
//	Code map[string]Interface //代码插件
//	File map[string]*FileInfo //文件插件
//	RPC  map[string]*RpcInfo  //网络插件
//}

type Printer interface {
	Plugin
	Print()
}

/*
Plugin
Install 就是下载文件
Uninstall 就是删除文件
*/
type Plugin interface {
	Info() Info     //插件信息
	Enable() error  //启用
	Disable() error //禁用
	Closed() bool   //是否关闭
	Err() error     //错误信息

}

type Plugin2 interface {
	Plugin
	Log() <-chan Log //日志信息
	Status() Status  //状态信息
}

// Log 日志信息
type Log struct {
	Time    time.Time `json:"time"`
	Level   int       `json:"level"`
	Tag     string    `json:"tag"`
	Message string    `json:"message"`
}

func (this Log) String() string {
	if len(this.Tag) == 0 {
		return fmt.Sprintf("%s %s", this.Time.Format("2006-01-02 15:04:05"), this.Message)
	}
	return fmt.Sprintf("%s [%s] %s", this.Time.Format("2006-01-02 15:04:05"), this.Tag, this.Message)
}

// Status 状态信息
type Status struct {
	Status  int    `json:"status"`  //状态枚举 0:未启用 1:启用
	Message string `json:"message"` //状态消息
}

// Info 插件信息
type Info struct {
	Name    string `json:"name"`    //插件名称
	Memo    string `json:"memo"`    //说明
	Author  string `json:"author"`  //作者
	Version string `json:"version"` //版本
	Website string `json:"website"` //官网
	License string `json:"license"` //许可证
}
