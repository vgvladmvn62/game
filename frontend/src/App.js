/*jshint esversion: 6 */

import React, { Component } from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import { Provider } from "react-redux";

import Questionnaire from "./components/questionnaire/Questionnaire";
import Products from "./components/products/Products";
import Landing from "./components/Landing";
import Banner from "./components/Banner";

import store from "./store";
import * as banner from "./actions/bannerActions";

class App extends Component {
    componentDidMount() {
	store.dispatch(banner.plain());
    }
    render() {
	const state = store.getState();
	console.log(state);
	return (
	    <Provider store={store}>
	      <div>
                <Banner />

		<div className="fd-ui fd-ui--fundamental">
		  <div className="fd-ui__header">
		    <nav className="fd-global-nav">
		      <div className="fd-global-nav__group fd-global-nav__group--left">
			<div className="fd-global-nav__logo" />

				  <a href={"/"}><div className="fd-global-nav__product-name">Bullseye v2</div></a>
		      </div>
		    </nav>
		  </div>

		  <div className="fd-ui__app">
		    <Router>
		      <div
			className="fd-container fd-container--fluid"
            id="main-container"
              >
			<Route exact path="/" component={Landing} />
			<Route path="/questionnaire" component={Questionnaire} />
			<Route path="/products" component={Products} />
		      </div>
		    </Router>
		  </div>
		</div>
	      </div>
	    </Provider>
	);
    }
}

export default App;
