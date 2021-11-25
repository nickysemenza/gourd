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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// MealRecipe is an object representing the database table.
type MealRecipe struct {
	MealID     string            `boil:"meal_id" json:"meal_id" toml:"meal_id" yaml:"meal_id"`
	RecipeID   string            `boil:"recipe_id" json:"recipe_id" toml:"recipe_id" yaml:"recipe_id"`
	Multiplier types.NullDecimal `boil:"multiplier" json:"multiplier,omitempty" toml:"multiplier" yaml:"multiplier,omitempty"`

	R *mealRecipeR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L mealRecipeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MealRecipeColumns = struct {
	MealID     string
	RecipeID   string
	Multiplier string
}{
	MealID:     "meal_id",
	RecipeID:   "recipe_id",
	Multiplier: "multiplier",
}

var MealRecipeTableColumns = struct {
	MealID     string
	RecipeID   string
	Multiplier string
}{
	MealID:     "meal_recipe.meal_id",
	RecipeID:   "meal_recipe.recipe_id",
	Multiplier: "meal_recipe.multiplier",
}

// Generated where

type whereHelpertypes_NullDecimal struct{ field string }

func (w whereHelpertypes_NullDecimal) EQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_NullDecimal) NEQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_NullDecimal) LT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_NullDecimal) LTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_NullDecimal) GT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_NullDecimal) GTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpertypes_NullDecimal) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_NullDecimal) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

var MealRecipeWhere = struct {
	MealID     whereHelperstring
	RecipeID   whereHelperstring
	Multiplier whereHelpertypes_NullDecimal
}{
	MealID:     whereHelperstring{field: "\"meal_recipe\".\"meal_id\""},
	RecipeID:   whereHelperstring{field: "\"meal_recipe\".\"recipe_id\""},
	Multiplier: whereHelpertypes_NullDecimal{field: "\"meal_recipe\".\"multiplier\""},
}

// MealRecipeRels is where relationship names are stored.
var MealRecipeRels = struct {
	Meal   string
	Recipe string
}{
	Meal:   "Meal",
	Recipe: "Recipe",
}

// mealRecipeR is where relationships are stored.
type mealRecipeR struct {
	Meal   *Meal   `boil:"Meal" json:"Meal" toml:"Meal" yaml:"Meal"`
	Recipe *Recipe `boil:"Recipe" json:"Recipe" toml:"Recipe" yaml:"Recipe"`
}

// NewStruct creates a new relationship struct
func (*mealRecipeR) NewStruct() *mealRecipeR {
	return &mealRecipeR{}
}

// mealRecipeL is where Load methods for each relationship are stored.
type mealRecipeL struct{}

var (
	mealRecipeAllColumns            = []string{"meal_id", "recipe_id", "multiplier"}
	mealRecipeColumnsWithoutDefault = []string{"meal_id", "recipe_id"}
	mealRecipeColumnsWithDefault    = []string{"multiplier"}
	mealRecipePrimaryKeyColumns     = []string{"meal_id", "recipe_id"}
)

type (
	// MealRecipeSlice is an alias for a slice of pointers to MealRecipe.
	// This should almost always be used instead of []MealRecipe.
	MealRecipeSlice []*MealRecipe
	// MealRecipeHook is the signature for custom MealRecipe hook methods
	MealRecipeHook func(context.Context, boil.ContextExecutor, *MealRecipe) error

	mealRecipeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	mealRecipeType                 = reflect.TypeOf(&MealRecipe{})
	mealRecipeMapping              = queries.MakeStructMapping(mealRecipeType)
	mealRecipePrimaryKeyMapping, _ = queries.BindMapping(mealRecipeType, mealRecipeMapping, mealRecipePrimaryKeyColumns)
	mealRecipeInsertCacheMut       sync.RWMutex
	mealRecipeInsertCache          = make(map[string]insertCache)
	mealRecipeUpdateCacheMut       sync.RWMutex
	mealRecipeUpdateCache          = make(map[string]updateCache)
	mealRecipeUpsertCacheMut       sync.RWMutex
	mealRecipeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var mealRecipeBeforeInsertHooks []MealRecipeHook
var mealRecipeBeforeUpdateHooks []MealRecipeHook
var mealRecipeBeforeDeleteHooks []MealRecipeHook
var mealRecipeBeforeUpsertHooks []MealRecipeHook

var mealRecipeAfterInsertHooks []MealRecipeHook
var mealRecipeAfterSelectHooks []MealRecipeHook
var mealRecipeAfterUpdateHooks []MealRecipeHook
var mealRecipeAfterDeleteHooks []MealRecipeHook
var mealRecipeAfterUpsertHooks []MealRecipeHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *MealRecipe) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *MealRecipe) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *MealRecipe) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *MealRecipe) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *MealRecipe) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *MealRecipe) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *MealRecipe) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *MealRecipe) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *MealRecipe) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealRecipeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMealRecipeHook registers your hook function for all future operations.
func AddMealRecipeHook(hookPoint boil.HookPoint, mealRecipeHook MealRecipeHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		mealRecipeBeforeInsertHooks = append(mealRecipeBeforeInsertHooks, mealRecipeHook)
	case boil.BeforeUpdateHook:
		mealRecipeBeforeUpdateHooks = append(mealRecipeBeforeUpdateHooks, mealRecipeHook)
	case boil.BeforeDeleteHook:
		mealRecipeBeforeDeleteHooks = append(mealRecipeBeforeDeleteHooks, mealRecipeHook)
	case boil.BeforeUpsertHook:
		mealRecipeBeforeUpsertHooks = append(mealRecipeBeforeUpsertHooks, mealRecipeHook)
	case boil.AfterInsertHook:
		mealRecipeAfterInsertHooks = append(mealRecipeAfterInsertHooks, mealRecipeHook)
	case boil.AfterSelectHook:
		mealRecipeAfterSelectHooks = append(mealRecipeAfterSelectHooks, mealRecipeHook)
	case boil.AfterUpdateHook:
		mealRecipeAfterUpdateHooks = append(mealRecipeAfterUpdateHooks, mealRecipeHook)
	case boil.AfterDeleteHook:
		mealRecipeAfterDeleteHooks = append(mealRecipeAfterDeleteHooks, mealRecipeHook)
	case boil.AfterUpsertHook:
		mealRecipeAfterUpsertHooks = append(mealRecipeAfterUpsertHooks, mealRecipeHook)
	}
}

// One returns a single mealRecipe record from the query.
func (q mealRecipeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*MealRecipe, error) {
	o := &MealRecipe{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for meal_recipe")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all MealRecipe records from the query.
func (q mealRecipeQuery) All(ctx context.Context, exec boil.ContextExecutor) (MealRecipeSlice, error) {
	var o []*MealRecipe

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to MealRecipe slice")
	}

	if len(mealRecipeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all MealRecipe records in the query.
func (q mealRecipeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count meal_recipe rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q mealRecipeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if meal_recipe exists")
	}

	return count > 0, nil
}

// Meal pointed to by the foreign key.
func (o *MealRecipe) Meal(mods ...qm.QueryMod) mealQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.MealID),
	}

	queryMods = append(queryMods, mods...)

	query := Meals(queryMods...)
	queries.SetFrom(query.Query, "\"meals\"")

	return query
}

// Recipe pointed to by the foreign key.
func (o *MealRecipe) Recipe(mods ...qm.QueryMod) recipeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.RecipeID),
	}

	queryMods = append(queryMods, mods...)

	query := Recipes(queryMods...)
	queries.SetFrom(query.Query, "\"recipes\"")

	return query
}

// LoadMeal allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (mealRecipeL) LoadMeal(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMealRecipe interface{}, mods queries.Applicator) error {
	var slice []*MealRecipe
	var object *MealRecipe

	if singular {
		object = maybeMealRecipe.(*MealRecipe)
	} else {
		slice = *maybeMealRecipe.(*[]*MealRecipe)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &mealRecipeR{}
		}
		args = append(args, object.MealID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &mealRecipeR{}
			}

			for _, a := range args {
				if a == obj.MealID {
					continue Outer
				}
			}

			args = append(args, obj.MealID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`meals`),
		qm.WhereIn(`meals.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Meal")
	}

	var resultSlice []*Meal
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Meal")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for meals")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for meals")
	}

	if len(mealRecipeAfterSelectHooks) != 0 {
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
		object.R.Meal = foreign
		if foreign.R == nil {
			foreign.R = &mealR{}
		}
		foreign.R.MealRecipes = append(foreign.R.MealRecipes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.MealID == foreign.ID {
				local.R.Meal = foreign
				if foreign.R == nil {
					foreign.R = &mealR{}
				}
				foreign.R.MealRecipes = append(foreign.R.MealRecipes, local)
				break
			}
		}
	}

	return nil
}

// LoadRecipe allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (mealRecipeL) LoadRecipe(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMealRecipe interface{}, mods queries.Applicator) error {
	var slice []*MealRecipe
	var object *MealRecipe

	if singular {
		object = maybeMealRecipe.(*MealRecipe)
	} else {
		slice = *maybeMealRecipe.(*[]*MealRecipe)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &mealRecipeR{}
		}
		args = append(args, object.RecipeID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &mealRecipeR{}
			}

			for _, a := range args {
				if a == obj.RecipeID {
					continue Outer
				}
			}

			args = append(args, obj.RecipeID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`recipes`),
		qm.WhereIn(`recipes.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Recipe")
	}

	var resultSlice []*Recipe
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Recipe")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for recipes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for recipes")
	}

	if len(mealRecipeAfterSelectHooks) != 0 {
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
		object.R.Recipe = foreign
		if foreign.R == nil {
			foreign.R = &recipeR{}
		}
		foreign.R.MealRecipes = append(foreign.R.MealRecipes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RecipeID == foreign.ID {
				local.R.Recipe = foreign
				if foreign.R == nil {
					foreign.R = &recipeR{}
				}
				foreign.R.MealRecipes = append(foreign.R.MealRecipes, local)
				break
			}
		}
	}

	return nil
}

// SetMeal of the mealRecipe to the related item.
// Sets o.R.Meal to related.
// Adds o to related.R.MealRecipes.
func (o *MealRecipe) SetMeal(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Meal) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"meal_recipe\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"meal_id"}),
		strmangle.WhereClause("\"", "\"", 2, mealRecipePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.MealID, o.RecipeID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.MealID = related.ID
	if o.R == nil {
		o.R = &mealRecipeR{
			Meal: related,
		}
	} else {
		o.R.Meal = related
	}

	if related.R == nil {
		related.R = &mealR{
			MealRecipes: MealRecipeSlice{o},
		}
	} else {
		related.R.MealRecipes = append(related.R.MealRecipes, o)
	}

	return nil
}

// SetRecipe of the mealRecipe to the related item.
// Sets o.R.Recipe to related.
// Adds o to related.R.MealRecipes.
func (o *MealRecipe) SetRecipe(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Recipe) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"meal_recipe\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"recipe_id"}),
		strmangle.WhereClause("\"", "\"", 2, mealRecipePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.MealID, o.RecipeID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.RecipeID = related.ID
	if o.R == nil {
		o.R = &mealRecipeR{
			Recipe: related,
		}
	} else {
		o.R.Recipe = related
	}

	if related.R == nil {
		related.R = &recipeR{
			MealRecipes: MealRecipeSlice{o},
		}
	} else {
		related.R.MealRecipes = append(related.R.MealRecipes, o)
	}

	return nil
}

// MealRecipes retrieves all the records using an executor.
func MealRecipes(mods ...qm.QueryMod) mealRecipeQuery {
	mods = append(mods, qm.From("\"meal_recipe\""))
	return mealRecipeQuery{NewQuery(mods...)}
}

// FindMealRecipe retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMealRecipe(ctx context.Context, exec boil.ContextExecutor, mealID string, recipeID string, selectCols ...string) (*MealRecipe, error) {
	mealRecipeObj := &MealRecipe{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"meal_recipe\" where \"meal_id\"=$1 AND \"recipe_id\"=$2", sel,
	)

	q := queries.Raw(query, mealID, recipeID)

	err := q.Bind(ctx, exec, mealRecipeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from meal_recipe")
	}

	if err = mealRecipeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return mealRecipeObj, err
	}

	return mealRecipeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MealRecipe) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no meal_recipe provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mealRecipeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	mealRecipeInsertCacheMut.RLock()
	cache, cached := mealRecipeInsertCache[key]
	mealRecipeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			mealRecipeAllColumns,
			mealRecipeColumnsWithDefault,
			mealRecipeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(mealRecipeType, mealRecipeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(mealRecipeType, mealRecipeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"meal_recipe\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"meal_recipe\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into meal_recipe")
	}

	if !cached {
		mealRecipeInsertCacheMut.Lock()
		mealRecipeInsertCache[key] = cache
		mealRecipeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the MealRecipe.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MealRecipe) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	mealRecipeUpdateCacheMut.RLock()
	cache, cached := mealRecipeUpdateCache[key]
	mealRecipeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			mealRecipeAllColumns,
			mealRecipePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update meal_recipe, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"meal_recipe\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, mealRecipePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(mealRecipeType, mealRecipeMapping, append(wl, mealRecipePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update meal_recipe row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for meal_recipe")
	}

	if !cached {
		mealRecipeUpdateCacheMut.Lock()
		mealRecipeUpdateCache[key] = cache
		mealRecipeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q mealRecipeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for meal_recipe")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for meal_recipe")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MealRecipeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealRecipePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"meal_recipe\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, mealRecipePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in mealRecipe slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all mealRecipe")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MealRecipe) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no meal_recipe provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mealRecipeColumnsWithDefault, o)

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

	mealRecipeUpsertCacheMut.RLock()
	cache, cached := mealRecipeUpsertCache[key]
	mealRecipeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			mealRecipeAllColumns,
			mealRecipeColumnsWithDefault,
			mealRecipeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			mealRecipeAllColumns,
			mealRecipePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert meal_recipe, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(mealRecipePrimaryKeyColumns))
			copy(conflict, mealRecipePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"meal_recipe\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(mealRecipeType, mealRecipeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(mealRecipeType, mealRecipeMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert meal_recipe")
	}

	if !cached {
		mealRecipeUpsertCacheMut.Lock()
		mealRecipeUpsertCache[key] = cache
		mealRecipeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single MealRecipe record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MealRecipe) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no MealRecipe provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), mealRecipePrimaryKeyMapping)
	sql := "DELETE FROM \"meal_recipe\" WHERE \"meal_id\"=$1 AND \"recipe_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from meal_recipe")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for meal_recipe")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q mealRecipeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no mealRecipeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from meal_recipe")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for meal_recipe")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MealRecipeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(mealRecipeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealRecipePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"meal_recipe\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mealRecipePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from mealRecipe slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for meal_recipe")
	}

	if len(mealRecipeAfterDeleteHooks) != 0 {
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
func (o *MealRecipe) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMealRecipe(ctx, exec, o.MealID, o.RecipeID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MealRecipeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MealRecipeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealRecipePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"meal_recipe\".* FROM \"meal_recipe\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mealRecipePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MealRecipeSlice")
	}

	*o = slice

	return nil
}

// MealRecipeExists checks if the MealRecipe row exists.
func MealRecipeExists(ctx context.Context, exec boil.ContextExecutor, mealID string, recipeID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"meal_recipe\" where \"meal_id\"=$1 AND \"recipe_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, mealID, recipeID)
	}
	row := exec.QueryRowContext(ctx, sql, mealID, recipeID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if meal_recipe exists")
	}

	return exists, nil
}