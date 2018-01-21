import React, { Component } from 'react';

import {
  fetchMealDetail,
  editMealMultiplier,
  editMealRecipe,
  saveMeal,
  addMealRecipe,
  deleteMealRecipe
} from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Breadcrumb, Header, Table, Input, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom';
import 'moment-timezone';
import RecipePicker from '../components/RecipePicker';

class MealDetailPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mealId: props.match.params.meal_id
    };
  }
  componentWillReceiveProps(nextProps) {
    if (nextProps.match !== this.props.match) {
      this.setState({ mealId: nextProps.match.params.meal_id });
      this.props.fetchMealDetail(nextProps.match.params.meal_id);
    }
  }
  componentDidMount() {
    this.props.fetchMealDetail(this.state.mealId);
  }

  render() {
    let thisMeal = this.props.meal_detail[this.state.mealId] || {
      recipe_meals: []
    };
    return (
      <div className="container">
        <Breadcrumb>
          <Breadcrumb.Section link as={Link} to="/">
            Home
          </Breadcrumb.Section>
          <Breadcrumb.Divider icon="right angle" />
          <Breadcrumb.Section link as={Link} to="/meals">
            Meals
          </Breadcrumb.Section>
          <Breadcrumb.Divider icon="right angle" />
          <Breadcrumb.Section active>
            Meal # {this.state.mealId}
          </Breadcrumb.Section>
        </Breadcrumb>

        <Header dividing content="Meal Detail" />

        <Table celled>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Recipe</Table.HeaderCell>
              <Table.HeaderCell>Multiplier</Table.HeaderCell>
              <Table.HeaderCell>Action</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {thisMeal.recipe_meals.map((eachMealRecipe, x) => (
              <Table.Row key={x}>
                <Table.Cell>
                  <Link
                    key={eachMealRecipe.recipe.id}
                    to={`/${eachMealRecipe.recipe.slug}`}
                  >
                    {eachMealRecipe.recipe.title}{' '}
                  </Link>
                  <RecipePicker
                    value={eachMealRecipe.recipe.slug}
                    onChange={(e, { value }) => {
                      console.log(value);
                      this.props.editMealRecipe(this.state.mealId, x, value);
                    }}
                  />
                </Table.Cell>
                <Table.Cell>
                  <Input
                    value={eachMealRecipe.multiplier}
                    onChange={e => {
                      this.props.editMealMultiplier(this.state.mealId, x, e);
                    }}
                    type="number"
                  />
                </Table.Cell>
                <Table.Cell>
                  <Button
                    icon="trash"
                    onClick={() =>
                      this.props.deleteMealRecipe(this.state.mealId, x)
                    }
                    content="delete"
                  />
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
        <Button
          icon="plus"
          onClick={() => this.props.addMealRecipe(this.state.mealId)}
          content="add item"
        />
        <Button
          icon="star"
          onClick={() => this.props.saveMeal(this.state.mealId)}
          content="save meal"
        />
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { meal_detail } = state.recipe;
  return { meal_detail };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchMealDetail,
      editMealMultiplier,
      editMealRecipe,
      saveMeal,
      addMealRecipe,
      deleteMealRecipe
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(MealDetailPage);
