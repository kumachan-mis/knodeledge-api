package domain

type ChapterEntity struct {
	id        ChapterIdObject
	name      ChapterNameObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
	authorId  UserIdObject
}

func NewChapterEntity(
	id ChapterIdObject,
	name ChapterNameObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
	authorId UserIdObject,
) *ChapterEntity {
	return &ChapterEntity{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
		authorId:  authorId,
	}
}

func (e *ChapterEntity) Id() *ChapterIdObject {
	return &e.id
}

func (e *ChapterEntity) Name() *ChapterNameObject {
	return &e.name
}

func (e *ChapterEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *ChapterEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}

func (e *ChapterEntity) AuthoredBy(userId *UserIdObject) bool {
	return e.authorId.Equals(userId)
}
