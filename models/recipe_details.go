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

// RecipeDetail is an object representing the database table.
type RecipeDetail struct {
	ID              string      `boil:"id" json:"id" toml:"id" yaml:"id"`
	RecipeID        string      `boil:"recipe_id" json:"recipe_id" toml:"recipe_id" yaml:"recipe_id"`
	Name            string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Equipment       null.String `boil:"equipment" json:"equipment,omitempty" toml:"equipment" yaml:"equipment,omitempty"`
	Source          null.JSON   `boil:"source" json:"source,omitempty" toml:"source" yaml:"source,omitempty"`
	Servings        null.Int    `boil:"servings" json:"servings,omitempty" toml:"servings" yaml:"servings,omitempty"`
	Quantity        null.Int    `boil:"quantity" json:"quantity,omitempty" toml:"quantity" yaml:"quantity,omitempty"`
	Unit            null.String `boil:"unit" json:"unit,omitempty" toml:"unit" yaml:"unit,omitempty"`
	Version         int         `boil:"version" json:"version" toml:"version" yaml:"version"`
	IsLatestVersion null.Bool   `boil:"is_latest_version" json:"is_latest_version,omitempty" toml:"is_latest_version" yaml:"is_latest_version,omitempty"`
	CreatedAt       time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *recipeDetailR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L recipeDetailL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RecipeDetailColumns = struct {
	ID              string
	RecipeID        string
	Name            string
	Equipment       string
	Source          string
	Servings        string
	Quantity        string
	Unit            string
	Version         string
	IsLatestVersion string
	CreatedAt       string
}{
	ID:              "id",
	RecipeID:        "recipe_id",
	Name:            "name",
	Equipment:       "equipment",
	Source:          "source",
	Servings:        "servings",
	Quantity:        "quantity",
	Unit:            "unit",
	Version:         "version",
	IsLatestVersion: "is_latest_version",
	CreatedAt:       "created_at",
}

var RecipeDetailTableColumns = struct {
	ID              string
	RecipeID        string
	Name            string
	Equipment       string
	Source          string
	Servings        string
	Quantity        string
	Unit            string
	Version         string
	IsLatestVersion string
	CreatedAt       string
}{
	ID:              "recipe_details.id",
	RecipeID:        "recipe_details.recipe_id",
	Name:            "recipe_details.name",
	Equipment:       "recipe_details.equipment",
	Source:          "recipe_details.source",
	Servings:        "recipe_details.servings",
	Quantity:        "recipe_details.quantity",
	Unit:            "recipe_details.unit",
	Version:         "recipe_details.version",
	IsLatestVersion: "recipe_details.is_latest_version",
	CreatedAt:       "recipe_details.created_at",
}

// Generated where

type whereHelpernull_Bool struct{ field string }

func (w whereHelpernull_Bool) EQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bool) NEQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bool) LT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bool) LTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bool) GT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bool) GTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Bool) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bool) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var RecipeDetailWhere = struct {
	ID              whereHelperstring
	RecipeID        whereHelperstring
	Name            whereHelperstring
	Equipment       whereHelpernull_String
	Source          whereHelpernull_JSON
	Servings        whereHelpernull_Int
	Quantity        whereHelpernull_Int
	Unit            whereHelpernull_String
	Version         whereHelperint
	IsLatestVersion whereHelpernull_Bool
	CreatedAt       whereHelpertime_Time
}{
	ID:              whereHelperstring{field: "\"recipe_details\".\"id\""},
	RecipeID:        whereHelperstring{field: "\"recipe_details\".\"recipe_id\""},
	Name:            whereHelperstring{field: "\"recipe_details\".\"name\""},
	Equipment:       whereHelpernull_String{field: "\"recipe_details\".\"equipment\""},
	Source:          whereHelpernull_JSON{field: "\"recipe_details\".\"source\""},
	Servings:        whereHelpernull_Int{field: "\"recipe_details\".\"servings\""},
	Quantity:        whereHelpernull_Int{field: "\"recipe_details\".\"quantity\""},
	Unit:            whereHelpernull_String{field: "\"recipe_details\".\"unit\""},
	Version:         whereHelperint{field: "\"recipe_details\".\"version\""},
	IsLatestVersion: whereHelpernull_Bool{field: "\"recipe_details\".\"is_latest_version\""},
	CreatedAt:       whereHelpertime_Time{field: "\"recipe_details\".\"created_at\""},
}

// RecipeDetailRels is where relationship names are stored.
var RecipeDetailRels = struct {
	Recipe         string
	RecipeSections string
}{
	Recipe:         "Recipe",
	RecipeSections: "RecipeSections",
}

// recipeDetailR is where relationships are stored.
type recipeDetailR struct {
	Recipe         *Recipe            `boil:"Recipe" json:"Recipe" toml:"Recipe" yaml:"Recipe"`
	RecipeSections RecipeSectionSlice `boil:"RecipeSections" json:"RecipeSections" toml:"RecipeSections" yaml:"RecipeSections"`
}

// NewStruct creates a new relationship struct
func (*recipeDetailR) NewStruct() *recipeDetailR {
	return &recipeDetailR{}
}

// recipeDetailL is where Load methods for each relationship are stored.
type recipeDetailL struct{}

var (
	recipeDetailAllColumns            = []string{"id", "recipe_id", "name", "equipment", "source", "servings", "quantity", "unit", "version", "is_latest_version", "created_at"}
	recipeDetailColumnsWithoutDefault = []string{"id", "recipe_id", "name", "equipment", "source", "servings", "quantity", "unit", "version"}
	recipeDetailColumnsWithDefault    = []string{"is_latest_version", "created_at"}
	recipeDetailPrimaryKeyColumns     = []string{"id"}
)

type (
	// RecipeDetailSlice is an alias for a slice of pointers to RecipeDetail.
	// This should almost always be used instead of []RecipeDetail.
	RecipeDetailSlice []*RecipeDetail
	// RecipeDetailHook is the signature for custom RecipeDetail hook methods
	RecipeDetailHook func(context.Context, boil.ContextExecutor, *RecipeDetail) error

	recipeDetailQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	recipeDetailType                 = reflect.TypeOf(&RecipeDetail{})
	recipeDetailMapping              = queries.MakeStructMapping(recipeDetailType)
	recipeDetailPrimaryKeyMapping, _ = queries.BindMapping(recipeDetailType, recipeDetailMapping, recipeDetailPrimaryKeyColumns)
	recipeDetailInsertCacheMut       sync.RWMutex
	recipeDetailInsertCache          = make(map[string]insertCache)
	recipeDetailUpdateCacheMut       sync.RWMutex
	recipeDetailUpdateCache          = make(map[string]updateCache)
	recipeDetailUpsertCacheMut       sync.RWMutex
	recipeDetailUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var recipeDetailBeforeInsertHooks []RecipeDetailHook
var recipeDetailBeforeUpdateHooks []RecipeDetailHook
var recipeDetailBeforeDeleteHooks []RecipeDetailHook
var recipeDetailBeforeUpsertHooks []RecipeDetailHook

var recipeDetailAfterInsertHooks []RecipeDetailHook
var recipeDetailAfterSelectHooks []RecipeDetailHook
var recipeDetailAfterUpdateHooks []RecipeDetailHook
var recipeDetailAfterDeleteHooks []RecipeDetailHook
var recipeDetailAfterUpsertHooks []RecipeDetailHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *RecipeDetail) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *RecipeDetail) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *RecipeDetail) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *RecipeDetail) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *RecipeDetail) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *RecipeDetail) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *RecipeDetail) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *RecipeDetail) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *RecipeDetail) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range recipeDetailAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRecipeDetailHook registers your hook function for all future operations.
func AddRecipeDetailHook(hookPoint boil.HookPoint, recipeDetailHook RecipeDetailHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		recipeDetailBeforeInsertHooks = append(recipeDetailBeforeInsertHooks, recipeDetailHook)
	case boil.BeforeUpdateHook:
		recipeDetailBeforeUpdateHooks = append(recipeDetailBeforeUpdateHooks, recipeDetailHook)
	case boil.BeforeDeleteHook:
		recipeDetailBeforeDeleteHooks = append(recipeDetailBeforeDeleteHooks, recipeDetailHook)
	case boil.BeforeUpsertHook:
		recipeDetailBeforeUpsertHooks = append(recipeDetailBeforeUpsertHooks, recipeDetailHook)
	case boil.AfterInsertHook:
		recipeDetailAfterInsertHooks = append(recipeDetailAfterInsertHooks, recipeDetailHook)
	case boil.AfterSelectHook:
		recipeDetailAfterSelectHooks = append(recipeDetailAfterSelectHooks, recipeDetailHook)
	case boil.AfterUpdateHook:
		recipeDetailAfterUpdateHooks = append(recipeDetailAfterUpdateHooks, recipeDetailHook)
	case boil.AfterDeleteHook:
		recipeDetailAfterDeleteHooks = append(recipeDetailAfterDeleteHooks, recipeDetailHook)
	case boil.AfterUpsertHook:
		recipeDetailAfterUpsertHooks = append(recipeDetailAfterUpsertHooks, recipeDetailHook)
	}
}

// One returns a single recipeDetail record from the query.
func (q recipeDetailQuery) One(ctx context.Context, exec boil.ContextExecutor) (*RecipeDetail, error) {
	o := &RecipeDetail{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for recipe_details")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all RecipeDetail records from the query.
func (q recipeDetailQuery) All(ctx context.Context, exec boil.ContextExecutor) (RecipeDetailSlice, error) {
	var o []*RecipeDetail

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RecipeDetail slice")
	}

	if len(recipeDetailAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all RecipeDetail records in the query.
func (q recipeDetailQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count recipe_details rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q recipeDetailQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if recipe_details exists")
	}

	return count > 0, nil
}

// Recipe pointed to by the foreign key.
func (o *RecipeDetail) Recipe(mods ...qm.QueryMod) recipeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.RecipeID),
	}

	queryMods = append(queryMods, mods...)

	query := Recipes(queryMods...)
	queries.SetFrom(query.Query, "\"recipes\"")

	return query
}

// RecipeSections retrieves all the recipe_section's RecipeSections with an executor.
func (o *RecipeDetail) RecipeSections(mods ...qm.QueryMod) recipeSectionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"recipe_sections\".\"recipe_detail\"=?", o.ID),
	)

	query := RecipeSections(queryMods...)
	queries.SetFrom(query.Query, "\"recipe_sections\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"recipe_sections\".*"})
	}

	return query
}

// LoadRecipe allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (recipeDetailL) LoadRecipe(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRecipeDetail interface{}, mods queries.Applicator) error {
	var slice []*RecipeDetail
	var object *RecipeDetail

	if singular {
		object = maybeRecipeDetail.(*RecipeDetail)
	} else {
		slice = *maybeRecipeDetail.(*[]*RecipeDetail)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &recipeDetailR{}
		}
		args = append(args, object.RecipeID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &recipeDetailR{}
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

	if len(recipeDetailAfterSelectHooks) != 0 {
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
		foreign.R.RecipeDetails = append(foreign.R.RecipeDetails, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RecipeID == foreign.ID {
				local.R.Recipe = foreign
				if foreign.R == nil {
					foreign.R = &recipeR{}
				}
				foreign.R.RecipeDetails = append(foreign.R.RecipeDetails, local)
				break
			}
		}
	}

	return nil
}

// LoadRecipeSections allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (recipeDetailL) LoadRecipeSections(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRecipeDetail interface{}, mods queries.Applicator) error {
	var slice []*RecipeDetail
	var object *RecipeDetail

	if singular {
		object = maybeRecipeDetail.(*RecipeDetail)
	} else {
		slice = *maybeRecipeDetail.(*[]*RecipeDetail)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &recipeDetailR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &recipeDetailR{}
			}

			for _, a := range args {
				if a == obj.ID {
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
		qm.From(`recipe_sections`),
		qm.WhereIn(`recipe_sections.recipe_detail in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load recipe_sections")
	}

	var resultSlice []*RecipeSection
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice recipe_sections")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on recipe_sections")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for recipe_sections")
	}

	if len(recipeSectionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.RecipeSections = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &recipeSectionR{}
			}
			foreign.R.RecipeSectionRecipeDetail = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.RecipeDetail {
				local.R.RecipeSections = append(local.R.RecipeSections, foreign)
				if foreign.R == nil {
					foreign.R = &recipeSectionR{}
				}
				foreign.R.RecipeSectionRecipeDetail = local
				break
			}
		}
	}

	return nil
}

// SetRecipe of the recipeDetail to the related item.
// Sets o.R.Recipe to related.
// Adds o to related.R.RecipeDetails.
func (o *RecipeDetail) SetRecipe(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Recipe) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"recipe_details\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"recipe_id"}),
		strmangle.WhereClause("\"", "\"", 2, recipeDetailPrimaryKeyColumns),
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

	o.RecipeID = related.ID
	if o.R == nil {
		o.R = &recipeDetailR{
			Recipe: related,
		}
	} else {
		o.R.Recipe = related
	}

	if related.R == nil {
		related.R = &recipeR{
			RecipeDetails: RecipeDetailSlice{o},
		}
	} else {
		related.R.RecipeDetails = append(related.R.RecipeDetails, o)
	}

	return nil
}

// AddRecipeSections adds the given related objects to the existing relationships
// of the recipe_detail, optionally inserting them as new records.
// Appends related to o.R.RecipeSections.
// Sets related.R.RecipeSectionRecipeDetail appropriately.
func (o *RecipeDetail) AddRecipeSections(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*RecipeSection) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.RecipeDetail = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"recipe_sections\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"recipe_detail"}),
				strmangle.WhereClause("\"", "\"", 2, recipeSectionPrimaryKeyColumns),
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

			rel.RecipeDetail = o.ID
		}
	}

	if o.R == nil {
		o.R = &recipeDetailR{
			RecipeSections: related,
		}
	} else {
		o.R.RecipeSections = append(o.R.RecipeSections, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &recipeSectionR{
				RecipeSectionRecipeDetail: o,
			}
		} else {
			rel.R.RecipeSectionRecipeDetail = o
		}
	}
	return nil
}

// RecipeDetails retrieves all the records using an executor.
func RecipeDetails(mods ...qm.QueryMod) recipeDetailQuery {
	mods = append(mods, qm.From("\"recipe_details\""))
	return recipeDetailQuery{NewQuery(mods...)}
}

// FindRecipeDetail retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRecipeDetail(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*RecipeDetail, error) {
	recipeDetailObj := &RecipeDetail{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"recipe_details\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, recipeDetailObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from recipe_details")
	}

	if err = recipeDetailObj.doAfterSelectHooks(ctx, exec); err != nil {
		return recipeDetailObj, err
	}

	return recipeDetailObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *RecipeDetail) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no recipe_details provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(recipeDetailColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	recipeDetailInsertCacheMut.RLock()
	cache, cached := recipeDetailInsertCache[key]
	recipeDetailInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			recipeDetailAllColumns,
			recipeDetailColumnsWithDefault,
			recipeDetailColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(recipeDetailType, recipeDetailMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(recipeDetailType, recipeDetailMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"recipe_details\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"recipe_details\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into recipe_details")
	}

	if !cached {
		recipeDetailInsertCacheMut.Lock()
		recipeDetailInsertCache[key] = cache
		recipeDetailInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the RecipeDetail.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *RecipeDetail) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	recipeDetailUpdateCacheMut.RLock()
	cache, cached := recipeDetailUpdateCache[key]
	recipeDetailUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			recipeDetailAllColumns,
			recipeDetailPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update recipe_details, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"recipe_details\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, recipeDetailPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(recipeDetailType, recipeDetailMapping, append(wl, recipeDetailPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update recipe_details row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for recipe_details")
	}

	if !cached {
		recipeDetailUpdateCacheMut.Lock()
		recipeDetailUpdateCache[key] = cache
		recipeDetailUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q recipeDetailQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for recipe_details")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for recipe_details")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RecipeDetailSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), recipeDetailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"recipe_details\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, recipeDetailPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in recipeDetail slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all recipeDetail")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *RecipeDetail) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no recipe_details provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(recipeDetailColumnsWithDefault, o)

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

	recipeDetailUpsertCacheMut.RLock()
	cache, cached := recipeDetailUpsertCache[key]
	recipeDetailUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			recipeDetailAllColumns,
			recipeDetailColumnsWithDefault,
			recipeDetailColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			recipeDetailAllColumns,
			recipeDetailPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert recipe_details, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(recipeDetailPrimaryKeyColumns))
			copy(conflict, recipeDetailPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"recipe_details\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(recipeDetailType, recipeDetailMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(recipeDetailType, recipeDetailMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert recipe_details")
	}

	if !cached {
		recipeDetailUpsertCacheMut.Lock()
		recipeDetailUpsertCache[key] = cache
		recipeDetailUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single RecipeDetail record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RecipeDetail) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no RecipeDetail provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), recipeDetailPrimaryKeyMapping)
	sql := "DELETE FROM \"recipe_details\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from recipe_details")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for recipe_details")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q recipeDetailQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no recipeDetailQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from recipe_details")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for recipe_details")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RecipeDetailSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(recipeDetailBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), recipeDetailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"recipe_details\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, recipeDetailPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from recipeDetail slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for recipe_details")
	}

	if len(recipeDetailAfterDeleteHooks) != 0 {
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
func (o *RecipeDetail) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRecipeDetail(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RecipeDetailSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RecipeDetailSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), recipeDetailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"recipe_details\".* FROM \"recipe_details\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, recipeDetailPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RecipeDetailSlice")
	}

	*o = slice

	return nil
}

// RecipeDetailExists checks if the RecipeDetail row exists.
func RecipeDetailExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"recipe_details\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if recipe_details exists")
	}

	return exists, nil
}
