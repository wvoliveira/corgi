import Link from "next/link";
import Router from "next/router";
import React, { useEffect } from "react";

import useSWR, { useSWRConfig } from "swr";
import LinkUser from "../../lib/api/user";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";


export default function Navbar() { 
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [isLoggedIn, setIsLoggedIn] = React.useState(false);

  const { mutate } = useSWRConfig()
  const { data: currentUser } = useSWR("user", storage);

  useEffect(() => {
    if (currentUser == undefined) {
      return;
    }

    console.log("current user: ", currentUser);

    const isLoggedIn = checkLogin(currentUser);
    setIsLoggedIn(isLoggedIn);
    return;
  }, currentUser);

  const handleLogout = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await LinkUser.logout();
      console.log("Data: ", data);
      console.log("Status: ", status);

      if (status !== 200) {
        setError(data);
        return
      }

      if (status === 200) {
      }

      setTimeout(() => {
        localStorage.removeItem("user");
        mutate("user")

        setIsLoggedIn(false);
        setLoading(false);

        Router.push("/");
      }, 1000);

    } catch (error) {
      console.error(error);
    } finally {
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
