# Build stage
FROM public.ecr.aws/amazonlinux/amazonlinux:2023.6.20241121.0 AS builder

# Install Node.js
RUN yum update -y && \
    yum install -y nodejs npm git && \
    yum clean all

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Fix permissions
RUN chmod -R 755 node_modules

# Build the application
RUN npm run build

# Production stage
FROM public.ecr.aws/amazonlinux/amazonlinux:2023.6.20241121.0

# Install nginx
RUN yum update -y && \
    yum install -y nginx && \
    yum clean all && \
    rm -rf /var/cache/yum

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]