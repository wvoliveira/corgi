const checkLogin = (currentUser: any) => {

  if (currentUser?.constructor === Object
    && Object.keys(currentUser).length !== 0) {
    return true;
  }

  return false;
};

export default checkLogin;