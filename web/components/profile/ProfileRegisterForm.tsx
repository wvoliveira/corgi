import Link from "next/link";
import Router from "next/router";
import React from "react";
import UserAPI from "../../lib/api/user";


const ProfileRegisterForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [name, setName] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [message, setMessage] = React.useState("");

  const handleNameChange = React.useCallback(
    (e: any) => setName(e.target.value),
    []
  );

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

    console.log("Name: ", name);
    console.log("Email: ", email);
    console.log("Password: ", password);

    try {
      const { data, status } = await UserAPI.register(name, email, password);
      console.log("Data: ", data);
      console.log("Status: ", status);

      if (status !== 200) {
        setError(data);
      }

      if (status == 200) {
        setMessage("OK. Redirecting to login page...");

        setTimeout(() => {
          Router.push("/login")
        }, 2000);

      }

    } catch (error) {
      var msg = error?.response?.data?.message;
      if (msg.length > 0) {
        msg = msg.charAt(0).toUpperCase() + msg.slice(1) + ".";
      }

      console.log("Error to register: ", error);
      console.log("Error from API:", msg);

      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit}>
        <p>
          <label htmlFor="name">Name</label><br/>
          <input type="text" onChange={handleNameChange} disabled={isLoading}/>
        </p>

        <p>
          <label htmlFor="email">Email</label><br/>
          <input type="email" onChange={handleEmailChange} disabled={isLoading}/>
        </p>

        <p>
          <label htmlFor="password">Password</label><br/>
          <input type="password" onChange={handlePasswordChange} disabled={isLoading}/>
        </p>

        <p>
          <button disabled={isLoading}>Create an account</button>
          { ' ' }
          {isLoading ? "Loading..." : null}
          {message ? message : null}
          {error ? error : null}
        </p>
      </form>

      <p>Already a user? <Link href="/login">Log In</Link></p>

    </>
  );
};

export default ProfileRegisterForm;