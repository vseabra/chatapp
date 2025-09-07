# Financial Chat Application

A real-time financial chat application with bot integration, built with Go backend, React frontend, and RabbitMQ message broker.

## Overview

This application provides a real-time chat experience where users can:
- Create and join chat rooms
- Send messages that are persisted to MongoDB
- Use bot commands (like `/stock=aapl`) to get financial data
- Receive real-time updates via WebSocket connections

## Architecture

The application follows a microservices architecture with the following components:

### Backend (`/backend`)
- **Technology**: Go with Gin framework
- **Database**: MongoDB for message and user persistence
- **Message Broker**: RabbitMQ for event-driven communication
- **Authentication**: JWT-based authentication
- **WebSocket**: Real-time communication hub

### Frontend (`/frontend`)
- **Technology**: React with TypeScript and Vite
- **UI Library**: shadcn ui
- **State Management**: React Context

### Stock Bot (`/stockbot`)
- **Technology**: Go microservice
- **Purpose**: Processes bot commands (e.g., `/stock aapl`)
- **Integration**: Consumes from RabbitMQ and publishes responses

## How It Works

1. **Message Flow**:
   - User sends a message through the frontend
   - Backend receives the message via websocket
   - Message is persisted to MongoDB
   - Message is published to RabbitMQ channel
   - WebSocket hub listener receives the message
   - Message is broadcast to all active WebSocket connections

2. **Bot Integration**:
   - User sends a bot command (e.g., `/stock aapl`)
   - This message does not get persisted to the database
   - Backend publishes a `bot.requested` event to RabbitMQ
   - Stock bot consumes the event and processes the command
   - Bot publishes response to `bot.response.submit` channel
   - Backend receives the response, persists it, and broadcasts to users

3. **Real-time Updates**:
   - WebSocket connections are managed by a central hub
   - All messages (user and bot) are broadcast to active connections
   - Frontend receives updates and updates the UI in real-time

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)
- Node.js 18+ (for local development)

## Quick Start

### Using Docker Compose (Recommended)


1. **Start all services**:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - RabbitMQ Management: http://localhost:15672 (guest/guest)
   - MongoDB: localhost:27017

### Services

The application consists of 4 main services:

- **`chat-api`** (Port 8080): Main backend API server
- **`frontend`** (Port 5173): React frontend application
- **`stockbot`** (Port 8081): Bot service for processing commands
- **`mongodb`** (Port 27017): Database for persistence
- **`rabbitmq`** (Ports 5672, 15672): Message broker

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Chat Rooms
- `GET /api/chatrooms` - List chat rooms
- `POST /api/chatrooms` - Create chat room
- `GET /api/chatrooms/:id/messages` - Get room messages

### Messages
- `POST /api/messages` - Send message
- `WebSocket /ws` - Real-time connection

## Bot Commands

The stock bot supports the following commands:

- `/stock <symbol>` - Get stock quote (e.g., `/stock aapl`)
- `/echo <message>` - Echo back a message
- `/help` - Show available commands

## Environment Variables

### Backend
- `MONGODB_URI`: MongoDB connection string
- `JWT_SECRET`: Secret for JWT token signing
- `PORT`: Server port (default: 8080)
- `RABBITMQ_URI`: RabbitMQ connection string

### Stock Bot
- `PORT`: Bot service port (default: 8181)
- `RABBITMQ_URI`: RabbitMQ connection string
- `RABBITMQ_EXCHANGE`: Exchange name (default: chat.events)
- `RABBITMQ_QUEUE`: Queue name (default: bot.stockbot)

## Database Schema

### Users
- `id`: Unique identifier
- `username`: User's username
- `email`: User's email
- `password`: Hashed password

### Chat Rooms
- `id`: Unique identifier
- `name`: Room name
- `created_by`: User ID of creator
- `created_at`: Creation timestamp

### Messages
- `id`: Unique identifier
- `room_id`: Chat room ID
- `user_id`: User ID (null for bot messages)
- `content`: Message content
- `type`: Message type (user/bot)
- `created_at`: Creation timestamp

## Message Broker

The application uses RabbitMQ with the following routing keys:

- `bot.requested`: Bot command requests
- `bot.response.submit`: Bot responses
- `message.broadcast`: Message broadcasting

## License

This project is licensed under the MIT License.
