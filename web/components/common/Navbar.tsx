import Link from "next/link";

import useSWR from "swr";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";


export default function Navbar() { 
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);

  console.log("current user: ", currentUser);
  console.log("is logged in: ", isLoggedIn);

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
