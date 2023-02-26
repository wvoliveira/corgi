import Router from "next/router";
import React from "react";
import { mutate } from "swr";

import ListErrors from "../common/ListErrors";
import APIAuthPassword from "../../lib/api/authPassword";

const LoginForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState([]);
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
        setErrors(data.errors);
      }

      if (data?.data?.user && data?.data?.tokens) {
        window.localStorage.setItem("corgi.user", JSON.stringify(data.data.user));
        window.localStorage.setItem("corgi.tokens", JSON.stringify(data.data.tokens));

        mutate("corgi.user", data.data.user);
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
      <ListErrors errors={errors} />

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
    </>
  );
};

export default LoginForm;
