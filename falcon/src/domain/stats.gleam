import gleam/string
import gleam/erlang/charlist.{type Charlist}

@external(erlang, "os", "cmd")
fn os_cmd(command: Charlist) -> Charlist

pub type SystemStats {
  SystemStats(cpu_usage: Float, mem_usage: Float, uptime_sec: Int)
}

pub fn get_stats() -> String {
  let uptime_out = os_cmd(charlist.from_string("cat /proc/uptime"))
  let mem_out = os_cmd(charlist.from_string("free -m | awk '/^Mem:/ {print $3\"/\"$2\" MB\"}'"))
  let cpu_out = os_cmd(charlist.from_string("top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'"))

  "System Stats: CPU " <> string.trim(charlist.to_string(cpu_out)) <> 
  "%, RAM " <> string.trim(charlist.to_string(mem_out)) <> 
  ", Uptime " <> string.trim(charlist.to_string(uptime_out))
}
