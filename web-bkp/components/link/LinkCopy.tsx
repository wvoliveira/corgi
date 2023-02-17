import React, { SyntheticEvent } from "react";

const LinkCopy = ({url}: any) => {
  const [buttonText, setButtonText] = React.useState("Copy URL");

  const handleCopy = async (e: any) => {
    e.preventDefault();
    
    await navigator.clipboard.writeText(url);
    setButtonText("Copied!")

    setTimeout(() => {
      setButtonText("Copy URL")
    }, 2000);
  };

  return (
    <>
      {url ? 
      <p>Your shortened URL: { ' ' }

        <input
          type="text"
          value={url}
          disabled={true}
        /> { ' ' }

        <button
          onClick={handleCopy}
        >
          {buttonText}
        </button>

      </p>
      : null}

    </>
  )
}

export default LinkCopy;
