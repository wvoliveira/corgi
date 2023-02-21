import React from "react";

import CustomLink from "../common/CustomLink";
import LoadingSpinner from "../common/LoadingSpinner";
import { usePageDispatch } from "../../lib/context/PageContext";
import useSWR from "swr";
import { SERVER_BASE_URL } from "../../lib/utils/constant";
import fetcher from "../../lib/utils/fetcher";
import ErrorMessage from "../common/ErrorMessage";

const Tags = () => {
  const setPage = usePageDispatch();
  const handleClick = React.useCallback(() => setPage(0), []);
  const { data, error } = useSWR(`${SERVER_BASE_URL}/tags`, fetcher);

  if (error) return <ErrorMessage message="Cannot load popular tags..." />;
  if (!data) return <LoadingSpinner />;

  const { tags } = data;
  return (
    <div>
      {tags?.map((tag) => (
        <a
          key={tag}
          href={`/?tag=${tag}`}
        >
          <span onClick={handleClick}>{tag}</span>
        </a>
      ))}
    </div>
  );
};

export default Tags;
