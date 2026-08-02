# ShadowChat

A privacy-focused end-to-end encrypted messaging application with a Go backend and Android client.

## Features

- 🔒 **End-to-End Encryption** - All messages are encrypted using AES-256-GCM
- 👥 **Anonymous Identity** - Create accounts without phone number or email
- 💬 **Real-time Messaging** - WebSocket-based instant messaging
- 📱 **Android Client** - Modern Jetpack Compose UI
- 🔑 **Secure Storage** - SQLCipher encrypted local database
- 🚀 **Self-Hosted** - Run your own backend
- 🏗️ **CI/CD** - Automatic builds on every push

## Quick Start - Download APK

### Pre-built APKs

Download the latest debug APK from GitHub Actions:

1. Go to [GitHub Actions](https://github.com/simon141404-gif/shadow-chat-/actions)
2. Click on the latest workflow run
3. Download the `debug-apk` artifact

### Build from Source

```bash
# Clone the repository
git clone https://github.com/simon141404-gif/shadow-chat-.git
cd shadow-chat-

# Build debug APK
cd android
./gradlew assembleDebug

# APK will be at: app/build/outputs/apk/debug/app-debug.apk
```

## Tech Stack

### Backend
- **Go 1.21** with Gin web framework
- **PostgreSQL** - Primary database
- **Redis** - Session management and caching
- **JWT** - Authentication tokens
- **WebSocket** - Real-time messaging

### Android
- **Kotlin** with Jetpack Compose
- **Material 3** design
- **Room + SQLCipher** - Encrypted local database
- **Retrofit** - Network layer
- **Android Keystore** - Key management
- **CameraX + ML Kit** - QR code scanning

## Project Structure

```
shadow-chat-/
├── cmd/server/           # Backend entry point
├── internal/
│   ├── config/          # Configuration
│   ├── db/               # Database connections
│   ├── model/            # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── http/             # HTTP handlers & routes
├── migrations/           # SQL migrations
└── android/              # Android application
    ├── app/              # Main app module
    └── core/             # Core libraries
        ├── common/
        ├── crypto/
        ├── database/
        ├── network/
        ├── testing/
        └── ui/
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Android Studio Ladybug+
- Gradle 8.2+

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/simon141404-gif/shadow-chat-.git
cd shadow-chat-
```

2. Configure environment variables:
```bash
export POSTGRES_URL="postgres://user:pass@localhost:5432/shadowchat?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key"
export UPLOAD_DIR="/data/uploads"
export ALLOWED_ORIGIN="*"
export PORT="8080"
export ENV="development"
```

3. Run database migrations:
```bash
psql $POSTGRES_URL -f migrations/001_init.sql
```

4. Start the server:
```bash
cd cmd/server
go run main.go
```

The API will be available at `http://localhost:8080`

### Android Setup

1. Open the `android/` directory in Android Studio

2. Configure the API base URL in `core/network/NetworkModule.kt`:
```kotlin
private const val BASE_URL = "http://10.0.2.2:8080/" // Android emulator
```

3. Build the debug APK:
```bash
cd android
./gradlew assembleDebug
```

The APK will be at `app/build/outputs/apk/debug/app-debug.apk`

### Release Build

For release builds, set up signing:

1. Create `keystore.properties`:
```properties
storeFile=/path/to/keystore.jks
storePassword=your_password
keyAlias=your_alias
keyPassword=your_key_password
```

2. Build release:
```bash
./gradlew assembleRelease
```

## API Endpoints

### Authentication
- `POST /v1/auth/anonymous` - Create anonymous identity
- `POST /v1/auth/refresh` - Refresh session

### Chats (Requires Auth)
- `GET /v1/chats` - List chats
- `POST /v1/chats` - Create chat
- `GET /v1/chats/:chatId` - Get chat
- `GET /v1/chats/:chatId/messages` - List messages
- `POST /v1/chats/:chatId/messages` - Send message

### Messages (Requires Auth)
- `PATCH /v1/messages/:messageId` - Edit message
- `DELETE /v1/messages/:messageId` - Delete message

### Contacts (Requires Auth)
- `GET /v1/contacts` - List contacts
- `POST /v1/contacts/share` - Share contact

### Profile (Requires Auth)
- `GET /v1/profile` - Get profile
- `PATCH /v1/profile` - Update profile

### Uploads (Requires Auth)
- `POST /v1/uploads` - Create upload
- `GET /v1/uploads/:uploadId` - Get upload

### WebSocket
- `GET /v1/ws?chatId={chatId}` - Real-time messaging

## Security

- **Encryption**: AES-256-GCM for message content
- **Key Storage**: Android Keystore for cryptographic keys
- **Database**: SQLCipher with secure passphrase
- **Transport**: HTTPS in production
- **Authentication**: JWT with short-lived access tokens

## Development

### Running Tests

Backend:
```bash
go test ./...
```

Android:
```bash
cd android
./gradlew test
```

### Code Generation

Generate mocks (Android):
```bash
cd android
./gradlew generateMocks
```

## CI/CD

### GitHub Actions

The project uses GitHub Actions for continuous integration:

- **Backend CI**: Runs tests, builds binary, validates migrations
- **Android CI**: Runs tests, builds debug APK, builds release APK

### Downloading APKs

1. Go to [GitHub Actions](https://github.com/simon141404-gif/shadow-chat-/actions)
2. Select the latest workflow run
3. Download the `debug-apk` artifact

### Release Builds

For release builds, set up these repository secrets:

- `ANDROID_KEYSTORE_BASE64` - Base64-encoded keystore
- `ANDROID_KEYSTORE_PASSWORD` - Keystore password
- `ANDROID_KEY_ALIAS` - Key alias
- `ANDROID_KEY_PASSWORD` - Key password

Release APKs are built automatically on push to `main`/`master`/`develop` branches.

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request
