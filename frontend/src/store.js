import { createStore, applyMiddleware, compose } from "redux";
import promise from "redux-promise-middleware";
import reducer from "./reducers";

const logger = store => next => action => {
  console.group(action.type);
  console.info("dispatching", action);
  let result = next(action);
  console.log("next state", store.getState());
  console.groupEnd();
  return result;
};

const middleware = applyMiddleware(promise(), logger);
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
export default createStore(reducer, composeEnhancers(middleware))
// export default createStore(reducer, middleware);
