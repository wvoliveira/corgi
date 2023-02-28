import Router from "next/router";
import React from "react";
import { mutate } from "swr";

import ListErrors from "../common/ListErrors";
import APIAuthPassword from "../../lib/api/authPassword";

const LoginForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  const handleEmailChange = React.useCallback(
    (e) => setEmail(e.target.value),
    []
  );
  const handlePasswordChange = React.useCallback(
    (e) => setPassword(e.target.value),
    []
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await APIAuthPassword.login(email, password);

      console.log(status);
      console.log(data);

      if (status !== 200) {
        setError(data?.message);
      }

      if (data?.data?.user && data?.data?.tokens) {
        window.localStorage.setItem("corgi.user", JSON.stringify(data.data.user));
        window.localStorage.setItem("corgi.tokens", JSON.stringify(data.data.tokens));

        await mutate("corgi.user", data.data.user);
        await Router.push("/");
      }
    } catch (error) {
      console.error("Error: ", error);
      setError(error);
    } finally {
      setLoading(false);
    }
  };

  if (error) {
    setTimeout(() => {
      setError("");
    }, 3000)
  }

  return (
    <>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Email"
          value={email}
          onChange={handleEmailChange}
        />
        {" "}
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={handlePasswordChange}
        />
        {" "}
        <button
          type="submit"
          disabled={isLoading}
        >
          Sign in
        </button>
      </form>

      {error && <ListErrors error={error} />}
    </>
  );
};

export default LoginForm;
