package notion

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/require"
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
				Properties: map[string]notionapi.Property{
					"Name": &notionapi.TitleProperty{
						Title: []notionapi.RichText{{Text: notionapi.Text{Content: "test1"}}},
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
				URL: "https://notion.so/test1",
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

type fakeBlock struct{}

var _ notionapi.BlockService = &fakeBlock{}

func (f *fakeBlock) GetChildren(context.Context, notionapi.BlockID, *notionapi.Pagination) (*notionapi.GetChildrenResponse, error) {
	return &notionapi.GetChildrenResponse{
		Results: []notionapi.Block{
			&notionapi.ImageBlock{
				Type: notionapi.BlockTypeImage,
				ID:   "block1",
				Image: notionapi.Image{
					File: &notionapi.FileObject{
						URL: "https://test.com/image.jpg",
					},
				},
			},
			&notionapi.CodeBlock{
				Type: notionapi.BlockTypeCode,
				ID:   "block2",
				Code: notionapi.Code{
					Text: []notionapi.RichText{
						{Text: notionapi.Text{Content: "name: toast"}},
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
		},
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

func TestSmoke(t *testing.T) {
	db := &fakeDB{}
	block := &fakeBlock{}
	c := Client{
		db:    db,
		block: block,
	}
	res, err := c.Dump(context.Background())
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, res[0].Title, "test1")
	spew.Dump(res)
	require.Len(t, res[0].Photos, 1)
}
