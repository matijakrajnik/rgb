import React, { useState } from 'react';

const AuthContext = React.createContext({
  token: null,
  loggedIn: false,
  login: (token) => {},
  logout: () => {},
});

export const AuthContextProvider = (props) => {
  const tokenKey = 'rgbToken';

  const [token, setToken] = useState(localStorage.getItem(tokenKey));

  const loggedIn = !!token;

  const loginHandler = (token) => {
    setToken(token);
    localStorage.setItem(tokenKey, token);
  };

  const logoutHandler = () => {
    setToken(null);
    localStorage.removeItem(tokenKey);
  };

  const contextValue = {
    token: token,
    loggedIn: loggedIn,
    login: loginHandler,
    logout: logoutHandler,
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {props.children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
