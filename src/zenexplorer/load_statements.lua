local function only_statements(steps, defaults)
  local statements = {}
  for k, _ in pairs(steps) do
    table.insert(statements, k)
  end
  return statements
end
local given_stms = only_statements(ZEN.given_steps)
local then_stms = only_statements(ZEN.then_steps)
local foreach_stms = only_statements(ZEN.foreach_steps)

local SCENARIOS = {
  "array",
  "bbs",
  "bitcoin",
  "credential",
  "data",
  "debug",
  "dictionary",
  "dp3t",
  "ecdh",
  "eddsa",
  "es256",
  "ethereum",
  "foreach",
  "fsp",
  "given",
  "hash",
  "http",
  "keyring",
  "pack",
  "petition",
  "planetmint",
  "pvss",
  "qp",
  "random",
  "reflow",
  "schnorr",
  "sd_jwt",
  "secshare",
  "table",
  "then",
  "time",
  "verify",
  "w3c",
  "when"
}

local when_stms = {}
when_stms["default"] = only_statements(ZEN.when_steps)
local if_stms = {}
if_stms["default"] = only_statements(ZEN.if_steps)

-- Load one scenario at a time
for _, scenario in ipairs(SCENARIOS) do
  ZEN.when_steps = {}
  ZEN.if_steps = {}
  load_scenario("zencode_" .. scenario)
  local statements = only_statements(ZEN.when_steps)
  if #statements > 0 then
    when_stms[scenario] = statements
  end
  local if_statements = only_statements(ZEN.if_steps)
  if #if_statements > 0 then
	if_stms[scenario] = if_statements
  end
end

print(JSON.encode({
  ["given"] = given_stms,
  ["then"] = then_stms,
  ["when"] = when_stms,
  ["if"] = if_stms,
  ["foreach"] = foreach_stms
}))
