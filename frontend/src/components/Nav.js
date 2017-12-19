import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import { Container, Menu } from 'semantic-ui-react';

export default class Nav extends Component {
  constructor(props) {
    super(props);

    this.toggle = this.toggle.bind(this);
    this.state = {
      isOpen: false
    };
  }
  toggle() {
    this.setState({
      isOpen: !this.state.isOpen
    });
  }
  render() {
    return (
      <Menu fixed="top" inverted>
        <Container>
          <Menu.Item as={NavLink} to="/" exact header>
            Nicky's Recipes
          </Menu.Item>
          <Menu.Item as={NavLink} to="about">
            About
          </Menu.Item>
        </Container>
      </Menu>
    );
  }
}
