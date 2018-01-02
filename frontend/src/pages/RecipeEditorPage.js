import React, { Component } from 'react';
import Recipe from '../components/Recipe';

import * as RecipeActionCreators from '../actions/recipe';
import { connect } from 'react-redux';
import {
  Grid,
  Form,
  Segment,
  Button,
  Input,
  Sticky,
  Label,
  Header
} from 'semantic-ui-react';
import AddRecipeNote from '../components/AddRecipeNote';
import { Link } from 'react-router-dom';
import ImageUploader from '../components/ImageUploader';
import RecipeCategoryEditor from '../components/RecipeCategoryEditor';

class EditorPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      slug: props.match.params.recipe_id
    };
  }
  handleContextRef = contextRef => this.setState({ contextRef });
  componentWillReceiveProps(nextProps) {
    console.log(nextProps, this.props);
    if (nextProps.match !== this.props.match) {
      this.setState({ slug: nextProps.match.params.recipe_id });
      this.props.fetchRecipeDetail(nextProps.match.params.recipe_id);
    }
  }
  componentDidMount() {
    this.fetchData();
  }
  fetchData() {
    this.props.fetchRecipeDetail(this.state.slug);
  }
  editTopLevelItem(fieldName, e) {
    let { value } = e.target;
    if (['quantity', 'servings', 'total_minutes'].includes(fieldName)) {
      value = parseFloat(value);
      if (isNaN(value)) value = 0;
    }
    this.props.editTopLevelItem(this.state.slug, fieldName, value);
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
    let { value } = e.target;
    if (field === 'grams' || field === 'amount') {
      value = parseFloat(value);
      if (isNaN(value)) value = 0;
    }
    this.props.editIngredient(
      this.state.slug,
      sectionNum,
      ingredientNum,
      field,
      value
    );
  }
  getCumulativeInstructionNum(sectionNum, instructionNum) {
    let r;
    r = this.props.recipe_detail[this.state.slug];
    if (!r) r = [];
    let num = 1;
    for (let x = 0; x < sectionNum; x++)
      num += r.sections[x].instructions.length;
    return num + instructionNum;
  }
  saveRecipe() {
    this.props.saveRecipe(this.state.slug);
  }
  render() {
    const { contextRef } = this.state;
    const recipe = this.props.recipe_detail[this.state.slug];
    if (!recipe) return <div>loading...</div>;
    if (recipe.error !== undefined) return <div>error! {recipe.error}</div>;
    return (
      <div>
        <Header as="h2" dividing>
          EDITING: {recipe.title}
          <Header.Subheader>
            <Link to={`/${recipe.slug}`}>Go to recipe</Link>
            <p>
              TODO::edit categories,{' '}
              {recipe.categories === null
                ? 'no categories'
                : recipe.categories.map(x => x.name).join(', ')}
            </p>
          </Header.Subheader>
        </Header>
        {recipe.sections.length === 0 ? (
          <Button
            icon="star"
            onClick={this.addSection.bind(this, 0)}
            content="initialize recipe"
          />
        ) : null}
        <Button
          icon="star"
          onClick={this.saveRecipe.bind(this)}
          content="save recipe"
        />
        <Grid columns={2}>
          <Grid.Column>
            <div ref={this.handleContextRef}>
              {/*EDITOR*/}
              <Form>
                <Form.Group>
                  <Form.Field width={8}>
                    <label>Title</label>
                    <input
                      type="text"
                      value={recipe.title}
                      onChange={this.editTopLevelItem.bind(this, 'title')}
                    />
                  </Form.Field>
                  <Form.Field width={8}>
                    <label>Source</label>
                    <input
                      type="text"
                      value={recipe.source}
                      onChange={this.editTopLevelItem.bind(this, 'source')}
                    />
                  </Form.Field>
                </Form.Group>
                <Form.Group>
                  <Form.Field width={8}>
                    <label>Quantity</label>
                    <input
                      type="number"
                      value={recipe.quantity}
                      onChange={this.editTopLevelItem.bind(this, 'quantity')}
                    />
                  </Form.Field>
                  <Form.Field width={8}>
                    <label>Quantity Unit</label>
                    <input
                      type="text"
                      value={recipe.unit}
                      onChange={this.editTopLevelItem.bind(this, 'unit')}
                    />
                  </Form.Field>
                </Form.Group>
                <Form.Group>
                  <Form.Field width={8}>
                    <label>Servings</label>
                    <input
                      type="number"
                      value={recipe.servings}
                      onChange={this.editTopLevelItem.bind(this, 'servings')}
                    />
                  </Form.Field>
                  <Form.Field width={8}>
                    <label>Total Minutes</label>
                    <input
                      type="number"
                      value={recipe.total_minutes}
                      onChange={this.editTopLevelItem.bind(
                        this,
                        'total_minutes'
                      )}
                    />
                  </Form.Field>
                </Form.Group>
              </Form>
              {recipe.sections.map((section, sectionNum) => {
                return (
                  <Segment key={sectionNum}>
                    <Label as="a" color="red" ribbon>
                      {String.fromCharCode(sectionNum + 65)}
                    </Label>
                    <Button.Group>
                      <Button
                        icon="arrow up"
                        onClick={this.addSection.bind(this, sectionNum)}
                        content="New Section"
                      />
                      <Button
                        icon="arrow down"
                        onClick={this.addSection.bind(this, sectionNum + 1)}
                        content="New Section"
                      />
                      <Button
                        icon="trash"
                        onClick={this.deleteSection.bind(this, sectionNum)}
                      />
                    </Button.Group>
                    <h2>Instructions</h2>
                    {section.instructions.map((instruction, instructionNum) => (
                      <Grid
                        padded={false}
                        key={`section-${sectionNum}-instruction-${instructionNum}`}
                      >
                        <Grid.Column width={2} className="shortGridColumn">
                          <Label horizontal>
                            {this.getCumulativeInstructionNum(
                              sectionNum,
                              instructionNum
                            )}
                          </Label>
                        </Grid.Column>
                        <Grid.Column width={8} className="shortGridColumn">
                          <Input
                            fluid
                            type="text"
                            value={instruction.name}
                            onChange={this.editInstruction.bind(
                              this,
                              sectionNum,
                              instructionNum
                            )}
                          />
                        </Grid.Column>
                        <Grid.Column className="shortGridColumn">
                          <Button.Group>
                            <Button
                              icon="arrow up"
                              onClick={this.addInstruction.bind(
                                this,
                                sectionNum,
                                instructionNum
                              )}
                              content="New"
                            />
                            <Button
                              icon="arrow down"
                              onClick={this.addInstruction.bind(
                                this,
                                sectionNum,
                                instructionNum + 1
                              )}
                              content="New"
                            />
                            <Button
                              icon="trash"
                              onClick={this.deleteInstruction.bind(
                                this,
                                sectionNum,
                                instructionNum
                              )}
                            />
                          </Button.Group>
                        </Grid.Column>
                      </Grid>
                    ))}
                    <h2>Ingredients</h2>
                    {section.ingredients.map((ingredient, ingredientNum) => (
                      <Segment
                        key={`section-${sectionNum}-ingredient-${ingredientNum}`}
                      >
                        <Label as="a" color="purple" ribbon>
                          {ingredient.item.name}
                        </Label>
                        <Form>
                          <Form.Group>
                            <Form.Field width={8}>
                              <label>Name</label>
                              <input
                                type="text"
                                value={ingredient.item.name}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'item'
                                )}
                              />
                            </Form.Field>
                            <Form.Field width={8}>
                              <label>Grams</label>
                              <Input
                                label={{ basic: true, content: 'g' }}
                                labelPosition="right"
                                type="number"
                                value={ingredient.grams}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'grams'
                                )}
                              />
                            </Form.Field>
                          </Form.Group>
                          <Form.Group>
                            <Form.Field width={8}>
                              <label>Amount</label>
                              <input
                                type="number"
                                value={ingredient.amount}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'amount'
                                )}
                              />
                            </Form.Field>
                            <Form.Field width={8}>
                              <label>Amount Unit</label>
                              <input
                                type="text"
                                value={ingredient.amount_unit}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'amount_unit'
                                )}
                              />
                            </Form.Field>
                          </Form.Group>
                          <Form.Group>
                            <Form.Field width={8}>
                              <label>Modifier</label>
                              <input
                                type="text"
                                value={ingredient.modifier}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'modifier'
                                )}
                              />
                            </Form.Field>
                            <Form.Field width={8}>
                              <label>Substitute</label>
                              <input
                                type="text"
                                value={ingredient.substitute}
                                onChange={this.editIngredient.bind(
                                  this,
                                  sectionNum,
                                  ingredientNum,
                                  'substitute'
                                )}
                              />
                            </Form.Field>
                          </Form.Group>
                          <Button.Group>
                            <Button
                              icon="arrow up"
                              onClick={this.addIngredient.bind(
                                this,
                                sectionNum,
                                ingredientNum
                              )}
                              content="New"
                            />
                            <Button
                              icon="arrow down"
                              onClick={this.addIngredient.bind(
                                this,
                                sectionNum,
                                ingredientNum + 1
                              )}
                              content="New"
                            />
                            <Button
                              icon="trash"
                              onClick={this.deleteIngredient.bind(
                                this,
                                sectionNum,
                                ingredientNum
                              )}
                            />
                          </Button.Group>
                        </Form>
                      </Segment>
                    ))}
                  </Segment>
                );
              })}
            </div>
          </Grid.Column>
          <Grid.Column>
            <Sticky context={contextRef}>
              <div style={{ marginTop: '8em' }}>
                <Recipe recipe={recipe} />
                <Header as="h2" dividing content="Add Images" />
                <ImageUploader
                  slug={this.state.slug}
                  onSuccessfulUpload={this.fetchData.bind(this)}
                />
                <Header as="h2" dividing content="Add Note" />
                <AddRecipeNote slug={this.state.slug} />
                <Header as="h2" dividing content="Categories" />
                <RecipeCategoryEditor slug={this.state.slug} />
              </div>
            </Sticky>
          </Grid.Column>
        </Grid>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    recipe_detail: state.recipe.recipe_detail
  };
}
export default connect(mapStateToProps, RecipeActionCreators)(EditorPage);
