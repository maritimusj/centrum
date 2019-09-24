package app

import (
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/model"
)
import log "github.com/sirupsen/logrus"

func eventProcessor() {
	ch := event.Sub(Ctx,
		event.User,
		event.Device,
		event.Equipment)

	for {
		select {
		case <-Ctx.Done():
			return
		case data := <-ch:
			if data != nil {
				processLog(data)
			}
		}
	}
}

func processLog(data *event.Data) {
	switch data.Code {
	case event.User:
		processUserLog(data.Values)
	case event.Device:
		processDeviceLog(data.Values)
	case event.Equipment:
		processEquipmentLog(data.Values)
	}
}

func processDeviceLog(data map[string]interface{}) {
	var err error
	var user model.User
	//data.Clone后，由于json.Marshal的原因，所有int64都变成了float64
	if userID, ok := data["userID"].(float64); ok {
		user, err = Store().GetUser(int64(userID))
		if err != nil {
			log.Error("processDeviceLog: ", err)
		}
	}
	var device model.Device
	if deviceID, ok := data["deviceID"].(float64); ok {
		device, err = Store().GetDevice(int64(deviceID))
		if err != nil {
			log.Error("processDeviceLog: ", err)
		}
	}

	if v, ok := data["event"].(string); ok {
		switch v {
		case event.Created:
			log.Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
			user.Logger().Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
			device.Logger().Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
		case event.Updated:
			log.Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
			user.Logger().Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
			device.Logger().Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
		case event.Deleted:
			title, _ := data["title"].(string)
			log.Warn(lang.Str(lang.UserDeleteDeviceOk, user.Name(), title))
			user.Logger().Warn(lang.Str(lang.UserDeleteDeviceOk, user.Name(), title))
		}
	}
}

func processUserLog(data map[string]interface{}) {
	var err error
	var admin model.User
	//data.Clone后，由于json.Marshal的原因，所有int64都变成了float64
	if userID, ok := data["adminID"].(float64); ok {
		admin, err = Store().GetUser(int64(userID))
		if err != nil {
			log.Error("processUserLog: ", err)
		}
	}
	var user model.User
	if userID, ok := data["userID"].(float64); ok {
		user, err = Store().GetUser(int64(userID))
		if err != nil {
			log.Error("processUserLog: ", err)
		}
	}

	if v, ok := data["event"].(string); ok {
		switch v {
		case event.Created:
			log.Info(lang.Str(lang.AdminCreateUserOk, admin.Name(), user.Title()))
			admin.Logger().Info(lang.Str(lang.AdminCreateUserOk, admin.Name(), user.Title()))
			user.Logger().Info(lang.Str(lang.AdminCreateUserOk, admin.Name(), user.Title()))
		case event.Updated:
			log.Info(lang.Str(lang.AdminUpdateUserOk, admin.Name(), user.Title()))
			admin.Logger().Info(lang.Str(lang.AdminUpdateUserOk, admin.Name(), user.Title()))
			user.Logger().Info(lang.Str(lang.AdminUpdateUserOk, admin.Name(), user.Title()))
		case event.Deleted:
			name, _ := data["name"].(string)
			log.Warn(lang.Str(lang.AdminDeleteUserOk, admin.Name(), name))
			admin.Logger().Warn(lang.Str(lang.AdminDeleteUserOk, admin.Name(), name))
		}
	}
}

func processEquipmentLog(data map[string]interface{}) {
	var err error
	var user model.User
	//data.Clone后，由于json.Marshal的原因，所有int64都变成了float64
	if userID, ok := data["userID"].(float64); ok {
		user, err = Store().GetUser(int64(userID))
		if err != nil {
			log.Error("processEquipmentLog: ", err)
		}
	}
	var equipment model.Equipment
	if equipmentID, ok := data["equipmentID"].(float64); ok {
		equipment, err = Store().GetEquipment(int64(equipmentID))
		if err != nil {
			log.Error("processEquipmentLog: ", err)
		}
	}

	if v, ok := data["event"].(string); ok {
		switch v {
		case event.Created:
			log.Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
			user.Logger().Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
			equipment.Logger().Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
		case event.Updated:
			log.Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
			user.Logger().Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
			equipment.Logger().Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
		case event.Deleted:
			title, _ := data["title"].(string)
			log.Warn(lang.Str(lang.UserDeleteEquipmentOk, user.Name(), title))
			user.Logger().Warn(lang.Str(lang.UserDeleteEquipmentOk, user.Name(), title))
		}
	}
}
