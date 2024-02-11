package domain

type ProjectEntity struct {
	name        ProjectNameObject
	description ProjectDescriptionObject
}

func NewProjectEntity(name ProjectNameObject, description ProjectDescriptionObject) (*ProjectEntity, error) {
	return &ProjectEntity{name: name, description: description}, nil
}

func (e *ProjectEntity) Name() ProjectNameObject {
	return e.name
}

func (e *ProjectEntity) Description() ProjectDescriptionObject {
	return e.description
}
