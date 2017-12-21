import React, { Component } from 'react';
import Recipe from '../components/Recipe';

import * as RecipeActionCreators from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Grid, Form, Segment, Button, Input, Sticky } from 'semantic-ui-react';
import RecipeTable from '../components/RecipeTable';

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
    const { contextRef } = this.state;
    const recipe = this.props.recipe_detail[this.state.slug];
    if (!recipe) return <div>loading...</div>;
    if (recipe.error !== undefined) return <div>error! {recipe.error}</div>;
    return (
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
                    onChange={this.editTopLevelItem.bind(this, 'total_minutes')}
                  />
                </Form.Field>
              </Form.Group>
            </Form>
            {recipe.sections.map((section, sectionNum) => {
              return (
                <Segment key={sectionNum}>
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
                      key={`section-${sectionNum}-instruction-${instructionNum}`}
                    >
                      <Grid.Column width={8}>
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
                      <Grid.Column>
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
            {/*<pre>{JSON.stringify(recipe,null, 2)}</pre>*/}
          </div>
        </Grid.Column>
        <Grid.Column>
          {/*<RecipeTable recipe={recipe} scale={1} />*/}
          <Sticky context={contextRef}>
            <div style={{ marginTop: '8em' }}>
              <Recipe recipe={recipe} />
            </div>
          </Sticky>
        </Grid.Column>
      </Grid>
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
