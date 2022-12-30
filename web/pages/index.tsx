import Head from 'next/head'
import Link from 'next/link'
import Image from 'next/image'
import { Inter } from '@next/font/google'
import styles from '../styles/Home.module.css'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <>
      <Head>
        <title>Corgi</title>
        <meta name="description" content="Corgi | A shortener system" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <Link href="/">
          <h1>Corgi</h1>
        </Link>

        <p>
          <Link href="/login">Login</Link>
          { ' ' } | { ' ' }
          <Link href="/register">Register</Link>
        </p>

        <p>Create a anonymous short link:
          { ' ' } <input type="text" placeholder="https://" />
          { ' ' } <button>Create</button>
        </p>

      </main>
    </>
  )
}
