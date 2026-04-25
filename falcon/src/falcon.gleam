import gleam/erlang/process
import mist
import ports/connect_server
import wisp
import wisp/wisp_mist

pub fn main() {
  wisp.configure_logger()

  let assert Ok(_) =
    wisp_mist.handler(connect_server.handle_request, "secret")
    |> mist.new
    |> mist.port(8080)
    |> mist.start

  process.sleep_forever()
}
