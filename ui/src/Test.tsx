import React, { useEffect, useState } from "react";
import { ConfigData } from "./api/openapi-fetch";
import Login from "./components/Login";
import { getConfig } from "./config";

const Test: React.FC = () => {
  const [config, setConfig] = useState<ConfigData>();
  useEffect(() => {
    getConfig().then((data) => setConfig(data));
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Gourd Recipe Database
          </h2>
          <div className="mt-2 text-center text-sm text-gray-600">
            {config && config.google_client_id && <Login config={config} />}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Test;
