const isEqual = require("lodash.isequal");

function validatePostBody(req, res, next) {
  const body = req.body;
  if (!body) {
    return res.status(400).send("Validation failed");
  }

  const requiredBodyAttributes = ["published", "content", "metadata"].sort();
  const bodyAttributes = Object.keys(body).sort();

  if (!isEqual(bodyAttributes, requiredBodyAttributes)) {
    return res.status(400).send("Validation failed");
  }

  const requiredMetadataAttributes = ["slug", "title"];
  const metadataAttributes = Object.keys(body.metadata);

  requiredMetadataAttributes.forEach((required) => {
    if (!metadataAttributes.includes(required)) {
      return res.status(400).send("Validation failed");
    }
  });

  if ("location" in body.metadata) {
    const location = body.metadata.location;
    if (!Array.isArray(location)) {
      return res.status(400).send("Validation failed");
    }
    if (location.length != 2) {
      return res.status(400).send("Validation failed");
    }
    if (typeof location[0] != "number" || typeof location[1] != "number") {
      return res.status(400).send("Validation failed");
    }
  }

  next();
}

module.exports = {
  validatePostBody,
};
