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
  Header,
  Divider
} from 'semantic-ui-react';
import AddRecipeNote from '../components/AddRecipeNote';
import { Link } from 'react-router-dom';
import ImageUploader from '../components/ImageUploader';
import RecipeCategoryEditor from '../components/RecipeCategoryEditor';
import RecipeEditorBasicInfo from '../components/RecipeEditorBasicInfo';
import RecipeEditorIngredientItem from '../components/RecipeEditorIngredientItem';
import RecipeEditorInstructionItem from '../components/RecipeEditorInstructionItem';

class EditorPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      slug: props.match.params.recipe_id
    };
    this.editTopLevelItem = this.editTopLevelItem.bind(this);
    this.editIngredient = this.editIngredient.bind(this);
    this.addIngredient = this.addIngredient.bind(this);
    this.deleteIngredient = this.deleteIngredient.bind(this);
    this.editInstruction = this.editInstruction.bind(this);
    this.addInstruction = this.addInstruction.bind(this);
    this.deleteInstruction = this.deleteInstruction.bind(this);
    this.getCumulativeInstructionNum = this.getCumulativeInstructionNum.bind(
      this
    );
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
        <RecipeEditorBasicInfo
          recipe={recipe}
          editTopLevelItem={this.editTopLevelItem}
        />
        <Divider />
        <Grid columns={2}>
          <Grid.Column>
            <div ref={this.handleContextRef}>
              {/*EDITOR*/}
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
                      <RecipeEditorInstructionItem
                        key={`section-${sectionNum}-instruction-${instructionNum}`}
                        sectionNum={sectionNum}
                        instructionNum={instructionNum}
                        instruction={instruction}
                        editInstruction={this.editInstruction}
                        addInstruction={this.addInstruction}
                        deleteInstruction={this.deleteInstruction}
                        getCumulativeInstructionNum={
                          this.getCumulativeInstructionNum
                        }
                      />
                    ))}
                    <h2>Ingredients</h2>
                    {section.ingredients.map((ingredient, ingredientNum) => (
                      <RecipeEditorIngredientItem
                        key={`section-${sectionNum}-ingredient-${ingredientNum}`}
                        sectionNum={sectionNum}
                        ingredientNum={ingredientNum}
                        ingredient={ingredient}
                        editIngredient={this.editIngredient}
                        addIngredient={this.addIngredient}
                        deleteIngredient={this.deleteIngredient}
                      />
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
