import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config/api';

function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [bio, setBio] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');
  
  const navigate = useNavigate();
  
  const validateForm = () => {
    // Reset error message
    setError('');
    
    // Check required fields
    if (!username.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      setError('Username, email, and password are required');
      return false;
    }
    
    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError('Please enter a valid email address');
      return false;
    }
    
    // Validate password length (minimum 8 characters)
    if (password.length < 8) {
      setError('Password must be at least 8 characters long');
      return false;
    }
    
    // Check if passwords match
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return false;
    }
    
    return true;
  };
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setIsLoading(true);
    
    try {
      // Call the registration API
      const response = await axios.post(`${API_BASE_URL}/profiles/register`, {
        username,
        email,
        password,
        firstName,
        lastName,
        bio
      });
      
      setSuccessMessage('Registration successful! You can now log in.');
      
      // Clear form fields
      setUsername('');
      setEmail('');
      setPassword('');
      setConfirmPassword('');
      setFirstName('');
      setLastName('');
      setBio('');
      
      // Redirect to login page after 2 seconds
      setTimeout(() => {
        navigate('/login');
      }, 2000);
      
    } catch (err) {
      console.error('Registration error:', err);
      setError(err.response?.data || err.message || 'An error occurred during registration');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <section className="section">
      <div className="container">
        <div className="columns is-centered">
          <div className="column is-two-thirds-tablet is-half-desktop">
            <div className="box">
              <h1 className="title has-text-centered">Create Account</h1>
              
              {successMessage && (
                <div className="notification is-success">
                  <button className="delete" onClick={() => setSuccessMessage('')}></button>
                  {successMessage}
                </div>
              )}
              
              {error && (
                <div className="notification is-danger">
                  <button className="delete" onClick={() => setError('')}></button>
                  {error}
                </div>
              )}
              
              <form onSubmit={handleSubmit}>
                <div className="field">
                  <label className="label">Username*</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="text"
                      placeholder="Enter a username"
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
                  <label className="label">Email*</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="email"
                      placeholder="Enter your email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                    <span className="icon is-small is-left">
                      <i className="fas fa-envelope"></i>
                    </span>
                  </div>
                </div>
                
                <div className="field">
                  <label className="label">Password* (minimum 8 characters)</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="password"
                      placeholder="Enter a password"
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
                
                <div className="field">
                  <label className="label">Confirm Password*</label>
                  <div className="control has-icons-left">
                    <input
                      className="input"
                      type="password"
                      placeholder="Confirm your password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                    <span className="icon is-small is-left">
                      <i className="fas fa-lock"></i>
                    </span>
                  </div>
                </div>
                
                <div className="field">
                  <label className="label">First Name</label>
                  <div className="control">
                    <input
                      className="input"
                      type="text"
                      placeholder="Enter your first name"
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                      disabled={isLoading}
                    />
                  </div>
                </div>
                
                <div className="field">
                  <label className="label">Last Name</label>
                  <div className="control">
                    <input
                      className="input"
                      type="text"
                      placeholder="Enter your last name"
                      value={lastName}
                      onChange={(e) => setLastName(e.target.value)}
                      disabled={isLoading}
                    />
                  </div>
                </div>
                
                <div className="field">
                  <label className="label">Bio</label>
                  <div className="control">
                    <textarea
                      className="textarea"
                      placeholder="Tell us about yourself"
                      value={bio}
                      onChange={(e) => setBio(e.target.value)}
                      disabled={isLoading}
                    ></textarea>
                  </div>
                </div>
                
                <div className="field">
                  <div className="control">
                    <button
                      type="submit"
                      className={`button is-primary is-fullwidth ${isLoading ? 'is-loading' : ''}`}
                      disabled={isLoading}
                    >
                      Register
                    </button>
                  </div>
                </div>
                
                <div className="has-text-centered mt-4">
                  Already have an account? <Link to="/login">Login</Link>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default Register; 