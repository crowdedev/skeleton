## Skeleton

Skeleton is a boilerplate, RESTful generator based on [Bima](https://github.com/KejawenLab/bima)

### Video

Check the [video](https://www.youtube.com/watch?v=zZPpDizZGIM)

### Requirements

- Go 1.17 or above

- Git

- [Taskfile](taskfile.dev)

- RDBMS or MongoDB for database storage

- Elasticsearch (Optional)

- MongoDB (for Logging Extension - Optional)

- RabbitMQ (Optional)

### Basic Usage

- Download using skeleton using git by running `git clone https://github.com/KejawenLab/skeleton.git`

- Download dependencies using `task update` command

- Copy `env.example` to `.env` and changes some value

- Run using `task serve` like below

![Default Empty](assets/imgs/empty-run.png)

- Open your browser and open `http://localhost:7777` or port assigned by you

![Swagger](assets/imgs/empty-swagger.png)

### Create New Module
