import React from "react";
import LinkAPI from "../../lib/api/link";
import LinkList from "./LinkList";

export default function LinkSearchForm() {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [urlShort, setURLShort] = React.useState("");
  const [links, setLinks] = React.useState([]);

  const handleURLFullChange = React.useCallback(
    (e: any) => setURLShort(e.target.value),
    []
  );

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    try {
      const { data, status } = await LinkAPI.list(urlShort);
      if (status !== 200) {
        setError(data.message);
      }

      if (data?.data) {
        console.log(data.data);
        setLinks(data.data);
      }

    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit}>
          <input 
            type="text" placeholder="Type a domain or keyword to search"
            value={urlShort}
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

      {error ? 
        <>
          Error: {error}
        </>
      : null}

      <LinkList links={links} />

    </>
  )
}
