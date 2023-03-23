# Home Drive

A file storage service built in Go that allows users to upload and store files of any type and size, and generate a shareable URL for each file. For image files, the service also generates a blurHash that can be used to display a low-resolution preview of the image before it fully loads.

The service supports image resizing, allowing users to specify the desired width and height of the image in the URL query parameters. For other file types, these query parameters will have no effect.

## Features
1. Upload and store files of any type and size
1. Generate shareable URLs for each file
1. Generate blurHash previews for image files
1. Support image resizing and cropping with URL query parameters

## Installation
1. Install Go on your machine
1. Clone the repository
1. Run go build to build the binary
1. Run the binary to start the server
## Usage
Once the server is running, you can use the following endpoints:

- `POST /api/upload` : Upload a file to the server
- `GET <FILE URL RECEIVE FROM UPLOAD>`: Retrieve a file by its filename
- `GET <FILE URL RECEIVE FROM UPLOAD>?s=<size>`: Retrieve an image file with the specified size in square format
- `GET <FILE URL RECEIVE FROM UPLOAD>?w=<width>`: Retrieve an image file with the specified width and auto height based on aspect ratio
- `GET <FILE URL RECEIVE FROM UPLOAD>?w=<width>&h=<height>`: Retrieve an image file with the specified width and height
## License
This project is licensed under the [MIT License](./LICENSE).
