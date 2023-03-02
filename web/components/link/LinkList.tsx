import {useRouter} from "next/router";
import React, {useEffect} from "react";
import useSWR from "swr";

import ErrorMessage from "../common/ErrorMessage";
import LoadingSpinner from "../common/LoadingSpinner";
import Maybe from "../common/Maybe";
import Pagination from "../common/Pagination";
import {SERVER_BASE_URL} from "../../lib/utils/constant";
import fetcher from "../../lib/utils/fetcher";

const LinkList = () => {
    const [protocol, setProtocol] = React.useState("http");
    const router = useRouter();
    const {query} = router;

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

    let fetchURL = `${SERVER_BASE_URL}/links?page=${query.page}&offset=${query.offset}`;
    console.debug("fetchURL: ", fetchURL);

    const {data: content, error} = useSWR(fetchURL, fetcher);
    console.error(error);

    if (error) {
        return (
            <div>
                <ErrorMessage message="Cannot load recent links..."/>
            </div>
        );
    }

    if (!content) return <LoadingSpinner/>;

    console.debug("Data: ", content)
    console.debug("Error: ", error)

    const {data} = content;

    if (data.links && data.links.length === 0) {
        return <div>No links are here... yet.</div>;
    }

    // @ts-ignore
    return (
        <>
            <table>
                <tr>
                    <th className="table-td-id">ID</th>
                    <th>Full URL</th>
                    <th>Short URL</th>
                    <th>Clicks</th>
                </tr>
            {data.links?.map((link, index) => {
                const shortURL = `${protocol}//${link.domain}/${link.keyword}`
                return (
                    <>
                        <tr key={link.id}>
                            <td className="table-td-id">{link.id}</td>
                            <td className="table-td-full-url">{link.url}</td>
                             <td><a
                                 key={link.id}
                                 target="_blank"
                                 href={shortURL}
                                 rel={shortURL}>{shortURL}
                             </a></td>
                            <td>{link?.clicks?.total}</td>
                        </tr>
                    </>
                )
            })}
            </table>

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
}
    ;

    export default LinkList;
