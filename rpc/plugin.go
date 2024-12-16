package rpc

import (
	"context"
	"github.com/injoyai/conv"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/logs"
	"time"
)

type Function func(args map[string]interface{}) (*Response, error)

func NewPlugin(system, key, name, version string) *Plugin {
	return &Plugin{
		System:     system,
		Key:        key,
		Name:       name,
		Version:    version,
		Functions:  nil,
		Registered: false,
		Enable:     false,
		client:     nil,
	}
}

type Plugin struct {
	System     string              `json:"system"`     //系统
	Key        string              `json:"key"`        //插件key
	Name       string              `json:"name"`       //名称
	Version    string              `json:"version"`    //版本
	Functions  map[string]Function `json:"functions"`  //插件函数
	Registered bool                `json:"registered"` //是否注册
	Enable     bool                `json:"enable"`     //是否启用

	//
	client *io.ClientManage //客户端

}

// Register 函数注册
func (this *Plugin) Register(method string, f Function) {
	this.Functions[method] = f
}

func (this *Plugin) Start(ctx context.Context, cfg *DiscoverConfig) {
	for this.client.GetClientLen() == 0 {
		<-time.After(conv.SelectDuration(cfg.Interval >= time.Second, cfg.Interval, 5*time.Second))
		//广播发现包
		serverAddress := cfg.broadcast(ctx, this.System, this.Key)
		for _, v := range serverAddress {
			logs.Trace("发现服务: ", v.Address)
			//连接服务
			if this.client.GetClient(v.Address) == nil {
				logs.Trace("连接服务: ", v.Address)
				c, err := this.client.DialClient(dial.WithTCP(v.Address))
				if err != nil {
					logs.Err(err)
					continue
				}
				logs.Trace("注册插件到服务")
				c.Redial(func(c *io.Client) {
					c.Debug(false)
					c.SetReadWriteWithPkg()
					c.WriteRead(conv.Bytes(&Register{
						Key:       this.Key,
						Name:      this.Name,
						Version:   this.Version,
						Functions: this.Functions,
					}))
				})
			}
		}
	}
}
