package api

import (
	"github.com/outcatcher/anwil/domains/auth"
	services "github.com/outcatcher/anwil/domains/services/dto"
	"github.com/outcatcher/anwil/domains/users"
)

// serviceID - ID of the service.
type serviceID string

// ServiceDefinition IDs.
const (
	serviceAuth  serviceID = "authentication"
	serviceUsers serviceID = "users"
)

var preparedServices = map[serviceID]services.Service{
	serviceAuth:  auth.New(),
	serviceUsers: users.New(),
}
