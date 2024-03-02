import fetch from "node-fetch";
import isEqual from "lodash.isequal";

const HOST = "http://localhost:9000";
const USERNAME = "test";
const PASSWORD = "test";

async function run() {
  // Fetch auth token
  const authResp = await fetch(`${HOST}/auth`, {
    method: "POST",
    headers: {
      Authorization: `Basic ${btoa(`${USERNAME}:${PASSWORD}`)}`,
    },
  });
  const authToken = (await authResp.json()).token;

  // Post hash list
  const hashes = {
    keys: {
      key1: random(),
    },
  };
  const hashPostResp = await fetch(`${HOST}/hashes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
    body: JSON.stringify(hashes),
  });
  if (hashPostResp.status !== 200) {
    console.error("Posting hashes failed:", await hashPostResp.text());
  }

  // Fetch hash list
  const hashResp = await fetch(`${HOST}/hashes`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  const hashesRes = await hashResp.json();

  // Test hashes
  if (!isEqual(hashes, hashesRes)) {
    console.error("Hashes response does not match:", hashes, hashesRes);
  }
}

// Generate a small, random string
function random() {
  return Math.random().toString(36).substring(7);
}

run();
