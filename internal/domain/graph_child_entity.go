package domain

type GraphChildEntity struct {
	name        GraphNameObject
	relation    GraphRelationObject
	description GraphDescriptionObject
	children    GraphChildrenEntity
}

func NewGraphChildEntity(
	name GraphNameObject,
	relation GraphRelationObject,
	description GraphDescriptionObject,
	children GraphChildrenEntity,
) *GraphChildEntity {
	return &GraphChildEntity{
		name:        name,
		relation:    relation,
		description: description,
		children:    children,
	}
}

func (e *GraphChildEntity) Name() *GraphNameObject {
	return &e.name
}

func (e *GraphChildEntity) Relation() *GraphRelationObject {
	return &e.relation
}

func (e *GraphChildEntity) Description() *GraphDescriptionObject {
	return &e.description
}

func (e *GraphChildEntity) Children() *GraphChildrenEntity {
	return &e.children
}
