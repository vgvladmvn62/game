import React from "react";

export default function AcceptButton(props) {
  return (
    <button
      onClick={props.onClick}
      className="fd-button--positive questionnaireAccept"
    >
      Accept
    </button>
  );
}
