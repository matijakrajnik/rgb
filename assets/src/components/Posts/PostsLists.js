import Post from "./Post";

const PostsList = (props) => {
  return (
    <ul>
      {props.posts.map((post) => (
        <Post
          onEditPost={props.onEditPost}
          onDeletePost={props.onDeletePost}
          key={post.ID}
          post={post} />
      ))}
    </ul>
  );
};

export default PostsList;
