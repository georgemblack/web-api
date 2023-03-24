import isEqual from "lodash.isequal";

function validatePostBody(req, res, next) {
  const body = req.body;
  if (!body) {
    return res.status(400).send("Validation failed");
  }

  const requiredBodyAttributes = ["title", "slug", "draft", "published", "content"].sort();
  const bodyAttributes = Object.keys(body).sort();

  requiredBodyAttributes.forEach((required) => {
    if (!bodyAttributes.includes(required)) {
      return res.status(400).send("Validation failed");
    }
  });

  /**
   * Location should be formatted as array with two strings, i.e. ["12.34", "56.78"]
   * Both strings should be coercible to numbers.
   */
  if ("location" in body) {
    const location = body.location;
    if (!Array.isArray(location)) {
      return res.status(400).send("Validation failed");
    }
    if (location.length != 2) {
      return res.status(400).send("Validation failed");
    }
    if (typeof location[0] != "string" || typeof location[1] != "string") {
      return res.status(400).send("Validation failed");
    }
    if (!Number(location[0]) || !Number(location[1])) {
      return res.status(400).send("Validation failed");
    }
  }

  next();
}

function validateLikeBody(req, res, next) {
  const body = req.body;
  if (!body) {
    return res.status(400).send("Validation failed");
  }

  const requiredBodyAttributes = ["title", "url"].sort();
  const bodyAttributes = Object.keys(body).sort();

  if (!isEqual(bodyAttributes, requiredBodyAttributes)) {
    return res.status(400).send("Validation failed");
  }
  if (typeof body.title !== "string") {
    return res.status(400).send("Validation failed");
  }
  if (body.title === "") {
    return res.status(400).send("Validation failed");
  }
  if (typeof body.url !== "string") {
    return res.status(400).send("Validation failed");
  }
  if (body.url === "") {
    return res.status(400).send("Validation failed");
  }

  next();
}

export default {
  validatePostBody,
  validateLikeBody,
};
