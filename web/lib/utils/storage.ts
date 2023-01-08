const storage = async key => {
  const value = localStorage.getItem(key);

  console.log("Key: ", key);
  console.log("Value: ", value);

  return !!value ? JSON.parse(value) : undefined;
};

export default storage;