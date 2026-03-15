# ---- Dev stage: Vite HMR ----
FROM node:22-alpine AS dev
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
CMD ["npm", "run", "dev", "--", "--host"]

# ---- Builder stage: production build ----
FROM dev AS builder
RUN npm run build

# ---- Prod stage: serve with Nginx ----
FROM nginx:alpine AS prod
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
