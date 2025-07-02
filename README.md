# IT Security Policy Chatbot

A modern, AI-powered chatbot application designed to help users navigate IT security policies and onboarding procedures. Built with Next.js 15, React 19, TypeScript for the frontend and Go with Gin for the backend.

## 🚀 Features

- **Security Onboarding Assistance**: Get guided help with security awareness topics including passwords, VPN setup, email security, and data protection policies
- **Policy Search**: Quickly find and access specific security policies, procedures, and guidelines relevant to your role
- **Instant Answers**: Get immediate responses with relevant policy documents, best practices, and actionable security guidance
- **AI-Powered Responses**: Uses Ollama (Llama 3.1) or Hugging Face APIs for intelligent, context-aware responses
- **Modern UI**: Beautiful, responsive interface built with Tailwind CSS and Radix UI components
- **Real-time Chat**: Interactive chat interface for seamless user experience

## 🛠️ Technology Stack

### Frontend

- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **UI Components**: Radix UI + Custom components
- **State Management**: TanStack React Query
- **Icons**: Lucide React
- **Development**: ESLint, Turbopack

### Backend

- **Language**: Go
- **Framework**: Gin (HTTP web framework)
- **AI Integration**: Ollama API + Hugging Face API
- **CORS**: Enabled for frontend integration
- **Docker**: Containerized deployment

## 📦 Installation & Setup

### Prerequisites

- Node.js 18+
- Go 1.21+
- Docker (optional)

### Frontend Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd chatbotawareness
```

2. Install frontend dependencies:

```bash
npm install
```

3. Set up environment variables (optional):

For production deployments, you can set these environment variables:

**Frontend (.env.local):**

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

**Backend (system environment or .env file):**

```bash
OLLAMA_URL=http://your-ollama-server:11434
HF_TOKEN=your_hugging_face_token
```

4. Run the frontend development server:

```bash
npm run dev
```

The frontend will be available at [http://localhost:3000](http://localhost:3000)

### Backend Setup

1. Navigate to the backend directory:

```bash
cd backend
```

2. Install Go dependencies:

```bash
go mod tidy
```

3. Set up environment variables (optional):

```bash
# For Ollama integration (preferred)
export OLLAMA_URL=http://your-ollama-server:11434

# For Hugging Face fallback
export HF_TOKEN=your_hugging_face_token
```

4. Run the backend server:

```bash
go run main.go
```

The backend API will be available at [http://localhost:8080](http://localhost:8080)

### Docker Setup (Alternative)

You can also run the backend using Docker:

```bash
cd backend
docker-compose up --build
```

## 🏗️ Project Structure

```
├── app/                    # Next.js App Router pages
│   ├── layout.tsx         # Root layout component
│   ├── page.tsx           # Home page with chat interface
│   ├── providers.tsx      # React Query and other providers
│   └── globals.css        # Global styles
├── components/            # Reusable React components
│   ├── chatbot/          # Chatbot-specific components
│   │   ├── ChatInterface.tsx
│   │   └── ChatMessage.tsx
│   └── ui/               # UI components (cards, buttons, etc.)
├── lib/                  # Utility functions and configurations
├── public/               # Static assets
├── backend/              # Go backend server
│   ├── main.go          # Main server file with API endpoints
│   ├── go.mod           # Go module dependencies
│   ├── Dockerfile       # Docker container configuration
│   ├── docker-compose.yml
│   └── setup_colab_integration.md
└── package.json         # Frontend dependencies
```

## 🔌 API Endpoints

### Backend API (Port 8080)

- `POST /api/chat` - Main chat endpoint

  - Body: `{"message": "string", "type": "onboarding|policy_search"}`
  - Response: `{"response": "string", "type": "string", "policy_files": [...]}`

- `GET /api/policies` - Get all available policies

  - Response: Array of policy objects

- `GET /health` - Health check endpoint

## 🤖 AI Integration

The chatbot supports multiple AI backends:

1. **Ollama (Preferred)**: Local Llama 3.1 model via Google Colab or local setup
2. **Hugging Face API**: Fallback using DialoGPT model
3. **Mock Responses**: Default fallback for development

### Setting up Ollama

See `backend/setup_colab_integration.md` for detailed instructions on setting up Ollama with Google Colab.

## 🚀 Available Scripts

### Frontend

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build the application for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint for code quality

### Backend

- `go run main.go` - Start the Go server
- `go mod tidy` - Clean up dependencies
- `docker-compose up` - Run with Docker

## 🎨 UI Components

The application uses a custom UI component library built with:

- **Radix UI**: Accessible, unstyled components
- **Tailwind CSS**: Utility-first CSS framework
- **Class Variance Authority**: Type-safe component variants
- **Lucide React**: Beautiful, customizable icons

## 🔧 Configuration

### Frontend Configuration

- `next.config.ts` - Next.js configuration
- `tailwind.config.js` - Tailwind CSS configuration
- `tsconfig.json` - TypeScript configuration
- `eslint.config.mjs` - ESLint configuration
- `components.json` - UI components configuration

### Backend Configuration

- `go.mod` - Go module dependencies
- `docker-compose.yml` - Docker services configuration
- Environment variables for AI API integration

## 📱 Responsive Design

The application is fully responsive and optimized for:

- Desktop computers
- Tablets
- Mobile devices

## 🔒 Security Features

- **Policy Database**: Built-in security policies for common topics
- **Context-Aware Responses**: AI responses tailored to security onboarding
- **CORS Protection**: Configured for secure frontend-backend communication
- **Input Validation**: Proper request validation and sanitization

## 🚀 Deployment

### Frontend Deployment

The Next.js app can be deployed to:

- Vercel (recommended)
- Netlify
- Any Node.js hosting platform

### Backend Deployment

The Go backend can be deployed to:

- Docker containers
- Cloud platforms (AWS, GCP, Azure)
- Traditional VPS servers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- Check the documentation
- Open an issue in the repository
- Contact your IT team for security-specific questions

---

**Note**: This chatbot is designed to assist with IT security policies and onboarding. For critical security decisions, always consult with your organization's IT security team.
