import { randomUUID } from "node:crypto";
import { chunkPages } from "./retrieval.js";

export class DocumentStore {
  #documents = new Map();

  constructor({ ttlMs = 60 * 60 * 1000, now = () => Date.now(), maxDocuments = 50 } = {}) {
    this.ttlMs = ttlMs;
    this.now = now;
    this.maxDocuments = maxDocuments;
  }

  put({ filename, pages }) {
    this.cleanup();
    if (this.#documents.size >= this.maxDocuments) {
      const oldest = [...this.#documents.values()].sort((a, b) => a.createdAt - b.createdAt)[0];
      if (oldest) this.#documents.delete(oldest.id);
    }
    const id = randomUUID();
    const document = {
      id,
      filename,
      pageCount: pages.length,
      chunks: chunkPages(pages),
      createdAt: this.now(),
    };
    this.#documents.set(id, document);
    return document;
  }

  get(id) {
    const document = this.#documents.get(id);
    if (!document) return null;
    if (this.now() - document.createdAt > this.ttlMs) {
      this.#documents.delete(id);
      return null;
    }
    return document;
  }

  delete(id) {
    return this.#documents.delete(id);
  }

  cleanup() {
    for (const [id, document] of this.#documents) {
      if (this.now() - document.createdAt > this.ttlMs) this.#documents.delete(id);
    }
  }
}
