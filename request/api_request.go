package request

import (
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
)

type Api struct {
	resource resource.Resource
	action   resource.Action
}

func (r *Api) Resource() resource.Resource {
	return r.resource
}

func (r *Api) Action() resource.Action {
	return r.action
}

func NewApiRequest(store store.Store, routerName string, method string) (model.Request, error) {
	res, err := store.GetApiResource(routerName, method)
	if err != nil {
		return nil, err
	}

	var action resource.Action
	switch method {
	case "GET":
		action = resource.View
	default:
		action = resource.Ctrl
	}

	return &Api{
		resource: res,
		action:   action,
	}, nil
}
