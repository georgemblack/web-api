const fetch = require("node-fetch");

const METADATA_SERVER_TOKEN_URL =
  "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=";
const SERVICE_URL = "https://web-builder-zjxddraycq-ue.a.run.app";

async function postBuild() {
  // fetch token
  let tokenResponse = await fetch(METADATA_SERVER_TOKEN_URL + SERVICE_URL, {
    headers: {
      "Metadata-Flavor": "Google",
    },
  });
  token = await tokenResponse.text();

  // start build
  let buildResponse = await fetch(SERVICE_URL, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  console.log(token);
  console.log(buildResponse);
  return await buildResponse.json();
}

module.exports = {
  postBuild,
};
