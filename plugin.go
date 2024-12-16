package plugin

import (
	"context"
	"errors"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
)

var _ Interface = (*Demo)(nil)

func NewDemo() Interface {
	return &Demo{runner: safe.NewRunner(func(ctx context.Context) error {

		return nil
	})}
}

type Demo struct {
	runner  *safe.Runner
	natures map[string]interface{}
}

func (this *Demo) Group() string {
	return "server"
}

func (this *Demo) Version() string {
	return "v1.0.0"
}

func (this *Demo) Name() string {
	return "演示"
}

func (this *Demo) Memo() string {
	return "演示插件"
}

func (this *Demo) Func(name string, args map[string]interface{}) (map[string]interface{}, error) {
	switch name {
	case "get":
		key := conv.String(args["key"])
		val := this.natures[key]
		return map[string]interface{}{"value": val}, nil
	}
	return nil, errors.New("未知方法")
}

func (this *Demo) Run() (<-chan RunInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (this *Demo) Running() bool {
	return this.runner.Running()
}
