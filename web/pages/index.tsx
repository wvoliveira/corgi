import LinkForm from "components/link/LinkForm";
import Head from "next/head";
import React from "react";

import Banner from "../components/home/Banner";
import MainView from "../components/home/MainView";
import Tags from "../components/home/Tags";
import LinkList from "../components/link/LinkList";

const Home = () => (
  <>
    <Head>
      <title>Corgi</title>
      <meta
        name="description"
        content="A URL shortener app."
      />
    </Head>
    <div>
      <div>
        <LinkForm />
        <LinkList />

        <div>
          <MainView />
          {/* <div>
            <div>
              <p>Popular Tags</p>
              <Tags />
            </div>
          </div> */}
        </div>
      </div>
    </div>
  </>
);

export default Home;
