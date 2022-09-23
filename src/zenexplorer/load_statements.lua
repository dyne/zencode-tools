local function only_statements(steps, defaults)
  local statements = {}
  for k, _ in pairs(steps) do
    table.insert(statements, k)
  end
  return statements
end
local given_stms = only_statements(ZEN.given_steps)
local then_stms = only_statements(ZEN.then_steps)

local function default_statements(steps)
  local statements = {}
  for k, _ in pairs(steps) do
    statements[k] = true
  end
end

local DEFAULT_SCENARIOS = {}
for scenario in pairs(SCENARIOS) do
  local _, _, s = string.find(scenario, "zencode_(.*)")
  table.insert(DEFAULT_SCENARIOS, s)
end

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
local when_defaults = default_statements(ZEN.when_steps)

-- Delete all defaults
ZEN.when_steps = {}

local when_stms = {}
-- Load one scenario at a time
for _, scenario in ipairs(SCENARIOS) do
  load_scenario("zencode_" .. scenario)
  when_stms[scenario] = only_statements(ZEN.when_steps)
end

print(JSON.encode({
  ["given"] = given_stms,
  ["then"] = then_stms,
  ["when"] = when_stms,
  ["default_scenarios"] = DEFAULT_SCENARIOS,
}))
