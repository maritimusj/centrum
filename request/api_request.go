package request

import (
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/store"
)

type Api struct {
	resource model.Resource
	action   model.Action
}

func (r *Api) Resource() model.Resource {
	return r.resource
}

func (r *Api) Action() model.Action {
	return r.action
}

func NewApiRequest(store store.Store, routerName string, method string) (model.Request, error) {
	resource, err := store.GetApiResource(routerName, method)
	if err != nil {
		return nil, err
	}

	var action model.Action
	switch method {
	case "GET":
		action = model.View
	default:
		action = model.Ctrl
	}

	return &Api{
		resource: resource,
		action:   action,
	}, nil
}
