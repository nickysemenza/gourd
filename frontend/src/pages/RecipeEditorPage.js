import React, { Component } from 'react';
import Recipe from '../components/Recipe';

import * as RecipeActionCreators from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

class EditorPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      slug: props.match.params.recipe_id
    };
  }
  componentWillReceiveProps(nextProps) {
    console.log(nextProps, this.props);
    if (nextProps.match !== this.props.match) {
      this.setState({ slug: nextProps.match.params.recipe_id });
      this.props.fetchRecipeDetail(nextProps.match.params.recipe_id);
    }
  }
  componentDidMount() {
    this.props.fetchRecipeDetail(this.state.slug);
  }
  editTopLevelItem(fieldName, e) {
    this.props.editTopLevelItem(this.state.slug, fieldName, e.target.value);
  }
  deleteSection(index) {
    this.props.deleteSectionByIndex(this.state.slug, index);
  }
  addSection(index) {
    this.props.addSection(this.state.slug, index);
  }
  deleteInstruction(sectionNum, instructionNum) {
    this.props.deleteInstruction(this.state.slug, sectionNum, instructionNum);
  }
  addInstruction(sectionNum, instructionNum) {
    this.props.addInstruction(this.state.slug, sectionNum, instructionNum);
  }
  editInstruction(sectionNum, instructionNum, e) {
    this.props.editInstruction(
      this.state.slug,
      sectionNum,
      instructionNum,
      e.target.value
    );
  }
  deleteIngredient(sectionNum, ingredientNum) {
    this.props.deleteIngredient(this.state.slug, sectionNum, ingredientNum);
  }
  addIngredient(sectionNum, ingredientNum) {
    this.props.addIngredient(this.state.slug, sectionNum, ingredientNum);
  }
  editIngredient(sectionNum, ingredientNum, field, e) {
    this.props.editIngredient(
      this.state.slug,
      sectionNum,
      ingredientNum,
      field,
      e.target.value
    );
  }
  render() {
    const recipe = this.props.recipe_detail[this.state.slug];
    if (!recipe) return <div>loading...</div>;
    if (recipe.error !== undefined) return <div>error! {recipe.error}</div>;
    return (
      <div className="row">
        <div className="col-md-6">
          {/*EDITOR*/}
          <div className="form-group row">
            <label
              htmlFor="example-text-input"
              className="col-2 col-form-label"
            >
              Title
            </label>
            <div className="col-10">
              <input
                className="form-control"
                type="text"
                value={recipe.title}
                onChange={this.editTopLevelItem.bind(this, 'title')}
              />
            </div>
          </div>

          <div className="form-group row">
            <label
              htmlFor="example-text-input"
              className="col-2 col-form-label"
            >
              source
            </label>
            <div className="col-10">
              <input
                className="form-control"
                type="text"
                value={recipe.source}
                onChange={this.editTopLevelItem.bind(this, 'source')}
              />
            </div>
          </div>

          <div className="form-group row">
            <label
              htmlFor="example-text-input"
              className="col-2 col-form-label"
            >
              quantity
            </label>
            <div className="col-10">
              <input
                className="form-control"
                type="number"
                value={recipe.quantity}
                onChange={this.editTopLevelItem.bind(this, 'quantity')}
              />
            </div>
          </div>

          <div className="form-group row">
            <label
              htmlFor="example-text-input"
              className="col-2 col-form-label"
            >
              servings
            </label>
            <div className="col-10">
              <input
                className="form-control"
                type="number"
                value={recipe.servings}
                onChange={this.editTopLevelItem.bind(this, 'servings')}
              />
            </div>
          </div>

          <div className="form-group row">
            <label
              htmlFor="example-text-input"
              className="col-2 col-form-label"
            >
              total_minutes
            </label>
            <div className="col-10">
              <input
                className="form-control"
                type="number"
                value={recipe.total_minutes}
                onChange={this.editTopLevelItem.bind(this, 'total_minutes')}
              />
            </div>
          </div>
          {recipe.sections.map((section, sectionNum) => {
            return (
              <div key={sectionNum}>
                <button onClick={this.addSection.bind(this, sectionNum)}>
                  Add New Section Before
                </button>
                <button onClick={this.addSection.bind(this, sectionNum + 1)}>
                  Add New Section After
                </button>
                <button onClick={this.deleteSection.bind(this, sectionNum)}>
                  Delete Section
                </button>
                <h2>Instructions</h2>
                {section.instructions.map((instruction, instructionNum) => (
                  <div
                    key={`section-${sectionNum}-instruction-${instructionNum}`}
                  >
                    <input
                      type="text"
                      value={instruction.name}
                      onChange={this.editInstruction.bind(
                        this,
                        sectionNum,
                        instructionNum
                      )}
                    />
                    <a
                      onClick={this.deleteInstruction.bind(
                        this,
                        sectionNum,
                        instructionNum
                      )}
                    >
                      delete
                    </a>
                    &nbsp; | &nbsp;
                    <a
                      onClick={this.addInstruction.bind(
                        this,
                        sectionNum,
                        instructionNum
                      )}
                    >
                      add before
                    </a>
                    &nbsp; | &nbsp;
                    <a
                      onClick={this.addInstruction.bind(
                        this,
                        sectionNum,
                        instructionNum + 1
                      )}
                    >
                      add after
                    </a>
                  </div>
                ))}
                <h2>Ingredients</h2>
                {section.ingredients.map((ingredient, ingredientNum) => (
                  <div
                    key={`section-${sectionNum}-ingredient-${ingredientNum}`}
                  >
                    <div className="form-group row">
                      <label
                        htmlFor="example-text-input"
                        className="col-2 col-form-label"
                      >
                        name
                      </label>
                      <div className="col-10">
                        <input
                          className="form-control"
                          type="text"
                          value={ingredient.item.name}
                          onChange={this.editIngredient.bind(
                            this,
                            sectionNum,
                            ingredientNum,
                            'item'
                          )}
                        />
                      </div>
                    </div>
                    <div className="form-group row">
                      <label
                        htmlFor="example-text-input"
                        className="col-2 col-form-label"
                      >
                        Grams
                      </label>
                      <div className="col-10">
                        <input
                          className="form-control"
                          type="number"
                          value={ingredient.grams}
                          onChange={this.editIngredient.bind(
                            this,
                            sectionNum,
                            ingredientNum,
                            'grams'
                          )}
                        />
                      </div>
                    </div>
                    <div className="form-group row">
                      <label
                        htmlFor="example-text-input"
                        className="col-2 col-form-label"
                      >
                        unit
                      </label>
                      <div className="col-10">
                        <input
                          className="form-control"
                          type="text"
                          value={ingredient.amount_unit}
                          onChange={this.editIngredient.bind(
                            this,
                            sectionNum,
                            ingredientNum,
                            'amount_unit'
                          )}
                        />
                      </div>
                    </div>
                    <div className="form-group row">
                      <label
                        htmlFor="example-text-input"
                        className="col-2 col-form-label"
                      >
                        amount
                      </label>
                      <div className="col-10">
                        <input
                          className="form-control"
                          type="number"
                          value={ingredient.amount}
                          onChange={this.editIngredient.bind(
                            this,
                            sectionNum,
                            ingredientNum,
                            'amount'
                          )}
                        />
                      </div>
                    </div>
                    {/*<pre>{JSON.stringify(ingredient,null, 2)} {ingredientNum}</pre>*/}
                    <a
                      onClick={this.deleteIngredient.bind(
                        this,
                        sectionNum,
                        ingredientNum
                      )}
                    >
                      delete
                    </a>
                    &nbsp; | &nbsp;
                    <a
                      onClick={this.addIngredient.bind(
                        this,
                        sectionNum,
                        ingredientNum
                      )}
                    >
                      add before
                    </a>
                    &nbsp; | &nbsp;
                    <a
                      onClick={this.addIngredient.bind(
                        this,
                        sectionNum,
                        ingredientNum + 1
                      )}
                    >
                      add after
                    </a>
                  </div>
                ))}

                <hr />
              </div>
            );
          })}
          <hr />
          {/*<pre>{JSON.stringify(recipe,null, 2)}</pre>*/}
        </div>
        <div className="col-md-6" style={{ position: 'fixed', left: '50%' }}>
          <Recipe recipe={recipe} />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    recipe_detail: state.recipe.recipe_detail
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(RecipeActionCreators, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(EditorPage);
