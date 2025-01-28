# MeubleHub

MeubleHub is a microservices-based application designed to help families manage and distribute belongings after a loved one's passing. The platform allows for systematic cataloging of items room by room, making the distribution process more organized and transparent.

## Overview

The application helps families to:
- Create and manage homes
- Register rooms within each home
- Catalog objects/belongings
- Track object reservations by family members
- Manage user accounts and authentication

## Architecture

MeubleHub consists of four microservices:

1. **Home Service** (Port: 8081)
   - Manages home entities
   - Uses SQLite for data persistence
   - Handles home creation and listing

2. **Room Service** (Port: 8082)
   - Manages rooms within homes
   - Uses SQLite for data persistence
   - Links rooms to specific homes

3. **Object Service** (Port: 8080)
   - Manages individual objects/belongings
   - Uses DragonflyDB for data storage
   - Handles object creation, listing, and reservation

4. **User Service** (Port: 8083)
   - Manages user accounts and authentication
   - Uses SQLite for data persistence
   - Handles user registration and login

## Prerequisites

- Docker
- Docker Compose
- Git

## Setup and Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/meubleHub.git
cd meubleHub
```

2. Create necessary data directories:
```bash
mkdir -p backend/home-service/data
mkdir -p backend/room-service/data
mkdir -p backend/user-service/data
```

3. Create `.env` file in the root directory:
```env
HOME_PORT=8081
HOME_DB_PATH=/app/data/home.db

OBJECT_PORT=8080
OBJECT_DRAGONFLY_HOST=dragonfly
OBJECT_DRAGONFLY_PORT=6379

ROOM_PORT=8082
ROOM_DB_PATH=/app/data/room.db

USER_PORT=8083
USER_DB_PATH=/app/data/user.db
```

4. Build and start the services:
```bash
docker-compose up --build
```

## API Endpoints

### Home Service (`localhost:8081`)
- `POST /homes` - Create a new home
- `GET /homes` - List all homes

### Room Service (`localhost:8082`)
- `POST /rooms` - Create a new room
- `GET /rooms?home_id=<id>` - List rooms for a specific home

### Object Service (`localhost:8080`)
- `POST /objects` - Create a new object
- `GET /objects` - List all objects
- `PATCH /objects/:id/reserve` - Reserve an object
- `GET /objects/reserved` - List reserved objects

### User Service (`localhost:8083`)
- `POST /users` - Create a new user
- `POST /login` - User login
- `GET /users` - List all users

## Development

### Branch Management
- Create feature branches from main:
```bash
git checkout -b feature/your-feature-name
```
- Submit changes through Pull Requests

### Testing
Each service includes unit tests and integration tests. Run tests for individual services:
```bash
cd backend/service-name
go test ./... -v
```

## Database Management

- SQLite databases are automatically created in the `data` directory of each service
- DragonflyDB is used for object data and runs in a separate container

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Version Information

Current version: v0.1
- All services are tagged with version v0.1
- Docker images are named: `hexagon/service-name:v0.1`

## Maintenance

### Backup
- SQLite databases are stored in volume mounts and should be backed up regularly
- Database files are in the `data` directory of each service

### Logs
- Each service uses structured logging with logrus
- Logs can be viewed using `docker-compose logs service-name`

## Troubleshooting

Common issues and solutions:

1. Database connection issues:
   - Ensure data directories exist and have proper permissions
   - Check volume mounts in docker-compose.yml

2. Service connectivity:
   - Verify all ports are correctly mapped
   - Check network configuration in docker-compose.yml

3. DragonflyDB connection:
   - Ensure DragonflyDB container is healthy
   - Verify connection settings in object service

## Contact

For any questions or support, please open an issue in the repository.