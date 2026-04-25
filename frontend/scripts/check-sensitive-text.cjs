const fs = require("fs");
const path = require("path");

const appRoot = path.resolve(__dirname, "..");
const scanRoots = ["src", "scripts"];
const allowedTypeFile = normalizePath("src/lib/api/types.ts");
const guardFile = normalizePath("scripts/check-sensitive-text.cjs");
const smokeFile = normalizePath("scripts/smoke-admin-browser.cjs");
const allowMarker = "sensitive-text-allowlist";
const blockedTokens = [
  "payload_json",
  "capability_profile",
  "provider_account_id",
  "provider_credential",
  "provider_credentials",
  "raw_payload",
  "raw_response",
  "encrypted_payload_ref",
  "access_token",
  "api_key",
  "secret",
  "secret_version",
  "credential",
  "credentials",
  "token",
];

function main() {
  const files = scanRoots.flatMap((root) => listFiles(path.join(appRoot, root)));
  const findings = [];
  for (const file of files) {
    const relativePath = normalizePath(path.relative(appRoot, file));
    const body = fs.readFileSync(file, "utf8");
    const lines = body.split(/\r?\n/);
    lines.forEach((line, index) => {
      for (const token of blockedTokens) {
        if (!line.toLowerCase().includes(token)) continue;
        if (isAllowed(relativePath, line, token)) continue;
        findings.push({ file: relativePath, line: index + 1, token });
      }
    });
  }

  if (findings.length > 0) {
    console.error("Frontend sensitive text guard failed:");
    for (const finding of findings) {
      console.error(`- ${finding.file}:${finding.line} contains "${finding.token}"`);
    }
    console.error("Move backend-only field names to approved API types or add a narrow redaction-test allowlist marker.");
    process.exit(1);
  }

  console.log(`Frontend sensitive text guard passed: scanned ${files.length} file(s).`);
}

function listFiles(root) {
  if (!fs.existsSync(root)) return [];
  const entries = fs.readdirSync(root, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    if (entry.name === "node_modules" || entry.name === ".next") continue;
    const fullPath = path.join(root, entry.name);
    if (entry.isDirectory()) {
      files.push(...listFiles(fullPath));
      continue;
    }
    if (isScannable(fullPath)) {
      files.push(fullPath);
    }
  }
  return files;
}

function isScannable(file) {
  return [".cjs", ".js", ".jsx", ".mjs", ".ts", ".tsx"].includes(path.extname(file));
}

function isAllowed(relativePath, line, token) {
  if (relativePath === guardFile) {
    return true;
  }
  if (relativePath === allowedTypeFile && ["provider_account_id", "capability_profile"].includes(token)) {
    return true;
  }
  return relativePath === smokeFile && line.includes(allowMarker);
}

function normalizePath(value) {
  return value.split(path.sep).join("/");
}

main();
