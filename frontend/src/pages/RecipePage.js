import React, { Component } from "react";
import Recipe from "../components/Recipe";
export default class RecipePage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      recipe: null,
      slug: props.match.params.recipe_id
    };
  }
  componentWillReceiveProps(nextProps) {
    console.log(nextProps, this.props);
    if (nextProps.match !== this.props.match) {
      this.setState({ slug: nextProps.match.params.recipe_id });
      this.fetchData(nextProps.match.params.recipe_id);
    }
  }
  componentDidMount() {
    this.fetchData(this.state.slug);
  }
  fetchData(r) {
    fetch(`/api/recipes/${r}`, { accept: "application/json" })
      .then(response => response.json())
      .then(json => this.setState({ recipe: json }));
  }
  render() {
    return (
      <div className="container">
        <Recipe recipe={this.state.recipe} slug={this.state.slug} />
      </div>
    );
  }
}
