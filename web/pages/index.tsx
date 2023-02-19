import LinkForm from "components/link/LinkForm";
import Head from "next/head";
import React from "react";

import Banner from "../components/home/Banner";
import MainView from "../components/home/MainView";
import Tags from "../components/home/Tags";

const Home = () => (
  <>
    <Head>
      <title>Corgi</title>
      <meta
        name="description"
        content="A shortener app."
      />
    </Head>
    <div className="home-page">
      {/* <Banner /> */}
      <div className="container page">
        <LinkForm />

        <div className="row">
          <MainView />
          <div className="col-md-3">
            <div className="sidebar">
              <p>Popular Tags</p>
              <Tags />
            </div>
          </div>
        </div>
      </div>
    </div>
  </>
);

export default Home;
