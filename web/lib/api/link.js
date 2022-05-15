import axios from "axios";

import { SERVER_BASE_API_URL } from "../utils/constant";

const APILink = {
    create: async (payload) => {
        console.log("HERE");
        console.log(payload);
        try {
            const response = await axios.post(
                `${SERVER_BASE_API_URL}/links`,
                JSON.stringify(payload)
            );
            return response;
        } catch (error) {
            return error.response;
        }
    },
    delete: async (slug, commentId) => {
        try {
            const response = await axios.delete(
                `${SERVER_BASE_API_URL}/articles/${slug}/comments/${commentId}`
            );
            return response;
        } catch (error) {
            return error.response;
        }
    },

    forArticle: (slug) =>
        axios.get(`${SERVER_BASE_API_URL}/articles/${slug}/comments`),
};

export default APILink;
