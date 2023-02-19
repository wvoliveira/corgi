import Router from "next/router";
import React from "react";
import { mutate } from "swr";

import ListErrors from "../common/ListErrors";
import APIAuthPassword from "../../lib/api/authPassword";

const LoginForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState([]);
  const [usernameEmail, setUsernameEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  const handleUsernameEmailChange = React.useCallback(
    (e) => setUsernameEmail(e.target.value),
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
      const { data, status } = await APIAuthPassword.login(usernameEmail, password);

      console.log(status);
      console.log(data);

      if (status !== 200) {
        setErrors(data.errors);
      }

      if (data?.data?.user && data?.data?.tokens) {
        window.localStorage.setItem("user", JSON.stringify(data.data.user));
        window.localStorage.setItem("tokens", JSON.stringify(data.data.tokens));

        mutate("user", data.data.user);
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
          placeholder="Username or email"
          value={usernameEmail}
          onChange={handleUsernameEmailChange}
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
