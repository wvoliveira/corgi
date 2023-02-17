import Head from 'next/head'
import Link from 'next/link'
import Image from 'next/image'
import { Inter } from '@next/font/google'
import styles from '../styles/Home.module.css'
import LinkForm from '../components/link/LinkCreateForm'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <>
      <LinkForm />
    </>
  )
}
