const fetch = require("node-fetch");
const config = require("config");

const { GoogleAuth } = require("google-auth-library");
const auth = new GoogleAuth();

const SERVICE_URL = config.get("buildServiceEndpoint");
let client;

async function postBuild() {
  if (!client) client = await auth.getIdTokenClient(SERVICE_URL);
  const clientHeaders = await client.getRequestHeaders();

  // start build
  let buildResponse = await fetch(SERVICE_URL, {
    method: "POST",
    headers: {
      Authorization: clientHeaders["Authorization"],
    },
  });
  return await buildResponse.json();
}

module.exports = {
  postBuild,
};
