import React, { createContext, useState, useContext, useEffect } from 'react';

// Create auth context
const AuthContext = createContext(null);

// Auth provider component
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check if user is already logged in
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = localStorage.getItem('auth_token');
        if (!token) {
          setLoading(false);
          return;
        }
        
        // Validate token with backend or use JWT decode
        // For now, we'll just assume the token presence means logged in
        setUser({
          id: localStorage.getItem('user_id'),
          name: localStorage.getItem('user_name'),
        });
      } catch (err) {
        console.error('Auth check error:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  // Login function
  const login = async (credentials) => {
    try {
      setLoading(true);
      setError(null);
      
      // This would be an actual API call in production
      // Simulating an API call for now
      const response = await new Promise((resolve) => {
        setTimeout(() => {
          if (credentials.username === 'admin' && credentials.password === 'password') {
            resolve({
              success: true,
              data: {
                token: 'fake-jwt-token',
                user: {
                  id: '1',
                  name: 'Admin User'
                }
              }
            });
          } else {
            resolve({
              success: false,
              error: 'Invalid credentials'
            });
          }
        }, 500);
      });

      if (!response.success) {
        throw new Error(response.error);
      }

      const { token, user: userData } = response.data;
      
      // Store auth info in localStorage
      localStorage.setItem('auth_token', token);
      localStorage.setItem('user_id', userData.id);
      localStorage.setItem('user_name', userData.name);
      
      setUser(userData);
      return true;
    } catch (err) {
      setError(err.message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  // Logout function
  const logout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('user_name');
    setUser(null);
  };

  // Auth context value
  const value = {
    user,
    isAuthenticated: !!user,
    loading,
    error,
    login,
    logout
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

// Custom hook to use auth context
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
} 