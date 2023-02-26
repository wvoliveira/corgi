import axios from "axios";

const updateOptions = () => {
  if (typeof window === "undefined") return {};
  if (!window.localStorage.getItem("corgi.tokens")) return {};
  if (Object.keys(window.localStorage.getItem("corgi.tokens")).length === 0) return {};

  const tokens = JSON.parse(window.localStorage.getItem("corgi.tokens"));

  if (!!tokens.access_token) {
    return {
      headers: {
        Authorization: `Bearer ${tokens.access_token}`,
      },
    };
  }
};

export default async function (url) {
  const { data } = await axios.get(url, updateOptions());
  return data;
}
