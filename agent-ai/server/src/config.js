const integer = (value, fallback, minimum = 1) => {
  const parsed = Number.parseInt(value ?? "", 10);
  return Number.isFinite(parsed) && parsed >= minimum ? parsed : fallback;
};

const boolean = (value, fallback = false) => {
  if (value == null || value === "") return fallback;
  if (value === "true") return true;
  if (value === "false") return false;
  throw new Error("Boolean configuration values must be 'true' or 'false'.");
};

export function loadConfig(env = process.env) {
  return {
    port: integer(env.PORT, 5001),
    clientOrigins: (env.CLIENT_ORIGIN ?? "http://localhost:5173")
      .split(",")
      .map((origin) => origin.trim())
      .filter(Boolean),
    trustProxy: boolean(env.TRUST_PROXY),
    maxUploadBytes: integer(env.MAX_UPLOAD_BYTES, 10 * 1024 * 1024),
    documentTtlMs: integer(env.DOCUMENT_TTL_MS, 60 * 60 * 1000),
    openAiApiKey: env.OPENAI_API_KEY?.trim() || null,
    openAiModel: env.OPENAI_MODEL?.trim() || "gpt-5",
    serpApiKey: env.SERPAPI_KEY?.trim() || null,
  };
}
