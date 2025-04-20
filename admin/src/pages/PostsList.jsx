import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import postService from '../services/postService';

function PostsList() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('all');
  const [pagination, setPagination] = useState({
    current_page: 1,
    per_page: 10,
    total_pages: 1,
    total_items: 0
  });
  
  const navigate = useNavigate();

  useEffect(() => {
    fetchPosts();
  }, [activeTab, pagination.current_page]);

  const fetchPosts = async () => {
    setLoading(true);
    setError(null);
    
    try {
      let response;
      
      if (activeTab === 'all') {
        response = await postService.getMyPosts(pagination.current_page, pagination.per_page);
      } else if (activeTab === 'published') {
        response = await postService.getMyPostsByPublishedStatus(true, pagination.current_page, pagination.per_page);
      } else if (activeTab === 'drafts') {
        response = await postService.getMyPostsByPublishedStatus(false, pagination.current_page, pagination.per_page);
      }
      
      setPosts(response.data || []);
      setPagination(response.pagination || pagination);
    } catch (err) {
      console.error('Error fetching posts:', err);
      setError('Failed to load posts. Please try again later.');
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (tab) => {
    setActiveTab(tab);
    setPagination({ ...pagination, current_page: 1 });
  };

  const handlePageChange = (page) => {
    setPagination({ ...pagination, current_page: page });
  };

  const handlePublishToggle = async (post) => {
    try {
      if (post.published) {
        await postService.unpublishPost(post.id);
      } else {
        await postService.publishPost(post.id);
      }
      
      // Refresh the post list
      fetchPosts();
    } catch (err) {
      console.error('Error toggling publish status:', err);
      setError('Failed to update publish status. Please try again.');
    }
  };

  const handleDelete = async (postId) => {
    if (window.confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
      try {
        await postService.deletePost(postId);
        
        // Refresh the post list
        fetchPosts();
      } catch (err) {
        console.error('Error deleting post:', err);
        setError('Failed to delete post. Please try again.');
      }
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // Generate pagination buttons
  const renderPagination = () => {
    const { current_page, total_pages } = pagination;
    
    // If there's only one page, don't show pagination
    if (total_pages <= 1) return null;
    
    let pages = [];
    const maxPagesToShow = 5; // Max number of page buttons to show
    
    // Calculate start and end page
    let startPage = Math.max(1, current_page - Math.floor(maxPagesToShow / 2));
    let endPage = Math.min(total_pages, startPage + maxPagesToShow - 1);
    
    // Adjust startPage if endPage is maxed out
    if (endPage === total_pages) {
      startPage = Math.max(1, endPage - maxPagesToShow + 1);
    }
    
    // Add "Previous" button
    pages.push(
      <li key="prev">
        <a 
          className="pagination-link" 
          onClick={() => handlePageChange(Math.max(1, current_page - 1))}
          disabled={current_page === 1}
        >
          <span className="icon">
            <i className="fas fa-chevron-left"></i>
          </span>
        </a>
      </li>
    );
    
    // Add ellipsis if needed before page numbers
    if (startPage > 1) {
      pages.push(
        <li key="ellipsis-start">
          <span className="pagination-ellipsis">&hellip;</span>
        </li>
      );
    }
    
    // Add page number buttons
    for (let i = startPage; i <= endPage; i++) {
      pages.push(
        <li key={i}>
          <a 
            className={`pagination-link ${i === current_page ? 'is-current' : ''}`}
            aria-label={`Go to page ${i}`}
            onClick={() => handlePageChange(i)}
          >
            {i}
          </a>
        </li>
      );
    }
    
    // Add ellipsis if needed after page numbers
    if (endPage < total_pages) {
      pages.push(
        <li key="ellipsis-end">
          <span className="pagination-ellipsis">&hellip;</span>
        </li>
      );
    }
    
    // Add "Next" button
    pages.push(
      <li key="next">
        <a 
          className="pagination-link" 
          onClick={() => handlePageChange(Math.min(total_pages, current_page + 1))}
          disabled={current_page === total_pages}
        >
          <span className="icon">
            <i className="fas fa-chevron-right"></i>
          </span>
        </a>
      </li>
    );
    
    return (
      <nav className="pagination is-centered" role="navigation" aria-label="pagination">
        <ul className="pagination-list">
          {pages}
        </ul>
      </nav>
    );
  };

  return (
    <section className="section">
      <div className="container">
        <div className="level">
          <div className="level-left">
            <div className="level-item">
              <h1 className="title">Manage Posts 1</h1>
            </div>
          </div>
          <div className="level-right">
            <div className="level-item">
              <Link to="/posts/create" className="button is-primary">
                <span className="icon">
                  <i className="fas fa-plus"></i>
                </span>
                <span>New Post</span>
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

        <div className="tabs">
          <ul>
            <li className={activeTab === 'all' ? 'is-active' : ''}>
              <a onClick={() => handleTabChange('all')}>All Posts</a>
            </li>
            <li className={activeTab === 'published' ? 'is-active' : ''}>
              <a onClick={() => handleTabChange('published')}>Published</a>
            </li>
            <li className={activeTab === 'drafts' ? 'is-active' : ''}>
              <a onClick={() => handleTabChange('drafts')}>Drafts</a>
            </li>
          </ul>
        </div>

        <div className="card">
          <div className="card-content">
            {loading ? (
              <div className="has-text-centered py-6">
                <button className="button is-loading is-large is-white"></button>
                <p className="mt-3">Loading posts...</p>
              </div>
            ) : posts.length === 0 ? (
              <div className="has-text-centered py-6">
                <p className="is-size-5">No posts found.</p>
                <Link to="/posts/create" className="button is-primary mt-4">
                  Create your first post
                </Link>
              </div>
            ) : (
              <div className="table-container">
                <table className="table is-fullwidth is-hoverable">
                  <thead>
                    <tr>
                      <th>Title</th>
                      <th>Type</th>
                      <th>Status</th>
                      <th>Created</th>
                      <th>Updated</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {posts.map((post) => (
                      <tr key={post.id}>
                        <td>
                          <Link to={`/posts/edit/${post.id}`} className="has-text-weight-semibold">
                            {post.title}
                          </Link>
                        </td>
                        <td>{post.type}</td>
                        <td>
                          <span className={`tag ${post.published ? 'is-success' : 'is-warning'}`}>
                            {post.published ? 'Published' : 'Draft'}
                          </span>
                        </td>
                        <td>{formatDate(post.createdAt)}</td>
                        <td>{formatDate(post.updatedAt)}</td>
                        <td>
                          <div className="buttons are-small">
                            <Link 
                              to={`/posts/edit/${post.id}`} 
                              className="button is-link is-outlined"
                              title="Edit"
                            >
                              <span className="icon">
                                <i className="fas fa-edit"></i>
                              </span>
                            </Link>
                            <button 
                              className={`button ${post.published ? 'is-warning' : 'is-success'} is-outlined`}
                              title={post.published ? 'Unpublish' : 'Publish'}
                              onClick={() => handlePublishToggle(post)}
                            >
                              <span className="icon">
                                <i className={`fas ${post.published ? 'fa-eye-slash' : 'fa-eye'}`}></i>
                              </span>
                            </button>
                            <button 
                              className="button is-danger is-outlined"
                              title="Delete"
                              onClick={() => handleDelete(post.id)}
                            >
                              <span className="icon">
                                <i className="fas fa-trash-alt"></i>
                              </span>
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
          
          {!loading && pagination.total_pages > 1 && (
            <div className="card-footer">
              <div className="card-footer-item">
                {renderPagination()}
              </div>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}

export default PostsList;