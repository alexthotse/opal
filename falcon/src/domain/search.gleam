import gleam/string
import gleam/list

pub type SearchResult {
  SearchResult(file: String, line: Int, match: String)
}

pub type SearchEngine {
  SearchEngine(results: List(SearchResult), time_ms: Int)
}

pub fn quick_search() -> String {
  let engine = SearchEngine(
    results: [
      SearchResult("src/main.rs", 42, "fn main() {"),
      SearchResult("src/utils.rs", 10, "pub fn helper() {")
    ],
    time_ms: 15
  )

  let results_str = engine.results 
    |> list.map(fn(r) { r.file <> ":" <> string.inspect(r.line) <> " - " <> r.match })
    |> string.join("\n")

  "Quick Search completed in " <> string.inspect(engine.time_ms) <> "ms:\n" <> results_str
}
