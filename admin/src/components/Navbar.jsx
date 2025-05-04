import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

function Navbar() {
  const [isActive, setIsActive] = useState(false);
  const { user, isAuthenticated, logout } = useAuth();
  const location = useLocation();

  const toggleNavbar = () => {
    setIsActive(!isActive);
  };

  const closeNavbar = () => {
    setIsActive(false);
  };

  const handleLogout = (e) => {
    e.preventDefault();
    logout();
    closeNavbar();
  };

  // Check if a route is active
  const isRouteActive = (path) => {
    return location.pathname === path;
  };

  return (
    <nav className="navbar is-link" role="navigation" aria-label="main navigation">
      <div className="container">
        <div className="navbar-brand">
          <Link className="navbar-item" to="/" onClick={closeNavbar}>
            <strong>OKBlog Admin</strong>
          </Link>

          <a
            role="button"
            className={`navbar-burger ${isActive ? 'is-active' : ''}`}
            aria-label="menu"
            aria-expanded="false"
            onClick={toggleNavbar}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>

        <div className={`navbar-menu ${isActive ? 'is-active' : ''}`}>
          {isAuthenticated && (
            <div className="navbar-start">
              <Link
                className={`navbar-item ${isRouteActive('/') ? 'is-active' : ''}`}
                to="/"
                onClick={closeNavbar}
              >
                <span className="icon-text">
                  <span className="icon">
                    <i className="fas fa-list"></i>
                  </span>
                  <span>Posts</span>
                </span>
              </Link>
              <Link
                className={`navbar-item ${isRouteActive('/posts/create') ? 'is-active' : ''}`}
                to="/posts/create"
                onClick={closeNavbar}
              >
                <span className="icon-text">
                  <span className="icon">
                    <i className="fas fa-plus"></i>
                  </span>
                  <span>New Post</span>
                </span>
              </Link>
            </div>
          )}

          <div className="navbar-end">
            <div className="navbar-item">
              <div className="buttons">
                {isAuthenticated ? (
                  <>
                    <div className="navbar-item">
                      <span className="has-text-white">Hi, {localStorage.getItem('user_firstName')} {localStorage.getItem('user_lastName')}</span>
                    </div>
                    <button className="button is-light" onClick={handleLogout}>
                      <span className="icon">
                        <i className="fas fa-sign-out-alt"></i>
                      </span>
                      <span>Logout</span>
                    </button>
                  </>
                ) : (
                  <Link
                    className="button is-light"
                    to="/login"
                    onClick={closeNavbar}
                  >
                    <span className="icon">
                      <i className="fas fa-sign-in-alt"></i>
                    </span>
                    <span>Log in</span>
                  </Link>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}

export default Navbar; 