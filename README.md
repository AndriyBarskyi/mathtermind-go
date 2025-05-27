# MathTermind (Go Version)

A Go-based implementation of the Mathtermind educational platform for learning mathematics and informatics.

## Overview

MathTermind is an educational platform that combines interactive learning modules with gamification elements to create an engaging learning experience.

## Features

- Interactive learning modules for mathematics and informatics
- Gamification elements including achievements and rewards
- Adaptive learning system that adjusts to user's progress
- Progress tracking and analytics
- Mathematical tools and visualizations
- Course management system

## Project Structure

```
mathtermind-go/
├── cmd/           # Backend server
│   └── server/    # Server implementation
├── internal/      # Internal packages
│   ├── core/      # Core application logic
│   ├── db/        # Database operations
│   ├── models/    # Data models
│   └── services/  # Business services
├── main.go        # Entry point
└── go.mod         # Go module definition
```

## Installation

1. Clone the repository
2. Ensure you have Go installed
3. Run `make setup` to install dependencies
4. Run `make run` to run the application
