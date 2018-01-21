import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Dropdown } from 'semantic-ui-react';
import { connect } from 'react-redux';
import { fetchRecipes } from '../actions/recipe';
import { bindActionCreators } from 'redux';

class RecipePicker extends Component {
  componentDidMount() {
    this.props.fetchRecipes();
  }
  render() {
    let options = this.props.recipe_list.map(x => {
      return {
        key: x.slug,
        text: x.title,
        value: x.slug
      };
    });
    return (
      <Dropdown
        placeholder="Select Recipe"
        fluid
        search
        selection
        options={options}
        value={this.props.value}
        onChange={this.props.onChange}
      />
    );
  }
}
RecipePicker.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired
};
function mapStateToProps(state) {
  let { recipe_list } = state.recipe;
  return { recipe_list };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchRecipes
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(RecipePicker);
