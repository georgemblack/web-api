import React from "react";
import * as ReactDOMServer from "react-dom/server";
import Markdoc from "@markdoc/markdoc";

import Image from "./tags/Image";
import Border from "./tags/Border";
import { image, border } from "./schema";

const config = {
  tags: {
    image,
    border,
  },
};

const options = {
  components: {
    Image,
    Border,
  },
};

function generate(raw) {
  // Parse content
  const ast = Markdoc.parse(raw);
  let content = Markdoc.transform(ast, config);
  content = content.children; // Remove root <article> node added by Markdoc

  // Render full content as HTML
  const nodes = Markdoc.renderers.react(content, React, options);
  const html = ReactDOMServer.renderToStaticMarkup(nodes);

  // Render content preview
  // Ignore all content after 'border' tag is identified
  const borderIndex = content.findIndex((tag) => tag.name === "Border");
  if (borderIndex !== -1) {
    content = content.slice(0, borderIndex);
  }
  const previewNodes = Markdoc.renderers.react(content, React, options);
  const htmlPreview = ReactDOMServer.renderToStaticMarkup(previewNodes);

  return { html, htmlPreview };
}

export default generate;
