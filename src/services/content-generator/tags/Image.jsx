import React from "react";

function Image({ urls, alt, children }) {
  const avifUrl = urls.find((url) => url.endsWith(".avif"));
  const jpgUrl = urls.find((url) => url.endsWith(".jpg"));
  const pngUrl = urls.find((url) => url.endsWith(".png"));

  return (
    <figure>
      <picture>
        {avifUrl && (
          <source
            srcSet={"https://george.black/assets" + avifUrl}
            type="image/avif"
          />
        )}
        {jpgUrl && (
          <img src={"https://george.black/assets" + jpgUrl} alt={alt} />
        )}
        {pngUrl && (
          <img src={"https://george.black/assets" + pngUrl} alt={alt} />
        )}
      </picture>
      {children && <figcaption>{children}</figcaption>}
    </figure>
  );
}

export default Image;
