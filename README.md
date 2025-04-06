
# 🔍 Apica In-Memory Search Engine – Full Stack (Golang + React)

## 📃 Project Overview
This project is a blazing-fast, in-memory log search engine built from scratch using **Golang** for the backend and **React (Vite)** for the frontend. It allows users to search through `.parquet` log files without relying on traditional databases or search services.

The app supports:
- Real-time keyword-based log searching
- Paginated and non-paginated views
- Performance metrics (duration, total matches)
- Dynamic Parquet file uploads
- A clean, responsive, and intuitive user interface

Designed as part of the **Apica Fullstack Assignment**, this project showcases modular code architecture, real-time UI interaction, and optimized API integration.

## 🛠 Technologies Used

### 🔧 Backend (Golang)
- **Go 1.20+**
- `net/http`, `httprouter`, `alice`
- `xitongsys/parquet-go` for reading Parquet files
- Custom structured logging
- Graceful shutdown with `context`

### 🌐 Frontend (React + Vite)
- **React 19 + Vite**
- `axios` for API communication
- `react-toastify` for notifications
- Dynamic table rendering, modals, and pagination
- Inline CSS (no external UI libraries used)

## 📂 Folder Structure

```
search-engine-golang-react/
├── backend/
│   ├── main.go
│   ├── routes/
│   ├── model/
│   ├── build.sh
│   └── go.mod / go.sum
│
├── frontend/
│   ├── src/
│   │   ├── App.jsx
│   │   ├── ResultsTable.jsx
│   │   └── apiUrls.js
│   ├── package.json
│   └── vite.config.js
```

## 🔄 How to Run

### 🛋️ Backend (Golang)
```bash
   cd backend
   go mod tidy
   export PARQUET_FILE_DIRECTORT=./parquet-files  # Choose any directory to store the parquet files
   go run main.go
```
Or using the shell script:
```bash
  sh build.sh
```
Server runs on: `http://localhost:9000`

### 🌐 Frontend (React + Vite)
```bash
  cd frontend
  npm install
  npm run dev
```
Frontend runs on: `http://localhost:5173`

## 🔢 Backend Logic Explained
- **main.go** initializes the HTTP server and routes with graceful shutdown handling.
- **router.go** defines route wrappers and panic recovery using `httprouter` and `alice`.
- **parquet.go** handles endpoints:
  - `GET /v1/list/log/entries`: Paginated log search
  - `GET /v1/fetch/all/records`: Full list (no pagination)
  - `POST /v1/upload/perquet`: Upload `.parquet` file
- **common.go** includes:
  - API logging
  - Response marshaling (success/failure structure)
  - Request/response timing logic
- **logentry.go** defines all log fields like `Message`, `StructuredData`, `NanoTimeStamp` from the Parquet schema.

### Techniques used:
- In-memory data storage using Go slices
- Parquet decoding using struct tags
- Dynamic filtering and substring match on multiple fields
- Efficient response marshalling for large datasets

## 🌈 Frontend Logic Explained

### `App.jsx`
- Manages all global states:
  - Search query, file upload, pagination, loading status
- Makes dynamic API calls via Axios
- Supports two search modes:
  - **TYPE_1**: Paginated (limit/offset) – **More efficient and faster**, especially for large datasets, as it fetches only required chunks.
  - **TYPE_2**: Full fetch – **Slower compared to TYPE_1**, as it fetches **all matching records** from the `.parquet` files without pagination.
- Implements real-time feedback using React Toasts
- Uploads `.parquet` files using multipart form data

### `ResultsTable.jsx`
- Accepts props like `data`, `pagination`, `query`
- Truncates long fields and shows full JSON in modals
- Highlights matching keywords using RegEx
- Implements chunked rendering for performance in TYPE_2
- Custom pagination and limit dropdown at bottom left

### `apiUrls.js`
- Stores API routes using env-based `BASE_URL`

## 📊 Performance & UI Features
- Keyword highlight
- Search duration and match count
- Truncated fields with tooltips
- Modal popup for JSON fields
- Smooth pagination
- File upload with validation

## 🔹 API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ping` | Health check |
| GET | `/v1/list/log/entries` | Paginated search |
| GET | `/v1/fetch/all/records` | Full list fetch |
| POST | `/v1/upload/perquet` | Upload Parquet file |

## 📈 Pros
- Blazing fast — all in-memory
- Simple and modular design
- No DB, No Elasticsearch
- Beautiful UI with instant feedback

## ❌ Cons
- In-memory = not scalable to TB scale
- No data persistence
- Lacks auth/security (for demo only)

## 📄 Screenshots
![Search UI](./assets/Screenshot%202025-04-06%20225255.png)
![Result Table](./assets/Screenshot%202025-04-06%20225304.png)
![Duration](./assets/Screenshot%202025-04-06%20225327.png)

## 🌟 Conclusion
This full-stack project demonstrates how to build a scalable, responsive, and super-fast log search engine using just Golang and React. The code is modular, clean, and highly performant for log-level diagnostics and real-time search analytics.

> 🔥 No DB. No search engine. Just pure logic and performance.
