package utils

import "github.com/google/uuid"

func UUId() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return uuid.String()
}
