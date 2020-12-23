package common

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func UUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Errorf("failed to make uuid: %w", err))
	}
	return u.String()
}
