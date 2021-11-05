package manager

import (
	"context"

	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
	"github.com/nickysemenza/gourd/image"
	"github.com/nickysemenza/gourd/notion"
	"github.com/nickysemenza/gourd/photos"
	"github.com/nickysemenza/gourd/rs_client"
	"gopkg.in/guregu/null.v4/zero"
)

// Manager manages recipes
type Manager struct {
	db         *db.Client
	Google     *google.Client
	Photos     *photos.Photos
	Auth       *auth.Auth
	R          *rs_client.Client
	Notion     *notion.Client
	ImageStore image.Store
}

func New(db *db.Client, g *google.Client, auth *auth.Auth, r *rs_client.Client, notion *notion.Client, imageStore image.Store) *Manager {
	return &Manager{db: db,
		Google:     g,
		Auth:       auth,
		Photos:     photos.New(db, g, imageStore),
		R:          r,
		Notion:     notion,
		ImageStore: imageStore,
	}
}

func (m *Manager) DB() *db.Client {
	return m.db
}

func (m *Manager) ProcessGoogleAuth(ctx context.Context, code string) (jwt string, rawUser map[string]interface{}, err error) {
	err = m.Google.Finish(ctx, code)
	if err != nil {
		return
	}
	user, err := m.Google.GetUserInfo(ctx)
	if err != nil {
		return
	}
	rawUser = map[string]interface{}{"raw": user}

	jwt, err = m.Auth.GetJWT(user)
	return
}

func (m *Manager) SyncNotionToMeals(ctx context.Context) error {
	nMeals, err := m.Notion.Dump(ctx)
	if err != nil {
		return err
	}

	var foo []db.NotionRecipe
	var bar []db.NotionImage
	var img []db.Image
	for _, nMeal := range nMeals {
		foo = append(foo, db.NotionRecipe{
			PageID:    nMeal.PageID,
			PageTitle: nMeal.Title,
			AteAt:     zero.TimeFromPtr(nMeal.Time),
		})
		for _, nPhoto := range nMeal.Photos {
			bh, image, err := image.GetBlurHash(ctx, nPhoto.URL)
			if err != nil {
				return err
			}
			id := common.ID("notion_image")
			err = m.ImageStore.SaveImage(ctx, id, image)
			if err != nil {
				return err
			}

			i := db.Image{
				ID:       id,
				BlurHash: bh,
				Source:   "notion",
			}
			img = append(img, i)
			bar = append(bar, db.NotionImage{
				PageID:  nMeal.PageID,
				BlockID: nPhoto.BlockID,
				ImageID: id,
			})
		}
	}

	err = m.db.SaveImage(ctx, img)
	if err != nil {
		return err
	}

	err = m.db.SaveNotionRecipes(ctx, foo)
	if err != nil {
		return err
	}
	err = m.db.UpsertNotionImages(ctx, bar)
	if err != nil {
		return err
	}
	return nil
}
