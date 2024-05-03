# dbstore

dbstore is a Go project for dumping MySQL databases, encrypting them with a password, and uploading them to an S3 bucket.

## installation and configuration

When running the container, you can adjust its behavior by passing one or more environment variables. Below is a list of the available variables and their descriptions:

| Variable Name            | Description                                             |
| ------------------------ | ------------------------------------------------------- |
| DATABASE_USER            | The username for connecting to the database.            |
| DATABASE_NAME            | The name of the database to use.                        |
| DATABASE_PASSWORD        | The password for connecting to the database.            |
| DATABASE_HOST            | The host address of the database.                       |
| BUCKET_NAME              | The name of the storage bucket.                         |
| BUCKET_ENDPOINT          | The endpoint URL for the storage bucket.                |
| BUCKET_ACCESS_KEY_ID     | The access key ID for accessing the storage bucket.     |
| BUCKET_SECRET_ACCESS_KEY | The secret access key for accessing the storage bucket. |
| ATTEMPTS                 | The number of attempts to retry failed operations.      |
| PASSWORD                 | A general password for the zip.                         |
| INTERVAL                 | The interval between dumps.                             |

Make sure to set these variables according to your environment and requirements.

## example usage

```
docker run -e DATABASE_USER=myuser \
           -e DATABASE_NAME=mydb \
           -e DATABASE_PASSWORD=mypassword \
           -e DATABASE_HOST=db.example.com \
           -e BUCKET_NAME=mybucket \
           -e BUCKET_ENDPOINT=https://example.com \
           -e BUCKET_ACCESS_KEY_ID=myaccesskey \
           -e BUCKET_SECRET_ACCESS_KEY=mysecretkey \
           -e ATTEMPTS=3 \
           -e PASSWORD=securepassword \
           -e INTERVAL=10 \
           qnton/dbstore:latest
```
