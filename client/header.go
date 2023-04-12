package client

import (
	"fmt"
	"strings"

	"github.com/wazofski/storz/store/options"
)

type restHeaderOption struct {
	Function options.OptionFunction
}

type headerOption interface {
	options.Option
	options.GetOption
	options.CreateOption
	options.UpdateOption
	options.DeleteOption
	options.ListOption
}

func Header(key string, val string) headerOption {
	return restHeaderOption{
		Function: func(options options.OptionHolder) error {
			restOpts, ok := options.(*restOptions)
			if !ok {
				// log.Printf("cannot apply REST client specific header option")
				return nil
			}
			if len(strings.Split(key, " ")) > 1 {
				return fmt.Errorf("invalid header name [%s]", key)
			}
			restOpts.Headers[key] = val
			// log.Printf("header option %s%s: [%s]", strings.ToUpper(key[:1]), key[1:], val)
			return nil
		},
	}
}

func (d restHeaderOption) ApplyFunction() options.OptionFunction {
	return d.Function
}
func (d restHeaderOption) GetCreateOption() options.Option {
	return d
}
func (d restHeaderOption) GetDeleteOption() options.Option {
	return d
}
func (d restHeaderOption) GetGetOption() options.Option {
	return d
}
func (d restHeaderOption) GetUpdateOption() options.Option {
	return d
}
func (d restHeaderOption) GetListOption() options.Option {
	return d
}
func (d restHeaderOption) GetHeaderOption() options.Option {
	return d
}
