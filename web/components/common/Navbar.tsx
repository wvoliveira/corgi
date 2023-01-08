import Link from "next/link";
import React from "react";

import useSWR, { useSWRConfig } from "swr";
import LinkUser from "../../lib/api/user";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";


export default function Navbar() { 
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");

  const { mutate } = useSWRConfig()
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);

  console.log("current user: ", currentUser);
  console.log("is logged in: ", isLoggedIn);

  const handleLogout = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await LinkUser.logout();
      console.log("Data: ", data);
      console.log("Status: ", status);

      if (status !== 200) {
        setError(data);
      }

    } catch (error) {
      console.error(error);
    } finally {

      setTimeout(() => {
        window.localStorage.removeItem("user");
        mutate("user")
        setLoading(false);
      }, 1000);
    }
  };

  return (
    <>

      <p>
        <Link href="/">Home</Link>
        { ' ' } | { ' ' }
        <Link href="/search">Search</Link>
        { ' ' } | { ' ' }

        {isLoggedIn ? 
          <>
            <Link href="/profile">Profile</Link>
            { ' ' } | { ' ' }
            {isLoading ? <>Logout...</> : <Link href="/logout" onClick={handleLogout}>Logout</Link>}
          </> 
        : 
          <>
            <Link href="/login">Login</Link>
            { ' ' } | { ' ' }
            <Link
            href="/register">Register</Link>
          </>
        }
      </p>

    </>
  )
}
