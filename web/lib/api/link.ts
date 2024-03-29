import axios from "axios";

import {SERVER_BASE_URL} from "../utils/constant";
import {getQuery} from "../utils/getQuery";

const APILink = {
    all: (page, limit = 10) =>
        axios.get(`${SERVER_BASE_URL}/articles?${getQuery(limit, page)}`),

    byAuthor: (author, page = 0, limit = 5) =>
        axios.get(
            `${SERVER_BASE_URL}/articles?author=${encodeURIComponent(
                author
            )}&${getQuery(limit, page)}`
        ),

    byTag: (tag, page = 0, limit = 10) =>
        axios.get(
            `${SERVER_BASE_URL}/articles?tag=${encodeURIComponent(tag)}&${getQuery(
                limit,
                page
            )}`
        ),

    delete: (id, token) =>
        axios.delete(`${SERVER_BASE_URL}/articles/${id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }),

    favorite: (slug) =>
        axios.post(`${SERVER_BASE_URL}/articles/${slug}/favorite`),

    favoritedBy: (author, page) =>
        axios.get(
            `${SERVER_BASE_URL}/articles?favorited=${encodeURIComponent(
                author
            )}&${getQuery(10, page)}`
        ),

    feed: (page, limit = 10) =>
        axios.get(`${SERVER_BASE_URL}/articles/feed?${getQuery(limit, page)}`),

    get: (slug) => axios.get(`${SERVER_BASE_URL}/articles/${slug}`),

    unfavorite: (slug) =>
        axios.delete(`${SERVER_BASE_URL}/articles/${slug}/favorite`),

    update: async (article, token) => {
        const {data, status} = await axios.put(
            `${SERVER_BASE_URL}/articles/${article.slug}`,
            JSON.stringify({article}),
            {
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${encodeURIComponent(token)}`,
                },
            }
        );
        return {
            data,
            status,
        };
    },

    create: async (url) => {
        const tokens: any = JSON.parse(window.localStorage.getItem("corgi.tokens"));
        const access_token = tokens?.access_token;

        const {data, status} = await axios.post(
            `${SERVER_BASE_URL}/links`,
            JSON.stringify({url: url}),
            {
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${encodeURIComponent(access_token)}`,
                },
            }
        );
        return {
            data,
            status,
        };
    },
};

export default APILink;
