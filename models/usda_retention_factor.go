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

// UsdaRetentionFactor is an object representing the database table.
type UsdaRetentionFactor struct {
	ID          int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Code        null.String `boil:"code" json:"code,omitempty" toml:"code" yaml:"code,omitempty"`
	FoodGroupID null.Int    `boil:"food_group_id" json:"food_group_id,omitempty" toml:"food_group_id" yaml:"food_group_id,omitempty"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`

	R *usdaRetentionFactorR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L usdaRetentionFactorL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UsdaRetentionFactorColumns = struct {
	ID          string
	Code        string
	FoodGroupID string
	Description string
}{
	ID:          "id",
	Code:        "code",
	FoodGroupID: "food_group_id",
	Description: "description",
}

var UsdaRetentionFactorTableColumns = struct {
	ID          string
	Code        string
	FoodGroupID string
	Description string
}{
	ID:          "usda_retention_factor.id",
	Code:        "usda_retention_factor.code",
	FoodGroupID: "usda_retention_factor.food_group_id",
	Description: "usda_retention_factor.description",
}

// Generated where

var UsdaRetentionFactorWhere = struct {
	ID          whereHelperint
	Code        whereHelpernull_String
	FoodGroupID whereHelpernull_Int
	Description whereHelpernull_String
}{
	ID:          whereHelperint{field: "\"usda_retention_factor\".\"id\""},
	Code:        whereHelpernull_String{field: "\"usda_retention_factor\".\"code\""},
	FoodGroupID: whereHelpernull_Int{field: "\"usda_retention_factor\".\"food_group_id\""},
	Description: whereHelpernull_String{field: "\"usda_retention_factor\".\"description\""},
}

// UsdaRetentionFactorRels is where relationship names are stored.
var UsdaRetentionFactorRels = struct {
	FoodGroup string
}{
	FoodGroup: "FoodGroup",
}

// usdaRetentionFactorR is where relationships are stored.
type usdaRetentionFactorR struct {
	FoodGroup *UsdaFoodCategory `boil:"FoodGroup" json:"FoodGroup" toml:"FoodGroup" yaml:"FoodGroup"`
}

// NewStruct creates a new relationship struct
func (*usdaRetentionFactorR) NewStruct() *usdaRetentionFactorR {
	return &usdaRetentionFactorR{}
}

// usdaRetentionFactorL is where Load methods for each relationship are stored.
type usdaRetentionFactorL struct{}

var (
	usdaRetentionFactorAllColumns            = []string{"id", "code", "food_group_id", "description"}
	usdaRetentionFactorColumnsWithoutDefault = []string{"id"}
	usdaRetentionFactorColumnsWithDefault    = []string{"code", "food_group_id", "description"}
	usdaRetentionFactorPrimaryKeyColumns     = []string{"id"}
	usdaRetentionFactorGeneratedColumns      = []string{}
)

type (
	// UsdaRetentionFactorSlice is an alias for a slice of pointers to UsdaRetentionFactor.
	// This should almost always be used instead of []UsdaRetentionFactor.
	UsdaRetentionFactorSlice []*UsdaRetentionFactor
	// UsdaRetentionFactorHook is the signature for custom UsdaRetentionFactor hook methods
	UsdaRetentionFactorHook func(context.Context, boil.ContextExecutor, *UsdaRetentionFactor) error

	usdaRetentionFactorQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	usdaRetentionFactorType                 = reflect.TypeOf(&UsdaRetentionFactor{})
	usdaRetentionFactorMapping              = queries.MakeStructMapping(usdaRetentionFactorType)
	usdaRetentionFactorPrimaryKeyMapping, _ = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, usdaRetentionFactorPrimaryKeyColumns)
	usdaRetentionFactorInsertCacheMut       sync.RWMutex
	usdaRetentionFactorInsertCache          = make(map[string]insertCache)
	usdaRetentionFactorUpdateCacheMut       sync.RWMutex
	usdaRetentionFactorUpdateCache          = make(map[string]updateCache)
	usdaRetentionFactorUpsertCacheMut       sync.RWMutex
	usdaRetentionFactorUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var usdaRetentionFactorAfterSelectHooks []UsdaRetentionFactorHook

var usdaRetentionFactorBeforeInsertHooks []UsdaRetentionFactorHook
var usdaRetentionFactorAfterInsertHooks []UsdaRetentionFactorHook

var usdaRetentionFactorBeforeUpdateHooks []UsdaRetentionFactorHook
var usdaRetentionFactorAfterUpdateHooks []UsdaRetentionFactorHook

var usdaRetentionFactorBeforeDeleteHooks []UsdaRetentionFactorHook
var usdaRetentionFactorAfterDeleteHooks []UsdaRetentionFactorHook

var usdaRetentionFactorBeforeUpsertHooks []UsdaRetentionFactorHook
var usdaRetentionFactorAfterUpsertHooks []UsdaRetentionFactorHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UsdaRetentionFactor) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UsdaRetentionFactor) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UsdaRetentionFactor) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UsdaRetentionFactor) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UsdaRetentionFactor) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UsdaRetentionFactor) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UsdaRetentionFactor) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UsdaRetentionFactor) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UsdaRetentionFactor) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaRetentionFactorAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUsdaRetentionFactorHook registers your hook function for all future operations.
func AddUsdaRetentionFactorHook(hookPoint boil.HookPoint, usdaRetentionFactorHook UsdaRetentionFactorHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		usdaRetentionFactorAfterSelectHooks = append(usdaRetentionFactorAfterSelectHooks, usdaRetentionFactorHook)
	case boil.BeforeInsertHook:
		usdaRetentionFactorBeforeInsertHooks = append(usdaRetentionFactorBeforeInsertHooks, usdaRetentionFactorHook)
	case boil.AfterInsertHook:
		usdaRetentionFactorAfterInsertHooks = append(usdaRetentionFactorAfterInsertHooks, usdaRetentionFactorHook)
	case boil.BeforeUpdateHook:
		usdaRetentionFactorBeforeUpdateHooks = append(usdaRetentionFactorBeforeUpdateHooks, usdaRetentionFactorHook)
	case boil.AfterUpdateHook:
		usdaRetentionFactorAfterUpdateHooks = append(usdaRetentionFactorAfterUpdateHooks, usdaRetentionFactorHook)
	case boil.BeforeDeleteHook:
		usdaRetentionFactorBeforeDeleteHooks = append(usdaRetentionFactorBeforeDeleteHooks, usdaRetentionFactorHook)
	case boil.AfterDeleteHook:
		usdaRetentionFactorAfterDeleteHooks = append(usdaRetentionFactorAfterDeleteHooks, usdaRetentionFactorHook)
	case boil.BeforeUpsertHook:
		usdaRetentionFactorBeforeUpsertHooks = append(usdaRetentionFactorBeforeUpsertHooks, usdaRetentionFactorHook)
	case boil.AfterUpsertHook:
		usdaRetentionFactorAfterUpsertHooks = append(usdaRetentionFactorAfterUpsertHooks, usdaRetentionFactorHook)
	}
}

// One returns a single usdaRetentionFactor record from the query.
func (q usdaRetentionFactorQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UsdaRetentionFactor, error) {
	o := &UsdaRetentionFactor{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for usda_retention_factor")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UsdaRetentionFactor records from the query.
func (q usdaRetentionFactorQuery) All(ctx context.Context, exec boil.ContextExecutor) (UsdaRetentionFactorSlice, error) {
	var o []*UsdaRetentionFactor

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UsdaRetentionFactor slice")
	}

	if len(usdaRetentionFactorAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UsdaRetentionFactor records in the query.
func (q usdaRetentionFactorQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count usda_retention_factor rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q usdaRetentionFactorQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if usda_retention_factor exists")
	}

	return count > 0, nil
}

// FoodGroup pointed to by the foreign key.
func (o *UsdaRetentionFactor) FoodGroup(mods ...qm.QueryMod) usdaFoodCategoryQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.FoodGroupID),
	}

	queryMods = append(queryMods, mods...)

	return UsdaFoodCategories(queryMods...)
}

// LoadFoodGroup allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (usdaRetentionFactorL) LoadFoodGroup(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUsdaRetentionFactor interface{}, mods queries.Applicator) error {
	var slice []*UsdaRetentionFactor
	var object *UsdaRetentionFactor

	if singular {
		object = maybeUsdaRetentionFactor.(*UsdaRetentionFactor)
	} else {
		slice = *maybeUsdaRetentionFactor.(*[]*UsdaRetentionFactor)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &usdaRetentionFactorR{}
		}
		if !queries.IsNil(object.FoodGroupID) {
			args = append(args, object.FoodGroupID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &usdaRetentionFactorR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.FoodGroupID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.FoodGroupID) {
				args = append(args, obj.FoodGroupID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`usda_food_category`),
		qm.WhereIn(`usda_food_category.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UsdaFoodCategory")
	}

	var resultSlice []*UsdaFoodCategory
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UsdaFoodCategory")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for usda_food_category")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for usda_food_category")
	}

	if len(usdaRetentionFactorAfterSelectHooks) != 0 {
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
		object.R.FoodGroup = foreign
		if foreign.R == nil {
			foreign.R = &usdaFoodCategoryR{}
		}
		foreign.R.FoodGroupUsdaRetentionFactors = append(foreign.R.FoodGroupUsdaRetentionFactors, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.FoodGroupID, foreign.ID) {
				local.R.FoodGroup = foreign
				if foreign.R == nil {
					foreign.R = &usdaFoodCategoryR{}
				}
				foreign.R.FoodGroupUsdaRetentionFactors = append(foreign.R.FoodGroupUsdaRetentionFactors, local)
				break
			}
		}
	}

	return nil
}

// SetFoodGroup of the usdaRetentionFactor to the related item.
// Sets o.R.FoodGroup to related.
// Adds o to related.R.FoodGroupUsdaRetentionFactors.
func (o *UsdaRetentionFactor) SetFoodGroup(ctx context.Context, exec boil.ContextExecutor, insert bool, related *UsdaFoodCategory) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"usda_retention_factor\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"food_group_id"}),
		strmangle.WhereClause("\"", "\"", 2, usdaRetentionFactorPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.FoodGroupID, related.ID)
	if o.R == nil {
		o.R = &usdaRetentionFactorR{
			FoodGroup: related,
		}
	} else {
		o.R.FoodGroup = related
	}

	if related.R == nil {
		related.R = &usdaFoodCategoryR{
			FoodGroupUsdaRetentionFactors: UsdaRetentionFactorSlice{o},
		}
	} else {
		related.R.FoodGroupUsdaRetentionFactors = append(related.R.FoodGroupUsdaRetentionFactors, o)
	}

	return nil
}

// RemoveFoodGroup relationship.
// Sets o.R.FoodGroup to nil.
// Removes o from all passed in related items' relationships struct.
func (o *UsdaRetentionFactor) RemoveFoodGroup(ctx context.Context, exec boil.ContextExecutor, related *UsdaFoodCategory) error {
	var err error

	queries.SetScanner(&o.FoodGroupID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("food_group_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.FoodGroup = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.FoodGroupUsdaRetentionFactors {
		if queries.Equal(o.FoodGroupID, ri.FoodGroupID) {
			continue
		}

		ln := len(related.R.FoodGroupUsdaRetentionFactors)
		if ln > 1 && i < ln-1 {
			related.R.FoodGroupUsdaRetentionFactors[i] = related.R.FoodGroupUsdaRetentionFactors[ln-1]
		}
		related.R.FoodGroupUsdaRetentionFactors = related.R.FoodGroupUsdaRetentionFactors[:ln-1]
		break
	}
	return nil
}

// UsdaRetentionFactors retrieves all the records using an executor.
func UsdaRetentionFactors(mods ...qm.QueryMod) usdaRetentionFactorQuery {
	mods = append(mods, qm.From("\"usda_retention_factor\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"usda_retention_factor\".*"})
	}

	return usdaRetentionFactorQuery{q}
}

// FindUsdaRetentionFactor retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUsdaRetentionFactor(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*UsdaRetentionFactor, error) {
	usdaRetentionFactorObj := &UsdaRetentionFactor{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"usda_retention_factor\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, usdaRetentionFactorObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from usda_retention_factor")
	}

	if err = usdaRetentionFactorObj.doAfterSelectHooks(ctx, exec); err != nil {
		return usdaRetentionFactorObj, err
	}

	return usdaRetentionFactorObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UsdaRetentionFactor) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_retention_factor provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaRetentionFactorColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	usdaRetentionFactorInsertCacheMut.RLock()
	cache, cached := usdaRetentionFactorInsertCache[key]
	usdaRetentionFactorInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			usdaRetentionFactorAllColumns,
			usdaRetentionFactorColumnsWithDefault,
			usdaRetentionFactorColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"usda_retention_factor\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"usda_retention_factor\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into usda_retention_factor")
	}

	if !cached {
		usdaRetentionFactorInsertCacheMut.Lock()
		usdaRetentionFactorInsertCache[key] = cache
		usdaRetentionFactorInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UsdaRetentionFactor.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UsdaRetentionFactor) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	usdaRetentionFactorUpdateCacheMut.RLock()
	cache, cached := usdaRetentionFactorUpdateCache[key]
	usdaRetentionFactorUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			usdaRetentionFactorAllColumns,
			usdaRetentionFactorPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update usda_retention_factor, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"usda_retention_factor\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, usdaRetentionFactorPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, append(wl, usdaRetentionFactorPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update usda_retention_factor row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for usda_retention_factor")
	}

	if !cached {
		usdaRetentionFactorUpdateCacheMut.Lock()
		usdaRetentionFactorUpdateCache[key] = cache
		usdaRetentionFactorUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q usdaRetentionFactorQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for usda_retention_factor")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for usda_retention_factor")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UsdaRetentionFactorSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaRetentionFactorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"usda_retention_factor\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, usdaRetentionFactorPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in usdaRetentionFactor slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all usdaRetentionFactor")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UsdaRetentionFactor) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_retention_factor provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaRetentionFactorColumnsWithDefault, o)

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

	usdaRetentionFactorUpsertCacheMut.RLock()
	cache, cached := usdaRetentionFactorUpsertCache[key]
	usdaRetentionFactorUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			usdaRetentionFactorAllColumns,
			usdaRetentionFactorColumnsWithDefault,
			usdaRetentionFactorColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			usdaRetentionFactorAllColumns,
			usdaRetentionFactorPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert usda_retention_factor, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(usdaRetentionFactorPrimaryKeyColumns))
			copy(conflict, usdaRetentionFactorPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"usda_retention_factor\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(usdaRetentionFactorType, usdaRetentionFactorMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert usda_retention_factor")
	}

	if !cached {
		usdaRetentionFactorUpsertCacheMut.Lock()
		usdaRetentionFactorUpsertCache[key] = cache
		usdaRetentionFactorUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UsdaRetentionFactor record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UsdaRetentionFactor) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UsdaRetentionFactor provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), usdaRetentionFactorPrimaryKeyMapping)
	sql := "DELETE FROM \"usda_retention_factor\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from usda_retention_factor")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for usda_retention_factor")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q usdaRetentionFactorQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no usdaRetentionFactorQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usda_retention_factor")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_retention_factor")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UsdaRetentionFactorSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(usdaRetentionFactorBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaRetentionFactorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"usda_retention_factor\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaRetentionFactorPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usdaRetentionFactor slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_retention_factor")
	}

	if len(usdaRetentionFactorAfterDeleteHooks) != 0 {
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
func (o *UsdaRetentionFactor) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUsdaRetentionFactor(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsdaRetentionFactorSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UsdaRetentionFactorSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaRetentionFactorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"usda_retention_factor\".* FROM \"usda_retention_factor\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaRetentionFactorPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UsdaRetentionFactorSlice")
	}

	*o = slice

	return nil
}

// UsdaRetentionFactorExists checks if the UsdaRetentionFactor row exists.
func UsdaRetentionFactorExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"usda_retention_factor\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if usda_retention_factor exists")
	}

	return exists, nil
}
