import Router from "next/router";
import React from "react";
import useSWR from "swr";

import ListErrors from "../../components/common/ListErrors";
import ArticleAPI from "../../lib/api/article";
import storage from "../../lib/utils/storage";
import editorReducer from "../../lib/utils/editorReducer";

const PublishArticleEditor = () => {
  const initialState = {
    title: "",
    description: "",
    body: "",
    tagList: [],
  };

  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState(null);
  const [posting, dispatch] = React.useReducer(editorReducer, initialState);
  const { data: currentUser } = useSWR("user", storage);

  const handleTitle = (e) =>
    dispatch({ type: "SET_TITLE", text: e.target.value });
  const handleDescription = (e) =>
    dispatch({ type: "SET_DESCRIPTION", text: e.target.value });
  const handleBody = (e) =>
    dispatch({ type: "SET_BODY", text: e.target.value });
  const addTag = (tag) => dispatch({ type: "ADD_TAG", tag: tag });
  const removeTag = (tag) => dispatch({ type: "REMOVE_TAG", tag: tag });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const { data, status } = await ArticleAPI.create(
      posting,
      currentUser?.token
    );

    setLoading(false);

    if (status !== 200) {
      setError(data?.message);
    }

    await Router.push("/");
  };

  return (
    <>
      <form>
        <p />
        <input
          type="text"
          placeholder="Article Title"
          value={posting.title}
          onChange={handleTitle}
        />
        <p />
        <input
          type="text"
          placeholder="What's this article about?"
          value={posting.description}
          onChange={handleDescription}
        />
        <p />
        <textarea
          rows={8}
          cols={21}
          placeholder="Write your article (in markdown)"
          value={posting.body}
          onChange={handleBody}
        />
        <p />
        <p />
        <button
          type="button"
          disabled={isLoading}
          onClick={handleSubmit}
        >
          Publish Article
        </button>
      </form>

      {error && <ListErrors error={error} />}
    </>
  );
};

export default PublishArticleEditor;
