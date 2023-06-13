package utils

import (
	"crypto/sha256"
	"encoding/json"
)

// HashCompute computes the SHA256 hash of the specified object.
func HashCompute(obj interface{}) (string, error) {
	// Create a new SHA256 hash instance.
	hashIns := sha256.New()
	// Encode the object as JSON.
	encodeJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	// Write the JSON-encoded object to the hash instance.
	_, err = hashIns.Write(encodeJson)
	if err != nil {
		return "", err
	}
	// Compute the final hash and return it as a string.
	return string(hashIns.Sum(nil)), nil
}
