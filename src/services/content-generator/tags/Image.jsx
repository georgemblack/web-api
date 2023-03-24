import React from "react";

function Image({ urls, alt, children }) {
  const avifUrl = urls.find((url) => url.endsWith(".avif"));
  const jpgUrl = urls.find((url) => url.endsWith(".jpg"));
  const pngUrl = urls.find((url) => url.endsWith(".png"));

  return (
    <figure>
      <picture>
        {avifUrl && <source srcSet={"/assets" + avifUrl} type="image/avif" />}
        {jpgUrl && <img src={"/assets" + jpgUrl} alt={alt} />}
        {pngUrl && <img src={"/assets" + pngUrl} alt={alt} />}
      </picture>
      {children && <figcaption>{children}</figcaption>}
    </figure>
  );
}

export default Image;
