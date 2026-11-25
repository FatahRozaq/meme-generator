// Cek cepat port dev sebelum process-compose dijalankan. Sisa proses dari
// sesi sebelumnya sering masih memegang port, membuat server baru gagal
// bind dan browser menggantung di server lama.
import net from "node:net";

const ports = [3000, 8080, 8099];
const hosts = ["127.0.0.1", "::1"];

function probe(port, host) {
  return new Promise((resolve) => {
    const sock = net.connect({ port, host });
    const done = (busy) => {
      sock.destroy();
      resolve(busy);
    };
    sock.setTimeout(400);
    sock.once("connect", () => done(true));
    sock.once("timeout", () => done(false));
    sock.once("error", () => done(false));
  });
}

const busy = [];
for (const port of ports) {
  for (const host of hosts) {
    if (await probe(port, host)) {
      busy.push(port);
      break;
    }
  }
}

if (busy.length > 0) {
  console.error(`[dev] port sudah dipakai: ${busy.join(", ")}.`);
  console.error("[dev] hentikan proses dev lama dulu: make dev-stop");
  process.exit(1);
}
