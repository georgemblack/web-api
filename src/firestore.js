const admin = require("firebase-admin");
const config = require("config");
const bowser = require("bowser");
const uuid = require("uuid/v4");

const VIEW_COLLECTION_NAME = config.get("viewCollectionName");
const LIKE_COLLECTION_NAME = config.get("likeCollectionName");
const POST_COLLECTION_NAME = config.get("postCollectionName");

// Firestore connection
admin.initializeApp({
  credential: admin.credential.applicationDefault(),
});
const db = admin.firestore();

async function getViews() {
  const date = new Date();
  date.setDate(date.getDate() - 30);

  const snapshot = await db
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

function postView(payload) {
  const docRef = db.collection(VIEW_COLLECTION_NAME).doc(uuid());
  docRef.set(payload);
}

async function deleteView(id) {
  const docRef = db.collection(VIEW_COLLECTION_NAME).doc(id);
  await docRef.delete();
}

async function getLikes() {
  const date = new Date();
  date.setDate(date.getDate() - 30);

  const snapshot = await db
    .collection(LIKE_COLLECTION_NAME)
    .where("timestamp", ">", date)
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

async function postLike(payload) {
  const docRef = db.collection(LIKE_COLLECTION_NAME).doc(uuid());
  docRef.set(payload);
}

async function deleteLike(id) {
  const docRef = db.collection(LIKE_COLLECTION_NAME).doc(id);
  await docRef.delete();
}

async function getPosts() {
  const snapshot = await db
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

async function getPost(id) {
  const doc = await db.collection(POST_COLLECTION_NAME).doc(id).get();
  const payload = doc.data();
  return {
    id: doc.id,
    ...payload,
  };
}

async function postPost(payload) {
  const docRef = db.collection(POST_COLLECTION_NAME).doc(uuid());
  await docRef.set(payload);
}

async function putPost(id, payload) {
  const docRef = db.collection(POST_COLLECTION_NAME).doc(id);
  await docRef.set(payload);
}

async function deletePost(id) {
  const docRef = db.collection(POST_COLLECTION_NAME).doc(id);
  await docRef.delete();
}

/**
 * Legacy – to be removed
 */
async function getBookmarks() {
  const snapshot = await db
    .collection(LIKE_COLLECTION_NAME)
    .orderBy("timestamp", "desc")
    .get();

  const bookmarks = snapshot.docs.map((doc) => {
    const payload = doc.data();
    return {
      id: doc.id,
      timestamp: payload.timestamp._seconds,
      title: payload.title,
      url: payload.url,
    };
  });

  return {
    bookmarks,
  };
}

module.exports = {
  getViews,
  postView,
  deleteView,
  getPosts,
  getPost,
  postPost,
  putPost,
  deletePost,
  getLikes,
  postLike,
  deleteLike,
  getBookmarks,
};
