import React, { Component } from 'react';
import {GramMeasurement, VolumeMeasurement} from './RecipeIngredientMeasurement';
export default class Recipe extends Component {
    constructor(props) {
        super(props);
        this.state = {
            scale: 1.0
        };
        this.handleScaleChange = this.handleScaleChange.bind(this);
    }
    componentDidMount() {

    }
    handleScaleChange(event) {
        this.setState({scale: parseFloat(event.target.value)});
    }
    render () {
        let recipe = this.props.recipe;
        if(!recipe) return(<div>loading...</div>);
        let totalWeight = 0;
        let instructionNum = 1;
        let tableRows = recipe.sections.map((recipeSection,num)=>{
            let ingredientList = [], weightList = [], volumeList = [];

            recipeSection.ingredients.forEach((i,n)=>{
                ingredientList.push(<div className="ingredientCellItem" key={n}>{i.name}</div>);
                if(i.grams) totalWeight+=i.grams;
                weightList.push(<div className="ingredientCellItem" key={n}>
                    <GramMeasurement grams={i.grams} scale={this.state.scale}/>
                </div>);
                volumeList.push(<div className="ingredientCellItem" key={n}>
                    <VolumeMeasurement measurement={i.measurement} scale={this.state.scale}/>
                </div>);
            });

            let instructionList = recipeSection.instructions.map((i,n)=><div key={n}>{`${instructionNum++}. ${i}`}</div>);

            return <tr key={`r-${num}`}>
                <td style={{verticalAlign: 'middle'}}>{String.fromCharCode(num+65)}.</td>
                <td>{ingredientList}</td>
                <td>{weightList}</td>
                <td>{volumeList}</td>{/*<td style={{"borderRight":"solid 1px grey"}}>*/}
                <td>{instructionList}</td>
            </tr>
        });

        return (
            <div>
                <h1>{recipe.title}</h1>
                <h4>{recipe.source}. Makes <i>{recipe.quantity * this.state.scale}</i> {recipe.unit}</h4>
                <div className="row">
                    <div className="col-9">
                        <table className="table table-sm ">
                            <thead className="thead-default">
                            <tr>
                                <td>&nbsp;</td>
                                <th>ingredients</th>
                                <th>weight</th>
                                <th>volume</th>
                                <th>steps</th>
                            </tr>
                            </thead>
                            <tbody>
                            {tableRows}
                            </tbody>
                        </table>
                    </div>


                    <div className="col">
                        <div className="card">
                            <div className="card-block">
                                <h4 className="card-title">Scaling</h4>
                                <p className="card-text">Some quick example text to build on the card title and make up the bulk of the card's content.</p>
                                <p>approx weight: <b>{totalWeight*this.state.scale}g</b></p>
                                {/*Scaling ratio: <input type="text" value={this.state.scale} onChange={this.handleScaleChange} />*/}

                                <div className="form-group row">
                                    <label htmlFor="example-text-input" className="col-4 col-form-label">Multiplier</label>
                                    <div className="col-7">
                                        <input className="form-control" type="number"  min="0" value={this.state.scale} onChange={this.handleScaleChange} id="example-text-input"/>
                                    </div>
                                </div>
                            </div>
                            <ul className="list-group list-group-flush">
                                <li className="list-group-item"><b>Total minutes:</b>&nbsp;{recipe.totalMinutes}</li>
                                <li className="list-group-item"><b>Equipment:</b>&nbsp;{recipe.equipment.join(', ')}</li>
                            </ul>
                            <div className="card-block">
                                <a href="#" className="card-link">Card link</a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

        );
    }
}

