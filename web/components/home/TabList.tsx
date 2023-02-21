import { useRouter } from "next/router";
import React from "react";
import useSWR from "swr";

import CustomLink from "../common/CustomLink";
import Maybe from "../common/Maybe";
import NavLink from "../common/NavLink";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";

const TabList = () => {
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);
  const router = useRouter();
  const {
    query: { tag },
  } = router;

  if (!isLoggedIn) {
    return (
      <div>
        <li>
          <a href="/">
            Global Feed
          </a>
        </li>

        <Maybe test={!!tag}>
          <li>
            <a href={`/?tag=${tag}`}>
              {tag}
            </a>
          </li>
        </Maybe>
      </div>
    );
  }

  return (
    <ul>
      <li>
        <a href={`/?follow=${currentUser?.username}`}>
          Your Feed
        </a>
      </li>

      <li>
        <a href="/">
          Global Feed
        </a>
      </li>

      <Maybe test={!!tag}>
        <li >
          <a
            href={`/?tag=${tag}`}
          >
           {tag}
          </a>
        </li>
      </Maybe>
    </ul>
  );
};

export default TabList;
