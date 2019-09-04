const default_state = {
    blured: false,
    visible: true
}
export default function reducer(state = default_state, action) {
    switch(action.type) {
    case "PLAIN":
	return {
	    blured: false,
	    visible: true
	}
    case "BLUR":
	return {
	    blured: true,
	    visible: true
	}
    case "HIDE":
	return {
	    blured: false,
	    visible: false
	}
    default:
	return state
    }
}
