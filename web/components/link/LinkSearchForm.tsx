import React, { useEffect } from "react";
import { useRouter } from "next/router";

import LinkAPI from "../../lib/api/link";
import LinkList from "./LinkList";

export default function LinkSearchForm() {
  const router = useRouter();

  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [searchText, setSearchText] = React.useState("");
  const [links, setLinks] = React.useState([]);

  const handleURLFullChange = React.useCallback(
    (e: any) => setSearchText(e.target.value),
    []
  );

  const getLinks = async () => {
    let text = router.query.q?.toString();

    setLoading(true);

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

    setLoading(true);
    setLinks([]);
    setError("");

    router.replace({
      query: { ...router.query, q: searchText },
    });

    getLinks();

  };

  useEffect(() => {
    let paramQ = router.query.q?.toString()
    setSearchText(paramQ ? paramQ : "");
    getLinks();
  }, [router.isReady])

  if (!router.isReady) {
    return <>Loading...</>
  }

  return (
    <>
      <form onSubmit={handleSubmit} method="get">
          <input
            type="text" placeholder="Type a text to search."
            value={searchText}
            onChange={handleURLFullChange}
            required={false}
          />

          { ' ' }

          <button 
            type="submit"
            disabled={isLoading}
          >
            Search
          </button>

      </form>

      {error != "" ?
        <>
          Error: {error}
        </>
      : null}

      <LinkList links={links} />

    </>
  )
}
