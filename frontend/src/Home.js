import React, { Component } from "react";
import { Link } from "react-router-dom";
export default class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      recipeList: []
    };
  }
  componentDidMount() {
    fetch(`/api/recipes`, { accept: "application/json" })
      .then(response => response.json())
      .then(json => this.setState({ recipeList: json }));
  }
  render() {
    let a = this.state.recipeList.map(a => (
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
