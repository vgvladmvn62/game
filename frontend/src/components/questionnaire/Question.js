import React from "react";
import Answer from "./Answer";

export default function Question(props) {
  return (
    <div>
      <div className="fd-panel__header">
        <div className="fd-panel__head">
          <h1 className="fd-panel__title">{props.text}</h1>
        </div>
      </div>
      <div className="fd-panel__body">
        <fieldset id={props.num} className="fd-form__set">
          {props.answers.map((a, i) => (
            <Answer
              key={"answer" + i}
              onChange={() => {
                props.setAnswer(i);
              }}
              text={a}
              question={props.num}
              answer={i}
              checked={i === props.selected}
            />
          ))}
        </fieldset>
      </div>
    </div>
  );
}
