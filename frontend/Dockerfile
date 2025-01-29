FROM node:20.11.1-alpine3.19 AS build

WORKDIR /app

COPY package*.json .
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve the React application using Nginx
FROM nginx:latest

# Remove the default Nginx configuration file
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom Nginx configuration file
COPY nginx.conf /etc/nginx/conf.d/

# Copy the built frontend files to Nginx directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose Nginx port
EXPOSE 80

# Command to start Nginx when the container starts
CMD ["nginx", "-g", "daemon off;"]