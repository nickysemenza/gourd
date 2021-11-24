package notion

import (
	"context"
	"fmt"
	"testing"

	"github.com/jomei/notionapi"
)

type fakeDB struct{}

var _ notionapi.DatabaseService = &fakeDB{}

func (f *fakeDB) Get(context.Context, notionapi.DatabaseID) (*notionapi.Database, error) {
	return nil, fmt.Errorf("not implemented")
}
func (f *fakeDB) List(context.Context, *notionapi.Pagination) (*notionapi.DatabaseListResponse, error) {
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
						Title: []notionapi.RichText{{Text: notionapi.Text{Content: "page1title"}}},
					},
					"Date": &notionapi.DateProperty{
						Date: notionapi.DateObject{Start: &time},
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

func NewFakeNotion(t *testing.T) *Client {
	t.Helper()

	db := &fakeDB{}
	block := &fakeBlock{
		children: map[notionapi.BlockID][]notionapi.Block{
			"foo": []notionapi.Block{
				&notionapi.ImageBlock{
					Type: notionapi.BlockTypeImage,
					ID:   "block1",
					Image: notionapi.Image{
						File: &notionapi.FileObject{
							URL: "https://picsum.photos/200/400",
						},
					},
				},
				&notionapi.CodeBlock{
					Type: notionapi.BlockTypeCode,
					ID:   "block2",
					Code: notionapi.Code{
						Text: []notionapi.RichText{
							{Text: notionapi.Text{Content: `name: toast
---
1 recipe bar
; toast`}},
						},
					},
				},
				&notionapi.CodeBlock{
					Type: notionapi.BlockTypeCode,
					ID:   "block3",
					Code: notionapi.Code{
						Text: []notionapi.RichText{
							{Text: notionapi.Text{Content: "not a arecipe"}},
						},
					},
				},
				&notionapi.ChildPageBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeChildPage,
					ID:     "block4child",
					ChildPage: struct {
						Title string "json:\"title\""
					}{
						Title: "child page name",
					},
				},
			},
			"foo2": []notionapi.Block{
				&notionapi.CodeBlock{
					Type: notionapi.BlockTypeCode,
					ID:   "block2",
					Code: notionapi.Code{
						Text: []notionapi.RichText{
							{Text: notionapi.Text{Content: `name: bar
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
			"block4child": &notionapi.Page{
				Object: notionapi.ObjectTypePage,
				ID:     "foo2",
				Parent: notionapi.Parent{
					Type: "page_id",
				},
				Properties: map[string]notionapi.Property{
					"title": &notionapi.TitleProperty{
						Title: []notionapi.RichText{{Text: notionapi.Text{Content: "subpagetitle"}}},
					},
				},
			},
		},
	}
	c := Client{
		db:    db,
		block: block,
		page:  page,
	}
	return &c
}
