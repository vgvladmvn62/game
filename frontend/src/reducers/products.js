const default_state = {
    loaded: false,
    loading: false,
    highlightedProduct: 0,
    matched: []
}

export default function reducer(state = default_state, action) {
    switch (action.type) {
        case "SUBMIT_QUESTIONNAIRE_REJECTED": {
            alert("Cannot submit questionnaire!");
            return state;
        }
        case "SUBMIT_QUESTIONNAIRE_FULFILLED": {
            return {
                ...state,
                loaded: true,
                matched: action.payload.data.matched
            };
        }
        case "SUBMIT_QUESTIONNAIRE_PENDING": {
            console.log("Waiting for products");
            return {
                ...state,
                loading: true
            }
        }
        case "RESET": {
            console.log("Resetting to main page");
            return default_state;
        }
        default: {
            return state;
        }
    }
}
