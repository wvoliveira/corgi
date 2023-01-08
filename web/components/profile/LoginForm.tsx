import Router from "next/router";
import React from "react";
import { mutate } from "swr";

// import ListErrors from "../common/ListErrors";
import UserAPI from "../../lib/api/user";
import Link from "next/link";

const LoginForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState([]);
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  const handleEmailChange = React.useCallback(
    (e: any) => setEmail(e.target.value),
    []
  );
  const handlePasswordChange = React.useCallback(
    (e: any) => setPassword(e.target.value),
    []
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    console.log("Email: ", email);
    console.log("Password: ", password);

    try {
      const { data, status } = await UserAPI.login(email, password);
      if (status !== 200) {
        setErrors(data);
      }

      if (data?.data) {
        window.localStorage.setItem("user", JSON.stringify(data.data));
        console.log(data);

        mutate("user", data?.data);
        Router.push("/");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {/* <ListErrors errors={errors} /> */}

      {/* <form onSubmit={handleSubmit}>
      </form> */}

      <form onSubmit={handleSubmit}>
        <p>
          <label htmlFor="email">Email</label><br/>
          <input type="email" onChange={handleEmailChange}/>
        </p>

        <p>
          <label htmlFor="password">Password</label><br/>
          <input type="password" onChange={handlePasswordChange}/>
        </p>

        <p><button>Login</button></p>
      </form>

      <p>Don&apos;t have an account? <Link href="/register">Register</Link></p>
      <p><Link href="/password/reset">Forgot password?</Link></p>

    </>
  );
};

export default LoginForm;