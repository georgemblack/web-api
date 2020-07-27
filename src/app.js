const express = require("express");
const config = require("config");
const auth = require("./auth");
const firestore = require("./firestore");
const rateLimiter = require("./rateLimiter");

ALLOWED_ORIGINS = config.get("originWhitelist");

// Express setup
const app = express();
app.use(express.json());
const port = process.env.PORT || 9000;

/**
 * Standardized headers for all requests
 */
app.use((req, res, next) => {
  const origin = ALLOWED_ORIGINS.includes(req.headers.origin)
    ? req.headers.origin
    : "https://georgeblack.me";

  res.header("Access-Control-Allow-Origin", origin);
  res.header("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE");
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
  "/admin/auth",
  rateLimiter.intenseRateLimit,
  auth.validateBasicAuth,
  async (req, res) => {
    return res.status(200).send({ token: auth.generateToken() });
  }
);

app.get(
  "/admin/views",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getViews());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
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
    firestore.postView(docPayload);
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
  return res.status(200).send("Thanks for visiting :)");
});

app.delete(
  "/admin/views/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deleteView(req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/admin/likes",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getLikes());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/admin/likes",
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
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.delete(
  "/admin/likes/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deleteLike(req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/admin/posts",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getPosts());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/admin/posts/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getPost(req.params.id));
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/admin/posts",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    const docPayload = {
      published: new Date(req.body.published),
      metadata: req.body.metadata,
      content: req.body.content,
    };

    try {
      await firestore.postPost(docPayload);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.put(
  "/admin/posts/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    const docPayload = {
      published: new Date(req.body.published),
      metadata: req.body.metadata,
      content: req.body.content,
    };

    try {
      await firestore.putPost(req.params.id, docPayload);
      return res.status(200).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.delete(
  "/admin/posts/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deletePost(req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

/**
 * Legacy – to be removed
 */
app.get("/bookmarks", async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    return res.status(200).send(await firestore.getBookmarks());
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
});

app.listen(port, () => console.log(`Listening on port ${port}`));
