import React, { useState, useEffect } from 'react';
import ReactQuill from 'react-quill';
import 'quill/dist/quill.snow.css';

function PostForm({ onSubmit, isLoading, initialData, isEdit = false }) {
  const [title, setTitle] = useState('');
  const [slug, setSlug] = useState('');
  const [type, setType] = useState('POST');
  const [content, setContent] = useState('');
  const [excerpt, setExcerpt] = useState('');
  const [tagsInput, setTagsInput] = useState('');
  const [published, setPublished] = useState(false);
  const [autoSlug, setAutoSlug] = useState(true);
  const [autoExcerpt, setAutoExcerpt] = useState(true);

  // Quill editor modules and formats
  const modules = {
    toolbar: [
      [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
      ['bold', 'italic', 'underline', 'strike', 'blockquote'],
      [{ 'list': 'ordered' }, { 'list': 'bullet' }],
      ['link', 'image', 'code-block'],
      ['clean']
    ],
  };
  
  const formats = [
    'header',
    'bold', 'italic', 'underline', 'strike', 'blockquote',
    'list', 'bullet',
    'link', 'image', 'code-block'
  ];

  // Initialize form with data if provided
  useEffect(() => {
    if (initialData) {
      setTitle(initialData.title || '');
      setSlug(initialData.slug || '');
      setType(initialData.type || 'POST');
      setContent(initialData.content || '');
      setExcerpt(initialData.excerpt || '');
      setTagsInput(initialData.tags ? Array.from(initialData.tags).join(', ') : '');
      setPublished(initialData.published || false);
      
      // Disable auto generation if we have data
      if (initialData.slug) setAutoSlug(false);
      if (initialData.excerpt) setAutoExcerpt(false);
    }
  }, [initialData]);

  // Auto-generate slug from title
  useEffect(() => {
    if (autoSlug && title) {
      const generatedSlug = title
        .toLowerCase()
        .replace(/[^\w\s-]/g, '') // Remove special characters
        .replace(/\s+/g, '-') // Replace spaces with hyphens
        .replace(/-+/g, '-'); // Replace multiple hyphens with a single one
      
      setSlug(generatedSlug);
    }
  }, [title, autoSlug]);

  // Auto-generate excerpt from content
  useEffect(() => {
    if (autoExcerpt && content) {
      // Take first 150 characters of content without HTML tags
      const textContent = content.replace(/<[^>]*>/g, '');
      const generatedExcerpt = textContent.substring(0, 150) + (textContent.length > 150 ? '...' : '');
      
      setExcerpt(generatedExcerpt);
    }
  }, [content, autoExcerpt]);

  const handleSubmit = (e) => {
    e.preventDefault();
    
    // Convert tags string to array
    const tagsArray = tagsInput
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag); // Remove empty tags
    
    const postData = {
      title,
      slug: slug || undefined, // Don't send empty string
      type,
      content,
      excerpt: excerpt || undefined, // Don't send empty string
      tags: tagsArray.length > 0 ? tagsArray : undefined,
      published
    };
    
    onSubmit(postData);
  };

  const handleSlugChange = (e) => {
    setAutoSlug(false);
    setSlug(e.target.value);
  };

  const handleExcerptChange = (e) => {
    setAutoExcerpt(false);
    setExcerpt(e.target.value);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="card">
        <div className="card-content">
          <div className="field">
            <label className="label">Title</label>
            <div className="control">
              <input 
                className="input" 
                type="text" 
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Post title"
                required
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="field">
            <label className="label">
              Slug
              <span className="has-text-grey is-size-7 ml-2">
                (URL-friendly identifier)
              </span>
            </label>
            <div className="control">
              <input 
                className="input" 
                type="text" 
                value={slug}
                onChange={handleSlugChange}
                placeholder="post-url-slug"
                disabled={isLoading}
              />
            </div>
            <p className="help">
              <label className="checkbox">
                <input 
                  type="checkbox" 
                  checked={autoSlug}
                  onChange={() => setAutoSlug(!autoSlug)}
                  disabled={isLoading}
                />
                <span className="ml-2">Auto-generate from title</span>
              </label>
            </p>
          </div>

          <div className="columns">
            <div className="column">
              <div className="field">
                <label className="label">Type</label>
                <div className="control">
                  <div className="select is-fullwidth">
                    <select 
                      value={type}
                      onChange={(e) => setType(e.target.value)}
                      disabled={isLoading}
                    >
                      <option value="POST">Post</option>
                      <option value="PAGE">Page</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
            <div className="column">
              <div className="field">
                <label className="label">Tags</label>
                <div className="control">
                  <input 
                    className="input" 
                    type="text" 
                    value={tagsInput}
                    onChange={(e) => setTagsInput(e.target.value)}
                    placeholder="tag1, tag2, tag3"
                    disabled={isLoading}
                  />
                </div>
                <p className="help">Separate tags with commas</p>
              </div>
            </div>
          </div>

          <div className="field">
            <label className="label">Content</label>
            <div className="control">
              <ReactQuill 
                theme="snow"
                value={content}
                onChange={setContent}
                modules={modules}
                formats={formats}
                placeholder="Write your post content here..."
                readOnly={isLoading}
                className="content-editor"
                style={{ minHeight: "300px", marginBottom: "40px" }}
              />
            </div>
          </div>

          <div className="field">
            <label className="label">
              Excerpt
              <span className="has-text-grey is-size-7 ml-2">
                (Short summary of the post)
              </span>
            </label>
            <div className="control">
              <textarea 
                className="textarea" 
                value={excerpt}
                onChange={handleExcerptChange}
                placeholder="Brief description for previews and search results"
                rows="3"
                disabled={isLoading}
              />
            </div>
            <p className="help">
              <label className="checkbox">
                <input 
                  type="checkbox" 
                  checked={autoExcerpt}
                  onChange={() => setAutoExcerpt(!autoExcerpt)}
                  disabled={isLoading}
                />
                <span className="ml-2">Auto-generate from content</span>
              </label>
            </p>
          </div>

          <div className="field">
            <div className="control">
              <label className="checkbox">
                <input 
                  type="checkbox" 
                  checked={published}
                  onChange={() => setPublished(!published)}
                  disabled={isLoading}
                />
                <span className="ml-2">Publish immediately</span>
              </label>
            </div>
          </div>
        </div>

        <div className="card-footer">
          <div className="card-footer-item">
            <div className="buttons">
              <button 
                type="submit" 
                className={`button is-primary ${isLoading ? 'is-loading' : ''}`}
                disabled={isLoading}
              >
                <span className="icon">
                  <i className="fas fa-save"></i>
                </span>
                <span>{isEdit ? 'Update' : 'Save'}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
}

export default PostForm; 