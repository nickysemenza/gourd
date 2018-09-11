import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import RecipeTable from './RecipeTable';
import {
  Button,
  Card,
  Divider,
  Form,
  Grid,
  Header,
  Image,
  Segment
} from 'semantic-ui-react';
import Moment from 'react-moment';
import 'moment-timezone';
export default class Recipe extends Component {
  constructor(props) {
    super(props);
    this.state = {
      scale: 1.0,
      showRaw: false
    };
    this.handleScaleChange = this.handleScaleChange.bind(this);
  }
  componentDidMount() {}
  handleScaleChange(event) {
    this.setState({ scale: parseFloat(event.target.value) });
  }
  handleToggleShowRaw = () => this.setState({ showRaw: !this.state.showRaw });
  render() {
    let recipe = this.props.recipe;
    if (!recipe) return <div>loading...</div>;
    if (recipe.error !== undefined) return <div>error! {recipe.error}</div>;

    let header = (
      <Header as="h1">
        {recipe.title}
        <Header.Subheader>
          {`From ${recipe.source}. Makes ${parseFloat(
            (recipe.quantity * this.state.scale).toFixed(1)
          )} ${recipe.unit}`}
        </Header.Subheader>
      </Header>
    );
    let rightSidebar = (
      <Segment>
        <Header as="h3">Scaling</Header>
        <p>
          <b>Approx weight:</b>&nbsp; TODO::CALC
        </p>
        <p>
          <b>Total minutes:</b>&nbsp;{recipe.totalMinutes}
        </p>
        <p>
          <b>Total minutes (calculated):</b>&nbsp;TODO::CALC
        </p>
        <p>
          <b>Equipment:</b>&nbsp;{recipe.equipment}
        </p>
        <Form>
          <Form.Field width={4}>
            <label>Multiplier</label>
            <input
              type="number"
              min="0"
              max="10"
              step=".1"
              value={this.state.scale}
              onChange={this.handleScaleChange}
            />
          </Form.Field>
        </Form>
        <Divider />
        <Header as="h3">Misc</Header>
        <Button
          as={Link}
          to={`/editor/${this.props.slug}`}
          content="edit recipe"
          icon="edit"
        />
        <Button
          content="show raw"
          onClick={this.handleToggleShowRaw.bind(this)}
        />
        <Divider />
        <Image src="http://via.placeholder.com/2000x1200" />
      </Segment>
    );
    return (
      <div>
        {header}
        <Grid>
          <Grid.Column width={12}>
            <RecipeTable recipe={recipe} scale={this.state.scale} />
            <Header as="h1">Notes</Header>
            {recipe.notes.map(eachNote => (
              <Card key={eachNote.id}>
                <Card.Content description={eachNote.body} />
                <Card.Content extra>
                  <Moment
                    tz="America/Los_Angeles"
                    format="ddd MMM Do YYYY, h:mm a"
                  >
                    {eachNote.created_at}
                  </Moment>
                </Card.Content>
              </Card>
            ))}

            <Header as="h1">Images</Header>
            {(recipe.images === null ? [] : recipe.images).map(eachImage => (
              <Card
                key={eachImage.id}
                image={eachImage.url}
                extra={
                  <Moment
                    tz="America/Los_Angeles"
                    format="ddd MMM Do YYYY, h:mm a"
                  >
                    {eachImage.created_at}
                  </Moment>
                }
              />
            ))}
            {this.state.showRaw ? (
              <div>
                <Header as="h1">raw JSON</Header>
                <Card fluid>
                  <pre>{JSON.stringify(recipe, true, 2)}</pre>
                </Card>
              </div>
            ) : null}
          </Grid.Column>
          <Grid.Column width={4}>{rightSidebar}</Grid.Column>
        </Grid>
      </div>
    );
  }
}
