const STOP_WORDS = new Set([
  "a", "an", "and", "are", "as", "at", "be", "by", "for", "from", "has",
  "have", "how", "i", "in", "is", "it", "of", "on", "or", "that", "the",
  "this", "to", "was", "were", "what", "when", "where", "which", "who", "why",
  "will", "with", "you", "your",
]);

export function tokenize(value) {
  return (value.toLocaleLowerCase().match(/[\p{L}\p{N}][\p{L}\p{N}'_-]*/gu) ?? [])
    .filter((token) => token.length > 1 && !STOP_WORDS.has(token));
}

export function chunkPages(pages, { chunkSize = 900, overlap = 120 } = {}) {
  const chunks = [];
  for (const page of pages) {
    const normalized = page.text.replace(/\s+/g, " ").trim();
    if (!normalized) continue;
    let start = 0;
    while (start < normalized.length) {
      let end = Math.min(start + chunkSize, normalized.length);
      if (end < normalized.length) {
        const boundary = normalized.lastIndexOf(" ", end);
        if (boundary > start + chunkSize * 0.6) end = boundary;
      }
      const text = normalized.slice(start, end).trim();
      if (text) chunks.push({ id: `${page.page}-${chunks.length}`, page: page.page, text });
      if (end >= normalized.length) break;
      start = Math.max(end - overlap, start + 1);
    }
  }
  return chunks;
}

export function rankChunks(chunks, question, limit = 4) {
  const queryTokens = [...new Set(tokenize(question))];
  if (queryTokens.length === 0) return chunks.slice(0, limit).map((chunk) => ({ ...chunk, score: 0 }));

  const documentFrequency = new Map();
  for (const chunk of chunks) {
    const seen = new Set(tokenize(chunk.text));
    for (const token of queryTokens) {
      if (seen.has(token)) documentFrequency.set(token, (documentFrequency.get(token) ?? 0) + 1);
    }
  }

  return chunks
    .map((chunk) => {
      const tokens = tokenize(chunk.text);
      const counts = new Map();
      for (const token of tokens) counts.set(token, (counts.get(token) ?? 0) + 1);
      let score = 0;
      for (const token of queryTokens) {
        const frequency = counts.get(token) ?? 0;
        if (!frequency) continue;
        const idf = Math.log(1 + chunks.length / (1 + (documentFrequency.get(token) ?? 0)));
        score += (1 + Math.log(frequency)) * idf;
      }
      const phrase = question.trim().toLocaleLowerCase();
      if (phrase.length > 4 && chunk.text.toLocaleLowerCase().includes(phrase)) score += 3;
      return { ...chunk, score };
    })
    .filter((chunk) => chunk.score > 0)
    .sort((a, b) => b.score - a.score || a.page - b.page)
    .slice(0, limit);
}

export function extractiveAnswer(question, rankedChunks) {
  if (rankedChunks.length === 0) {
    return "I couldn't find that information in the uploaded document.";
  }
  const queryTokens = new Set(tokenize(question));
  const sentences = rankedChunks.flatMap((chunk) =>
    chunk.text
      .split(/(?<=[.!?])\s+/)
      .map((text) => ({ text: text.trim(), page: chunk.page }))
      .filter(({ text }) => text.length >= 25)
  );
  const selected = sentences
    .map((sentence) => ({
      ...sentence,
      score: tokenize(sentence.text).reduce((sum, token) => sum + (queryTokens.has(token) ? 1 : 0), 0),
    }))
    .sort((a, b) => b.score - a.score)
    .slice(0, 3)
    .map(({ text }) => text);
  return selected.length ? selected.join(" ") : rankedChunks[0].text.slice(0, 700);
}
