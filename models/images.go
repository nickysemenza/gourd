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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Image is an object representing the database table.
type Image struct {
	ID       string `boil:"id" json:"id" toml:"id" yaml:"id"`
	BlurHash string `boil:"blur_hash" json:"blur_hash" toml:"blur_hash" yaml:"blur_hash"`
	Source   string `boil:"source" json:"source" toml:"source" yaml:"source"`

	R *imageR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L imageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ImageColumns = struct {
	ID       string
	BlurHash string
	Source   string
}{
	ID:       "id",
	BlurHash: "blur_hash",
	Source:   "source",
}

var ImageTableColumns = struct {
	ID       string
	BlurHash string
	Source   string
}{
	ID:       "images.id",
	BlurHash: "images.blur_hash",
	Source:   "images.source",
}

// Generated where

var ImageWhere = struct {
	ID       whereHelperstring
	BlurHash whereHelperstring
	Source   whereHelperstring
}{
	ID:       whereHelperstring{field: "\"images\".\"id\""},
	BlurHash: whereHelperstring{field: "\"images\".\"blur_hash\""},
	Source:   whereHelperstring{field: "\"images\".\"source\""},
}

// ImageRels is where relationship names are stored.
var ImageRels = struct {
	GphotosPhotos string
	NotionImages  string
}{
	GphotosPhotos: "GphotosPhotos",
	NotionImages:  "NotionImages",
}

// imageR is where relationships are stored.
type imageR struct {
	GphotosPhotos GphotosPhotoSlice `boil:"GphotosPhotos" json:"GphotosPhotos" toml:"GphotosPhotos" yaml:"GphotosPhotos"`
	NotionImages  NotionImageSlice  `boil:"NotionImages" json:"NotionImages" toml:"NotionImages" yaml:"NotionImages"`
}

// NewStruct creates a new relationship struct
func (*imageR) NewStruct() *imageR {
	return &imageR{}
}

// imageL is where Load methods for each relationship are stored.
type imageL struct{}

var (
	imageAllColumns            = []string{"id", "blur_hash", "source"}
	imageColumnsWithoutDefault = []string{"id", "blur_hash", "source"}
	imageColumnsWithDefault    = []string{}
	imagePrimaryKeyColumns     = []string{"id"}
	imageGeneratedColumns      = []string{}
)

type (
	// ImageSlice is an alias for a slice of pointers to Image.
	// This should almost always be used instead of []Image.
	ImageSlice []*Image
	// ImageHook is the signature for custom Image hook methods
	ImageHook func(context.Context, boil.ContextExecutor, *Image) error

	imageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	imageType                 = reflect.TypeOf(&Image{})
	imageMapping              = queries.MakeStructMapping(imageType)
	imagePrimaryKeyMapping, _ = queries.BindMapping(imageType, imageMapping, imagePrimaryKeyColumns)
	imageInsertCacheMut       sync.RWMutex
	imageInsertCache          = make(map[string]insertCache)
	imageUpdateCacheMut       sync.RWMutex
	imageUpdateCache          = make(map[string]updateCache)
	imageUpsertCacheMut       sync.RWMutex
	imageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var imageAfterSelectHooks []ImageHook

var imageBeforeInsertHooks []ImageHook
var imageAfterInsertHooks []ImageHook

var imageBeforeUpdateHooks []ImageHook
var imageAfterUpdateHooks []ImageHook

var imageBeforeDeleteHooks []ImageHook
var imageAfterDeleteHooks []ImageHook

var imageBeforeUpsertHooks []ImageHook
var imageAfterUpsertHooks []ImageHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Image) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Image) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Image) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Image) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Image) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Image) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Image) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Image) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Image) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddImageHook registers your hook function for all future operations.
func AddImageHook(hookPoint boil.HookPoint, imageHook ImageHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		imageAfterSelectHooks = append(imageAfterSelectHooks, imageHook)
	case boil.BeforeInsertHook:
		imageBeforeInsertHooks = append(imageBeforeInsertHooks, imageHook)
	case boil.AfterInsertHook:
		imageAfterInsertHooks = append(imageAfterInsertHooks, imageHook)
	case boil.BeforeUpdateHook:
		imageBeforeUpdateHooks = append(imageBeforeUpdateHooks, imageHook)
	case boil.AfterUpdateHook:
		imageAfterUpdateHooks = append(imageAfterUpdateHooks, imageHook)
	case boil.BeforeDeleteHook:
		imageBeforeDeleteHooks = append(imageBeforeDeleteHooks, imageHook)
	case boil.AfterDeleteHook:
		imageAfterDeleteHooks = append(imageAfterDeleteHooks, imageHook)
	case boil.BeforeUpsertHook:
		imageBeforeUpsertHooks = append(imageBeforeUpsertHooks, imageHook)
	case boil.AfterUpsertHook:
		imageAfterUpsertHooks = append(imageAfterUpsertHooks, imageHook)
	}
}

// One returns a single image record from the query.
func (q imageQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Image, error) {
	o := &Image{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for images")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Image records from the query.
func (q imageQuery) All(ctx context.Context, exec boil.ContextExecutor) (ImageSlice, error) {
	var o []*Image

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Image slice")
	}

	if len(imageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Image records in the query.
func (q imageQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count images rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q imageQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if images exists")
	}

	return count > 0, nil
}

// GphotosPhotos retrieves all the gphotos_photo's GphotosPhotos with an executor.
func (o *Image) GphotosPhotos(mods ...qm.QueryMod) gphotosPhotoQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"gphotos_photos\".\"image_id\"=?", o.ID),
	)

	return GphotosPhotos(queryMods...)
}

// NotionImages retrieves all the notion_image's NotionImages with an executor.
func (o *Image) NotionImages(mods ...qm.QueryMod) notionImageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"notion_image\".\"image_id\"=?", o.ID),
	)

	return NotionImages(queryMods...)
}

// LoadGphotosPhotos allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (imageL) LoadGphotosPhotos(ctx context.Context, e boil.ContextExecutor, singular bool, maybeImage interface{}, mods queries.Applicator) error {
	var slice []*Image
	var object *Image

	if singular {
		object = maybeImage.(*Image)
	} else {
		slice = *maybeImage.(*[]*Image)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &imageR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &imageR{}
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
		qm.From(`gphotos_photos`),
		qm.WhereIn(`gphotos_photos.image_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load gphotos_photos")
	}

	var resultSlice []*GphotosPhoto
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice gphotos_photos")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on gphotos_photos")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for gphotos_photos")
	}

	if len(gphotosPhotoAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.GphotosPhotos = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &gphotosPhotoR{}
			}
			foreign.R.Image = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ImageID {
				local.R.GphotosPhotos = append(local.R.GphotosPhotos, foreign)
				if foreign.R == nil {
					foreign.R = &gphotosPhotoR{}
				}
				foreign.R.Image = local
				break
			}
		}
	}

	return nil
}

// LoadNotionImages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (imageL) LoadNotionImages(ctx context.Context, e boil.ContextExecutor, singular bool, maybeImage interface{}, mods queries.Applicator) error {
	var slice []*Image
	var object *Image

	if singular {
		object = maybeImage.(*Image)
	} else {
		slice = *maybeImage.(*[]*Image)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &imageR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &imageR{}
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
		qm.From(`notion_image`),
		qm.WhereIn(`notion_image.image_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load notion_image")
	}

	var resultSlice []*NotionImage
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice notion_image")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on notion_image")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for notion_image")
	}

	if len(notionImageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.NotionImages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &notionImageR{}
			}
			foreign.R.Image = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ImageID {
				local.R.NotionImages = append(local.R.NotionImages, foreign)
				if foreign.R == nil {
					foreign.R = &notionImageR{}
				}
				foreign.R.Image = local
				break
			}
		}
	}

	return nil
}

// AddGphotosPhotos adds the given related objects to the existing relationships
// of the image, optionally inserting them as new records.
// Appends related to o.R.GphotosPhotos.
// Sets related.R.Image appropriately.
func (o *Image) AddGphotosPhotos(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*GphotosPhoto) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ImageID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"gphotos_photos\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"image_id"}),
				strmangle.WhereClause("\"", "\"", 2, gphotosPhotoPrimaryKeyColumns),
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

			rel.ImageID = o.ID
		}
	}

	if o.R == nil {
		o.R = &imageR{
			GphotosPhotos: related,
		}
	} else {
		o.R.GphotosPhotos = append(o.R.GphotosPhotos, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &gphotosPhotoR{
				Image: o,
			}
		} else {
			rel.R.Image = o
		}
	}
	return nil
}

// AddNotionImages adds the given related objects to the existing relationships
// of the image, optionally inserting them as new records.
// Appends related to o.R.NotionImages.
// Sets related.R.Image appropriately.
func (o *Image) AddNotionImages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*NotionImage) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ImageID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"notion_image\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"image_id"}),
				strmangle.WhereClause("\"", "\"", 2, notionImagePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.BlockID, rel.PageID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ImageID = o.ID
		}
	}

	if o.R == nil {
		o.R = &imageR{
			NotionImages: related,
		}
	} else {
		o.R.NotionImages = append(o.R.NotionImages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &notionImageR{
				Image: o,
			}
		} else {
			rel.R.Image = o
		}
	}
	return nil
}

// Images retrieves all the records using an executor.
func Images(mods ...qm.QueryMod) imageQuery {
	mods = append(mods, qm.From("\"images\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"images\".*"})
	}

	return imageQuery{q}
}

// FindImage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindImage(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Image, error) {
	imageObj := &Image{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"images\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, imageObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from images")
	}

	if err = imageObj.doAfterSelectHooks(ctx, exec); err != nil {
		return imageObj, err
	}

	return imageObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Image) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no images provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(imageColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	imageInsertCacheMut.RLock()
	cache, cached := imageInsertCache[key]
	imageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			imageAllColumns,
			imageColumnsWithDefault,
			imageColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(imageType, imageMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"images\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"images\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into images")
	}

	if !cached {
		imageInsertCacheMut.Lock()
		imageInsertCache[key] = cache
		imageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Image.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Image) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	imageUpdateCacheMut.RLock()
	cache, cached := imageUpdateCache[key]
	imageUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			imageAllColumns,
			imagePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update images, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"images\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, imagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, append(wl, imagePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update images row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for images")
	}

	if !cached {
		imageUpdateCacheMut.Lock()
		imageUpdateCache[key] = cache
		imageUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q imageQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for images")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ImageSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"images\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, imagePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in image slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all image")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Image) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no images provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(imageColumnsWithDefault, o)

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

	imageUpsertCacheMut.RLock()
	cache, cached := imageUpsertCache[key]
	imageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			imageAllColumns,
			imageColumnsWithDefault,
			imageColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			imageAllColumns,
			imagePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert images, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(imagePrimaryKeyColumns))
			copy(conflict, imagePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"images\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(imageType, imageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert images")
	}

	if !cached {
		imageUpsertCacheMut.Lock()
		imageUpsertCache[key] = cache
		imageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Image record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Image) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Image provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), imagePrimaryKeyMapping)
	sql := "DELETE FROM \"images\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for images")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q imageQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no imageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for images")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ImageSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(imageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"images\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, imagePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from image slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for images")
	}

	if len(imageAfterDeleteHooks) != 0 {
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
func (o *Image) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindImage(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ImageSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ImageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"images\".* FROM \"images\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, imagePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ImageSlice")
	}

	*o = slice

	return nil
}

// ImageExists checks if the Image row exists.
func ImageExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"images\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if images exists")
	}

	return exists, nil
}
