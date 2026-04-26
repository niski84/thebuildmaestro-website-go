# From RAG to Ambient Intelligence: Why I Moved the Work to the Front of the Pipeline

I built a RAG system in Go a couple of years back. The Hive. It did all the things you are supposed to do. Document ingestion, embedding, vector search, a chat interface on top. It worked. People could ask it questions and get reasonable answers most of the time.

But I never actually used it. Not really. Not the way I had imagined I would.

That bothered me. So I spent some time figuring out why.

## The Quiet Failure of Retrieval-First Systems

The promise of RAG is that you stuff everything into a vector database and the model figures out what matters at query time. The hidden assumption is that the user will ask the right question, in the right way, at the right time.

That assumption is wrong almost always.

People do not know what to ask. They do not know what is in there. They do not remember to look. The interface is fundamentally reactive, and human attention is not.

Worse, the system has no understanding of what any of the data means. It is just a pile of vectors. The model at the front sees what the search returns and tries to make sense of it on the fly. Sometimes it works. Often it does not.

## Where the Real Work Should Live

The thing I kept coming back to is what Nick B Jones talks about. You have to know what your data means before you store it. Not after. The cost of classifying intent at ingestion time is small. The cost of trying to recover meaning from a flat vector index at query time is enormous, and it gets paid every single query forever.

So I started over with a different premise. What if the system did the work upfront? What if every piece of data that came in was classified, tagged, and structured by intent at the moment it arrived?

That changes everything downstream. Retrieval becomes simple because the index already knows what is in it. The model does less work because the structure is already there. And, more importantly, the system can act without being asked, because it knows what it has.

## Nogura

Nogura is what came out of that thinking. It is an ambient intelligence platform written in Go. It listens to speech and text, classifies intent at ingestion, and acts when its confidence is high enough.

A few of the ideas:

**Semantic action vectors.** Every event entering the system is mapped not just into an embedding, but into a structured representation of intent, entities, and possible actions. The vector tells you what the data means and what it might trigger.

**Temporal claim tiering.** Some facts are lasting. Some are session-bound. Some are transient. Nogura tracks the lifetime of every claim and lets things expire when they should.

**Rolling entity cache.** A short-lived cache of recently mentioned entities. When someone says "she" or "that file" or "the project," Nogura resolves the reference without searching the whole history.

**Knowledge graph extraction.** Subject-relation-object triples extracted from incoming streams and stored alongside the vector representation. Semantic similarity plus structured query.

**Confidence-gated execution.** The system only acts when it is confident in the intent. Low confidence becomes an observation, not an action. This is the line between an assistant that helps and one that breaks things.

## Why This Matters

For a long time, AI engineering was about making the model bigger or making the prompts smarter. Both of those have diminishing returns. The interesting work now is in the data layer. What you store, how you classify it, when you act on it.

If you are building anything that needs to behave intelligently in the background instead of waiting for a user to ask the right question, the work has to live at the front of the pipeline. Nogura is one attempt at making that real. The lessons from The Hive made it possible.

If you want to talk about any of this, my email is in the footer.
