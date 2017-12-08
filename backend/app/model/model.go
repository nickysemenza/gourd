package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Recipe struct {
	Model
	Slug         string    `json:"slug" gorm:"unique"`
	Title        string    `json:"title"`
	TotalMinutes uint      `json:"total_minutes"`
	Equipment    string    `json:"equipment"`
	Source       string    `json:"source"`
	Servings     uint      `json:"servings"`
	Unit         string    `json:"unit"`
	Quantity     uint      `json:"quantity"`
	Sections     []Section `json:"sections"`
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

//func GetIngredientByName(db *gorm.DB, name string) Ingredient {
//	ingredient := Ingredient{}
//	db.FirstOrCreate(&ingredient, Ingredient{Name: name})
//	return ingredient
//}

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

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Section{}, &SectionInstruction{}, &SectionIngredient{}, &Recipe{}, &Ingredient{})
	db.Model(&Section{}).AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionInstruction{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("section_id", "sections(id)", "RESTRICT", "RESTRICT")
	db.Model(&SectionIngredient{}).AddForeignKey("item_id", "ingredients(id)", "RESTRICT", "RESTRICT")
	return db
}
func DBReset(db *gorm.DB) *gorm.DB {
	db.DropTable(&SectionIngredient{})
	db.DropTable(&SectionInstruction{})
	db.DropTable(&Section{})
	db.DropTable(&Ingredient{})
	db.DropTable(&Recipe{})

	return db
}
