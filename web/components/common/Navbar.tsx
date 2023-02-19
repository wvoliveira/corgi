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
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);

  const handleClick = React.useCallback(() => setPage(0), []);

  return (
    <>
      <a href="/" onClick={handleClick}>
        Corgi
      </a>
      { " · " }
      <a href="/" onClick={handleClick}>
        Home
      </a>
      { " · " }
      <Maybe test={isLoggedIn}>
        <a href="/editor/new">
          New Post
        </a>
        { " · " }
        <a href="/user/settings">
          Settings
        </a>
        { " · " }
        <a href={`/profile/${currentUser?.username}`} onClick={handleClick}>
          {currentUser?.name}
        </a>
      </Maybe>
      <Maybe test={!isLoggedIn}>
        <a href="/user/login">
          Sign in
        </a>
        { " · " }
        <a href="/user/register">
          Sign up
        </a>
      </Maybe>
    </>
  );
};

export default Navbar;
