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

// GphotosAlbum is an object representing the database table.
type GphotosAlbum struct {
	ID      string `boil:"id" json:"id" toml:"id" yaml:"id"`
	Usecase string `boil:"usecase" json:"usecase" toml:"usecase" yaml:"usecase"`

	R *gphotosAlbumR `boil:"rel" json:"rel" toml:"rel" yaml:"rel"`
	L gphotosAlbumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GphotosAlbumColumns = struct {
	ID      string
	Usecase string
}{
	ID:      "id",
	Usecase: "usecase",
}

var GphotosAlbumTableColumns = struct {
	ID      string
	Usecase string
}{
	ID:      "gphotos_albums.id",
	Usecase: "gphotos_albums.usecase",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var GphotosAlbumWhere = struct {
	ID      whereHelperstring
	Usecase whereHelperstring
}{
	ID:      whereHelperstring{field: "\"gphotos_albums\".\"id\""},
	Usecase: whereHelperstring{field: "\"gphotos_albums\".\"usecase\""},
}

// GphotosAlbumRels is where relationship names are stored.
var GphotosAlbumRels = struct {
	AlbumGphotosPhotos string
}{
	AlbumGphotosPhotos: "AlbumGphotosPhotos",
}

// gphotosAlbumR is where relationships are stored.
type gphotosAlbumR struct {
	AlbumGphotosPhotos GphotosPhotoSlice `boil:"AlbumGphotosPhotos" json:"AlbumGphotosPhotos" toml:"AlbumGphotosPhotos" yaml:"AlbumGphotosPhotos"`
}

// NewStruct creates a new relationship struct
func (*gphotosAlbumR) NewStruct() *gphotosAlbumR {
	return &gphotosAlbumR{}
}

// gphotosAlbumL is where Load methods for each relationship are stored.
type gphotosAlbumL struct{}

var (
	gphotosAlbumAllColumns            = []string{"id", "usecase"}
	gphotosAlbumColumnsWithoutDefault = []string{"id", "usecase"}
	gphotosAlbumColumnsWithDefault    = []string{}
	gphotosAlbumPrimaryKeyColumns     = []string{"id"}
	gphotosAlbumGeneratedColumns      = []string{}
)

type (
	// GphotosAlbumSlice is an alias for a slice of pointers to GphotosAlbum.
	// This should almost always be used instead of []GphotosAlbum.
	GphotosAlbumSlice []*GphotosAlbum
	// GphotosAlbumHook is the signature for custom GphotosAlbum hook methods
	GphotosAlbumHook func(context.Context, boil.ContextExecutor, *GphotosAlbum) error

	gphotosAlbumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	gphotosAlbumType                 = reflect.TypeOf(&GphotosAlbum{})
	gphotosAlbumMapping              = queries.MakeStructMapping(gphotosAlbumType)
	gphotosAlbumPrimaryKeyMapping, _ = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, gphotosAlbumPrimaryKeyColumns)
	gphotosAlbumInsertCacheMut       sync.RWMutex
	gphotosAlbumInsertCache          = make(map[string]insertCache)
	gphotosAlbumUpdateCacheMut       sync.RWMutex
	gphotosAlbumUpdateCache          = make(map[string]updateCache)
	gphotosAlbumUpsertCacheMut       sync.RWMutex
	gphotosAlbumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var gphotosAlbumAfterSelectHooks []GphotosAlbumHook

var gphotosAlbumBeforeInsertHooks []GphotosAlbumHook
var gphotosAlbumAfterInsertHooks []GphotosAlbumHook

var gphotosAlbumBeforeUpdateHooks []GphotosAlbumHook
var gphotosAlbumAfterUpdateHooks []GphotosAlbumHook

var gphotosAlbumBeforeDeleteHooks []GphotosAlbumHook
var gphotosAlbumAfterDeleteHooks []GphotosAlbumHook

var gphotosAlbumBeforeUpsertHooks []GphotosAlbumHook
var gphotosAlbumAfterUpsertHooks []GphotosAlbumHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *GphotosAlbum) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *GphotosAlbum) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *GphotosAlbum) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *GphotosAlbum) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *GphotosAlbum) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *GphotosAlbum) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *GphotosAlbum) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *GphotosAlbum) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *GphotosAlbum) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range gphotosAlbumAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGphotosAlbumHook registers your hook function for all future operations.
func AddGphotosAlbumHook(hookPoint boil.HookPoint, gphotosAlbumHook GphotosAlbumHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		gphotosAlbumAfterSelectHooks = append(gphotosAlbumAfterSelectHooks, gphotosAlbumHook)
	case boil.BeforeInsertHook:
		gphotosAlbumBeforeInsertHooks = append(gphotosAlbumBeforeInsertHooks, gphotosAlbumHook)
	case boil.AfterInsertHook:
		gphotosAlbumAfterInsertHooks = append(gphotosAlbumAfterInsertHooks, gphotosAlbumHook)
	case boil.BeforeUpdateHook:
		gphotosAlbumBeforeUpdateHooks = append(gphotosAlbumBeforeUpdateHooks, gphotosAlbumHook)
	case boil.AfterUpdateHook:
		gphotosAlbumAfterUpdateHooks = append(gphotosAlbumAfterUpdateHooks, gphotosAlbumHook)
	case boil.BeforeDeleteHook:
		gphotosAlbumBeforeDeleteHooks = append(gphotosAlbumBeforeDeleteHooks, gphotosAlbumHook)
	case boil.AfterDeleteHook:
		gphotosAlbumAfterDeleteHooks = append(gphotosAlbumAfterDeleteHooks, gphotosAlbumHook)
	case boil.BeforeUpsertHook:
		gphotosAlbumBeforeUpsertHooks = append(gphotosAlbumBeforeUpsertHooks, gphotosAlbumHook)
	case boil.AfterUpsertHook:
		gphotosAlbumAfterUpsertHooks = append(gphotosAlbumAfterUpsertHooks, gphotosAlbumHook)
	}
}

// One returns a single gphotosAlbum record from the query.
func (q gphotosAlbumQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GphotosAlbum, error) {
	o := &GphotosAlbum{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for gphotos_albums")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all GphotosAlbum records from the query.
func (q gphotosAlbumQuery) All(ctx context.Context, exec boil.ContextExecutor) (GphotosAlbumSlice, error) {
	var o []*GphotosAlbum

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to GphotosAlbum slice")
	}

	if len(gphotosAlbumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all GphotosAlbum records in the query.
func (q gphotosAlbumQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count gphotos_albums rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q gphotosAlbumQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if gphotos_albums exists")
	}

	return count > 0, nil
}

// AlbumGphotosPhotos retrieves all the gphotos_photo's GphotosPhotos with an executor via album_id column.
func (o *GphotosAlbum) AlbumGphotosPhotos(mods ...qm.QueryMod) gphotosPhotoQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"gphotos_photos\".\"album_id\"=?", o.ID),
	)

	return GphotosPhotos(queryMods...)
}

// LoadAlbumGphotosPhotos allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (gphotosAlbumL) LoadAlbumGphotosPhotos(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGphotosAlbum interface{}, mods queries.Applicator) error {
	var slice []*GphotosAlbum
	var object *GphotosAlbum

	if singular {
		object = maybeGphotosAlbum.(*GphotosAlbum)
	} else {
		slice = *maybeGphotosAlbum.(*[]*GphotosAlbum)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &gphotosAlbumR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &gphotosAlbumR{}
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
		qm.WhereIn(`gphotos_photos.album_id in ?`, args...),
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
		object.R.AlbumGphotosPhotos = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &gphotosPhotoR{}
			}
			foreign.R.Album = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.AlbumID {
				local.R.AlbumGphotosPhotos = append(local.R.AlbumGphotosPhotos, foreign)
				if foreign.R == nil {
					foreign.R = &gphotosPhotoR{}
				}
				foreign.R.Album = local
				break
			}
		}
	}

	return nil
}

// AddAlbumGphotosPhotos adds the given related objects to the existing relationships
// of the gphotos_album, optionally inserting them as new records.
// Appends related to o.R.AlbumGphotosPhotos.
// Sets related.R.Album appropriately.
func (o *GphotosAlbum) AddAlbumGphotosPhotos(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*GphotosPhoto) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AlbumID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"gphotos_photos\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"album_id"}),
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

			rel.AlbumID = o.ID
		}
	}

	if o.R == nil {
		o.R = &gphotosAlbumR{
			AlbumGphotosPhotos: related,
		}
	} else {
		o.R.AlbumGphotosPhotos = append(o.R.AlbumGphotosPhotos, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &gphotosPhotoR{
				Album: o,
			}
		} else {
			rel.R.Album = o
		}
	}
	return nil
}

// GphotosAlbums retrieves all the records using an executor.
func GphotosAlbums(mods ...qm.QueryMod) gphotosAlbumQuery {
	mods = append(mods, qm.From("\"gphotos_albums\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"gphotos_albums\".*"})
	}

	return gphotosAlbumQuery{q}
}

// FindGphotosAlbum retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGphotosAlbum(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*GphotosAlbum, error) {
	gphotosAlbumObj := &GphotosAlbum{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"gphotos_albums\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, gphotosAlbumObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from gphotos_albums")
	}

	if err = gphotosAlbumObj.doAfterSelectHooks(ctx, exec); err != nil {
		return gphotosAlbumObj, err
	}

	return gphotosAlbumObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GphotosAlbum) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no gphotos_albums provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(gphotosAlbumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	gphotosAlbumInsertCacheMut.RLock()
	cache, cached := gphotosAlbumInsertCache[key]
	gphotosAlbumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			gphotosAlbumAllColumns,
			gphotosAlbumColumnsWithDefault,
			gphotosAlbumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"gphotos_albums\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"gphotos_albums\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into gphotos_albums")
	}

	if !cached {
		gphotosAlbumInsertCacheMut.Lock()
		gphotosAlbumInsertCache[key] = cache
		gphotosAlbumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the GphotosAlbum.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GphotosAlbum) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	gphotosAlbumUpdateCacheMut.RLock()
	cache, cached := gphotosAlbumUpdateCache[key]
	gphotosAlbumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			gphotosAlbumAllColumns,
			gphotosAlbumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update gphotos_albums, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"gphotos_albums\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, gphotosAlbumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, append(wl, gphotosAlbumPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update gphotos_albums row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for gphotos_albums")
	}

	if !cached {
		gphotosAlbumUpdateCacheMut.Lock()
		gphotosAlbumUpdateCache[key] = cache
		gphotosAlbumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q gphotosAlbumQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for gphotos_albums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for gphotos_albums")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GphotosAlbumSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gphotosAlbumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"gphotos_albums\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, gphotosAlbumPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in gphotosAlbum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all gphotosAlbum")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GphotosAlbum) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no gphotos_albums provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(gphotosAlbumColumnsWithDefault, o)

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

	gphotosAlbumUpsertCacheMut.RLock()
	cache, cached := gphotosAlbumUpsertCache[key]
	gphotosAlbumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			gphotosAlbumAllColumns,
			gphotosAlbumColumnsWithDefault,
			gphotosAlbumColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			gphotosAlbumAllColumns,
			gphotosAlbumPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert gphotos_albums, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(gphotosAlbumPrimaryKeyColumns))
			copy(conflict, gphotosAlbumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"gphotos_albums\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(gphotosAlbumType, gphotosAlbumMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert gphotos_albums")
	}

	if !cached {
		gphotosAlbumUpsertCacheMut.Lock()
		gphotosAlbumUpsertCache[key] = cache
		gphotosAlbumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single GphotosAlbum record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GphotosAlbum) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no GphotosAlbum provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), gphotosAlbumPrimaryKeyMapping)
	sql := "DELETE FROM \"gphotos_albums\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from gphotos_albums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for gphotos_albums")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q gphotosAlbumQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no gphotosAlbumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from gphotos_albums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for gphotos_albums")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GphotosAlbumSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(gphotosAlbumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gphotosAlbumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"gphotos_albums\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gphotosAlbumPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from gphotosAlbum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for gphotos_albums")
	}

	if len(gphotosAlbumAfterDeleteHooks) != 0 {
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
func (o *GphotosAlbum) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGphotosAlbum(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GphotosAlbumSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GphotosAlbumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), gphotosAlbumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"gphotos_albums\".* FROM \"gphotos_albums\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, gphotosAlbumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GphotosAlbumSlice")
	}

	*o = slice

	return nil
}

// GphotosAlbumExists checks if the GphotosAlbum row exists.
func GphotosAlbumExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"gphotos_albums\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if gphotos_albums exists")
	}

	return exists, nil
}
