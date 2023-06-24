import React, { useState } from "react";
import { Blurhash } from "react-blurhash";
import { Photo } from "../../api/react-query/gourdApiSchemas";

const ProgressiveImage: React.FC<{
  photo: Photo;
  className?: string;
}> = ({ photo, className = "w-60" }) => {
  const [loaded, setLoaded] = useState(false);
  const { blur_hash, base_url, id, source } = photo;
  return (
    <div className="transition-opacity">
      {blur_hash && (!loaded || base_url === "") && (
        <Blurhash
          className={className}
          hash={blur_hash}
          width="100%"
          resolutionX={32}
          resolutionY={32}
          punch={1}
        />
      )}
      {base_url !== "" && (
        <img
          className={className}
          loading="lazy"
          onLoad={() => setLoaded(true)}
          key={id}
          // https://developers.google.com/photos/library/guides/access-media-items#image-base-urls
          src={
            base_url
            // base_url.includes("notion") ? base_url : `${base_url}=w${maxWidth}`
          }
          // width={maxWidth}
          // height={scaledHeight}
          alt={`todo - ${source} ${id}`}
        />
      )}
    </div>
  );
};
export default ProgressiveImage;
