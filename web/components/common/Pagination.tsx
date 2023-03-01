import React from "react";
// @ts-ignore
import { trigger } from "swr";

import { getRange, getPageInfo } from "../../lib/utils/calculatePagination";
import { usePageDispatch, usePageState } from "../../lib/context/PageContext";
import Maybe from "./Maybe";
import Link from "next/link";

interface PaginationProps {
  total: number;
  limit: number;
  pageCount: number;
  currentPage: number;
  lastIndex: number;
  fetchURL: string;
}

const Pagination = ({
    page,
    pages,
    limit,
    total,
}) => {
    const previousPage = page === 1 && 1 || page-1
    const nextPage = page === pages && page || page+1
  return (
    <div>
        <br/>
        Page: {page} {" "}
        Pages: {pages} {" "}
        Limit: {limit} {" "}
        Total: {total} {" "}
        <Link href={`?page=${previousPage}`}>
          {`<`}
        </Link>
        {" | "}
        <Link href={`?page=${nextPage}`}>
          {`>`}
        </Link>
    </div>
  );
};

export default Pagination;
