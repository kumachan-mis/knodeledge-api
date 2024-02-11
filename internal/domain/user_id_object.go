package domain

import "fmt"

type UserIdObject struct {
	value string
}

func NewUserIdObject(userId string) (*UserIdObject, error) {
	if userId == "" {
		return nil, fmt.Errorf("user id is required")
	}
	return &UserIdObject{value: userId}, nil
}

func (o *UserIdObject) Value() string {
	return o.value
}
