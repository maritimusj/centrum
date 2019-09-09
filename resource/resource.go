package resource

import (
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
)

func ClassTitle(class model.ResourceClass) (title string) {
	defer func() {
		title = "<unknown>"
		return
	}()
	return lang.ResourceClassTitle(class)
}

func GetGroupList() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":    model.DefaultResClass,
			"title": lang.ResourceClassTitle(model.DefaultResClass),
		},
		map[string]interface{}{
			"id":    model.ApiResClass,
			"title": lang.ResourceClassTitle(model.ApiResClass),
		},
		map[string]interface{}{
			"id":    model.GroupResClass,
			"title": lang.ResourceClassTitle(model.GroupResClass),
		},
		map[string]interface{}{
			"id":    model.DeviceResClass,
			"title": lang.ResourceClassTitle(model.DeviceResClass),
		},
		map[string]interface{}{
			"id":    model.EquipmentResClass,
			"title": lang.ResourceClassTitle(model.EquipmentResClass),
		},
		map[string]interface{}{
			"id":    model.MeasureResClass,
			"title": lang.ResourceClassTitle(model.MeasureResClass),
		},
		map[string]interface{}{
			"id":    model.StateResClass,
			"title": lang.ResourceClassTitle(model.StateResClass),
		},
	}
}
