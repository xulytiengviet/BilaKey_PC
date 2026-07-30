"use strict";

const path = require("path");
const load = (name) => require(path.join(__dirname, name));
const BASE_ROWS = [1, 2, 3, 4, 5, 6].flatMap((n) => load(`base_${String(n).padStart(2, "0")}.json`));
const PATCH_54 = [1, 2, 3].flatMap((n) => load(`patch54_${String(n).padStart(2, "0")}.json`));
const PATCH_56 = load("patch56.json");
const INITIAL_PAIRS = load("initials.json");
const CANONICAL_PREFERRED = load("canonical.json");
const CRITICAL_COLLISIONS = load("critical.json");
const BASE_VOWEL_GROUPS = load("base_vowels.json");
const SPECIAL_CHARS = load("special_chars.json");
const DATA_CORRECTIONS = load("corrections.json");
const LEGACY_FORWARD_ALIASES = load("legacy_aliases.json");
const META = load("metadata.json");

const VERSION = META.version;
const SOURCE_VERSION = META.sourceVersion;
const RULESET = META.ruleset;
const I_TO_Y = new Map([["i", "y"], ["ì", "ỳ"], ["ỉ", "ỷ"], ["ĩ", "ỹ"], ["í", "ý"], ["ị", "ỵ"]]);
const TONE_GROUPS = ["uùủũúụ", "yỳỷỹýỵ"];
const CVN_INITIALS = [...INITIAL_PAIRS].sort((a, b) => b[0].length - a[0].length);
const CVN_CODE_INITIALS = [...new Set(INITIAL_PAIRS.map(([, code]) => code))].sort((a, b) => b.length - a.length);
const CQN_TO_CVN_INITIAL = new Map(INITIAL_PAIRS);

function getBaseVowel(ch = "") {
  const lower = ch.toLocaleLowerCase("vi-VN");
  for (const group of BASE_VOWEL_GROUPS) if (group.includes(lower)) return group[0];
  return lower;
}
function addCandidate(index, code, cqn, source) {
  const list = index.get(code) || [];
  if (!list.some((entry) => entry.cqn === cqn)) list.push(Object.freeze({cqn, source}));
  index.set(code, list);
}
function addForward(index, cqn, code, source) {
  const old = index.get(cqn);
  if (old && old.code !== code) throw new Error(`CVNSS conflict ${cqn}: ${old.code} vs ${code}`);
  if (!old) index.set(cqn, Object.freeze({code, source}));
}
function buildIndexes() {
  const cqnToCvn = new Map(), cqnToCvss = new Map(), cvnCandidates = new Map(), cvssCandidates = new Map(), cvnToCvss = new Map(), cvssToCvn = new Map();
  for (const [cqn, cvn, cvss] of BASE_ROWS) {
    addForward(cqnToCvn, cqn, cvn, "base"); addForward(cqnToCvss, cqn, cvss, "base");
    addCandidate(cvnCandidates, cvn, cqn, "base"); addCandidate(cvssCandidates, cvss, cqn, "base");
    if (!cvnToCvss.has(cvn)) cvnToCvss.set(cvn, cvss); if (!cvssToCvn.has(cvss)) cvssToCvn.set(cvss, cvn);
  }
  for (const [legacy, code] of LEGACY_FORWARD_ALIASES) { addForward(cqnToCvn, legacy, code, "legacy-alias"); addForward(cqnToCvss, legacy, code, "legacy-alias"); }
  const applyPatch = (rows, source) => { for (const [cqn, code] of rows) { addForward(cqnToCvn, cqn, code, source); addForward(cqnToCvss, cqn, code, source); addCandidate(cvnCandidates, code, cqn, source); addCandidate(cvssCandidates, code, cqn, source); if (!cvnToCvss.has(code)) cvnToCvss.set(code, code); if (!cvssToCvn.has(code)) cvssToCvn.set(code, code); } };
  applyPatch(PATCH_54, "patch54"); applyPatch(PATCH_56, "patch56");
  return {cqnToCvn, cqnToCvss, cvnCandidates, cvssCandidates, cvnToCvss, cvssToCvn};
}
const INDEX = buildIndexes();
function isAllUpper(word) { return /\p{L}/u.test(word) && word === word.toLocaleUpperCase("vi-VN"); }
function applyCase(source, value) { if (!value) return value; if (isAllUpper(source)) return value.toLocaleUpperCase("vi-VN"); const first = source[0] || ""; if (first === first.toLocaleUpperCase("vi-VN") && first !== first.toLocaleLowerCase("vi-VN")) return value[0].toLocaleUpperCase("vi-VN") + value.slice(1); return value; }
function tokenize(input) { return String(input ?? "").normalize("NFC").match(/[\p{L}]+|[^\p{L}]+/gu) || []; }
function splitCqnWord(word) { const lower = word.toLocaleLowerCase("vi-VN").normalize("NFC"); for (const [initial] of CVN_INITIALS) if (lower.startsWith(initial)) return {initial, vowel: lower.slice(initial.length)}; return {initial: "", vowel: lower}; }
function normalizeCqnParts(initial, vowel) { let onset = initial, v = vowel; if (onset === "gi" && v !== "a" && INDEX.cqnToCvn.has("i" + v)) v = "i" + v; if (onset === "qu" && getBaseVowel(v[0]) === "y") v = "u" + v; if (onset === "g" && getBaseVowel(v[0]) === "i") onset = "gi"; if (onset === "g" && "eê".includes(getBaseVowel(v[0]))) onset = "gh"; if (onset === "c" && "iêe".includes(getBaseVowel(v[0]))) onset = "k"; if (onset === "ng" && "iêe".includes(getBaseVowel(v[0]))) onset = "ngh"; return {initial: onset, vowel: v}; }
function restoreCqnInitial(code, vowel) { const base = getBaseVowel(vowel[0]); if (code === "q") return "qu"; if (code === "j") return "gi"; if (code === "z") return "d"; if (code === "d") return "đ"; if (code === "f") return "ph"; if (code === "k") return "kh"; if (code === "w") return "iêe".includes(base) ? "ngh" : "ng"; if (code === "g") return "iêe".includes(base) ? "gh" : "g"; if (code === "c") return "iêe".includes(base) ? "k" : "c"; return code; }
function migrateQuGlideTone(vowel) { const chars = Array.from(vowel); if (chars.length < 2 || getBaseVowel(chars[0]) !== "u" || getBaseVowel(chars[1]) !== "y") return vowel; const tone = TONE_GROUPS[0].indexOf(chars[0]); if (tone <= 0) return vowel; chars[0] = "u"; chars[1] = TONE_GROUPS[1][tone]; return chars.join(""); }
function cleanCqnAfterJoin(initial, vowel) { let onset = initial, v = vowel; if (!onset && I_TO_Y.has(v[0]) && !v.startsWith("ia")) v = I_TO_Y.get(v[0]) + v.slice(1); if (onset === "gi" && getBaseVowel(v[0]) === "i") { if (/^[iìỉĩíị](?:m|n|p|t|c|ch|ng|nh)?$/u.test(v)) return "g" + v; v = v.slice(1) || v; } if (onset === "qu" && getBaseVowel(v[0]) === "u" && v.length > 1) { v = migrateQuGlideTone(v); v = Array.from(v).slice(1).join(""); } return onset + v; }
function buildCqnWord(initial, vowel) { return cleanCqnAfterJoin(restoreCqnInitial(initial, vowel), vowel); }
function splitCodeWord(word, index) { const lower = word.toLocaleLowerCase("vi-VN").normalize("NFC"); if (index.has(lower)) return {initial: "", codeVowel: lower}; for (const initial of CVN_CODE_INITIALS) { if (!lower.startsWith(initial)) continue; const rest = lower.slice(initial.length); if (rest && index.has(rest)) return {initial, codeVowel: rest}; } for (const initial of CVN_CODE_INITIALS) if (lower.startsWith(initial)) return {initial, codeVowel: lower.slice(initial.length)}; return {initial: "", codeVowel: lower}; }
function chooseVowel(code, candidates) { if (!candidates.length) return code; if (candidates.length === 1) return candidates[0].cqn; const preferred = CANONICAL_PREFERRED[code]; return candidates.some((entry) => entry.cqn === preferred) ? preferred : candidates[0].cqn; }
function codeToAll(word, source) { const cvss = source === "cvss"; const candidates = cvss ? INDEX.cvssCandidates : INDEX.cvnCandidates; const toOther = cvss ? INDEX.cvssToCvn : INDEX.cvnToCvss; const {initial, codeVowel} = splitCodeWord(word, candidates); const cqn = buildCqnWord(initial, chooseVowel(codeVowel, candidates.get(codeVowel) || [])); const original = initial + codeVowel; const other = initial + (toOther.get(codeVowel) ?? codeVowel); return cvss ? {cqn: applyCase(word, cqn), cvn: applyCase(word, other), cvss: applyCase(word, original)} : {cqn: applyCase(word, cqn), cvn: applyCase(word, original), cvss: applyCase(word, other)}; }
function cqnToAll(word) { const split = splitCqnWord(word); const normalized = normalizeCqnParts(split.initial, split.vowel); const initial = CQN_TO_CVN_INITIAL.get(normalized.initial) ?? normalized.initial; const cvn = INDEX.cqnToCvn.get(normalized.vowel)?.code ?? normalized.vowel; const cvss = INDEX.cqnToCvss.get(normalized.vowel)?.code ?? cvn; return {cqn: applyCase(word, word.normalize("NFC")), cvn: applyCase(word, initial + cvn), cvss: applyCase(word, initial + cvss)}; }
function convert(input, mode = "cqn") { if (!["cqn", "cvn", "cvss"].includes(mode)) throw new Error("mode must be cqn, cvn or cvss"); const result = {cqn: [], cvn: [], cvss: []}; for (const token of tokenize(input)) { if (!/^[\p{L}]+$/u.test(token)) { result.cqn.push(token); result.cvn.push(token); result.cvss.push(token); continue; } const value = mode === "cqn" ? cqnToAll(token) : codeToAll(token, mode); result.cqn.push(value.cqn); result.cvn.push(value.cvn); result.cvss.push(value.cvss); } return {cqn: result.cqn.join(""), cvn: result.cvn.join(""), cvss: result.cvss.join("")}; }
function candidatesForWord(word, source = "cvss") { const index = source === "cvss" ? INDEX.cvssCandidates : INDEX.cvnCandidates; const {initial, codeVowel} = splitCodeWord(word, index); const words = []; for (const entry of index.get(codeVowel) || []) { const value = applyCase(word, buildCqnWord(initial, entry.cqn)); if (!words.includes(value)) words.push(value); } return words.length ? words : [applyCase(word, buildCqnWord(initial, codeVowel))]; }
function inspectWord(word, source = "cvss") { const index = source === "cvss" ? INDEX.cvssCandidates : INDEX.cvnCandidates; const split = splitCodeWord(word, index); const entries = index.get(split.codeVowel) || []; const candidates = candidatesForWord(word, source); return {input: word, source, initial: split.initial, codeVowel: split.codeVowel, ambiguous: candidates.length > 1, critical: Object.hasOwn(CRITICAL_COLLISIONS, split.codeVowel), candidates, selected: codeToAll(word, source).cqn, preferredVowel: CANONICAL_PREFERRED[split.codeVowel] ?? null, sources: entries.map((entry) => entry.source)}; }
function audit() { const groups = [...INDEX.cvssCandidates.values()].filter((entries) => entries.length > 1); return {version: VERSION, sourceVersion: SOURCE_VERSION, baseRows: BASE_ROWS.length, patch54Entries: PATCH_54.length, patch56Entries: PATCH_56.length, totalPatchEntries: PATCH_54.length + PATCH_56.length, ambiguityGroups: groups.length, explicitCanonicalPolicies: Object.keys(CANONICAL_PREFERRED).length, criticalCollisionGroups: Object.keys(CRITICAL_COLLISIONS).length, maxCandidateCount: Math.max(...groups.map((entries) => entries.length)), dataCorrections: DATA_CORRECTIONS.length, legacyForwardAliases: LEGACY_FORWARD_ALIASES.length, silentReverseOverwrite: 0}; }
function selfTest() { const tests = [["toiy", "tôi"], ["iwy", "yêu"], ["vidf", "việt"], ["ily", "yên"], ["idb", "yết"], ["tizb", "tiếng"], ["wizy", "nghiêng"], ["ses", "sẽ"], ["hed", "hề"], ["mod", "mồ"], ["mof", "mộ"], ["vos", "võ"], ["qyl", "quỳ"], ["qyz", "quỷ"], ["qys", "quỹ"], ["qyj", "quý"], ["qyr", "quỵ"]]; const failures = []; for (const [code, expected] of tests) { const actual = convert(code, "cvss").cqn; if (actual !== expected) failures.push({type: "decode", code, expected, actual}); } const roundTrips = ["sẽ", "hề", "mồ", "mộ", "võ", "quỳ", "quỷ", "quỹ", "quý", "quỵ"]; for (const cqn of roundTrips) { const encoded = convert(cqn, "cqn").cvss; const decoded = convert(encoded, "cvss").cqn; if (decoded !== cqn) failures.push({type: "roundtrip", cqn, encoded, decoded}); } const a = audit(); const expected = {baseRows: 758, totalPatchEntries: 336, ambiguityGroups: 56, explicitCanonicalPolicies: 56, criticalCollisionGroups: 5, silentReverseOverwrite: 0}; for (const [field, value] of Object.entries(expected)) if (a[field] !== value) failures.push({type: "audit", field, expected: value, actual: a[field]}); return {ok: failures.length === 0, tests: tests.length + roundTrips.length + Object.keys(expected).length, failures}; }

module.exports = Object.freeze({VERSION, SOURCE_VERSION, RULESET, convert, fromCqn: (input) => convert(input, "cqn"), fromCvn: (input) => convert(input, "cvn"), fromCvss: (input) => convert(input, "cvss"), candidatesForWord, inspectWord, audit, selfTest, corrections: DATA_CORRECTIONS, legacyForwardAliases: LEGACY_FORWARD_ALIASES, criticalCollisions: CRITICAL_COLLISIONS, specialChars: SPECIAL_CHARS});
