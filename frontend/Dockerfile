# Stage 1: Build the Vue.js application
FROM node:20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json and install dependencies
COPY frontend/package*.json ./
RUN npm install

# Copy the rest of the application source code
COPY frontend/ .

# Build the application for production
RUN npm run build


# Stage 2: Serve the application with Nginx
FROM nginx:stable-alpine

# Copy the built static files from the 'builder' stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy the production Nginx configuration
# We need to go up one level to access the 'deploy' directory
COPY deploy/nginx.prod.conf /etc/nginx/conf.d/default.conf

# Expose port 80 for the Nginx server
EXPOSE 80

# The default command for the nginx image is to start the server
CMD ["nginx", "-g", "daemon off;"] 