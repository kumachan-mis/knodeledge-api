package record

type ProjectEntry struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description,omitempty"`
	UserId      string `firestore:"userId"`
}
