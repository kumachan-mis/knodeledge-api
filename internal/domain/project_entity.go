package domain

type ProjectEntity struct {
	id          ProjectIdObject
	name        ProjectNameObject
	description ProjectDescriptionObject
}

func NewProjectEntity(
	id ProjectIdObject,
	name ProjectNameObject,
	description ProjectDescriptionObject,
) (*ProjectEntity, error) {
	return &ProjectEntity{id: id, name: name, description: description}, nil
}

func (e *ProjectEntity) Id() ProjectIdObject {
	return e.id
}

func (e *ProjectEntity) Name() ProjectNameObject {
	return e.name
}

func (e *ProjectEntity) Description() ProjectDescriptionObject {
	return e.description
}
