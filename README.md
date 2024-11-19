Archiva is a simple web service built with Go that allows users to upload ZIP files to read their contents and create new ZIP archives from uploaded files.

Features
Read ZIP Files: Upload a ZIP file to see the contents, including file names and sizes.
Create ZIP Archives: Combine multiple files into a ZIP archive for download.
File Type Validation: Supports specific file types only:
DOCX (application/vnd.openxmlformats-officedocument.wordprocessingml.document)
XML (application/xml)
JPEG (image/jpeg)
PNG (image/png)
Requirements
Go 1.18+
Postman (or similar tool) for testing
Usage
API Endpoints
Inspect a ZIP File
Endpoint: POST /api/archive
Content-Type: multipart/form-data
Field: file
Upload a single ZIP file to get a list of the contents. Use Postman to send a POST request with the file in the file field.

Response Example:

json
Копировать код
{
  "filename": "example.zip",
  "total_size": 4.55,
  "total_files": 3,
  "files": [
    {
      "file_path": "document.docx",
      "size": 0.05,
      "mimetype": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    },
    {
      "file_path": "image.jpg",
      "size": 4.43,
      "mimetype": "image/jpeg"
    }
  ]
}
Create a ZIP Archive
Endpoint: POST /api/archive/files
Content-Type: multipart/form-data
Field: files[]
Upload multiple files to create a ZIP archive. Each file should be added to files[]. The ZIP archive will be returned for download.

Response: Returns a downloadable ZIP file.

Testing with Postman
Inspect ZIP:

Create a POST request to http://localhost:8080/api/archive.
Use form-data, add a field called file, and upload a ZIP file.
Send the request and see the JSON response with ZIP details.
Create ZIP:

Create a POST request to http://localhost:8080/api/archive/files.
Use form-data, add multiple files[] fields for each file you want to include in the ZIP.
Send the request and download the ZIP file returned in the response.
