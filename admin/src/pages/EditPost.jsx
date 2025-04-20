import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import postService from '../services/postService';
import PostForm from '../components/PostForm';

function EditPost() {
  const { id } = useParams();
  const navigate = useNavigate();
  
  const [post, setPost] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);

  useEffect(() => {
    fetchPost();
  }, [id]);

  const fetchPost = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await postService.getPostById(id);
      setPost(response.data);
    } catch (err) {
      console.error(`Error fetching post with ID ${id}:`, err);
      setError('Failed to load post. Please try again later.');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (postData) => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);
    
    try {
      await postService.updatePost(id, postData);
      setSuccessMessage('Post updated successfully');
      fetchPost(); // Refresh post data
    } catch (err) {
      console.error(`Error updating post with ID ${id}:`, err);
      setError(err.response?.data?.message || 'Failed to update post. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handlePublish = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);
    
    try {
      await postService.publishPost(id);
      setSuccessMessage('Post published successfully');
      fetchPost(); // Refresh post data
    } catch (err) {
      console.error(`Error publishing post with ID ${id}:`, err);
      setError(err.response?.data?.message || 'Failed to publish post. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleUnpublish = async () => {
    setSaving(true);
    setError(null);
    setSuccessMessage(null);
    
    try {
      await postService.unpublishPost(id);
      setSuccessMessage('Post unpublished successfully');
      fetchPost(); // Refresh post data
    } catch (err) {
      console.error(`Error unpublishing post with ID ${id}:`, err);
      setError(err.response?.data?.message || 'Failed to unpublish post. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
      setSaving(true);
      setError(null);
      
      try {
        await postService.deletePost(id);
        navigate('/');
      } catch (err) {
        console.error(`Error deleting post with ID ${id}:`, err);
        setError(err.response?.data?.message || 'Failed to delete post. Please try again.');
        setSaving(false);
      }
    }
  };

  if (loading) {
    return (
      <section className="section">
        <div className="container has-text-centered">
          <button className="button is-loading is-large is-white"></button>
          <p className="mt-3">Loading post...</p>
        </div>
      </section>
    );
  }

  if (error && !post) {
    return (
      <section className="section">
        <div className="container">
          <div className="notification is-danger">
            <p>{error}</p>
          </div>
          <Link to="/" className="button is-link">
            Back to Posts
          </Link>
        </div>
      </section>
    );
  }

  return (
    <section className="section">
      <div className="container">
        <div className="level">
          <div className="level-left">
            <div className="level-item">
              <h1 className="title">Edit Post</h1>
            </div>
          </div>
          <div className="level-right">
            <div className="level-item">
              <Link to="/" className="button">
                <span className="icon">
                  <i className="fas fa-arrow-left"></i>
                </span>
                <span>Back to Posts</span>
              </Link>
            </div>
          </div>
        </div>

        {error && (
          <div className="notification is-danger">
            <button className="delete" onClick={() => setError(null)}></button>
            {error}
          </div>
        )}

        {successMessage && (
          <div className="notification is-success">
            <button className="delete" onClick={() => setSuccessMessage(null)}></button>
            {successMessage}
          </div>
        )}

        {post && (
          <>
            <div className="box">
              <div className="level">
                <div className="level-left">
                  <div className="level-item">
                    <div className="tags has-addons">
                      <span className="tag is-medium">Status</span>
                      <span className={`tag is-medium ${post.isPublished ? 'is-success' : 'is-warning'}`}>
                        {post.isPublished ? 'Published' : 'Draft'}
                      </span>
                    </div>
                  </div>
                </div>
                <div className="level-right">
                  <div className="level-item">
                    <div className="buttons">
                      {post.isPublished ? (
                        <button 
                          className={`button is-warning ${saving ? 'is-loading' : ''}`}
                          onClick={handleUnpublish}
                          disabled={saving}
                        >
                          <span className="icon">
                            <i className="fas fa-eye-slash"></i>
                          </span>
                          <span>Unpublish</span>
                        </button>
                      ) : (
                        <button 
                          className={`button is-success ${saving ? 'is-loading' : ''}`}
                          onClick={handlePublish}
                          disabled={saving}
                        >
                          <span className="icon">
                            <i className="fas fa-eye"></i>
                          </span>
                          <span>Publish</span>
                        </button>
                      )}
                      <button 
                        className={`button is-danger ${saving ? 'is-loading' : ''}`}
                        onClick={handleDelete}
                        disabled={saving}
                      >
                        <span className="icon">
                          <i className="fas fa-trash-alt"></i>
                        </span>
                        <span>Delete</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <PostForm 
              onSubmit={handleSubmit}
              isLoading={saving}
              initialData={post}
              isEdit={true}
            />
          </>
        )}
      </div>
    </section>
  );
}

export default EditPost; 