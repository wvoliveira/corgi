import React from "react";
import UserAPI from "../../lib/api/user";


const LoginResetPassword = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState([]);
  const [email, setEmail] = React.useState("");
  const [message, setMessage] = React.useState("");

  const handleEmailChange = React.useCallback(
    (e: any) => setEmail(e.target.value),
    []
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    console.log("Email: ", email);

    try {
      const { data, status } = await UserAPI.reset(email);
      console.log(data);

      if (status !== 200) {
        setError(data);
      }

      if (status == 200) {
        setMessage("Some text from API response");
      }

    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>

      <p>Type your e-mail bellow to reset password.</p>

      <form onSubmit={handleSubmit}>
        <p>
          <label htmlFor="email">Email</label><br/>
          <input type="email" onChange={handleEmailChange} placeholder="user@email.com"/>
        </p>

        <p><button>Login</button></p>
      </form>

    </>
  );
};

export default LoginResetPassword;