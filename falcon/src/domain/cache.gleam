import gleam/string

pub type CacheMetrics {
  CacheMetrics(hits: Int, misses: Int, size_mb: Float)
}

pub fn cached_microcompact() -> String {
  let metrics = CacheMetrics(hits: 1405, misses: 23, size_mb: 45.5)
  let ratio = 1405.0 /. 1428.0 *. 100.0
  "Microcompact Cache stats: " <> string.inspect(metrics.size_mb) <> "MB. Hits: " <> string.inspect(metrics.hits) <> ", Misses: " <> string.inspect(metrics.misses) <> ". Hit ratio: " <> string.inspect(ratio) <> "%"
}
