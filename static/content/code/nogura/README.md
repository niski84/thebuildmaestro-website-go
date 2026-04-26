# Nogura

Nogura is an ambient intelligence platform I am building in Go. It listens to speech and text streams, classifies intent at the moment data arrives, and acts on it through confidence-gated execution.

## Why Nogura Exists

Most AI systems treat retrieval as the hard part. You stuff everything into a vector database and hope the query at the other end finds what matters. That worked for early RAG demos. It does not hold up when you want software that actually behaves intelligently in the background.

The lesson I took from building earlier RAG systems is that the work belongs at the front of the pipeline, not the back. If you classify intent and extract structured meaning at ingestion time, retrieval becomes simple, and the system can act without being asked.

## Core Ideas

**Semantic action vectors.** Every event that enters the system is mapped into a vector space that represents what it means and what it might trigger. This is not just embeddings of text. It is a structured representation of intent, entities, and possible actions.

**Temporal claim tiering.** Not every fact is equal. Some are lasting, some are session-bound, some are transient. Nogura tracks the lifetime of every claim it learns and discards what should not stick.

**Rolling entity cache.** Pronouns and references resolve against a short-lived cache of recently mentioned entities. The system knows what "she" or "that file" or "the project" refers to without searching the entire history.

**Knowledge graph extraction.** Subject-relation-object triples are extracted from incoming streams and stored alongside the vector representation, giving the system both semantic similarity and structured query.

**Confidence-gated execution.** The system only acts when its confidence in the intent is high enough. Low-confidence events become observations, not actions. This is the safety boundary between an assistant that helps and one that does damage.

**Dual-pass transcription and an MCP server.** Audio capture runs through Whisper, with a refinement pass for accuracy. The MCP server lets external agents query and act on Nogura state.

## Stack

- Go (single binary)
- Local Whisper for transcription
- Local LLMs via Ollama
- Embedded vector store
- Event-driven pub/sub internally

## Status

Active development. Private repository for now. Public writeups and demos coming as the system stabilizes.
