import axios from "axios";

const UserAPI = {
  current: async () => {
    const user: any = window.localStorage.getItem("user");
    const token = user?.token;
    try {
      const response = await axios.get(`/user`, {
        headers: {
          Authorization: `Token ${encodeURIComponent(token)}`,
        },
      });
      return response;
    } catch (error) {
      return error.response;
    }
  },
  login: async (email, password) => {
    try {
      const response = await axios.post(
        `/api/users/login`,
        JSON.stringify({ user: { email, password } }),
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
  register: async (username, email, password) => {
    try {
      const response = await axios.post(
        `/api/users`,
        JSON.stringify({ user: { username, email, password } }),
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
  save: async (user) => {
    try {
      const response = await axios.put(
        `/api/user`,
        JSON.stringify({ user }),
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
  follow: async (username) => {
    const user: any = JSON.parse(window.localStorage.getItem("user"));
    const token = user?.token;
    try {
      const response = await axios.post(
        `/api/profiles/${username}/follow`,
        {},
        {
          headers: {
            Authorization: `Token ${encodeURIComponent(token)}`,
          },
        }
      );
      return response;
    } catch (error) {
      return error.response;
    }
  },
  unfollow: async (username) => {
    const user: any = JSON.parse(window.localStorage.getItem("user"));
    const token = user?.token;
    try {
      const response = await axios.delete(
        `/api/profiles/${username}/follow`,
        {
          headers: {
            Authorization: `Token ${encodeURIComponent(token)}`,
          },
        }
      );
      return response;
    } catch (error) {
      return error.response;
    }
  },
  get: async (username) => axios.get(`/api/profiles/${username}`),
};

export default UserAPI;
