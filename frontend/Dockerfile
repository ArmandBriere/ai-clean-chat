# Step 1: Build the application
FROM oven/bun AS deps

# Set the working directory in the container
WORKDIR /app

# Copy all the application files to the container
COPY package.json .
COPY bun.* .

# Install dependencies and build the application
RUN bun i

# Step 2: Build the bun application
FROM oven/bun AS builder

# Set the working directory for the builder container
WORKDIR /app

COPY --from=deps /app /app
COPY . .

RUN bun run build

# Step 3: Create a smaller image for running the application
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
CMD ["bun", "run", "dev", "--host", "0.0.0.0", "--port", "3000"]
