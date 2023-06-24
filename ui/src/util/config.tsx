import { Configuration } from "../api/openapi-fetch";
import { fetchGetConfig } from "../api/react-query/gourdApiComponents";
import { getJWT } from "../auth/auth";
import { getAPIURL } from "./urls";

export const getConfig = async () => {
  const config = await fetchGetConfig({});
  return config;
};

export const getOpenapiFetchConfig = () => {
  const c = new Configuration({ basePath: getAPIURL(), accessToken: getJWT() });
  return c;
};
