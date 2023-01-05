import Link from "next/link";
import React from "react";

export default function Navbar() {

  return (
    <>

      <p>
        <Link href="/">Home</Link>
        { ' ' } | { ' ' }
        <Link href="/search">Search</Link>
        { ' ' } | { ' ' }
        <Link href="/login">Login</Link>
        { ' ' } | { ' ' }
        <Link
         href="/register">Register</Link>
      </p>

    </>
  )
}
