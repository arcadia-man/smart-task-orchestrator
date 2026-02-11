# Multi-stage build for React frontend
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code (excluding node_modules)
COPY public ./public
COPY src ./src
COPY index.html .
COPY vite.config.js .
COPY tsconfig.json .
COPY tsconfig.node.json .
COPY tailwind.config.js .
COPY postcss.config.js .

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built app to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]