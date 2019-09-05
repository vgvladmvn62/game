import React from "react";
import { connect } from "react-redux";
import Question from "./Question";
import AcceptButton from "./AcceptButton";
import * as actions from "../../actions/questionnaireActions";

const mapState = state => {
  return state.questionnaire;
};

const mapDispach = dispatch => {
  return {
    selectAnswer: (q, a) => dispatch(actions.selectAnswer(q, a)),
    nextStep: () => dispatch(actions.nextStep()),
    previousStep: () => dispatch(actions.previousStep()),
    submit: (a, q) => dispatch(actions.submit(a, q))
  };
};

export default connect(
  mapState,
  mapDispach
)(function(props) {
  const q = props.questions[props.num];
  return (
    <div>
      <Question
        num={props.num}
        text={q.text}
        answers={q.answers}
        selected={
          props.answers[props.num] === null ? -1 : props.answers[props.num]
        }
        setAnswer={selected => {
            props.selectAnswer(props.num, selected);
            props.nextStep();
        }}
      />
      <div className="fd-panel__footer questionnaireFooter">
        <div className="questionnaireNavigation fd-button-group" role="group">
          <button
            className="fd-button--grouped questionnaireNavBtn"
            onClick={() => {
              props.previousStep();
            }}
            aria-disabled={props.step === 0}
          >
            Previous
          </button>
          <button
            className="fd-button--grouped questionnaireNavBtn"
            onClick={() => {
              props.nextStep();
            }}
            aria-disabled={props.step === props.questions.length - 1}
          >
            Next
          </button>
        </div>
        {props.filled ? (
          <AcceptButton
            onClick={() => {
              props.submit(props.answers, props.questions)
              props.history.push("/products");
            }}
          />
        ) : (
          ""
        )}
      </div>
    </div>
  );
});
