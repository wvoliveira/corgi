import React, { useEffect } from "react";
import useSWR from "swr";

import logger from "../../lib/logger/logger";
import storage from "../../lib/utils/storage";


const SettingsForm = () => {
  // logger.info("Profile/Setting page");
  const { data, error } = useSWR("user", storage);

  if (error) return <div>Failed to load</div>
  if (!data) return <div>Loading...</div>

  console.log(data);

  return (
    <>
      <p>Profile / Settings</p>
    </>
  );
};

export default SettingsForm;
