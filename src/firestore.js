const admin = require("firebase-admin");
const uuid = require("uuid/v4");

const VIEW_COLLECTION_NAME = "web-views";
const BOOKMARK_COLLECTION_NAME = "web-bookmarks";
const POST_COLLECTION_NAME = "web-posts";

// Firestore connection
admin.initializeApp({
  credential: admin.credential.applicationDefault(),
});
const db = admin.firestore();

function writeView(payload) {
  const docRef = db.collection(VIEW_COLLECTION_NAME).doc(uuid());
  docRef.set(payload);
}

async function getBookmarks() {
  const snapshot = await db
    .collection(BOOKMARK_COLLECTION_NAME)
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

async function getPosts() {
  const snapshot = await db
    .collection(POST_COLLECTION_NAME)
    .orderBy("published", "desc")
    .get();

  const posts = snapshot.docs.map((doc) => {
    const payload = doc.data();
    return {
      id: doc.id,
      published: payload.published._seconds,
      metadata: payload.metadata,
      content: payload.content,
    };
  });

  return {
    posts,
  };
}

module.exports = {
  writeView,
  getBookmarks,
  getPosts,
};
