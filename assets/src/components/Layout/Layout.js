import { Fragment } from 'react';

import './Layout.css';

import NavigationBar from './NavigationBar';

const Layout = (props) => {
  return (
    <Fragment>
      <NavigationBar/>
      <main>
        <div className="container">{props.children}</div>
      </main>
    </Fragment>
  );
};

export default Layout;
