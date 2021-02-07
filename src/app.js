const express = require("express");
const config = require("config");

const auth = require("./auth");
const firestore = require("./services/firestore");
const build = require("./services/build");
const rateLimiter = require("./rateLimiter");

const ALLOWED_ORIGIN = config.get("allowedOrigin");
const VIEW_COLLECTION = config.get("viewCollectionName");
const LIKE_COLLECTION = config.get("likeCollectionName");
const POST_COLLECTION = config.get("postCollectionName");
const LINK_BIN_COLLECTION = config.get("linkBinCollectionName");

// Express setup
const app = express();
app.use(express.json());
const port = process.env.PORT || 9000;

/**
 * Standardized headers for all requests
 */
app.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", ALLOWED_ORIGIN);
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
  "/auth",
  rateLimiter.intenseRateLimit,
  auth.validateBasicAuth,
  async (req, res) => {
    return res.status(200).send({ token: auth.generateToken() });
  }
);

app.get(
  "/views",
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

app.post("/stats/views", auth.validatePrivateAccessToken, async (req, res) => {
  let document = req.body;

  // timestamp -> date object
  if (!document.timestamp) {
    return res.status(400).send("Bad request");
  }
  document.timestamp = new Date(document.timestamp);

  try {
    await firestore.postItem(VIEW_COLLECTION, document);
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
  return res.status(200).send();
});

app.delete(
  "/views/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deleteItem(VIEW_COLLECTION, req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/likes",
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
      await firestore.postItem(LIKE_COLLECTION, docPayload);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.delete(
  "/likes/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deleteItem(LIKE_COLLECTION, req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/posts",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      if ("published" in req.query) {
        return res.status(200).send(await firestore.getPublishedPosts());
      }
      return res.status(200).send(await firestore.getPosts());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/posts/:id",
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
  "/posts",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    const docPayload = {
      published: new Date(req.body.published),
      metadata: req.body.metadata,
      content: req.body.content,
    };

    try {
      await firestore.postItem(POST_COLLECTION, docPayload);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.put(
  "/posts/:id",
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
  "/posts/:id",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    try {
      await firestore.deleteItem(POST_COLLECTION, req.params.id);
      return res.status(201).send("Done");
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.get(
  "/bin/links",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getLinkBin());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/bin/links",
  rateLimiter.rateLimit,
  auth.validatePrivateAccessToken,
  async (req, res) => {
    let document = req.body;
    document.timestamp = new Date();

    if (
      !document.url ||
      typeof document.url !== "string" ||
      Object.keys(document).length != 2
    ) {
      return res.status(400).send("Validation failed");
    }

    try {
      await firestore.postItem(LINK_BIN_COLLECTION, document);
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
    return res.status(201).send("Added to bin!");
  }
);

app.post(
  "/builds",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await build.postBuild());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/backups",
  rateLimiter.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.createBackup());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.listen(port, () => console.log(`Listening on port ${port}`));
