package notion

import (
	"context"
	"fmt"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/volatiletech/null/v8"
)

type fakeDB struct{}

var _ notionapi.DatabaseService = &fakeDB{}

func (f *fakeDB) Get(context.Context, notionapi.DatabaseID) (*notionapi.Database, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f *fakeDB) Query(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
	var time notionapi.Date
	return &notionapi.DatabaseQueryResponse{
		Results: []notionapi.Page{
			{
				Object: notionapi.ObjectTypePage,
				ID:     "foo",
				Parent: notionapi.Parent{
					Type: "database_id",
				},
				Properties: map[string]notionapi.Property{
					"Name": &notionapi.TitleProperty{
						Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "page1title"}}},
					},
					"Date": &notionapi.DateProperty{
						Date: &notionapi.DateObject{Start: &time},
					},
					"Source": &notionapi.URLProperty{
						URL: "https://test.com",
					},
					"Tags": &notionapi.MultiSelectProperty{
						MultiSelect: []notionapi.Option{{Name: "dinner"}},
					},
				},
				URL: "https://notion.so/abc",
			},
		},
	}, nil
}
func (f *fakeDB) Update(context.Context, notionapi.DatabaseID, *notionapi.DatabaseUpdateRequest) (*notionapi.Database, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakeDB) Create(ctx context.Context, request *notionapi.DatabaseCreateRequest) (*notionapi.Database, error) {
	return nil, fmt.Errorf("not implemented")
}

type fakeBlock struct {
	children map[notionapi.BlockID][]notionapi.Block
}

var _ notionapi.BlockService = &fakeBlock{}

func (f *fakeBlock) GetChildren(_ context.Context, id notionapi.BlockID, _ *notionapi.Pagination) (*notionapi.GetChildrenResponse, error) {
	item, ok := f.children[id]
	if !ok {
		return nil, fmt.Errorf("not found: %s", id)
	}
	return &notionapi.GetChildrenResponse{
		Results: item,
	}, nil
}
func (f *fakeBlock) AppendChildren(context.Context, notionapi.BlockID, *notionapi.AppendBlockChildrenRequest) (*notionapi.AppendBlockChildrenResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakeBlock) Get(context.Context, notionapi.BlockID) (notionapi.Block, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakeBlock) Update(ctx context.Context, id notionapi.BlockID, request *notionapi.BlockUpdateRequest) (notionapi.Block, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakeBlock) Delete(ctx context.Context, id notionapi.BlockID) (notionapi.Block, error) {
	return nil, fmt.Errorf("not implemented")
}

type fakePage struct {
	details map[notionapi.PageID]*notionapi.Page
}

var _ notionapi.PageService = &fakePage{}

func (f *fakePage) Get(_ context.Context, id notionapi.PageID) (*notionapi.Page, error) {
	item, ok := f.details[id]
	if !ok {
		return nil, fmt.Errorf("not found: %s", id)
	}
	return item, nil
}
func (f *fakePage) Create(context.Context, *notionapi.PageCreateRequest) (*notionapi.Page, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakePage) Update(context.Context, notionapi.PageID, *notionapi.PageUpdateRequest) (*notionapi.Page, error) {
	return nil, fmt.Errorf("not implemented")
}

// NewFakeNotion makes a fake
func NewFakeNotion(t *testing.T) *rawClient {
	t.Helper()

	db := &fakeDB{}
	block := &fakeBlock{
		children: map[notionapi.BlockID][]notionapi.Block{
			"foo": {
				&notionapi.ImageBlock{
					BasicBlock: notionapi.BasicBlock{
						Type: notionapi.BlockTypeImage,
						ID:   "block1",
					},
					Image: notionapi.Image{
						File: &notionapi.FileObject{
							URL: "https://picsum.photos/200/400",
						},
					},
				},
				&notionapi.CodeBlock{
					BasicBlock: notionapi.BasicBlock{
						Type: notionapi.BlockTypeCode,
						ID:   "block2",
					},
					Code: notionapi.Code{
						RichText: []notionapi.RichText{
							{Text: &notionapi.Text{Content: `name: toast
---
1 recipe bar
; toast`}},
						},
					},
				},
				&notionapi.CodeBlock{
					BasicBlock: notionapi.BasicBlock{
						Type: notionapi.BlockTypeCode,
						ID:   "block3",
					},
					Code: notionapi.Code{
						RichText: []notionapi.RichText{
							{Text: &notionapi.Text{Content: "not a arecipe"}},
						},
					},
				},
				&notionapi.ChildPageBlock{
					BasicBlock: notionapi.BasicBlock{
						Object: notionapi.ObjectTypeBlock,
						Type:   notionapi.BlockTypeChildPage,
						ID:     "block4child",
					},

					ChildPage: struct {
						Title string "json:\"title\""
					}{
						Title: "child page name",
					},
				},
			},
			"foo2": {
				&notionapi.CodeBlock{
					BasicBlock: notionapi.BasicBlock{
						Type: notionapi.BlockTypeCode,
						ID:   "block2",
					},
					Code: notionapi.Code{
						RichText: []notionapi.RichText{
							{Text: &notionapi.Text{Content: `name: bar
---
1 slice bread
; eat`}},
						},
					},
				},
			},
		},
	}
	page := &fakePage{
		details: map[notionapi.PageID]*notionapi.Page{
			"block4child": {
				Object: notionapi.ObjectTypePage,
				ID:     "foo2",
				Parent: notionapi.Parent{
					Type: "page_id",
				},
				Properties: map[string]notionapi.Property{
					"title": &notionapi.TitleProperty{
						Type:  notionapi.PropertyTypeTitle,
						Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "subpagetitle"}}},
					},
					"ID": &notionapi.UniqueIDProperty{
						Type:     notionapi.PropertyTypeUniqueID,
						UniqueID: notionapi.UniqueID{Prefix: null.StringFrom("MEAL").Ptr(), Number: 452},
					},
				},
			},
		},
	}
	c := rawClient{
		db:    db,
		block: block,
		page:  page,
	}
	return &c
}
