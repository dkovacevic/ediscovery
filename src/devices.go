package main

import (
	"encoding/json"
	"net/http"
)

type Device struct {
	JID    string `json:"jid"`
	Name   string `json:"name"`
	Device string `json:"device"`
}

func getDevices(writer http.ResponseWriter, _ *http.Request) {
	// Get all devices from the container
	devices, err := container.GetAllDevices()
	if err != nil {
		http.Error(writer, "Unable to fetch devices", http.StatusInternalServerError)
		return
	}

	// Prepare a slice to hold the device details
	var deviceDetails []Device

	// Iterate over the devices and populate the device details slice
	for _, deviceStore := range devices {
		deviceDetails = append(deviceDetails, Device{
			JID:    deviceStore.ID.User,
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
