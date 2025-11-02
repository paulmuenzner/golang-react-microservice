# app/frontend/public/Dockerfile

ARG NODE_VERSION

###############################################
################# DEVELOPMENT #################
###############################################
FROM docker.io/library/node:${NODE_VERSION}-alpine AS dev

WORKDIR /build

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source
COPY . .

# Hot reload env variables
ENV WATCHPACK_POLLING=true
ENV CHOKIDAR_USEPOLLING=true

# Expose port
EXPOSE 3000

# Start dev server
CMD ["npm", "run", "dev"]

###############################################
########### PRODUCTION BUILD STAGE ############
###############################################
FROM docker.io/library/node:${NODE_VERSION}-alpine AS builder

WORKDIR /build

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

###############################################
############## PRODUCTION STAGE ###############
###############################################
FROM docker.io/library/node:${NODE_VERSION}-alpine AS production

WORKDIR /build

ENV NODE_ENV=production

# Copy built app
COPY --from=builder /build/package*.json ./
COPY --from=builder /build/.next ./.next
COPY --from=builder /build/public ./public
COPY --from=builder /build/node_modules ./node_modules

EXPOSE 3000

CMD ["npm", "start"]