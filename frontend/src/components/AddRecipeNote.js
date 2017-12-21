import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Button, Form, Segment, TextArea } from 'semantic-ui-react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { addRecipeNote } from '../actions/recipe';

class AddRecipeNote extends Component {
  constructor(props) {
    super(props);
    this.state = { value: '' };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }
  handleSubmit() {
    this.props.addRecipeNote(this.props.slug, this.state.value);
    //reset the text area
    this.setState({ value: '' });
  }
  render() {
    return (
      <Segment>
        <Form>
          <Form.Field
            control={TextArea}
            label="Note"
            placeholder="text"
            value={this.state.value}
            onChange={this.handleChange}
          />
          <Form.Field
            control={Button}
            content="Add Note"
            onClick={this.handleSubmit}
          />
        </Form>
      </Segment>
    );
  }
}
AddRecipeNote.propTypes = {
  slug: PropTypes.string.isRequired
};

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      addRecipeNote
    },
    dispatch
  );
};

export default connect(false, mapDispatchToProps)(AddRecipeNote);
