import React, { Component } from 'react';

import {
  fetchMealDetail,
  editMealMultiplier,
  editMealRecipe,
  saveMeal,
  addMealRecipe,
  deleteMealRecipe,
  editMealField
} from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import {
  Breadcrumb,
  Header,
  Table,
  Input,
  Button,
  Form
} from 'semantic-ui-react';
import { Link } from 'react-router-dom';
import 'moment-timezone';
import RecipePicker from '../components/RecipePicker';
import { getScaledMeasurementString } from '../components/RecipeIngredientMeasurement';
import DatePicker from 'react-datepicker';
import moment from 'moment';

import 'react-datepicker/dist/react-datepicker.css';

class MealDetailPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mealId: props.match.params.meal_id
    };
  }
  componentWillReceiveProps(nextProps) {
    if (nextProps.match !== this.props.match) {
      this.setState({ mealId: nextProps.match.params.meal_id });
      this.props.fetchMealDetail(nextProps.match.params.meal_id);
    }
  }
  componentDidMount() {
    this.props.fetchMealDetail(this.state.mealId);
  }

  render() {
    let thisMeal = this.props.meal_detail[this.state.mealId] || {
      recipe_meals: []
    };
    const options = [
      { key: 'l', text: 'Lunch', value: 'lunch' },
      { key: 'd', text: 'Dinner', value: 'dinner' },
      { key: 's', text: 'Snack', value: 'snack' }
    ];
    let a = {};

    thisMeal.recipe_meals.forEach(r => {
      let { sections, title } = r.recipe;
      let scale = r.multiplier;
      if (sections)
        sections.forEach(s => {
          s.ingredients.forEach(i => {
            let { name } = i.item;
            let x = `${getScaledMeasurementString(i, scale)} [${title}]`;
            if (a[name]) a[name].push(x);
            else {
              a[name] = [x];
            }
            console.log({ i });
          });
        });
    });
    return (
      <div className="container">
        <Breadcrumb>
          <Breadcrumb.Section link as={Link} to="/">
            Home
          </Breadcrumb.Section>
          <Breadcrumb.Divider icon="right angle" />
          <Breadcrumb.Section link as={Link} to="/meals">
            Meals
          </Breadcrumb.Section>
          <Breadcrumb.Divider icon="right angle" />
          <Breadcrumb.Section active>
            Meal # {this.state.mealId}
          </Breadcrumb.Section>
        </Breadcrumb>
        <pre>{JSON.stringify(a, true, 2)}</pre>
        <Header dividing content="Meal Detail" />

        <Table celled>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Recipe</Table.HeaderCell>
              <Table.HeaderCell>Multiplier</Table.HeaderCell>
              <Table.HeaderCell>Action</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {thisMeal.recipe_meals.map((eachMealRecipe, x) => (
              <Table.Row key={x}>
                <Table.Cell>
                  <RecipePicker
                    value={eachMealRecipe.recipe.slug}
                    onChange={(e, { value }) => {
                      console.log(value);
                      this.props.editMealRecipe(this.state.mealId, x, value);
                    }}
                  />
                </Table.Cell>
                <Table.Cell>
                  <Input
                    value={eachMealRecipe.multiplier}
                    onChange={e => {
                      this.props.editMealMultiplier(this.state.mealId, x, e);
                    }}
                    type="number"
                  />
                </Table.Cell>
                <Table.Cell>
                  <Button
                    icon="trash"
                    onClick={() =>
                      this.props.deleteMealRecipe(this.state.mealId, x)
                    }
                    content="delete"
                  />
                  <Button
                    icon="info circle"
                    as={Link}
                    to={`/${eachMealRecipe.recipe.slug}`}
                    content={eachMealRecipe.recipe.title}
                  />
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
        <Form>
          <Form.Input
            fluid
            label="Name"
            placeholder="Name"
            value={thisMeal.name}
            onChange={e =>
              this.props.editMealField(
                this.state.mealId,
                'name',
                e.target.value
              )
            }
          />
          <Form.Select
            options={options}
            fluid
            label="Type"
            placeholder="Type"
            value={thisMeal.type}
            onChange={(e, { value }) =>
              this.props.editMealField(this.state.mealId, 'type', value)
            }
          />
          <Form.Input
            fluid
            label="Description"
            placeholder="Description"
            value={thisMeal.description}
            onChange={e =>
              this.props.editMealField(
                this.state.mealId,
                'description',
                e.target.value
              )
            }
          />
          {/*<Form.Select fluid label='Gender' options={options} placeholder='Gender' />*/}
          <DatePicker
            selected={moment(thisMeal.time)}
            onChange={date => {
              this.props.editMealField(this.state.mealId, 'time', date);
            }}
          />
          <br />
          <Button
            icon="plus"
            onClick={() => this.props.addMealRecipe(this.state.mealId)}
            content="add item"
          />
          <Button
            icon="star"
            onClick={() => this.props.saveMeal(this.state.mealId)}
            content="save meal"
          />
        </Form>
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { meal_detail } = state.recipe;
  return { meal_detail };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchMealDetail,
      editMealMultiplier,
      editMealRecipe,
      saveMeal,
      addMealRecipe,
      deleteMealRecipe,
      editMealField
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(MealDetailPage);
