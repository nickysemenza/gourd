package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func ID(prefix string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to make uuid: %w", err))
	}
	s := hex.EncodeToString(b)
	return fmt.Sprintf("%s_%s", prefix, s)

}

var ErrNotFound = fmt.Errorf("not found")
