import React from "react";

export default function(props) {
  return (
    <div className="fd-form__item fd-form__item--check" onChange={props.onChange}>
      <input
        className="fd-form__control"
        type="radio"
        id={"radio-" + props.question + "-" + props.answer}
        value=""
        name={"radio-name-" + props.question}
        checked={props.checked}
      />
      <label
        className="fd-form__label"
        htmlFor={"radio-" + props.question + "-" + props.answer}
      >
        {props.text}
      </label>
    </div>
  );
}
