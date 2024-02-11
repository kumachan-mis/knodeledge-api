package record

type ProjectEntry struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	UserId      string `firestore:"userId"`
}
