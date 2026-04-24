import gleam/string

pub type SystemStats {
  SystemStats(cpu_usage: Float, mem_usage: Float, uptime_sec: Int)
}

pub fn get_stats() -> String {
  let stats = SystemStats(cpu_usage: 45.2, mem_usage: 68.9, uptime_sec: 3600)
  "System Stats: CPU " <> string.inspect(stats.cpu_usage) <> "%, RAM " <> string.inspect(stats.mem_usage) <> "%, Uptime " <> string.inspect(stats.uptime_sec) <> "s"
}
