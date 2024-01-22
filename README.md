# **This project is of an educational nature, created for personal use and to learn about DevOps tools and technologies.**
This application is a secure solution for encrypting, transferring, and temporarily storing sensitive data. Developed with Go for backend operations and Vue.js for the frontend interface, it ensures robust security for handling confidential information.

Using the Advanced Encryption Standard (AES), the backend encrypts text data, managing all processes with unique UUIDs for each data entry. The frontend provides an intuitive interface for users to input and send text for encryption, returning a one-time access link to the decrypted data, enhancing data privacy.

Deployed using Docker and managed with Docker Compose, the application components are isolated for efficient operation. Nginx acts as a reverse proxy, facilitating secure interactions between the frontend and backend.

# **Key DevOps Features from the Docker Compose File:**

**Containerization and Image Management:**

The application uses Docker containers for each component (backend, frontend, database, and Nginx), specified with version-controlled images like igorsky888/vur:vur-app-v1.0.3 for the app service. This ensures consistent and replicable environments across development and production setups.
The container names like my_vur_backend and my_vur_frontend indicate clear identification and management of each service.

**Service Dependencies and Environment Variables:**

The depends_on attribute ensures orderly startup by managing dependencies, like the app service depending on the database.
Environment variables such as HEX_KEY, VUR_USER, VUR_PASS, and VUR_DB are used for secure and flexible configuration, promoting best practices in sensitive data management.

**Persistent Storage and Logging:**

The MySQL database service utilizes mounted volumes for data persistence (/root/vur/data) and logging (/root/vur/logs). This setup allows for data durability and simplified log management, crucial for long-term maintenance and debugging.

**Network Configuration:**

Defined backend and frontend networks in Docker Compose file ensure segregated and secure communication channels between different components of the application.
Nginx is configured as a reverse proxy, facilitating efficient request handling and enhanced security.

**GeoIP Integration with Nginx:**

Custom Nginx image (igorsky888/vur:nginx-geoip-v1.0.0) includes GeoIP databases, indicating advanced traffic analysis and geolocation capabilities.

Backend Technical Details (Go)
**1 Encryption:**

   AES (Advanced Encryption Standard) is used for text encryption. AES is a widely used and reliable encryption algorithm.
   Encryption is performed using CFB (Cipher Feedback) mode, which ensures reliability and security.
   The encryption key is stored in an encrypted form in the HEX_KEY environment variable and is converted to bytes before use.

**2 Request Processing:**

   When receiving text from a user via a POST request to /save, the backend encrypts the text and saves it in a MySQL database with a unique UUID.
   After encryption, a unique link with the UUID is generated and returned to the client.

**3 Decryption and Deletion:**

   When accessing /text/<UUID>, the backend retrieves and decrypts the corresponding text using the same AES algorithm.
   The decrypted text is provided to the user, after which the corresponding record is deleted from the database, ensuring one-time access to the information.

**Network Configuration and Routing**
All application components (frontend, backend, and database) are launched in separate Docker containers and managed using Docker Compose.
Nginx is used as a reverse proxy to redirect requests between frontend and backend, allowing to hide internal ports and simplify access to the application.
**Frontend (Vue.js)**
The interface, developed in Vue.js, allows users to enter text, which is then sent to the backend for encryption.
After processing the text, the backend returns a unique link to the frontend for accessing the decrypted text.
This system provides a reliable and secure solution for transmitting and temporarily storing encrypted information, with the possibility of its one-time viewing and automatic deletion after access.