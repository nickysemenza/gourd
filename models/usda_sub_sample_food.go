// Code generated by SQLBoiler 4.10.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UsdaSubSampleFood is an object representing the database table.
type UsdaSubSampleFood struct {
	FDCID             int      `boil:"fdc_id" json:"fdc_id" toml:"fdc_id" yaml:"fdc_id"`
	FDCIDOfSampleFood null.Int `boil:"fdc_id_of_sample_food" json:"fdc_id_of_sample_food,omitempty" toml:"fdc_id_of_sample_food" yaml:"fdc_id_of_sample_food,omitempty"`

	R *usdaSubSampleFoodR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L usdaSubSampleFoodL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UsdaSubSampleFoodColumns = struct {
	FDCID             string
	FDCIDOfSampleFood string
}{
	FDCID:             "fdc_id",
	FDCIDOfSampleFood: "fdc_id_of_sample_food",
}

var UsdaSubSampleFoodTableColumns = struct {
	FDCID             string
	FDCIDOfSampleFood string
}{
	FDCID:             "usda_sub_sample_food.fdc_id",
	FDCIDOfSampleFood: "usda_sub_sample_food.fdc_id_of_sample_food",
}

// Generated where

var UsdaSubSampleFoodWhere = struct {
	FDCID             whereHelperint
	FDCIDOfSampleFood whereHelpernull_Int
}{
	FDCID:             whereHelperint{field: "\"usda_sub_sample_food\".\"fdc_id\""},
	FDCIDOfSampleFood: whereHelpernull_Int{field: "\"usda_sub_sample_food\".\"fdc_id_of_sample_food\""},
}

// UsdaSubSampleFoodRels is where relationship names are stored.
var UsdaSubSampleFoodRels = struct {
	FDC                       string
	FDCIDOfSampleFoodUsdaFood string
}{
	FDC:                       "FDC",
	FDCIDOfSampleFoodUsdaFood: "FDCIDOfSampleFoodUsdaFood",
}

// usdaSubSampleFoodR is where relationships are stored.
type usdaSubSampleFoodR struct {
	FDC                       *UsdaFood `boil:"FDC" json:"FDC" toml:"FDC" yaml:"FDC"`
	FDCIDOfSampleFoodUsdaFood *UsdaFood `boil:"FDCIDOfSampleFoodUsdaFood" json:"FDCIDOfSampleFoodUsdaFood" toml:"FDCIDOfSampleFoodUsdaFood" yaml:"FDCIDOfSampleFoodUsdaFood"`
}

// NewStruct creates a new relationship struct
func (*usdaSubSampleFoodR) NewStruct() *usdaSubSampleFoodR {
	return &usdaSubSampleFoodR{}
}

// usdaSubSampleFoodL is where Load methods for each relationship are stored.
type usdaSubSampleFoodL struct{}

var (
	usdaSubSampleFoodAllColumns            = []string{"fdc_id", "fdc_id_of_sample_food"}
	usdaSubSampleFoodColumnsWithoutDefault = []string{"fdc_id"}
	usdaSubSampleFoodColumnsWithDefault    = []string{"fdc_id_of_sample_food"}
	usdaSubSampleFoodPrimaryKeyColumns     = []string{"fdc_id"}
	usdaSubSampleFoodGeneratedColumns      = []string{}
)

type (
	// UsdaSubSampleFoodSlice is an alias for a slice of pointers to UsdaSubSampleFood.
	// This should almost always be used instead of []UsdaSubSampleFood.
	UsdaSubSampleFoodSlice []*UsdaSubSampleFood
	// UsdaSubSampleFoodHook is the signature for custom UsdaSubSampleFood hook methods
	UsdaSubSampleFoodHook func(context.Context, boil.ContextExecutor, *UsdaSubSampleFood) error

	usdaSubSampleFoodQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	usdaSubSampleFoodType                 = reflect.TypeOf(&UsdaSubSampleFood{})
	usdaSubSampleFoodMapping              = queries.MakeStructMapping(usdaSubSampleFoodType)
	usdaSubSampleFoodPrimaryKeyMapping, _ = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, usdaSubSampleFoodPrimaryKeyColumns)
	usdaSubSampleFoodInsertCacheMut       sync.RWMutex
	usdaSubSampleFoodInsertCache          = make(map[string]insertCache)
	usdaSubSampleFoodUpdateCacheMut       sync.RWMutex
	usdaSubSampleFoodUpdateCache          = make(map[string]updateCache)
	usdaSubSampleFoodUpsertCacheMut       sync.RWMutex
	usdaSubSampleFoodUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var usdaSubSampleFoodAfterSelectHooks []UsdaSubSampleFoodHook

var usdaSubSampleFoodBeforeInsertHooks []UsdaSubSampleFoodHook
var usdaSubSampleFoodAfterInsertHooks []UsdaSubSampleFoodHook

var usdaSubSampleFoodBeforeUpdateHooks []UsdaSubSampleFoodHook
var usdaSubSampleFoodAfterUpdateHooks []UsdaSubSampleFoodHook

var usdaSubSampleFoodBeforeDeleteHooks []UsdaSubSampleFoodHook
var usdaSubSampleFoodAfterDeleteHooks []UsdaSubSampleFoodHook

var usdaSubSampleFoodBeforeUpsertHooks []UsdaSubSampleFoodHook
var usdaSubSampleFoodAfterUpsertHooks []UsdaSubSampleFoodHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UsdaSubSampleFood) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UsdaSubSampleFood) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UsdaSubSampleFood) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UsdaSubSampleFood) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UsdaSubSampleFood) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UsdaSubSampleFood) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UsdaSubSampleFood) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UsdaSubSampleFood) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UsdaSubSampleFood) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaSubSampleFoodAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUsdaSubSampleFoodHook registers your hook function for all future operations.
func AddUsdaSubSampleFoodHook(hookPoint boil.HookPoint, usdaSubSampleFoodHook UsdaSubSampleFoodHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		usdaSubSampleFoodAfterSelectHooks = append(usdaSubSampleFoodAfterSelectHooks, usdaSubSampleFoodHook)
	case boil.BeforeInsertHook:
		usdaSubSampleFoodBeforeInsertHooks = append(usdaSubSampleFoodBeforeInsertHooks, usdaSubSampleFoodHook)
	case boil.AfterInsertHook:
		usdaSubSampleFoodAfterInsertHooks = append(usdaSubSampleFoodAfterInsertHooks, usdaSubSampleFoodHook)
	case boil.BeforeUpdateHook:
		usdaSubSampleFoodBeforeUpdateHooks = append(usdaSubSampleFoodBeforeUpdateHooks, usdaSubSampleFoodHook)
	case boil.AfterUpdateHook:
		usdaSubSampleFoodAfterUpdateHooks = append(usdaSubSampleFoodAfterUpdateHooks, usdaSubSampleFoodHook)
	case boil.BeforeDeleteHook:
		usdaSubSampleFoodBeforeDeleteHooks = append(usdaSubSampleFoodBeforeDeleteHooks, usdaSubSampleFoodHook)
	case boil.AfterDeleteHook:
		usdaSubSampleFoodAfterDeleteHooks = append(usdaSubSampleFoodAfterDeleteHooks, usdaSubSampleFoodHook)
	case boil.BeforeUpsertHook:
		usdaSubSampleFoodBeforeUpsertHooks = append(usdaSubSampleFoodBeforeUpsertHooks, usdaSubSampleFoodHook)
	case boil.AfterUpsertHook:
		usdaSubSampleFoodAfterUpsertHooks = append(usdaSubSampleFoodAfterUpsertHooks, usdaSubSampleFoodHook)
	}
}

// One returns a single usdaSubSampleFood record from the query.
func (q usdaSubSampleFoodQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UsdaSubSampleFood, error) {
	o := &UsdaSubSampleFood{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for usda_sub_sample_food")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UsdaSubSampleFood records from the query.
func (q usdaSubSampleFoodQuery) All(ctx context.Context, exec boil.ContextExecutor) (UsdaSubSampleFoodSlice, error) {
	var o []*UsdaSubSampleFood

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UsdaSubSampleFood slice")
	}

	if len(usdaSubSampleFoodAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UsdaSubSampleFood records in the query.
func (q usdaSubSampleFoodQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count usda_sub_sample_food rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q usdaSubSampleFoodQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if usda_sub_sample_food exists")
	}

	return count > 0, nil
}

// FDC pointed to by the foreign key.
func (o *UsdaSubSampleFood) FDC(mods ...qm.QueryMod) usdaFoodQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"fdc_id\" = ?", o.FDCID),
	}

	queryMods = append(queryMods, mods...)

	return UsdaFoods(queryMods...)
}

// FDCIDOfSampleFoodUsdaFood pointed to by the foreign key.
func (o *UsdaSubSampleFood) FDCIDOfSampleFoodUsdaFood(mods ...qm.QueryMod) usdaFoodQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"fdc_id\" = ?", o.FDCIDOfSampleFood),
	}

	queryMods = append(queryMods, mods...)

	return UsdaFoods(queryMods...)
}

// LoadFDC allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (usdaSubSampleFoodL) LoadFDC(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUsdaSubSampleFood interface{}, mods queries.Applicator) error {
	var slice []*UsdaSubSampleFood
	var object *UsdaSubSampleFood

	if singular {
		object = maybeUsdaSubSampleFood.(*UsdaSubSampleFood)
	} else {
		slice = *maybeUsdaSubSampleFood.(*[]*UsdaSubSampleFood)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &usdaSubSampleFoodR{}
		}
		args = append(args, object.FDCID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &usdaSubSampleFoodR{}
			}

			for _, a := range args {
				if a == obj.FDCID {
					continue Outer
				}
			}

			args = append(args, obj.FDCID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`usda_food`),
		qm.WhereIn(`usda_food.fdc_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UsdaFood")
	}

	var resultSlice []*UsdaFood
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UsdaFood")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for usda_food")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for usda_food")
	}

	if len(usdaSubSampleFoodAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.FDC = foreign
		if foreign.R == nil {
			foreign.R = &usdaFoodR{}
		}
		foreign.R.FDCUsdaSubSampleFood = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.FDCID == foreign.FDCID {
				local.R.FDC = foreign
				if foreign.R == nil {
					foreign.R = &usdaFoodR{}
				}
				foreign.R.FDCUsdaSubSampleFood = local
				break
			}
		}
	}

	return nil
}

// LoadFDCIDOfSampleFoodUsdaFood allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (usdaSubSampleFoodL) LoadFDCIDOfSampleFoodUsdaFood(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUsdaSubSampleFood interface{}, mods queries.Applicator) error {
	var slice []*UsdaSubSampleFood
	var object *UsdaSubSampleFood

	if singular {
		object = maybeUsdaSubSampleFood.(*UsdaSubSampleFood)
	} else {
		slice = *maybeUsdaSubSampleFood.(*[]*UsdaSubSampleFood)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &usdaSubSampleFoodR{}
		}
		if !queries.IsNil(object.FDCIDOfSampleFood) {
			args = append(args, object.FDCIDOfSampleFood)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &usdaSubSampleFoodR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.FDCIDOfSampleFood) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.FDCIDOfSampleFood) {
				args = append(args, obj.FDCIDOfSampleFood)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`usda_food`),
		qm.WhereIn(`usda_food.fdc_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UsdaFood")
	}

	var resultSlice []*UsdaFood
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UsdaFood")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for usda_food")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for usda_food")
	}

	if len(usdaSubSampleFoodAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.FDCIDOfSampleFoodUsdaFood = foreign
		if foreign.R == nil {
			foreign.R = &usdaFoodR{}
		}
		foreign.R.FDCIDOfSampleFoodUsdaSubSampleFoods = append(foreign.R.FDCIDOfSampleFoodUsdaSubSampleFoods, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.FDCIDOfSampleFood, foreign.FDCID) {
				local.R.FDCIDOfSampleFoodUsdaFood = foreign
				if foreign.R == nil {
					foreign.R = &usdaFoodR{}
				}
				foreign.R.FDCIDOfSampleFoodUsdaSubSampleFoods = append(foreign.R.FDCIDOfSampleFoodUsdaSubSampleFoods, local)
				break
			}
		}
	}

	return nil
}

// SetFDC of the usdaSubSampleFood to the related item.
// Sets o.R.FDC to related.
// Adds o to related.R.FDCUsdaSubSampleFood.
func (o *UsdaSubSampleFood) SetFDC(ctx context.Context, exec boil.ContextExecutor, insert bool, related *UsdaFood) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"usda_sub_sample_food\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"fdc_id"}),
		strmangle.WhereClause("\"", "\"", 2, usdaSubSampleFoodPrimaryKeyColumns),
	)
	values := []interface{}{related.FDCID, o.FDCID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.FDCID = related.FDCID
	if o.R == nil {
		o.R = &usdaSubSampleFoodR{
			FDC: related,
		}
	} else {
		o.R.FDC = related
	}

	if related.R == nil {
		related.R = &usdaFoodR{
			FDCUsdaSubSampleFood: o,
		}
	} else {
		related.R.FDCUsdaSubSampleFood = o
	}

	return nil
}

// SetFDCIDOfSampleFoodUsdaFood of the usdaSubSampleFood to the related item.
// Sets o.R.FDCIDOfSampleFoodUsdaFood to related.
// Adds o to related.R.FDCIDOfSampleFoodUsdaSubSampleFoods.
func (o *UsdaSubSampleFood) SetFDCIDOfSampleFoodUsdaFood(ctx context.Context, exec boil.ContextExecutor, insert bool, related *UsdaFood) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"usda_sub_sample_food\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"fdc_id_of_sample_food"}),
		strmangle.WhereClause("\"", "\"", 2, usdaSubSampleFoodPrimaryKeyColumns),
	)
	values := []interface{}{related.FDCID, o.FDCID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.FDCIDOfSampleFood, related.FDCID)
	if o.R == nil {
		o.R = &usdaSubSampleFoodR{
			FDCIDOfSampleFoodUsdaFood: related,
		}
	} else {
		o.R.FDCIDOfSampleFoodUsdaFood = related
	}

	if related.R == nil {
		related.R = &usdaFoodR{
			FDCIDOfSampleFoodUsdaSubSampleFoods: UsdaSubSampleFoodSlice{o},
		}
	} else {
		related.R.FDCIDOfSampleFoodUsdaSubSampleFoods = append(related.R.FDCIDOfSampleFoodUsdaSubSampleFoods, o)
	}

	return nil
}

// RemoveFDCIDOfSampleFoodUsdaFood relationship.
// Sets o.R.FDCIDOfSampleFoodUsdaFood to nil.
// Removes o from all passed in related items' relationships struct.
func (o *UsdaSubSampleFood) RemoveFDCIDOfSampleFoodUsdaFood(ctx context.Context, exec boil.ContextExecutor, related *UsdaFood) error {
	var err error

	queries.SetScanner(&o.FDCIDOfSampleFood, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("fdc_id_of_sample_food")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.FDCIDOfSampleFoodUsdaFood = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.FDCIDOfSampleFoodUsdaSubSampleFoods {
		if queries.Equal(o.FDCIDOfSampleFood, ri.FDCIDOfSampleFood) {
			continue
		}

		ln := len(related.R.FDCIDOfSampleFoodUsdaSubSampleFoods)
		if ln > 1 && i < ln-1 {
			related.R.FDCIDOfSampleFoodUsdaSubSampleFoods[i] = related.R.FDCIDOfSampleFoodUsdaSubSampleFoods[ln-1]
		}
		related.R.FDCIDOfSampleFoodUsdaSubSampleFoods = related.R.FDCIDOfSampleFoodUsdaSubSampleFoods[:ln-1]
		break
	}
	return nil
}

// UsdaSubSampleFoods retrieves all the records using an executor.
func UsdaSubSampleFoods(mods ...qm.QueryMod) usdaSubSampleFoodQuery {
	mods = append(mods, qm.From("\"usda_sub_sample_food\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"usda_sub_sample_food\".*"})
	}

	return usdaSubSampleFoodQuery{q}
}

// FindUsdaSubSampleFood retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUsdaSubSampleFood(ctx context.Context, exec boil.ContextExecutor, fDCID int, selectCols ...string) (*UsdaSubSampleFood, error) {
	usdaSubSampleFoodObj := &UsdaSubSampleFood{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"usda_sub_sample_food\" where \"fdc_id\"=$1", sel,
	)

	q := queries.Raw(query, fDCID)

	err := q.Bind(ctx, exec, usdaSubSampleFoodObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from usda_sub_sample_food")
	}

	if err = usdaSubSampleFoodObj.doAfterSelectHooks(ctx, exec); err != nil {
		return usdaSubSampleFoodObj, err
	}

	return usdaSubSampleFoodObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UsdaSubSampleFood) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_sub_sample_food provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaSubSampleFoodColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	usdaSubSampleFoodInsertCacheMut.RLock()
	cache, cached := usdaSubSampleFoodInsertCache[key]
	usdaSubSampleFoodInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			usdaSubSampleFoodAllColumns,
			usdaSubSampleFoodColumnsWithDefault,
			usdaSubSampleFoodColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"usda_sub_sample_food\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"usda_sub_sample_food\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into usda_sub_sample_food")
	}

	if !cached {
		usdaSubSampleFoodInsertCacheMut.Lock()
		usdaSubSampleFoodInsertCache[key] = cache
		usdaSubSampleFoodInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UsdaSubSampleFood.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UsdaSubSampleFood) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	usdaSubSampleFoodUpdateCacheMut.RLock()
	cache, cached := usdaSubSampleFoodUpdateCache[key]
	usdaSubSampleFoodUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			usdaSubSampleFoodAllColumns,
			usdaSubSampleFoodPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update usda_sub_sample_food, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"usda_sub_sample_food\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, usdaSubSampleFoodPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, append(wl, usdaSubSampleFoodPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update usda_sub_sample_food row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for usda_sub_sample_food")
	}

	if !cached {
		usdaSubSampleFoodUpdateCacheMut.Lock()
		usdaSubSampleFoodUpdateCache[key] = cache
		usdaSubSampleFoodUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q usdaSubSampleFoodQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for usda_sub_sample_food")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for usda_sub_sample_food")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UsdaSubSampleFoodSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaSubSampleFoodPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"usda_sub_sample_food\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, usdaSubSampleFoodPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in usdaSubSampleFood slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all usdaSubSampleFood")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UsdaSubSampleFood) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_sub_sample_food provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaSubSampleFoodColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	usdaSubSampleFoodUpsertCacheMut.RLock()
	cache, cached := usdaSubSampleFoodUpsertCache[key]
	usdaSubSampleFoodUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			usdaSubSampleFoodAllColumns,
			usdaSubSampleFoodColumnsWithDefault,
			usdaSubSampleFoodColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			usdaSubSampleFoodAllColumns,
			usdaSubSampleFoodPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert usda_sub_sample_food, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(usdaSubSampleFoodPrimaryKeyColumns))
			copy(conflict, usdaSubSampleFoodPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"usda_sub_sample_food\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(usdaSubSampleFoodType, usdaSubSampleFoodMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert usda_sub_sample_food")
	}

	if !cached {
		usdaSubSampleFoodUpsertCacheMut.Lock()
		usdaSubSampleFoodUpsertCache[key] = cache
		usdaSubSampleFoodUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UsdaSubSampleFood record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UsdaSubSampleFood) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UsdaSubSampleFood provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), usdaSubSampleFoodPrimaryKeyMapping)
	sql := "DELETE FROM \"usda_sub_sample_food\" WHERE \"fdc_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from usda_sub_sample_food")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for usda_sub_sample_food")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q usdaSubSampleFoodQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no usdaSubSampleFoodQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usda_sub_sample_food")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_sub_sample_food")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UsdaSubSampleFoodSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(usdaSubSampleFoodBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaSubSampleFoodPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"usda_sub_sample_food\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaSubSampleFoodPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usdaSubSampleFood slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_sub_sample_food")
	}

	if len(usdaSubSampleFoodAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UsdaSubSampleFood) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUsdaSubSampleFood(ctx, exec, o.FDCID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsdaSubSampleFoodSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UsdaSubSampleFoodSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaSubSampleFoodPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"usda_sub_sample_food\".* FROM \"usda_sub_sample_food\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaSubSampleFoodPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UsdaSubSampleFoodSlice")
	}

	*o = slice

	return nil
}

// UsdaSubSampleFoodExists checks if the UsdaSubSampleFood row exists.
func UsdaSubSampleFoodExists(ctx context.Context, exec boil.ContextExecutor, fDCID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"usda_sub_sample_food\" where \"fdc_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, fDCID)
	}
	row := exec.QueryRowContext(ctx, sql, fDCID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if usda_sub_sample_food exists")
	}

	return exists, nil
}
