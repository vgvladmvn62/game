import axios from "axios";

export default () => {
    axios.get(window["configuration"].backendURL + "/turnofflights")
}