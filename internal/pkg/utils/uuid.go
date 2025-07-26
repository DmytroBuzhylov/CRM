package utils

import "github.com/google/uuid"

func GenerateInvitationCode() string {
	return uuid.New().String()
}
