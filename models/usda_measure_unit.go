// Code generated by SQLBoiler 4.8.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// UsdaMeasureUnit is an object representing the database table.
type UsdaMeasureUnit struct {
	ID   int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`

	R *usdaMeasureUnitR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L usdaMeasureUnitL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UsdaMeasureUnitColumns = struct {
	ID   string
	Name string
}{
	ID:   "id",
	Name: "name",
}

var UsdaMeasureUnitTableColumns = struct {
	ID   string
	Name string
}{
	ID:   "usda_measure_unit.id",
	Name: "usda_measure_unit.name",
}

// Generated where

var UsdaMeasureUnitWhere = struct {
	ID   whereHelperint
	Name whereHelpernull_String
}{
	ID:   whereHelperint{field: "\"usda_measure_unit\".\"id\""},
	Name: whereHelpernull_String{field: "\"usda_measure_unit\".\"name\""},
}

// UsdaMeasureUnitRels is where relationship names are stored.
var UsdaMeasureUnitRels = struct {
	MeasureUnitUsdaFoodPortions string
}{
	MeasureUnitUsdaFoodPortions: "MeasureUnitUsdaFoodPortions",
}

// usdaMeasureUnitR is where relationships are stored.
type usdaMeasureUnitR struct {
	MeasureUnitUsdaFoodPortions UsdaFoodPortionSlice `boil:"MeasureUnitUsdaFoodPortions" json:"MeasureUnitUsdaFoodPortions" toml:"MeasureUnitUsdaFoodPortions" yaml:"MeasureUnitUsdaFoodPortions"`
}

// NewStruct creates a new relationship struct
func (*usdaMeasureUnitR) NewStruct() *usdaMeasureUnitR {
	return &usdaMeasureUnitR{}
}

// usdaMeasureUnitL is where Load methods for each relationship are stored.
type usdaMeasureUnitL struct{}

var (
	usdaMeasureUnitAllColumns            = []string{"id", "name"}
	usdaMeasureUnitColumnsWithoutDefault = []string{"id", "name"}
	usdaMeasureUnitColumnsWithDefault    = []string{}
	usdaMeasureUnitPrimaryKeyColumns     = []string{"id"}
)

type (
	// UsdaMeasureUnitSlice is an alias for a slice of pointers to UsdaMeasureUnit.
	// This should almost always be used instead of []UsdaMeasureUnit.
	UsdaMeasureUnitSlice []*UsdaMeasureUnit
	// UsdaMeasureUnitHook is the signature for custom UsdaMeasureUnit hook methods
	UsdaMeasureUnitHook func(context.Context, boil.ContextExecutor, *UsdaMeasureUnit) error

	usdaMeasureUnitQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	usdaMeasureUnitType                 = reflect.TypeOf(&UsdaMeasureUnit{})
	usdaMeasureUnitMapping              = queries.MakeStructMapping(usdaMeasureUnitType)
	usdaMeasureUnitPrimaryKeyMapping, _ = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, usdaMeasureUnitPrimaryKeyColumns)
	usdaMeasureUnitInsertCacheMut       sync.RWMutex
	usdaMeasureUnitInsertCache          = make(map[string]insertCache)
	usdaMeasureUnitUpdateCacheMut       sync.RWMutex
	usdaMeasureUnitUpdateCache          = make(map[string]updateCache)
	usdaMeasureUnitUpsertCacheMut       sync.RWMutex
	usdaMeasureUnitUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var usdaMeasureUnitBeforeInsertHooks []UsdaMeasureUnitHook
var usdaMeasureUnitBeforeUpdateHooks []UsdaMeasureUnitHook
var usdaMeasureUnitBeforeDeleteHooks []UsdaMeasureUnitHook
var usdaMeasureUnitBeforeUpsertHooks []UsdaMeasureUnitHook

var usdaMeasureUnitAfterInsertHooks []UsdaMeasureUnitHook
var usdaMeasureUnitAfterSelectHooks []UsdaMeasureUnitHook
var usdaMeasureUnitAfterUpdateHooks []UsdaMeasureUnitHook
var usdaMeasureUnitAfterDeleteHooks []UsdaMeasureUnitHook
var usdaMeasureUnitAfterUpsertHooks []UsdaMeasureUnitHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UsdaMeasureUnit) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UsdaMeasureUnit) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UsdaMeasureUnit) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UsdaMeasureUnit) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UsdaMeasureUnit) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UsdaMeasureUnit) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UsdaMeasureUnit) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UsdaMeasureUnit) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UsdaMeasureUnit) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaMeasureUnitAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUsdaMeasureUnitHook registers your hook function for all future operations.
func AddUsdaMeasureUnitHook(hookPoint boil.HookPoint, usdaMeasureUnitHook UsdaMeasureUnitHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		usdaMeasureUnitBeforeInsertHooks = append(usdaMeasureUnitBeforeInsertHooks, usdaMeasureUnitHook)
	case boil.BeforeUpdateHook:
		usdaMeasureUnitBeforeUpdateHooks = append(usdaMeasureUnitBeforeUpdateHooks, usdaMeasureUnitHook)
	case boil.BeforeDeleteHook:
		usdaMeasureUnitBeforeDeleteHooks = append(usdaMeasureUnitBeforeDeleteHooks, usdaMeasureUnitHook)
	case boil.BeforeUpsertHook:
		usdaMeasureUnitBeforeUpsertHooks = append(usdaMeasureUnitBeforeUpsertHooks, usdaMeasureUnitHook)
	case boil.AfterInsertHook:
		usdaMeasureUnitAfterInsertHooks = append(usdaMeasureUnitAfterInsertHooks, usdaMeasureUnitHook)
	case boil.AfterSelectHook:
		usdaMeasureUnitAfterSelectHooks = append(usdaMeasureUnitAfterSelectHooks, usdaMeasureUnitHook)
	case boil.AfterUpdateHook:
		usdaMeasureUnitAfterUpdateHooks = append(usdaMeasureUnitAfterUpdateHooks, usdaMeasureUnitHook)
	case boil.AfterDeleteHook:
		usdaMeasureUnitAfterDeleteHooks = append(usdaMeasureUnitAfterDeleteHooks, usdaMeasureUnitHook)
	case boil.AfterUpsertHook:
		usdaMeasureUnitAfterUpsertHooks = append(usdaMeasureUnitAfterUpsertHooks, usdaMeasureUnitHook)
	}
}

// One returns a single usdaMeasureUnit record from the query.
func (q usdaMeasureUnitQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UsdaMeasureUnit, error) {
	o := &UsdaMeasureUnit{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for usda_measure_unit")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UsdaMeasureUnit records from the query.
func (q usdaMeasureUnitQuery) All(ctx context.Context, exec boil.ContextExecutor) (UsdaMeasureUnitSlice, error) {
	var o []*UsdaMeasureUnit

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UsdaMeasureUnit slice")
	}

	if len(usdaMeasureUnitAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UsdaMeasureUnit records in the query.
func (q usdaMeasureUnitQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count usda_measure_unit rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q usdaMeasureUnitQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if usda_measure_unit exists")
	}

	return count > 0, nil
}

// MeasureUnitUsdaFoodPortions retrieves all the usda_food_portion's UsdaFoodPortions with an executor via measure_unit_id column.
func (o *UsdaMeasureUnit) MeasureUnitUsdaFoodPortions(mods ...qm.QueryMod) usdaFoodPortionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"usda_food_portion\".\"measure_unit_id\"=?", o.ID),
	)

	query := UsdaFoodPortions(queryMods...)
	queries.SetFrom(query.Query, "\"usda_food_portion\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"usda_food_portion\".*"})
	}

	return query
}

// LoadMeasureUnitUsdaFoodPortions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (usdaMeasureUnitL) LoadMeasureUnitUsdaFoodPortions(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUsdaMeasureUnit interface{}, mods queries.Applicator) error {
	var slice []*UsdaMeasureUnit
	var object *UsdaMeasureUnit

	if singular {
		object = maybeUsdaMeasureUnit.(*UsdaMeasureUnit)
	} else {
		slice = *maybeUsdaMeasureUnit.(*[]*UsdaMeasureUnit)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &usdaMeasureUnitR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &usdaMeasureUnitR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`usda_food_portion`),
		qm.WhereIn(`usda_food_portion.measure_unit_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load usda_food_portion")
	}

	var resultSlice []*UsdaFoodPortion
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice usda_food_portion")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on usda_food_portion")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for usda_food_portion")
	}

	if len(usdaFoodPortionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.MeasureUnitUsdaFoodPortions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &usdaFoodPortionR{}
			}
			foreign.R.MeasureUnit = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.MeasureUnitID) {
				local.R.MeasureUnitUsdaFoodPortions = append(local.R.MeasureUnitUsdaFoodPortions, foreign)
				if foreign.R == nil {
					foreign.R = &usdaFoodPortionR{}
				}
				foreign.R.MeasureUnit = local
				break
			}
		}
	}

	return nil
}

// AddMeasureUnitUsdaFoodPortions adds the given related objects to the existing relationships
// of the usda_measure_unit, optionally inserting them as new records.
// Appends related to o.R.MeasureUnitUsdaFoodPortions.
// Sets related.R.MeasureUnit appropriately.
func (o *UsdaMeasureUnit) AddMeasureUnitUsdaFoodPortions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UsdaFoodPortion) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.MeasureUnitID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"usda_food_portion\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"measure_unit_id"}),
				strmangle.WhereClause("\"", "\"", 2, usdaFoodPortionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.MeasureUnitID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &usdaMeasureUnitR{
			MeasureUnitUsdaFoodPortions: related,
		}
	} else {
		o.R.MeasureUnitUsdaFoodPortions = append(o.R.MeasureUnitUsdaFoodPortions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &usdaFoodPortionR{
				MeasureUnit: o,
			}
		} else {
			rel.R.MeasureUnit = o
		}
	}
	return nil
}

// SetMeasureUnitUsdaFoodPortions removes all previously related items of the
// usda_measure_unit replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.MeasureUnit's MeasureUnitUsdaFoodPortions accordingly.
// Replaces o.R.MeasureUnitUsdaFoodPortions with related.
// Sets related.R.MeasureUnit's MeasureUnitUsdaFoodPortions accordingly.
func (o *UsdaMeasureUnit) SetMeasureUnitUsdaFoodPortions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UsdaFoodPortion) error {
	query := "update \"usda_food_portion\" set \"measure_unit_id\" = null where \"measure_unit_id\" = $1"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.MeasureUnitUsdaFoodPortions {
			queries.SetScanner(&rel.MeasureUnitID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.MeasureUnit = nil
		}

		o.R.MeasureUnitUsdaFoodPortions = nil
	}
	return o.AddMeasureUnitUsdaFoodPortions(ctx, exec, insert, related...)
}

// RemoveMeasureUnitUsdaFoodPortions relationships from objects passed in.
// Removes related items from R.MeasureUnitUsdaFoodPortions (uses pointer comparison, removal does not keep order)
// Sets related.R.MeasureUnit.
func (o *UsdaMeasureUnit) RemoveMeasureUnitUsdaFoodPortions(ctx context.Context, exec boil.ContextExecutor, related ...*UsdaFoodPortion) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.MeasureUnitID, nil)
		if rel.R != nil {
			rel.R.MeasureUnit = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("measure_unit_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.MeasureUnitUsdaFoodPortions {
			if rel != ri {
				continue
			}

			ln := len(o.R.MeasureUnitUsdaFoodPortions)
			if ln > 1 && i < ln-1 {
				o.R.MeasureUnitUsdaFoodPortions[i] = o.R.MeasureUnitUsdaFoodPortions[ln-1]
			}
			o.R.MeasureUnitUsdaFoodPortions = o.R.MeasureUnitUsdaFoodPortions[:ln-1]
			break
		}
	}

	return nil
}

// UsdaMeasureUnits retrieves all the records using an executor.
func UsdaMeasureUnits(mods ...qm.QueryMod) usdaMeasureUnitQuery {
	mods = append(mods, qm.From("\"usda_measure_unit\""))
	return usdaMeasureUnitQuery{NewQuery(mods...)}
}

// FindUsdaMeasureUnit retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUsdaMeasureUnit(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*UsdaMeasureUnit, error) {
	usdaMeasureUnitObj := &UsdaMeasureUnit{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"usda_measure_unit\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, usdaMeasureUnitObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from usda_measure_unit")
	}

	if err = usdaMeasureUnitObj.doAfterSelectHooks(ctx, exec); err != nil {
		return usdaMeasureUnitObj, err
	}

	return usdaMeasureUnitObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UsdaMeasureUnit) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_measure_unit provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaMeasureUnitColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	usdaMeasureUnitInsertCacheMut.RLock()
	cache, cached := usdaMeasureUnitInsertCache[key]
	usdaMeasureUnitInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			usdaMeasureUnitAllColumns,
			usdaMeasureUnitColumnsWithDefault,
			usdaMeasureUnitColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"usda_measure_unit\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"usda_measure_unit\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into usda_measure_unit")
	}

	if !cached {
		usdaMeasureUnitInsertCacheMut.Lock()
		usdaMeasureUnitInsertCache[key] = cache
		usdaMeasureUnitInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UsdaMeasureUnit.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UsdaMeasureUnit) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	usdaMeasureUnitUpdateCacheMut.RLock()
	cache, cached := usdaMeasureUnitUpdateCache[key]
	usdaMeasureUnitUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			usdaMeasureUnitAllColumns,
			usdaMeasureUnitPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update usda_measure_unit, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"usda_measure_unit\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, usdaMeasureUnitPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, append(wl, usdaMeasureUnitPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update usda_measure_unit row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for usda_measure_unit")
	}

	if !cached {
		usdaMeasureUnitUpdateCacheMut.Lock()
		usdaMeasureUnitUpdateCache[key] = cache
		usdaMeasureUnitUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q usdaMeasureUnitQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for usda_measure_unit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for usda_measure_unit")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UsdaMeasureUnitSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaMeasureUnitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"usda_measure_unit\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, usdaMeasureUnitPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in usdaMeasureUnit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all usdaMeasureUnit")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UsdaMeasureUnit) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_measure_unit provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaMeasureUnitColumnsWithDefault, o)

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

	usdaMeasureUnitUpsertCacheMut.RLock()
	cache, cached := usdaMeasureUnitUpsertCache[key]
	usdaMeasureUnitUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			usdaMeasureUnitAllColumns,
			usdaMeasureUnitColumnsWithDefault,
			usdaMeasureUnitColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			usdaMeasureUnitAllColumns,
			usdaMeasureUnitPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert usda_measure_unit, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(usdaMeasureUnitPrimaryKeyColumns))
			copy(conflict, usdaMeasureUnitPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"usda_measure_unit\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(usdaMeasureUnitType, usdaMeasureUnitMapping, ret)
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
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert usda_measure_unit")
	}

	if !cached {
		usdaMeasureUnitUpsertCacheMut.Lock()
		usdaMeasureUnitUpsertCache[key] = cache
		usdaMeasureUnitUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UsdaMeasureUnit record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UsdaMeasureUnit) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UsdaMeasureUnit provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), usdaMeasureUnitPrimaryKeyMapping)
	sql := "DELETE FROM \"usda_measure_unit\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from usda_measure_unit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for usda_measure_unit")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q usdaMeasureUnitQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no usdaMeasureUnitQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usda_measure_unit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_measure_unit")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UsdaMeasureUnitSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(usdaMeasureUnitBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaMeasureUnitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"usda_measure_unit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaMeasureUnitPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usdaMeasureUnit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_measure_unit")
	}

	if len(usdaMeasureUnitAfterDeleteHooks) != 0 {
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
func (o *UsdaMeasureUnit) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUsdaMeasureUnit(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsdaMeasureUnitSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UsdaMeasureUnitSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaMeasureUnitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"usda_measure_unit\".* FROM \"usda_measure_unit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaMeasureUnitPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UsdaMeasureUnitSlice")
	}

	*o = slice

	return nil
}

// UsdaMeasureUnitExists checks if the UsdaMeasureUnit row exists.
func UsdaMeasureUnitExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"usda_measure_unit\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if usda_measure_unit exists")
	}

	return exists, nil
}