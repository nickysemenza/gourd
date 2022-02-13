// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// UsdaFoodCategory is an object representing the database table.
type UsdaFoodCategory struct {
	ID          int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Code        null.String `boil:"code" json:"code,omitempty" toml:"code" yaml:"code,omitempty"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`

	R *usdaFoodCategoryR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L usdaFoodCategoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UsdaFoodCategoryColumns = struct {
	ID          string
	Code        string
	Description string
}{
	ID:          "id",
	Code:        "code",
	Description: "description",
}

var UsdaFoodCategoryTableColumns = struct {
	ID          string
	Code        string
	Description string
}{
	ID:          "usda_food_category.id",
	Code:        "usda_food_category.code",
	Description: "usda_food_category.description",
}

// Generated where

var UsdaFoodCategoryWhere = struct {
	ID          whereHelperint
	Code        whereHelpernull_String
	Description whereHelpernull_String
}{
	ID:          whereHelperint{field: "\"usda_food_category\".\"id\""},
	Code:        whereHelpernull_String{field: "\"usda_food_category\".\"code\""},
	Description: whereHelpernull_String{field: "\"usda_food_category\".\"description\""},
}

// UsdaFoodCategoryRels is where relationship names are stored.
var UsdaFoodCategoryRels = struct {
	FoodGroupUsdaRetentionFactors string
}{
	FoodGroupUsdaRetentionFactors: "FoodGroupUsdaRetentionFactors",
}

// usdaFoodCategoryR is where relationships are stored.
type usdaFoodCategoryR struct {
	FoodGroupUsdaRetentionFactors UsdaRetentionFactorSlice `boil:"FoodGroupUsdaRetentionFactors" json:"FoodGroupUsdaRetentionFactors" toml:"FoodGroupUsdaRetentionFactors" yaml:"FoodGroupUsdaRetentionFactors"`
}

// NewStruct creates a new relationship struct
func (*usdaFoodCategoryR) NewStruct() *usdaFoodCategoryR {
	return &usdaFoodCategoryR{}
}

// usdaFoodCategoryL is where Load methods for each relationship are stored.
type usdaFoodCategoryL struct{}

var (
	usdaFoodCategoryAllColumns            = []string{"id", "code", "description"}
	usdaFoodCategoryColumnsWithoutDefault = []string{"id", "code", "description"}
	usdaFoodCategoryColumnsWithDefault    = []string{}
	usdaFoodCategoryPrimaryKeyColumns     = []string{"id"}
	usdaFoodCategoryGeneratedColumns      = []string{}
)

type (
	// UsdaFoodCategorySlice is an alias for a slice of pointers to UsdaFoodCategory.
	// This should almost always be used instead of []UsdaFoodCategory.
	UsdaFoodCategorySlice []*UsdaFoodCategory
	// UsdaFoodCategoryHook is the signature for custom UsdaFoodCategory hook methods
	UsdaFoodCategoryHook func(context.Context, boil.ContextExecutor, *UsdaFoodCategory) error

	usdaFoodCategoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	usdaFoodCategoryType                 = reflect.TypeOf(&UsdaFoodCategory{})
	usdaFoodCategoryMapping              = queries.MakeStructMapping(usdaFoodCategoryType)
	usdaFoodCategoryPrimaryKeyMapping, _ = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, usdaFoodCategoryPrimaryKeyColumns)
	usdaFoodCategoryInsertCacheMut       sync.RWMutex
	usdaFoodCategoryInsertCache          = make(map[string]insertCache)
	usdaFoodCategoryUpdateCacheMut       sync.RWMutex
	usdaFoodCategoryUpdateCache          = make(map[string]updateCache)
	usdaFoodCategoryUpsertCacheMut       sync.RWMutex
	usdaFoodCategoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var usdaFoodCategoryAfterSelectHooks []UsdaFoodCategoryHook

var usdaFoodCategoryBeforeInsertHooks []UsdaFoodCategoryHook
var usdaFoodCategoryAfterInsertHooks []UsdaFoodCategoryHook

var usdaFoodCategoryBeforeUpdateHooks []UsdaFoodCategoryHook
var usdaFoodCategoryAfterUpdateHooks []UsdaFoodCategoryHook

var usdaFoodCategoryBeforeDeleteHooks []UsdaFoodCategoryHook
var usdaFoodCategoryAfterDeleteHooks []UsdaFoodCategoryHook

var usdaFoodCategoryBeforeUpsertHooks []UsdaFoodCategoryHook
var usdaFoodCategoryAfterUpsertHooks []UsdaFoodCategoryHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UsdaFoodCategory) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UsdaFoodCategory) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UsdaFoodCategory) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UsdaFoodCategory) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UsdaFoodCategory) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UsdaFoodCategory) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UsdaFoodCategory) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UsdaFoodCategory) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UsdaFoodCategory) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usdaFoodCategoryAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUsdaFoodCategoryHook registers your hook function for all future operations.
func AddUsdaFoodCategoryHook(hookPoint boil.HookPoint, usdaFoodCategoryHook UsdaFoodCategoryHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		usdaFoodCategoryAfterSelectHooks = append(usdaFoodCategoryAfterSelectHooks, usdaFoodCategoryHook)
	case boil.BeforeInsertHook:
		usdaFoodCategoryBeforeInsertHooks = append(usdaFoodCategoryBeforeInsertHooks, usdaFoodCategoryHook)
	case boil.AfterInsertHook:
		usdaFoodCategoryAfterInsertHooks = append(usdaFoodCategoryAfterInsertHooks, usdaFoodCategoryHook)
	case boil.BeforeUpdateHook:
		usdaFoodCategoryBeforeUpdateHooks = append(usdaFoodCategoryBeforeUpdateHooks, usdaFoodCategoryHook)
	case boil.AfterUpdateHook:
		usdaFoodCategoryAfterUpdateHooks = append(usdaFoodCategoryAfterUpdateHooks, usdaFoodCategoryHook)
	case boil.BeforeDeleteHook:
		usdaFoodCategoryBeforeDeleteHooks = append(usdaFoodCategoryBeforeDeleteHooks, usdaFoodCategoryHook)
	case boil.AfterDeleteHook:
		usdaFoodCategoryAfterDeleteHooks = append(usdaFoodCategoryAfterDeleteHooks, usdaFoodCategoryHook)
	case boil.BeforeUpsertHook:
		usdaFoodCategoryBeforeUpsertHooks = append(usdaFoodCategoryBeforeUpsertHooks, usdaFoodCategoryHook)
	case boil.AfterUpsertHook:
		usdaFoodCategoryAfterUpsertHooks = append(usdaFoodCategoryAfterUpsertHooks, usdaFoodCategoryHook)
	}
}

// One returns a single usdaFoodCategory record from the query.
func (q usdaFoodCategoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UsdaFoodCategory, error) {
	o := &UsdaFoodCategory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for usda_food_category")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UsdaFoodCategory records from the query.
func (q usdaFoodCategoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (UsdaFoodCategorySlice, error) {
	var o []*UsdaFoodCategory

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UsdaFoodCategory slice")
	}

	if len(usdaFoodCategoryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UsdaFoodCategory records in the query.
func (q usdaFoodCategoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count usda_food_category rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q usdaFoodCategoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if usda_food_category exists")
	}

	return count > 0, nil
}

// FoodGroupUsdaRetentionFactors retrieves all the usda_retention_factor's UsdaRetentionFactors with an executor via food_group_id column.
func (o *UsdaFoodCategory) FoodGroupUsdaRetentionFactors(mods ...qm.QueryMod) usdaRetentionFactorQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"usda_retention_factor\".\"food_group_id\"=?", o.ID),
	)

	query := UsdaRetentionFactors(queryMods...)
	queries.SetFrom(query.Query, "\"usda_retention_factor\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"usda_retention_factor\".*"})
	}

	return query
}

// LoadFoodGroupUsdaRetentionFactors allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (usdaFoodCategoryL) LoadFoodGroupUsdaRetentionFactors(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUsdaFoodCategory interface{}, mods queries.Applicator) error {
	var slice []*UsdaFoodCategory
	var object *UsdaFoodCategory

	if singular {
		object = maybeUsdaFoodCategory.(*UsdaFoodCategory)
	} else {
		slice = *maybeUsdaFoodCategory.(*[]*UsdaFoodCategory)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &usdaFoodCategoryR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &usdaFoodCategoryR{}
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
		qm.From(`usda_retention_factor`),
		qm.WhereIn(`usda_retention_factor.food_group_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load usda_retention_factor")
	}

	var resultSlice []*UsdaRetentionFactor
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice usda_retention_factor")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on usda_retention_factor")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for usda_retention_factor")
	}

	if len(usdaRetentionFactorAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.FoodGroupUsdaRetentionFactors = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &usdaRetentionFactorR{}
			}
			foreign.R.FoodGroup = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.FoodGroupID) {
				local.R.FoodGroupUsdaRetentionFactors = append(local.R.FoodGroupUsdaRetentionFactors, foreign)
				if foreign.R == nil {
					foreign.R = &usdaRetentionFactorR{}
				}
				foreign.R.FoodGroup = local
				break
			}
		}
	}

	return nil
}

// AddFoodGroupUsdaRetentionFactors adds the given related objects to the existing relationships
// of the usda_food_category, optionally inserting them as new records.
// Appends related to o.R.FoodGroupUsdaRetentionFactors.
// Sets related.R.FoodGroup appropriately.
func (o *UsdaFoodCategory) AddFoodGroupUsdaRetentionFactors(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UsdaRetentionFactor) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.FoodGroupID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"usda_retention_factor\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"food_group_id"}),
				strmangle.WhereClause("\"", "\"", 2, usdaRetentionFactorPrimaryKeyColumns),
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

			queries.Assign(&rel.FoodGroupID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &usdaFoodCategoryR{
			FoodGroupUsdaRetentionFactors: related,
		}
	} else {
		o.R.FoodGroupUsdaRetentionFactors = append(o.R.FoodGroupUsdaRetentionFactors, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &usdaRetentionFactorR{
				FoodGroup: o,
			}
		} else {
			rel.R.FoodGroup = o
		}
	}
	return nil
}

// SetFoodGroupUsdaRetentionFactors removes all previously related items of the
// usda_food_category replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.FoodGroup's FoodGroupUsdaRetentionFactors accordingly.
// Replaces o.R.FoodGroupUsdaRetentionFactors with related.
// Sets related.R.FoodGroup's FoodGroupUsdaRetentionFactors accordingly.
func (o *UsdaFoodCategory) SetFoodGroupUsdaRetentionFactors(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UsdaRetentionFactor) error {
	query := "update \"usda_retention_factor\" set \"food_group_id\" = null where \"food_group_id\" = $1"
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
		for _, rel := range o.R.FoodGroupUsdaRetentionFactors {
			queries.SetScanner(&rel.FoodGroupID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.FoodGroup = nil
		}

		o.R.FoodGroupUsdaRetentionFactors = nil
	}
	return o.AddFoodGroupUsdaRetentionFactors(ctx, exec, insert, related...)
}

// RemoveFoodGroupUsdaRetentionFactors relationships from objects passed in.
// Removes related items from R.FoodGroupUsdaRetentionFactors (uses pointer comparison, removal does not keep order)
// Sets related.R.FoodGroup.
func (o *UsdaFoodCategory) RemoveFoodGroupUsdaRetentionFactors(ctx context.Context, exec boil.ContextExecutor, related ...*UsdaRetentionFactor) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.FoodGroupID, nil)
		if rel.R != nil {
			rel.R.FoodGroup = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("food_group_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.FoodGroupUsdaRetentionFactors {
			if rel != ri {
				continue
			}

			ln := len(o.R.FoodGroupUsdaRetentionFactors)
			if ln > 1 && i < ln-1 {
				o.R.FoodGroupUsdaRetentionFactors[i] = o.R.FoodGroupUsdaRetentionFactors[ln-1]
			}
			o.R.FoodGroupUsdaRetentionFactors = o.R.FoodGroupUsdaRetentionFactors[:ln-1]
			break
		}
	}

	return nil
}

// UsdaFoodCategories retrieves all the records using an executor.
func UsdaFoodCategories(mods ...qm.QueryMod) usdaFoodCategoryQuery {
	mods = append(mods, qm.From("\"usda_food_category\""))
	return usdaFoodCategoryQuery{NewQuery(mods...)}
}

// FindUsdaFoodCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUsdaFoodCategory(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*UsdaFoodCategory, error) {
	usdaFoodCategoryObj := &UsdaFoodCategory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"usda_food_category\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, usdaFoodCategoryObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from usda_food_category")
	}

	if err = usdaFoodCategoryObj.doAfterSelectHooks(ctx, exec); err != nil {
		return usdaFoodCategoryObj, err
	}

	return usdaFoodCategoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UsdaFoodCategory) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_food_category provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaFoodCategoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	usdaFoodCategoryInsertCacheMut.RLock()
	cache, cached := usdaFoodCategoryInsertCache[key]
	usdaFoodCategoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			usdaFoodCategoryAllColumns,
			usdaFoodCategoryColumnsWithDefault,
			usdaFoodCategoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"usda_food_category\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"usda_food_category\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into usda_food_category")
	}

	if !cached {
		usdaFoodCategoryInsertCacheMut.Lock()
		usdaFoodCategoryInsertCache[key] = cache
		usdaFoodCategoryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UsdaFoodCategory.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UsdaFoodCategory) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	usdaFoodCategoryUpdateCacheMut.RLock()
	cache, cached := usdaFoodCategoryUpdateCache[key]
	usdaFoodCategoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			usdaFoodCategoryAllColumns,
			usdaFoodCategoryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update usda_food_category, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"usda_food_category\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, usdaFoodCategoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, append(wl, usdaFoodCategoryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update usda_food_category row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for usda_food_category")
	}

	if !cached {
		usdaFoodCategoryUpdateCacheMut.Lock()
		usdaFoodCategoryUpdateCache[key] = cache
		usdaFoodCategoryUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q usdaFoodCategoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for usda_food_category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for usda_food_category")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UsdaFoodCategorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaFoodCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"usda_food_category\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, usdaFoodCategoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in usdaFoodCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all usdaFoodCategory")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UsdaFoodCategory) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no usda_food_category provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usdaFoodCategoryColumnsWithDefault, o)

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

	usdaFoodCategoryUpsertCacheMut.RLock()
	cache, cached := usdaFoodCategoryUpsertCache[key]
	usdaFoodCategoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			usdaFoodCategoryAllColumns,
			usdaFoodCategoryColumnsWithDefault,
			usdaFoodCategoryColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			usdaFoodCategoryAllColumns,
			usdaFoodCategoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert usda_food_category, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(usdaFoodCategoryPrimaryKeyColumns))
			copy(conflict, usdaFoodCategoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"usda_food_category\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(usdaFoodCategoryType, usdaFoodCategoryMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert usda_food_category")
	}

	if !cached {
		usdaFoodCategoryUpsertCacheMut.Lock()
		usdaFoodCategoryUpsertCache[key] = cache
		usdaFoodCategoryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UsdaFoodCategory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UsdaFoodCategory) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UsdaFoodCategory provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), usdaFoodCategoryPrimaryKeyMapping)
	sql := "DELETE FROM \"usda_food_category\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from usda_food_category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for usda_food_category")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q usdaFoodCategoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no usdaFoodCategoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usda_food_category")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_food_category")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UsdaFoodCategorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(usdaFoodCategoryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaFoodCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"usda_food_category\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaFoodCategoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from usdaFoodCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for usda_food_category")
	}

	if len(usdaFoodCategoryAfterDeleteHooks) != 0 {
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
func (o *UsdaFoodCategory) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUsdaFoodCategory(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsdaFoodCategorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UsdaFoodCategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usdaFoodCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"usda_food_category\".* FROM \"usda_food_category\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usdaFoodCategoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UsdaFoodCategorySlice")
	}

	*o = slice

	return nil
}

// UsdaFoodCategoryExists checks if the UsdaFoodCategory row exists.
func UsdaFoodCategoryExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"usda_food_category\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if usda_food_category exists")
	}

	return exists, nil
}
