# CodeRoulette 🎯

A real-time programming battle platform where developers compete in live coding challenges with skill cards and detailed battle reports.

## 🚀 Features

- **Real-time Battles**: Compete against other programmers in live coding challenges
- **Skill Cards**: Use special abilities to gain advantages during battles
- **Detailed Reports**: Get comprehensive battle reports and performance analytics
- **Spectator Mode**: Watch other battles and learn from top players
- **Multiple Languages**: Support for Go, Python, and JavaScript
- **Leaderboard**: Track your ranking and compete with the best

## 🏗️ Architecture

### Backend (Go)
- **Gin Framework**: High-performance HTTP web framework
- **PostgreSQL**: Primary database for user data, problems, and match history
- **Redis**: Caching and real-time matchmaking queue
- **WebSocket**: Real-time communication for battles
- **Docker**: Safe code execution environment

### Frontend (React + TypeScript)
- **React 18**: Modern React with hooks and functional components
- **Monaco Editor**: VS Code-like code editor
- **WebSocket Client**: Real-time battle communication
- **Responsive Design**: Works on desktop and mobile

## 🛠️ Tech Stack

### Backend
- Go 1.21
- Gin Web Framework
- PostgreSQL 15
- Redis 7
- Docker
- WebSocket (Gorilla)

### Frontend
- React 18
- TypeScript
- Monaco Editor
- Socket.io Client
- CSS3 with modern features

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/coderoulette.git
   cd coderoulette
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PostgreSQL: localhost:5432
   - Redis: localhost:6379

### Local Development

1. **Start databases**
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Backend setup**
   ```bash
   cd backend
   cp env.example .env
   # Edit .env with your database credentials
   go mod download
   go run main.go
   ```

3. **Frontend setup**
   ```bash
   cd frontend
   npm install
   npm start
   ```

## 📁 Project Structure

```
coderoulette/
├── backend/                 # Go backend
│   ├── internal/
│   │   ├── config/         # Configuration
│   │   ├── database/       # Database models and connection
│   │   ├── handlers/       # HTTP handlers
│   │   └── services/       # Business logic
│   ├── main.go
│   └── Dockerfile
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── App.tsx
│   │   └── index.tsx
│   ├── public/
│   └── Dockerfile
├── docker-compose.yml      # Development environment
└── README.md
```

## 🎮 How to Play

1. **Start a Battle**: Choose difficulty and programming language
2. **Get Matched**: System finds an opponent with similar skill level
3. **Code**: Solve the given problem in the time limit
4. **Use Skill Cards**: Deploy special abilities to gain advantages
5. **Submit**: Submit your solution and see results
6. **View Report**: Get detailed battle analysis and performance metrics

## 🃏 Skill Cards

- **Code Peek**: View one line of opponent's code
- **Hint**: Get a hint for the current problem
- **Time Boost**: Get 30 seconds extra time
- **Test Swap**: Swap one test case with opponent
- **Code Lock**: Lock opponent's code for 10 seconds
- **Perfect Score**: Guarantee 100% score on next submission

## 🔧 API Endpoints

### Matches
- `POST /api/v1/matches/queue` - Queue for matchmaking
- `GET /api/v1/matches/status/:id` - Get match status
- `GET /api/v1/matches/queue-status` - Get queue status

### Problems
- `GET /api/v1/problems/random` - Get random problem
- `GET /api/v1/problems/:id` - Get specific problem
- `POST /api/v1/problems` - Create new problem

### Submissions
- `POST /api/v1/submissions` - Submit code
- `GET /api/v1/submissions/:id` - Get submission details
- `GET /api/v1/submissions/match/:matchId` - Get match submissions

### Reports
- `GET /api/v1/reports/:matchId` - Get match report
- `GET /api/v1/reports/user/:userId` - Get user reports
- `GET /api/v1/reports/leaderboard` - Get leaderboard

### Skill Cards
- `GET /api/v1/skill-cards` - Get available cards
- `GET /api/v1/skill-cards/player/:playerId` - Get player cards
- `POST /api/v1/skill-cards/use` - Use skill card

### WebSocket
- `GET /ws/match/:roomId` - Join match room

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 🚀 Deployment

### Production Build
```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret
- `PORT`: Server port (default: 8080)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by competitive programming platforms
- Built with modern web technologies
- Designed for developer education and entertainment

## 📞 Support

If you have any questions or need help, please:
- Open an issue on GitHub
- Check the documentation
- Join our community discussions

---

**Happy Coding! 🎯**
