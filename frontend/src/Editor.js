import React, { Component } from "react";
import Recipe from './Recipe';
import YAML from 'yamljs';
class About extends Component {
  constructor(props) {
    super(props);
    this.state = {"id":5,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"slug":"double-chocolate-chunk-cookies","title":"Double Chocolate Chunk Cookies","total_minutes":100,"equipment":"sheet pan, oven","source":"from Bouchon Bakery Cookbook","servings":24,"unit":"cookies","quantity":24,"sections":[{"id":24,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"sort_order":0,"ingredients":[{"id":52,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"item":{"id":20,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"butter"},"item_id":20,"grams":334,"amount":0,"amount_unit":"","substitute":"","modifier":"","optional":false,"section_id":24}],"instructions":[{"id":37,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"name":"cream in mixer","section_id":24}],"recipe_id":5,"minutes":0},{"id":25,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"sort_order":0,"ingredients":[{"id":53,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"item":{"id":21,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"dark brown sugar"},"item_id":21,"grams":268.1,"amount":1,"amount_unit":"cup","substitute":"light brown sugar","modifier":"","optional":false,"section_id":25},{"id":54,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:39Z","deleted_at":null,"item":{"id":22,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"molasses"},"item_id":22,"grams":24,"amount":0,"amount_unit":"","substitute":"","modifier":"","optional":false,"section_id":25},{"id":55,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":23,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"sugar"},"item_id":23,"grams":208,"amount":0,"amount_unit":"","substitute":"","modifier":"","optional":false,"section_id":25}],"instructions":[{"id":38,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"add to mixer for 4 min","section_id":25}],"recipe_id":5,"minutes":0},{"id":26,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"sort_order":0,"ingredients":[{"id":56,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":17,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"eggs"},"item_id":17,"grams":120,"amount":0,"amount_unit":"","substitute":"","modifier":"","optional":false,"section_id":26}],"instructions":[{"id":39,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"add to mixer","section_id":26},{"id":40,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"mix for 30 seconds","section_id":26}],"recipe_id":5,"minutes":0},{"id":27,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"sort_order":0,"ingredients":[{"id":57,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":1,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"flour"},"item_id":1,"grams":380,"amount":0,"amount_unit":"","substitute":"","modifier":"sifted","optional":false,"section_id":27},{"id":58,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":27,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"cocoa powder"},"item_id":27,"grams":96,"amount":0,"amount_unit":"","substitute":"","modifier":"sifted","optional":false,"section_id":27},{"id":59,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":24,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"baking soda"},"item_id":24,"grams":0,"amount":1,"amount_unit":"tsp","substitute":"","modifier":"","optional":false,"section_id":27},{"id":60,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":13,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"salt"},"item_id":13,"grams":0,"amount":2,"amount_unit":"tsp","substitute":"","modifier":"","optional":false,"section_id":27}],"instructions":[{"id":41,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"add in 2 additions to mixer, mix for 40 seconds inbetween","section_id":27}],"recipe_id":5,"minutes":0},{"id":28,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"sort_order":0,"ingredients":[{"id":61,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":25,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"chocolate chunks"},"item_id":25,"grams":214,"amount":0,"amount_unit":"","substitute":"","modifier":"coarsely chopped","optional":false,"section_id":28},{"id":62,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"item":{"id":26,"created_at":"2017-11-30T01:43:39Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"chocolate chips"},"item_id":26,"grams":214,"amount":0,"amount_unit":"","substitute":"","modifier":"","optional":false,"section_id":28}],"instructions":[{"id":42,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"mix to incorporate","section_id":28},{"id":43,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"fridge for 1-24 hours","section_id":28},{"id":44,"created_at":"2017-11-30T01:43:40Z","updated_at":"2017-11-30T01:43:40Z","deleted_at":null,"name":"bake at 325 for 20 min","section_id":28}],"recipe_id":5,"minutes":0}]}
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
      sections[sectionNum].instructions.splice(instructionNum,0,{name: ""});
      this.setState({sections})
  }
  editInstruction(sectionNum,instructionNum,e) {
      let sections = this.state.sections;
      sections[sectionNum].instructions[instructionNum]["name"] = e.target.value;
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
          item: {name: 'name'},
          grams: 0,
          amount_unit: "cup",
          amount: 1,
          substitute: "",
          modifier: "",
          optional: false
      });
      this.setState({sections})
  }
  editIngredient(sectionNum,ingredientNum,field,e) {
      let sections = this.state.sections;
      if(field==="item")
          sections[sectionNum].ingredients[ingredientNum]["item"]["name"] = e.target.value;
      else
          sections[sectionNum].ingredients[ingredientNum][field] = e.target.value;
      // sections[sectionNum].ingredients[ingredientNum] = e.target.value;
      this.setState({sections})
  }
  render() {
      const recipe = this.state;
    return (
      <div className="row">
          <div className="col-md-6">
          {/*EDITOR*/}
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
          <label htmlFor="example-text-input" className="col-2 col-form-label">total_minutes</label>
          <div className="col-10">
            <input className="form-control" type="number" value={this.state.total_minutes} onChange={this.editTopLevelItem.bind(this,'total_minutes')} />
          </div>
        </div>
          { this.state.sections.map((section, sectionNum)=> {
            return (<div key={sectionNum}>
              <button onClick={this.addSection.bind(this,sectionNum)}>Add New Section Before</button>
              <button onClick={this.addSection.bind(this,sectionNum+1)}>Add New Section After</button>
              <button onClick={this.deleteSection.bind(this,sectionNum)}>Delete Section</button>
              <h2>Instructions</h2>
                    { section.instructions.map((instruction, instructionNum) => <div key={`section-${sectionNum}-instruction-${instructionNum}`}>
                      <input type="text" value={instruction.name} onChange={this.editInstruction.bind(this,sectionNum,instructionNum)}/>
                      <a onClick={this.deleteInstruction.bind(this,sectionNum,instructionNum)}>delete</a>
                      &nbsp; | &nbsp;
                      <a onClick={this.addInstruction.bind(this,sectionNum,instructionNum)}>add before</a>
                      &nbsp; | &nbsp;
                      <a onClick={this.addInstruction.bind(this,sectionNum,instructionNum+1)}>add after</a>
                    </div>)}
              <h2>Ingredients</h2>
                { section.ingredients.map((ingredient, ingredientNum) => <div key={`section-${sectionNum}-ingredient-${ingredientNum}`}>

                    <div className="form-group row">
                        <label htmlFor="example-text-input" className="col-2 col-form-label">name</label>
                        <div className="col-10">
                            <input className="form-control" type="text" value={ingredient.item.name} onChange={this.editIngredient.bind(this,sectionNum,ingredientNum,"item")} />
                        </div>
                    </div>

                    <div className="form-group row">
                        <label htmlFor="example-text-input" className="col-2 col-form-label">Grams</label>
                        <div className="col-10">
                            <input className="form-control" type="number" value={ingredient.grams} onChange={this.editIngredient.bind(this,sectionNum,ingredientNum,"grams")} />
                        </div>
                    </div>

                    <div className="form-group row">
                        <label htmlFor="example-text-input" className="col-2 col-form-label">unit</label>
                        <div className="col-10">
                            <input className="form-control" type="text" value={ingredient.amount_unit} onChange={this.editIngredient.bind(this,sectionNum,ingredientNum,"amount_unit")} />
                        </div>
                    </div>

                    <div className="form-group row">
                        <label htmlFor="example-text-input" className="col-2 col-form-label">amount</label>
                        <div className="col-10">
                            <input className="form-control" type="number" value={ingredient.amount} onChange={this.editIngredient.bind(this,sectionNum,ingredientNum,"amount")} />
                        </div>
                    </div>


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
              <hr/>
          <pre>{JSON.stringify(recipe,null, 2)}</pre>
          </div>
          <div className="col-md-6" style={{position: 'fixed', left: '50%'}}>
              <Recipe recipe={recipe}/>

          </div>
      </div>
    );
  }
}

export default About;
