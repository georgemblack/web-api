import Firestore from "@google-cloud/firestore";

/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatPostPayload(requestBody) {
  const docPayload = {
    title: requestBody.title,
    slug: requestBody.slug,
    published: new Date(requestBody.published),
    content: requestBody.content,
    draft: requestBody.draft,
    tags: requestBody.tags || [],
  };

  // If location provided, convert to Firestore geopoint
  if ("location" in requestBody.metadata) {
    const lat = requestBody.metadata.location[0];
    const lon = requestBody.metadata.location[1];
    docPayload.location = new Firestore.GeoPoint(
      Number(lat),
      Number(lon)
    );
  }

  return docPayload;
}

/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatLikePayload(requestBody) {
  return {
    title: requestBody.title,
    url: requestBody.url,
    timestamp: new Date(),
  };
}

export default {
  formatPostPayload,
  formatLikePayload,
};
