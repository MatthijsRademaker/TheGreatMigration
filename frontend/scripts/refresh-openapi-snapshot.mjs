import { writeFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";

const snapshotPath = fileURLToPath(
	new URL("../openapi-snapshot.json", import.meta.url),
);
const openApiUrl =
	process.env.OPENAPI_URL || "http://localhost:8080/openapi.json";

const response = await fetch(openApiUrl);

if (!response.ok) {
	throw new Error(
		`Failed to fetch ${openApiUrl}: ${response.status} ${response.statusText}`,
	);
}

const snapshot = await response.text();
await writeFile(snapshotPath, `${snapshot.trim()}\n`, "utf8");
