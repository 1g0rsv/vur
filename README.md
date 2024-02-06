# VUR Application Overview

This project is of an educational nature, created for personal use and to learn about DevOps tools and technologies. This application is a secure solution for encrypting, transferring, and temporarily storing sensitive data. Developed with Go for backend operations and Vue.js for the frontend interface, it ensures robust security for handling confidential information.

## Key Features

- **Encryption**: Uses the Advanced Encryption Standard (AES) for encrypting text data, managing all processes with unique UUIDs for each data entry.
- **Frontend Interface**: Provides an intuitive interface for users to input and send text for encryption, returning a one-time access link to the decrypted data, enhancing data privacy.
- **Deployment**: Deployed using Docker and managed with Docker Compose, the application components are isolated for efficient operation.
- **Nginx**: Acts as a reverse proxy, facilitating secure interactions between the frontend and backend.

### DevOps Features from the Docker Compose File

#### Containerization and Image Management
- Uses Docker containers for each component (backend, frontend, database, and Nginx), specified with version-controlled images like `igorsky888/vur:vur-app-v1.0.3`.
- Container names like `my_vur_backend` and `my_vur_frontend` for clear identification and management.

#### Service Dependencies and Environment Variables
- `depends_on` attribute manages dependencies, ensuring orderly startup.
- Environment variables such as `HEX_KEY`, `VUR_USER`, `VUR_PASS`, and `VUR_DB` are used for configuration.

#### Persistent Storage and Logging
- Mounted volumes for data persistence (`/root/vur/data`) and logging (`/root/vur/logs`).

#### Network Configuration
- Defined backend and frontend networks ensure segregated and secure communication channels.
- Nginx configured as a reverse proxy for efficient request handling and security.

#### GeoIP Integration with Nginx
- Custom Nginx image (`igorsky888/vur:nginx-geoip-v1.0.0`) includes GeoIP databases for advanced traffic analysis.

### Backend Technical Details (Go)
- **Encryption**: Utilizes AES in CFB mode for reliable and secure encryption.
- **Request Processing**: Encrypted text is saved in a MySQL database with a unique UUID.
- **Decryption and Deletion**: Decrypted text is provided for one-time access and then automatically deleted.

### Frontend (Vue.js)
- Developed in Vue.js, allowing users to securely transmit text for encryption and receive a unique link for decrypted text access.

This system provides a reliable and secure solution for transmitting and temporarily storing encrypted information, with the possibility of its one-time viewing and automatic deletion after access.
