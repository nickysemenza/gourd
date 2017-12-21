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
    const extra = (
      <a>
        <Icon name="clock" />
        16 Minutes
      </a>
    );
    let a = this.props.recipe_list.map(eachSlug => (
      <Card
        as={Link}
        to={`/${eachSlug}`}
        image="http://via.placeholder.com/2000x1200"
        header="Elliot Baker"
        meta="Friend"
        description="Elliot is a sound engineer living in Nashville who enjoys playing guitar and hanging with his cat."
        extra={extra}
      />
    ));
    return <Card.Group>{a}</Card.Group>;
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
