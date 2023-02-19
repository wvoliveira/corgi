import Router from "next/router";
import React from "react";
import useSWR, { mutate } from "swr";

import storage from "../../lib/utils/storage";
import ListErrors from "../common/ListErrors";
import APILink from "../../lib/api/link";

const LinkForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState([]);
  const [url, setURL] = React.useState("");
  const [response, setResponse] = React.useState({});

  const { data: tokens } = useSWR("tokens", storage);

  const handleURLChange = React.useCallback(
    (e) => setURL(e.target.value),
    []
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await APILink.create(url, tokens.access_token);

      console.log(status);
      console.log(data);

      if (status !== 201) {
        setErrors(data.errors);
      }

      if (data?.data) {
        setResponse(data.data)
        console.log(response)
        // mutate("user", data.data.user);
        // Router.push("/");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <ListErrors errors={errors} />

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="URL"
          value={url}
          onChange={handleURLChange}
        />

        { " " }

        <button
          type="submit"
          disabled={isLoading}
        >
          Create
        </button>
      </form>
    </>
  );
};

export default LinkForm;
