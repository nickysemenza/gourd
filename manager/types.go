package manager

import (
	"github.com/nickysemenza/food/db"
	"gopkg.in/guregu/null.v3/zero"
)

// Recipe is the presentation for a Recipe
type Recipe struct {
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	TotalMinutes int64     `json:"total_minutes"`
	Equipment    string    `json:"equipment"`
	Source       string    `json:"source"`
	Quantity     int64     `json:"quantity"`
	Servings     int64     `json:"servings"`
	Unit         string    `json:"unit"`
	Sections     []Section `json:"sections"`
}

// Section is the presentation for a Sections
type Section struct {
	Minutes      int64               `json:"minutes"`
	Ingredients  []SectionIngredient `json:"ingredients"`
	Instructions []Instruction       `json:"instructions"`
}

// SectionIngredient is the presentation for a Ingredient
type SectionIngredient struct {
	Name      string  `json:"name"`
	Grams     float64 `json:"grams"`
	Amount    float64
	Unit      string
	Adjective string
	Optional  bool
}

// Instruction is the presentation for a Instruction
type Instruction struct {
	Instruction string `json:"instruction"`
}

func (r *Recipe) toDB() *db.Recipe {
	dbr := db.Recipe{
		UUID:         r.UUID,
		Name:         r.Name,
		TotalMinutes: zero.IntFrom(r.TotalMinutes),
		Equipment:    zero.StringFrom(r.Equipment),
		Source:       zero.StringFrom(r.Source),
		Servings:     zero.IntFrom(r.Servings),
		Quantity:     zero.IntFrom(r.Quantity),
		Unit:         zero.StringFrom(r.Unit),
	}
	for _, s := range r.Sections {
		dbs := db.Section{
			Minutes: zero.IntFrom(s.Minutes),
		}
		for _, i := range s.Instructions {
			dbs.Instructions = append(dbs.Instructions, db.SectionInstruction{
				Instruction: i.Instruction,
			})

		}
		for _, i := range s.Ingredients {
			dbs.Ingredients = append(dbs.Ingredients, db.SectionIngredient{
				Name:      i.Name,
				Grams:     zero.FloatFrom(i.Grams),
				Amount:    zero.FloatFrom(i.Amount),
				Unit:      zero.StringFrom(i.Unit),
				Adjective: zero.StringFrom(i.Adjective),
				Optional:  zero.BoolFrom(i.Optional),
			})
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}
	return &dbr
}

// FromRecipe transforms a db record
func FromRecipe(dbr *db.Recipe) *Recipe {
	r := Recipe{
		UUID:         dbr.UUID,
		Name:         dbr.Name,
		TotalMinutes: dbr.TotalMinutes.Int64,
		Equipment:    dbr.Equipment.String,
		Source:       dbr.Source.String,
		Servings:     dbr.Servings.Int64,
		Quantity:     dbr.Quantity.Int64,
		Unit:         dbr.Unit.String,
	}
	for _, dbs := range dbr.Sections {
		s := Section{
			Minutes: dbs.Minutes.Int64,
		}
		for _, i := range dbs.Instructions {
			s.Instructions = append(s.Instructions, Instruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range dbs.Ingredients {
			s.Ingredients = append(s.Ingredients, SectionIngredient{
				Grams:     i.Grams.Float64,
				Name:      i.Name,
				Amount:    i.Amount.Float64,
				Unit:      i.Unit.String,
				Adjective: i.Adjective.String,
				Optional:  i.Optional.Bool,
			})
		}

		r.Sections = append(r.Sections, s)
	}
	return &r
}
