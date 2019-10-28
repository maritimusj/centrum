package app

import (
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/edge"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func initEvent() error {
	eventsMap := map[string]interface{}{
		event.ApiServerStarted: eventApiServerStarted,
		event.UserCreated:      eventUserCreated,
		event.UserUpdated:      eventUserUpdated,
		event.UserDeleted:      eventUserDeleted,

		event.DeviceCreated: eventDeviceCreated,
		event.DeviceUpdated: eventDeviceUpdated,
		event.DeviceDeleted: eventDeviceDeleted,

		event.EquipmentCreated: eventEquipmentCreated,
		event.EquipmentUpdated: eventEquipmentUpdated,
		event.EquipmentDeleted: eventEquipmentDeleted,
	}

	for e, fn := range eventsMap {
		err := Event.SubscribeAsync(e, fn, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func eventApiServerStarted() {
	err := BootAllDevices()
	if err != nil {
		log.Error("[BootAllDevices] ", err)
	}
}

func eventUserCreated(userID int64, newUserID int64) {
	adminUser, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventUserCreated: ", err)
		return
	}
	user, err := Store().GetUser(newUserID)
	if err != nil {
		log.Error("eventUserCreated: ", err)
		return
	}

	log.Info(lang.Str(lang.AdminCreateUserOk, adminUser.Name(), user.Title()))
	adminUser.Logger().Info(lang.Str(lang.AdminCreateUserOk, adminUser.Name(), user.Title()))
	user.Logger().Info(lang.Str(lang.AdminCreateUserOk, adminUser.Name(), user.Title()))
}

func eventUserUpdated(userID int64, newUserID int64) {
	adminUser, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventUserUpdated: ", err)
		return
	}
	user, err := Store().GetUser(newUserID)
	if err != nil {
		log.Error("eventUserUpdated: ", err)
		return
	}

	log.Info(lang.Str(lang.AdminUpdateUserOk, adminUser.Name(), user.Title()))
	adminUser.Logger().Info(lang.Str(lang.AdminUpdateUserOk, adminUser.Name(), user.Title()))
	user.Logger().Info(lang.Str(lang.AdminUpdateUserOk, adminUser.Name(), user.Title()))
}

func eventUserDeleted(userID int64, name string) {
	adminUser, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventUserDeleted: ", err)
		return
	}

	log.Warn(lang.Str(lang.AdminDeleteUserOk, adminUser.Name(), name))
	adminUser.Logger().Warn(lang.Str(lang.AdminDeleteUserOk, adminUser.Name(), name))
}

func eventDeviceCreated(userID int64, deviceID int64) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventDeviceCreated: ", err)
		return
	}
	device, err := Store().GetDevice(deviceID)
	if err != nil {
		log.Error("eventDeviceCreated: ", err)
		return
	}

	err = edge.ActiveDevice(device)
	if err != nil {
		log.Error("eventDeviceCreated: active device: ", err)
	}

	log.Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
	user.Logger().Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
	device.Logger().Info(lang.Str(lang.UserCreateDeviceOk, user.Name(), device.Title()))
}

func eventDeviceUpdated(userID int64, deviceID int64) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventDeviceUpdated: ", err)
		return
	}
	device, err := Store().GetDevice(deviceID)
	if err != nil {
		log.Error("eventDeviceUpdated: ", err)
		return
	}

	err = edge.ActiveDevice(device)
	if err != nil {
		log.Error("eventDeviceUpdated: active device: ", err)
	}

	log.Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
	user.Logger().Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
	device.Logger().Info(lang.Str(lang.UserUpdateDeviceOk, user.Name(), device.Title()))
}

func eventDeviceDeleted(userID int64, id int64, uid string, title string) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventDeviceDeleted: ", uid, err)
		return
	}

	edge.Remove(strconv.FormatInt(id, 10))

	log.Warn(lang.Str(lang.UserDeleteDeviceOk, user.Name(), title))
	user.Logger().Warn(lang.Str(lang.UserDeleteDeviceOk, user.Name(), title))
}

func eventEquipmentCreated(userID int64, equipmentID int64) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventEquipmentCreated: ", err)
		return
	}
	equipment, err := Store().GetEquipment(equipmentID)
	if err != nil {
		log.Error("eventEquipmentCreated: ", err)
		return
	}

	log.Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
	user.Logger().Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
	equipment.Logger().Info(lang.Str(lang.UserCreateEquipmentOk, user.Name(), equipment.Title()))
}

func eventEquipmentUpdated(userID int64, equipmentID int64) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventEquipmentUpdated: ", err)
		return
	}
	equipment, err := Store().GetEquipment(equipmentID)
	if err != nil {
		log.Error("eventEquipmentUpdated: ", err)
		return
	}

	log.Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
	user.Logger().Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
	equipment.Logger().Info(lang.Str(lang.UserUpdateEquipmentOk, user.Name(), equipment.Title()))
}

func eventEquipmentDeleted(userID int64, title string) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventEquipmentDeleted: ", err)
		return
	}

	log.Warn(lang.Str(lang.UserDeleteEquipmentOk, user.Name(), title))
	user.Logger().Warn(lang.Str(lang.UserDeleteEquipmentOk, user.Name(), title))
}
