# News Alligator CLI App

The News Alligator CLI App is a command-line application
that aggregates news from various sources based on the provided flags.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed on your system.

## Pulling the Docker Image

The Docker image for the News Alligator CLI App is available on Docker Hub. You can pull the image using the following
command:

```sh
docker pull antohachaban/cli-alligator:0.0.1
```

## Running the Docker Container

To run the News Alligator CLI App, you can use the following command:

```sh
docker run --rm antohachaban/cli-alligator:0.0.1
```

### Running with Flags

The News Alligator CLI App can be run with different flags to customize its behavior. For example:

```sh
docker run --rm antohachaban/cli-alligator:0.0.1 -sources=bbc -keywords=Ukraine
```

## Viewing the Help Message

To see the available flags and options, you can run the following command:

```sh
docker run --rm antohachaban/cli-alligator:0.0.1 --help
```

## Stopping the Docker Container

The --rm option in the run command ensures that the container is automatically removed after it stops.