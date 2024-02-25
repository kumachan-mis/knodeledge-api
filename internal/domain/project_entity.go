package domain

type ProjectEntity struct {
	id          ProjectIdObject
	name        ProjectNameObject
	description ProjectDescriptionObject
	createdAt   CreatedAtObject
	updatedAt   UpdatedAtObject
}

func NewProjectEntity(
	id ProjectIdObject,
	name ProjectNameObject,
	description ProjectDescriptionObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *ProjectEntity {
	return &ProjectEntity{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
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
