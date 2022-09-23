local function only_statements(steps, defaults)
  local statements = {}
  for k, _ in pairs(steps) do
    table.insert(statements, k)
  end
  return statements
end
local given_stms = only_statements(ZEN.given_steps)
local then_stms = only_statements(ZEN.then_steps)

local SCENARIOS = {
  "array",
  "bitcoin",
  "credential",
  "data",
  "debug",
  "dictionary",
  "dp3t",
  "ecdh",
  "eddsa",
  "ethereum",
  "given",
  "hash",
  "http",
  "keyring",
  "pack",
  "petition",
  "qp",
  "random",
  "reflow",
  "schnorr",
  "secshare",
  "then",
  "verify",
  "w3c",
  "when"
}

local when_stms = {}
when_stms["default"] = only_statements(ZEN.when_steps)
-- Load one scenario at a time
for _, scenario in ipairs(SCENARIOS) do
  ZEN.when_steps = {}
  load_scenario("zencode_" .. scenario)
  local statements = only_statements(ZEN.when_steps)
  if #statements > 0 then
    when_stms[scenario] = statements
  end
end

print(JSON.encode({
  ["given"] = given_stms,
  ["then"] = then_stms,
  ["when"] = when_stms,
}))
