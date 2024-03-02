import express from "express";
import pino from "pino-http";
import config from "config";

import format from "./format";
import auth from "./middlewares/auth";
import rateLimit from "./middlewares/rateLimit";
import validate from "./middlewares/validate";
import generate from "./services/content-generator/index";
import firestore from "./services/firestore";
import build from "./services/build";

const ALLOWED_ORIGIN: string = config.get("allowedOrigin");
const LIKE_COLLECTION: string = config.get("likeCollectionName");
const POST_COLLECTION: string = config.get("postCollectionName");

// Express setup
const app = express();
const logger = pino();
app.use(express.json());
app.use(logger);
const port = process.env.PORT || 9000;

/**
 * Standardized headers for all requests
 */
app.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", ALLOWED_ORIGIN);
  res.header("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE");
  res.header("Access-Control-Allow-Headers", "Content-Type, Authorization");
  next();
});

/**
 * Top-level routes
 */
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
  rateLimit.intenseRateLimit,
  auth.validateBasicAuth,
  async (req, res) => {
    return res.status(200).send({ token: auth.generateToken() });
  }
);

app.get("/likes", rateLimit.rateLimit, auth.validateToken, async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    return res.status(200).send(await firestore.getLikes());
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
});

app.post(
  "/likes",
  rateLimit.rateLimit,
  auth.validateToken,
  validate.validateLikeBody,
  async (req, res) => {
    try {
      const docPayload = format.formatLikePayload(req.body);
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
  rateLimit.rateLimit,
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

app.get("/posts", rateLimit.rateLimit, auth.validateToken, async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    let result;
    if ("published" in req.query) {
      result = await firestore.getPublishedPosts();
    } else {
      result = await firestore.getPosts();
    }

    // Transform result
    const processedPosts = result.posts.map((post) => {
      const { html, htmlPreview } = generate(post.content);
      post.contentHtml = html;
      post.contentHtmlPreview = htmlPreview;
      return post;
    });
    result.posts = processedPosts;

    return res.status(200).send(result);
  } catch (err) {
    console.log(err);
    return res.status(500).send("Internal error");
  }
});

app.get(
  "/posts/:id",
  rateLimit.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      let post = await firestore.getPost(req.params.id);

      const { html, htmlPreview } = generate(post.content);
      post.contentHtml = html;
      post.contentHtmlPreview = htmlPreview;

      return res.status(200).send(post);
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/posts",
  rateLimit.rateLimit,
  auth.validateToken,
  validate.validatePostBody,
  async (req, res) => {
    try {
      const docPayload = format.formatPostPayload(req.body);
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
  rateLimit.rateLimit,
  auth.validateToken,
  validate.validatePostBody,
  async (req, res) => {
    try {
      const docPayload = format.formatPostPayload(req.body);
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
  rateLimit.rateLimit,
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

app.post("/content", auth.validateToken, async (req, res) => {
  res.header("Content-Type", "application/json");
  try {
    const result = generate(req.body.content);
    return res.status(200).send(result);
  } catch (err) {
    console.log(err);
    return res.status(400).send("Bad request");
  }
});

app.post(
  "/builds",
  rateLimit.rateLimit,
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
  rateLimit.rateLimit,
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

app.get(
  "/hashes",
  rateLimit.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(200).send(await firestore.getHashes());
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.post(
  "/hashes",
  rateLimit.rateLimit,
  auth.validateToken,
  async (req, res) => {
    res.header("Content-Type", "application/json");
    try {
      return res.status(201).send(await firestore.postHashes(req.body));
    } catch (err) {
      console.log(err);
      return res.status(500).send("Internal error");
    }
  }
);

app.listen(port, () => console.log(`Listening on port ${port}`));
