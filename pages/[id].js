import { useRouter } from "next/router";
import useSWR from "swr";

import styles from "../styles/Home.module.css";

const fetcher = async (url) => {
  const res = await fetch(url);
  const data = await res.json();

  if (res.status !== 200) {
    throw new Error(data.message);
  }
  return data;
};

export default function URL() {
  const router = useRouter();
  const { query } = useRouter();

  if (query.id == "api") {
    return (
      <div className={styles.container}>
        <main className={styles.main}>
          <h1 className={styles.title}>ELGA</h1>
          <p className={styles.description}>Redirect "API".</p>
        </main>
      </div>
    );
  }

  const { data, error } = useSWR(() => query.id && `/api/${query.id}`, fetcher);

  if (error) return <div>{error.message}</div>;
  if (!data) return <div>Loading...</div>;

  // return <h1>{data.url}</h1>;
  router.push(data.url, '', '')
}
