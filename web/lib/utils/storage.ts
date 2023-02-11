const storage = async key => {
  const value = localStorage.getItem(key);
  return !!value ? JSON.parse(value) : JSON.parse("");
};

export default storage;