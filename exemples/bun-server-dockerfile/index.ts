import pg from "pg";

function renderHtmlResponse(body: string): Response {
  return new Response(`<div style="text-align: center;">${body}</div>`, {
    headers: {
      "content-type": "text/html; charset=UTF-8",
    },
  });
}

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

      return renderHtmlResponse("<h1>Postgres connected ! 🎉</h1>");
    } catch (error) {
      console.error(error);
    }
    return renderHtmlResponse("<h1>Postgres is not connected ! 😢</h1>");
  },
});

console.log("Server is running on http://localhost:" + Bun.env.PORT);
