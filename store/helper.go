package store

import "github.com/maritimusj/centrum/model"

type Option struct {
	Limit   int64
	Offset  int64
	Kind    model.MeasureKind
	Class   model.ResourceClass
	Group   int64
	Keyword string
}

type OptionFN func(*Option)

func Page(page, pageSize int64) OptionFN {
	return func(i *Option) {
		i.Offset = (page - 1) * pageSize
		i.Limit = pageSize
	}
}

func Limit(limit int64) OptionFN {
	return func(i *Option) {
		i.Limit = limit
	}
}

func Offset(offset int64) OptionFN {
	return func(i *Option) {
		i.Offset = offset
	}
}

func Kind(kind model.MeasureKind) OptionFN {
	return func(i *Option) {
		i.Kind = kind
	}
}

func Class(class model.ResourceClass) OptionFN {
	return func(i *Option) {
		i.Class = class
	}
}

func Keyword(keyword string) OptionFN {
	return func(i *Option) {
		i.Keyword = keyword
	}
}
