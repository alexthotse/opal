-module(jido_ffi).
-export([dispatch_jido/1]).

dispatch_jido(Action) ->
    %% Call the Elixir Falcon.JidoAgent
    'Elixir.Falcon.JidoAgent':process_action(Action).
