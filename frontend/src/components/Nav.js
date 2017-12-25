import React from 'react';
import { NavLink } from 'react-router-dom';
import { Container, Menu } from 'semantic-ui-react';

const Nav = () => (
  <Menu fixed="top" inverted>
    <Container>
      <Menu.Item as={NavLink} to="/" exact header>
        Nicky's Recipes
      </Menu.Item>
      <Menu.Item as={NavLink} to="about">
        About
      </Menu.Item>
      <Menu.Item as={NavLink} to="images">
        Images
      </Menu.Item>
    </Container>
  </Menu>
);
export default Nav;
