import { createClient } from "@hey-api/openapi-ts";

await createClient({
	input: "./openapi-snapshot.json",
	output: "./src/client",
	client: "@hey-api/client-fetch",
	plugins: [{ name: "@pinia/colada" }],
});
