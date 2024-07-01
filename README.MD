# News Alligator Web App

The News Alligator Web App is a web application that aggregates news from various sources and provides a convenient interface for browsing articles.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed on your system.

## Pulling the Docker Image

The Docker image for the News Alligator Web App is available on Docker Hub. You can pull the image using the following command:

```sh
docker pull antohachaban/news-alligator-web
```

## Running the Docker Container

To run the News Alligator Web App, you can use the following command:

```sh
docker run -d -v news-aggregator-backups:/root/backups -p 8080:8080 antohachaban/news-alligator-web
```

This command will start the News Alligator Web App in a Docker container and expose it on port 8080. 
The `-v news-aggregator-backups:/root/backups` option mounts a volume to store backups of the database.

## Accessing the Web App

Once the container is running, you can access the News Alligator Web App by navigating to `https://localhost:8080` in your web browser.

## Stopping the Docker Container

To stop the News Alligator Web App container, you can use the following command:

```sh
docker stop <container_id>
```

You can get the container ID by running `docker ps` and looking for the container running the `antohachaban/news-alligator-web` image.

## Removing the Docker Container

To remove the News Alligator Web App container, you can use the following command:

```sh
docker rm <container_id>
```

You can get the container ID by running `docker ps -a` and looking for the container with the `Exited` status.

## Removing the Docker Image

To remove the News Alligator Web App image, you can use the following command:

```sh
docker rmi antohachaban/news-alligator-web
```

This will remove the image from your local Docker registry.