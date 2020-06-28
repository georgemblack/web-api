const express = require("express");
const config = require("config");
const auth = require("./auth");
const firestore = require("./firestore");
const rateLimiter = require("./rateLimiter");

DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN = config.get("accessControlAllowOrigin");

// Express setup
const app = express();
app.use(express.json());
const port = process.env.PORT || 8080;

/**
 * Standardized headers for all requests
 */
app.use((req, res, next) => {
  if (req.hostname == "admin.georgeblack.me") {
    res.header("Access-Control-Allow-Origin", "https://admin.georgeblack.me");
  } else {
    res.header(
      "Access-Control-Allow-Origin",
      DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN
    );
  }
  res.header("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
  res.header("Access-Control-Allow-Headers", "Content-Type, Authorization");
  res.header("Accept-CH", "UA, Platform, Model, Arch, Viewport-Width, Width");
  res.header("Accept-CH-Lifetime", "2592000");
  next();
});

app.options((req, res) => {
  res.sendStatus(200);
});

app.get("/", (req, res) => {
  res.status(200).send("Howdy!");
});

/**
 * Generate token for client, auth with username and password
 */
app.post(
  "/auth",
  rateLimiter.rateLimit,
  auth.validateBasicAuth,
  async (req, res) => {
    return res.status(200).send({ token: auth.generateToken() });
  }
);

app.post("/views", async (req, res) => {
  // validate payload
  if (
    typeof req.body.hostname !== "string" ||
    req.body.hostname === "" ||
    typeof req.body.pathname !== "string" ||
    req.body.pathname === "" ||
    typeof req.body.referrer !== "string" ||
    typeof req.body.windowInnerWidth !== "number" ||
    !Number.isInteger(req.body.windowInnerWidth) ||
    typeof req.body.timezone !== "string" ||
    req.body.timezone === ""
  ) {
    return res.status(400).send("Validation failed");
  }

  // build client hints if available
  const clientHints = {};
  if (req.get("Sec-CH-UA")) clientHints.userAgent = req.get("Sec-CH-UA");
  if (req.get("Sec-CH-Platform"))
    clientHints.platform = req.get("Sec-CH-Platform");
  if (req.get("Sec-CH-Model")) clientHints.model = req.get("Sec-CH-Model");
  if (req.get("Sec-CH-Arch")) clientHints.arch = req.get("Sec-CH-Arch");
  if (req.get("Sec-CH-Viewport-Width"))
    clientHints.viewportWidth = req.get("Sec-CH-Viewport-Width");
  if (req.get("Sec-CH-Width")) clientHints.width = req.get("Sec-CH-Width");

  // primary payload
  const docPayload = {
    hostname: req.body.hostname,
    pathname: req.body.pathname,
    windowInnerWidth: req.body.windowInnerWidth,
    timezone: req.body.timezone,
    timestamp: new Date(),
  };

  // append possibly empty items
  if (req.get("User-Agent")) docPayload.userAgent = req.get("User-Agent");
  if (req.body.referrer) docPayload.referrer = req.body.referrer;
  if (Object.keys(clientHints).length !== 0)
    docPayload.clientHints = clientHints;

  // write to firestore
  try {
    firestore.writeView(docPayload);
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
  return res.status(200).send("Thanks for visiting :)");
});

app.get("/bookmarks", async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    return res.status(200).send(await firestore.getBookmarks());
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
});

app.get(
  "/likes",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getBookmarks());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/likes",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    if (
      typeof req.body.title !== "string" ||
      req.body.title === "" ||
      typeof req.body.url !== "string" ||
      req.body.url === ""
    ) {
      return res.status(400).send("Validation failed");
    }

    const docPayload = {
      title: req.body.title,
      url: req.body.url,
      timestamp: new Date(),
    };

    try {
      await firestore.postLike(docPayload);
      return res.status(200).send("Done");
    } catch (err) {
      return res.status(500).send("Internal error");
    }
  }
);

app.get("/posts", async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    return res.status(200).send(await firestore.getPosts());
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
});

app.listen(port, () => console.log(`Listening on port ${port}`));
