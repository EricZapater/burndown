# 🔥 Burndown

> **Don't track what you ate. Track what you have left.**

**Burndown** is a minimalist, AI-first nutrition tracker built for efficiency. Unlike traditional apps that act as a historical log, Burndown operates on a "reverse budget" logic, treating your daily caloric and macronutrient intake as a **sprint backlog** that must be managed down to zero.

![Go](https://img.shields.io/badge/Backend-Go-00ADD8?style=flat-square&logo=go)
![React Native](https://img.shields.io/badge/Frontend-React_Native-61DAFB?style=flat-square&logo=react)
![Expo](https://img.shields.io/badge/Mobile-Expo-000020?style=flat-square&logo=expo)
![Gemini AI](https://img.shields.io/badge/AI-Gemini_1.5_Flash-8E75B2?style=flat-square&logo=google)

## 🚀 The Concept

Most nutrition apps tell you that you overate *after* it happens. Burndown focuses on **Capacity Planning**:
1.  **Set a Daily Budget:** Your "Sprint Goal" (e.g., 2500 kcal).
2.  **Snap & Burn:** Take a photo of your meal. AI analyzes it and deducts it from your remaining balance.
3.  **Consultant Mode:** Not sure what to eat? The AI analyzes your remaining "budget" (e.g., *you lack 30g of protein but only have 200kcal left*) and recommends the best option between two choices.

## 🛠 Tech Stack

This project is structured as a **Monorepo** to maintain context across the stack and facilitate AI-assisted development.

### Backend (`/api`)
* **Language:** Go (Golang) 1.22+
* **Database:** PostgreSQL
* **AI Integration:** Google Gemini 1.5 Flash (via `google-generative-ai` SDK) using **Structured Output (JSON Schema)** for deterministic responses.

### Frontend (`/app`)
* **Framework:** React Native (via Expo SDK 50+)
* **Language:** TypeScript
* **State Management:** Zustand
* **UI Library:** Gluestack-UI / Tamagui
* **Features:** Camera integration, Local Notifications (Hydration), Haptic Feedback.

## 📂 Project Structure

```bash
burndown/
├── api/                # Go Backend
│   ├── cmd/server/     # Entry point
│   ├── internal/       # Business logic & Models
│   └── db/             # SQLite database file
├── app/                # Expo Frontend
│   ├── src/components/
│   ├── src/screens/
│   └── src/stores/     # Zustand stores
├── docker-compose.yml  # Development environment
└── Makefile            # Shortcut commands
⚡️ Getting Started
Prerequisites
Go 1.22+ installed

Node.js & npm/bun

Gemini API Key (Get it from Google AI Studio)

1. Clone & Setup
Bash
git clone [https://github.com/yourusername/burndown.git](https://github.com/yourusername/burndown.git)
cd burndown
2. Configure Environment
Create a .env file in the /api directory:

Fragment de codi
PORT=8080
GEMINI_API_KEY=your_api_key_here
DB_PATH=./burndown.db
3. Run Backend (Go)
Bash
cd api
go mod download
go run cmd/server/main.go
# Server running on http://localhost:8080
4. Run Frontend (Expo)
In a new terminal:

Bash
cd app
npm install
npx expo start
# Scan the QR code with your phone (Expo Go app)
🧠 AI Workflow
Burndown uses Gemini 1.5 Flash for high-speed, low-latency image analysis.

Image Analysis: The backend sends the image byte stream to Gemini with a strict JSON Schema enforcing fields: calories, protein_g, carbs_g, fat_g.

Decision Engine: In "Consultant Mode", the backend injects the user's current database state ("remaining macros") into the prompt context to generate personalized advice.

✅ Roadmap
[ ] MVP: Core "Burndown" chart & Manual/AI logging.

[ ] Hydration: Local notifications for water intake.

[ ] Consultant Mode: "Option A vs Option B" logic.

[ ] History: Simple list of daily committed logs.

[ ] User Auth: Basic PIN or Biometric protection.

📄 License
Distributed under the MIT License.
