FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install --only=production

COPY . .

# build the app
RUN npm run build


# stage -> The runner
FROM node:18-alpine AS runtime

WORKDIR /app

# Copy built application from builder stage
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY package.json ./

EXPOSE 3000
