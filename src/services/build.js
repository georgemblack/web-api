const fetch = require("node-fetch");

const METADATA_SERVER_TOKEN_URL =
  "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=";
const RECEIVING_SERVICE_URL = "https://web-builder-zjxddraycq-ue.a.run.app";

async function postBuild() {
  let response = await fetch(
    METADATA_SERVER_TOKEN_URL + RECEIVING_SERVICE_URL,
    {
      headers: {
        "Metadata-Flavor": "Google",
      },
    }
  );

  responseBody = await response.json();
  console.log(response);
  console.log(responseBody);
  return responseBody;
}

module.exports = {
  postBuild,
};
