import React, { Component } from 'react';
import { NavLink } from 'react-router-dom'
import {NavbarToggler, Collapse } from 'reactstrap';
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
    render () {
        return (
            <nav className="navbar navbar-toggleable-md navbar-light bg-faded">
                <a className="navbar-brand" href="#">Recipes <NavbarToggler right onClick={this.toggle} /></a>


                <Collapse isOpen={this.state.isOpen} navbar>
                    <ul className="navbar-nav mr-auto">
                        <li className="nav-item">
                            <NavLink exact to="/" className="nav-link" activeClassName="active">Home</NavLink>
                        </li>
                        <li className="nav-item">
                            {/*<a className="nav-link active" href="#">Link</a>*/}
                            <NavLink to="/about" className="nav-link" activeClassName="active">About</NavLink>

                        </li>
                        <li className="nav-item">
                            <NavLink to="/chocolate" className="nav-link" activeClassName="active">Chocolate</NavLink>
                        </li>
                    </ul>
                    <form className="form-inline my-2 my-lg-0">
                        <input className="form-control mr-sm-2" type="text" placeholder="Search"/>
                            <button className="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
                    </form>
                </Collapse>
            </nav>
        );
    }
}
