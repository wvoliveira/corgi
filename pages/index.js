import styles from "../styles/styles.module.css";

export default function Home() {
  return (
    <div className={styles.container}>

      <main className={styles.main}>
        <h1 className={styles.title}>ELGA</h1>
        <p className={styles.description}>Design & Technology</p>
        <p className={styles.description}>URL shortener.</p>
      </main>
    </div>
  );
}
