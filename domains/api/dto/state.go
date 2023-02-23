package dto

import (
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
)

type State interface {
	usersDTO.WithUsers
	authDTO.WithAuth
	configDTO.WithConfig
	storageDTO.WithStorage
}
