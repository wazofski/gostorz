package options

import (
	"errors"
	"log"
)

type Option interface {
	ApplyFunction() OptionFunction
}

type CreateOption interface {
	Option
	GetCreateOption() Option
}

type DeleteOption interface {
	Option
	GetDeleteOption() Option
}

type GetOption interface {
	Option
	GetGetOption() Option
}

type UpdateOption interface {
	Option
	GetUpdateOption() Option
}

type ListOption interface {
	Option
	GetListOption() Option
}

type OptionHolder interface {
	CommonOptions() *CommonOptionHolder
}

type OptionFunction func(OptionHolder) error

type PropFilterSetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyFilterSetting []string

type CommonOptionHolder struct {
	PropFilter       *PropFilterSetting
	KeyFilter        *KeyFilterSetting
	OrderBy          string
	OrderIncremental bool
	PageSize         int
	PageOffset       int
}

func (d *CommonOptionHolder) CommonOptions() *CommonOptionHolder {
	return d
}

func CommonOptionHolderFactory() CommonOptionHolder {
	return CommonOptionHolder{
		PropFilter:       nil,
		KeyFilter:        nil,
		OrderBy:          "",
		OrderIncremental: true,
		PageSize:         0,
		PageOffset:       0,
	}
}

func PropFilter(prop string, val string) ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.PropFilter != nil {
				return errors.New("prop filter option already set")
			}

			commonOptions.PropFilter = &PropFilterSetting{
				Key:   prop,
				Value: val,
			}

			// opstr, _ := json.Marshal(*commonOptions.Filter)
			// log.Printf("filter option %s", string(opstr))
			return nil
		},
	}
}

func KeyFilter(keys ...string) ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			if len(keys) == 0 {
				log.Printf("ignoring empty key filter")
				return nil
			}

			commonOptions := options.CommonOptions()
			if commonOptions.KeyFilter != nil {
				return errors.New("key filter option already set")
			}

			commonOptions.KeyFilter = (*KeyFilterSetting)(&keys)

			// opstr, _ := json.Marshal(*commonOptions.Filter)
			// log.Printf("filter option %s", string(opstr))
			return nil
		},
	}
}

func PageSize(ps int) ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.PageSize > 0 {
				return errors.New("page size option has already been set")
			}
			commonOptions.PageSize = ps
			// log.Printf("pagination size option %d", ps)
			return nil
		},
	}
}

func PageOffset(po int) ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.PageOffset > 0 {
				return errors.New("page offset option has already been set")
			}
			if commonOptions.PageOffset < 0 {
				return errors.New("page offset cannot be negative")
			}
			commonOptions.PageOffset = po
			// log.Printf("pagination offset option %d", po)
			return nil
		},
	}
}

func OrderBy(field string) ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			commonOptions := options.CommonOptions()
			if len(commonOptions.OrderBy) > 0 {
				return errors.New("order by option has already been set")
			}
			commonOptions.OrderBy = field
			// log.Printf("order by option: %s", field)
			return nil
		},
	}
}

func OrderDescending() ListOption {
	return listOption{
		Function: func(options OptionHolder) error {
			commonOptions := options.CommonOptions()
			if !commonOptions.OrderIncremental {
				return errors.New("order incremental option has already been set")
			}
			commonOptions.OrderIncremental = false
			// log.Printf("order incremental option %v", val)
			return nil
		},
	}
}

type listOption struct {
	Function OptionFunction
}

func (d listOption) GetListOption() Option {
	return d
}

func (d listOption) ApplyFunction() OptionFunction {
	return d.Function
}
