package domain

type UserEntity struct {
	id UserIdObject
}

func NewUserEntity(id UserIdObject) (*UserEntity, error) {
	return &UserEntity{id: id}, nil
}

func (e *UserEntity) Id() *UserIdObject {
	return &e.id
}
