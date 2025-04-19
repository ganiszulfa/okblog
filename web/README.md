# OKBlog Web Frontend

A Vue.js SSR frontend for the OKBlog application. This project uses Nuxt.js for server-side rendering.

## Features

- Server-side rendering for improved SEO and performance
- Responsive design with Tailwind CSS
- Blog post listing with pagination
- Detailed post view with author information
- View count tracking

## Prerequisites

- Node.js (v14 or later recommended)
- npm or yarn

## Setup

1. Clone the repository
2. Navigate to the web directory
3. Install dependencies:

```bash
npm install
# or
yarn install
```

## Configuration

The application is configured to communicate with the API at `http://localhost:8080` by default. You can customize this by setting the `API_URL` environment variable:

```bash
export API_URL=http://your-api-url
```

## Development

To start the development server:

```bash
npm run dev
# or
yarn dev
```

The application will be available at `http://localhost:3000`.

## Production Build

To build for production:

```bash
npm run build
# or
yarn build
```

To start the production server:

```bash
npm run start
# or
yarn start
```

## Static Generation

If you prefer to generate a static site:

```bash
npm run generate
# or
yarn generate
```

This will generate static HTML files in the `dist` directory.

## Docker Support

A Dockerfile is included for containerized deployment. To build the Docker image:

```bash
docker build -t okblog-web .
```

To run the container:

```bash
docker run -p 3000:3000 okblog-web
``` 