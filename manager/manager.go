package manager

import (
	"context"
	"database/sql"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v1"
)

// Ingredient is a globally-scoped ingredient
type Ingredient struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Recipe represents a recipe
type Recipe struct {
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	TotalMinutes uint      `json:"total_minutes"`
	Equipment    string    `json:"equipment"`
	Source       string    `json:"source"`
	Servings     uint      `json:"servings"`
	Quantity     string    `json:"quantity"`
	Unit         string    `json:"unit"`
	Sections     []Section `json:"sections"`
}

// Section represents a Section
type Section struct {
	UUID         string               `json:"uuid"`
	RecipeUUID   string               `json:"recipe_uuid"`
	Minutes      uint                 `json:"minutes"`
	Ingredients  []SectionIngredient  `json:"ingredients"`
	Instructions []SectionInstruction `json:"instructions"`
}

// SectionIngredient represents a SectionIngredient
type SectionIngredient struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"` // todo: use this to load an Ingredient
	Grams       float32 `json:"grams"`
	SectionUUID string  `json:"section_uuid"`
	Amount      float32
	Unit        string
	Adjective   string
	Optional    bool

	// one of the following:
	RecipeUUID     string `json:"recipe_uuid"`
	IngredientUUID string `json:"ingredient_uuid"`
}

// SectionInstruction represents a SectionInstruction
type SectionInstruction struct {
	UUID        string `json:"uuid"`
	Instruction string `json:"instruction"`
	SectionUUID string `json:"recipe_uuid"`
}

// Manager manages recipes
type Manager struct {
	db *sqlx.DB
}

// New creates a new Manager
func New(db *sqlx.DB) *Manager {
	return &Manager{db: db}
}

// LoadFromFile loads a recipe from a file
func (m *Manager) LoadFromFile(ctx context.Context, filename string) (*Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := &Recipe{}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	if err := m.AssignUUIDs(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}
func setUUID(val string) string {
	if val != "" {
		return val
	}
	return uuid.NewV4().String()
}

// AssignUUIDs adds uuids where missing
func (m *Manager) AssignUUIDs(ctx context.Context, r *Recipe) error {
	r.UUID = setUUID(r.UUID)
	for x := range r.Sections {
		r.Sections[x].UUID = setUUID(r.Sections[x].UUID)
		r.Sections[x].RecipeUUID = r.UUID
		for y := range r.Sections[x].Ingredients {
			r.Sections[x].Ingredients[y].UUID = setUUID(r.Sections[x].Ingredients[y].UUID)
			r.Sections[x].Ingredients[y].SectionUUID = r.Sections[x].UUID
			ing, err := m.IngredientByName(ctx, r.Sections[x].Ingredients[y].Name)
			if err != nil {
				return err
			}
			r.Sections[x].Ingredients[y].IngredientUUID = ing.UUID
		}
		for y := range r.Sections[x].Instructions {
			r.Sections[x].Instructions[y].UUID = setUUID(r.Sections[x].Instructions[y].UUID)
			r.Sections[x].Instructions[y].SectionUUID = r.Sections[x].UUID
		}
	}
	return nil
}
func (m *Manager) GetRecipe(ctx context.Context, name, uuid string) (*Recipe, error) {
	r := &Recipe{}
	err := m.db.Get(r, "SELECT * FROM recipes WHERE name=$1 OR uuid = $2", name, uuid)
	if err != nil {
		return nil, err
	}

}
func (m *Manager) SaveRecipe(ctx context.Context, r *Recipe) error {
	sectionUUIDs := make([]string, len(r.Sections))
	for x, section := range r.Sections {
		sectionUUIDs[x] = section.UUID
	}
	_, err := m.db.ExecContext(ctx, `DELETE FROM recipe_section_instructions WHERE section IN ($1)`, pq.Array(sectionUUIDs))
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `DELETE FROM recipe_section_ingredients WHERE section IN ($1)`, pq.Array(sectionUUIDs))
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `DELETE FROM recipe_sections WHERE recipe = $1`, r.UUID)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `DELETE FROM recipes WHERE uuid = $1`, r.UUID)
	if err != nil {
		return err
	}

	// now recreate
	_, err = m.db.ExecContext(ctx, `INSERT INTO recipes 
(uuid, name, total_minutes, equipment, source, servings, quantity, unit)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8)`,
		r.UUID, r.Name, r.TotalMinutes, r.Equipment,
		r.Source, r.Servings, r.Quantity, r.Unit,
	)
	if err != nil {
		return err
	}
	for x, section := range r.Sections {
		_, err = m.db.ExecContext(ctx, `INSERT INTO recipe_sections 
(uuid, recipe, sort, minutes)
VALUES
($1, $2, $3, $4)`,
			section.UUID, r.UUID, x, section.Minutes,
		)
		if err != nil {
			return err
		}
		for y, i := range section.Ingredients {
			_, err = m.db.ExecContext(ctx, `INSERT INTO recipe_section_ingredients 
(uuid, section, sort, ingredient, grams, amount, unit, adjective, optional)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				i.UUID, i.SectionUUID, y, i.IngredientUUID, i.Grams,
				i.Amount, i.Unit, i.Adjective, i.Optional)
			if err != nil {
				return err
			}
		}
		for y, i := range section.Instructions {
			_, err = m.db.ExecContext(ctx, `INSERT INTO recipe_section_instructions 
(uuid, section, sort, instruction)
VALUES
($1, $2, $3, $4)`,
				i.UUID, i.SectionUUID, y, i.Instruction)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

// IngredientByName retrieves an ingredient by name, creating it if it does not exist
func (m *Manager) IngredientByName(ctx context.Context, name string) (*Ingredient, error) {
	ingredient := &Ingredient{}
	err := m.db.GetContext(ctx, ingredient, `SELECT * FROM ingredients
	WHERE name = $1 LIMIT 1`, name)
	if err == sql.ErrNoRows {
		_, err = m.db.ExecContext(ctx, `INSERT INTO ingredients (uuid, name) VALUES ($1, $2)`, setUUID(""), name)
		if err != nil {
			return nil, err
		}
		return m.IngredientByName(ctx, name)
	}
	return ingredient, err

}
