import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button, Form, Header } from 'semantic-ui-react';
import { Redirect } from 'react-router-dom';
import { loginFromJWT } from '../actions/users';

class Auth extends Component {
  constructor(props) {
    super(props);
  }
  componentDidMount() {
    this.props.loginFromJWT(this.props.match.params.jwt);
  }
  render() {
    if (this.props.user.authenticated) return <Redirect to="/" />;
    return (
      <div className="container">
        <Header dividing content="Loading..." />
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    user: state.user
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      loginFromJWT
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(Auth);
