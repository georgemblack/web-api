const fetch = require("node-fetch");

const METADATA_SERVER_TOKEN_URL =
  "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=";
const SERVICE_URL = "https://web-builder-zjxddraycq-ue.a.run.app";

async function postBuild() {
  let response = await fetch(
    METADATA_SERVER_TOKEN_URL + SERVICE_URL,
    {
      headers: {
        "Metadata-Flavor": "Google",
      },
    }
  );

  responseBody = await response.text();
  console.log(response);
  console.log(responseBody);
  return responseBody;
}

module.exports = {
  postBuild,
};
