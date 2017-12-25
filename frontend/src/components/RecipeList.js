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
            key={eachRecipe.id}
            as={Link}
            to={`/${eachRecipe.slug}`}
            image={
              eachRecipe.images === null
                ? 'http://via.placeholder.com/2000x1200'
                : eachRecipe.images[0].url
            }
            header={eachRecipe.title}
            meta="todo::categories"
            description="todo::description"
            extra={
              <p>
                <Icon name="clock" />
                {eachRecipe.total_minutes} Minutes
              </p>
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
