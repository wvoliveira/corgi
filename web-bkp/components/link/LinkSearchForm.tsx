import React, { useEffect } from "react";
import { useRouter } from "next/router";

import LinkAPI from "../../lib/api/link";
import LinkList from "./LinkList";

export default function LinkSearchForm() {
  const router = useRouter();

  const [isLoading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");
  const [searchText, setSearchText] = React.useState("");
  const [links, setLinks] = React.useState([]);

  const handleSearchTextChange = React.useCallback(
    (e: any) => setSearchText(e.target.value),
    []
  );

  const getLinks = async (text="") => {
    setLoading(true);

    text = text ?? searchText

    try {
      const { data, status, statusText } = await LinkAPI.FindAll(text);

      if (status == 500) {
        console.log("Error 500");
        console.log("Data: ", data);
        console.log("statusText: ", statusText);
        setError(statusText);
      }

      if (status !== 200 && status !== 500) {
        setError(data.message);
      }

      if (data?.data) {
        setLinks(data.data);
      }

    } catch (error: any) {
      setError(error);
      console.log("Error data: ", error.data);
      console.log("Error status: ", error.status);
      console.log("Error headers: ", error.headers);

    } finally {
      setLoading(false);
    }
  }

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    router.replace({
      query: { ...router.query, q: searchText },
    });
  };

  useEffect(() => {
    if (!router.isReady || !router.query) {
      return
    }
    const { q } = router.query;

    let text = q?.toString()
    setSearchText(text);
    getLinks(text);
  }, [router])

  if (isLoading) {
    return <>Loading...</>
  }

  return (
    <>
      <form onSubmit={handleSubmit} method="get">
          <input
            type="text" placeholder="Type a text to search."
            value={searchText}
            onChange={handleSearchTextChange}
            required={false}
          />

          { ' ' }

          <button type="submit">Search</button>
      </form>

      {error ?
        <>
          Error: {error}
        </>
      : null}

      <LinkList links={links} />
    </>
  )
}
