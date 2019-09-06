import React, { Component } from "react";
import { connect } from "react-redux";
import Loading from "../Loading";
import * as banner from "../../actions/bannerActions";
import ProductPanel from "./ProductPanel";
import { reset } from "../../actions/productsActions";
import turnOffLightsAction from "../../actions/turnOffLightsAction";


class Products extends Component {
    componentDidMount() {
        this.props.banner();
    }
    loadingOrProducts() {
        if (this.props.loaded) {
            return this.props.matched.map((_,i)=><ProductPanel key={i} product={i}/>);
        } else {
            return <Loading />;
        }
    }
    reset() {
        this.props.reset();
        this.props.history.push("/")
    }
    render() {
        return <div className="App productsPanel fd-panel">
            <div className="reset-button-container">
            <button
                onClick={ () => {
                    turnOffLightsAction()
                    this.reset()
                } }
                className="fd-button--negative back"
            >
                Go back
            </button>
            </div>
            {this.loadingOrProducts()}
        </div>
    }
}

const mapState = state => {
    return {
	...state.products
    };
};

const mapDispatch = dispatch => {
    return {
        banner: ()=> dispatch(banner.hide()),
        reset: () => {dispatch(reset())}
    };
};

export default connect(mapState, mapDispatch)(Products);
