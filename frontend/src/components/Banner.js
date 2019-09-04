import React, { Component } from "react";
import { connect } from "react-redux";

import "./Banner.css";

class Banner extends Component {
    render() {
	return <div className="bannerCrop">
		 <img alt="banner" src="/bullseye.jpg"
		      className={"banner " + (this.props.blured ? "blured ":"") +
				 (this.props.visible ? "":"hidden ")} />
	       </div>;
    } 
}

const mapState = state => {
    return state.banner;
};

const mapDispach = dispatch => {
    return {
    };
};

export default connect(mapState, mapDispach)(Banner);
