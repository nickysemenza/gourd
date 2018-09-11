import React from 'react';
import { NavLink } from 'react-router-dom';
import { Container, Menu } from 'semantic-ui-react';
import { connect } from 'react-redux';
import { API_BASE_URL } from '../config';

const Nav = ({ user }) => (
  <Menu fixed="top" inverted>
    <Container>
      <Menu.Item as={NavLink} to="/" exact header>
        Nicky's Recipes
      </Menu.Item>
      <Menu.Item as={NavLink} to="/about">
        About
      </Menu.Item>
      <Menu.Item as={NavLink} to="/images">
        Images
      </Menu.Item>
      <Menu.Item as={NavLink} to="/meals">
        Meals
      </Menu.Item>

      <Menu.Menu position="right">
        <Menu.Item as={NavLink} to="/settings">
          Settings
        </Menu.Item>
        {user.authenticated ? (
          <Menu.Item name="signup">{user.me.first_name}</Menu.Item>
        ) : (
          <Menu.Item as="a" href={`${API_BASE_URL}/auth/facebook/login`}>
            LOGIN
          </Menu.Item>
        )}
      </Menu.Menu>
    </Container>
  </Menu>
);

function mapStateToProps(state) {
  return { user: state.user };
}
export default connect(mapStateToProps, {})(Nav);
