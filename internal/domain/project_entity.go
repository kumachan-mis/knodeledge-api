package domain

type ProjectEntity struct {
	id          ProjectIdObject
	name        ProjectNameObject
	description ProjectDescriptionObject
	createdAt   CreatedAtObject
	updatedAt   UpdatedAtObject
	authorId    UserIdObject
}

func NewProjectEntity(
	id ProjectIdObject,
	name ProjectNameObject,
	description ProjectDescriptionObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
	authorId UserIdObject,
) *ProjectEntity {
	return &ProjectEntity{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		authorId:    authorId,
	}
}

func (e *ProjectEntity) Id() *ProjectIdObject {
	return &e.id
}

func (e *ProjectEntity) Name() *ProjectNameObject {
	return &e.name
}

func (e *ProjectEntity) Description() *ProjectDescriptionObject {
	return &e.description
}

func (e *ProjectEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *ProjectEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}

func (e *ProjectEntity) AuthoredBy(userId *UserIdObject) bool {
	return e.authorId.Equals(userId)
}
