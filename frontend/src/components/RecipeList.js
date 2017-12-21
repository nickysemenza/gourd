import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import { Card, Icon } from 'semantic-ui-react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchRecipes } from '../actions/recipe';

class RecipeList extends Component {
  componentDidMount() {
    this.props.fetchRecipes();
  }
  render() {
    return (
      <Card.Group>
        {this.props.recipe_list.map(eachRecipe => (
          <Card
            as={Link}
            to={`/${eachRecipe.slug}`}
            image="http://via.placeholder.com/2000x1200"
            header={eachRecipe.title}
            meta="todo::categories"
            description="todo::description"
            extra={
              <a>
                <Icon name="clock" />
                {eachRecipe.total_minutes} Minutes
              </a>
            }
          />
        ))}
      </Card.Group>
    );
  }
}

function mapStateToProps(state) {
  return {
    recipe_list: state.recipe.recipe_list
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchRecipes
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(RecipeList);
