import pg from "pg";

Bun.serve({
  port: Bun.env.PORT,
  async fetch(req) {
    try {
      const client = new pg.Client({
        host: Bun.env.POSTGRES_HOSTNAME,
        database: Bun.env.POSTGRES_DB,
        user: Bun.env.POSTGRES_USER,
        password: Bun.env.POSTGRES_PASSWORD,
      });

      await client.connect();
      const result = await client.query("SELECT NOW()");
      console.log(result);

      return new Response(
        "🚀 Postgres is connected ! 🎉\nResult: " +
          JSON.stringify(result.fields)
      );
    } catch (error) {
      console.error(error);
    }
    return new Response("Postgres is not connected ! 😢");
  },
});

console.log("Server is running on http://localhost:" + Bun.env.PORT);
