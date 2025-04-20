# OKBlog Admin

A React-based admin interface for managing blog posts in the OKBlog platform.

## Features

- View all posts with pagination
- Filter posts by published/draft status
- Create new posts
- Edit existing posts
- Publish/unpublish posts
- Delete posts
- Tag management
- Post type selection (Post or Page)

## Technologies Used

- React 19+
- React Router v6
- Bulma CSS Framework
- Axios for API communication
- Font Awesome icons

## Development Setup

1. Install dependencies:
   ```
   npm install
   ```

2. Start the development server:
   ```
   npm start
   ```

3. Build for production:
   ```
   npm run build
   ```

## Authentication

The admin interface uses JWT token-based authentication. The current implementation includes a simple mock authentication for development purposes. In a production environment, this should be connected to a proper authentication service.

Default development credentials:
- Username: `admin`
- Password: `password`

## API Connection

The application connects to the backend API through a proxy configuration in the Vite setup. The default configuration assumes the API is running at `http://localhost:8080`.

You can modify the API endpoint in the `vite.config.js` file if needed.

## Project Structure

- `src/components/` - Reusable UI components
- `src/contexts/` - React context providers
- `src/pages/` - Page components
- `src/services/` - API service functions

## License

Copyright © 2023 OKBlog 