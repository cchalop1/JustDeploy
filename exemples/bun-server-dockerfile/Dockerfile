FROM oven/bun

# Copy the lock and package file
COPY bun.lockb . 
COPY package.json . 

# Install dependencies
RUN bun install --frozen-lockfile

COPY index.ts .
COPY tsconfig.json .

RUN bun check

CMD ["bun", "run", "index.ts"]
