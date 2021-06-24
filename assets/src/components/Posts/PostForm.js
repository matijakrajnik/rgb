import { useState, useContext, useEffect, useCallback } from 'react';

import AuthContext from '../../store/auth-context';
import Errors from '../Errors/Errors';

const PostForm = (props) => {
  const authContext = useContext(AuthContext);

  const [titleValue, setTitleValue] = useState('');
  const [contentValue, setContentValue] = useState('');

  const [errors, setErrors] = useState({});

  const populateFields = useCallback(() => {
    if (props.post) {
      setTitleValue(props.post.Title);
      setContentValue(props.post.Content);
    }
  }, [props.post]);

  useEffect(() => {
    populateFields();
  }, [populateFields]);

  async function submitHandler(event) {
    event.preventDefault();
    setErrors({});

    try {
      const method = props.onEditPost ? 'PUT' : 'POST';
      let body = {
        Title: titleValue,
        Content: contentValue,
      }
      if (props.onEditPost) {
        body.ID = props.post.ID;
      }
      const response = await fetch('api/posts',
        {
          method: method,
          body: JSON.stringify(body),
          headers: {
            'Content-Type': 'application/json',
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
        setTitleValue('');
        setContentValue('');
        if (props.onAddPost) {
          props.onAddPost(data.data);
        }
        if (props.onEditPost) {
          props.onEditPost(data.data);
        }
      }
    } catch (error) {
      setErrors({ "error": error.message });
    }
  };

  const titleChangeHandler = (event) => { setTitleValue(event.target.value) }
  const contentChangeHandler = (event) => { setContentValue(event.target.value) }

  const errorContent = Object.keys(errors).length === 0 ? null : Errors(errors);
  const submitButtonText = props.onEditPost ? 'Update Post' : 'Add Post';

  return (
    <section>
      <div className="container w-75 pb-4">
        <form onSubmit={submitHandler}>
          <div className="form-group pb-3">
            <label htmlFor="title">Title</label>
            <input id="title" type="text" className="form-control" required value={titleValue} onChange={titleChangeHandler}></input>
          </div>
          <div className="form-group pb-3">
            <label htmlFor="content">Content</label>
            <textarea id="content" className="form-control" rows="5" required value={contentValue} onChange={contentChangeHandler}></textarea>
          </div>
          <button type="submit" className="btn btn-success">{submitButtonText}</button>
        </form>
        {errorContent}
      </div>
    </section>
  );
}

export default PostForm;
