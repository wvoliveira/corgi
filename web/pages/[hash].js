import redirect from "nextjs-redirect";

import styles from "../styles/styles.module.css";
import { urls } from "../constants/urls.js";

// This function gets called at build time
export async function getStaticPaths() {
  // Get urls from urls data
  const paths = urls.map((url) => ({
    params: { id: url.id, hash: url.hash },
  }));

  // We'll pre-render only these paths at build time.
  // { fallback: false } means other routes should 404.
  return { paths, fallback: false };
}

// This also gets called at build time
export async function getStaticProps({ params }) {
  const hash = searchHash(params.hash);

  // Pass post data to the page via props
  return { props: { hash } };
}

function searchHash(hash) {
  const filtered = urls.filter((p) => p.hash === hash);

  var payload = {
    id: "",
    url: "",
    status: 200,
    message: "successful",
  };

  if (filtered.length > 0) {
    const data = filtered[0];
    payload["id"] = data.id;
    payload["url"] = data.url;
  } else {
    (payload["status"] = 404),
      (payload["message"] = `url with id ${hash} not found.`);
  }
  return payload;
}

function URL({ hash }) {
  const Redirect = redirect(hash.url);

  return (
    <Redirect>
      <div className={styles.container}>
        <main className={styles.main}>
          <h1 className={styles.title}>ELGA</h1>
          <p className={styles.description}>Redirect to {hash.url}...</p>
        </main>
      </div>
    </Redirect>
  );
}

export default URL;
