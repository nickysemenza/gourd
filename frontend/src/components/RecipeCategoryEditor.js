import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import {
  addCategoryToRecipe,
  fetchCategories,
  removeCategoryFromRecipe
} from '../actions/recipe';
import { Button, List } from 'semantic-ui-react';

class RecipeCategoryEditor extends Component {
  componentDidMount() {
    this.props.fetchCategories();
    this.doesRecipeHaveCategoryID = this.doesRecipeHaveCategoryID.bind(this);
  }
  doesRecipeHaveCategoryID = categoryId =>
    this.props.recipe_detail[this.props.slug].categories.filter(
      x => x.id === categoryId
    ).length === 1;
  render() {
    return (
      <div>
        <List divided verticalAlign="middle">
          {this.props.category_list.map((c, categoryListIndex) => (
            <List.Item>
              <List.Content floated="right">
                <Button
                  disabled={!this.doesRecipeHaveCategoryID(c.id)}
                  onClick={() =>
                    this.props.removeCategoryFromRecipe(this.props.slug, c.id)
                  }
                >
                  Remove
                </Button>
              </List.Content>
              <List.Content floated="right">
                <Button
                  disabled={this.doesRecipeHaveCategoryID(c.id)}
                  onClick={() =>
                    this.props.addCategoryToRecipe(this.props.slug, c.id)
                  }
                >
                  Add
                </Button>
              </List.Content>
              <List.Content>{c.name}</List.Content>
            </List.Item>
          ))}
        </List>
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { category_list, recipe_detail } = state.recipe;
  return { category_list, recipe_detail };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchCategories,
      removeCategoryFromRecipe,
      addCategoryToRecipe
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(
  RecipeCategoryEditor
);
RecipeCategoryEditor.propTypes = {
  slug: PropTypes.string.isRequired
};
