# syntax=docker/dockerfile:1

# Set the version of Node (for base image) and Bun
ARG NODE_VERSION=20.11.0
ARG BUN_VERSION=1

################################################################################
# Use the official Bun image instead of Node as base for all stages.
FROM oven/bun:${BUN_VERSION} as base

# Set the working directory for all build stages.
WORKDIR /usr/src/app

################################################################################
# Create a stage for installing production dependencies with Bun.
FROM base as deps

# Download production dependencies using Bun's package manager.
# We use --frozen-lockfile to ensure Bun uses the exact versions in bun.lockb
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=bun.lockb,target=bun.lockb \
    bun install --frozen-lockfile --production

################################################################################
# Create a stage for building the application.
FROM deps as build

# Download additional development dependencies before building the application.
# This step ensures that Bun installs development dependencies for building.
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=bun.lockb,target=bun.lockb \
    bun install --frozen-lockfile

# Copy the source files into the image.
COPY . .

# Run the build command (adjust to your actual build command if needed).
RUN bun run build

################################################################################
# Create a final stage to run the application with minimal runtime dependencies.
FROM base as final

# Install any production dependencies that are required at runtime (like serve).
RUN bun add serve

# Use production environment by default.
ENV NODE_ENV production

# Run the application as a non-root user.
USER bun

# Copy package.json so that package manager commands can be used.
COPY package.json .

# Copy the production dependencies from the deps stage and also
# the built application from the build stage into the final image.
COPY --from=deps /usr/src/app/node_modules ./node_modules
COPY --from=build /usr/src/app/dist/ ./dist/

# Expose the port that the application listens on (use 5173 if it's a Vite app or 8080 as in the original example).
EXPOSE 8080

# Run the application using the serve package.
CMD ["bun", "run", "serve", "-s", "-l", "8080", "dist"]
