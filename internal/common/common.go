package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/types"
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

func DecimalFromFloat(f float64) types.Decimal {
	d := decimal.WithContext(types.DecimalContext)
	d.SetFloat64(f)
	return types.NewDecimal(d)
}
func NullDecimalFromFloat(f *float64) types.NullDecimal {
	if f == nil {
		return types.NewNullDecimal(nil)
	}
	d := decimal.WithContext(types.DecimalContext)
	d.SetFloat64(*f)
	return types.NewNullDecimal(d)
}
