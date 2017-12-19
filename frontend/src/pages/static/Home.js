import React, { Component } from "react";
import { Link } from "react-router-dom";

import {
    fetchRecipes
} from '../../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux'

//TODO: make RecipeList component

class Home extends Component {
  componentDidMount() {
      this.props.fetchRecipes();
  }
  render() {
    let a = this.props.recipe_list.map(a => (
      <div key={a}><Link to={`/${a}`}>{a}</Link></div>
    ));
    return (
      <div>
        <h2>Nicky's Recipe Stash</h2>
        <hr />
        <div className="row">
          <div className="col">
            {a}
          </div>
          <div className="w-100 hidden-md-up" />
          <div className="col">
            <iframe
              src="https://www.nicky.photos/frame/slideshow?key=NrsTC5&autoStart=1&captions=1&navigation=1&playButton=0&randomize=1&speed=3&transition=fade&transitionSpeed=1"
              width="100%"
              height="600"
              frameBorder="no"
              scrolling="no"
              title="smugmug gallery"
            />
          </div>
        </div>
        <hr />
        View this project on github:
        {" "}
        <a href="https://github.com/nickysemenza/food" target="_blank" rel="noopener noreferrer">
          github.com/nickysemenza/food
        </a>
        <br />
        <a href="http://www.nicky.photos/Food/My-Food/" target="_blank" rel="noopener noreferrer">
          View all photos
        </a>
      </div>
    );
  }
}


function mapStateToProps (state) {
    return {
        recipe_list: state.recipe.recipe_list,
    };
}

const mapDispatchToProps = (dispatch) => {
    return bindActionCreators({
        fetchRecipes
    }, dispatch)
};

export default connect(mapStateToProps, mapDispatchToProps)(Home);