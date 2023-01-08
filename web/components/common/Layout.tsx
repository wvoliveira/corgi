import Head from "next/head";
import Link from "next/link";
import React from "react";
import Footer from "./Footer";
import Navbar from "./Navbar";

export default function Layout({ children }:{children:any}) {

  return (
    <>

      <Head>
        <title>Corgi</title>
        <meta name="description" content="Corgi | A shortener system" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Link href="/">
        <h1>Corgi</h1>
      </Link>

      <Navbar />
      { children }
      <Footer />
    </>
  )
}
