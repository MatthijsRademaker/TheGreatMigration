// Refresh the committed snapshot from a running backend with:
//   npm run refresh:openapi-snapshot
// Regenerate the committed client artifacts with:
//   npm run generate:api
// Or run both steps with:
//   npm run regen:api

export default {
	input: "./openapi-snapshot.json",
	output: {
		path: "./src/client",
		clean: true,
	},
	client: "@hey-api/client-fetch",
	plugins: [{ name: "@pinia/colada" }],
};
