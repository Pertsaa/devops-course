const server = Bun.serve({
  port: 8081,
  async fetch(req) {
    const { pathname } = new URL(req.url);
    if (pathname === "/") {
      if (req.method !== "GET") {
        return new Response("Method Not Allowed", { status: 405 });
      }
      try {
        const hostname = await getHostname();
        const uptime = await getUptime();
        const diskInfo = await getDiskInfo();
        const processInfo = await getProcessInfo();
        return Response.json({
          hostname,
          uptime,
          disk_info: diskInfo,
          process_info: processInfo,
        });
      } catch (error) {
        console.error(error);
        return Response.json(
          { error: "failed to fetch system data" },
          { status: 500 }
        );
      }
    }
    return new Response("Not Found", { status: 404 });
  },
});

console.log(`API listening on port ${server.port}...`);

async function getHostname() {
  const { stdout } = Bun.spawn(["hostname"]);
  const output = await new Response(stdout).text();
  return output.trim().replaceAll("\n", "");
}

async function getUptime() {
  const { stdout } = Bun.spawn(["uptime"]);
  const output = await new Response(stdout).text();
  return output.trim().replaceAll("\n", "");
}

async function getDiskInfo() {
  const { stdout } = Bun.spawn(["df", "-h"]);
  const output = await new Response(stdout).text();
  return output.trim().replaceAll("\n", "");
}

async function getProcessInfo() {
  const { stdout } = Bun.spawn(["ps", "-ax"]);
  const output = await new Response(stdout).text();
  return output.trim().replaceAll("\n", "");
}
