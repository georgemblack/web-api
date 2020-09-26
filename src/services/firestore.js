const { Firestore } = require("@google-cloud/firestore");
const config = require("config");
const bowser = require("bowser");
const uuid = require("uuid");

const VIEW_COLLECTION_NAME = config.get("viewCollectionName");
const LIKE_COLLECTION_NAME = config.get("likeCollectionName");
const POST_COLLECTION_NAME = config.get("postCollectionName");

const firestore = new Firestore();

async function postItem(collection, payload) {
  const doc = firestore.doc(`${collection}/${uuid.v4()}`);
  await doc.set(payload);
}

async function deleteItem(collection, id) {
  const doc = firestore.doc(`${collection}/${id}`);
  await doc.delete();
}

async function getViews() {
  const date = new Date();
  date.setDate(date.getDate() - 30);

  const snapshot = await firestore
    .collection(VIEW_COLLECTION_NAME)
    .where("timestamp", ">", date)
    .orderBy("timestamp", "desc")
    .get();

  const views = snapshot.docs.map((doc) => {
    const payload = doc.data();
    const browser = bowser.getParser(payload.userAgent);
    const browserName = browser.getBrowserName();

    return {
      id: doc.id,
      timestamp: payload.timestamp._seconds,
      pathname: payload.pathname,
      referrer: payload.referrer || "",
      userAgent: payload.userAgent,
      windowInnerWidth: payload.windowInnerWidth,
      timezone: payload.timezone,
      hostname: payload.hostname,
      userAgent: payload.userAgent,
      browser: browserName,
    };
  });

  return {
    views,
  };
}

async function getLikes() {
  const snapshot = await firestore
    .collection(LIKE_COLLECTION_NAME)
    .orderBy("timestamp", "desc")
    .get();

  const likes = snapshot.docs.map((doc) => {
    const payload = doc.data();
    return {
      id: doc.id,
      ...payload,
    };
  });

  return {
    likes,
  };
}

async function getPosts() {
  const snapshot = await firestore
    .collection(POST_COLLECTION_NAME)
    .orderBy("published", "desc")
    .get();

  const posts = snapshot.docs.map((doc) => {
    const payload = doc.data();
    return {
      id: doc.id,
      ...payload,
    };
  });

  return {
    posts,
  };
}

async function getPublishedPosts() {
  const snapshot = await firestore
    .collection(POST_COLLECTION_NAME)
    .orderBy("published", "desc")
    .get();

  let posts = snapshot.docs.map((doc) => {
    const payload = doc.data();
    return {
      id: doc.id,
      ...payload,
    };
  });

  // filter
  posts = posts.filter((post) => {
    if (!("metadata" in post)) return false;
    if (!("draft" in post.metadata)) return false;
    return !post.metadata.draft;
  });

  return {
    posts,
  };
}

async function getPost(id) {
  const doc = await firestore.doc(`${POST_COLLECTION_NAME}/${id}`).get();
  const payload = doc.data();
  return {
    id: doc.id,
    ...payload,
  };
}

async function putPost(id, payload) {
  const doc = firestore.doc(`${POST_COLLECTION_NAME}/${id}`);
  await doc.set(payload);
}

module.exports = {
  postItem,
  deleteItem,
  getViews,
  getPosts,
  getPublishedPosts,
  getPost,
  putPost,
  getLikes,
};
