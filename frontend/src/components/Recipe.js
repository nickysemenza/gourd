import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import RecipeTable from './RecipeTable';
export default class Recipe extends Component {
  constructor(props) {
    super(props);
    this.state = {
      scale: 1.0
    };
    this.handleScaleChange = this.handleScaleChange.bind(this);
  }
  componentDidMount() {}
  handleScaleChange(event) {
    this.setState({ scale: parseFloat(event.target.value) });
  }
  render() {
    let recipe = this.props.recipe;
    if (!recipe) return <div>loading...</div>;
    if (recipe.error !== undefined) return <div>error! {recipe.error}</div>;

    return (
      <div>
        <h1>{recipe.title}</h1>
        <h4>
          From {recipe.source} . Makes{' '}
          <i>{parseFloat((recipe.quantity * this.state.scale).toFixed(1))}</i>{' '}
          {recipe.unit}
        </h4>
        <div className="row">
          <div className="col col-sm-12">
            <RecipeTable recipe={recipe} scale={this.state.scale} />
          </div>
        </div>
        <div className="row">
          <div className="col col-sm-6">
            <div className="card">
              <div className="card-block">
                <h4 className="card-title">Scaling</h4>
                <p className="card-text">
                  approx weight:&nbsp; TODOCALCg
                  <Link to={`/editor/${this.props.slug}`}>
                    <button>edit</button>
                  </Link>
                </p>
                <hr />
                <div className="form-group row">
                  <label htmlFor="example-text-input">Multiplier</label>
                  <br />
                  <input
                    className="form-control"
                    type="number"
                    min="0"
                    max="10"
                    step=".1"
                    value={this.state.scale}
                    onChange={this.handleScaleChange}
                    id="example-text-input"
                  />
                </div>
              </div>
              <hr />
              <ul className="list-group list-group-flush">
                <li className="list-group-item">
                  <b>Total minutes:</b>&nbsp;{recipe.totalMinutes}
                  <b>Total minutes (calculated):</b>&nbsp;TODOCALC
                </li>
                <li className="list-group-item">
                  <b>Equipment:</b>&nbsp;{recipe.equipment}
                </li>
              </ul>
            </div>
          </div>
          <div className="col col-sm-6">
            <div className="card">
              <p className="card-text">
                {recipe.image ? (
                  <img
                    src={recipe.image}
                    alt=""
                    style={{
                      maxHeight: '350px',
                      maxWidth: '350px',
                      height: 'auto',
                      width: 'auto'
                    }}
                  />
                ) : null}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
