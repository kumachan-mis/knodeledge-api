package domain

type UserEntity struct {
	id UserIdObject
}

func NewUserEntity(id UserIdObject) *UserEntity {
	return &UserEntity{id: id}
}

func (e *UserEntity) Id() *UserIdObject {
	return &e.id
}
