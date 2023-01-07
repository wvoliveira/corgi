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
      return error;
    }
  },
  login: async (email: string, password: string) => {
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
      return error;
    }
  },
  register: async (username: string, email: string, password: string) => {
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
      return error;
    }
  },
  save: async (user: any) => {
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
      return error;
    }
  }
};

export default UserAPI;
