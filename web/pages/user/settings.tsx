import Router from "next/router";
import React from "react";
import useSWR, { mutate, trigger } from "swr";

import SettingsForm from "../../components/profile/SettingsForm";
import checkLogin from "../../lib/utils/checkLogin";
import storage from "../../lib/utils/storage";

const Settings = ({ res }) => {
  const { data: currentUser } = useSWR("user", storage);
  const isLoggedIn = checkLogin(currentUser);

  if (!isLoggedIn) {
    if (res) {
      res.writeHead(302, {
        Location: "/",
      });
      res.end();
    }
    Router.push(`/`);
  }

  const handleLogout = async (e) => {
    e.preventDefault();
    window.localStorage.removeItem("user");
    mutate("user", null);
    Router.push(`/`).then(() => trigger("user"));
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

Settings.getInitialProps = async ({ res }) => {
  return {
    res,
  };
};

export default Settings;
