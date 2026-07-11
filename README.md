# Sany

Sany is a small file-sharing service. It accepts a multipart file, stores its metadata and content in SQLite, and returns a public link for downloading it.

## Run

Set `SERVER_HOST`, `SERVER_PORT`, and `DATASOURCE_PATH` in `.env`, then start the server:

```sh
make run
```

Upload with `POST /api/v1/files` using the multipart field `file`. Download with the returned `GET /api/v1/files?id=<id>` link.

## Commands

Run `make help` to see available commands. Use `make api` for the end-to-end API checks after starting the server. The current system design is in [design/system/v1.1.png](design/system/v1.1.png).
