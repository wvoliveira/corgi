import Router, {useRouter} from "next/router";
import React from "react";
import useSWR, { mutate } from "swr";

import storage from "../../lib/utils/storage";
import ListErrors from "../common/ListErrors";
import APILink from "../../lib/api/link";

const LinkForm = () => {
  const [isLoading, setLoading] = React.useState(false);
  const [errors, setErrors] = React.useState([]);
  const [fullURL, setFullURL] = React.useState("");
  const [shortURL, setShortURL] = React.useState("");
  const [response, setResponse] = React.useState(null);

  const handleURLChange = React.useCallback(
    (e) => setFullURL(e.target.value),
    []
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await APILink.create(fullURL);

      console.debug("STATUS: ", status);
      console.debug("DATA: ", data);

      if (status !== 201) {
        setErrors(data.errors);
      }

      if (data?.data) {
        setResponse(data.data);

        const link = window.location.protocol + "//" + data.data?.domain + "/" + data.data?.keyword;
        setShortURL(link);

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
      <br/>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="URL"
          value={fullURL}
          onChange={handleURLChange}
          required={true}
          disabled={isLoading}
        />

        { " " }

        <button
          type="submit"
          disabled={isLoading}
        >
          Create
        </button>
      </form>

      <ListErrors errors={errors} />

      {/*Protocol: {window.location.protocol}*/}
      {console.log("RESPONSE: ", response)}
      {response &&
      <p>
        Link: <a href={shortURL}>{shortURL}</a>
      </p>
      }
    </>
  );
};

export default LinkForm;
