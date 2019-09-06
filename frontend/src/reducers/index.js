import { combineReducers } from "redux";

import products from "./products";
import questionnaire from "./questionnaire";
import banner from "./banner";

export default combineReducers({
    products,
    questionnaire,
    banner
});
