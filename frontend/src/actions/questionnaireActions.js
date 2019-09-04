import axios from "axios";

export function fetchQuestions() {
    return {
        type: "FETCH_QUESTIONS",
        payload: axios.get(window["configuration"].backendURL + "/questions")
    };
}

export function selectAnswer(question, answer) {
    return {
        type: "SELECT_ANSWER",
        payload: {
            question: question,
            answer: answer
        }
    };
}

export function nextStep() {
    return {
        type: "NEXT_STEP"
    };
}

export function submit(answers, questions) {
    console.log(questions);
    console.log(answers);
    return {
        type: "SUBMIT_QUESTIONNAIRE",
        payload: axios.post(window["configuration"].backendURL + "/results", {
            answers: answers.map((a,q)=>{return questions[q].answers[a]} )
        })
    };
}

export function previousStep() {
    return {
        type: "PREVIOUS_STEP"
    };
}
