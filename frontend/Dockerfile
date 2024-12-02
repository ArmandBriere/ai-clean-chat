# Step 1: Build the application
FROM oven/bun AS builder

# Set the working directory in the container
WORKDIR /app

# Copy all the application files to the container
COPY . .

# Install dependencies and build the application
RUN bun i
RUN bun run build

# Step 2: Create a smaller image for running the application
FROM oven/bun

# Set the working directory for the runtime container
WORKDIR /app

# Copy the production build and required files
COPY --from=builder /app/.svelte-kit ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/bun.lockb ./bun.lockb

# Expose the port the application will run on
EXPOSE 3000

# Start the application
CMD ["bun", "run", "preview"]
