import { client } from "@/client/client.gen";

export function resolveApiBaseUrl(baseUrl = import.meta.env.VITE_API_BASE_URL) {
	return baseUrl ?? "";
}

export function configureApiClient(
	options: { baseUrl?: string; fetch?: typeof fetch } = {},
) {
	const nextConfig: {
		baseUrl: string;
		fetch?: typeof fetch;
	} = {
		baseUrl: resolveApiBaseUrl(options.baseUrl),
	};

	if (options.fetch) {
		nextConfig.fetch = options.fetch;
	}

	return client.setConfig(nextConfig);
}
