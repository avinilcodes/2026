# My Website - Angular Frontend

This is an Angular frontend application that connects to the Go backend API.

## Project Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── components/
│   │   │   ├── home/           # Home page component
│   │   │   ├── profile/        # Profile component
│   │   │   ├── blogs/          # Blogs list component
│   │   │   └── blog-detail/    # Single blog detail component
│   │   ├── services/
│   │   │   └── api.service.ts  # API service for backend calls
│   │   ├── models/
│   │   │   └── models.ts       # TypeScript interfaces
│   │   ├── app.component.*     # Root component
│   │   └── app.routes.ts       # Routing configuration
│   ├── main.ts                 # Application entry point
│   └── styles.css              # Global styles
├── angular.json                # Angular configuration
├── tsconfig.json               # TypeScript configuration
└── package.json                # Dependencies
```

## Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

## Running the Application

Start the development server:
```bash
npm start
```

The application will be available at `http://localhost:4200`

## API Configuration

The API service is configured to connect to the Go backend at `http://localhost:8080/api`.

If your backend is running on a different host/port, update the `apiUrl` in [src/app/services/api.service.ts](src/app/services/api.service.ts):

```typescript
private apiUrl = 'http://localhost:8080/api';
```

## Available Routes

- `/` - Home page
- `/profile` - Your profile information
- `/blogs` - List of all blog posts
- `/blog/:id` - Individual blog post detail

## Features

### Home Component
Displays welcome message and portfolio introduction

### Profile Component
Shows your professional profile with:
- Name and title
- Bio
- Contact email
- Skills list
- Social media links (GitHub, LinkedIn)

### Blogs Component
Displays all blog posts in a grid layout with:
- Blog title
- Author and publication date
- Content preview
- Tags
- "Read More" link to full post

### Blog Detail Component
Shows the full blog post with:
- Title
- Author and date
- Tags
- Full content

## Building for Production

```bash
npm run build
```

The production build will be generated in the `dist/` directory.

## Technologies Used

- Angular 17
- TypeScript 5.2
- RxJS 7.8
- Angular Router
