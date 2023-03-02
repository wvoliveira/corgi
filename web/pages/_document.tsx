import Document, { Html, Head, Main, NextScript } from 'next/document'
import React from "react";

class MyDocument extends Document {

  render() {
    return (
      <Html lang="en">
        <Head>
          <meta httpEquiv="Content-Type" content="text/html; charset=utf-8" />
          <meta httpEquiv="X-UA-Compatible" content="IE=edge" />
          <meta name="robots" content="index, follow" />
          <meta key="googlebot" name="googlebot" content="index,follow" />
          <meta name="google" content="notranslate" />
          <meta name="mobile-web-app-capable" content="yes" />
          <meta name="apple-mobile-web-app-capable" content="yes" />
          <meta name="keywords" content="corgi, shortener, elga" />
          <meta property="og:locale" content="en_US" />
          <meta property="og:site_name" content="corgi" />
          <meta property="og:title" content="Corgi. A URL shortener app" />
          <meta
            property="og:description"
            content="Corgi. A simple and fast URL shortener app"
          />
          <meta property="og:url" content="https://corgi.app/" />
          <meta
            property="og:image"
            content="https://corgi.app/images/share-link.png"
          />
          <meta property="twitter:card" content="corgi-shortener" />
          <meta
            property="twitter:url"
            content="https://corgi.app/"
          />
          <meta
            property="twitter:title"
            content="Corgi. A URL shortener app"
          />
          <meta
            property="twitter:description"
            content="Corgi. A URL shortener app"
          />
          <meta
            property="twitter:image"
            content="https://corgi.app/images/talk-link.jpg"
          />
          <meta name="msapplication-TileColor" content="#000" />
          <meta
            name="msapplication-TileImage"
            content="/images/ms-icon-144x144.png"
          />
          <meta name="theme-color" content="#000" />
          <link
            rel="apple-touch-icon"
            sizes="57x57"
            href="/images/apple-icon-57x57.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="60x60"
            href="/images/apple-icon-60x60.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="72x72"
            href="/images/apple-icon-72x72.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="76x76"
            href="/images/apple-icon-76x76.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="114x114"
            href="/images/apple-icon-114x114.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="120x120"
            href="/images/apple-icon-120x120.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="144x144"
            href="/images/apple-icon-144x144.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="152x152"
            href="/images/apple-icon-152x152.png"
          />
          <link
            rel="apple-touch-icon"
            sizes="180x180"
            href="/images/apple-icon-180x180.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="192x192"
            href="/images/android-icon-192x192.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="32x32"
            href="/images/favicon-32x32.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="96x96"
            href="/images/favicon-96x96.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="16x16"
            href="/images/favicon-16x16.png"
          />
          <script
            type="application/ld+json"
            dangerouslySetInnerHTML={{
              __html: `{
              "@context": "https://schema.org/",
              "@type": "Organization",
              "url": "https://corgi.app/",
              "logo": "https://corgi.app/images/share-link.png",
              "sameAs": [
                "https://elga.io",
              ],
            }`,
            }}
          />
          <link rel="manifest" href="/manifest.json" />
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default MyDocument;
