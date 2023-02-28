const storage = async key => {
  const value = localStorage.getItem(key);
  const data = !!value ? JSON.parse(value) : undefined;
  console.debug("Data: ", data);
  return data;
};

export default storage;
