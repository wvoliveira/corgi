import axios from "axios";

import { SERVER_BASE_URL } from "../utils/constant";

const APIAuthPassword = {
  login: async (email, password) => {
    try {
      const response = await axios.post(
        `${SERVER_BASE_URL}/auth/password/login`,
        JSON.stringify({ email: email, password: password }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      return response;
    } catch (error) {
      return error.response;
    }
  },
  register: async (name, email, password) => {
    try {
      const response = await axios.post(
        `${SERVER_BASE_URL}/auth/password/register`,
        JSON.stringify({ name: name, email: email, password: password }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      return response;
    } catch (error) {
      return error.response;
    }
  },
};

export default APIAuthPassword;
