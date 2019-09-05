/*jshint esversion: 6 */

import React, { Component } from "react";
import { NavLink } from "react-router-dom";
import { connect } from "react-redux";
import * as banner from "../actions/bannerActions";
import turnOffLightsAction from "../actions/turnOffLightsAction"

import "./Landing.css";

class Landing extends Component {
    componentDidMount() {
        this.props.banner();
    }
    render() {
        return (
            <div>
                <div className="landingPanel fd-panel">
                    <div className="landingBody">
                        <h1>Welcome to Bullseye!</h1>
                        <NavLink className="link" to="/questionnaire">
                            <button href="/questionnaire" className="fd-button" onClick={turnOffLightsAction}>
                                Start
                            </button>
                        </NavLink>
                    </div>
                </div>
            </div>
        );
    }
}

const mapDispach = dispatch => {
    return {
        banner: () => dispatch(banner.plain())
    };
};

export default connect(
    undefined,
    mapDispach
)(Landing);
