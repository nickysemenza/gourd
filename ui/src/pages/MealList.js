import React, { Component } from 'react';

import { fetchMealList } from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Header, List, Table } from 'semantic-ui-react';
import { Link } from 'react-router-dom';
import Moment from 'react-moment';
import 'moment-timezone';
import MealCalendar from '../components/MealCalendar';

class MealList extends Component {
  componentDidMount() {
    this.props.fetchMealList();
  }
  buildRecipeWithMultiplierListForMeal = meal =>
    (meal.recipe_meals === null ? [] : meal.recipe_meals).map(eachRM => {
      let { recipe } = eachRM;
      return (
        <List.Item key={recipe.id} as={Link} to={`/${recipe.slug}`}>
          {recipe.title} @ {eachRM.multiplier}x
        </List.Item>
      );
    });

  render() {
    return (
      <div className="container">
        <Header dividing content="Meals" />
        <Table celled>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Meal</Table.HeaderCell>
              <Table.HeaderCell>Date</Table.HeaderCell>
              <Table.HeaderCell>Recipes</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {this.props.meal_list.map(meal => (
              <Table.Row key={meal.id}>
                <Table.Cell>
                  <Header
                    as={Link}
                    to={`/meal/${meal.id}`}
                    content={meal.name}
                  />
                  <pre>{meal.description}</pre>
                </Table.Cell>
                <Table.Cell>
                  <Moment tz="America/Los_Angeles" format="ddd MMM Do YYYY">
                    {meal.time}
                  </Moment>
                </Table.Cell>
                <Table.Cell>
                  <List link>
                    {this.buildRecipeWithMultiplierListForMeal(meal)}
                  </List>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
        <Header dividing content="Calendar" />
        <MealCalendar
          meal_list={this.props.meal_list}
          buildRecipeWithMultiplierListForMeal={
            this.buildRecipeWithMultiplierListForMeal
          }
        />
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { meal_list } = state.recipe;
  return { meal_list };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchMealList
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(MealList);
