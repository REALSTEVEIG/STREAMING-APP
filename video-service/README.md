# Video Service

The **Video Service** is a high-performance backend service built with Go (Golang) that handles video uploads, transcoding, metadata management, and storage. It integrates with AWS S3 for cloud storage and MongoDB for metadata storage. This service is part of a larger microservices-based architecture for a streaming application.

## Features

- **Video Upload and Storage**: Users can upload videos directly.
- **Transcoding**: Converts uploaded videos to HLS/DASH using FFmpeg for adaptive streaming.
- **Metadata Management**: Stores and retrieves metadata like title, tags, and duration.
- **Short Video Processing**: Supports previews and thumbnails for short videos.

---

## Prerequisites

1. **Go**: Ensure [Go](https://golang.org/doc/install) is installed.
2. **MongoDB**: Running MongoDB instance for metadata storage.
3. **AWS S3**: Configure S3 credentials for cloud storage.
4. **Air**: Hot reload development tool. Install via:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

---

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
MONGO_URI=mongodb://localhost:27017/xxxxxx
AWS_ACCESS_KEY_ID=xxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxx
AWS_REGION=eu-north-1
AWS_S3_BUCKET=xxxxx

```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd video-service
```

### 2. Run the Service with Air

Start the service with live reloading for development:

```bash
air
```

### 3. Run the Service Manually

Alternatively, you can run the service directly with:

```bash
go run main.go
```

---

## Swagger Documentation

After starting the server, visit the Swagger UI for API documentation:

[http://localhost:8080/api/docs/index.html](http://localhost:8080/api/docs/index.html)

---

## Endpoints

### Upload Video

- **Method**: `POST`
- **Path**: `/api/videos/upload`
- **Description**: Upload a video, store it in S3, and save metadata in MongoDB.
- **Request**:
  - `title` (formData string, required): The title of the video.
  - `tags` (formData array, optional): Tags for the video.
  - `file` (formData file, required): The video file to upload.

### Get Video Metadata

- **Method**: `GET`
- **Path**: `/api/videos/{id}`
- **Description**: Retrieve video metadata by its ID.

---

## Architecture

- **Programming Language**: Go (Golang)
- **Database**: MongoDB for metadata storage.
- **Cloud Storage**: AWS S3 for video files.
- **Transcoding**: FFmpeg for adaptive bitrate streaming (HLS/DASH).

---

## Contribution

Contributions are welcome! Please open an issue or submit a pull request.

---

## License

This project is licensed under the MIT License.

---