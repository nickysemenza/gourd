import { fetchGetConfig } from "../api/react-query/gourdApiComponents";

export const getConfig = async () => {
  const config = await fetchGetConfig({});
  return config;
};
