const { Firestore } = require("@google-cloud/firestore");

/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatPostPayload(requestBody) {
  const docPayload = {
    published: new Date(requestBody.published),
    metadata: requestBody.metadata,
    content: requestBody.content,
  };

  // If location provided, convert to Firestore geopoint
  if ("location" in docPayload.metadata) {
    const lat = docPayload.metadata.location[0];
    const lon = docPayload.metadata.location[1];
    docPayload.metadata.location = new Firestore.GeoPoint(
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

/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatLinkBinPayload(requestBody) {
  return {
    url: requestBody.url,
    timestamp: new Date(),
  };
}

module.exports = {
  formatPostPayload,
  formatLikePayload,
  formatLinkBinPayload,
};
