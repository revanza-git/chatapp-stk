# 🚀 Cloud Deployment Guide

This guide covers deploying your IT Security Policy Chatbot & Document Management System to various free cloud platforms.

## 📋 Pre-Deployment Checklist

1. **Push your code to GitHub**:
   ```bash
   git add .
   git commit -m "Ready for cloud deployment"
   git push origin main
   ```

2. **Generate a secure JWT secret**:
   ```bash
   # Generate a 32-character random string
   openssl rand -base64 32
   ```

## 🏆 Option 1: Railway (Recommended)

### Why Railway?
- ✅ $5/month free credits (sufficient for small apps)
- ✅ Native Docker support
- ✅ Auto-detects docker-compose.yml
- ✅ Built-in PostgreSQL
- ✅ Automatic SSL certificates

### Steps:
1. **Sign up**: Go to [railway.app](https://railway.app)
2. **Connect GitHub**: Link your repository
3. **Create Project**: Railway auto-detects your Docker setup
4. **Configure Environment**:
   - Set `JWT_SECRET` to your generated secret
   - Railway automatically handles database connection
   - Frontend will auto-detect backend URL
5. **Deploy**: Click "Deploy" - that's it!

### Railway Environment Variables:
```env
JWT_SECRET=your-generated-secret-here
AI_ENABLED=false
NODE_ENV=production
```

## 🎯 Option 2: Render

### Why Render?
- ✅ 750 hours/month free per service
- ✅ Docker support
- ✅ Free PostgreSQL (90 days)
- ✅ Easy GitHub integration

### Steps:
1. **Sign up**: Go to [render.com](https://render.com)
2. **Create Services**:
   - **Web Service** for Frontend (Dockerfile: `./Dockerfile.frontend`)
   - **Web Service** for Backend (Dockerfile: `./backend/Dockerfile`)
   - **PostgreSQL Database**
3. **Configure Environment Variables**:
   ```env
   # Frontend
   NODE_ENV=production
   NEXT_PUBLIC_API_URL=https://your-backend-url.onrender.com

   # Backend
   JWT_SECRET=your-generated-secret-here
   AI_ENABLED=false
   DATABASE_URL=postgresql://... (from database service)
   ```
4. **Deploy**: Services auto-deploy on push

## ⚡ Option 3: Fly.io

### Why Fly.io?
- ✅ 3 shared VMs free
- ✅ Docker-first platform
- ✅ Global edge deployment
- ✅ Great performance

### Steps:
1. **Install Fly CLI**:
   ```bash
   # Windows
   iwr https://fly.io/install.ps1 -useb | iex

   # macOS/Linux
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login & Initialize**:
   ```bash
   fly auth login
   fly launch --dockerfile Dockerfile.frontend
   ```

3. **Deploy Backend & Database**:
   ```bash
   cd backend
   fly launch --dockerfile Dockerfile
   fly postgres create --name chatapp-db
   ```

## 🔧 Environment Variables Setup

### Required Variables:
```env
# Security
JWT_SECRET=your-32-character-secret-here

# AI Features (keep disabled for free tier)
AI_ENABLED=false

# Database (cloud platforms provide this)
DATABASE_URL=postgresql://...

# Production mode
NODE_ENV=production
```

### How to Generate JWT Secret:
```bash
# Option 1: OpenSSL
openssl rand -base64 32

# Option 2: Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"

# Option 3: Online (use trusted source)
# Visit: https://generate-secret.vercel.app/32
```

## 📊 Platform Comparison

| Platform | Free Tier | Docker Support | Database | Ease of Use | Best For |
|----------|-----------|----------------|----------|-------------|----------|
| **Railway** | $5 credits/month | ⭐⭐⭐⭐⭐ | PostgreSQL included | ⭐⭐⭐⭐⭐ | **Beginners** |
| **Render** | 750hrs/service | ⭐⭐⭐⭐ | 90 days free | ⭐⭐⭐⭐ | Multiple services |
| **Fly.io** | 3 VMs | ⭐⭐⭐⭐⭐ | Paid after limit | ⭐⭐⭐ | Advanced users |

## 🚨 Important Notes

### For Free Tier Optimization:
1. **AI Disabled**: Saves resources (`AI_ENABLED=false`)
2. **File Uploads**: Limited by platform storage
3. **Database**: Use platform-provided PostgreSQL
4. **SSL**: Enable for database connections (`DB_SSLMODE=require`)

### Security Checklist:
- [ ] Change default JWT secret
- [ ] Use strong database passwords
- [ ] Enable SSL for database connections
- [ ] Review CORS settings for production domains

## 🔄 Continuous Deployment

All platforms support automatic deployment on git push:

```bash
# Make changes to your code
git add .
git commit -m "Update feature"
git push origin main
# 🚀 Automatic deployment triggers!
```

## 🆘 Troubleshooting

### Common Issues:

**Build Failures:**
- Check Dockerfile paths in configuration
- Ensure all dependencies are in package.json
- Verify environment variables are set

**Database Connection:**
- Enable SSL mode for cloud databases
- Check DATABASE_URL format
- Verify network policies allow connections

**Frontend-Backend Communication:**
- Update NEXT_PUBLIC_API_URL to backend URL
- Check CORS configuration
- Verify both services are running

## 📞 Need Help?

1. **Railway**: [docs.railway.app](https://docs.railway.app)
2. **Render**: [render.com/docs](https://render.com/docs)
3. **Fly.io**: [fly.io/docs](https://fly.io/docs)

---

**🎉 Ready to deploy? Start with Railway for the easiest experience!** 