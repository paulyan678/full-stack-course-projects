# README demonstration assets

These assets support the root and project-level README product tours.

| Application | GIF provenance | Demonstrated path |
| --- | --- | --- |
| Agent AI | Live local browser capture | Upload PDF → ask a grounded question → inspect the answer and page source |
| OnlineOrder | Live local browser capture using the H2 demo profile | Sign in → browse menu → add quantity → review total → complete demo checkout |
| SocialAI | Live local browser capture using the default adapters | Generate local SVG → publish → browse collection → open the post |
| Spotify Local | Source-driven rendering, visibly labeled in every scene | Home → Midnight Drive → favorite and play Night Signals → Favorites |

Each application has a looping `*-demo.gif` and a matching `*-poster.png`. The browser-based animations were composed from screenshots of the running applications; the source images contained no API keys or persisted private data.

Spotify requires an Android runtime that was not available on the capture machine. Its walkthrough is therefore generated from the real Compose states and checked-in JSON fixtures rather than represented as emulator footage. Rebuild it from the repository root with:

```bash
python3 -m pip install -r spotify/scripts/requirements-demo.txt
python3 spotify/scripts/render_readme_demo.py
```

The renderer validates the named fixture path before writing the GIF and poster, so changes to the showcased album or song fail instead of silently producing a misleading demonstration.
