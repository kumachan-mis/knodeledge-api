package record

type HelloWorldEntry struct {
	Name    string `firestore:"name"`
	Message string `firestore:"message"`
}
