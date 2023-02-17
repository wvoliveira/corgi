import React from "react";

import { APP_NAME } from "../../lib/utils/constant";

const Banner = () => (
  <div className="banner">
    <div className="container">
      <h1 className="logo-font">{APP_NAME.toLowerCase()}</h1>
      <p>A shortener app.</p>
    </div>
  </div>
);

export default Banner;
