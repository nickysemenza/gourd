import React, { Component } from 'react';
import {Link} from 'react-router-dom';
export default class Home extends Component {

    constructor(props) {
        super(props);
        this.state = {
            recipeList: []
        }
    }
    componentDidMount() {
        fetch(`/api/recipes`, {accept: 'application/json'})
            .then((response) => response.json())
            .then((json) => this.setState({recipeList: json}));
    }
    render () {

        let a = this.state.recipeList.map(a=><li key={a}><Link to={`/${a}`}>{a}</Link></li>);
        return (
            <div>
                <h2>Recipes:</h2>
                {a}
            </div>
        );
    }
}

