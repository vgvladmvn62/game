import React, { Component } from "react";
import "./Questionnaire.css";
import QuestionPanel from "./QuestionPanel";
import Loading from "../Loading";
import { connect } from "react-redux";
import * as actions from "../../actions/questionnaireActions";
import * as banner from "../../actions/bannerActions";

class Questionnaire extends Component {
    componentDidMount() {
        this.props.banner();
        console.log(this.props);
        this.props.fetchData();
    }

    questionOrLoading() {
        if (this.props.loaded) {
            return <QuestionPanel history={this.props.history} num={this.props.step} />;
        } else {
            return <Loading />;
        }
    }

    render() {
        return (
            <div className="App fd-panel">
                <div className="questionnaireRoot fd-panel">
                    {this.questionOrLoading()}
                </div>
            </div>
        );
    }
}

const mapState = state => {
    return state.questionnaire;
};

const mapDispatch = dispatch => {
    return {
        fetchData: () => dispatch(actions.fetchQuestions()),
        banner: () => dispatch(banner.blur())
    };
};

export default connect(
    mapState,
    mapDispatch
)(Questionnaire);
