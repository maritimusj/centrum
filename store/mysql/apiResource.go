package mysqlStore

import (
	"github.com/maritimusj/centrum/resource"
)

type ApiResource struct {
	id    int64
	name  string
	title string
	desc  string

	store *mysqlStore
}

func NewApiResource(store *mysqlStore, id int64) *ApiResource {
	return &ApiResource{
		id:    id,
		store: store,
	}
}

func (res *ApiResource) GetID() int64 {
	return res.id
}

func (res *ApiResource) Title() string {
	return res.title
}

func (res *ApiResource) Desc() string {
	return res.desc
}

func (res *ApiResource) ResourceUID() string {
	return res.name
}

func (res *ApiResource) ResourceClass() resource.Class {
	return resource.Api
}

func (res *ApiResource) ResourceTitle() string {
	return res.title
}

func (res *ApiResource) ResourceDesc() string {
	return res.desc
}
