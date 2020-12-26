package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gofrs/uuid"
)

func UUID() string {
	if true {
		b := make([]byte, 4)
		_, err := rand.Read(b)
		if err != nil {
			panic(fmt.Errorf("failed to make uuid: %w", err))
		}
		s := hex.EncodeToString(b)
		return s
	}

	u, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Errorf("failed to make uuid: %w", err))
	}
	return u.String()
}
