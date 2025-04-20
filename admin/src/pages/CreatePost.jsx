import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import postService from '../services/postService';
import PostForm from '../components/PostForm';

function CreatePost() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleSubmit = async (postData) => {
    setLoading(true);
    setError(null);
    
    try {
      await postService.createPost(postData);
      navigate('/');
    } catch (err) {
      console.error('Error creating post:', err);
      setError(err.response?.data?.message || 'Failed to create post. Please try again.');
      setLoading(false);
    }
  };

  return (
    <section className="section">
      <div className="container">
        <h1 className="title">Create New Post</h1>

        {error && (
          <div className="notification is-danger">
            <button className="delete" onClick={() => setError(null)}></button>
            {error}
          </div>
        )}

        <PostForm 
          onSubmit={handleSubmit}
          isLoading={loading}
          initialData={{
            type: 'POST',
            title: '',
            content: '',
            tags: [],
            isPublished: false,
            slug: '',
            excerpt: ''
          }}
        />
      </div>
    </section>
  );
}

export default CreatePost;