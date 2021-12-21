import Link from 'next/link'
import Layout from '../components/Layout'

const Index = () => (
  <Layout title="Corgi | Home">
    <h1>Corgi 🐕</h1>
    <p>
      <Link href="/about">
        <a>About</a>
      </Link>
    </p>
  </Layout>
)

export default Index
