import React, { useState } from "react";
import { Blurhash } from "react-blurhash";
import { Photo } from "../api/openapi-hooks/api";

const ProgressiveImage: React.FC<{ photo: Photo; maxWidth?: number }> = ({
  photo,
  maxWidth = 40,
}) => {
  const [loaded, setLoaded] = useState(false);
  const { blur_hash, width, height, base_url, id, source } = photo;
  const scalingRatio = maxWidth / width;
  const scaledHeight = scalingRatio * height;
  return (
    <div style={{ width: maxWidth, height: scaledHeight }}>
      {blur_hash && (!loaded || base_url === "") && (
        <Blurhash
          hash={blur_hash}
          width={maxWidth}
          height={scaledHeight}
          resolutionX={32}
          resolutionY={32}
          punch={1}
        />
      )}
      {base_url !== "" && (
        <img
          loading="lazy"
          onLoad={() => setLoaded(true)}
          key={id}
          // https://developers.google.com/photos/library/guides/access-media-items#image-base-urls
          src={
            base_url
            // base_url.includes("notion") ? base_url : `${base_url}=w${maxWidth}`
          }
          width={maxWidth}
          height={scaledHeight}
          alt={`todo - ${source} ${id}`}
        />
      )}
    </div>
  );
};
export default ProgressiveImage;
