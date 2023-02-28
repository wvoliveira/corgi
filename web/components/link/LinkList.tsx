import { useRouter } from "next/router";
import React, {useEffect} from "react";
import useSWR from "swr";

import ErrorMessage from "../common/ErrorMessage";
import LoadingSpinner from "../common/LoadingSpinner";
import Maybe from "../common/Maybe";
import Pagination from "../common/Pagination";
import { SERVER_BASE_URL, DEFAULT_LIMIT } from "../../lib/utils/constant";
import fetcher from "../../lib/utils/fetcher";

const LinkList = () => {
  const [protocol, setProtocol] = React.useState("http");
  const router = useRouter();
  const { asPath, pathname, query } = router;

  useEffect(() => {
    setProtocol(window.location.protocol);
  });

  if (query.page === undefined) {
    // @ts-ignore
    query.page = 1
  }

  if (query.offset == undefined) {
    // @ts-ignore
    query.offset = 0
  }

  console.debug("Protocol: ", protocol);

  let fetchURL = `${SERVER_BASE_URL}/links?page=${query.page}&offset=${query.offset}`;
  console.debug("fetchURL: ", fetchURL);

  const { data: content, error } = useSWR(fetchURL, fetcher);
  console.error(error);

  if (error) {
    return (
      <div>
        <ErrorMessage message="Cannot load recent links..." />
      </div>
    );
  }

  if (!content) return <LoadingSpinner />;

  console.debug("Data: ", content)
  console.debug("Error: ", error)

  const { data } = content;

  if (data.links && data.links.length === 0) {
    return <div>No links are here... yet.</div>;
  }

  // @ts-ignore
  return (
    <>
      {data.links?.map((link, index) => {
        const shortURL = `${protocol}//${link.domain}/${link.keyword}`
        return (
          <p key={link.id} title={link.id}>
            {link.id.substring(0, 5)} | {" "}
            <a key={link.id} target="_blank" href={shortURL} rel={shortURL}>{shortURL.substring(0, 15)}</a>... | {" "}
            {link.url}
          </p>
        )
      })}

      <Maybe test={data.total && data.total > 10}>
        <Pagination
          page={data.page}
          pages={data.pages}
          limit={data.limit}
          total={data.total}
        />
      </Maybe>
    </>
  );
};

export default LinkList;
