import Router from "next/router";
import React from "react";
import { mutate } from "swr";

// import ListErrors from "../common/ListErrors";
import UserAPI from "../../lib/api/user";
import Link from "next/link";

const LoginForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [message, setMessage] = React.useState("");
  const [error, setError] = React.useState("");
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

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    setLoading(true);
    setError("");
    setMessage("");

    try {
      const { data, status } = await UserAPI.login(email, password);

      if (status == 401) {
        setError("Email or password invalid!");
        console.log("Message: ", data?.message);
        return
      }

      if (status == 200 && data?.data) {
        window.localStorage.setItem("user", JSON.stringify(data.data));

        setMessage("Authenticated! Reloading...")

        // TODO: set user to storage when auth and reload the page.
        setTimeout(() => {
          mutate("user", data.data);
          setLoading(false);
          Router.push("/");
        }, 1000);
      }
    } catch (error) {
      console.error("Error: ", error);
    } finally {
      console.log("Finally handle.")
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
          <input type="email" autoComplete="true" onChange={handleEmailChange} readOnly={isLoading}/>
        </p>

        <p>
          <label htmlFor="password">Password</label><br/>
          <input type="password" autoComplete="true" onChange={handlePasswordChange} readOnly={isLoading}/>
        </p>

        <p><button disabled={isLoading}>Login</button></p>
      </form>

      {error != "" ?
        <>
          Error: {error}
        </>
      : null}

      {message != "" ?
        <>
          {message}
        </>
      : null}

      <p>Don&apos;t have an account? <Link href="/register">Register</Link></p>
      <p><Link href="/login/reset">Forgot password?</Link></p>

    </>
  );
};

export default LoginForm;