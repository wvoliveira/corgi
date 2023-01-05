import Link from "next/link"
import Head from "next/head"

export default function Login() {
  return (
    <>
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