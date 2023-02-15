const storage = async key => {
  const value = localStorage.getItem(key);

  if (value) {
    return JSON.parse(value);
  }

  return "";
};

export default storage;