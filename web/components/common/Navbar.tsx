import React from "react";
import useSWR from "swr";

import CustomLink from "./CustomLink";
import Maybe from "./Maybe";
import NavLink from "./NavLink";
import { usePageDispatch } from "../../lib/context/PageContext";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";

const Navbar = () => {
  const setPage = usePageDispatch();
  const { data: currentUser } = useSWR("corgi.user", storage);
  const isLoggedIn = checkLogin(currentUser);

  const handleClick = React.useCallback(() => setPage(0), []);

  return (
    <>
      Corgi
      {" · "}
      <CustomLink href="/" as="/">
        <span onClick={handleClick}>Home</span>
      </CustomLink>
      {" · "}
      <Maybe test={isLoggedIn}>
        <CustomLink href="/editor/new" as="/editor/new">
          New Link
        </CustomLink>
        {" · "}
        <CustomLink href="/user/settings" as="/user/settings">
          Settings
        </CustomLink>
        {" · "}
        <CustomLink 
          href={`/profile/${currentUser?.username}`}
          as={`/profile/${currentUser?.username}`}
        >
          <span onClick={handleClick}>{currentUser?.name}</span>
        </CustomLink>
      </Maybe>
      <Maybe test={!isLoggedIn}>
        <CustomLink href="/user/login" as="/user/login">
          Login
        </CustomLink>
        {" · "}
        <CustomLink href="/user/register" as="/user/register">
          Register
        </CustomLink>
      </Maybe>
    </>
  );
};

export default Navbar;
