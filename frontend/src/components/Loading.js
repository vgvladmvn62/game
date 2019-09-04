import React from "react";

export default function(props) {
    return <div className={ "fd-spinner " + props.className } aria-hidden="false" aria-label="Loading">
        <div></div>
    </div>
}
