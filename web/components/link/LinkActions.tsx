import Router, { useRouter } from "next/router";
import React from "react";
// @ts-ignore
import useSWR, { trigger } from "swr";

import CustomLink from "../common/CustomLink";
import checkLogin from "../../lib/utils/checkLogin";
import ArticleAPI from "../../lib/api/article";
import { SERVER_BASE_URL } from "../../lib/utils/constant";
import storage from "../../lib/utils/storage";
import Maybe from "../common/Maybe";

const LinkActions = ({ article }) => {
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);
  const router = useRouter();
  const {
    query: { pid },
  } = router;

  const handleDelete = async () => {
    if (!isLoggedIn) return;

    const result = window.confirm("Do you really want to delete it?");

    if (!result) return;

    await ArticleAPI.delete(pid, currentUser?.token);
    trigger(`${SERVER_BASE_URL}/articles/${pid}`);
    Router.push(`/`);
  };

  const canModify =
    isLoggedIn && currentUser?.username === article?.author?.username;

  return (
    <Maybe test={canModify}>
      <span>
        <a
          href="/editor/[pid]"
        >
          <i /> Edit Article
        </a>

        <button
          onClick={handleDelete}
        >
          <i /> Delete Article
        </button>
      </span>
    </Maybe>
  );
};

export default LinkActions;
