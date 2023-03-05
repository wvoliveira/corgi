import Router from "next/router";
import React from "react";
import { mutate } from "swr";

import ListErrors from "../common/ListErrors";
import APIAuthPassword from "../../lib/api/authPassword";

const RegisterForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState([]);
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
      const { data, status } = await APIAuthPassword.register(
        email,
        password
      );
      if (status !== 200 && data?.errors) {
        setError(data.errors);
      }
      if (data?.user) {
        window.localStorage.setItem("corgi.user", JSON.stringify(data.user));
        mutate("user", data.user);
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
      <ListErrors error={error} />

      <form onSubmit={handleSubmit}>
        <input
          type="email"
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
          Sign up
        </button>
      </form>
    </>
  );
};

export default RegisterForm;
