# IT Security Policy Chatbot

A modern, AI-powered chatbot application designed to help users navigate IT security policies and onboarding procedures. Built with Next.js 15, React 19, TypeScript for the frontend and Go with Gin for the backend.

## 🚀 Features

- **Security Onboarding Assistance**: Get guided help with security awareness topics including passwords, VPN setup, email security, and data protection policies
- **Policy Search**: Quickly find and access specific security policies, procedures, and guidelines relevant to your role
- **Intelligent Chat Modes**: Choose between "Security Onboarding" and "Policy Search" for targeted assistance
- **AI-Powered Responses**: Uses Ollama (Llama 3.1) or Hugging Face APIs for intelligent, context-aware responses
- **Modern UI**: Beautiful, responsive interface built with Tailwind CSS v4 and Radix UI components
- **Real-time Chat**: Interactive chat interface with message history and typing indicators
- **Robust Fallback System**: Automatic fallback from Ollama → Hugging Face → Mock responses

## 🛠️ Technology Stack

### Frontend
- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS v4
- **UI Components**: Radix UI primitives with custom styling
- **State Management**: TanStack React Query v5
- **Icons**: Lucide React
- **Development**: ESLint 9, Turbopack

### Backend
- **Language**: Go 1.23+
- **Framework**: Gin (HTTP web framework)
- **AI Integration**: 
  - Ollama API (via Google Colab or local)
  - Hugging Face Inference API (DialoGPT-medium)
  - Mock responses for development
- **CORS**: Configured for cross-origin requests
- **Docker**: Full containerization support

## 📦 Installation & Setup

### Prerequisites
- Node.js 18+
- Go 1.23+
- Docker (optional)

### Quick Start

1. **Clone the repository**:
```bash
git clone <repository-url>
cd chatapp-stk
```

2. **Install frontend dependencies**:
```bash
npm install
```

3. **Start the frontend** (in one terminal):
```bash
npm run dev
```

4. **Start the backend** (in another terminal):
```bash
cd backend
go mod tidy
go run main.go
```

The application will be available at:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend API: [http://localhost:8080](http://localhost:8080)

### Environment Configuration

#### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

#### Backend Environment Variables
```bash
# Primary AI Backend (Ollama via Google Colab)
OLLAMA_URL=https://your-ngrok-url.ngrok.io

# Fallback AI Backend (Hugging Face)
HF_TOKEN=hf_your_hugging_face_token_here
```

### AI Backend Setup

#### Option 1: Google Colab + Ollama (Recommended)
1. Upload `backend/colab_ollama_setup.ipynb` to Google Colab
2. Run all cells to setup Ollama with Llama 3.1
3. Copy the ngrok URL and set it as `OLLAMA_URL`

See `backend/setup_colab_integration.md` for detailed instructions.

#### Option 2: Local Docker Setup
```bash
cd backend
docker-compose up --build
```

#### Option 3: Hugging Face Only
Set only the `HF_TOKEN` environment variable.

## 🏗️ Project Structure

```
chatapp-stk/
├── app/                          # Next.js App Router
│   ├── layout.tsx               # Root layout with providers
│   ├── page.tsx                 # Main chat page
│   ├── providers.tsx            # React Query provider
│   └── globals.css              # Global Tailwind styles
│
├── components/                   # React components
│   ├── chatbot/
│   │   ├── ChatInterface.tsx    # Main chat component
│   │   └── ChatMessage.tsx      # Individual message component
│   └── ui/                      # Reusable UI components
│       ├── badge.tsx
│       ├── button.tsx
│       ├── card.tsx
│       ├── dialog.tsx
│       ├── input.tsx
│       ├── scroll-area.tsx
│       └── separator.tsx
│
├── lib/                         # Utility libraries
│   ├── api.ts                   # API client functions
│   ├── types.ts                 # TypeScript type definitions
│   └── utils.ts                 # Utility functions
│
├── backend/                     # Go backend server
│   ├── main.go                  # Main server with API endpoints
│   ├── go.mod                   # Go dependencies
│   ├── go.sum                   # Go dependency checksums
│   ├── Dockerfile               # Container configuration
│   ├── docker-compose.yml       # Multi-service setup
│   ├── colab_ollama_setup.ipynb # Google Colab setup notebook
│   └── setup_colab_integration.md # Setup documentation
│
└── Configuration files
    ├── next.config.ts           # Next.js configuration
    ├── package.json             # Frontend dependencies
    ├── tsconfig.json            # TypeScript configuration
    ├── components.json          # UI components config
    └── postcss.config.mjs       # PostCSS configuration
```

## 🔌 API Endpoints

### Backend API (Port 8080)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chat` | Send chat message and get AI response |
| `GET` | `/api/policies` | Retrieve all available security policies |
| `GET` | `/health` | Backend health check |

#### Chat API Request/Response
```typescript
// Request
{
  "message": "What is our password policy?",
  "type": "onboarding" | "policy_search"
}

// Response
{
  "response": "Password policy details...",
  "type": "onboarding",
  "policy_files": [...]  // Optional policy references
}
```

## 🤖 AI Integration

The chatbot uses a robust three-tier AI system:

1. **Ollama (Primary)**: Llama 3.1 8B model via Google Colab or local deployment
2. **Hugging Face (Fallback)**: DialoGPT-medium for basic conversations
3. **Mock Responses (Development)**: Predefined responses for offline development

### Built-in Policy Database

The backend includes a comprehensive policy database covering:
- **Authentication**: Password policies, MFA requirements
- **Data Protection**: Classification, encryption standards
- **Remote Work**: VPN usage, device security
- **Incident Response**: Reporting procedures, escalation matrix

## 🚀 Available Scripts

### Frontend
```bash
npm run dev      # Development server with Turbopack
npm run build    # Production build
npm run start    # Production server
npm run lint     # ESLint code quality check
```

### Backend
```bash
go run main.go   # Start development server
go mod tidy      # Clean up dependencies
go build         # Build binary
```

### Docker
```bash
docker-compose up --build    # Full stack with Ollama
docker-compose down          # Stop all services
```

## 🎨 UI Features

- **Responsive Design**: Optimized for desktop, tablet, and mobile
- **Modern Typography**: Geist Sans and Geist Mono fonts
- **Accessible Components**: Built with Radix UI primitives
- **Smooth Animations**: Tailwind CSS animations
- **Dark/Light Mode Ready**: CSS variables for theming

## 🔧 Configuration Files

| File | Purpose |
|------|---------|
| `next.config.ts` | Next.js framework configuration |
| `tailwind.config.js` | Tailwind CSS customization |
| `tsconfig.json` | TypeScript compiler options |
| `eslint.config.mjs` | Code linting rules |
| `components.json` | UI component library settings |
| `postcss.config.mjs` | PostCSS processing |

## 🚀 Deployment

### Frontend Deployment
Deploy to any of these platforms:
- **Vercel** (recommended for Next.js)
- **Netlify**
- **AWS Amplify**
- **Any Node.js hosting**

### Backend Deployment
Deploy using:
- **Docker containers** (recommended)
- **Google Cloud Run**
- **AWS ECS/Fargate**
- **Traditional VPS**

### Production Environment Variables
```bash
# Frontend
NEXT_PUBLIC_API_URL=https://your-api-domain.com/api

# Backend
OLLAMA_URL=https://your-ollama-instance.com
HF_TOKEN=your_production_hf_token
PORT=8080
```

## 🔒 Security Features

- **Input Validation**: Comprehensive request validation and sanitization
- **CORS Configuration**: Properly configured cross-origin resource sharing
- **Environment Variables**: Sensitive data stored in environment variables
- **Error Handling**: Graceful error handling with fallback responses
- **Policy Database**: Built-in security policy knowledge base

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow TypeScript best practices
- Use conventional commit messages
- Ensure all tests pass
- Update documentation for new features

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support & Troubleshooting

### Common Issues

**Frontend won't start:**
- Ensure Node.js 18+ is installed
- Run `npm install` to install dependencies
- Check for port conflicts (3000)

**Backend connection errors:**
- Verify Go 1.23+ is installed
- Check backend is running on port 8080
- Verify API URL in frontend configuration

**AI responses not working:**
- Check environment variables are set
- Verify Ollama/Colab setup (see setup guide)
- Check console logs for API errors

### Getting Help
- Check the documentation
- Review the setup guides in `backend/`
- Open an issue in the repository
- Contact your IT team for security-specific questions

---

**Note**: This chatbot is designed to assist with IT security policies and onboarding. For critical security decisions, always consult with your organization's IT security team.
