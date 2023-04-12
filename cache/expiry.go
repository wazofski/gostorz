package cache

import (
	"fmt"
	"time"

	"github.com/wazofski/storz/store/options"
)

type cacheExpiryOption struct {
	Function options.OptionFunction
}

type expiryOption interface {
	options.Option
	// options.GetOption
	options.CreateOption
	options.UpdateOption
}

func Expire(duration time.Duration) expiryOption {
	return &cacheExpiryOption{
		Function: func(options options.OptionHolder) error {
			cacheOpts, ok := options.(*cacheOptions)
			if !ok {
				// log.Printf("cannot apply Cache specific option")
				return nil
			}

			if duration < 0 {
				return fmt.Errorf("invalid duration [%d]", duration)
			}

			cacheOpts.Expiration = duration
			return nil
		},
	}
}

func (d *cacheExpiryOption) ApplyFunction() options.OptionFunction {
	return d.Function
}
func (d *cacheExpiryOption) GetCreateOption() options.Option {
	return d
}
func (d *cacheExpiryOption) GetUpdateOption() options.Option {
	return d
}

// func (d *cacheExpiryOption) GetGetOption() options.Option {
// 	return d
// }
