import React, { Component } from 'react';

import { fetchMealList } from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Header, Image, List, Table } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

class MealList extends Component {
  componentDidMount() {
    this.props.fetchMealList();
  }
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
                  <Header as="h3" content={meal.name} />
                  {meal.description}
                </Table.Cell>
                <Table.Cell>date</Table.Cell>
                <Table.Cell>
                  <List link>
                    {meal.recipe_meals.map(eachRM => {
                      let { recipe } = eachRM;
                      return (
                        <List.Item
                          key={recipe.id}
                          as={Link}
                          to={`/${recipe.slug}`}
                        >
                          {recipe.title} @ {eachRM.multiplier}x
                        </List.Item>
                      );
                    })}
                  </List>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
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
