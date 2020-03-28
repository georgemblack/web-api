const admin = require("firebase-admin");
const uuid = require("uuid/v4");

const VIEW_COLLECTION_NAME = "web-views";
const BOOKMARK_COLLECTION_NAME = "web-bookmarks";

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
  const date = new Date();
  date.setDate(date.getDate() - 20);

  const snapshot = await db
    .collection(BOOKMARKS_COLLECTION_NAME)
    .where("timestamp", ">", date)
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

  const aggregates = {
    total: views.length,
  };

  return {
    bookmarks,
  };
}

module.exports = {
  writeView,
  getBookmarks,
};
