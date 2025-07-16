# IT Security Policy Chatbot & Document Management System

A comprehensive, AI-powered platform designed to help organizations manage IT security policies, user onboarding, and document management. Built with Next.js 15, React 19, TypeScript for the frontend and Go with Gin for the backend, fully containerized with Docker.

🌟 **NEW**: Deploy to the cloud for **FREE** with Railway! See [Cloud Deployment Guide](DEPLOYMENT.md) for details.

## 🚀 Features

### 🔐 User Authentication & Authorization
- **Secure Login/Registration**: JWT-based authentication system
- **Role-Based Access Control**: User, Admin, HR, and IT Security roles
- **User Profile Management**: Update personal information and credentials
- **Protected Routes**: Secure access to sensitive features

### 🗄️ Document Management
- **File Upload System**: Support for PDF, DOC, DOCX files up to 10MB
- **Document Processing**: Automatic text extraction and indexing
- **Smart Downloads**: Original filename preservation with clean naming
- **Document Search**: Find documents by content and metadata
- **CRUD Operations**: Create, read, update, and delete documents

### 🤖 AI-Powered Chat Assistant
- **Security Onboarding Assistance**: Guided help with security awareness topics
- **Policy Search**: Intelligent document and policy retrieval
- **Context-Aware Responses**: AI understands uploaded documents
- **Multiple AI Backends**: Ollama (Llama 3.1) with Hugging Face fallback
- **Chat History**: Persistent conversation tracking

### 👥 Admin Dashboard
- **User Management**: Create, update, deactivate user accounts
- **Role Assignment**: Manage user permissions and access levels
- **Audit Logging**: Track all system activities and user actions
- **System Monitoring**: View platform usage and health metrics

### 🎨 Modern User Interface
- **Responsive Design**: Optimized for desktop, tablet, and mobile
- **Dashboard Layout**: Tabbed interface for different functions
- **File Upload**: Drag-and-drop with progress indicators
- **Real-time Updates**: Live chat and status updates

## 🏗️ System Architecture

The platform follows a modern microservices architecture with Docker containerization:

```mermaid
graph TB
    subgraph "Client Layer"
        U[👤 Users<br/>Web Browser]
    end

    subgraph "Docker Container Network"
        subgraph "Frontend Container<br/>:3000"
            FE[🖥️ Next.js Frontend<br/>React 19 + TypeScript<br/>Tailwind CSS + Radix UI]
        end

        subgraph "Backend Container<br/>:8080"
            BE[⚙️ Go Backend<br/>Gin Framework<br/>JWT Auth + RBAC<br/>File Processing]
        end

        subgraph "Database Container<br/>:5433"
            DB[(🗄️ PostgreSQL 15<br/>GORM ORM<br/>User Data<br/>Documents<br/>Audit Logs)]
        end

        subgraph "AI Container<br/>:11434"
            AI[🤖 Ollama Server<br/>Llama 3.1 Model<br/>AI Chat Engine]
        end
    end

    subgraph "Persistent Storage"
        PV1[📦 postgres-data<br/>Database Files]
        PV2[📦 ollama-data<br/>AI Models]
        PV3[📦 uploads-data<br/>Document Files]
    end

    subgraph "External Services"
        HF[🤗 Hugging Face API<br/>Fallback AI Service]
    end

    %% User interactions
    U -->|HTTPS Requests| FE
    FE -->|API Calls| BE

    %% Backend connections
    BE -->|SQL Queries| DB
    BE -->|AI Requests| AI
    BE -->|Fallback AI| HF
    BE -->|File Storage| PV3

    %% Database storage
    DB -.->|Persist| PV1

    %% AI model storage
    AI -.->|Persist| PV2

    %% User flows
    U -.->|"1. Login/Register"| FE
    FE -.->|"2. Authentication"| BE
    BE -.->|"3. User Verification"| DB
    U -.->|"4. Upload Documents"| FE
    FE -.->|"5. File Processing"| BE
    BE -.->|"6. Store Metadata"| DB
    U -.->|"7. Chat with AI"| FE
    FE -.->|"8. Process Chat"| BE
    BE -.->|"9. AI Response"| AI

    %% Styling
    classDef containerStyle fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef storageStyle fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef externalStyle fill:#fff3e0,stroke:#e65100,stroke-width:2px
    classDef userStyle fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px

    class FE,BE,DB,AI containerStyle
    class PV1,PV2,PV3 storageStyle
    class HF externalStyle
    class U userStyle
```

### 🔄 Data Flow
1. **User Authentication**: Users access the frontend and authenticate via JWT tokens
2. **Document Management**: Files are uploaded, processed, and stored with metadata in PostgreSQL
3. **AI Integration**: Chat requests are processed by Ollama with Hugging Face as fallback
4. **Audit Logging**: All user activities are tracked in the database for compliance
5. **Persistent Storage**: Data persists across container restarts using Docker volumes

### 🌐 Network Architecture
- **Frontend (Port 3000)**: Serves the React application with SSR/SSG capabilities
- **Backend (Port 8080)**: RESTful API with authentication, file processing, and AI integration
- **Database (Port 5433)**: PostgreSQL with ACID compliance and backup capabilities
- **AI Service (Port 11434)**: Ollama server for local AI inference with model caching

## 🛠️ Technology Stack

### Frontend
- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS v4
- **UI Components**: Radix UI primitives with custom styling
- **State Management**: TanStack React Query v5
- **Authentication**: JWT with React Context
- **File Upload**: React Dropzone
- **Icons**: Lucide React

### Backend
- **Language**: Go 1.24+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens with middleware
- **File Processing**: PDF, DOC, DOCX text extraction
- **AI Integration**: 
  - Ollama API (Llama 3.1)
  - Hugging Face Inference API
  - Fallback response system
- **CORS**: Multi-origin support for Docker/local development

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL 15 with persistent volumes
- **AI Model Server**: Ollama container with GPU support
- **File Storage**: Volume-mounted uploads with backup support
- **Health Checks**: Container health monitoring

## 📦 Quick Start with Docker (Recommended)

### Prerequisites
- [Docker Desktop](https://docs.docker.com/get-docker/) installed
- At least 4GB of available RAM

### 🚀 One-Click Setup

#### For Windows:
```cmd
# Run the startup script
docker-start.bat
```

#### For Linux/macOS:
```bash
# Make the script executable and run
chmod +x docker-start.sh
./docker-start.sh
```

#### Manual Docker Setup:
```bash
# Clone the repository
git clone <repository-url>
cd chatapp-stk

# Copy environment file
cp .env.example .env

# Start all services
docker-compose up --build -d

# Setup AI model (optional)
docker-compose exec ollama ollama pull llama2
```

### 🌐 Access Your Application
Once running, visit:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5433
- **Ollama AI**: http://localhost:11434

## 🏗️ Project Structure

```
chatapp-stk/
├── app/                          # Next.js App Router
│   ├── dashboard/
│   │   └── page.tsx             # Main dashboard with document management
│   ├── login/
│   │   └── page.tsx             # User login page
│   ├── register/
│   │   └── page.tsx             # User registration page
│   ├── layout.tsx               # Root layout with providers
│   ├── page.tsx                 # Landing page
│   ├── providers.tsx            # React Query & Auth providers
│   └── globals.css              # Global Tailwind styles
│
├── components/                   # React components
│   ├── admin/
│   │   ├── AuditLogs.tsx        # System audit logging
│   │   └── UserManagement.tsx   # User administration
│   ├── auth/
│   │   ├── ProtectedRoute.tsx   # Route protection wrapper
│   │   └── UserProfile.tsx      # User profile management
│   ├── chatbot/
│   │   ├── ChatInterface.tsx    # Main chat component
│   │   └── ChatMessage.tsx      # Individual message component
│   └── ui/                      # Reusable UI components
│       ├── badge.tsx
│       ├── button.tsx
│       ├── card.tsx
│       ├── dialog.tsx
│       ├── file-upload.tsx      # Drag-and-drop file upload
│       ├── input.tsx
│       ├── scroll-area.tsx
│       └── separator.tsx
│
├── lib/                         # Utility libraries
│   ├── api.ts                   # API client functions
│   ├── auth.tsx                 # Authentication context & utilities
│   ├── types.ts                 # TypeScript type definitions
│   └── utils.ts                 # Utility functions
│
├── backend/                     # Go backend server
│   ├── main.go                  # Main server with all API endpoints
│   ├── go.mod                   # Go dependencies (1.24+)
│   ├── go.sum                   # Go dependency checksums
│   ├── init-db.sql              # Database initialization
│   ├── uploads/                 # File storage directory
│   ├── Dockerfile               # Backend container configuration
│   ├── docker-compose.yml       # Legacy compose file
│   ├── colab_ollama_setup.ipynb # Google Colab setup notebook
│   └── setup_colab_integration.md # AI setup documentation
│
├── Docker Configuration
│   ├── docker-compose.yml       # Multi-service orchestration
│   ├── Dockerfile.frontend      # Frontend container
│   ├── docker-start.bat         # Windows startup script
│   ├── docker-start.sh          # Linux/macOS startup script
│   ├── DOCKER.md               # Comprehensive Docker guide
│   ├── .dockerignore           # Frontend build optimization
│   ├── backend/.dockerignore   # Backend build optimization
│   └── .env.example            # Environment configuration template
│
└── Configuration Files
    ├── next.config.ts           # Next.js configuration
    ├── package.json             # Frontend dependencies
    ├── tsconfig.json            # TypeScript configuration
    ├── components.json          # UI components config
    └── postcss.config.mjs       # PostCSS configuration
```

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/login` | User login |
| `POST` | `/api/auth/register` | User registration |

### Chat & AI
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chat` | Send chat message and get AI response |

### Document Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/documents` | List all documents (paginated) |
| `POST` | `/api/documents` | Upload new document |
| `GET` | `/api/documents/:id` | Get document details |
| `PUT` | `/api/documents/:id` | Update document |
| `DELETE` | `/api/documents/:id` | Delete document |
| `GET` | `/api/documents/:id/download` | Download original file |

### User Management (Admin)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | List all users |
| `POST` | `/api/users` | Create new user |
| `PUT` | `/api/users/:id` | Update user |
| `DELETE` | `/api/users/:id` | Deactivate user |

### Audit & Monitoring
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/audit-logs` | View system audit logs |
| `GET` | `/api/health` | Backend health check |

## 🐳 Docker Services

The application runs as a multi-container setup:

| Service | Port | Description |
|---------|------|-------------|
| **Frontend** | 3000 | Next.js React application |
| **Backend** | 8080 | Go API server |
| **PostgreSQL** | 5433 | Database with persistent storage |
| **Ollama** | 11434 | AI model server (DISABLED by default) |

### ⚡ AI Features Status
**AI features are currently DISABLED** for faster deployment and cloud hosting compatibility:
- ✅ **Chat Interface**: Works with intelligent mock responses
- ✅ **Document Management**: Full functionality available  
- ✅ **User Authentication**: Complete system operational
- 🔄 **AI Integration**: Can be enabled when needed (see DOCKER.md)

### Data Persistence
- **postgres-data**: Database files and user data
- **ollama-data**: Downloaded AI models
- **uploads-data**: Uploaded documents and files

## 👥 User Roles & Permissions

### User (Default)
- ✅ Access chat interface
- ✅ Upload and download documents
- ✅ View own profile
- ❌ User management
- ❌ Audit logs

### HR
- ✅ All User permissions
- ✅ View all documents
- ✅ User onboarding assistance
- ❌ System administration

### IT Security
- ✅ All HR permissions
- ✅ Security policy management
- ✅ Advanced chat features
- ❌ User account management

### Admin
- ✅ All permissions
- ✅ User management (create, update, deactivate)
- ✅ Audit log access
- ✅ System monitoring

## 🚀 Development Setup

### Local Development (without Docker)

1. **Prerequisites**:
   - Node.js 18+
   - Go 1.24+
   - PostgreSQL 15+

2. **Database Setup**:
   ```bash
   # Create database
   createdb chatbot_db
   psql chatbot_db < backend/init-db.sql
   ```

3. **Backend Setup**:
   ```bash
   cd backend
   go mod tidy
   
   # Set environment variables
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=chatbot_user
   export DB_PASSWORD=chatbot_password
   export DB_NAME=chatbot_db
   
   go run main.go
   ```

4. **Frontend Setup**:
   ```bash
   npm install
   npm run dev
   ```

### Environment Configuration

#### Local Development Environment (.env)
```bash
# Frontend Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080

# Backend Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=chatbot_user
DB_PASSWORD=chatbot_password
DB_NAME=chatbot_db
DB_SSLMODE=disable

# Security
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# AI Configuration (optional)
AI_ENABLED=false
OLLAMA_HOST=ollama:11434
HF_TOKEN=your-hugging-face-token-here

# Server Configuration
PORT=8080
```

#### Cloud Deployment Environment Variables
For cloud platforms like Railway, set these environment variables directly in your platform's dashboard:

**Frontend Service:**
```bash
NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
NODE_ENV=production
```

**Backend Service:**
```bash
# Database (automatically provided by Railway PostgreSQL)
DATABASE_URL=postgresql://username:password@host:port/database

# Security (generate secure values)
JWT_SECRET=your-super-secure-random-jwt-secret-minimum-32-characters

# AI Configuration
AI_ENABLED=false
OLLAMA_HOST=
HF_TOKEN=

# Server
PORT=8080

# CORS
FRONTEND_URL=https://your-frontend-url.railway.app
```

## 🤖 AI Integration

The chatbot uses a robust multi-tier AI system:

1. **Ollama (Primary)**: Llama 3.1 model running in Docker container
2. **Hugging Face (Fallback)**: Cloud-based inference API
3. **Mock Responses (Development)**: Predefined responses for offline development

### Document Processing
- **PDF Text Extraction**: Full-text indexing of PDF documents
- **Microsoft Office**: DOC and DOCX file processing
- **Context Integration**: AI responses include relevant document content
- **Smart Search**: Semantic search across document content

## 🔒 Security Features

### Authentication & Authorization
- **JWT Tokens**: Secure, stateless authentication
- **Password Hashing**: bcrypt with salt rounds
- **Role-Based Access**: Granular permission system
- **Session Management**: Token refresh and logout

### Data Protection
- **Input Validation**: Comprehensive request sanitization
- **File Upload Security**: Type validation, size limits, virus scanning
- **CORS Configuration**: Secure cross-origin resource sharing
- **SQL Injection Prevention**: Parameterized queries with GORM

### Audit & Monitoring
- **Activity Logging**: All user actions tracked
- **Security Events**: Login attempts, permission changes
- **System Monitoring**: Health checks and performance metrics
- **Data Backup**: Automated database and file backups

## 🚀 Local Development with Docker

### Docker Management
```bash
# View service status
docker-compose ps

# View logs
docker-compose logs -f [service-name]

# Restart services
docker-compose restart [service-name]

# Reset everything
docker-compose down -v
docker-compose up --build -d
```

### Health Checks
```bash
# Backend health
curl http://localhost:8080/api/health

# Frontend
curl http://localhost:3000

# Database connection
docker-compose exec postgres psql -U chatbot_user -d chatbot_db -c "SELECT 1;"
```

## ☁️ Cloud Deployment (FREE)

Deploy your application to Railway for **FREE**! The platform is optimized for cloud deployment with automatic scaling and SSL certificates.

### 🚀 **Quick Deploy with Railway (Recommended)**

#### Prerequisites
- GitHub account with your code pushed
- Railway account (free signup)

#### Deployment Steps
1. **Push your code to GitHub**:
   ```bash
   git add .
   git commit -m "Ready for Railway deployment"
   git push origin main
   ```

2. **Deploy on Railway**:
   - Visit [railway.app](https://railway.app) and sign up
   - Click "New Project" → "Deploy from GitHub repo"
   - Select your repository
   - Railway will automatically detect your Docker setup

3. **Configure Services**:
   Railway will create two services:
   - **Frontend**: Next.js application
   - **Backend**: Go API server
   - **Database**: PostgreSQL (automatically provisioned)

4. **Set Environment Variables**:
   
   **Frontend Service Variables:**
   ```env
   NEXT_PUBLIC_API_URL=https://your-backend-service.railway.app
   NODE_ENV=production
   ```
   
   **Backend Service Variables:**
   ```env
   JWT_SECRET=VElJSUlJSUlJSUlJSUlJSUlJSUlJSUlJSUlJSUlJSUk=
   AI_ENABLED=false
   FRONTEND_URL=https://your-frontend-service.railway.app
   PORT=8080
   ```

5. **Deploy**: Railway automatically builds and deploys! 🎉

#### After Deployment
- Your frontend will be available at: `https://your-frontend-service.railway.app`
- Your backend API will be at: `https://your-backend-service.railway.app`
- Database is automatically connected and managed

### 🔧 **Alternative Cloud Platforms**

| Platform | Free Tier | Database | Difficulty | Best For |
|----------|-----------|----------|------------|----------|
| **Railway** ⭐ | $5 credits/month | PostgreSQL included | ⭐ Easy | Docker apps |
| **Render** | 750hrs/service | 90 days free | ⭐⭐ Medium | Multiple services |
| **Fly.io** | 3 shared VMs | Paid after limit | ⭐⭐⭐ Advanced | Performance |

### 🛡️ **Cloud Optimizations**
- ✅ **AI Disabled by default**: Saves compute resources for free tier
- ✅ **Multi-stage Docker builds**: Efficient containers
- ✅ **Automatic SSL/HTTPS**: Secure connections included
- ✅ **Database included**: PostgreSQL automatically provisioned
- ✅ **Auto-scaling**: Handles traffic spikes automatically

**📖 For detailed deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md)**

## 🛠️ Troubleshooting

### Common Issues

#### Local Development
1. **Port conflicts**: Update ports in docker-compose.yml
2. **Build failures**: Check Docker logs with `docker-compose logs -f`
3. **Database connection**: Verify PostgreSQL container is running
4. **File uploads failing**: Check volume permissions and disk space

#### Cloud Deployment
1. **Frontend calling localhost**: Ensure `NEXT_PUBLIC_API_URL` is set correctly in cloud platform
2. **Authentication errors**: Verify `JWT_SECRET` is set and matches between frontend/backend
3. **CORS issues**: Ensure `FRONTEND_URL` is set in backend environment variables
4. **Database connection**: Check if DATABASE_URL is automatically provided by your cloud platform

### Debug Commands
```bash
# Local development
docker-compose logs -f [service-name]
docker-compose ps
curl http://localhost:8080/api/health

# Check environment variables (local)
docker-compose exec frontend printenv | grep NEXT_PUBLIC
docker-compose exec backend printenv | grep JWT_SECRET
```

## 📄 Documentation

- **[Docker Setup Guide](DOCKER.md)**: Comprehensive Docker documentation
- **[Cloud Deployment Guide](DEPLOYMENT.md)**: Railway and other cloud platform deployment
- **Environment Setup**: Use `.env.example` as template for local development

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and test thoroughly
4. Commit with conventional commit messages
5. Push to your branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Guidelines
- Follow TypeScript/Go best practices
- Write tests for new features
- Update documentation
- Ensure Docker builds work
- Test with different user roles

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

### Getting Help
1. Check the [Docker Guide](DOCKER.md) for setup issues
2. Review logs with `docker-compose logs -f`
3. Open an issue in the repository
4. Contact your IT team for security-specific questions

### System Requirements
- **Minimum**: 4GB RAM, 2 CPU cores, 10GB disk space
- **Recommended**: 8GB RAM, 4 CPU cores, 50GB disk space
- **Docker**: Version 20.10+ with Compose V2

---

**Note**: This platform is designed for IT security policy management and user onboarding. For production deployments, always use secure JWT secrets and follow your organization's security guidelines.

Built with ❤️ using Next.js, Go, and Docker. Optimized for Railway deployment.

# Railway redeploy trigger