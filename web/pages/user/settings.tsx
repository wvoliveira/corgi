import { useRouter } from 'next/router'
import React, {useEffect} from "react";
// @ts-ignore
import useSWR, { mutate, trigger } from "swr";

import SettingsForm from "../../components/profile/SettingsForm";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";

const Settings = () => {
  const router = useRouter()
  const { data: currentUser } = useSWR("corgi.user", storage);
  const isLoggedIn = checkLogin(currentUser);

  useEffect(() => {
    if (!isLoggedIn) {
      router.push(`/`);
    }
  }, [])

  const handleLogout = async (e) => {
    e.preventDefault();
    window.localStorage.removeItem("corgi.user");
    window.localStorage.removeItem("corgi.tokens");
    await mutate("corgi.user", null);
    await mutate("corgi.tokens", null);
    router.push(`/`).then(() => trigger("user"));
  };

  return (
    <div>
      <div>
        <div>
          <div>
            <h1>Your Settings</h1>
            <SettingsForm />
            <hr />
            <button onClick={handleLogout}>
              Or click here to logout.
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

// Settings.getInitialProps = async ({ res }) => {
//   return {
//     res,
//   };
// };

export default Settings;
