import axios from "axios";

type Link = {
  id: string;
	createdAt: Date;
	updatedAt: Date;

	domain:  string;
	keyword: string;
	url:     string;
	title:   string;
	active:  string;
};

const LinkAPI = {
  create: async (payload: any) => {
    try {
      const response = await axios.post(
        `/api/links`,
        JSON.stringify(payload),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      return response;
    } catch (error: any) {
      return error.response;
    }
  },
  list: async (url = "") => {
    var uri = `/api/links`;
    if (url != "") {
      uri = `${uri}?u=${url}`
    }

    try {
      const response = await axios.get(uri, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response;
    } catch (error: any) {
      return error.response;
    }
  },
  find: async (id: string) => {
    try {
      const response = await axios.get(`/api/links/${id}`, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response;
    } catch (error: any) {
      return error.response;
    }
  },
  save: async (id: string, link: any) => {
    try {
      const response = await axios.patch(
        `/api/links/${id}`,
        JSON.stringify({ link }),
        {
          headers: {
            "Content-Type": "application/json",
          },
      });
      return response;
    } catch (error: any) {
      return error.response;
    }
  },
  delete: async (id: string) => {
    try {
      const response = await axios.delete(`/api/links/${id}`, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response;
    } catch (error: any) {
      return error.response;
    }
  },
};

export default LinkAPI;
