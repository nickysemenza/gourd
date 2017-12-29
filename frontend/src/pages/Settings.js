import React, { Component } from 'react';

import { createRecipe } from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button, Form, Header } from 'semantic-ui-react';

class Settings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      newRecipeTitle: '',
      newRecipeSlug: ''
    };
  }
  makeProperSlug = value =>
    value
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^0-9a-z\-]/g, '');

  editNewRecipeTitle = e => {
    let { value } = e.target;
    let newRecipeSlug = this.makeProperSlug(value);
    this.setState({ newRecipeTitle: value, newRecipeSlug });
  };
  editNewRecipeSlug = e => {
    let newRecipeSlug = this.makeProperSlug(e.target.value);
    this.setState({ newRecipeSlug });
  };
  createRecipe = () => {
    this.props.createRecipe(
      this.state.newRecipeSlug,
      this.state.newRecipeTitle
    );
  };
  render() {
    return (
      <div className="container">
        <Header dividing content="Add New" />
        <Form>
          <Form.Group>
            <Form.Field width={8}>
              <label>Title</label>
              <input
                type="text"
                value={this.state.newRecipeTitle}
                onChange={this.editNewRecipeTitle.bind(this)}
              />
            </Form.Field>
            <Form.Field width={8}>
              <label>Slug</label>
              <input
                type="text"
                value={this.state.newRecipeSlug}
                onChange={this.editNewRecipeSlug.bind(this)}
              />
            </Form.Field>
          </Form.Group>
          <Button type="submit" onClick={this.createRecipe.bind(this)}>
            Add new Recipe
          </Button>
        </Form>
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
      createRecipe
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(Settings);
