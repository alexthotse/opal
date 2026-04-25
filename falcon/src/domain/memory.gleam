import gleam/string
import gleam/list

pub type MemoryBlock {
  MemoryBlock(id: String, content: String, importance: Int)
}

pub type MemoryCache {
  MemoryCache(blocks: List(MemoryBlock))
}

pub fn extract_memories() -> String {
  let cache = MemoryCache(
    blocks: [
      MemoryBlock("M1", "User prefers dark mode", 8),
      MemoryBlock("M2", "Last project was in Rust", 5)
    ]
  )
  
  let blocks_str = cache.blocks
    |> list.map(fn(b) { "[" <> b.id <> "|Lvl:" <> string.inspect(b.importance) <> "] " <> b.content })
    |> string.join("\n")

  "EXTRACT_MEMORIES: Synthesizing context into SQLite cache:\n" <> blocks_str
}

pub fn compaction_reminders() -> String {
  let blocks_to_compress = 15
  "COMPACTION_REMINDERS: Context limit approaching, consider compressing " <> string.inspect(blocks_to_compress) <> " old memory blocks."
}
