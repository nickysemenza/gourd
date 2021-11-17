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

func testUsdaFoodAttributes(t *testing.T) {
	t.Parallel()

	query := UsdaFoodAttributes()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUsdaFoodAttributesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
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

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaFoodAttributesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := UsdaFoodAttributes().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaFoodAttributesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UsdaFoodAttributeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsdaFoodAttributesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UsdaFoodAttributeExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if UsdaFoodAttribute exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UsdaFoodAttributeExists to return true, but got false.")
	}
}

func testUsdaFoodAttributesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	usdaFoodAttributeFound, err := FindUsdaFoodAttribute(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if usdaFoodAttributeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUsdaFoodAttributesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = UsdaFoodAttributes().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUsdaFoodAttributesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := UsdaFoodAttributes().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUsdaFoodAttributesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	usdaFoodAttributeOne := &UsdaFoodAttribute{}
	usdaFoodAttributeTwo := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, usdaFoodAttributeOne, usdaFoodAttributeDBTypes, false, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}
	if err = randomize.Struct(seed, usdaFoodAttributeTwo, usdaFoodAttributeDBTypes, false, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = usdaFoodAttributeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = usdaFoodAttributeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UsdaFoodAttributes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUsdaFoodAttributesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	usdaFoodAttributeOne := &UsdaFoodAttribute{}
	usdaFoodAttributeTwo := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, usdaFoodAttributeOne, usdaFoodAttributeDBTypes, false, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}
	if err = randomize.Struct(seed, usdaFoodAttributeTwo, usdaFoodAttributeDBTypes, false, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = usdaFoodAttributeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = usdaFoodAttributeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func usdaFoodAttributeBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func usdaFoodAttributeAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UsdaFoodAttribute) error {
	*o = UsdaFoodAttribute{}
	return nil
}

func testUsdaFoodAttributesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &UsdaFoodAttribute{}
	o := &UsdaFoodAttribute{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute object: %s", err)
	}

	AddUsdaFoodAttributeHook(boil.BeforeInsertHook, usdaFoodAttributeBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeBeforeInsertHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.AfterInsertHook, usdaFoodAttributeAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeAfterInsertHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.AfterSelectHook, usdaFoodAttributeAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeAfterSelectHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.BeforeUpdateHook, usdaFoodAttributeBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeBeforeUpdateHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.AfterUpdateHook, usdaFoodAttributeAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeAfterUpdateHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.BeforeDeleteHook, usdaFoodAttributeBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeBeforeDeleteHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.AfterDeleteHook, usdaFoodAttributeAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeAfterDeleteHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.BeforeUpsertHook, usdaFoodAttributeBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeBeforeUpsertHooks = []UsdaFoodAttributeHook{}

	AddUsdaFoodAttributeHook(boil.AfterUpsertHook, usdaFoodAttributeAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	usdaFoodAttributeAfterUpsertHooks = []UsdaFoodAttributeHook{}
}

func testUsdaFoodAttributesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsdaFoodAttributesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(usdaFoodAttributeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsdaFoodAttributeToOneUsdaFoodUsingFDC(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UsdaFoodAttribute
	var foreign UsdaFood

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, usdaFoodDBTypes, false, usdaFoodColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFood struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.FDCID, foreign.FDCID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.FDC().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.FDCID, foreign.FDCID) {
		t.Errorf("want: %v, got %v", foreign.FDCID, check.FDCID)
	}

	slice := UsdaFoodAttributeSlice{&local}
	if err = local.L.LoadFDC(ctx, tx, false, (*[]*UsdaFoodAttribute)(&slice), nil); err != nil {
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

func testUsdaFoodAttributeToOneUsdaFoodAttributeTypeUsingFoodAttributeType(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UsdaFoodAttribute
	var foreign UsdaFoodAttributeType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, usdaFoodAttributeTypeDBTypes, false, usdaFoodAttributeTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttributeType struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.FoodAttributeTypeID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.FoodAttributeType().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := UsdaFoodAttributeSlice{&local}
	if err = local.L.LoadFoodAttributeType(ctx, tx, false, (*[]*UsdaFoodAttribute)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FoodAttributeType == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.FoodAttributeType = nil
	if err = local.L.LoadFoodAttributeType(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.FoodAttributeType == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUsdaFoodAttributeToOneSetOpUsdaFoodUsingFDC(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UsdaFoodAttribute
	var b, c UsdaFood

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, usdaFoodAttributeDBTypes, false, strmangle.SetComplement(usdaFoodAttributePrimaryKeyColumns, usdaFoodAttributeColumnsWithoutDefault)...); err != nil {
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

		if x.R.FDCUsdaFoodAttributes[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.FDCID, x.FDCID) {
			t.Error("foreign key was wrong value", a.FDCID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FDCID))
		reflect.Indirect(reflect.ValueOf(&a.FDCID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.FDCID, x.FDCID) {
			t.Error("foreign key was wrong value", a.FDCID, x.FDCID)
		}
	}
}

func testUsdaFoodAttributeToOneRemoveOpUsdaFoodUsingFDC(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UsdaFoodAttribute
	var b UsdaFood

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, usdaFoodAttributeDBTypes, false, strmangle.SetComplement(usdaFoodAttributePrimaryKeyColumns, usdaFoodAttributeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, usdaFoodDBTypes, false, strmangle.SetComplement(usdaFoodPrimaryKeyColumns, usdaFoodColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetFDC(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveFDC(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.FDC().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.FDC != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.FDCID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.FDCUsdaFoodAttributes) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testUsdaFoodAttributeToOneSetOpUsdaFoodAttributeTypeUsingFoodAttributeType(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UsdaFoodAttribute
	var b, c UsdaFoodAttributeType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, usdaFoodAttributeDBTypes, false, strmangle.SetComplement(usdaFoodAttributePrimaryKeyColumns, usdaFoodAttributeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, usdaFoodAttributeTypeDBTypes, false, strmangle.SetComplement(usdaFoodAttributeTypePrimaryKeyColumns, usdaFoodAttributeTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, usdaFoodAttributeTypeDBTypes, false, strmangle.SetComplement(usdaFoodAttributeTypePrimaryKeyColumns, usdaFoodAttributeTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*UsdaFoodAttributeType{&b, &c} {
		err = a.SetFoodAttributeType(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.FoodAttributeType != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.FoodAttributeTypeUsdaFoodAttributes[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.FoodAttributeTypeID, x.ID) {
			t.Error("foreign key was wrong value", a.FoodAttributeTypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FoodAttributeTypeID))
		reflect.Indirect(reflect.ValueOf(&a.FoodAttributeTypeID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.FoodAttributeTypeID, x.ID) {
			t.Error("foreign key was wrong value", a.FoodAttributeTypeID, x.ID)
		}
	}
}

func testUsdaFoodAttributeToOneRemoveOpUsdaFoodAttributeTypeUsingFoodAttributeType(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UsdaFoodAttribute
	var b UsdaFoodAttributeType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, usdaFoodAttributeDBTypes, false, strmangle.SetComplement(usdaFoodAttributePrimaryKeyColumns, usdaFoodAttributeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, usdaFoodAttributeTypeDBTypes, false, strmangle.SetComplement(usdaFoodAttributeTypePrimaryKeyColumns, usdaFoodAttributeTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetFoodAttributeType(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveFoodAttributeType(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.FoodAttributeType().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.FoodAttributeType != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.FoodAttributeTypeID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.FoodAttributeTypeUsdaFoodAttributes) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testUsdaFoodAttributesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
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

func testUsdaFoodAttributesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UsdaFoodAttributeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUsdaFoodAttributesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UsdaFoodAttributes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	usdaFoodAttributeDBTypes = map[string]string{`ID`: `integer`, `FDCID`: `integer`, `SeqNum`: `integer`, `FoodAttributeTypeID`: `integer`, `Name`: `text`, `Value`: `text`}
	_                        = bytes.MinRead
)

func testUsdaFoodAttributesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(usdaFoodAttributePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(usdaFoodAttributeAllColumns) == len(usdaFoodAttributePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUsdaFoodAttributesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(usdaFoodAttributeAllColumns) == len(usdaFoodAttributePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UsdaFoodAttribute{}
	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, usdaFoodAttributeDBTypes, true, usdaFoodAttributePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(usdaFoodAttributeAllColumns, usdaFoodAttributePrimaryKeyColumns) {
		fields = usdaFoodAttributeAllColumns
	} else {
		fields = strmangle.SetComplement(
			usdaFoodAttributeAllColumns,
			usdaFoodAttributePrimaryKeyColumns,
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

	slice := UsdaFoodAttributeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUsdaFoodAttributesUpsert(t *testing.T) {
	t.Parallel()

	if len(usdaFoodAttributeAllColumns) == len(usdaFoodAttributePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := UsdaFoodAttribute{}
	if err = randomize.Struct(seed, &o, usdaFoodAttributeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UsdaFoodAttribute: %s", err)
	}

	count, err := UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, usdaFoodAttributeDBTypes, false, usdaFoodAttributePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UsdaFoodAttribute struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UsdaFoodAttribute: %s", err)
	}

	count, err = UsdaFoodAttributes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
