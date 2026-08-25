import { mkdir, copyFile } from "node:fs/promises";
await mkdir("dist", { recursive: true });
await copyFile("index.html", "dist/index.html");
await copyFile("app.js", "dist/app.js");
