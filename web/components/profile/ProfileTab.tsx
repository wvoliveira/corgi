import React from "react";

import NavLink from "../common/NavLink";
import { usePageDispatch } from "../../lib/context/PageContext";

const ProfileTab = ({ profile }) => {
  const setPage = usePageDispatch();
  return (
    <ul className="nav nav-pills outline-active">
      <li className="nav-item">
        <NavLink
          href="/profile/[pid]"
          as={`/profile/${encodeURIComponent(profile?.username)}`}
        >
          <span onClick={() => setPage(0)}>My links</span>
        </NavLink>
      </li>
      <li className="nav-item">
        <NavLink
          href="/profile/[pid]?favorite=true"
          as={`/profile/${encodeURIComponent(profile?.username)}?favorite=true`}
        >
          <span onClick={() => setPage(0)}>Favorite links</span>
        </NavLink>
      </li>
    </ul>
  );
};

export default ProfileTab;
