import Firestore from "@google-cloud/firestore";
import config from "config";
import { v4 as uuidv4 } from "uuid";

const LIKE_COLLECTION_NAME = config.get("likeCollectionName");
const POST_COLLECTION_NAME = config.get("postCollectionName");

const COLLECTIONS_FOR_BACKUP = [LIKE_COLLECTION_NAME, POST_COLLECTION_NAME];
const BACKUP_BUCKET_NAME = config.get("backupBucketName");
const GCLOUD_PROJECT_ID = config.get("gcloudProjectID");

const firestoreService = new Firestore();
const admin = new Firestore.v1.FirestoreAdminClient();

async function postItem(collection, payload) {
  const doc = firestoreService.doc(`${collection}/${uuidv4()}`);
  await doc.set(payload);
}

async function deleteItem(collection, id) {
  const doc = firestoreService.doc(`${collection}/${id}`);
  await doc.delete();
}

async function getLikes() {
  const snapshot = await firestoreService
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
  const snapshot = await firestoreService
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
  const snapshot = await firestoreService
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
  const doc = await firestoreService.doc(`${POST_COLLECTION_NAME}/${id}`).get();
  const payload = doc.data();
  return {
    id: doc.id,
    ...payload,
  };
}

async function putPost(id, payload) {
  const doc = firestoreService.doc(`${POST_COLLECTION_NAME}/${id}`);
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

export default {
  postItem,
  deleteItem,
  getLikes,
  getPosts,
  getPublishedPosts,
  getPost,
  putPost,
  createBackup,
};
