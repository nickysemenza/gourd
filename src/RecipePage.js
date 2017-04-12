import React, { Component } from 'react';
import Recipe from './Recipe';
export default class RecipePage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            recipe: null,
            slug: props.match.params.recipe_id
        }
    }
    componentDidMount() {
        fetch(`/api/recipes/${this.state.slug}`, {accept: 'application/json'})
            .then((response) => response.json())
            .then((json) => this.setState({recipe: json}));
    }

    render () {
        // let slug = this.props.match.params.recipe_id;
        return (
            <div>
                <Recipe recipe={this.state.recipe} slug={this.state.slug}/>
            </div>
        );
    }
}
