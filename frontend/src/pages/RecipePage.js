import React, { Component } from "react";
import Recipe from "../components/Recipe";
import {API_BASE_URL} from "../config";

import {
    fetchRecipeDetail
} from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux'


class RecipePage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      slug: props.match.params.recipe_id
    };
  }
  componentWillReceiveProps(nextProps) {
    console.log(nextProps, this.props);
    if (nextProps.match !== this.props.match) {
      this.setState({ slug: nextProps.match.params.recipe_id });
        this.props.fetchRecipeDetail(nextProps.match.params.recipe_id);
    }
  }
  componentDidMount() {
      this.props.fetchRecipeDetail(this.state.slug);
  }
  render() {
      let {slug} = this.state;
    return (
      <div className="container">
        <Recipe recipe={this.props.recipe_detail[slug]} slug={slug} />
      </div>
    );
  }
}


function mapStateToProps (state) {
    return {
        recipe_detail: state.recipe.recipe_detail,
    };
}

const mapDispatchToProps = (dispatch) => {
    return bindActionCreators({
        fetchRecipeDetail
    }, dispatch)
};

export default connect(mapStateToProps, mapDispatchToProps)(RecipePage);
