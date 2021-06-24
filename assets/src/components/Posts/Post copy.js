import { useState, useContext } from 'react';

import AuthContext from '../../store/auth-context';
import Errors from '../Errors/Errors';

const Post = (props) => {
  const [errors, setErrors] = useState({});

  const authContext = useContext(AuthContext);

  async function deleteHandler(event) {
    try {
      const response = await fetch('api/posts/' + props.ID,
        {
          method: 'DELETE',
          headers: {
            'Authorization': 'Bearer ' + authContext.token,
          },
        }
      );
      const data = await response.json();
      if (!response.ok) {
        let errorText = 'Failed to add new post.';
        if (!data.hasOwnProperty('error')) {
          throw new Error(errorText);
        }
        if ((typeof data['error'] === 'string')) {
          setErrors({ 'unknown': data['error'] })
        } else {
          setErrors(data['error']);
        }
      } else {
        props.onDeletePost(props.ID);
      }
    } catch (error) {
      setErrors({ "error": error.message });
    }
  };

  const errorContent = Object.keys(errors).length === 0 ? null : Errors(errors);

  return (
    <div className="card mb-5 pb-2">
      <div className="card-header">{props.Title}</div>
      <div className="card-body">{props.Content}</div>
      <div className="container">
        <button type="button" className="btn btn-danger float-right mx-3" onClick={deleteHandler}>Delete</button>
      </div>
      {errorContent}
    </div>
  );
};

export default Post;
