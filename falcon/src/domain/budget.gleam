import gleam/string

pub type BudgetContext {
  BudgetContext(allocated: Float, spent: Float, currency: String)
}

pub fn check_budget() -> String {
  let ctx = BudgetContext(allocated: 1000.0, spent: 450.5, currency: "USD")
  let remaining = ctx.allocated -. ctx.spent
  "Budget Check: " <> string.inspect(ctx.spent) <> " " <> ctx.currency <> " spent out of " <> string.inspect(ctx.allocated) <> " " <> ctx.currency <> ". Remaining: " <> string.inspect(remaining) <> " " <> ctx.currency
}
