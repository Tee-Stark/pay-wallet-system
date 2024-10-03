package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateRef() string {
	// generate unique payment reference
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // handle error appropriately in production code
	}
	randomString := hex.EncodeToString(randomBytes)

	// Combine the timestamp and the random string to form the reference
	ref := fmt.Sprintf("pref_%d_%s", timestamp, randomString)
	return ref
}

func GenerateUUID() string {
	return uuid.New().String()
}
