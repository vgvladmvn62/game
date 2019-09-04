const default_state = {
    step: 0,
    loaded: false,
    questions: [],
    answers: [],
    filled: false
}

export default function reducer(
    state = default_state,
    action
) {
    switch (action.type) {
        case "FETCH_QUESTIONS_PENDING": {
            console.log("Waiting for questions");
            return state;
        }
        case "FETCH_QUESTIONS_REJECTED": {
            console.log("Request failed, ", action.payload);
            alert("Server not responding");
            return state;
        }
        case "FETCH_QUESTIONS_FULFILLED": {
            return {
                ...state,
                loaded: true,
                questions: action.payload.data.questions,
                answers: Array(action.payload.data.questions.length).fill(-1)
            };
        }
        case "SELECT_ANSWER": {
            const newState = {
                ...state,
                answers: Object.assign([...state.answers], {
                    [action.payload.question]: action.payload.answer
                })
            };
            newState.filled = newState.answers.every(a => a >= 0);

            return newState;
        }

        case "NEXT_STEP": {
            return {
                ...state,
                step: state.step + (state.step < state.questions.length - 1)
            };
        }
        case "PREVIOUS_STEP": {
            return {
                ...state,
                step: state.step - (state.step > 0)
            };
        }
        case "RESET": {
            return default_state;
        }

        default: {
            return state;
        }
    }
}
