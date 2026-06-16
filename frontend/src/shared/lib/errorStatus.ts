/**
 * Extracts an HTTP status code from a caught error value.
 *
 * The generated Hey API client with throwOnError:true throws the parsed JSON
 * error body (ErrorModel) on non-2xx responses.  ErrorModel.shape is
 * { status?: number; title?: string; detail?: string }.
 *
 * Returns the numeric status if the error looks like an ErrorModel, otherwise
 * undefined (e.g. network errors, plain Error instances, or string errors).
 */
export function getHttpErrorStatus(err: unknown): number | undefined {
	if (err && typeof err === "object" && "status" in err) {
		const status = (err as { status: unknown }).status;
		if (typeof status === "number") {
			return status;
		}
	}
	return undefined;
}
