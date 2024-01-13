package db

const COLLECTION = "hello_world"

func GetHelloWorld(name string) (string, error) {
	iter := FirestoreClient.Collection(COLLECTION).
		Where("name", "==", name).
		Limit(1).
		Documents(FirestoreContext)

	snapshot, err := iter.Next()
	if err != nil {
		return "", err
	}

	return snapshot.Data()["message"].(string), nil
}

func SetHelloWorld(name string, message string) (string, error) {
	reference, _, err := FirestoreClient.
		Collection(COLLECTION).
		Add(FirestoreContext, map[string]interface{}{
			"name":    name,
			"message": message,
		})
	if err != nil {
		return "", err
	}

	snapshot, err := reference.Get(FirestoreContext)
	if err != nil {
		return "", err
	}

	return snapshot.Data()["message"].(string), nil
}
