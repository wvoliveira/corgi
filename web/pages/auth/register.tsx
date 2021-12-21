import Link from 'next/link'
import Layout from '../../components/Layout'

const RegisterPage = () => (
    <Layout title="Corgi | Register">
        <h1>Register</h1>
        <p>
            <Link href="/about">
                <a>About</a>
            </Link>
        </p>
    </Layout>
)

export default RegisterPage
