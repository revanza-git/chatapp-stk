# IT Security Policy Chatbot

A modern, AI-powered chatbot application designed to help users navigate IT security policies and onboarding procedures. Built with Next.js 15, React 19, and TypeScript.

## 🚀 Features

- **Security Onboarding Assistance**: Get guided help with security awareness topics including passwords, VPN setup, email security, and data protection policies
- **Policy Search**: Quickly find and access specific security policies, procedures, and guidelines relevant to your role
- **Instant Answers**: Get immediate responses with relevant policy documents, best practices, and actionable security guidance
- **Modern UI**: Beautiful, responsive interface built with Tailwind CSS and Radix UI components
- **Real-time Chat**: Interactive chat interface for seamless user experience

## 🛠️ Technology Stack

- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **UI Components**: Radix UI + Custom components
- **State Management**: TanStack React Query
- **Icons**: Lucide React
- **Development**: ESLint, Turbopack

## 📦 Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd frontend
```

2. Install dependencies:

```bash
npm install
# or
yarn install
# or
pnpm install
```

3. Run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser to see the application.

## 🏗️ Project Structure

```
├── app/                    # Next.js App Router pages
│   ├── layout.tsx         # Root layout component
│   ├── page.tsx           # Home page with chat interface
│   ├── providers.tsx      # React Query and other providers
│   └── globals.css        # Global styles
├── components/            # Reusable components
│   ├── chatbot/          # Chatbot-specific components
│   └── ui/               # UI components (cards, buttons, etc.)
├── lib/                  # Utility functions and configurations
└── public/               # Static assets
```

## 🚀 Available Scripts

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build the application for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint for code quality

## 🎨 UI Components

The application uses a custom UI component library built with:

- **Radix UI**: Accessible, unstyled components
- **Tailwind CSS**: Utility-first CSS framework
- **Class Variance Authority**: Type-safe component variants
- **Lucide React**: Beautiful, customizable icons

## 🔧 Configuration

The project includes several configuration files:

- `next.config.ts` - Next.js configuration
- `tailwind.config.js` - Tailwind CSS configuration
- `tsconfig.json` - TypeScript configuration
- `eslint.config.mjs` - ESLint configuration
- `components.json` - UI components configuration

## 📱 Responsive Design

The application is fully responsive and optimized for:

- Desktop computers
- Tablets
- Mobile devices

## 🔒 Security Considerations

- Always verify critical security information with your IT team
- The chatbot provides guidance but should not replace official security protocols
- Sensitive information should be handled according to your organization's policies

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
