import React from "react";
import LinkAPI from "../../lib/api/link";
import LinkCopy from "./LinkCopy";

export default function LinkForm() {
  const [isLoading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");
  const [urlFull, setURLFull] = React.useState("");
  const [urlShort, setURLShort] = React.useState("");

  const handleURLFullChange = React.useCallback(
    (e: any) => setURLFull(e.target.value),
    []
  );

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    let payload = {
      url: urlFull,
    }

    try {
      const { data, status } = await LinkAPI.create(payload);
      if (status !== 200) {
        setError(data.message);
      }

      if (data?.data) {
        let url = data.data.domain + "/" + data.data.keyword
        setURLShort(url)
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
        <p>Paste the URL to be shortened:

          { ' ' }

          <input 
            type="text" placeholder="https://"
            value={urlFull}
            onChange={handleURLFullChange}
            required={true}
          />

          { ' ' }

          <button 
            type="submit"
            disabled={isLoading}
          >
            Shorten URL
          </button>

        </p>
      </form>

      <LinkCopy url={urlShort} />

      {error ? 
        <>
          Error: {error}
        </>
      : null}

    </>
  )
}
