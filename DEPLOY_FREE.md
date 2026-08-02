# 🚀 Free Deployment Guide

## Option 1: Render.com (Recommended - 100% Free)

### Steps:
1. Go to https://render.com and sign up with GitHub
2. Click "New +" > "Web Service"
3. Connect your GitHub: `simon141404-gif/shadow-chat-`
4. Select the `master` branch
5. Configure:
   - **Name**: `shadow-chat-backend`
   - **Environment**: `Go`
   - **Region**: `Oregon` (or closest to you)
   - **Build Command**: `go build -o bin/shadowchat-backend ./cmd/server`
   - **Start Command**: `./bin/shadowchat-backend`

6. Add Environment Variables:
   ```
   DATABASE_URL=postgres://user:pass@host:5432/db
   REDIS_URL=redis://host:6379
   JWT_SECRET=your_random_secret_key
   ALLOWED_ORIGIN=*
   ```

7. Click "Create Web Service"

### Free Database:
- Create free PostgreSQL: https://render.com/docs/new-postgres
- Create free Redis: https://render.com/docs/new-redis

---

## Option 2: Fly.io (Also Free)

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Launch
fly launch --image=your-docker-image
```

---

## Option 3: Railway (Sleeps after inactivity)

Railway has a free tier but sleeps after 5 minutes of inactivity. It wakes up on first request but takes 30-60 seconds.

---

## After Backend is Ready:

Update Android app to point to your backend:

File: `android/app/src/main/java/com/shadowchat/di/NetworkModule.kt`

Change the `BASE_URL` to your backend URL.
