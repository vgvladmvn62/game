import React, { Component } from "react";
import { connect } from "react-redux";
import "./ProductPanel.css"

class ProductPanel extends Component {
    render() {
        const productData = this.props.matched[this.props.product];
        const product = this.props.matched[this.props.product].product;
        return (
            <div className="fd-container fd-container--fluid fd-container--flex">
                <div className="fd-col">
                    <div className="fd-tile productTile">
                        <div className="fd-tile__content productContent">
                            <div className="fd-container fd-container--fluid">
                                <div className="fd-tile__title productTitle">

                                    <div className="productName">
                                        <span>{product.name}</span>
                                    </div>

                                    <div className="productPrice">
                                        <span>{product.price.formatted_value}</span>
                                    </div>

                                </div>
                            </div>
                            <div className="fd-product-tile__media">
                                <img src={product.image} className="responsive" alt="Product" />
                            </div>

                            <div className="productTags">
                                {productData.attributes.map((attr, index) => <span key={index} className={"icon-" + attr.found}></span>)}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

const mapState = state => {
    return { ...state.products, answers: state.questionnaire.answers }
};

export default connect(mapState)(ProductPanel)
