import { PDFParse } from "pdf-parse";

export async function parsePdf(buffer, { maxPages = 500, maxCharacters = 2_000_000 } = {}) {
  const parser = new PDFParse({ data: new Uint8Array(buffer) });
  try {
    const result = await parser.getText();
    if (result.pages.length > maxPages) {
      const error = new Error(`The PDF exceeds the ${maxPages}-page processing limit.`);
      error.status = 422;
      throw error;
    }
    let extractedCharacters = 0;
    const pages = result.pages
      .map((page) => {
        const text = page.text.trim();
        extractedCharacters += text.length;
        if (extractedCharacters > maxCharacters) {
          const error = new Error("The PDF contains too much extractable text.");
          error.status = 422;
          throw error;
        }
        return { page: page.num, text };
      })
      .filter((page) => page.text.length > 0);
    if (pages.length === 0) {
      const error = new Error("The PDF contains no extractable text.");
      error.status = 422;
      throw error;
    }
    return pages;
  } finally {
    await parser.destroy();
  }
}
