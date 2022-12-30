import Link from "next/link"
import Head from "next/head"

export default function Register() {
  return (
    <>
      <Head>
        <title>Corgi | Register</title>
        <meta name="description" content="Corgi | A shortener system" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Link href="/">
          <h1>Corgi</h1>
      </Link>

      <p>
        <Link href="/login">Login</Link>
      </p>

      <p><input type="email" placeholder="user@email.com"/></p>
      <p><input type="password" placeholder="password"/></p>
      <p><input type="password" placeholder="password"/></p>
      <p><button>Register</button></p>

    </>
  )
}