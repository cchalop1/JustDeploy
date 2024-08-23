FROM oven/bun

# Copy the lock and package file
COPY bun.lockb . 
COPY package.json . 

# Install dependencies
RUN bun install --frozen-lockfile

COPY index.ts .
EXPOSE 80

CMD ["bun", "run", "index.ts"]
