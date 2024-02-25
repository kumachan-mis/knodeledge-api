package domain

type ProjectWithoutAutofieldEntity struct {
	name        ProjectNameObject
	description ProjectDescriptionObject
}

func NewProjectWithoutAutofieldEntity(
	name ProjectNameObject,
	description ProjectDescriptionObject,
) *ProjectWithoutAutofieldEntity {
	return &ProjectWithoutAutofieldEntity{name: name, description: description}
}

func (e *ProjectWithoutAutofieldEntity) Name() *ProjectNameObject {
	return &e.name
}

func (e *ProjectWithoutAutofieldEntity) Description() *ProjectDescriptionObject {
	return &e.description
}
