# dbstore

dbstore is a Go project for dumping MySQL databases, encrypting them with a password, and uploading them to an S3 bucket.

## Features

- Dump MySQL databases securely.
- Encrypt dumps with a password.
- Upload encrypted dumps to an S3 bucket.
- Configurable retry mechanism for failed operations.

## Getting Started

### Docker-compose

```yaml
services:
  dbstore:
    image: qnton/dbstore
    environment:
      DB_USER: <your_db_user>
      DB_NAME: <your_db_name>
      DB_PASSWORD: <your_db_password>
      DB_HOST: <your_db_host>
      S3_NAME: <your_s3_name>
      S3_ENDPOINT: <your_s3_endpoint>
      S3_ACCESS_KEY_ID: <your_access_key_id>
      S3_SECRET_ACCESS_KEY: <your_secret_access_key>
      ATTEMPTS: <your_attempts>
      PASSWORD: <your_password>
      INTERVAL: <your_interval>
```

### Build yourself

Clone the repository:

```bash
git clone https://github.com/qnton/dbstore
```

Navigate to the project directory:

```bash
cd dbstore
```

Create Docker-compose:

```yaml
services:
  app:
    build: .
    environment:
    DB_USER: <your_db_user>
    DB_NAME: <your_db_name>
    DB_PASSWORD: <your_db_password>
    DB_HOST: <your_db_host>
    S3_NAME: <your_s3_name>
    S3_ENDPOINT: <your_s3_endpoint>
    S3_ACCESS_KEY_ID: <your_access_key_id>
    S3_SECRET_ACCESS_KEY: <your_secret_access_key>
    ATTEMPTS: <your_attempts>
    PASSWORD: <your_password>
    INTERVAL: <your_interval>
```

```bash
docker-compose up --build
```

Or build your own Docker image:

```bash
docker build -t dbstore .
docker run -d --name dbstore-container \
  -e DB_USER=$DB_USER \
  -e DB_NAME=$DB_NAME \
  -e DB_PASSWORD=$DB_PASSWORD \
  -e DB_HOST=$DB_HOST \
  -e S3_NAME=$S3_NAME \
  -e S3_ENDPOINT=$S3_ENDPOINT \
  -e S3_ACCESS_KEY_ID=$S3_ACCESS_KEY_ID \
  -e S3_SECRET_ACCESS_KEY=$S3_SECRET_ACCESS_KEY \
  -e ATTEMPTS=$ATTEMPTS \
  -e PASSWORD=$PASSWORD \
  -e INTERVAL=$INTERVAL \
  dbstore
```

## Configuration

Configure the behavior of dbstore using the following environment variables:

| Variable Name        | Description                                    |
| -------------------- | ---------------------------------------------- |
| DB_USER              | Database username.                             |
| DB_NAME              | Database name.                                 |
| DB_PASSWORD          | Database password.                             |
| DB_HOST              | Database host address.                         |
| S3_NAME              | Name of the storage bucket.                    |
| S3_ENDPOINT          | Endpoint URL for the storage bucket.           |
| S3_ACCESS_KEY_ID     | Access key ID for the storage bucket.          |
| S3_SECRET_ACCESS_KEY | Secret access key for the storage bucket.      |
| ATTEMPTS             | Number of attempts to retry failed operations. |
| PASSWORD             | Password for encrypting the dump.              |
| INTERVAL             | Interval between dumps (in seconds).           |

Ensure that you set these variables appropriately for your environment and requirements.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
