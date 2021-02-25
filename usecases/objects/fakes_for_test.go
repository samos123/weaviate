//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package objects

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/entities/filters"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/entities/search"
	"github.com/semi-technologies/weaviate/usecases/projector"
	"github.com/semi-technologies/weaviate/usecases/traverser"
	"github.com/stretchr/testify/mock"
)

type fakeSchemaManager struct {
	CalledWith struct {
		fromClass string
		property  string
		toClass   string
	}
	GetSchemaResponse schema.Schema
}

func (f *fakeSchemaManager) UpdatePropertyAddDataType(ctx context.Context, principal *models.Principal,
	fromClass, property, toClass string) error {
	f.CalledWith = struct {
		fromClass string
		property  string
		toClass   string
	}{
		fromClass: fromClass,
		property:  property,
		toClass:   toClass,
	}
	return nil
}

func (f *fakeSchemaManager) GetSchema(principal *models.Principal) (schema.Schema, error) {
	return f.GetSchemaResponse, nil
}

type fakeLocks struct{}

func (f *fakeLocks) LockConnector() (func() error, error) {
	return func() error { return nil }, nil
}

func (f *fakeLocks) LockSchema() (func() error, error) {
	return func() error { return nil }, nil
}

type fakeVectorizerProvider struct {
	vectorizer *fakeVectorizer
}

func (f *fakeVectorizerProvider) Vectorizer(modName, className string) (Vectorizer, error) {
	return f.vectorizer, nil
}

type fakeVectorizer struct {
	mock.Mock
}

func (f *fakeVectorizer) UpdateObject(ctx context.Context, object *models.Object) error {
	args := f.Called(object)
	object.Vector = args.Get(0).([]float32)
	return args.Error(1)
}

func (f *fakeVectorizer) Corpi(ctx context.Context, corpi []string) ([]float32, error) {
	panic("not implemented")
}

type fakeAuthorizer struct{}

func (f *fakeAuthorizer) Authorize(principal *models.Principal, verb, resource string) error {
	return nil
}

type fakeVectorRepo struct {
	mock.Mock
}

func (f *fakeVectorRepo) Exists(ctx context.Context,
	id strfmt.UUID) (bool, error) {
	args := f.Called(id)
	return args.Bool(0), args.Error(1)
}

func (f *fakeVectorRepo) ObjectByID(ctx context.Context,
	id strfmt.UUID, props traverser.SelectProperties, additional traverser.AdditionalProperties) (*search.Result, error) {
	args := f.Called(id, props, additional)
	return args.Get(0).(*search.Result), args.Error(1)
}

func (f *fakeVectorRepo) ObjectSearch(ctx context.Context, limit int,
	filters *filters.LocalFilter, additional traverser.AdditionalProperties) (search.Results, error) {
	args := f.Called(limit, filters, additional)
	return args.Get(0).([]search.Result), args.Error(1)
}

func (f *fakeVectorRepo) PutObject(ctx context.Context,
	concept *models.Object, vector []float32) error {
	args := f.Called(concept, vector)
	return args.Error(0)
}

func (f *fakeVectorRepo) BatchPutObjects(ctx context.Context, batch BatchObjects) (BatchObjects, error) {
	args := f.Called(batch)
	return batch, args.Error(0)
}

func (f *fakeVectorRepo) AddBatchReferences(ctx context.Context, batch BatchReferences) (BatchReferences, error) {
	args := f.Called(batch)
	return batch, args.Error(0)
}

func (f *fakeVectorRepo) Merge(ctx context.Context, merge MergeDocument) error {
	args := f.Called(merge)
	return args.Error(0)
}

func (f *fakeVectorRepo) DeleteObject(ctx context.Context,
	className string, id strfmt.UUID) error {
	args := f.Called(className, id)
	return args.Error(0)
}

func (f *fakeVectorRepo) AddReference(ctx context.Context,
	class string, source strfmt.UUID, prop string,
	ref *models.SingleRef) error {
	args := f.Called(source, prop, ref)
	return args.Error(0)
}

type fakeExtender struct {
	single *search.Result
	multi  []search.Result
}

func (f *fakeExtender) Single(ctx context.Context, in *search.Result, limit *int) (*search.Result, error) {
	return f.single, nil
}

func (f *fakeExtender) Multi(ctx context.Context, in []search.Result, limit *int) ([]search.Result, error) {
	return f.multi, nil
}

type fakeProjector struct {
	multi []search.Result
}

func (f *fakeProjector) Reduce(in []search.Result, params *projector.Params) ([]search.Result, error) {
	return f.multi, nil
}