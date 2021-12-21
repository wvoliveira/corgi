import Link from 'next/link'
import Layout from '../../components/Layout'
import {login} from "../../services/login_service";
import {useState} from "react";

export type LoginInputs = {
    email: string
    password: string
}

function Login() {
    // these values are hardcoded since our main.go api only accepts this auth combo
    const initialValues: LoginInputs = {email: "", password: ""};

    const [inputs, setInputs] = useState(initialValues);
    const [error, setError] = useState("");

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        const res = await login(inputs);
        if (res) setError(res);
    };

    const handleInputChange = (e: React.ChangeEvent<any>) => {
        e.persist();
        setInputs({
            ...inputs,
            [e.target.name]: e.target.value,
        });
    };

    return <>
        <Layout title="Corgi | Login">
            <h1>Login</h1>
            <p>
                <Link href="http://localhost:8081/api/auth/google/login">Google</Link>{' '}
                <Link href="http://localhost:8081/api/auth/facebook/login">Facebook</Link>
            </p>

            <form onSubmit={handleSubmit}>
                <div>
                    <label>Email: </label>
                    <input type="email" id="email" name="email" onChange={handleInputChange} value={inputs.email}
                           placeholder="user@example.com"/>
                </div>
                <div>
                    <label>Password: </label>
                    <input type="password" id="password" name="password" onChange={handleInputChange}
                           value={inputs.password} placeholder="********"/>
                </div>

                {error ? <p>Error: {error}</p> : null}

                <button type="submit">Login</button>
            </form>

        </Layout>
    </>;
}

export default Login