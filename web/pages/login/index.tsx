import Link from "next/link"
import Head from "next/head"

export default function Login() {
  return (
    <>
      <Head>
        <title>Corgi | Login</title>
        <meta name="description" content="Corgi | A shortener system" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Link href="/">
          <h1>Corgi</h1>
      </Link>

      <form>
        <p>
          <label htmlFor="email">Email</label><br/>
          <input type="email"/>
        </p>

        <p>
          <label htmlFor="password">Password</label><br/>
          <input type="password"/>
        </p>

        <p><button>Login</button></p>
      </form>

      <p>Don't have an account? <Link href="/register">Register</Link></p>
      <p><Link href="/password/reset">Forgot password?</Link></p>

    </>
  )
}