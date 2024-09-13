package handlers

import (
	"encoding/json"
	"lh-whatsapp/src/meow"
	"lh-whatsapp/src/models"
	"net/http"
)

func GetUsers(writer http.ResponseWriter, _ *http.Request) {
	// Get all devices from the container
	devices, err := meow.GetAllDevices()
	if err != nil {
		http.Error(writer, "Unable to fetch devices", http.StatusInternalServerError)
		return
	}

	// Prepare a slice to hold the device details
	var deviceDetails []models.User

	// Iterate over the devices and populate the device details slice
	for _, deviceStore := range devices {
		deviceDetails = append(deviceDetails, models.User{
			JID:    deviceStore.ID.String(),
			User:   deviceStore.ID.User,
			Name:   deviceStore.PushName,
			Device: deviceStore.Platform,
		})
	}

	// Convert the device details to JSON
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(deviceDetails)
	if err != nil {
		http.Error(writer, "Unable to encode devices to JSON", http.StatusInternalServerError)
		return
	}
}
