import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [loginError, setLoginError] = useState('');
  
  const { login, error } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  
  // Get the page to redirect to after login
  const from = location.state?.from?.pathname || '/';
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!username.trim() || !password.trim()) {
      setLoginError('Username and password are required');
      return;
    }
    
    setIsLoading(true);
    setLoginError('');
    
    try {
      const success = await login({ username, password });
      if (success) {
        navigate(from, { replace: true });
      } else {
        setLoginError('Login failed. Please check your credentials.');
      }
    } catch (err) {
      setLoginError(err.message || 'An error occurred during login');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <section className="section">
      <div className="container">
        <div className="columns is-centered">
          <div className="column is-one-third">
            <div className="box">
              <h1 className="title has-text-centered">Login</h1>
              
              <form onSubmit={handleSubmit}>
                <div className="field">
                  <label className="label">Username</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="text"
                      placeholder="Enter your username"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                    <span className="icon is-small is-left">
                      <i className="fas fa-user"></i>
                    </span>
                  </div>
                </div>
                
                <div className="field">
                  <label className="label">Password</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="password"
                      placeholder="Enter your password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                    <span className="icon is-small is-left">
                      <i className="fas fa-lock"></i>
                    </span>
                  </div>
                </div>
                
                {(loginError || error) && (
                  <div className="notification is-danger">
                    <button className="delete" onClick={() => setLoginError('')}></button>
                    {loginError || error}
                  </div>
                )}
                
                <div className="field">
                  <div className="control">
                    <button
                      type="submit"
                      className={`button is-primary is-fullwidth ${isLoading ? 'is-loading' : ''}`}
                      disabled={isLoading}
                    >
                      Login
                    </button>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default Login; 