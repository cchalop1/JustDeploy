import { createClient } from "redis";

Bun.serve({
  port: 80,
  async fetch(req) {
    try {
      const client = await createClient({
        url: "redis://localhost:6379",
        password: "password",
        username: "default",
      })
        .on("error", (err) => console.log("Redis Client Error", err))
        .connect();

      await client.set("key", "connected");
      const value = await client.get("key");
      await client.disconnect();
      return new Response("Postgres is connected ! 🎉\nResult: " + value);
    } catch (error) {
      console.error(error);
    }
    return new Response("Redis is not connected ! 😢");
  },
});

console.log("Server is running on http://localhost:80");
