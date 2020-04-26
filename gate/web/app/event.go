package app

import (
	"strconv"

	"github.com/maritimusj/centrum/gate/event"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore"
	"github.com/maritimusj/centrum/gate/web/edge"
	log "github.com/sirupsen/logrus"
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
	BootAllDevices()
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

	log.WithField("src", logStore.SystemLog).Info(lang.AdminCreateUserOk.Str(adminUser.Name(), user.Title()))
	adminUser.Logger().Info(lang.AdminCreateUserOk.Str(adminUser.Name(), user.Title()))
	user.Logger().Info(lang.AdminCreateUserOk.Str(adminUser.Name(), user.Title()))
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

	log.WithField("src", logStore.SystemLog).Info(lang.AdminUpdateUserOk.Str(adminUser.Name(), user.Title()))
	adminUser.Logger().Info(lang.AdminUpdateUserOk.Str(adminUser.Name(), user.Title()))
	user.Logger().Info(lang.AdminUpdateUserOk.Str(adminUser.Name(), user.Title()))
}

func eventUserDeleted(userID int64, name string) {
	adminUser, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventUserDeleted: ", err)
		return
	}

	log.WithField("src", logStore.SystemLog).Warn(lang.AdminDeleteUserOk.Str(adminUser.Name(), name))
	adminUser.Logger().Warn(lang.AdminDeleteUserOk.Str(adminUser.Name(), name))
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

	err = edge.ActiveDevice(device, Config)
	if err != nil {
		log.Error("eventDeviceCreated: active device: ", err)
	}

	log.WithField("src", logStore.SystemLog).Info(lang.UserCreateDeviceOk.Str(user.Name(), device.Title()))
	user.Logger().Info(lang.UserCreateDeviceOk.Str(user.Name(), device.Title()))
	device.Logger().Info(lang.UserCreateDeviceOk.Str(user.Name(), device.Title()))
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

	err = edge.ActiveDevice(device, Config)
	if err != nil {
		log.Error("eventDeviceUpdated: active device: ", err)
	}

	log.WithField("src", logStore.SystemLog).Info(lang.UserUpdateDeviceOk.Str(user.Name(), device.Title()))
	user.Logger().Info(lang.UserUpdateDeviceOk.Str(user.Name(), device.Title()))
	device.Logger().Info(lang.UserUpdateDeviceOk.Str(user.Name(), device.Title()))
}

func eventDeviceDeleted(userID int64, id int64, uid string, title string) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventDeviceDeleted: ", uid, err)
		return
	}

	edge.Remove(strconv.FormatInt(id, 10))

	log.WithField("src", logStore.SystemLog).Warn(lang.UserDeleteDeviceOk.Str(user.Name(), title))
	user.Logger().Warn(lang.UserDeleteDeviceOk.Str(user.Name(), title))
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

	log.WithField("src", logStore.SystemLog).Info(lang.UserCreateEquipmentOk.Str(user.Name(), equipment.Title()))
	user.Logger().Info(lang.UserCreateEquipmentOk.Str(user.Name(), equipment.Title()))
	equipment.Logger().Info(lang.UserCreateEquipmentOk.Str(user.Name(), equipment.Title()))
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

	log.WithField("src", logStore.SystemLog).Info(lang.UserUpdateEquipmentOk.Str(user.Name(), equipment.Title()))
	user.Logger().Info(lang.UserUpdateEquipmentOk.Str(user.Name(), equipment.Title()))
	equipment.Logger().Info(lang.UserUpdateEquipmentOk.Str(user.Name(), equipment.Title()))
}

func eventEquipmentDeleted(userID int64, title string) {
	user, err := Store().GetUser(userID)
	if err != nil {
		log.Error("eventEquipmentDeleted: ", err)
		return
	}

	log.WithField("src", logStore.SystemLog).Warn(lang.UserDeleteEquipmentOk.Str(user.Name(), title))
	user.Logger().Warn(lang.UserDeleteEquipmentOk.Str(user.Name(), title))
}
