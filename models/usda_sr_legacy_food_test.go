// Code generated by SQLBoiler 4.8.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testUsdaSRLegacyFoods(t *testing.T) {
	t.Parallel()

	query := UsdaSRLegacyFoods()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUsdaSRLegacyFoodsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaSRLegacyFoodsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := UsdaSRLegacyFoods().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaSRLegacyFoodsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UsdaSRLegacyFoodSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaSRLegacyFoodsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UsdaSRLegacyFoodExists(ctx, tx, o.FDCID)
	if err != nil {
		t.Errorf("Unable to check if UsdaSRLegacyFood exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UsdaSRLegacyFoodExists to return true, but got false.")
	}
}

func testUsdaSRLegacyFoodsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	usdaSRLegacyFoodFound, err := FindUsdaSRLegacyFood(ctx, tx, o.FDCID)
	if err != nil {
		t.Error(err)
	}

	if usdaSRLegacyFoodFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUsdaSRLegacyFoodsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = UsdaSRLegacyFoods().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUsdaSRLegacyFoodsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := UsdaSRLegacyFoods().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUsdaSRLegacyFoodsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	usdaSRLegacyFoodOne := &UsdaSRLegacyFood{}
	usdaSRLegacyFoodTwo := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, usdaSRLegacyFoodOne, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}
	if err = randomize.Struct(seed, usdaSRLegacyFoodTwo, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = usdaSRLegacyFoodOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = usdaSRLegacyFoodTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UsdaSRLegacyFoods().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUsdaSRLegacyFoodsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	usdaSRLegacyFoodOne := &UsdaSRLegacyFood{}
	usdaSRLegacyFoodTwo := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, usdaSRLegacyFoodOne, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}
	if err = randomize.Struct(seed, usdaSRLegacyFoodTwo, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = usdaSRLegacyFoodOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = usdaSRLegacyFoodTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func usdaSRLegacyFoodBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func usdaSRLegacyFoodAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaSRLegacyFood) error {
	*o = UsdaSRLegacyFood{}
	return nil
}

func testUsdaSRLegacyFoodsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &UsdaSRLegacyFood{}
	o := &UsdaSRLegacyFood{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood object: %s", err)
	}

	AddUsdaSRLegacyFoodHook(boil.BeforeInsertHook, usdaSRLegacyFoodBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodBeforeInsertHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.AfterInsertHook, usdaSRLegacyFoodAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodAfterInsertHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.AfterSelectHook, usdaSRLegacyFoodAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodAfterSelectHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.BeforeUpdateHook, usdaSRLegacyFoodBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodBeforeUpdateHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.AfterUpdateHook, usdaSRLegacyFoodAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodAfterUpdateHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.BeforeDeleteHook, usdaSRLegacyFoodBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodBeforeDeleteHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.AfterDeleteHook, usdaSRLegacyFoodAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodAfterDeleteHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.BeforeUpsertHook, usdaSRLegacyFoodBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodBeforeUpsertHooks = []UsdaSRLegacyFoodHook{}

	AddUsdaSRLegacyFoodHook(boil.AfterUpsertHook, usdaSRLegacyFoodAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	usdaSRLegacyFoodAfterUpsertHooks = []UsdaSRLegacyFoodHook{}
}

func testUsdaSRLegacyFoodsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsdaSRLegacyFoodsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(usdaSRLegacyFoodColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsdaSRLegacyFoodToOneUsdaFoodUsingFDC(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UsdaSRLegacyFood
	var foreign UsdaFood

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, usdaFoodDBTypes, false, usdaFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFood struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.FDCID = foreign.FDCID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.FDC().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.FDCID != foreign.FDCID {
		t.Errorf("want: %v, got %v", foreign.FDCID, check.FDCID)
	}

	slice := UsdaSRLegacyFoodSlice{&local}
	if err = local.L.LoadFDC(ctx, tx, false, (*[]*UsdaSRLegacyFood)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FDC == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.FDC = nil
	if err = local.L.LoadFDC(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FDC == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUsdaSRLegacyFoodToOneSetOpUsdaFoodUsingFDC(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UsdaSRLegacyFood
	var b, c UsdaFood

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, usdaSRLegacyFoodDBTypes, false, strmangle.SetComplement(usdaSRLegacyFoodPrimaryKeyColumns, usdaSRLegacyFoodColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, usdaFoodDBTypes, false, strmangle.SetComplement(usdaFoodPrimaryKeyColumns, usdaFoodColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, usdaFoodDBTypes, false, strmangle.SetComplement(usdaFoodPrimaryKeyColumns, usdaFoodColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*UsdaFood{&b, &c} {
		err = a.SetFDC(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.FDC != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.FDCUsdaSRLegacyFood != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.FDCID != x.FDCID {
			t.Error("foreign key was wrong value", a.FDCID)
		}

		if exists, err := UsdaSRLegacyFoodExists(ctx, tx, a.FDCID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}

func testUsdaSRLegacyFoodsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUsdaSRLegacyFoodsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UsdaSRLegacyFoodSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUsdaSRLegacyFoodsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UsdaSRLegacyFoods().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	usdaSRLegacyFoodDBTypes = map[string]string{`FDCID`: `integer`, `NDBNumber`: `integer`}
	_                       = bytes.MinRead
)

func testUsdaSRLegacyFoodsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(usdaSRLegacyFoodPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(usdaSRLegacyFoodAllColumns) == len(usdaSRLegacyFoodPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUsdaSRLegacyFoodsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(usdaSRLegacyFoodAllColumns) == len(usdaSRLegacyFoodPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, usdaSRLegacyFoodDBTypes, true, usdaSRLegacyFoodPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(usdaSRLegacyFoodAllColumns, usdaSRLegacyFoodPrimaryKeyColumns) {
		fields = usdaSRLegacyFoodAllColumns
	} else {
		fields = strmangle.SetComplement(
			usdaSRLegacyFoodAllColumns,
			usdaSRLegacyFoodPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := UsdaSRLegacyFoodSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUsdaSRLegacyFoodsUpsert(t *testing.T) {
	t.Parallel()

	if len(usdaSRLegacyFoodAllColumns) == len(usdaSRLegacyFoodPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := UsdaSRLegacyFood{}
	if err = randomize.Struct(seed, &o, usdaSRLegacyFoodDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UsdaSRLegacyFood: %s", err)
	}

	count, err := UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, usdaSRLegacyFoodDBTypes, false, usdaSRLegacyFoodPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaSRLegacyFood struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UsdaSRLegacyFood: %s", err)
	}

	count, err = UsdaSRLegacyFoods().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
