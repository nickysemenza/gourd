package model

import (
	"time"

	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Recipe struct {
	Model
	Slug         string       `json:"slug" gorm:"unique"`
	Title        string       `json:"title"`
	TotalMinutes uint         `json:"total_minutes"`
	Equipment    string       `json:"equipment"`
	Source       string       `json:"source"`
	Servings     uint         `json:"servings"`
	Unit         string       `json:"unit"`
	Quantity     uint         `json:"quantity"`
	Sections     []Section    `json:"sections"`
	Notes        []RecipeNote `json:"notes"`
	Images       []Image      `json:"images" gorm:"many2many:recipe_images;"`
	Categories   []Category   `json:"categories" gorm:"many2many:recipe_categories;"`
	RecipeMeal   []RecipeMeal `json:"recipe_meals"`
}
type Section struct {
	Model
	SortOrder    uint                 `json:"sort_order"`
	Ingredients  []SectionIngredient  `json:"ingredients"`
	Instructions []SectionInstruction `json:"instructions"`
	RecipeID     uint                 `json:"recipe_id"`
	Minutes      uint                 `json:"minutes"`
}
type SectionInstruction struct {
	Model
	Name      string `json:"name"`
	SectionID uint   `json:"section_id"`
}
type SectionIngredient struct {
	Model
	Item       Ingredient `json:"item"`
	ItemID     uint       `json:"item_id"`
	Grams      float32    `json:"grams"`
	Amount     float32    `json:"amount"`
	Unit       string     `json:"amount_unit"`
	Substitute string     `json:"substitute"`
	Modifier   string     `json:"modifier"`
	Optional   bool       `json:"optional"`
	SectionID  uint       `json:"section_id"`
}
type Ingredient struct {
	Model
	Name string `json:"name"`
}
type RecipeNote struct {
	Model
	Body     string `json:"body" sql:"type:text"`
	RecipeID uint   `json:"recipe_id"`
}
type Image struct {
	Model
	Path             string   `json:"path"`
	OriginalFileName string   `json:"original_name"`
	IsInS3           bool     `json:"in_s3"`
	Recipes          []Recipe `json:"recipes" gorm:"many2many:recipe_images;"`
	Md5Hash          string   `json:"md5"`
}
type Category struct {
	Model
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes" gorm:"many2many:recipe_categories;"`
}
type Meal struct {
	Model
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	RecipeMeal  []RecipeMeal `json:"recipe_meals"`
	Time        time.Time    `json:"time"`
}
type RecipeMeal struct {
	Recipe     Recipe `json:"recipe"`
	RecipeID   uint
	Meal       Meal `json:"meal"`
	MealID     uint
	Multiplier uint `json:"multiplier" gorm:"default:1"`
}

func (i *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	var url string
	if i.IsInS3 {
		url = "https://" + os.Getenv("S3_BUCKET") + ".s3.amazonaws.com/" + i.Path
	} else {
		url = os.Getenv("API_PUBLIC_URL") + "/public/" + i.Path
	}
	return json.Marshal(&struct {
		FullURL string `json:"url"`
		*Alias
	}{
		FullURL: url,
		Alias:   (*Alias)(i),
	})
}

//add the recipe to a Meal, with specified multiplier
func (recipe *Recipe) AddToMeal(db *gorm.DB, meal *Meal, multiplier uint) {
	recipemeal := RecipeMeal{}
	recipemeal.Multiplier = multiplier
	recipemeal.RecipeID = recipe.ID
	recipemeal.MealID = meal.ID
	db.Create(&recipemeal)
}
func (ingredient *Ingredient) FindOrCreateUsingName(db *gorm.DB) {
	db.FirstOrCreate(&ingredient, Ingredient{Name: ingredient.Name})
}
func (ingredient *Ingredient) GetFresh(db *gorm.DB) {
	db.First(&ingredient, ingredient.ID)
}

func (ingredient *Ingredient) AfterCreate() (err error) {
	log.Printf("[ingredient] created: %s, %d", ingredient.Name, ingredient.ID)
	return
}
func GetRecipeIDFromSlug(db *gorm.DB, slug string) (ID uint, err error) {
	recipe := Recipe{}
	if err := db.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		return 0, err
	}
	return recipe.ID, nil
}

func (updatedRecipe Recipe) CreateOrUpdate(db *gorm.DB, recursivelyStripIDs bool) {
	//todo: ensure that we aren't overwriting something with same slug, by checking for presence of ID
	for x := range updatedRecipe.Sections {
		eachSection := &updatedRecipe.Sections[x]
		if recursivelyStripIDs {
			eachSection.ID = 0
		}
		for y := range eachSection.Ingredients {
			eachSectionIngredient := &eachSection.Ingredients[y]
			eachItem := &eachSectionIngredient.Item
			if recursivelyStripIDs {
				eachSectionIngredient.ID = 0
				eachItem.ID = 0
			}
			if eachItem.ID == 0 {
				//	new ingredient!
				//	find by name, to see if we have existing
				eachItem.FindOrCreateUsingName(db)
				//eachItem = model.GetIngredientByName(e.DB, eachItem.Name)
				log.Printf("[ingredient] %s does not have an ID, giving it %d: ", eachItem.Name, eachItem.ID)
			} else {
				//	get fresh obj via eachIngredient.ID
				fresh := *eachItem
				fresh.GetFresh(db)
				//	if eachIngredient.Name != fresh.Name IT WAS MUTATED AAH!
				if eachItem.Name != fresh.Name {
					log.Printf("[ingredient] name of #%d was mutated! (%s->%s)", eachItem.ID, eachItem.Name, fresh.Name)
					//we want to preserve the original eachItem; create new w/ eachItem.Name

					// find by name, or create new
					newItem := Ingredient{Name: eachItem.Name}
					newItem.FindOrCreateUsingName(db)
					eachSectionIngredient.Item = newItem
				}
			}

		}
		for y := range eachSection.Instructions {
			eachSectionInstruction := &eachSection.Instructions[y]
			if recursivelyStripIDs {
				eachSectionInstruction.ID = 0
			}
		}
		//update Ingredients and Instructions relations
		db.Model(&eachSection).Association("Ingredients").Replace(eachSection.Ingredients)
		db.Model(&eachSection).Association("Instructions").Replace(eachSection.Instructions)
	}
	if recursivelyStripIDs {
		updatedRecipe.ID = 0
	}
	//update Categories and Sections relations
	db.Model(&updatedRecipe).Association("Categories").Replace(updatedRecipe.Categories)
	db.Model(&updatedRecipe).Association("Sections").Replace(updatedRecipe.Sections)

	if err := db.Save(&updatedRecipe).Error; err != nil {
		log.Println(err)
	}
}

var ModelsInOrder = []interface{}{
	&Recipe{},
	&Ingredient{},
	&Image{},
	&RecipeNote{},
	&Meal{},
	&RecipeMeal{},
	&Section{},
	&SectionIngredient{},
	&SectionInstruction{},
	&Category{},
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	for y := range ModelsInOrder {
		db.AutoMigrate(ModelsInOrder[y])
	}
	db.Model(&Section{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&RecipeNote{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionInstruction{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("item_id", "ingredients(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_categories").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_categories").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_images").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_images").AddForeignKey("image_id", "images(id)", "RESTRICT", "RESTRICT")

	db.Table("recipe_meals").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Table("recipe_meals").AddForeignKey("meal_id", "meals(id)", "RESTRICT", "RESTRICT")
	return db
}
func DBReset(db *gorm.DB) *gorm.DB {
	db.Exec("SET foreign_key_checks = 0;")
	for i := len(ModelsInOrder) - 1; i >= 0; i-- {
		db.DropTable(ModelsInOrder[i])
	}
	db.Exec("SET foreign_key_checks = 1;")
	return db
}
