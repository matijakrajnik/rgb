import { useState, useContext } from 'react';

import AuthContext from '../../store/auth-context';
import Errors from '../Errors/Errors';
import PostForm from "./PostForm";

const Post = (props) => {
  const [editing, setEditing] = useState(false);
  const [errors, setErrors] = useState({});

  const authContext = useContext(AuthContext);

  const switchModeHandler = () => {
    setEditing((prevState) => !prevState);
    setErrors({});
  };

  async function deleteHandler() {
    try {
      const response = await fetch('api/posts/' + props.post.ID,
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
        props.onDeletePost(props.post.ID);
      }
    } catch (error) {
      setErrors({ "error": error.message });
    }
  };

  const editPostHandler = () => {
    setEditing(false);
    props.onEditPost();
  }

  const cardTitle = editing ? 'Edit post' : props.post.Title;
  const cardBody = editing ? <PostForm post={props.post} onEditPost={editPostHandler} editing={true}/> : props.post.Content;
  const switchModeButtonText = editing ? 'Cancel' : 'Edit';
  const cardButtons = editing ?
    <div className="container">
      <button type="button" className="btn btn-link" onClick={switchModeHandler}>{switchModeButtonText}</button>
      <button type="button" className="btn btn-danger float-right mx-3" onClick={deleteHandler}>Delete</button>
    </div>
    :
    <div className="container">
      <button type="button" className="btn btn-link" onClick={switchModeHandler}>{switchModeButtonText}</button>
      <button type="button" className="btn btn-danger float-right mx-3" onClick={deleteHandler}>Delete</button>
    </div>
  const errorContent = Object.keys(errors).length === 0 ? null : Errors(errors);

  return (
    <div className="card mb-5 pb-2">
      <div className="card-header">{cardTitle}</div>
      <div className="card-body">{cardBody}</div>
      {cardButtons}
      {errorContent}
    </div>
  );
};

export default Post;
