package domain

import "time"

type UpdatedAtObject struct {
	value time.Time
}

func NewUpdatedAtObject(updatedAt time.Time) (*UpdatedAtObject, error) {
	return &UpdatedAtObject{value: updatedAt}, nil
}

func (o *UpdatedAtObject) Value() time.Time {
	return o.value
}
