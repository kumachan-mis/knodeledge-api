package domain

import "time"

type CreatedAtObject struct {
	value time.Time
}

func NewCreatedAtObject(createdAt time.Time) (*CreatedAtObject, error) {
	return &CreatedAtObject{value: createdAt}, nil
}

func (o *CreatedAtObject) Value() time.Time {
	return o.value
}
