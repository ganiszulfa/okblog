import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import Navbar from './components/Navbar';
import PostsList from './pages/PostsList';
import EditPost from './pages/EditPost';
import CreatePost from './pages/CreatePost';
import Login from './pages/Login';
import { AuthProvider, useAuth } from './contexts/AuthContext';

// Protected route component
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return (
      <div className="section">
        <div className="container has-text-centered">
          <div className="button is-loading is-large"></div>
          <p className="mt-3">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
};

function AppContent() {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  return (
    <>
      <Navbar />
      <main>
        <Routes>
          <Route 
            path="/login" 
            element={isAuthenticated ? <Navigate to="/" replace /> : <Login />} 
          />
          <Route 
            path="/" 
            element={
              <ProtectedRoute>
                <PostsList />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/posts/create" 
            element={
              <ProtectedRoute>
                <CreatePost />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/posts/edit/:id" 
            element={
              <ProtectedRoute>
                <EditPost />
              </ProtectedRoute>
            } 
          />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App; 