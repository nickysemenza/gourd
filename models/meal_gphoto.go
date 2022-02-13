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

// MealGphoto is an object representing the database table.
type MealGphoto struct {
	MealID            string      `boil:"meal_id" json:"meal_id" toml:"meal_id" yaml:"meal_id"`
	GphotosID         string      `boil:"gphotos_id" json:"gphotos_id" toml:"gphotos_id" yaml:"gphotos_id"`
	HighlightRecipeID null.String `boil:"highlight_recipe_id" json:"highlight_recipe_id,omitempty" toml:"highlight_recipe_id" yaml:"highlight_recipe_id,omitempty"`

	R *mealGphotoR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L mealGphotoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MealGphotoColumns = struct {
	MealID            string
	GphotosID         string
	HighlightRecipeID string
}{
	MealID:            "meal_id",
	GphotosID:         "gphotos_id",
	HighlightRecipeID: "highlight_recipe_id",
}

var MealGphotoTableColumns = struct {
	MealID            string
	GphotosID         string
	HighlightRecipeID string
}{
	MealID:            "meal_gphoto.meal_id",
	GphotosID:         "meal_gphoto.gphotos_id",
	HighlightRecipeID: "meal_gphoto.highlight_recipe_id",
}

// Generated where

var MealGphotoWhere = struct {
	MealID            whereHelperstring
	GphotosID         whereHelperstring
	HighlightRecipeID whereHelpernull_String
}{
	MealID:            whereHelperstring{field: "\"meal_gphoto\".\"meal_id\""},
	GphotosID:         whereHelperstring{field: "\"meal_gphoto\".\"gphotos_id\""},
	HighlightRecipeID: whereHelpernull_String{field: "\"meal_gphoto\".\"highlight_recipe_id\""},
}

// MealGphotoRels is where relationship names are stored.
var MealGphotoRels = struct {
	Gphoto          string
	HighlightRecipe string
	Meal            string
}{
	Gphoto:          "Gphoto",
	HighlightRecipe: "HighlightRecipe",
	Meal:            "Meal",
}

// mealGphotoR is where relationships are stored.
type mealGphotoR struct {
	Gphoto          *GphotosPhoto `boil:"Gphoto" json:"Gphoto" toml:"Gphoto" yaml:"Gphoto"`
	HighlightRecipe *Recipe       `boil:"HighlightRecipe" json:"HighlightRecipe" toml:"HighlightRecipe" yaml:"HighlightRecipe"`
	Meal            *Meal         `boil:"Meal" json:"Meal" toml:"Meal" yaml:"Meal"`
}

// NewStruct creates a new relationship struct
func (*mealGphotoR) NewStruct() *mealGphotoR {
	return &mealGphotoR{}
}

// mealGphotoL is where Load methods for each relationship are stored.
type mealGphotoL struct{}

var (
	mealGphotoAllColumns            = []string{"meal_id", "gphotos_id", "highlight_recipe_id"}
	mealGphotoColumnsWithoutDefault = []string{"meal_id", "gphotos_id", "highlight_recipe_id"}
	mealGphotoColumnsWithDefault    = []string{}
	mealGphotoPrimaryKeyColumns     = []string{"meal_id", "gphotos_id"}
	mealGphotoGeneratedColumns      = []string{}
)

type (
	// MealGphotoSlice is an alias for a slice of pointers to MealGphoto.
	// This should almost always be used instead of []MealGphoto.
	MealGphotoSlice []*MealGphoto
	// MealGphotoHook is the signature for custom MealGphoto hook methods
	MealGphotoHook func(context.Context, boil.ContextExecutor, *MealGphoto) error

	mealGphotoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	mealGphotoType                 = reflect.TypeOf(&MealGphoto{})
	mealGphotoMapping              = queries.MakeStructMapping(mealGphotoType)
	mealGphotoPrimaryKeyMapping, _ = queries.BindMapping(mealGphotoType, mealGphotoMapping, mealGphotoPrimaryKeyColumns)
	mealGphotoInsertCacheMut       sync.RWMutex
	mealGphotoInsertCache          = make(map[string]insertCache)
	mealGphotoUpdateCacheMut       sync.RWMutex
	mealGphotoUpdateCache          = make(map[string]updateCache)
	mealGphotoUpsertCacheMut       sync.RWMutex
	mealGphotoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var mealGphotoAfterSelectHooks []MealGphotoHook

var mealGphotoBeforeInsertHooks []MealGphotoHook
var mealGphotoAfterInsertHooks []MealGphotoHook

var mealGphotoBeforeUpdateHooks []MealGphotoHook
var mealGphotoAfterUpdateHooks []MealGphotoHook

var mealGphotoBeforeDeleteHooks []MealGphotoHook
var mealGphotoAfterDeleteHooks []MealGphotoHook

var mealGphotoBeforeUpsertHooks []MealGphotoHook
var mealGphotoAfterUpsertHooks []MealGphotoHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *MealGphoto) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *MealGphoto) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *MealGphoto) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *MealGphoto) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *MealGphoto) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *MealGphoto) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *MealGphoto) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *MealGphoto) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *MealGphoto) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range mealGphotoAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMealGphotoHook registers your hook function for all future operations.
func AddMealGphotoHook(hookPoint boil.HookPoint, mealGphotoHook MealGphotoHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		mealGphotoAfterSelectHooks = append(mealGphotoAfterSelectHooks, mealGphotoHook)
	case boil.BeforeInsertHook:
		mealGphotoBeforeInsertHooks = append(mealGphotoBeforeInsertHooks, mealGphotoHook)
	case boil.AfterInsertHook:
		mealGphotoAfterInsertHooks = append(mealGphotoAfterInsertHooks, mealGphotoHook)
	case boil.BeforeUpdateHook:
		mealGphotoBeforeUpdateHooks = append(mealGphotoBeforeUpdateHooks, mealGphotoHook)
	case boil.AfterUpdateHook:
		mealGphotoAfterUpdateHooks = append(mealGphotoAfterUpdateHooks, mealGphotoHook)
	case boil.BeforeDeleteHook:
		mealGphotoBeforeDeleteHooks = append(mealGphotoBeforeDeleteHooks, mealGphotoHook)
	case boil.AfterDeleteHook:
		mealGphotoAfterDeleteHooks = append(mealGphotoAfterDeleteHooks, mealGphotoHook)
	case boil.BeforeUpsertHook:
		mealGphotoBeforeUpsertHooks = append(mealGphotoBeforeUpsertHooks, mealGphotoHook)
	case boil.AfterUpsertHook:
		mealGphotoAfterUpsertHooks = append(mealGphotoAfterUpsertHooks, mealGphotoHook)
	}
}

// One returns a single mealGphoto record from the query.
func (q mealGphotoQuery) One(ctx context.Context, exec boil.ContextExecutor) (*MealGphoto, error) {
	o := &MealGphoto{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for meal_gphoto")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all MealGphoto records from the query.
func (q mealGphotoQuery) All(ctx context.Context, exec boil.ContextExecutor) (MealGphotoSlice, error) {
	var o []*MealGphoto

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to MealGphoto slice")
	}

	if len(mealGphotoAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all MealGphoto records in the query.
func (q mealGphotoQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count meal_gphoto rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q mealGphotoQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if meal_gphoto exists")
	}

	return count > 0, nil
}

// Gphoto pointed to by the foreign key.
func (o *MealGphoto) Gphoto(mods ...qm.QueryMod) gphotosPhotoQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.GphotosID),
	}

	queryMods = append(queryMods, mods...)

	query := GphotosPhotos(queryMods...)
	queries.SetFrom(query.Query, "\"gphotos_photos\"")

	return query
}

// HighlightRecipe pointed to by the foreign key.
func (o *MealGphoto) HighlightRecipe(mods ...qm.QueryMod) recipeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.HighlightRecipeID),
	}

	queryMods = append(queryMods, mods...)

	query := Recipes(queryMods...)
	queries.SetFrom(query.Query, "\"recipes\"")

	return query
}

// Meal pointed to by the foreign key.
func (o *MealGphoto) Meal(mods ...qm.QueryMod) mealQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.MealID),
	}

	queryMods = append(queryMods, mods...)

	query := Meals(queryMods...)
	queries.SetFrom(query.Query, "\"meals\"")

	return query
}

// LoadGphoto allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (mealGphotoL) LoadGphoto(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMealGphoto interface{}, mods queries.Applicator) error {
	var slice []*MealGphoto
	var object *MealGphoto

	if singular {
		object = maybeMealGphoto.(*MealGphoto)
	} else {
		slice = *maybeMealGphoto.(*[]*MealGphoto)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &mealGphotoR{}
		}
		args = append(args, object.GphotosID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &mealGphotoR{}
			}

			for _, a := range args {
				if a == obj.GphotosID {
					continue Outer
				}
			}

			args = append(args, obj.GphotosID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`gphotos_photos`),
		qm.WhereIn(`gphotos_photos.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load GphotosPhoto")
	}

	var resultSlice []*GphotosPhoto
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice GphotosPhoto")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for gphotos_photos")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for gphotos_photos")
	}

	if len(mealGphotoAfterSelectHooks) != 0 {
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
		object.R.Gphoto = foreign
		if foreign.R == nil {
			foreign.R = &gphotosPhotoR{}
		}
		foreign.R.GphotoMealGphotos = append(foreign.R.GphotoMealGphotos, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.GphotosID == foreign.ID {
				local.R.Gphoto = foreign
				if foreign.R == nil {
					foreign.R = &gphotosPhotoR{}
				}
				foreign.R.GphotoMealGphotos = append(foreign.R.GphotoMealGphotos, local)
				break
			}
		}
	}

	return nil
}

// LoadHighlightRecipe allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (mealGphotoL) LoadHighlightRecipe(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMealGphoto interface{}, mods queries.Applicator) error {
	var slice []*MealGphoto
	var object *MealGphoto

	if singular {
		object = maybeMealGphoto.(*MealGphoto)
	} else {
		slice = *maybeMealGphoto.(*[]*MealGphoto)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &mealGphotoR{}
		}
		if !queries.IsNil(object.HighlightRecipeID) {
			args = append(args, object.HighlightRecipeID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &mealGphotoR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.HighlightRecipeID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.HighlightRecipeID) {
				args = append(args, obj.HighlightRecipeID)
			}

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

	if len(mealGphotoAfterSelectHooks) != 0 {
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
		object.R.HighlightRecipe = foreign
		if foreign.R == nil {
			foreign.R = &recipeR{}
		}
		foreign.R.HighlightRecipeMealGphotos = append(foreign.R.HighlightRecipeMealGphotos, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.HighlightRecipeID, foreign.ID) {
				local.R.HighlightRecipe = foreign
				if foreign.R == nil {
					foreign.R = &recipeR{}
				}
				foreign.R.HighlightRecipeMealGphotos = append(foreign.R.HighlightRecipeMealGphotos, local)
				break
			}
		}
	}

	return nil
}

// LoadMeal allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (mealGphotoL) LoadMeal(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMealGphoto interface{}, mods queries.Applicator) error {
	var slice []*MealGphoto
	var object *MealGphoto

	if singular {
		object = maybeMealGphoto.(*MealGphoto)
	} else {
		slice = *maybeMealGphoto.(*[]*MealGphoto)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &mealGphotoR{}
		}
		args = append(args, object.MealID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &mealGphotoR{}
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

	if len(mealGphotoAfterSelectHooks) != 0 {
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
		foreign.R.MealGphotos = append(foreign.R.MealGphotos, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.MealID == foreign.ID {
				local.R.Meal = foreign
				if foreign.R == nil {
					foreign.R = &mealR{}
				}
				foreign.R.MealGphotos = append(foreign.R.MealGphotos, local)
				break
			}
		}
	}

	return nil
}

// SetGphoto of the mealGphoto to the related item.
// Sets o.R.Gphoto to related.
// Adds o to related.R.GphotoMealGphotos.
func (o *MealGphoto) SetGphoto(ctx context.Context, exec boil.ContextExecutor, insert bool, related *GphotosPhoto) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"meal_gphoto\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"gphotos_id"}),
		strmangle.WhereClause("\"", "\"", 2, mealGphotoPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.MealID, o.GphotosID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GphotosID = related.ID
	if o.R == nil {
		o.R = &mealGphotoR{
			Gphoto: related,
		}
	} else {
		o.R.Gphoto = related
	}

	if related.R == nil {
		related.R = &gphotosPhotoR{
			GphotoMealGphotos: MealGphotoSlice{o},
		}
	} else {
		related.R.GphotoMealGphotos = append(related.R.GphotoMealGphotos, o)
	}

	return nil
}

// SetHighlightRecipe of the mealGphoto to the related item.
// Sets o.R.HighlightRecipe to related.
// Adds o to related.R.HighlightRecipeMealGphotos.
func (o *MealGphoto) SetHighlightRecipe(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Recipe) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"meal_gphoto\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"highlight_recipe_id"}),
		strmangle.WhereClause("\"", "\"", 2, mealGphotoPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.MealID, o.GphotosID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.HighlightRecipeID, related.ID)
	if o.R == nil {
		o.R = &mealGphotoR{
			HighlightRecipe: related,
		}
	} else {
		o.R.HighlightRecipe = related
	}

	if related.R == nil {
		related.R = &recipeR{
			HighlightRecipeMealGphotos: MealGphotoSlice{o},
		}
	} else {
		related.R.HighlightRecipeMealGphotos = append(related.R.HighlightRecipeMealGphotos, o)
	}

	return nil
}

// RemoveHighlightRecipe relationship.
// Sets o.R.HighlightRecipe to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *MealGphoto) RemoveHighlightRecipe(ctx context.Context, exec boil.ContextExecutor, related *Recipe) error {
	var err error

	queries.SetScanner(&o.HighlightRecipeID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("highlight_recipe_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.HighlightRecipe = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.HighlightRecipeMealGphotos {
		if queries.Equal(o.HighlightRecipeID, ri.HighlightRecipeID) {
			continue
		}

		ln := len(related.R.HighlightRecipeMealGphotos)
		if ln > 1 && i < ln-1 {
			related.R.HighlightRecipeMealGphotos[i] = related.R.HighlightRecipeMealGphotos[ln-1]
		}
		related.R.HighlightRecipeMealGphotos = related.R.HighlightRecipeMealGphotos[:ln-1]
		break
	}
	return nil
}

// SetMeal of the mealGphoto to the related item.
// Sets o.R.Meal to related.
// Adds o to related.R.MealGphotos.
func (o *MealGphoto) SetMeal(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Meal) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"meal_gphoto\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"meal_id"}),
		strmangle.WhereClause("\"", "\"", 2, mealGphotoPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.MealID, o.GphotosID}

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
		o.R = &mealGphotoR{
			Meal: related,
		}
	} else {
		o.R.Meal = related
	}

	if related.R == nil {
		related.R = &mealR{
			MealGphotos: MealGphotoSlice{o},
		}
	} else {
		related.R.MealGphotos = append(related.R.MealGphotos, o)
	}

	return nil
}

// MealGphotos retrieves all the records using an executor.
func MealGphotos(mods ...qm.QueryMod) mealGphotoQuery {
	mods = append(mods, qm.From("\"meal_gphoto\""))
	return mealGphotoQuery{NewQuery(mods...)}
}

// FindMealGphoto retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMealGphoto(ctx context.Context, exec boil.ContextExecutor, mealID string, gphotosID string, selectCols ...string) (*MealGphoto, error) {
	mealGphotoObj := &MealGphoto{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"meal_gphoto\" where \"meal_id\"=$1 AND \"gphotos_id\"=$2", sel,
	)

	q := queries.Raw(query, mealID, gphotosID)

	err := q.Bind(ctx, exec, mealGphotoObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from meal_gphoto")
	}

	if err = mealGphotoObj.doAfterSelectHooks(ctx, exec); err != nil {
		return mealGphotoObj, err
	}

	return mealGphotoObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MealGphoto) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no meal_gphoto provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mealGphotoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	mealGphotoInsertCacheMut.RLock()
	cache, cached := mealGphotoInsertCache[key]
	mealGphotoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			mealGphotoAllColumns,
			mealGphotoColumnsWithDefault,
			mealGphotoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(mealGphotoType, mealGphotoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(mealGphotoType, mealGphotoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"meal_gphoto\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"meal_gphoto\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into meal_gphoto")
	}

	if !cached {
		mealGphotoInsertCacheMut.Lock()
		mealGphotoInsertCache[key] = cache
		mealGphotoInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the MealGphoto.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MealGphoto) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	mealGphotoUpdateCacheMut.RLock()
	cache, cached := mealGphotoUpdateCache[key]
	mealGphotoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			mealGphotoAllColumns,
			mealGphotoPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update meal_gphoto, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"meal_gphoto\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, mealGphotoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(mealGphotoType, mealGphotoMapping, append(wl, mealGphotoPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update meal_gphoto row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for meal_gphoto")
	}

	if !cached {
		mealGphotoUpdateCacheMut.Lock()
		mealGphotoUpdateCache[key] = cache
		mealGphotoUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q mealGphotoQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for meal_gphoto")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for meal_gphoto")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MealGphotoSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealGphotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"meal_gphoto\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, mealGphotoPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in mealGphoto slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all mealGphoto")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MealGphoto) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no meal_gphoto provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mealGphotoColumnsWithDefault, o)

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

	mealGphotoUpsertCacheMut.RLock()
	cache, cached := mealGphotoUpsertCache[key]
	mealGphotoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			mealGphotoAllColumns,
			mealGphotoColumnsWithDefault,
			mealGphotoColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			mealGphotoAllColumns,
			mealGphotoPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert meal_gphoto, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(mealGphotoPrimaryKeyColumns))
			copy(conflict, mealGphotoPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"meal_gphoto\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(mealGphotoType, mealGphotoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(mealGphotoType, mealGphotoMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert meal_gphoto")
	}

	if !cached {
		mealGphotoUpsertCacheMut.Lock()
		mealGphotoUpsertCache[key] = cache
		mealGphotoUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single MealGphoto record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MealGphoto) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no MealGphoto provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), mealGphotoPrimaryKeyMapping)
	sql := "DELETE FROM \"meal_gphoto\" WHERE \"meal_id\"=$1 AND \"gphotos_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from meal_gphoto")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for meal_gphoto")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q mealGphotoQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no mealGphotoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from meal_gphoto")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for meal_gphoto")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MealGphotoSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(mealGphotoBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealGphotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"meal_gphoto\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mealGphotoPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from mealGphoto slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for meal_gphoto")
	}

	if len(mealGphotoAfterDeleteHooks) != 0 {
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
func (o *MealGphoto) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMealGphoto(ctx, exec, o.MealID, o.GphotosID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MealGphotoSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MealGphotoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mealGphotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"meal_gphoto\".* FROM \"meal_gphoto\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mealGphotoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MealGphotoSlice")
	}

	*o = slice

	return nil
}

// MealGphotoExists checks if the MealGphoto row exists.
func MealGphotoExists(ctx context.Context, exec boil.ContextExecutor, mealID string, gphotosID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"meal_gphoto\" where \"meal_id\"=$1 AND \"gphotos_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, mealID, gphotosID)
	}
	row := exec.QueryRowContext(ctx, sql, mealID, gphotosID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if meal_gphoto exists")
	}

	return exists, nil
}
