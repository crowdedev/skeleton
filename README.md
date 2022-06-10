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

- Create database

- Copy `env.example` to `.env` and changes some value

- Run using `task serve`

![Default Empty](assets/imgs/empty-run.png)

- Open your browser and open `http://localhost:7777` or port assigned by you

![Swagger](assets/imgs/empty-swagger.png)

### Create New Module


- Run `task module -- register`

- Follow the instructions 

![Module Register](assets/imgs/module-register.png)

- Bima will generate `todos` folder as your module space, creating `protos/todo.proto`, register your module in `configs/modules.yaml` and register your Dependency Injection defined in `todos/dics/todo.go` to `provider.go`

- Refresh your browser

![Module Swagger](assets/imgs/module-swagger.png)

## Register Request Filter
