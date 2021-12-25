import React, { useState } from "react";
import { Blurhash } from "react-blurhash";
import { Photo } from "../api/openapi-hooks/api";

const ProgressiveImage: React.FC<{
  photo: Photo;
  maxWidth?: number;
  className?: string;
}> = ({ photo, maxWidth = 40, className = "" }) => {
  const [loaded, setLoaded] = useState(false);
  const { blur_hash, base_url, id, source } = photo;
  return (
    <div>
      {blur_hash && (!loaded || base_url === "") && (
        <Blurhash
          className={className}
          hash={blur_hash}
          resolutionX={32}
          resolutionY={32}
          punch={1}
        />
      )}
      {base_url !== "" && (
        <img
          hidden={!loaded}
          className={className}
          loading="lazy"
          onLoad={() => setLoaded(true)}
          key={id}
          // https://developers.google.com/photos/library/guides/access-media-items#image-base-urls
          src={
            base_url
            // base_url.includes("notion") ? base_url : `${base_url}=w${maxWidth}`
          }
          alt={`todo - ${source} ${id}`}
        />
      )}
    </div>
  );
};
export default ProgressiveImage;
