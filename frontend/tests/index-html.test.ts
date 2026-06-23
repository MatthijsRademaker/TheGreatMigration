import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const indexHtml = readFileSync(
	new URL("../index.html", import.meta.url),
	"utf8",
);

describe("index.html", () => {
	it("uses logo.png as browser tab icon", () => {
		expect(indexHtml).toContain('rel="icon"');
		expect(indexHtml).toContain('href="/images/logo.png"');
	});
});
