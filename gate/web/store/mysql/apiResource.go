package mysqlStore

import (
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
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

func (res *ApiResource) OrganizationID() int64 {
	return 0
}

func (res *ApiResource) ResourceID() int64 {
	return res.id
}

func (res *ApiResource) ResourceClass() resource2.Class {
	return resource2.Api
}

func (res *ApiResource) ResourceTitle() string {
	return res.title
}

func (res *ApiResource) ResourceDesc() string {
	return res.desc
}

func (res *ApiResource) GetChildrenResources(options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	return []model2.Resource{}, 0, nil
}
