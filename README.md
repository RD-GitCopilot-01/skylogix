# Todo App - Golang Backend + React TypeScript Frontend

A simple fullstack Todo application built with Golang backend and React + TypeScript frontend.

## Features

- ✅ Create new todos
- ✅ Mark todos as completed/incomplete
- ✅ Delete todos
- ✅ View completion progress
- ✅ Clean, responsive UI with Tailwind CSS
- ✅ RESTful API with CORS support

## Tech Stack

### Backend (Golang)
- Go 1.18+
- Gorilla Mux for routing
- CORS middleware
- In-memory data storage
- RESTful API design

### Frontend (React + TypeScript)
- React 18
- TypeScript
- Vite build tool
- Tailwind CSS for styling
- shadcn/ui components
- Lucide React icons

## API Endpoints

- `GET /api/todos` - Get all todos
- `POST /api/todos` - Create a new todo
- `PUT /api/todos/{id}` - Update todo completion status
- `DELETE /api/todos/{id}` - Delete a todo

## Local Development

### Backend Setup
```bash
cd backend
go mod tidy
go run main.go
```
Backend will run on http://localhost:8000

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```
Frontend will run on http://localhost:5173

## Project Structure

```
todo-app/
├── backend/
│   ├── main.go          # Main Go server with API endpoints
│   ├── go.mod           # Go module dependencies
│   └── go.sum           # Go module checksums
├── frontend/
│   ├── src/
│   │   ├── App.tsx      # Main React component
│   │   ├── components/  # shadcn/ui components
│   │   └── ...
│   ├── package.json     # Node.js dependencies
│   └── vite.config.ts   # Vite configuration
└── README.md
```

## Features Implemented

- **CRUD Operations**: Full Create, Read, Update, Delete functionality
- **Real-time Updates**: Frontend updates immediately after API calls
- **Error Handling**: Proper error handling for API failures
- **Responsive Design**: Works on desktop and mobile devices
- **Type Safety**: Full TypeScript support in frontend
- **Clean Architecture**: Separation of concerns between frontend and backend

## Development Notes

- Backend uses in-memory storage (data resets on server restart)
- CORS is configured to allow frontend-backend communication
- Frontend uses environment variables for API URL configuration
- All API calls include proper error handling and loading states
