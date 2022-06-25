const { Firestore } = require("@google-cloud/firestore");
const config = require("config");
const uuid = require("uuid");

const LIKE_COLLECTION_NAME = config.get("likeCollectionName");
const POST_COLLECTION_NAME = config.get("postCollectionName");

const COLLECTIONS_FOR_BACKUP = [LIKE_COLLECTION_NAME, POST_COLLECTION_NAME];
const BACKUP_BUCKET_NAME = config.get("backupBucketName");
const GCLOUD_PROJECT_ID = config.get("gcloudProjectID");

const firestore = new Firestore();
const admin = new Firestore.v1.FirestoreAdminClient();

async function postItem(collection, payload) {
  const doc = firestore.doc(`${collection}/${uuid.v4()}`);
  await doc.set(payload);
}

async function deleteItem(collection, id) {
  const doc = firestore.doc(`${collection}/${id}`);
  await doc.delete();
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

async function createBackup() {
  try {
    const responses = await admin.exportDocuments({
      name: admin.databasePath(GCLOUD_PROJECT_ID, "(default)"),
      outputUriPrefix: BACKUP_BUCKET_NAME,
      collectionIds: COLLECTIONS_FOR_BACKUP,
    });

    const response = responses[0];
    return {
      backupID: response["name"],
      backupPrefix: response["metadata"]["outputUriPrefix"],
    };
  } catch (err) {
    throw new Error(`Create backup failed: ${err}`);
  }
}

module.exports = {
  postItem,
  deleteItem,
  getPosts,
  getPublishedPosts,
  getPost,
  putPost,
  getLikes,
  createBackup,
};
