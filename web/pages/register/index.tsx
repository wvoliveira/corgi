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

      <form>
        <p>
          <label htmlFor="name">Name</label><br/>
          <input type="text"/>
        </p>

        <p>
          <label htmlFor="email">Email</label><br/>
          <input type="email"/>
        </p>

        <p>
          <label htmlFor="password">Password</label><br/>
          <input type="password"/>
        </p>

        <p><button>Create an account</button></p>
      </form>

      <p>Already a user? <Link href="/login">Log In</Link></p>

    </>
  )
}