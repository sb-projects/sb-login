package models

import "time"

type (
	deviceStatus int
	Device       struct {
		ID               string
		RegistrationDate time.Time
		Status           deviceStatus
	}
)

const (
	Device_Active deviceStatus = iota + 1
	Device_Inactive
)
