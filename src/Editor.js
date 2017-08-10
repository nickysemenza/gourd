import React, { Component } from "react";
class About extends Component {
  constructor(props) {
    super(props);
    this.state = {"title":"Chocolate Chip Cookies","source":"from Bouchon Bakery Cookbook","quantity":24,"servings":24,"unit":"cookies","totalMinutes":100,"equipment":["sheet pan","oven"],"sections":[{"ingredients":[{"name":"butter","grams":334}],"instructions":["cream in mixer"]},{"ingredients":[{"name":"dark brown sugar","substitute":"light brown sugar","grams":268.1,"measurement":{"unit":"cup","amount":1}},{"name":"molasses","grams":24},{"name":"sugar","grams":208}],"instructions":["add to mixer for 4 min"]},{"ingredients":[{"name":"eggs","grams":120}],"instructions":["add to mixer","mix for 30 seconds"]},{"ingredients":[{"name":"flour","modifier":"sifted","optional":true,"grams":476.1},{"name":"baking soda","measurement":{"unit":"tsp","amount":1}},{"name":"salt","measurement":{"unit":"tsp","amount":2}}],"instructions":["add in 2 additions to mixer, mix for 40 seconds inbetween"]},{"ingredients":[{"name":"chocolate chunks","modifier":"coarsely chopped","grams":214},{"name":"chocolate chips","grams":214}],"instructions":["mix to incorporate","fridge for 1-24 hours","bake at 325 for 20 min"]}]}
  }
  editTopLevelItem(fieldName,e) {
      this.setState({[fieldName]: e.target.value});
  }
  deleteSection(index) {
      this.setState((prevState) => ({
          sections: [...prevState.sections.slice(0,index), ...prevState.sections.slice(index+1)]
      }));
  }
  addSection(index) {
      let sections = this.state.sections;
      sections.splice(index, 0, {ingredients: [], instructions: []});
      this.setState({sections})
  }
  deleteInstruction(sectionNum,instructionNum) {
      let sections = this.state.sections;
      sections[sectionNum].instructions.splice(instructionNum,1);
      this.setState({sections})
  }
  addInstruction(sectionNum,instructionNum) {
      let sections = this.state.sections;
      sections[sectionNum].instructions.splice(instructionNum,0,"");
      this.setState({sections})
  }
  editInstruction(sectionNum,instructionNum,e) {
      let sections = this.state.sections;
      sections[sectionNum].instructions[instructionNum] = e.target.value;
      this.setState({sections})
  }
  deleteIngredient(sectionNum,ingredientNum) {
      let sections = this.state.sections;
      sections[sectionNum].ingredients.splice(ingredientNum,1);
      this.setState({sections})
  }
  addIngredient(sectionNum,ingredientNum) {
      let sections = this.state.sections;
      sections[sectionNum].ingredients.splice(ingredientNum,0,{
          name: 'name',
          grams: 0,
          measurement: {
              "unit": "cup",
              "amount": 1
          }
      });
      this.setState({sections})
  }
  editIngredient(sectionNum,ingredientNum,e) {
      let sections = this.state.sections;
      sections[sectionNum].ingredients[ingredientNum] = e.target.value;
      this.setState({sections})
  }
  render() {
    return (
      <div>

        <div className="form-group row">
          <label htmlFor="example-text-input" className="col-2 col-form-label">Title</label>
          <div className="col-10">
            <input className="form-control" type="text" value={this.state.title} onChange={this.editTopLevelItem.bind(this,'title')} />
          </div>
        </div>

        <div className="form-group row">
          <label htmlFor="example-text-input" className="col-2 col-form-label">source</label>
          <div className="col-10">
            <input className="form-control" type="text" value={this.state.source} onChange={this.editTopLevelItem.bind(this,'source')} />
          </div>
        </div>

        <div className="form-group row">
          <label htmlFor="example-text-input" className="col-2 col-form-label">quantity</label>
          <div className="col-10">
            <input className="form-control" type="number" value={this.state.quantity} onChange={this.editTopLevelItem.bind(this,'quantity')} />
          </div>
        </div>

        <div className="form-group row">
          <label htmlFor="example-text-input" className="col-2 col-form-label">servings</label>
          <div className="col-10">
            <input className="form-control" type="number" value={this.state.servings} onChange={this.editTopLevelItem.bind(this,'servings')} />
          </div>
        </div>

        <div className="form-group row">
          <label htmlFor="example-text-input" className="col-2 col-form-label">totalMinutes</label>
          <div className="col-10">
            <input className="form-control" type="number" value={this.state.totalMinutes} onChange={this.editTopLevelItem.bind(this,'totalMinutes')} />
          </div>
        </div>
          { this.state.sections.map((section, sectionNum)=> {
            return (<div key={sectionNum}>
              <button onClick={this.addSection.bind(this,sectionNum)}>Add New Section Before</button>
              <button onClick={this.addSection.bind(this,sectionNum+1)}>Add New Section After</button>
              <button onClick={this.deleteSection.bind(this,sectionNum)}>Delete Section</button>
              <h2>Instructions</h2>
                    { section.instructions.map((instruction, instructionNum) => <div key={`section-${sectionNum}-instruction-${instructionNum}`}>
                      <input type="text" value={instruction} onChange={this.editInstruction.bind(this,sectionNum,instructionNum)}/>
                      <a onClick={this.deleteInstruction.bind(this,sectionNum,instructionNum)}>delete</a>
                      &nbsp; | &nbsp;
                      <a onClick={this.addInstruction.bind(this,sectionNum,instructionNum)}>add before</a>
                      &nbsp; | &nbsp;
                      <a onClick={this.addInstruction.bind(this,sectionNum,instructionNum+1)}>add after</a>
                    </div>)}
              <h2>Ingredients</h2>
                { section.ingredients.map((ingredient, ingredientNum) => <div key={`section-${sectionNum}-ingredient-${ingredientNum}`}>
                    {/*<input type="text" value={ingredient} onChange={this.editIngredient.bind(this,sectionNum,ingredientNum)}/>*/}
                    <pre>{JSON.stringify(ingredient,null, 2)} {ingredientNum}</pre>

                    <a onClick={this.deleteIngredient.bind(this,sectionNum,ingredientNum)}>delete</a>
                    &nbsp; | &nbsp;
                    <a onClick={this.addIngredient.bind(this,sectionNum,ingredientNum)}>add before</a>
                    &nbsp; | &nbsp;
                    <a onClick={this.addIngredient.bind(this,sectionNum,ingredientNum+1)}>add after</a>
                </div>)}

              <hr/>
            </div>);
          })}

        {/*<pre>{JSON.stringify(this.state,null, 2)}</pre>*/}
      </div>
    );
  }
}

export default About;
