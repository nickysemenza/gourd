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
	"github.com/volatiletech/strmangle"
)

// NotionImage is an object representing the database table.
type NotionImage struct {
	BlockID  string    `boil:"block_id" json:"block_id" toml:"block_id" yaml:"block_id"`
	PageID   string    `boil:"page_id" json:"page_id" toml:"page_id" yaml:"page_id"`
	LastSeen time.Time `boil:"last_seen" json:"last_seen" toml:"last_seen" yaml:"last_seen"`
	Image    string    `boil:"image" json:"image" toml:"image" yaml:"image"`

	R *notionImageR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L notionImageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var NotionImageColumns = struct {
	BlockID  string
	PageID   string
	LastSeen string
	Image    string
}{
	BlockID:  "block_id",
	PageID:   "page_id",
	LastSeen: "last_seen",
	Image:    "image",
}

var NotionImageTableColumns = struct {
	BlockID  string
	PageID   string
	LastSeen string
	Image    string
}{
	BlockID:  "notion_image.block_id",
	PageID:   "notion_image.page_id",
	LastSeen: "notion_image.last_seen",
	Image:    "notion_image.image",
}

// Generated where

var NotionImageWhere = struct {
	BlockID  whereHelperstring
	PageID   whereHelperstring
	LastSeen whereHelpertime_Time
	Image    whereHelperstring
}{
	BlockID:  whereHelperstring{field: "\"notion_image\".\"block_id\""},
	PageID:   whereHelperstring{field: "\"notion_image\".\"page_id\""},
	LastSeen: whereHelpertime_Time{field: "\"notion_image\".\"last_seen\""},
	Image:    whereHelperstring{field: "\"notion_image\".\"image\""},
}

// NotionImageRels is where relationship names are stored.
var NotionImageRels = struct {
	NotionImageImage string
	Page             string
}{
	NotionImageImage: "NotionImageImage",
	Page:             "Page",
}

// notionImageR is where relationships are stored.
type notionImageR struct {
	NotionImageImage *Image        `boil:"NotionImageImage" json:"NotionImageImage" toml:"NotionImageImage" yaml:"NotionImageImage"`
	Page             *NotionRecipe `boil:"Page" json:"Page" toml:"Page" yaml:"Page"`
}

// NewStruct creates a new relationship struct
func (*notionImageR) NewStruct() *notionImageR {
	return &notionImageR{}
}

// notionImageL is where Load methods for each relationship are stored.
type notionImageL struct{}

var (
	notionImageAllColumns            = []string{"block_id", "page_id", "last_seen", "image"}
	notionImageColumnsWithoutDefault = []string{"block_id", "page_id", "image"}
	notionImageColumnsWithDefault    = []string{"last_seen"}
	notionImagePrimaryKeyColumns     = []string{"block_id", "page_id"}
)

type (
	// NotionImageSlice is an alias for a slice of pointers to NotionImage.
	// This should almost always be used instead of []NotionImage.
	NotionImageSlice []*NotionImage
	// NotionImageHook is the signature for custom NotionImage hook methods
	NotionImageHook func(context.Context, boil.ContextExecutor, *NotionImage) error

	notionImageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	notionImageType                 = reflect.TypeOf(&NotionImage{})
	notionImageMapping              = queries.MakeStructMapping(notionImageType)
	notionImagePrimaryKeyMapping, _ = queries.BindMapping(notionImageType, notionImageMapping, notionImagePrimaryKeyColumns)
	notionImageInsertCacheMut       sync.RWMutex
	notionImageInsertCache          = make(map[string]insertCache)
	notionImageUpdateCacheMut       sync.RWMutex
	notionImageUpdateCache          = make(map[string]updateCache)
	notionImageUpsertCacheMut       sync.RWMutex
	notionImageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var notionImageBeforeInsertHooks []NotionImageHook
var notionImageBeforeUpdateHooks []NotionImageHook
var notionImageBeforeDeleteHooks []NotionImageHook
var notionImageBeforeUpsertHooks []NotionImageHook

var notionImageAfterInsertHooks []NotionImageHook
var notionImageAfterSelectHooks []NotionImageHook
var notionImageAfterUpdateHooks []NotionImageHook
var notionImageAfterDeleteHooks []NotionImageHook
var notionImageAfterUpsertHooks []NotionImageHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *NotionImage) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *NotionImage) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *NotionImage) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *NotionImage) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *NotionImage) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *NotionImage) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *NotionImage) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *NotionImage) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *NotionImage) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range notionImageAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddNotionImageHook registers your hook function for all future operations.
func AddNotionImageHook(hookPoint boil.HookPoint, notionImageHook NotionImageHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		notionImageBeforeInsertHooks = append(notionImageBeforeInsertHooks, notionImageHook)
	case boil.BeforeUpdateHook:
		notionImageBeforeUpdateHooks = append(notionImageBeforeUpdateHooks, notionImageHook)
	case boil.BeforeDeleteHook:
		notionImageBeforeDeleteHooks = append(notionImageBeforeDeleteHooks, notionImageHook)
	case boil.BeforeUpsertHook:
		notionImageBeforeUpsertHooks = append(notionImageBeforeUpsertHooks, notionImageHook)
	case boil.AfterInsertHook:
		notionImageAfterInsertHooks = append(notionImageAfterInsertHooks, notionImageHook)
	case boil.AfterSelectHook:
		notionImageAfterSelectHooks = append(notionImageAfterSelectHooks, notionImageHook)
	case boil.AfterUpdateHook:
		notionImageAfterUpdateHooks = append(notionImageAfterUpdateHooks, notionImageHook)
	case boil.AfterDeleteHook:
		notionImageAfterDeleteHooks = append(notionImageAfterDeleteHooks, notionImageHook)
	case boil.AfterUpsertHook:
		notionImageAfterUpsertHooks = append(notionImageAfterUpsertHooks, notionImageHook)
	}
}

// One returns a single notionImage record from the query.
func (q notionImageQuery) One(ctx context.Context, exec boil.ContextExecutor) (*NotionImage, error) {
	o := &NotionImage{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for notion_image")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all NotionImage records from the query.
func (q notionImageQuery) All(ctx context.Context, exec boil.ContextExecutor) (NotionImageSlice, error) {
	var o []*NotionImage

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to NotionImage slice")
	}

	if len(notionImageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all NotionImage records in the query.
func (q notionImageQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count notion_image rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q notionImageQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if notion_image exists")
	}

	return count > 0, nil
}

// NotionImageImage pointed to by the foreign key.
func (o *NotionImage) NotionImageImage(mods ...qm.QueryMod) imageQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.Image),
	}

	queryMods = append(queryMods, mods...)

	query := Images(queryMods...)
	queries.SetFrom(query.Query, "\"images\"")

	return query
}

// Page pointed to by the foreign key.
func (o *NotionImage) Page(mods ...qm.QueryMod) notionRecipeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"page_id\" = ?", o.PageID),
	}

	queryMods = append(queryMods, mods...)

	query := NotionRecipes(queryMods...)
	queries.SetFrom(query.Query, "\"notion_recipe\"")

	return query
}

// LoadNotionImageImage allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (notionImageL) LoadNotionImageImage(ctx context.Context, e boil.ContextExecutor, singular bool, maybeNotionImage interface{}, mods queries.Applicator) error {
	var slice []*NotionImage
	var object *NotionImage

	if singular {
		object = maybeNotionImage.(*NotionImage)
	} else {
		slice = *maybeNotionImage.(*[]*NotionImage)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &notionImageR{}
		}
		args = append(args, object.Image)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &notionImageR{}
			}

			for _, a := range args {
				if a == obj.Image {
					continue Outer
				}
			}

			args = append(args, obj.Image)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`images`),
		qm.WhereIn(`images.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Image")
	}

	var resultSlice []*Image
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Image")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for images")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for images")
	}

	if len(notionImageAfterSelectHooks) != 0 {
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
		object.R.NotionImageImage = foreign
		if foreign.R == nil {
			foreign.R = &imageR{}
		}
		foreign.R.NotionImages = append(foreign.R.NotionImages, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.Image == foreign.ID {
				local.R.NotionImageImage = foreign
				if foreign.R == nil {
					foreign.R = &imageR{}
				}
				foreign.R.NotionImages = append(foreign.R.NotionImages, local)
				break
			}
		}
	}

	return nil
}

// LoadPage allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (notionImageL) LoadPage(ctx context.Context, e boil.ContextExecutor, singular bool, maybeNotionImage interface{}, mods queries.Applicator) error {
	var slice []*NotionImage
	var object *NotionImage

	if singular {
		object = maybeNotionImage.(*NotionImage)
	} else {
		slice = *maybeNotionImage.(*[]*NotionImage)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &notionImageR{}
		}
		args = append(args, object.PageID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &notionImageR{}
			}

			for _, a := range args {
				if a == obj.PageID {
					continue Outer
				}
			}

			args = append(args, obj.PageID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`notion_recipe`),
		qm.WhereIn(`notion_recipe.page_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load NotionRecipe")
	}

	var resultSlice []*NotionRecipe
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice NotionRecipe")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for notion_recipe")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for notion_recipe")
	}

	if len(notionImageAfterSelectHooks) != 0 {
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
		object.R.Page = foreign
		if foreign.R == nil {
			foreign.R = &notionRecipeR{}
		}
		foreign.R.PageNotionImages = append(foreign.R.PageNotionImages, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PageID == foreign.PageID {
				local.R.Page = foreign
				if foreign.R == nil {
					foreign.R = &notionRecipeR{}
				}
				foreign.R.PageNotionImages = append(foreign.R.PageNotionImages, local)
				break
			}
		}
	}

	return nil
}

// SetNotionImageImage of the notionImage to the related item.
// Sets o.R.NotionImageImage to related.
// Adds o to related.R.NotionImages.
func (o *NotionImage) SetNotionImageImage(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Image) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"notion_image\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"image"}),
		strmangle.WhereClause("\"", "\"", 2, notionImagePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.BlockID, o.PageID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.Image = related.ID
	if o.R == nil {
		o.R = &notionImageR{
			NotionImageImage: related,
		}
	} else {
		o.R.NotionImageImage = related
	}

	if related.R == nil {
		related.R = &imageR{
			NotionImages: NotionImageSlice{o},
		}
	} else {
		related.R.NotionImages = append(related.R.NotionImages, o)
	}

	return nil
}

// SetPage of the notionImage to the related item.
// Sets o.R.Page to related.
// Adds o to related.R.PageNotionImages.
func (o *NotionImage) SetPage(ctx context.Context, exec boil.ContextExecutor, insert bool, related *NotionRecipe) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"notion_image\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"page_id"}),
		strmangle.WhereClause("\"", "\"", 2, notionImagePrimaryKeyColumns),
	)
	values := []interface{}{related.PageID, o.BlockID, o.PageID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PageID = related.PageID
	if o.R == nil {
		o.R = &notionImageR{
			Page: related,
		}
	} else {
		o.R.Page = related
	}

	if related.R == nil {
		related.R = &notionRecipeR{
			PageNotionImages: NotionImageSlice{o},
		}
	} else {
		related.R.PageNotionImages = append(related.R.PageNotionImages, o)
	}

	return nil
}

// NotionImages retrieves all the records using an executor.
func NotionImages(mods ...qm.QueryMod) notionImageQuery {
	mods = append(mods, qm.From("\"notion_image\""))
	return notionImageQuery{NewQuery(mods...)}
}

// FindNotionImage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindNotionImage(ctx context.Context, exec boil.ContextExecutor, blockID string, pageID string, selectCols ...string) (*NotionImage, error) {
	notionImageObj := &NotionImage{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"notion_image\" where \"block_id\"=$1 AND \"page_id\"=$2", sel,
	)

	q := queries.Raw(query, blockID, pageID)

	err := q.Bind(ctx, exec, notionImageObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from notion_image")
	}

	if err = notionImageObj.doAfterSelectHooks(ctx, exec); err != nil {
		return notionImageObj, err
	}

	return notionImageObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *NotionImage) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no notion_image provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(notionImageColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	notionImageInsertCacheMut.RLock()
	cache, cached := notionImageInsertCache[key]
	notionImageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			notionImageAllColumns,
			notionImageColumnsWithDefault,
			notionImageColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(notionImageType, notionImageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(notionImageType, notionImageMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"notion_image\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"notion_image\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into notion_image")
	}

	if !cached {
		notionImageInsertCacheMut.Lock()
		notionImageInsertCache[key] = cache
		notionImageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the NotionImage.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *NotionImage) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	notionImageUpdateCacheMut.RLock()
	cache, cached := notionImageUpdateCache[key]
	notionImageUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			notionImageAllColumns,
			notionImagePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update notion_image, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"notion_image\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, notionImagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(notionImageType, notionImageMapping, append(wl, notionImagePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update notion_image row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for notion_image")
	}

	if !cached {
		notionImageUpdateCacheMut.Lock()
		notionImageUpdateCache[key] = cache
		notionImageUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q notionImageQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for notion_image")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for notion_image")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o NotionImageSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notionImagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"notion_image\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, notionImagePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in notionImage slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all notionImage")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *NotionImage) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no notion_image provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(notionImageColumnsWithDefault, o)

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

	notionImageUpsertCacheMut.RLock()
	cache, cached := notionImageUpsertCache[key]
	notionImageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			notionImageAllColumns,
			notionImageColumnsWithDefault,
			notionImageColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			notionImageAllColumns,
			notionImagePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert notion_image, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(notionImagePrimaryKeyColumns))
			copy(conflict, notionImagePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"notion_image\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(notionImageType, notionImageMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(notionImageType, notionImageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert notion_image")
	}

	if !cached {
		notionImageUpsertCacheMut.Lock()
		notionImageUpsertCache[key] = cache
		notionImageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single NotionImage record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *NotionImage) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no NotionImage provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), notionImagePrimaryKeyMapping)
	sql := "DELETE FROM \"notion_image\" WHERE \"block_id\"=$1 AND \"page_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from notion_image")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for notion_image")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q notionImageQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no notionImageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from notion_image")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for notion_image")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o NotionImageSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(notionImageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notionImagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"notion_image\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, notionImagePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from notionImage slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for notion_image")
	}

	if len(notionImageAfterDeleteHooks) != 0 {
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
func (o *NotionImage) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindNotionImage(ctx, exec, o.BlockID, o.PageID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *NotionImageSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := NotionImageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notionImagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"notion_image\".* FROM \"notion_image\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, notionImagePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in NotionImageSlice")
	}

	*o = slice

	return nil
}

// NotionImageExists checks if the NotionImage row exists.
func NotionImageExists(ctx context.Context, exec boil.ContextExecutor, blockID string, pageID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"notion_image\" where \"block_id\"=$1 AND \"page_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, blockID, pageID)
	}
	row := exec.QueryRowContext(ctx, sql, blockID, pageID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if notion_image exists")
	}

	return exists, nil
}
