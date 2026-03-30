# ---- Build stage ----
# Uses Node to compile the Angular app into static HTML/CSS/JS
FROM node:22-alpine AS builder

WORKDIR /app

# Copy dependency files first (Docker caches this layer)
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code and build for production
COPY . .
RUN npx ng build --configuration=production

# ---- Run stage ----
# Uses nginx to serve the static files
FROM nginx:alpine

# Copy the built Angular files to nginx's serve directory
COPY --from=builder /app/dist/client/browser /usr/share/nginx/html

# Custom nginx config to handle Angular routing (SPA fallback)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
