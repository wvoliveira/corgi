import axios from "axios";

const updateOptions = () => {
  if (typeof window === "undefined") return {};

  if (!window.localStorage.user) return {};

  if (Object.keys(window.localStorage.user).length === 0) return {};

  const user = JSON.parse(window.localStorage.user);

  if (!!user.access_token) {
    return {
      headers: {
        Authorization: `Token ${user.access_token}`,
      },
    };
  }
};

export default async function (url) {
  const { data } = await axios.get(url, updateOptions());
  return data;
}
