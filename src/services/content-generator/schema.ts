export const image = {
  render: "Image",
  children: ["paragraph"],
  attributes: {
    urls: {
      type: Array,
      required: true,
      errorLevel: "critical",
    },
    alt: {
      type: String,
      required: false,
    },
  },
};

export const border = {
  render: "Border",
};
