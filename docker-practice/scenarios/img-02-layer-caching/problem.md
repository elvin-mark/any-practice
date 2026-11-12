# Task: Optimize Dockerfile Layer Caching

In this directory, you will find a simple web service project. The current `Dockerfile` invalidates its build cache on every single source code modification because `COPY . .` is executed before dependency installation.

### Objectives:
1. Examine `Dockerfile` and `package.json`.
2. Restructure `Dockerfile` to optimize Docker's build layer cache:
   - Copy `package*.json` separately first.
   - Run a simulated dependency install (`RUN echo "Installing deps..." && sleep 1 && touch /app/deps_installed`).
   - Copy the remaining source files (`COPY . .`).
3. Build and tag the optimized image as `cached-app:latest`.
4. Verify that rebuilding the image after touching a source file (e.g., `touch server.js`) reuses the cached dependency layer rather than re-running the installation step.
